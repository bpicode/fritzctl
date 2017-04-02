package completion

import (
	"testing"
	"bytes"

	"github.com/stretchr/testify/assert"
	"fmt"
)

// TestBourneAgainSimpleApp tests an app with usage 'myapp mycommand'.
func TestBourneAgainSimpleApp(t *testing.T) {
	bash := BourneAgain("myapp", []string{"mycommand"})
	buffer := new(bytes.Buffer)
	err := bash.Export(buffer)
	assert.NoError(t, err)
	exportCapture := buffer.String()
	assert.NotEmpty(t, exportCapture)
	assert.Contains(t, exportCapture, "_myapp()")
	assert.Contains(t, exportCapture, "complete -F _myapp myapp")
	fmt.Println("Exported:\n", exportCapture)
}
