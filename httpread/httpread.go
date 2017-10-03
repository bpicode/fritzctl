package httpread

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bpicode/fritzctl/logger"
	"github.com/pkg/errors"
)

var (
	httpStatusBuzzwords = map[string]int{"500 Internal Server Error": 500}
)

type stringDecoder struct {
	reader io.Reader
}

func (s *stringDecoder) Decode(v interface{}) error {
	bytesRead, err := ioutil.ReadAll(s.reader)
	sp := v.(*string)
	*sp = string(bytesRead)
	return err
}

// ReadFullyString reads a http response into a string.
// The response is checked for its status code and the http.Response.Body is closed.
func ReadFullyString(f func() (*http.Response, error)) (string, error) {
	body := ""
	err := readDecode(f, func(r io.Reader) decoder {
		return &stringDecoder{reader: r}
	}, &body)
	if err != nil {
		return "", err
	}
	sc, sp := guessStatusCode(body)
	if sc >= 400 {
		return "", fmt.Errorf("HTTP status code error (%d, guessed): remote replied with '%s'", sc, sp)
	}
	return body, nil
}

func guessStatusCode(body string) (int, string) {
	// There are web servers that send the wrong status code, but provide some hint in the text/html.
	for k, v := range httpStatusBuzzwords {
		if strings.Contains(strings.ToLower(body), strings.ToLower(k)) {
			return v, k
		}
	}
	return 0, ""
}

type decoder interface {
	Decode(v interface{}) error
}

type decoderFactory func(io.Reader) decoder

// ReadFullyXML reads a http response into a data container using an XML decoder.
// The response is checked for its status code and the http.Response.Body is closed.
func ReadFullyXML(f func() (*http.Response, error), v interface{}) error {
	return readDecode(f, func(r io.Reader) decoder {
		return xml.NewDecoder(r)
	}, v)
}

// ReadFullyJSON reads a http response into a data container using a json decoder.
// The response is checked for its status code and the http.Response.Body is closed.
func ReadFullyJSON(f func() (*http.Response, error), v interface{}) error {
	return readDecode(f, func(r io.Reader) decoder {
		return json.NewDecoder(r)
	}, v)
}

func readDecode(f func() (*http.Response, error), df decoderFactory, v interface{}) error {
	response, err := f()
	if err != nil {
		return errors.Wrap(err, "error obtaining HTTP response from remote")
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		return fmt.Errorf("HTTP status code error (%d): remote replied with '%s'", response.StatusCode, response.Status)
	}
	return decode(response.Body, df, v)
}

func decode(r io.Reader, df decoderFactory, v interface{}) error {
	buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)
	defer func() { logger.Debug("DATA:", buf) }()
	err := df(tee).Decode(v)
	if err != nil {
		return errors.Wrapf(err, "unable to decode remote response")
	}
	return nil
}
