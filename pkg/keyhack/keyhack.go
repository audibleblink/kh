package keyhack

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
)

func init() {
	data, err := ioutil.ReadFile("./keyhacks.yml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, Registry)
	if err != nil {
		panic(err)
	}
	registerCommands()
}

var Registry = make(map[string]*KeyHack)

type validator func(*http.Response) (bool, error)

type KeyHack struct {
	Name      string
	Request   request
	Validator validator
}

type request struct {
	Method  string
	URL     string
	Headers map[string]string
}

// Check is the main package function to which a user can pass both the service name
// and the token they wish to validate
func Check(service, token string) (ok bool, err error) {
	kh := Registry[service]
	if kh == nil {
		err = fmt.Errorf("Subcommand %s not configured", service)
		return
	}

	ok, err = kh.Validate(token)
	if err != nil {
		return
	}
	return
}

// Validate will take the configured properties and use them to send a request to
// the service whose token is attempting to be validated
func (kh *KeyHack) Validate(token string) (ok bool, err error) {
	log.Printf("Validating: %s", kh.Name)
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
func curl(req *request) (res *http.Response, err error) {
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
func fillTemplate(req *request, token string) *request {
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