package fritz

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bpicode/fritzctl/concurrent"
	"github.com/bpicode/fritzctl/logger"
	"github.com/bpicode/fritzctl/stringutils"
)

// ConcurrentHomeAutomationAPI allows to concurrently reconfigure AHA systems.
type ConcurrentHomeAutomationAPI interface {
	SwitchOn(names ...string) error
	SwitchOff(names ...string) error
	Toggle(names ...string) error
	ApplyTemperature(value float64, names ...string) error
}

// ConcurrentHomeAutomation creates a Fritz AHA API from a given base API an applies commands in parallel using the
// go concurrency programming model.
func ConcurrentHomeAutomation(homeAuto HomeAutomationAPI) ConcurrentHomeAutomationAPI {
	return &concurrentAhaHTTP{homeAuto: homeAuto}
}

type concurrentAhaHTTP struct {
	homeAuto HomeAutomationAPI
}

// ApplyTemperature sets the desired temperature of "HKR" devices.
func (aha *concurrentAhaHTTP) ApplyTemperature(value float64, names ...string) error {
	return aha.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return aha.homeAuto.ApplyTemperature(value, ain)
		}
	}, names...)
}

// SwitchOn switches devices on. The devices are identified by their names.
func (aha *concurrentAhaHTTP) SwitchOn(names ...string) error {
	return aha.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return aha.homeAuto.SwitchOn(ain)
		}
	}, names...)
}

// SwitchOff switches devices off. The devices are identified by their names.
func (aha *concurrentAhaHTTP) SwitchOff(names ...string) error {
	return aha.doConcurrently(func(ain string) func() (string, error) {
		return func() (string, error) {
			return aha.homeAuto.SwitchOff(ain)
		}
	}, names...)
}

// Toggle toggles the on/off state of devices.
func (aha *concurrentAhaHTTP) Toggle(names ...string) error {
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
	logger.Success("Successfully processed device '" + key + "'; response was: " + strings.TrimSpace(message))
	return concurrent.Result{Msg: message, Err: nil}
}

func genericErrorHandler(key, message string, err error) concurrent.Result {
	logger.Warn("Error while processing device '" + key + "'; error was: " + err.Error())
	return concurrent.Result{Msg: message, Err: fmt.Errorf("error toggling device '%s': %s", key, err.Error())}
}

func genericResult(results []concurrent.Result) error {
	if err := truncateToOne(results); err != nil {
		return errors.New("Not all devices could be processed! Nested errors are: " + err.Error())
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
			return nil, errors.New("No device found with name '" + name + "'. Available devices are " + strings.Join(quoted, ", "))
		}
		targets[name] = workFactory(ain)
	}
	return targets, nil
}
