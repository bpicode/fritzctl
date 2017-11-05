package httpread_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bpicode/fritzctl/httpread"
)

// JSON reads marshals the body of a http.Response to a struct.
func ExampleJSON() {
	f := func() (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(strings.NewReader(`{"x": 2, "y": 5}`))}, nil
	}
	var pt = struct {
		X int
		Y int
	}{}
	httpread.JSON(f, &pt)
	fmt.Println(pt)
	// output: {2 5}
}

// String reads the body of a http.Response to a string.
func ExampleString() {
	f := func() (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(strings.NewReader("text"))}, nil
	}
	s, _ := httpread.String(f)
	fmt.Println(s)
	// output: text
}

// XML reads marshals the body of a http.Response to a struct.
func ExampleXML() {
	f := func() (*http.Response, error) {
		return &http.Response{Body: ioutil.NopCloser(strings.NewReader(`<point><x>2</x><y>5</y></point>`))}, nil
	}
	var pt = struct {
		X int `xml:"x"`
		Y int `xml:"y"`
	}{}
	httpread.XML(f, &pt)
	fmt.Println(pt)
	// output: {2 5}
}
