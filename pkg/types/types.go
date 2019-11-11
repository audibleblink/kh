package types

import (
	"fmt"
	"net/http"
	"strings"
)

// Validator is a function that users define which establishes what a valid
// authenticated HTTP response looks like from a given service
type Validator func(*http.Response) (bool, error)

// KeyHack represents an API service definition that's read in from the config YAML
type KeyHack struct {
	Name      string
	Request   Request
	Validator Validator
}

// Request contains the necessary paramters for which KeyHacks use to
// validate a given API Token/Webhook
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
}

// Validate will take the configured properties and use them to send a request to
// the service whose token is attempting to be validated
func (kh *KeyHack) Validate(token string) (ok bool, err error) {
	req := fillTemplate(&kh.Request, token)
	res, err := curl(req)
	if err != nil {
		return
	}
	ok, err = kh.Validator(res)
	if err != nil {
		return
	}
	return
}

// curl is responsible for creating and sending the HTTP Request based on the parsed
// YAML block for a given KeyHack
func curl(req *Request) (res *http.Response, err error) {
	newReq, err := http.NewRequest(req.Method, req.URL, strings.NewReader(""))
	if err != nil {
		return
	}

	for k, v := range req.Headers {
		newReq.Header.Add(k, v)
	}

	client := &http.Client{}
	res, err = client.Do(newReq)
	return
}

// fillTemplate checks if string from the YAML configuration contains a format string
// and fills it with a token if it does
func fillTemplate(req *Request, token string) *Request {
	if strings.Contains(req.URL, "%s") {
		req.URL = fmt.Sprintf(req.URL, token)
	}
	for k, v := range req.Headers {
		if strings.Contains(req.Headers[k], "%s") {
			req.Headers[k] = fmt.Sprintf(v, token)
		}
	}
	return req
}
