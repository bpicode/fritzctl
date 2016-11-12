package httpread

import (
	"net/http"
	"strings"
	"testing"

	"io"

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
	_, err := ReadFullyString(clientPtr.Get("asfdnklfnlkanfknaf.afsajf.asfja"))
	assert.Error(t, err)
}

// TestReadFullyError400 considers 400 bad request as response.
func TestReadFullyError400(t *testing.T) {
	resp := &http.Response{StatusCode: 400, Status: "Bad Request", Body: dummyCloser{Reader: &strings.Reader{}}}
	_, err := ReadFullyString(resp, nil)
	assert.Error(t, err)
}

// TestReadFullySuccess follows the regular workflow.
func TestReadFullySuccess(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: dummyCloser{Reader: strings.NewReader("payload")}}
	body, err := ReadFullyString(resp, nil)
	assert.NoError(t, err)
	assert.Equal(t, "payload", body)
}

// TestReadFullyXMLErrorAtRequest reads from an error-prone source and asserts that the error is propagated.
func TestReadFullyXMLErrorAtRequest(t *testing.T) {
	clientPtr := &http.Client{}
	var payload string
	a, b := clientPtr.Get("asfdnklfnlkanfknaf.afsajf.asfja")
	err := ReadFullyXML(a, b, &payload)
	assert.Error(t, err)
}

// TestReadFullyXMLError400 considers 400 bad request as response.
func TestReadFullyXMLError400(t *testing.T) {
	resp := &http.Response{StatusCode: 400, Status: "Bad Request", Body: dummyCloser{Reader: &strings.Reader{}}}
	var payload string
	err := ReadFullyXML(resp, nil, &payload)
	assert.Error(t, err)
}

// TestReadFullyXMLDecodeError considers a malformed, non-XML payload.
func TestReadFullyXMLDecodeError(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: dummyCloser{Reader: strings.NewReader("no-xml")}}
	var payload string
	err := ReadFullyXML(resp, nil, &payload)
	assert.Error(t, err)
}

// TestReadFullyXMLSuccess follows the regular workflow.
func TestReadFullyXMLSuccess(t *testing.T) {
	resp := &http.Response{StatusCode: 200, Status: "OK", Body: dummyCloser{Reader: strings.NewReader("<dummy></dummy>")}}
	var payload string
	err := ReadFullyXML(resp, nil, &payload)
	assert.NoError(t, err)
}
