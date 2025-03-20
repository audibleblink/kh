package services

import (
	cli "github.com/audibleblink/kh/cmd"
)

// Service configuration
const (
	DiscordSubCmd = "discord"
	DiscordToken = "<token>"
)

// init registers the Discord service
func init() {
	// Create the command
	_ = cli.NewServiceCommand(DiscordSubCmd, DiscordToken)
	
	// No custom validator needed - uses default 200 OK check
}