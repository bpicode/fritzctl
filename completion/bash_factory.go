package completion

// Shell is the interface representing any shell having a completion feature.
type Shell interface {
	Print()
}

type Bash struct {
	commands []string
	appName string
}

func BourneAgain(appName string, commands []string) Shell {
	return &Bash{appName: appName, commands: commands}
}

func (bash *Bash) Print() {

}
