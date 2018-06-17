package manifest

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/bpicode/fritzctl/internal/console"
)

type ahaAPIApplier struct {
	fritz aha
}

type aha interface {
	On(names ...string) error
	Off(names ...string) error
	Temp(value float64, names ...string) error
}

// NewApplier is an Applier that performs changes to the AHA system via the HTTP API.
func NewApplier(f aha) Applier {
	return &ahaAPIApplier{fritz: f}
}

// Apply does only log the proposed changes.
func (a *ahaAPIApplier) Apply(src, target *Plan) error {
	planner := TargetBasedPlanner(reconfigureSwitch, reconfigureThermostat)
	actions, err := planner.Plan(src, target)
	if err != nil {
		return err
	}
	fanOutChan, wg := a.fanOut(actions)
	fanInChan := a.fanIn(fanOutChan)
	wg.Wait()
	close(fanOutChan)
	err = <-fanInChan
	close(fanInChan)
	return err
}

func (a *ahaAPIApplier) fanIn(fanOutChan chan error) chan error {
	fanInChan := make(chan error)
	go func() {
		var errMessages []string
		for err := range fanOutChan {
			errMessages = appendToErrorMessages(errMessages, err)
		}
		if len(errMessages) > 0 {
			fanInChan <- errors.New("the following operations failed:\n" + strings.Join(errMessages, "\n"))
		} else {
			fanInChan <- nil
		}
	}()
	return fanInChan
}

func appendToErrorMessages(errMsgs []string, err error) []string {
	if err != nil {
		return append(errMsgs, err.Error())
	}
	return errMsgs
}

func (a *ahaAPIApplier) fanOut(actions []Action) (chan error, *sync.WaitGroup) {
	var wg sync.WaitGroup
	fanOutChan := make(chan error)
	for _, action := range actions {
		wg.Add(1)
		go func(ac Action) {
			defer wg.Done()
			fanOutChan <- ac.Perform(a.fritz)
		}(action)
	}
	return fanOutChan, &wg
}

type reconfigureSwitchAction struct {
	before Switch
	after  Switch
}

func reconfigureSwitch(before, after Switch) Action {
	return &reconfigureSwitchAction{before: before, after: after}
}

// Perform applies the target state to a switch by turning it on/off.
func (a *reconfigureSwitchAction) Perform(f aha) (err error) {
	if a.before.State != a.after.State {
		if a.after.State {
			err = f.On(a.before.Name)
		} else {
			err = f.Off(a.before.Name)
		}
		if err == nil {
			fmt.Printf("\t[%s]\t'%s'\t%s\t⟶\t%s\n", console.Green("OK"), a.before.Name, console.Btoc(a.before.State), console.Btoc(a.after.State))
		} else {
			fmt.Printf("\t[%s]\t'%s'\t%s\t⟶\t%s\t%s\n", console.Red("FAIL"), a.before.Name, console.Btoc(a.before.State), console.Btoc(a.after.State), err.Error())
		}
	}
	return err
}

type reconfigureThermostatAction struct {
	before Thermostat
	after  Thermostat
}

func reconfigureThermostat(before, after Thermostat) Action {
	return &reconfigureThermostatAction{before: before, after: after}
}

// Perform applies the target state to a switch by turning it on/off.
func (a *reconfigureThermostatAction) Perform(f aha) (err error) {
	if a.before.Temperature != a.after.Temperature {
		err = f.Temp(a.after.Temperature, a.before.Name)
		if err == nil {
			fmt.Printf("\t[%s]\t'%s'\t%.1f°C\t⟶\t%.1f°C\n", console.Green("OK"), a.before.Name, a.before.Temperature, a.after.Temperature)
		} else {
			fmt.Printf("\t[%s]\t'%s'\t%.1f°C\t⟶\t%.1f°C\t%s\n", console.Red("FAIL"), a.before.Name, a.before.Temperature, a.after.Temperature, err.Error())
		}
	}
	return err
}
