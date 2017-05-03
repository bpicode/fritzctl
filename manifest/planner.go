package manifest

import (
	"fmt"
)

type Planner interface {
	Plan(src, target *Plan) ([]Action, error)
}

type Action struct {
}

func DifferentialPlanner() Planner {
	return &diffPlanner{}
}

type diffPlanner struct {
}

func (d *diffPlanner) Plan(src, target *Plan) ([]Action, error) {
	var actions []Action
	switchActions, err := d.PlanSwitches(src, target)
	if err != nil {
		return []Action{}, err
	}
	actions = append(switchActions, actions...)
	return actions, nil
}
func (d *diffPlanner) PlanSwitches(src, target *Plan) ([]Action, error) {
	var switchActions []Action
	for _, t := range target.Switches {
		_, ok := src.switchStateOf(t.Name)
		if !ok {
			return []Action{}, fmt.Errorf("unable to find device: '%s'", t.Name)
		}
		switchActions = append(switchActions, Action{})
	}
	return switchActions, nil
}
