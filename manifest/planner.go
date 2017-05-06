package manifest

import (
	"fmt"

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
func TargetBasedPlanner(scf switchCommandFactory, tcf thermostatCommandFactory) Planner {
	return &targetBasedPlanner{switchCommandFactory: scf, thermostatCommandFactory: tcf}
}

type targetBasedPlanner struct {
	switchCommandFactory     switchCommandFactory
	thermostatCommandFactory thermostatCommandFactory
}

type switchCommandFactory func(before, after Switch) Action

type thermostatCommandFactory func(before, after Thermostat) Action

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
		before, ok := src.switchNamed(t.Name)
		if !ok {
			return []Action{}, fmt.Errorf("unable to find device (switch): '%s'", t.Name)
		}
		switchActions = append(switchActions, d.switchCommandFactory(before, t))
	}
	return switchActions, nil
}

// PlanThermostats creates a partial execution plan (a slice of Actions) which shall be applied to the thermostats.
func (d *targetBasedPlanner) PlanThermostats(src, target *Plan) ([]Action, error) {
	var switchActions []Action
	for _, t := range target.Thermostats {
		before, ok := src.thermostatNamed(t.Name)
		if !ok {
			return []Action{}, fmt.Errorf("unable to find device (thermostat): '%s'", t.Name)
		}
		switchActions = append(switchActions, d.thermostatCommandFactory(before, t))
	}
	return switchActions, nil
}
