package fritz

import (
	"strings"

	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/stringutils"
	"github.com/pkg/errors"
)

// homeAutoConfigurator allows to reconfigure AHA systems.
type homeAutoConfigurator interface {
	on(names ...string) error
	off(names ...string) error
	toggle(names ...string) error
	temp(value float64, names ...string) error
}

// concurrentConfigurator creates a Fritz AHA API from a given base API an applies commands in parallel using the
// go concurrency programming model.
func concurrentConfigurator(homeAuto HomeAutomationAPI) homeAutoConfigurator {
	return &concurrentAhaHTTP{homeAuto: homeAuto}
}

type concurrentAhaHTTP struct {
	homeAuto HomeAutomationAPI
}

// temp sets the desired temperature of "HKR" devices.
func (aha *concurrentAhaHTTP) temp(value float64, names ...string) error {
	return aha.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return aha.homeAuto.ApplyTemperature(value, ain)
		}
	}, names...)
}

// on switches devices on. The devices are identified by their names.
func (aha *concurrentAhaHTTP) on(names ...string) error {
	return aha.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return aha.homeAuto.SwitchOn(ain)
		}
	}, names...)
}

// off switches devices off. The devices are identified by their names.
func (aha *concurrentAhaHTTP) off(names ...string) error {
	return aha.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return aha.homeAuto.SwitchOff(ain)
		}
	}, names...)
}

// toggle toggles the on/off state of devices.
func (aha *concurrentAhaHTTP) toggle(names ...string) error {
	return aha.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return aha.homeAuto.Toggle(ain)
		}
	}, names...)
}

func (aha *concurrentAhaHTTP) doConcurrently(workFactory func(string) func() (string, error), names ...string) error {
	targets, err := buildBacklog(aha.homeAuto, names, workFactory)
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
		return errors.Wrap(err, "not all operations could be completed")
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
		return errors.New(strings.Join(msgs, "; "))
	}
	return nil
}

func buildBacklog(aha HomeAutomationAPI, names []string, workFactory func(string) func() (string, error)) (map[string]func() (string, error), error) {
	namesAndAins, err := aha.NameToAinTable()
	if err != nil {
		return nil, err
	}
	targets := make(map[string]func() (string, error))
	for _, name := range names {
		ain, ok := namesAndAins[name]
		if ain == "" || !ok {
			quoted := stringutils.Quote(stringutils.StringKeys(namesAndAins))
			return nil, errors.Errorf("nothing found with name '%s'; choose one out of '%s'", name, strings.Join(quoted, ", "))
		}
		targets[name] = workFactory(ain)
	}
	return targets, nil
}
