package keyhack

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// mockHTTP is a helper to create mock HTTP responses
type mockHTTP struct {
	server *httptest.Server
}

// setupMockHTTP creates a test HTTP server with the given response
func setupMockHTTP(statusCode int, body string, headers map[string]string) *mockHTTP {
	mock := &mockHTTP{}
	
	// Create a test server with a handler that returns the specified status and body
	mock.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set response headers
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		// Set status code
		w.WriteHeader(statusCode)
		// Write body
		w.Write([]byte(body))
	}))
	
	return mock
}

// Close shuts down the mock HTTP server
func (m *mockHTTP) Close() {
	m.server.Close()
}

// URL returns the mock server's URL
func (m *mockHTTP) URL() string {
	return m.server.URL
}

func TestPrepareRequest(t *testing.T) {
	testCases := []struct {
		name         string
		keyHack      KeyHack
		token        string
		wantURL      string
		wantHeaders  map[string]string
		wantMethod   string
	}{
		{
			name: "URL Template",
			keyHack: KeyHack{
				Name: "test",
				Request: Request{
					Method: "GET",
					URL:    "https://api.example.com/%s/info",
					Headers: map[string]string{
						"Accept": "application/json",
					},
				},
			},
			token:       "abc123",
			wantURL:     "https://api.example.com/abc123/info",
			wantHeaders: map[string]string{"Accept": "application/json"},
			wantMethod:  "GET",
		},
		{
			name: "Header Template",
			keyHack: KeyHack{
				Name: "test",
				Request: Request{
					Method: "POST",
					URL:    "https://api.example.com/info",
					Headers: map[string]string{
						"Authorization": "Bearer %s",
						"Accept":        "application/json",
					},
				},
			},
			token:      "abc123",
			wantURL:    "https://api.example.com/info",
			wantHeaders: map[string]string{
				"Authorization": "Bearer abc123",
				"Accept":        "application/json",
			},
			wantMethod: "POST",
		},
		{
			name: "Both Templates",
			keyHack: KeyHack{
				Name: "test",
				Request: Request{
					Method: "PUT",
					URL:    "https://api.example.com/%s/info",
					Headers: map[string]string{
						"Authorization": "Bearer %s",
						"Accept":        "application/json",
					},
				},
			},
			token:      "abc123",
			wantURL:    "https://api.example.com/abc123/info",
			wantHeaders: map[string]string{
				"Authorization": "Bearer abc123",
				"Accept":        "application/json",
			},
			wantMethod: "PUT",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := tc.keyHack.prepareRequest(tc.token)
			
			// Check URL
			if req.URL != tc.wantURL {
				t.Errorf("URL = %q, want %q", req.URL, tc.wantURL)
			}
			
			// Check method
			if req.Method != tc.wantMethod {
				t.Errorf("Method = %q, want %q", req.Method, tc.wantMethod)
			}
			
			// Check headers
			if len(req.Headers) != len(tc.wantHeaders) {
				t.Errorf("Headers count = %d, want %d", len(req.Headers), len(tc.wantHeaders))
			}
			
			for k, v := range tc.wantHeaders {
				if req.Headers[k] != v {
					t.Errorf("Header[%q] = %q, want %q", k, req.Headers[k], v)
				}
			}
		})
	}
}

func TestDefaultValidator(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		want       bool
	}{
		{
			name:       "Success",
			statusCode: 200,
			want:       true,
		},
		{
			name:       "Failure",
			statusCode: 401,
			want:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp := &http.Response{
				StatusCode: tc.statusCode,
				Body:       io.NopCloser(strings.NewReader("")),
			}
			
			got, err := defaultValidator(resp)
			if err != nil {
				t.Errorf("defaultValidator() error = %v", err)
			}
			
			if got != tc.want {
				t.Errorf("defaultValidator() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestCheck(t *testing.T) {
	// Set up a service for testing
	service := &KeyHack{
		Name: "test",
		Request: Request{
			Method: "GET",
			URL:    "https://api.example.com/test",
		},
		Validator: Validator{
			Custom: true,
			Fn: func(resp *http.Response) (bool, error) {
				return resp.StatusCode == 200, nil
			},
		},
	}

	// Save the original Registry.GetService
	originalGetService := Registry.GetService
	
	// Set up a temporary mock registry
	mockRegistry := func(name string) (*KeyHack, bool) {
		if name == "test" {
			return service, true
		}
		return nil, false
	}
	
	// Install the mock
	Registry.GetService = mockRegistry
	
	// Restore the original registry when we're done
	defer func() {
		Registry.GetService = originalGetService
	}()

	// Set up a mock HTTP server
	mock := setupMockHTTP(200, `{"success": true}`, nil)
	defer mock.Close()
	
	// Update the service URL to point to our mock server
	service.URL = mock.URL()

	// Test successful Check
	got, err := Check("test", "dummy-token")
	if err != nil {
		t.Errorf("Check() error = %v", err)
	}
	if !got {
		t.Errorf("Check() = %v, want true", got)
	}

	// Test non-existent service
	_, err = Check("nonexistent", "dummy-token")
	if err == nil {
		t.Error("Check() with nonexistent service should error")
	}
}

func TestValidate(t *testing.T) {
	// Set up a mock HTTP server
	mockSuccess := setupMockHTTP(200, `{"success": true}`, nil)
	defer mockSuccess.Close()
	
	mockFailure := setupMockHTTP(401, `{"error": "unauthorized"}`, nil)
	defer mockFailure.Close()

	testCases := []struct {
		name    string
		keyHack KeyHack
		token   string
		want    bool
		wantErr bool
	}{
		{
			name: "Success - Default Validator",
			keyHack: KeyHack{
				Name: "test",
				Request: Request{
					Method: "GET",
					URL:    mockSuccess.URL(),
				},
				Validator: Validator{
					Custom: false,
				},
			},
			token:   "dummy-token",
			want:    true,
			wantErr: false,
		},
		{
			name: "Failure - Default Validator",
			keyHack: KeyHack{
				Name: "test",
				Request: Request{
					Method: "GET",
					URL:    mockFailure.URL(),
				},
				Validator: Validator{
					Custom: false,
				},
			},
			token:   "dummy-token",
			want:    false,
			wantErr: false,
		},
		{
			name: "Success - Custom Validator",
			keyHack: KeyHack{
				Name: "test",
				Request: Request{
					Method: "GET",
					URL:    mockSuccess.URL(),
				},
				Validator: Validator{
					Custom: true,
					Fn: func(resp *http.Response) (bool, error) {
						return true, nil
					},
				},
			},
			token:   "dummy-token",
			want:    true,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.keyHack.Validate(tc.token)
			
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			
			if got != tc.want {
				t.Errorf("Validate() = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestInvalidURL ensures sendRequest handles invalid URLs properly
func TestInvalidURL(t *testing.T) {
	kh := &KeyHack{
		Name: "test",
		Request: Request{
			Method: "GET",
			URL:    "://invalid-url", // Invalid URL format
		},
	}
	
	req := &Request{
		Method: "GET",
		URL:    "://invalid-url",
	}
	
	_, err := kh.sendRequest(req)
	if err == nil {
		t.Error("sendRequest() with invalid URL should return error")
	}
}