package httpread

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	httpStatusBuzzwords = map[string]int{"500 Internal Server Error": 500}
)

// HTTPStatusCodeError represents an 4xx client or a 5xx server error.
type HTTPStatusCodeError struct {
	error
}

func statusCodeError(code int, phrase string) *HTTPStatusCodeError {
	return &HTTPStatusCodeError{error: fmt.Errorf("HTTP status code error (%d): remote replied with %s", code, phrase)}
}

// ReadFullyString reads a http response into a string.
// The response is checked for its status code and the http.Response.Body is closed.
func ReadFullyString(response *http.Response, errorObtainResponse error) (string, error) {
	if errorObtainResponse != nil {
		return "", errorObtainResponse
	}
	defer response.Body.Close()
	bytesRead, errRead := ioutil.ReadAll(response.Body)
	body := string(bytesRead)
	statusCode, statusPhrase := guessStatusCode(response.StatusCode, response.Status, body)
	if statusCode >= 400 {
		return body, statusCodeError(statusCode, statusPhrase)
	}
	return body, errRead
}

func guessStatusCode(claimedCode int, claimedPhrase, body string) (int, string) {
	if claimedCode >= 400 {
		return claimedCode, claimedPhrase // This is already bad enough.
	}
	// There are web servers that send the wrong status code, but provide some hint in the text/html.
	for k, v := range httpStatusBuzzwords {
		if strings.Contains(strings.ToLower(body), strings.ToLower(k)) {
			return v, k
		}
	}
	return claimedCode, claimedPhrase
}

// DecodeError represents an error related to unmarshalling.
type DecodeError struct {
	error
}

func decodeError(err error) *DecodeError {
	return &DecodeError{error: fmt.Errorf("Unable to parse remote response as XML: %s", err.Error())}
}

// ReadFullyXML reads a http response into a data container using an XML decoder.
// The response is checked for its status code and the http.Response.Body is closed.
func ReadFullyXML(response *http.Response, errorObtainResponse error, v interface{}) error {
	return readDecode(response, errorObtainResponse, func(r io.Reader, v interface{}) error {
		return xml.NewDecoder(r).Decode(v)
	}, v)
}

// ReadFullyJSON reads a http response into a data container using a json decoder.
// The response is checked for its status code and the http.Response.Body is closed.
func ReadFullyJSON(response *http.Response, errorObtainResponse error, v interface{}) error {
	return readDecode(response, errorObtainResponse, func(r io.Reader, v interface{}) error {
		return json.NewDecoder(r).Decode(v)
	}, v)
}

func readDecode(response *http.Response, errorObtainResponse error, decode func(r io.Reader, v interface{}) error, v interface{}) error {
	if errorObtainResponse != nil {
		return errorObtainResponse
	}
	defer response.Body.Close()
	errDecode := decode(response.Body, v)
	if response.StatusCode >= 400 {
		return statusCodeError(response.StatusCode, response.Status)
	}
	if errDecode != nil {
		return decodeError(errDecode)
	}
	return nil
}
