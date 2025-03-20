package main

import (
	_ "embed"
	"fmt"
	"os"

	cli "github.com/audibleblink/kh/cmd"
	_ "github.com/audibleblink/kh/cmd/services"
	"github.com/audibleblink/kh/pkg/keyhack"
	"github.com/audibleblink/kh/pkg/registry"
)

//go:embed keyhacks.yml
var configData []byte

func init() {
	// Initialize the keyhack registry lookup
	keyhack.Registry.GetService = registry.GetService
}

func main() {
	if err := registry.LoadFromBytes(configData); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %s\n", err)
		os.Exit(1)
	}

	// Execute the CLI
	cli.Execute()
}
