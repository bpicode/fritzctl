package manifest

import (
	"fmt"

	"github.com/bpicode/fritzctl/internal/console"
)

type dryRunner struct {
}

type justLogSwitchAction struct {
	before Switch
	after  Switch
}

type justLogThermostatAction struct {
	before Thermostat
	after  Thermostat
}

// DryRunner is an Applier that only plans changes to the AHA system.
func DryRunner() Applier {
	return &dryRunner{}
}

// Apply does only log the proposed changes.
func (d *dryRunner) Apply(src, target *Plan) error {
	planner := TargetBasedPlanner(justLogSwitchState, justLogThermostat)
	actions, err := planner.Plan(src, target)
	if err != nil {
		return err
	}
	fmt.Println("\n\nThe following actions would be applied by the manifest:")
	for _, action := range actions {
		action.Perform(nil)
	}
	return nil
}

func justLogThermostat(before, after Thermostat) Action {
	return &justLogThermostatAction{before: before, after: after}
}

// Perform only logs changes.
func (a *justLogThermostatAction) Perform(f aha) error {
	if a.before.Temperature != a.after.Temperature {
		fmt.Printf("\t'%s'\t%.1f°C\t⟶\t%.1f°C\n", a.before.Name, a.before.Temperature, a.after.Temperature)
	}
	return nil
}

func justLogSwitchState(before, after Switch) Action {
	return &justLogSwitchAction{before: before, after: after}
}

// Perform only logs changes.
func (a *justLogSwitchAction) Perform(f aha) error {
	if a.before.State != a.after.State {
		fmt.Printf("\t'%s'\t%s\t⟶\t%s\n", a.before.Name, console.Btoc(a.before.State), console.Btoc(a.after.State))
	}
	return nil
}
