package man_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/bpicode/fritzctl/man"
	"github.com/spf13/cobra"
)

// GenerateManPage generates a single man page describing all commands and sub-commands.
func ExampleGenerateManPage() {
	root := cobra.Command{
		Use: "myApp",
	}
	options := man.Options{
		Header: man.Header{
			Section: "1",
			Manual:  "myApp's man page",
			Title:   "myApp",
		},
		Origin: man.Origin{
			Date:   time.Date(2006, time.January, 1, 8, 0, 0, 0, time.UTC),
			Source: "written by a monkey on a typewriter",
		},
		SeeAlso: []string{"strace(1)"},
	}
	buf := new(bytes.Buffer)
	man.GenerateManPage(&root, &options, buf)
	s := buf.String()
	fmt.Println(s[0:26])
	// output:
	// .TH "myApp" "1" "Jan 2006"
}
