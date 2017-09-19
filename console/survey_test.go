package console

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAskZeroQuestions tests the special case when no questions are asked.
func TestAskZeroQuestions(t *testing.T) {
	s := Survey{In: os.Stdin, Out: os.Stdout}
	err := s.Ask([]Question{}, &struct{}{})
	assert.NoError(t, err)
}

// TestAsk tests a survey of questions.
func TestAsk(t *testing.T) {
	response := struct {
		MyString string
		MyBool   bool
	}{}
	reader := bytes.NewReader([]byte(`some string
true
`))
	s := Survey{In: reader, Out: os.Stdout}
	err := s.Ask([]Question{
		ForString("myString", "Enter some string:", ""),
		ForBool("myBool", "Enter some bool:", false),
	}, &response)
	assert.NoError(t, err)
	assert.Equal(t, "some string", response.MyString)
	assert.Equal(t, true, response.MyBool)
}

// TestAskForPassword tests a survey for a password.
func TestAskForPassword(t *testing.T) {
	response := struct{ Password string }{}
	s := Survey{Out: os.Stdout}
	assert.NotPanics(t, func() {
		s.Ask([]Question{
			ForPassword("Password", "Enter password:"),
		}, &response)
	})
}

type faultyReader struct {
}

// Read does always fail on the faultyReader.
func (f faultyReader) Read(b []byte) (n int, err error) {
	return 0, errors.New("I always fail as a reader")
}

// TestAskWithFaultyReader tests a survey's error handling.
func TestAskWithFaultyReader(t *testing.T) {
	s := Survey{In: faultyReader{}, Out: os.Stdout}
	err := s.Ask([]Question{
		ForString("s", "Enter some string:", ""),
	}, &struct{}{})
	assert.Error(t, err)
}

// TestAskWithDefault tests a survey of with a default value.
func TestAskWithDefault(t *testing.T) {
	response := struct {
		MyString string
		MyBool   bool
	}{}
	reader := bytes.NewReader([]byte(`

`))
	s := Survey{In: reader, Out: os.Stdout}
	s.Ask([]Question{
		ForString("myString", "Enter some string:", "DEFAULT"),
		ForBool("myBool", "Enter some bool:", true),
	}, &response)
	assert.Equal(t, "DEFAULT", response.MyString)
	assert.Equal(t, true, response.MyBool)
}

// TestSequentialSurveys tests a survey followed by another survey.
func TestSequentialSurveys(t *testing.T) {
	response := struct {
		MyString1 string
		MyString2 string
	}{}

	s1 := Survey{In: bytes.NewReader([]byte("a\n")), Out: os.Stdout}
	s1.Ask([]Question{
		ForString("myString1", "Enter some string:", ""),
	}, &response)

	s2 := Survey{In: bytes.NewReader([]byte("b\n")), Out: os.Stdout}
	s2.Ask([]Question{
		ForString("myString2", "Enter some string:", ""),
	}, &response)
	assert.Equal(t, "a", response.MyString1)
	assert.Equal(t, "b", response.MyString2)
}
