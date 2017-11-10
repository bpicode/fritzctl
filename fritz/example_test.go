package fritz_test

import (
	"net/url"

	"github.com/bpicode/fritzctl/fritz"
)

// Interact with Home Automation API.
func ExampleHomeAuto() {
	h := fritz.NewHomeAuto(
		fritz.URL(&url.URL{Host: "fritz-box.home", Scheme: "https"}),
		fritz.SkipTLSVerify(),
		fritz.Credentials("", "password"))
	h.Login()
	h.Temp(24.0, "Th1", "Th2")
	// output:
}
