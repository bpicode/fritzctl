package httpread

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// HTTPStatusCodeError represents an 4xx client or a 5xx server error.
type HTTPStatusCodeError struct {
	error
}

func statusCodeError(response *http.Response) *HTTPStatusCodeError {
	return &HTTPStatusCodeError{error: fmt.Errorf("HTTP status code error: remote replied with %s", response.Status)}
}

// ReadFullyString reads a http response into a string. The response is checked for its status code.
func ReadFullyString(response *http.Response, errorObtainResponse error) (string, error) {
	if errorObtainResponse != nil {
		return "", errorObtainResponse
	}
	defer response.Body.Close()
	body, errRead := ioutil.ReadAll(response.Body)
	if response.StatusCode >= 400 {
		return string(body), statusCodeError(response)
	}
	return string(body), errRead
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
		return statusCodeError(response)
	}
	if errDecode != nil {
		return xmlDecodeError(errDecode)
	}
	return nil
}
