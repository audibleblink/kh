package registry

import (
	"fmt"

	"gopkg.in/yaml.v3"

	kh "github.com/audibleblink/kh/pkg/keyhack"
)

// ServiceRegistry holds configurations for service validation
type ServiceRegistry map[string]*kh.KeyHack

// registry is the global registry instance, accessed through accessor functions
var registry = make(ServiceRegistry)

// LoadFromBytes loads configuration from the provided bytes
func LoadFromBytes(configData []byte) error {
	if err := yaml.Unmarshal(configData, &registry); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	return nil
}

// GetService returns a service by name, or nil if not found
func GetService(name string) (*kh.KeyHack, bool) {
	service, exists := registry[name]
	return service, exists
}

// RegisterValidator registers a custom validator for a service
func RegisterValidator(serviceName string, validator kh.ValidatorFunc) error {
	service, exists := registry[serviceName]
	if !exists {
		return fmt.Errorf("service %q not configured", serviceName)
	}

	service.Validator.Fn = validator
	return nil
}
