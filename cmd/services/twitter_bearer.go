package services

import (
	cli "github.com/audibleblink/kh/cmd"
)

// Service configuration
const (
	TwitterBearerSubCmd = "twitter-bearer"
	TwitterBearerToken = "<token>"
)

// init registers the Twitter Bearer service
func init() {
	// Create the command
	_ = cli.NewServiceCommand(TwitterBearerSubCmd, TwitterBearerToken)
	
	// No custom validator needed - uses default 200 OK check
}