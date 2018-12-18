package fritz

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestButton_LastPressed test the LastPressed method of Button.
func TestButton_LastPressed(t *testing.T) {
	genesis := time.Unix(0, 0)
	eveningOfDec18 := time.Unix(1545160121, 0)
	tests := []struct {
		name string
		btn  Button
		want *time.Time
	}{
		{name: "no data", btn: Button{}, want: nil},
		{name: "nonsense data", btn: Button{LastPressedTimestamp: "lol"}, want: nil},
		{name: "genesis", btn: Button{LastPressedTimestamp: "0"}, want: &genesis},
		{name: "evening of December 18th 2018", btn: Button{LastPressedTimestamp: "1545160121"}, want: &eveningOfDec18},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			last := tc.btn.LastPressed()
			assert.Equal(t, tc.want, last)
		})
	}
}

// TestButton_FmtLastPressedCompact test the FmtLastPressedCompact method of Button.
func TestButton_FmtLastPressedCompact(t *testing.T) {
	genesis := time.Unix(0, 0)
	oneDayAfterGenesis := genesis.Add(time.Hour * 24)
	assert.Empty(t, (&Button{}).FmtLastPressedCompact(genesis))
	assert.NotEmpty(t, (&Button{LastPressedTimestamp: "0"}).FmtLastPressedCompact(genesis))
	assert.NotEmpty(t, (&Button{LastPressedTimestamp: "0"}).FmtLastPressedCompact(oneDayAfterGenesis))
	assert.NotEmpty(t, (&Button{LastPressedTimestamp: "1545160121"}).FmtLastPressedCompact(genesis))
}
