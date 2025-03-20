package keyhack

import (
	"context"
	"fmt"
	"maps"
	"net/http"
	"strings"
	"time"
)

// Registry provides access to the service registry
var Registry struct {
	GetService func(name string) (*KeyHack, bool)
}

// Check validates a token against the specified service
func Check(serviceName, token string) (bool, error) {
	service, exists := Registry.GetService(serviceName)
	if !exists {
		return false, fmt.Errorf("service %q not configured", serviceName)
	}

	ok, err := service.Validate(token)
	if err != nil {
		return false, fmt.Errorf("validation for service %q failed: %w", serviceName, err)
	}

	return ok, nil
}

// ValidatorFunc is a function that users define which establishes what a valid
// authenticated HTTP response looks like from a given service
type ValidatorFunc func(*http.Response) (bool, error)

// Validator holds validation configuration and logic for a service
type Validator struct {
	Custom bool
	Status int
	Fn     ValidatorFunc
}

// Request contains the necessary parameters for validating a token
type Request struct {
	Method  string
	URL     string
	Headers map[string]string
}

// KeyHack represents an API service definition from the config YAML
type KeyHack struct {
	Name string
	Request
	Validator
	Custom bool
}

// Validate sends an HTTP request with the given token and validates the response
func (kh *KeyHack) Validate(token string) (bool, error) {
	// Fill in the token template
	req := kh.prepareRequest(token)

	// Send the request
	res, err := kh.sendRequest(req)
	if err != nil {
		return false, fmt.Errorf("validation request failed: %w", err)
	}
	defer res.Body.Close()

	// Use default validator if none provided
	if !kh.Custom {
		kh.Fn = defaultValidator
	}

	// Run the validator
	ok, err := kh.Fn(res)
	if err != nil {
		return false, fmt.Errorf("validator function failed: %w", err)
	}

	return ok, nil
}

// defaultValidator checks for HTTP 200 OK status code
func defaultValidator(resp *http.Response) (bool, error) {
	return resp.StatusCode == 200, nil
}

// prepareRequest creates a Request with the token inserted in templates
func (kh *KeyHack) prepareRequest(token string) *Request {
	newReq := &Request{
		Method:  kh.Method,
		URL:     kh.URL,
		Headers: make(map[string]string, len(kh.Headers)),
	}

	// Copy headers to avoid modifying original
	maps.Copy(newReq.Headers, kh.Headers)

	// Fill URL template
	if strings.Contains(kh.URL, "%s") {
		newReq.URL = fmt.Sprintf(kh.URL, token)
	}

	// Fill header templates
	for k, v := range kh.Headers {
		if strings.Contains(v, "%s") {
			newReq.Headers[k] = fmt.Sprintf(v, token)
		}
	}

	return newReq
}

// sendRequest performs the HTTP request and returns the response
func (kh *KeyHack) sendRequest(req *Request) (*http.Response, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build the HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, req.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for k, v := range req.Headers {
		httpReq.Header.Add(k, v)
	}

	// Send the request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	return res, nil
}
