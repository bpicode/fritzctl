package httpread

import (
	"encoding/xml"
	"fmt"
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

// ReadFullyString reads a http response into a string. The response is checked for its status code.
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

// XMLDecodeError represents an error related to XML unmarshalling.
type XMLDecodeError struct {
	error
}

func xmlDecodeError(err error) *XMLDecodeError {
	return &XMLDecodeError{error: fmt.Errorf("Unable to parse remote response as XML: %s", err.Error())}
}

// ReadFullyXML reads a http response into a datat container using an XML decoder. The response is checked for its status code.
func ReadFullyXML(response *http.Response, errorObtainResponse error, v interface{}) error {
	if errorObtainResponse != nil {
		return errorObtainResponse
	}
	defer response.Body.Close()
	errDecode := xml.NewDecoder(response.Body).Decode(v)
	if response.StatusCode >= 400 {
		return statusCodeError(response.StatusCode, response.Status)
	}
	if errDecode != nil {
		return xmlDecodeError(errDecode)
	}
	return nil
}
