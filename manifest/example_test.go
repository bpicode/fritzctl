package manifest_test

import (
	"net/url"

	"github.com/bpicode/fritzctl/fritz"
	"github.com/bpicode/fritzctl/manifest"
)

// Using an Applier, one can change the config of devices in one run.
func ExampleNewApplier() {
	h := fritz.NewHomeAuto(fritz.URL(&url.URL{
		Host:   "fritz-box.home",
		Scheme: "https",
	}))
	a := manifest.NewApplier(h)
	a.Apply(&manifest.Plan{
		Switches: []manifest.Switch{{Name: "Switch1", State: false}, {Name: "Switch2", State: false}},
	}, &manifest.Plan{
		Switches: []manifest.Switch{{Name: "Switch1", State: false}, {Name: "Switch2", State: false}},
	})
	// output:
}
