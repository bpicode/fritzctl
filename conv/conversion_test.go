package conv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFleat64ToString(t *testing.T) {
	fs := []float64{1.2, -12, 4.14, 9.72, 6.666666}
	transformable := Float64Slice(fs)
	strs := transformable.String('f', 2)
	assert.NotNil(t, strs)
	assert.Len(t, strs, len(fs))
	assert.Equal(t, "1.20", strs[0])
	assert.Equal(t, "-12.00", strs[1])
	assert.Equal(t, "4.14", strs[2])
	assert.Equal(t, "9.72", strs[3])
	assert.Equal(t, "6.67", strs[4])
}
