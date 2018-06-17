package console_test

import (
	"bytes"
	"fmt"
	"os"

	"github.com/bpicode/fritzctl/internal/console"
)

// Checkmark represents a state of one the three: "✘", "✔", "?". It can be parsed from several primitives.
func ExampleCheckmark() {
	fmt.Println(console.Btoc(true))
	fmt.Println(console.Btoc(false))

	fmt.Println(console.Itoc(1))
	fmt.Println(console.Itoc(0))

	fmt.Println(console.Stoc("1"))
	fmt.Println(console.Stoc("0"))
	fmt.Println(console.Stoc(""))
	// output:
	// ✔
	// ✘
	// ✔
	// ✘
	// ✔
	// ✘
	// ?
}

// Table generally wraps borders and separators around structured data.
func ExampleTable() {
	t := console.NewTable(console.Headers("NAME", "AGE", "ADDRESS"))
	t.Append([]string{"John", "35", "14th Oak Rd"})
	t.Append([]string{"Jane", "24", "27th Elm St."})
	t.Append([]string{"Tom", "55", "4th Maple Rd."})
	t.Print(os.Stdout)
	// output:
	// +------+-----+---------------+
	// | NAME | AGE |    ADDRESS    |
	// +------+-----+---------------+
	// | John |  35 | 14th Oak Rd   |
	// | Jane |  24 | 27th Elm St.  |
	// | Tom  |  55 | 4th Maple Rd. |
	// +------+-----+---------------+
}

// Survey is a sequential print-read process to obtain data, e.g. from stdin.
func ExampleSurvey() {
	r := bytes.NewReader([]byte("example.com\n"))
	s := console.Survey{In: r, Out: os.Stdout}
	t := struct{ Host string }{}
	s.Ask([]console.Question{
		console.ForString("Host", "Enter host", "localhost"),
	}, &t)
	fmt.Println(t.Host)
	// output: ? Enter host [localhost]: example.com
}
