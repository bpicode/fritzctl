package completion

import (
	"io"
	"text/template"
	"strings"

	"github.com/bpicode/fritzctl/stringutils"
)

// ShellExporter is the interface representing any shell having a completion feature.
type ShellExporter interface {
	Export(w io.Writer) error
	Add(cmd string)
	AddFlag(flag string)
}

type bash struct {
	appName  string
	commands []string
	tpl      string
}

type commandChain struct {
	AppName                string
	RootCommands           []string
	ParentVsDirectChildren map[string][]string
}

// BourneAgain instantiates a bash completion exporter.
func BourneAgain(appName string, commands []string) ShellExporter {
	return &bash{appName: appName, commands: commands, tpl: bashCompletionTemplate}
}

// Add adds a command to the stash of command. The argument cmd is expected
// to have the format '<C1> <C2> ... <CN>', e.g. 'clean', 'clean build' or
// 'clean build deploy' etc.
func (bash *bash) Add(cmd string) {
	bash.commands = append(bash.commands, cmd)
}

// AddFlag adds a global flag to the stash of commands, e.g. '--help'.
func (bash *bash) AddFlag(flag string) {
	bash.commands = append(bash.commands, stringutils.Transform(bash.commands, func(cmd string) string {
		return flag + " " + cmd
	})...)
}

// Export exports the completion script by writing it to an io.Writer.
func (bash *bash) Export(w io.Writer) error {
	tpl, err := template.New("completion.bash." + bash.appName).Parse(bash.tpl)
	if err != nil {
		return err
	}
	data := expandCommands(bash)
	return tpl.Execute(w, data)
}

func expandCommands(b *bash) commandChain {
	data := commandChain{
		AppName:                b.appName,
		ParentVsDirectChildren: make(map[string][]string),
	}
	data.pairParentAndDirectChildren(b.commands)
	return data
}

func (t *commandChain) pairParentAndDirectChildren(commands []string) {
	for _, cmd := range commands {
		t.addPairsForCommand(cmd)
	}
}

func (t *commandChain) addPairsForCommand(cmd string) {
	fields := strings.Fields(cmd)
	if len(fields) > 0 {
		t.addRoot(fields[0])
		t.addAdjacentPairs(fields)
	}
}
func (t *commandChain) addRoot(root string) {
	t.RootCommands = stringutils.AppendIfAbsent(t.RootCommands, root)
}

func (t *commandChain) addAdjacentPairs(fields []string) {
	for i, field := range fields[1:] {
		t.ParentVsDirectChildren[fields[i]] = stringutils.AppendIfAbsent(t.ParentVsDirectChildren[fields[i]], field)
	}
}
