package printer

import (
	"encoding/json"
	"io"

	"github.com/bpicode/fritzctl/console"
)

// Print arbitrates between printable types.
// If the passed argument is of type *console.Table, we print the table.
// If the passed argument is of any other type, we encode it as json.
func Print(data interface{}, writer io.Writer) {
	if table, ok := data.(*console.Table); ok {
		table.Print(writer)
		return
	}
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	encoder.Encode(data)
}
