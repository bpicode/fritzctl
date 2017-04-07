package completion

import (
	"io"
	"text/template"
	"fmt"
)

// ShellExporter is the interface representing any shell having a completion feature.
type ShellExporter interface {
	Export(w io.Writer) error
}

type bash struct {
	appName  string
	commands []string
}

type command struct {
	Level    int
	Name     string
	Children []command
}

type applicationData struct {
	AppName         string
	Commands        []command
	LevelVsCommands map[int][]command
	Flags           []string
}

// BourneAgain instantiate a bash completion exporter.
func BourneAgain(appName string, commands []string) ShellExporter {
	return &bash{appName: appName, commands: commands}
}

// Export exports the completion script by writing it ro an io.Writer.
func (bash *bash) Export(w io.Writer) error {
	tpl, err := template.New(bash.appName + "_outer").Parse(bashCompletionFunctionDefinition)
	if err != nil {
		return err
	}
	data := expandCommands(bash)
	return tpl.Execute(w, data)
}
func expandCommands(bash *bash) applicationData {
	var commandTable [][]string
	fmt.Println("COMMAND TABLE", commandTable)
	xx := make(map[int][]command)
	xx[1] = make([]command, 0)
	xx[1] = append(xx[1], command{Name: "mycommand"})
	data := applicationData{AppName: bash.appName, LevelVsCommands: xx}
	return data
}
