package manifest

import (
	"fmt"

	"github.com/bpicode/fritzctl/console"
	"github.com/bpicode/fritzctl/fritz"
)

// Planner represents an execution planner, returning actions to transition from a src to a target state.
type Planner interface {
	Plan(src, target *Plan) ([]Action, error)
}

// Action is one operation on the home automation system.
type Action interface {
	Perform(f fritz.HomeAutomationApi) error
}

// TargetBasedPlanner creates a Planner that only focuses on target state. Devices in the source state that are not
// referenced in the target will be left untouched.
func TargetBasedPlanner() Planner {
	return &targetBasedPlanner{}
}

type targetBasedPlanner struct {
}

// Plan creates an execution plan (a slice of Actions) which shall be applied in oder to reach the target state.
func (d *targetBasedPlanner) Plan(src, target *Plan) ([]Action, error) {
	var actions []Action
	switchActions, err := d.PlanSwitches(src, target)
	if err != nil {
		return []Action{}, err
	}
	actions = append(actions, switchActions...)

	thermostatActions, err := d.PlanThermostats(src, target)
	if err != nil {
		return []Action{}, err
	}
	actions = append(actions, thermostatActions...)
	return actions, nil
}

// PlanSwitches creates a partial execution plan (a slice of Actions) which shall be applied to the switches.
func (d *targetBasedPlanner) PlanSwitches(src, target *Plan) ([]Action, error) {
	var switchActions []Action
	for _, t := range target.Switches {
		before, ok := src.switchStateOf(t.Name)
		if !ok {
			return []Action{}, fmt.Errorf("unable to find device: '%s'", t.Name)
		}
		switchActions = append(switchActions, justLogSwitchState(t.Name, before, t.State))
	}
	return switchActions, nil
}

// PlanThermostats creates a partial execution plan (a slice of Actions) which shall be applied to the thermostats.
func (d *targetBasedPlanner) PlanThermostats(src, target *Plan) ([]Action, error) {
	var switchActions []Action
	for _, t := range target.Thermostats {
		before, ok := src.temperatureOf(t.Name)
		if !ok {
			return []Action{}, fmt.Errorf("unable to find device: '%s'", t.Name)
		}
		switchActions = append(switchActions, justLogThermostat(t.Name, before, t.Temperature))
	}
	return switchActions, nil
}

type justLogSwitchAction struct {
	name   string
	before bool
	after  bool
}

func justLogSwitchState(name string, before, after bool) Action {
	return &justLogSwitchAction{name: name, before: before, after: after}
}

// Perform only logs changes.
func (a *justLogSwitchAction) Perform(f fritz.HomeAutomationApi) error {
	if a.before != a.after {
		fmt.Printf("\t'%s'\t%s\t⟶\t%s\n", a.name, console.Btoc(a.before), console.Btoc(a.after))
	}
	return nil
}

type justLogThermostatAction struct {
	name   string
	before float64
	after  float64
}

func justLogThermostat(name string, before, after float64) Action {
	return &justLogThermostatAction{name: name, before: before, after: after}
}

// Perform only logs changes.
func (a *justLogThermostatAction) Perform(f fritz.HomeAutomationApi) error {
	if a.before != a.after {
		fmt.Printf("\t'%s'\t%.1f°C\t⟶\t%.1f°C\n", a.name, a.before, a.after)
	}
	return nil
}
