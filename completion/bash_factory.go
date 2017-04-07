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
	level   int
	name    string
	parents []command
}

// BourneAgain instantiate a bash completion exporter.
func BourneAgain(appName string, commands []string) ShellExporter {
	return &bash{appName: appName, commands: commands}
}

// Export exports the completion script by writing it ro an io.Writer.
func (bash *bash) Export(w io.Writer) error {
	tmpl, err := template.New(bash.appName + "_outer").Parse(bashCompletionFunctionDefinition)
	if err != nil {
		return err
	}
	type applicationData struct {
		AppName  string
		Commands []string
		Flags    []string
	}

	var commandTable [][]string
	fmt.Println("COMMAND TABLE", commandTable)

	return tmpl.Execute(w, applicationData{AppName: bash.appName, Commands: bash.commands})
}
