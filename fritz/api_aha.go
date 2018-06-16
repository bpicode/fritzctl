package fritz

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bpicode/fritzctl/internal/errors"
	"github.com/bpicode/fritzctl/internal/stringutils"
	"github.com/bpicode/fritzctl/logger"
)

// HomeAuto is a client for the Home Automation HTTP Interface,
// see https://avm.de/fileadmin/user_upload/Global/Service/Schnittstellen/AHA-HTTP-Interface.pdf.
type HomeAuto interface {
	Login() error
	List() (*Devicelist, error)
	On(names ...string) error
	Off(names ...string) error
	Toggle(names ...string) error
	Temp(value float64, names ...string) error
}

// NewHomeAuto a HomeAuto that communicates with the FRITZ!Box by means of the Home Automation HTTP Interface.
func NewHomeAuto(options ...Option) HomeAuto {
	client := defaultClient()
	aha := newAinBased(client)
	homeAuto := homeAuto{
		client:  client,
		aha:     aha,
		caching: false,
	}
	for _, option := range options {
		option(&homeAuto)
	}
	return &homeAuto
}

// Option applies fine-grained configuration to the HomeAuto client.
type Option func(h *homeAuto)

// codebeat:disable[TOO_MANY_IVARS]
type homeAuto struct {
	client        *Client
	aha           ainBased
	caching       bool
	cacheLock     sync.Mutex
	cachedDevices *Devicelist
}

// codebeat:enable[TOO_MANY_IVARS]

// Login tries to authenticate against the FRITZ!Box. If not successful, an error is returned. This method should be
// called before any of the other methods unless authentication is turned off at the FRITZ!Box itself.
func (h *homeAuto) Login() error {
	return h.client.Login()
}

// List fetches the devices known at the FRITZ!Box. See Devicelist for details. If the devices could not be obtained,
// an error is returned.
func (h *homeAuto) List() (*Devicelist, error) {
	h.cacheLock.Lock()
	defer h.cacheLock.Unlock()
	if h.caching && h.cachedDevices != nil {
		logger.Debug("Device list cache hit")
		l := *h.cachedDevices
		return &l, nil
	}
	l, err := h.aha.listDevices()
	if h.caching {
		h.cachedDevices = l
	}
	return l, err
}

// On activates the given devices. Devices are identified by their name. If any of the operations does not succeed,
// an error is returned.
func (h *homeAuto) On(names ...string) error {
	return h.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return h.aha.switchOn(ain)
		}
	}, names...)
}

// Off deactivates the given devices. Devices are identified by their name. Inverse of On.
func (h *homeAuto) Off(names ...string) error {
	return h.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return h.aha.switchOff(ain)
		}
	}, names...)
}

// toggle switches the state of the given devices from ON to OFF and vice versa. Devices are identified by their name.
func (h *homeAuto) Toggle(names ...string) error {
	return h.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return h.aha.toggle(ain)
		}
	}, names...)
}

// Temp applies the temperature setting to the given devices. Devices are identified by their name.
func (h *homeAuto) Temp(value float64, names ...string) error {
	return h.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return h.aha.applyTemperature(value, ain)
		}
	}, names...)
}

func (h *homeAuto) doConcurrently(workFactory func(string) func() (string, error), names ...string) error {
	targets, err := buildBacklog(h, names, workFactory)
	if err != nil {
		return err
	}
	results := scatterGather(targets, genericSuccessHandler, genericErrorHandler)
	return genericResult(results)
}

func genericSuccessHandler(key, message string) result {
	logger.Success("Successfully processed '" + key + "'; response was: " + strings.TrimSpace(message))
	return result{msg: message, err: nil}
}

func genericErrorHandler(key, message string, err error) result {
	logger.Warn("Error while processing '" + key + "'; error was: " + err.Error())
	return result{msg: message, err: errors.Wrapf(err, "error operating '%s'", key)}
}

func genericResult(results []result) error {
	if err := truncateToOne(results); err != nil {
		return errors.Wrapf(err, "not all operations could be completed")
	}
	return nil
}

func truncateToOne(results []result) error {
	errs := make([]error, 0, len(results))
	for _, res := range results {
		if res.err != nil {
			errs = append(errs, res.err)
		}
	}
	if len(errs) > 0 {
		msgs := stringutils.ErrorMessages(errs)
		return fmt.Errorf(strings.Join(msgs, "; "))
	}
	return nil
}

func buildBacklog(h HomeAuto, names []string, workFactory func(string) func() (string, error)) (map[string]func() (string, error), error) {
	devList, err := h.List()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to list devices")
	}
	namesAndAins := devList.NamesAndAins()
	targets := make(map[string]func() (string, error))
	for _, name := range names {
		ain, ok := namesAndAins[name]
		if ain == "" || !ok {
			quoted := stringutils.Quote(stringutils.Keys(namesAndAins))
			return nil, fmt.Errorf("nothing found with name '%s'; choose one out of '%s'", name, strings.Join(quoted, ", "))
		}
		targets[name] = workFactory(ain)
	}
	return targets, nil
}
