package completion

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Contains(t, exportCapture, "mycommand")
	fmt.Println("Exported:\n", exportCapture)
}

// TestBourneAgainSimpleAppWithTwoArgs tests an app with usage 'myapp {mycommand|anothercomand}'.
func TestBourneAgainSimpleAppWithTwoArgs(t *testing.T) {
	buffer := new(bytes.Buffer)
	BourneAgain("myapp", []string{"mycommand", "anothercommand"}).Export(buffer)
	exportCapture := buffer.String()
	assert.Contains(t, exportCapture, "mycommand anothercommand")
	fmt.Println("Exported:\n", exportCapture)
}

// TestBourneAgainWithOneNested tests an app with usage 'myapp mycommand mysubcommand'.
func TestBourneAgainWithOneNested(t *testing.T) {
	buffer := new(bytes.Buffer)
	BourneAgain("myapp", []string{"mycommand", "mysubcommand"}).Export(buffer)
	exportCapture := buffer.String()
	assert.Contains(t, exportCapture, "mycommand")
	assert.Contains(t, exportCapture, "mysubcommand")
	fmt.Println("Exported:\n", exportCapture)
}

// TestBourneAgainWithLevel5 tests an app with usage 'myapp c1 c2 c3 c4 c5'.
func TestBourneAgainWithLevel5(t *testing.T) {
	buffer := new(bytes.Buffer)
	BourneAgain("myapp", []string{"c1 c2 c3 c4 c5"}).Export(buffer)
	exportCapture := buffer.String()
	assert.Contains(t, exportCapture, "c1")
	assert.Contains(t, exportCapture, "c2")
	assert.Contains(t, exportCapture, "c3")
	assert.Contains(t, exportCapture, "c4")
	assert.Contains(t, exportCapture, "c5")
	fmt.Println("Exported:\n", exportCapture)
}

// TestBourneAgainWithMultipleSubCommands tests an app with usage 'myapp c1 {c21|c22}'.
func TestBourneAgainWithMultipleSubCommands(t *testing.T) {
	buffer := new(bytes.Buffer)
	err := BourneAgain("myapp", []string{"c1 c21", "c1 c22"}).Export(buffer)
	assert.NoError(t, err)
	exportCapture := buffer.String()
	assert.Contains(t, exportCapture, "c1")
	assert.Contains(t, exportCapture, "c21 c22")
	fmt.Println("Exported:\n", exportCapture)
}

// TestNonValidTemplate tests the error handling of bash.
func TestNonValidTemplate(t *testing.T) {
	b := bash{appName: "app", commands: []string{}, tpl: "{{{{{}}nonsense"}
	err := b.Export(ioutil.Discard)
	assert.Error(t, err)
}
