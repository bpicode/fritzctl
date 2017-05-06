package manifest

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestExportToStdout tests the export of a plan to stdout.
func TestExportToStdout(t *testing.T) {
	exporter := ExporterTo(os.Stdout)
	assert.NotNil(t, exporter)
	plan := Plan{
		Switches:    []Switch{{Name: "s1", State: true}, {Name: "s2", State: false}},
		Thermostats: []Thermostat{{Name: "t1", Temperature: 20.0}, {Name: "t2", Temperature: 22.0}},
	}
	err := exporter.Export(&plan)
	assert.NoError(t, err)
}

// TestExportWithMarshalError tests the export of a plan when a marshalling error occurs.
func TestExportWithMarshalError(t *testing.T) {
	exporter := ExporterTo(os.Stdout)
	err := exporter.export(nil, func(in interface{}) ([]byte, error) {
		return nil, errors.New("cannot marshall")
	})
	assert.Error(t, err)
}
