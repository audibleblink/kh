package registry

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/audibleblink/kh/pkg/keyhack"
)

// getRegistrySize returns the number of services in the registry
func getRegistrySize() int {
	return len(registry)
}

// clearRegistry resets the registry to empty state
func clearRegistry() {
	registry = make(ServiceRegistry)
}

func TestLoadFromBytes(t *testing.T) {
	// Clear the registry before test
	clearRegistry()

	testCases := []struct {
		name      string
		yamlData  []byte
		wantSize  int
		wantError bool
	}{
		{
			name: "Valid YAML",
			yamlData: []byte(`
github-token:
  name: github-token
  request:
    method: GET
    url: 'https://api.github.com/users'
    headers:
      Authorization: "token %s"
slack-token:
  name: slack-token
  request:
    method: POST
    url: 'https://slack.com/api/auth.test?token=%s&pretty=1'
  validator:
    custom: true
`),
			wantSize:  2,
			wantError: false,
		},
		{
			name:      "Empty YAML",
			yamlData:  []byte(""),
			wantSize:  0,
			wantError: false,
		},
		{
			name: "Invalid YAML",
			yamlData: []byte(`
invalid:
  - this
  is not: valid
  yaml: [
`),
			wantSize:  0,
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Clear registry before each case
			clearRegistry()

			// Test LoadFromBytes
			err := LoadFromBytes(tc.yamlData)
			
			// Check error
			if (err != nil) != tc.wantError {
				t.Errorf("LoadFromBytes() error = %v, wantError %v", err, tc.wantError)
				return
			}
			
			// Check registry size
			gotSize := getRegistrySize()
			if gotSize != tc.wantSize {
				t.Errorf("LoadFromBytes() registry size = %d, want %d", gotSize, tc.wantSize)
			}
		})
	}
}

func TestGetService(t *testing.T) {
	// Clear registry and set up test data
	clearRegistry()
	
	// Add a test service to the registry
	registry["test-service"] = &keyhack.KeyHack{
		Name: "test-service",
		Request: keyhack.Request{
			Method: "GET",
			URL:    "https://example.com/api",
		},
	}

	testCases := []struct {
		name       string
		serviceName string
		wantExists bool
	}{
		{
			name:        "Service Exists",
			serviceName: "test-service",
			wantExists:  true,
		},
		{
			name:        "Service Doesn't Exist",
			serviceName: "nonexistent",
			wantExists:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service, exists := GetService(tc.serviceName)
			
			if exists != tc.wantExists {
				t.Errorf("GetService() exists = %v, want %v", exists, tc.wantExists)
			}
			
			if tc.wantExists && service == nil {
				t.Errorf("GetService() returned nil service but exists = true")
			}
			
			if tc.wantExists && service.Name != tc.serviceName {
				t.Errorf("GetService() service.Name = %s, want %s", service.Name, tc.serviceName)
			}
		})
	}
}

func TestRegisterValidator(t *testing.T) {
	// Clear registry and set up test data
	clearRegistry()
	
	// Add a test service to the registry
	registry["test-service"] = &keyhack.KeyHack{
		Name: "test-service",
		Request: keyhack.Request{
			Method: "GET",
			URL:    "https://example.com/api",
		},
	}

	// Create a validator function
	validatorFunc := func(resp *http.Response) (bool, error) {
		return resp.StatusCode == 200, nil
	}

	testCases := []struct {
		name        string
		serviceName string
		validator   keyhack.ValidatorFunc
		wantError   bool
	}{
		{
			name:        "Valid Service",
			serviceName: "test-service",
			validator:   validatorFunc,
			wantError:   false,
		},
		{
			name:        "Invalid Service",
			serviceName: "nonexistent",
			validator:   validatorFunc,
			wantError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := RegisterValidator(tc.serviceName, tc.validator)
			
			if (err != nil) != tc.wantError {
				t.Errorf("RegisterValidator() error = %v, wantError %v", err, tc.wantError)
				return
			}
			
			if !tc.wantError {
				// Verify validator was set
				service, exists := GetService(tc.serviceName)
				if !exists {
					t.Errorf("Service %s not found after registering validator", tc.serviceName)
					return
				}
				
				// Compare function pointers using reflection
				if reflect.ValueOf(service.Validator.Fn).Pointer() != reflect.ValueOf(tc.validator).Pointer() {
					t.Errorf("RegisterValidator() did not properly set the validator function")
				}
			}
		})
	}
}

func TestCompleteYAMLParsing(t *testing.T) {
	// This test ensures that all fields in the YAML are properly parsed
	clearRegistry()
	
	yamlData := []byte(`
github-token:
  name: github-token
  request:
    method: GET
    url: 'https://api.github.com/users'
    headers:
      Authorization: "token %s"
      Accept: "application/json"
  validator:
    custom: true
    status: 200
`)

	err := LoadFromBytes(yamlData)
	if err != nil {
		t.Fatalf("LoadFromBytes() failed with error: %v", err)
	}
	
	service, exists := GetService("github-token")
	if !exists {
		t.Fatalf("Service 'github-token' not found after loading YAML")
	}
	
	// Check service name
	if service.Name != "github-token" {
		t.Errorf("Service name = %s, want 'github-token'", service.Name)
	}
	
	// Check request method
	if service.Method != "GET" {
		t.Errorf("Request method = %s, want 'GET'", service.Method)
	}
	
	// Check request URL
	if service.URL != "https://api.github.com/users" {
		t.Errorf("Request URL = %s, want 'https://api.github.com/users'", service.URL)
	}
	
	// Check headers
	if len(service.Headers) != 2 {
		t.Errorf("Request headers count = %d, want 2", len(service.Headers))
	}
	
	if service.Headers["Authorization"] != "token %s" {
		t.Errorf("Authorization header = %s, want 'token %%s'", service.Headers["Authorization"])
	}
	
	if service.Headers["Accept"] != "application/json" {
		t.Errorf("Accept header = %s, want 'application/json'", service.Headers["Accept"])
	}
	
	// Check validator
	if !service.Validator.Custom {
		t.Errorf("Validator Custom = %v, want true", service.Validator.Custom)
	}
	
	if service.Status != 200 {
		t.Errorf("Validator Status = %d, want 200", service.Status)
	}
}