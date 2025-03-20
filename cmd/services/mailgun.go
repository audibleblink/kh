package services

import (
	cli "github.com/audibleblink/kh/cmd"
)

// Service configuration
const (
	MailgunSubCmd = "mailgun"
	MailgunToken = "<token>"
)

// init registers the Mailgun service
func init() {
	// Create the command
	_ = cli.NewServiceCommand(MailgunSubCmd, MailgunToken)
	
	// No custom validator needed - uses default 200 OK check
}