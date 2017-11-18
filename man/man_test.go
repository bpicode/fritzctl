package man

import (
	"bytes"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// TestGenerate tests the man page generation.
func TestGenerate(t *testing.T) {
	buf := new(bytes.Buffer)
	err := Generate(exampleCommand(), &Options{
		Header: Header{
			Title:   "myApp",
			Section: "1",
			Manual:  "myApp man page",
		},
		Origin: Origin{
			Source: "Generated in test",
			Date:   time.Now(),
		},
		SeeAlso: []string{"strace(1)"},
	}, buf)
	assert.NoError(t, err)

	captured := buf.String()
	assert.NotEmpty(t, captured)
}

func exampleCommand() *cobra.Command {
	root := cobra.Command{Short: "root", Long: "the root cmd", Use: "myApp"}
	sub := cobra.Command{
		Short:   "sub",
		Long:    "the sub cmd",
		Use:     "sub",
		Example: "myApp sub", Run: func(cmd *cobra.Command, args []string) {},
	}
	root.AddCommand(&sub)
	root.InitDefaultHelpFlag()
	root.InitDefaultHelpCmd()
	return &root
}
