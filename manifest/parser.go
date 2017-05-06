package manifest

import (
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// ParseFile opens a file and marshals the contents as a Plan, returning a pointer it.
func ParseFile(filename string) (*Plan, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return Parse(file)
}

// Parse takes an io.Reader and marshals the contents as a Plan, returning a pointer it.
func Parse(r io.Reader) (*Plan, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var state Plan
	err = yaml.Unmarshal(bytes, &state)
	return &state, err
}
