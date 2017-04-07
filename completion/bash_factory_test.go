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
	assert.Contains(t, exportCapture, "mycommand")
	fmt.Println("Exported:\n", exportCapture)
}

// TestBourneAgainSimpleAppWithTwoArgs tests an app with usage 'myapp {mycommand|anothercomand}'.
func TestBourneAgainSimpleAppWithTwoArgs(t *testing.T) {
	bash := BourneAgain("myapp", []string{"mycommand", "anothercommand"})
	buffer := new(bytes.Buffer)
	err := bash.Export(buffer)
	assert.NoError(t, err)
	exportCapture := buffer.String()
	assert.Contains(t, exportCapture, "mycommand anothercommand")
	fmt.Println("Exported:\n", exportCapture)
}

// TestBourneAgainWithOneNested tests an app with usage 'myapp mycommand mysubcommand'.
func TestBourneAgainWithOneNested(t *testing.T) {
	b := &bash{appName: "myapp", commands: []string{"mycommand mysubcommand"}}
	buffer := new(bytes.Buffer)
	err := b.exportByExpanding(buffer, func(b *bash) applicationData {
		commands := make(map[int][]command)
		mySubCommand := command{Name: "mysubcommand"}
		myCommand := command{Name: "mycommand", Children: []command{mySubCommand}}
		commands[1] = append(commands[1], myCommand)
		commands[2] = append(commands[2], myCommand)
		return applicationData{
			AppName:         b.appName,
			LevelVsCommands: commands,
		}
	})
	assert.NoError(t, err)
	exportCapture := buffer.String()
	assert.Contains(t, exportCapture, "mycommand")
	assert.Contains(t, exportCapture, "mysubcommand")
	fmt.Println("Exported:\n", exportCapture)
}

// TestBourneAgainWithLevel5 tests an app with usage 'myapp c1 c2 c3 c4 c5'.
func TestBourneAgainWithLevel5(t *testing.T) {
	b := &bash{appName: "myapp", commands: []string{"c1 c2 c3 c4 c5"}}
	buffer := new(bytes.Buffer)
	err := b.exportByExpanding(buffer, func(b *bash) applicationData {
		commands := make(map[int][]command)
		c5 := command{Name: "c5"}
		c4 := command{Name: "c4", Children: []command{c5}}
		c3 := command{Name: "c3", Children: []command{c4}}
		c2 := command{Name: "c2", Children: []command{c3}}
		c1 := command{Name: "c1", Children: []command{c2}}
		commands[1] = append(commands[1], c1)
		commands[2] = append(commands[2], c1)
		commands[3] = append(commands[3], c2)
		commands[4] = append(commands[4], c3)
		commands[5] = append(commands[5], c4)
		return applicationData{
			AppName:         b.appName,
			LevelVsCommands: commands,
		}
	})
	assert.NoError(t, err)
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
	b := &bash{appName: "myapp", commands: []string{"c1 c21", "c1 c22"}}
	buffer := new(bytes.Buffer)
	err := b.exportByExpanding(buffer, func(b *bash) applicationData {
		commands := make(map[int][]command)
		c22 := command{Name: "c22"}
		c21 := command{Name: "c21"}
		c1 := command{Name: "c1", Children: []command{c21, c22}}
		commands[1] = append(commands[1], c1)
		commands[2] = append(commands[2], c1)
		return applicationData{
			AppName:         b.appName,
			LevelVsCommands: commands,
		}
	})
	assert.NoError(t, err)
	exportCapture := buffer.String()
	assert.Contains(t, exportCapture, "c1")
	assert.Contains(t, exportCapture, "c21 c22")
	fmt.Println("Exported:\n", exportCapture)
}
