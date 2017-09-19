package console

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

// Question is a container for confronting the user to decide on an answer.
type Question struct {
	Key          string                 // Key identifies the target field name.
	Text         string                 // Text is presented to the user.
	Converter    Converter              // Converter should map text input to a target type.
	CustomSource func() (string, error) // CustomSource may replace the custom input source.
	Defaulter    func() interface{}     // Defaulter may supply a fallback value when an empty input is supplied.
}

// Converter converts strings to any type.
type Converter interface {
	// Convert performs the conversion.
	Convert(s string) (interface{}, error)
}

type toString struct {
}

// Convert performs the conversion.
func (a *toString) Convert(s string) (interface{}, error) {
	return s, nil
}

type toBool struct {
}

// Convert performs the conversion.
func (a *toBool) Convert(s string) (interface{}, error) {
	return strconv.ParseBool(s)
}

// ForString creates a Question with a string as target value.
func ForString(key, text, def string) Question {
	return Question{Key: key, Text: text, Converter: &toString{}, Defaulter: func() interface{} {
		return def
	}}
}

// ForBool creates a Question with a bool as target value.
func ForBool(key, text string, def bool) Question {
	return Question{Key: key, Text: text, Converter: &toBool{}, Defaulter: func() interface{} {
		return def
	}}
}

// ForPassword creates a Question with a sting as target value. The input from the terminal will not be echoed.
func ForPassword(key, text string) Question {
	return Question{Key: key, Text: text, Converter: &toString{}, CustomSource: func() (string, error) {
		pwBytes, err := terminal.ReadPassword(0)
		fmt.Println()
		return string(pwBytes), err
	}}
}

// Survey contains the configuration on how to present/obtain questions/answers to/from the user.
type Survey struct {
	In  io.Reader // In is the input source, e.g. os.Stdin.
	Out io.Writer // Out is the output sink, e.g. os.Stdout.
}

// Ask confronts the user with the passed questions. Questions are traversed in order.
func (s *Survey) Ask(qs []Question, v interface{}) error {
	scanner := bufio.NewScanner(s.In)
	as := make(map[string]interface{})
	for _, q := range qs {
		q.prompt(s.Out)
		a, err := q.obtain(scanner)
		if err != nil {
			return fmt.Errorf("could not complete survey for key '%s': %s", q.Key, err)
		}
		as[q.Key] = a
	}
	return s.writeTo(as, v)
}

func (s *Survey) writeTo(m map[string]interface{}, v interface{}) error {
	j, _ := json.Marshal(m)
	return json.Unmarshal(j, v)
}

func (q Question) prompt(w io.Writer) {
	hint := ""
	if q.Defaulter != nil {
		hint = fmt.Sprintf(" [%s]", q.Defaulter())
	}
	cyan := color.New(color.FgCyan).SprintfFunc()
	bold := color.New(color.Bold).SprintfFunc()
	fmt.Fprintf(w, "%s %s%s: ", cyan("?"), bold("%s", q.Text), hint)
}

func (q *Question) obtain(s *bufio.Scanner) (interface{}, error) {
	answer, err := q.getAnswer(s)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(answer) == "" && q.Defaulter != nil {
		return q.Defaulter(), nil
	}
	return q.Converter.Convert(answer)
}

func (q *Question) getAnswer(s *bufio.Scanner) (string, error) {
	if q.CustomSource != nil {
		return q.CustomSource()
	}
	ok := s.Scan()
	if !ok {
		return "", fmt.Errorf("could not scan: error or premature end of input: %s", s.Err())
	}
	return s.Text(), nil
}
