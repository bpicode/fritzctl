package mock_test

import (
	"github.com/bpicode/fritzctl/mock"
)

// Start a new mock server.
func Example() {
	m := mock.New()
	m.DeviceList = "path/to/devicelist.xml"
	m.Start()
	defer m.Close()
	// output:
}
