package services

import (
	cli "github.com/audibleblink/kh/cmd"
)

// Service configuration
const (
	TwitterSubCmd = "twitter"
	TwitterToken = "<token:secret>"
)

// init registers the Twitter service
func init() {
	// Create the command
	_ = cli.NewServiceCommand(TwitterSubCmd, TwitterToken)
	
	// No custom validator needed - uses default 200 OK check
}