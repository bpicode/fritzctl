package fritz

import (
	"encoding/json"
	"os"
	"testing"

	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestJSONUnmarshalTrafficData(t *testing.T) {
	var dd []TrafficMonitoringData
	f, err := os.Open("testdata/traffic_mon_answer.json")
	assert.NoError(t, err)
	defer f.Close()
	err = json.NewDecoder(f).Decode(&dd)
	assert.NoError(t, err)
	assert.NotNil(t, dd)
	assert.Len(t, dd, 1)
	d := dd[0]
	assert.NotNil(t, d)
	fmt.Printf("%+v\n", d)
}
