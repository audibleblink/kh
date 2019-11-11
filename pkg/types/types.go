package types

import (
	"fmt"
	"net/http"
	"strings"
)

// Validator is a function that users define which establishes what a valid
// authenticated HTTP response looks like from a given service
type validator func(*http.Response) (bool, error)

type Validator struct {
	Custom bool
	Status int
	Fn     validator
}

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

	if !kh.Validator.Custom {
		kh.Validator.Fn = defaultValidator
	}

	ok, err = kh.Validator.Fn(res)
	return
}

func defaultValidator(resp *http.Response) (ok bool, err error) {
	ok = resp.StatusCode == 200
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
	newReq := &Request{
		Method:  req.Method,
		URL:     req.URL,
		Headers: req.Headers,
	}

	if strings.Contains(req.URL, "%s") {
		newReq.URL = fmt.Sprintf(req.URL, token)
	}
	for k, v := range req.Headers {
		if strings.Contains(req.Headers[k], "%s") {
			newReq.Headers[k] = fmt.Sprintf(v, token)
		}
	}
	return newReq
}
