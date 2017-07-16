package man

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/cpuguy83/go-md2man/md2man"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Options specify the common properties of the man page.
type Options struct {
	Header  Header   // General data.
	Origin  Origin   // Data representing the origin of this man page.
	SeeAlso []string // Links to other man pages, e.g. "ps(1)"
}

// Origin conveys data of the man page origin.
type Origin struct {
	Source string    // Origin of the man page.
	Date   time.Time // The Date when the man page was written.
}

// Header conveys general data of the man page.
type Header struct {
	// Man page title.
	Title string
	// Man page section.
	// Use "1" for General commands.
	// Use "2" for System calls.
	// Use "3" for Library functions, covering in particular the C standard library.
	// Use "4" for Special files (usually devices, those found in /dev) and drivers.
	// Use "5" for File formats and conventions. Use "6" for Games and screensavers.
	// Use "7" for Miscellanea.
	// Use "8" for System administration commands and daemons.
	Section string
	// Manual title.
	Manual string
}

type mdBuffer struct {
	*bytes.Buffer
}

func (buf *mdBuffer) printfln(format string, a ...interface{}) (int, error) {
	return buf.WriteString(fmt.Sprintf(format+"\n", a...))
}

func (buf *mdBuffer) printline(s string) (int, error) {
	return buf.WriteString(s + "\n")
}

func (buf *mdBuffer) boldln(s string) (int, error) {
	return buf.printfln("**%s**", s)
}

func (buf *mdBuffer) header(title string) (int, error) {
	return buf.printfln("# %s", title)
}

func (buf *mdBuffer) code(s string) (int, error) {
	return buf.printfln("```\n%s\n```", s)
}

func newMdBuffer() *mdBuffer {
	buf := new(mdBuffer)
	buf.Buffer = new(bytes.Buffer)
	return buf
}

// GenerateManPage writes the man page, taking the given command as root, to a writer.
func GenerateManPage(cmd *cobra.Command, options *Options, w io.Writer) error {
	buf := newMdBuffer()
	writeMetadata(options, buf)
	writeName(cmd.Name(), cmd.Short, buf)
	writeSynopsis(cmd.UseLine(), buf)
	writeDescription(cmd.Long, buf)
	writeOptions(cmd.Flags(), buf)
	writeCommands(cmd, buf)
	writeExamples(cmd, buf)
	writeExitStatus(buf)
	writeSeeAlso(options.SeeAlso, buf)
	bytesMarkdown := buf.Bytes()
	bytesMan := md2man.Render(bytesMarkdown)
	_, err := w.Write(bytesMan)
	return err
}

func writeMetadata(options *Options, buf *mdBuffer) {
	buf.printfln("%% %s(%s)%s", options.Header.Title, options.Header.Section, options.Origin.Date.Format("Jan 2006"))
	buf.printfln("%% %s", options.Origin.Source)
	buf.printfln("%% %s", options.Header.Manual)
}

func writeName(name, short string, buf *mdBuffer) {
	buf.header("NAME")
	buf.printfln("%s \\- %s", name, short)
}

func writeSynopsis(use string, buf *mdBuffer) {
	buf.header("SYNOPSIS")
	buf.boldln(use)
}

func writeDescription(description string, buf *mdBuffer) {
	buf.header("DESCRIPTION")
	buf.printline(description)
}

func writeOptions(flags *pflag.FlagSet, buf *mdBuffer) {
	buf.header("OPTIONS")
	flags.VisitAll(func(flag *pflag.Flag) {
		writeOption(flag, buf)
	})
}

func writeOption(flag *pflag.Flag, buf *mdBuffer) {
	buf.printfln("**--%s**=%q\n\t%s\n", flag.Name, flag.DefValue, flag.Usage)
}

func writeCommands(cmd *cobra.Command, buf *mdBuffer) {
	buf.header("COMMANDS")
	buf.printline("The following commands are understood:\n")
	writeCommandRecursive(cmd, buf)
}

func writeCommandRecursive(cmd *cobra.Command, buf *mdBuffer) {
	if cmd.RunE != nil || cmd.Run != nil {
		buf.printfln("**%s** - %s\n", cmd.CommandPath(), cmd.Short)
		buf.printfln("\t\t%s\n", cmd.Long)
	}
	for _, sub := range cmd.Commands() {
		writeCommandRecursive(sub, buf)
	}
}

func writeExamples(cmd *cobra.Command, buf *mdBuffer) {
	buf.header("EXAMPLES")
	writeExampleRecursive(cmd, buf)
}

func writeExampleRecursive(cmd *cobra.Command, buf *mdBuffer) {
	if cmd.Example != "" {
		buf.boldln(cmd.Short)
		buf.code(cmd.Example)
	}
	for _, sub := range cmd.Commands() {
		writeExampleRecursive(sub, buf)
	}
}

func writeExitStatus(buf *mdBuffer) {
	buf.header("EXIT STATUS")
	buf.printline("On success, 0 is returned, a non-zero failure code otherwise. " +
		"If the return code is non-zero, look at the output to get a hint on what went wrong.")
}

func writeSeeAlso(sa []string, buf *mdBuffer) {
	buf.header("SEE ALSO")
	buf.boldln(strings.Join(sa, ", "))
}
