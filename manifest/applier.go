package manifest

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

func (d *dryRunner) Apply(src, target *Plan) error {
	planner := DifferentialPlanner()
	_, err := planner.Plan(src, target)
	if err != nil {
		return err
	}
	return nil
}
