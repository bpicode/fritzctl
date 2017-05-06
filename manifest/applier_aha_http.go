package manifest

import (
	"fmt"

	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
)

// AhaApiApplier is an Applier that performs changes to the AHA system via the HTTP API.
func AhaApiApplier(f fritz.HomeAutomationApi) Applier {
	return &ahaApiApplier{fritz: f}
}

type ahaApiApplier struct {
	fritz fritz.HomeAutomationApi
}

// Apply does only log the proposed changes.
func (a *ahaApiApplier) Apply(src, target *Plan) error {
	planner := TargetBasedPlanner(reconfigureSwitch, reconfigureThermostat)
	actions, err := planner.Plan(src, target)
	if err != nil {
		return err
	}
	for _, action := range actions {
		action.Perform(a.fritz)
	}
	return nil
}

type reconfigureSwitchAction struct {
	before Switch
	after  Switch
}

func reconfigureSwitch(before, after Switch) Action {
	return &reconfigureSwitchAction{before: before, after: after}
}

// Perform applies the target state to a switch by turning it on/off.
func (a *reconfigureSwitchAction) Perform(f fritz.HomeAutomationApi) (err error) {
	if a.before.State != a.after.State {
		if a.after.State {
			_, err = f.SwitchOn(a.before.ain)
		} else {
			_, err = f.SwitchOff(a.before.ain)
		}
		if err == nil {
			fmt.Printf("\tOK\t'%s'\t%s\t⟶\t%s\n", a.before.Name, console.Btoc(a.before.State), console.Btoc(a.after.State))
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
func (a *reconfigureThermostatAction) Perform(f fritz.HomeAutomationApi) (err error) {
	if a.before.Temperature != a.after.Temperature {
		_, err = f.ApplyTemperature(a.after.Temperature, a.before.ain)
		if err == nil {
			fmt.Printf("\tOK\t'%s'\t%.1f°C\t⟶\t%.1f°C\n", a.before.Name, a.before.Temperature, a.after.Temperature)
		}
		return err
	}
	return err
}
