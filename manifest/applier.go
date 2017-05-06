package manifest

// Applier defines the interface to apply a plan to the AHA system.
type Applier interface {
	// Apply performs the changes necessary to transition from src to target configuration. If the target plan
	// could be realized, it returns an error. If the plan was applied successfully, it returns null.
	Apply(src, target *Plan) error
}
