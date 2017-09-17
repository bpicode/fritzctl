package fritz

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bpicode/fritzctl/concurrent"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/stringutils"
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
	results := concurrent.ScatterGather(targets, genericSuccessHandler, genericErrorHandler)
	return genericResult(results)
}

func genericSuccessHandler(key, message string) concurrent.Result {
	logger.Success("Successfully processed '" + key + "'; response was: " + strings.TrimSpace(message))
	return concurrent.Result{Msg: message, Err: nil}
}

func genericErrorHandler(key, message string, err error) concurrent.Result {
	logger.Warn("Error while processing '" + key + "'; error was: " + err.Error())
	return concurrent.Result{Msg: message, Err: fmt.Errorf("error operating '%s': %s", key, err.Error())}
}

func genericResult(results []concurrent.Result) error {
	if err := truncateToOne(results); err != nil {
		return errors.New("not all operations could be completed! Nested errors are: " + err.Error())
	}
	return nil
}

func truncateToOne(results []concurrent.Result) error {
	errs := make([]error, 0, len(results))
	for _, res := range results {
		if res.Err != nil {
			errs = append(errs, res.Err)
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
			return nil, errors.New("nothing found with name '" + name + "'; choose one out of " + strings.Join(quoted, ", "))
		}
		targets[name] = workFactory(ain)
	}
	return targets, nil
}
