package httpread

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyCloser struct {
	io.Reader
}

func (dummyCloser) Close() error {
	return nil
}

// TestReadFullyErrorAtRequest reads from an error-prone source and asserts that the error is propagated.
func TestReadFullyErrorAtRequest(t *testing.T) {
	clientPtr := &http.Client{}
	_, err := ReadFullyString(func() (*http.Response, error) {
		return clientPtr.Get("asfdnklfnlkanfknaf.afsajf.asfja")
	})
	assert.Error(t, err)
}

// TestReadFullyError400 considers 400 bad request as response.
func TestReadFullyError400(t *testing.T) {
	resp := &http.Response{StatusCode: 400, Status: "Bad Request", Body: dummyCloser{Reader: &strings.Reader{}}}
	_, err := ReadFullyString(func() (*http.Response, error) {
		return resp, nil
	})
	assert.Error(t, err)
}

// TestReadFullySuccess follows the regular workflow.
func TestReadFullySuccess(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: dummyCloser{Reader: strings.NewReader("payload")}}
	body, err := ReadFullyString(func() (*http.Response, error) {
		return resp, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, "payload", body)
}

// TestReadFullyWithStatusCodeGuessing simulates a web server that does not handle status codes very well.
func TestReadFullyWithStatusCodeGuessing(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: dummyCloser{Reader: strings.NewReader("HTTP/1.0 500 Internal Server Error\nContent-Length: 0\nContent-Type: text/plain; charset=utf-8")}}
	_, err := ReadFullyString(func() (*http.Response, error) {
		return resp, nil
	})
	assert.Error(t, err)
}

// TestReadFullyXMLErrorAtRequest reads from an error-prone source and asserts that the error is propagated.
func TestReadFullyXMLErrorAtRequest(t *testing.T) {
	clientPtr := &http.Client{}
	var payload string
	err := ReadFullyXML(func() (*http.Response, error) {
		return clientPtr.Get("asfdnklfnlkanfknaf.afsajf.asfja")
	}, &payload)
	assert.Error(t, err)
}

// TestReadFullyXMLError400 considers 400 bad request as response.
func TestReadFullyXMLError400(t *testing.T) {
	resp := &http.Response{StatusCode: 400, Status: "Bad Request", Body: dummyCloser{Reader: &strings.Reader{}}}
	var payload string
	err := ReadFullyXML(func() (*http.Response, error) {
		return resp, nil
	}, &payload)
	assert.Error(t, err)
}

// TestReadFullyXMLDecodeError considers a malformed, non-XML payload.
func TestReadFullyXMLDecodeError(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: dummyCloser{Reader: strings.NewReader("no-xml")}}
	var payload string
	err := ReadFullyXML(func() (*http.Response, error) {
		return resp, nil
	}, &payload)
	assert.Error(t, err)
}

// TestReadFullyXMLSuccess follows the regular workflow.
func TestReadFullyXMLSuccess(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: dummyCloser{Reader: strings.NewReader("<dummy></dummy>")}}
	var payload string
	err := ReadFullyXML(func() (*http.Response, error) {
		return resp, nil
	}, &payload)
	assert.NoError(t, err)
}

// TestReadFullyJSON tests decoding into JSOn.
func TestReadFullyJSON(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: dummyCloser{Reader: strings.NewReader(`{"a":"b"}`)}}
	var payload struct {
		A string `json:"a"`
	}
	err := ReadFullyJSON(func() (*http.Response, error) {
		return resp, nil
	}, &payload)
	assert.NoError(t, err)
	assert.Equal(t, payload.A, "b")
}
