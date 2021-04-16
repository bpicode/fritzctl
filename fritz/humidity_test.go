package fritz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormattingOfRelativeHumidity(t *testing.T) {
	assert.Equal(t, "56", (&Humidity{RelHumidity: "56"}).FmtRelativeHumidity())
}
