package manifest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"os"
	"encoding/xml"
	"github.com/bpicode/fritzctl/fritz"
)

// TestConvertDevicelist tests the conversion from fritz.Devivelist.
func TestConvertDevicelist(t *testing.T) {
	f, err := os.Open("../testdata/devicelist_fritzos06.83.xml")
	defer f.Close()
	assert.NoError(t, err)

	var l fritz.Devicelist
	err = xml.NewDecoder(f).Decode(&l)
	assert.NoError(t, err)

	plan := ConvertDevicelist(&l)
	assert.NotNil(t, plan)

	temperature, ok := plan.temperatureOf("HKR_1")
	assert.True(t, ok)
	assert.InDelta(t, 126.5, temperature, 0.01)
}
