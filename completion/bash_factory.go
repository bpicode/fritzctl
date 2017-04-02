package completion

// Shell is the interface representing any shell having a completion feature.
type Shell interface {
	Print()
}

type Bash struct {
	commands []string
}

func BourneAgain(commands []string) Shell {
	return &Bash{commands: commands}
}

func (bash *Bash) Print() {

}

