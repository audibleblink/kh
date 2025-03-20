package cli

import (
	"slices"
	"testing"
)

func TestNewServiceCommand(t *testing.T) {
	// Save the size of serviceCommands before the test
	initialSize := len(serviceCommands)

	// Call NewServiceCommand to create a new command
	cmd := NewServiceCommand("test-service", "<token>")

	// Verify command properties
	if cmd.Use != "test-service <token>" {
		t.Errorf("Command Use = %q, want %q", cmd.Use, "test-service <token>")
	}

	if cmd.Short != "<token>" {
		t.Errorf("Command Short = %q, want %q", cmd.Short, "<token>")
	}

	// Verify the command was registered
	if len(serviceCommands) != initialSize+1 {
		t.Errorf("serviceCommands size = %d, want %d", len(serviceCommands), initialSize+1)
	}

	// Verify the command is in the serviceCommands slice
	if !slices.Contains(serviceCommands, cmd) {
		t.Errorf("The created command was not found in serviceCommands")
	}
}
