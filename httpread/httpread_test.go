package httpread

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type errorReader struct {
}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error")
}

// TestErrorAtRequest reads from an error-prone source and asserts that the error is propagated.
func TestErrorAtRequest(t *testing.T) {
	clientPtr := &http.Client{}
	_, err := String(func() (*http.Response, error) {
		return clientPtr.Get("asfdnklfnlkanfknaf.afsajf.asfja")
	})
	assert.Error(t, err)
}

// TestError400 considers 400 bad request as response.
func TestError400(t *testing.T) {
	resp := &http.Response{StatusCode: 400, Status: "Bad Request", Body: ioutil.NopCloser(&strings.Reader{})}
	_, err := String(func() (*http.Response, error) {
		return resp, nil
	})
	assert.Error(t, err)
}

// TestSuccess follows the regular workflow.
func TestSuccess(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: ioutil.NopCloser(strings.NewReader("payload"))}
	body, err := String(func() (*http.Response, error) {
		return resp, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "payload", body)
}

// TestStatusCodeGuessing simulates a web server that does not handle status codes very well.
func TestStatusCodeGuessing(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: ioutil.NopCloser(strings.NewReader("HTTP/1.0 500 Internal Server Error\nContent-Length: 0\nContent-Type: text/plain; charset=utf-8"))}
	_, err := String(func() (*http.Response, error) {
		return resp, nil
	})
	assert.Error(t, err)
}

// TestXMLErrorAtRequest reads from an error-prone source and asserts that the error is propagated.
func TestXMLErrorAtRequest(t *testing.T) {
	clientPtr := &http.Client{}
	var payload string
	err := XML(func() (*http.Response, error) {
		return clientPtr.Get("asfdnklfnlkanfknaf.afsajf.asfja")
	}, &payload)
	assert.Error(t, err)
}

// TestXMLError400 considers 400 bad request as response.
func TestXMLError400(t *testing.T) {
	resp := &http.Response{StatusCode: 400, Status: "Bad Request", Body: ioutil.NopCloser(&strings.Reader{})}
	var payload string
	err := XML(func() (*http.Response, error) {
		return resp, nil
	}, &payload)
	assert.Error(t, err)
}

// TestXMLDecodeError considers a malformed, non-XML payload.
func TestXMLDecodeError(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: ioutil.NopCloser(strings.NewReader("no-xml"))}
	var payload string
	err := XML(func() (*http.Response, error) {
		return resp, nil
	}, &payload)
	assert.Error(t, err)
}

// TestReadFullyXMLSuccess follows the regular workflow.
func TestReadFullyXMLSuccess(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: ioutil.NopCloser(strings.NewReader("<dummy></dummy>"))}
	var payload string
	err := XML(func() (*http.Response, error) {
		return resp, nil
	}, &payload)
	assert.NoError(t, err)
}

// TestReadFullyJSON tests decoding into JSON.
func TestReadFullyJSON(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: ioutil.NopCloser(strings.NewReader(`{"a":"b"}`))}
	var payload struct {
		A string `json:"a"`
	}
	err := JSON(func() (*http.Response, error) {
		return resp, nil
	}, &payload)
	assert.NoError(t, err)
	assert.Equal(t, payload.A, "b")
}

// TestStringDecoder tests decoding into strings.
func TestStringDecoder(t *testing.T) {
	assert.Error(t, (&stringDecoder{reader: errorReader{}}).Decode(new(string)))
	assert.Error(t, (&stringDecoder{reader: strings.NewReader("somevalue")}).Decode(&struct{}{}))
	assert.NoError(t, (&stringDecoder{reader: strings.NewReader("somevalue")}).Decode(new(string)))
}
