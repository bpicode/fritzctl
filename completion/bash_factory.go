package completion

import (
	"io"
	"text/template"
)

// ShellExporter is the interface representing any shell having a completion feature.
type ShellExporter interface {
	Export(w io.Writer) error
}

type bash struct {
	commands []string
	AppName  string
}

// BourneAgain instantiate a bash completion exporter.
func BourneAgain(appName string, commands []string) ShellExporter {
	return &bash{AppName: appName, commands: commands}
}

// Export exports the completion script by writing it ro an io.Writer.
func (bash *bash) Export(w io.Writer) error {
	tmpl, err := template.New(bash.AppName + "_outer").Parse(bashCompletionFunctionDefinition)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, bash)
}
