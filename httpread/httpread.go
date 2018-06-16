package httpread

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bpicode/fritzctl/internal/errors"
	"github.com/bpicode/fritzctl/logger"
)

var (
	httpStatusBuzzwords = map[string]int{"500 Internal Server Error": 500}
)

type stringDecoder struct {
	reader io.Reader
}

func (s *stringDecoder) Decode(v interface{}) error {
	bytesRead, err := ioutil.ReadAll(s.reader)
	if err != nil {
		return err
	}
	sp, ok := v.(*string)
	if !ok {
		return fmt.Errorf("cannot decode into string, call with string pointer")
	}
	*sp = string(bytesRead)
	return nil
}

// String reads a http response into a string.
// The response is checked for its status code and the http.Response.Body is closed.
func String(f func() (*http.Response, error)) (string, error) {
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

type csvDecoder struct {
	reader io.Reader
	comma  rune
}

func (c *csvDecoder) Decode(v interface{}) error {
	cr := csv.NewReader(c.reader)
	cr.Comma = c.comma
	cr.FieldsPerRecord = -1
	records, err := cr.ReadAll()
	if err != nil {
		return err
	}
	t, ok := v.(*[][]string)
	if !ok {
		return fmt.Errorf("cannot decode into csv, call with *[][]string slice")
	}
	*t = records
	return nil
}

// Csv reads a http response into a [][]string.
// The response is checked for its status code and the http.Response.Body is closed.
func Csv(f func() (*http.Response, error), comma rune) ([][]string, error) {
	var records [][]string
	err := readDecode(f, func(r io.Reader) decoder {
		return &csvDecoder{reader: r, comma: comma}
	}, &records)
	return records, err
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

// XML reads a http response into a data container using an XML decoder.
// The response is checked for its status code and the http.Response.Body is closed.
func XML(f func() (*http.Response, error), v interface{}) error {
	return readDecode(f, func(r io.Reader) decoder {
		return xml.NewDecoder(r)
	}, v)
}

// JSON reads a http response into a data container using a json decoder.
// The response is checked for its status code and the http.Response.Body is closed.
func JSON(f func() (*http.Response, error), v interface{}) error {
	return readDecode(f, func(r io.Reader) decoder {
		return json.NewDecoder(r)
	}, v)
}

func readDecode(f func() (*http.Response, error), df decoderFactory, v interface{}) error {
	response, err := f()
	if err != nil {
		return errors.Wrapf(err, "error obtaining HTTP response from remote")
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
