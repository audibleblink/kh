package keyhack

import (
	"fmt"

	"github.com/audibleblink/kh/pkg/registry"
)

func init() {
	registry.Build()
}

// Check is the main package function to which a user can pass both the service name
// and the token they wish to validate
func Check(service, token string) (ok bool, err error) {
	kh := registry.Registry[service]
	if kh == nil {
		err = fmt.Errorf("Subcommand %s not configured", service)
		return
	}

	ok, err = kh.Validate(token)
	return
}
