package console

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

var ansi = regexp.MustCompile("\033\\[(?:[0-9]{1,3}(?:;[0-9]{1,3})*)?[m|K]")

// Table contains the data to draw a table.
type Table struct {
	headers []string
	rows    [][]string
}

// Option allows to mutate the table.
type Option func(t *Table)

// NewTable creates a new table and applies the given Options to it.
func NewTable(opts ...Option) *Table {
	table := &Table{}
	for _, opt := range opts {
		opt(table)
	}
	return table
}

// Headers is an option to store the header texts in the table.
func Headers(hs ...string) Option {
	return func(t *Table) {
		t.headers = hs
	}
}

// Append extends the table body by one line. The length of the row should match the length of the headers as well as
// the other rows.
func (t *Table) Append(row []string) {
	t.rows = append(t.rows, row)
}

// Print writes the text representation of the table to a given io.Writer.
func (t *Table) Print(w io.Writer) {
	t.printHeaders(w)
	t.printRows(w)
	t.printRuler(w)
}

func (t *Table) printHeaders(w io.Writer) (int64, error) {
	buf := new(bytes.Buffer)
	t.printRuler(buf)
	t.printTitles(buf)
	t.printRuler(buf)
	return buf.WriteTo(w)
}

func (t *Table) printRuler(w io.Writer) (int64, error) {
	buf := new(bytes.Buffer)
	for i := range t.headers {
		width := t.widthOf(i)
		dashes := strings.Repeat("-", width+2)
		fmt.Fprintf(buf, "+%s", dashes)
	}
	fmt.Fprintln(buf, "+")
	return buf.WriteTo(w)
}

func (t *Table) printTitles(w io.Writer) {
	for i, h := range t.headers {
		width := t.widthOf(i)
		leftPad := strings.Repeat(" ", (width-len(h))/2)
		rightPad := strings.Repeat(" ", (width-len(h)+1)/2)
		fmt.Fprintf(w, "| %s%s%s ", leftPad, h, rightPad)
	}
	fmt.Fprintln(w, "|")
}

func (t *Table) printRows(w io.Writer) (int64, error) {
	buf := new(bytes.Buffer)
	for _, r := range t.rows {
		t.printRow(r, w)
	}
	return buf.WriteTo(w)
}

func (t *Table) printRow(row []string, w io.Writer) {
	for i, c := range row {
		t.printCol(i, c, w)
	}
	fmt.Fprintln(w, "|")
}

func (t *Table) printCol(i int, val string, w io.Writer) {
	width := t.widthOf(i)
	white := strings.Repeat(" ", max(width-runeLen(val), 0))
	if isNumeric(val) {
		fmt.Fprintf(w, "| %s%s ", white, val)
	} else {
		fmt.Fprintf(w, "| %s%s ", val, white)
	}
}

func isNumeric(s string) bool {
	fields := strings.Fields(s)
	if len(fields) == 0 {
		return false
	}
	_, err := strconv.ParseFloat(fields[0], 64)
	return err == nil
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func (t *Table) widthOf(i int) int {
	width := runeLen(t.headers[i])
	for _, row := range t.rows {
		if runeLen(row[i]) > width {
			width = runeLen(row[i])
		}
	}
	return width
}

func runeLen(s string) int {
	noAnsi := ansi.ReplaceAllLiteralString(s, "")
	l := 0
	for range noAnsi {
		l++
	}
	return l
}
