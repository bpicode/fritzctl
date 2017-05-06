package manifest

import (
	"fmt"

	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
)

// Applier defines the interface to apply a plan to the AHA system.
type Applier interface {
	// Apply performs the changes necessary to transition from src to target configuration. If the target plan
	// could be realized, it returns an error. If the plan was applied successfully, it returns null.
	Apply(src, target *Plan) error
}

// DryRunner is an Applier that only plans changes to the AHA system.
func DryRunner() Applier {
	return &dryRunner{}
}

type dryRunner struct {
}

// Apply does only log the proposed changes.
func (d *dryRunner) Apply(src, target *Plan) error {
	planner := TargetBasedPlanner(justLogSwitchState, justLogThermostat)
	actions, err := planner.Plan(src, target)
	if err != nil {
		return err
	}
	fmt.Println("\n\nThe following actions would be applied by the manifest:\n")
	for _, action := range actions {
		action.Perform(nil)
	}
	return nil
}

func justLogThermostat(before, after Thermostat) Action {
	return &justLogThermostatAction{before: before, after: after}
}

type justLogThermostatAction struct {
	before Thermostat
	after  Thermostat
}

// Perform only logs changes.
func (a *justLogThermostatAction) Perform(f fritz.HomeAutomationApi) error {
	if a.before.Temperature != a.after.Temperature {
		fmt.Printf("\t'%s'\t%.1f°C\t⟶\t%.1f°C\n", a.before.Name, a.before.Temperature, a.after.Temperature)
	}
	return nil
}

type justLogSwitchAction struct {
	before Switch
	after  Switch
}

func justLogSwitchState(before, after Switch) Action {
	return &justLogSwitchAction{before: before, after: after}
}

// Perform only logs changes.
func (a *justLogSwitchAction) Perform(f fritz.HomeAutomationApi) error {
	if a.before.State != a.after.State {
		fmt.Printf("\t'%s'\t%s\t⟶\t%s\n", a.before.Name, console.Btoc(a.before.State), console.Btoc(a.after.State))
	}
	return nil
}

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
func (a *reconfigureSwitchAction) Perform(f fritz.HomeAutomationApi) error {
	if a.before.State != a.after.State {
		if a.after.State {
			f.SwitchOn(a.before.ain)
		} else {
			f.SwitchOff(a.before.ain)
		}
	}
	return nil
}

type reconfigureThermostatAction struct {
	before Thermostat
	after  Thermostat
}

func reconfigureThermostat(before, after Thermostat) Action {
	return &reconfigureThermostatAction{before: before, after: after}
}

// Perform applies the target state to a switch by turning it on/off.
func (a *reconfigureThermostatAction) Perform(f fritz.HomeAutomationApi) error {
	if a.before.Temperature != a.after.Temperature {
		f.ApplyTemperature(a.after.Temperature, a.before.ain)
	}
	return nil
}
