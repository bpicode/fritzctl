package printer

import (
	"bytes"
	"testing"

	"github.com/bpicode/fritzctl/console"
	"github.com/stretchr/testify/assert"
)

// TestPrintJSON probes the JSON sector of Print.
func TestPrintJSON(t *testing.T) {
	capt := bytes.NewBuffer(nil)
	Print(struct{ A string }{A: "text"}, capt)
	assert.Contains(t, capt.String(), "text")
	assert.Contains(t, capt.String(), "A")
}

// TestPrintTable probes the table sector of Print.
func TestPrintTable(t *testing.T) {
	capt := bytes.NewBuffer(nil)
	table := console.NewTable(console.Headers("X"))
	Print(table, capt)
	assert.Contains(t, capt.String(), "X")
	assert.Contains(t, capt.String(), "+--")
}
