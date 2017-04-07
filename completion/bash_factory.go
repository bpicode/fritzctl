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
	appName  string
	commands []string
}

type command struct {
	Name     string
	Children []command
}

type applicationData struct {
	AppName         string
	Commands        []command
	LevelVsCommands map[int][]command
	Flags           []string
}

type commandExpander func(*bash) applicationData

// BourneAgain instantiate a bash completion exporter.
func BourneAgain(appName string, commands []string) ShellExporter {
	return &bash{appName: appName, commands: commands}
}

// Export exports the completion script by writing it ro an io.Writer.
func (bash *bash) Export(w io.Writer) error {
	return bash.exportByExpanding(w, expandCommands)
}

func (bash *bash) exportByExpanding(w io.Writer, e commandExpander) error {
	tpl, err := template.New("completion.bashbash." + bash.appName).Parse(bashCompletionFunctionDefinition)
	if err != nil {
		return err
	}
	data := e(bash)
	return tpl.Execute(w, data)
}

func expandCommands(b *bash) applicationData {
	commandMap := commandMap(b.commands)
	data := applicationData{AppName: b.appName, LevelVsCommands: commandMap}
	return data
}

func commandMap(commands []string) map[int][]command {
	cmdMap := make(map[int][]command)
	cmdMap[1] = make([]command, 0)
	for _, cmd := range commands {
		cmdMap[1] = append(cmdMap[1], command{Name: cmd})
	}
	return cmdMap
}
