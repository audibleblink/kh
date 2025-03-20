package services

import (
	"net/http"

	cli "github.com/audibleblink/kh/cmd"
	"github.com/audibleblink/kh/pkg/registry"
)

// Service configuration
const (
	SlackSubCmd = "slack-token"
	SlackToken  = "<token>"
)

// init registers the Slack service
func init() {
	// Create the command
	_ = cli.NewServiceCommand(SlackSubCmd, SlackToken)

	// Register the validator directly
	_ = registry.RegisterValidator(SlackSubCmd, validateSlack)
}

// validateSlack defines what a successful authentication looks like
// based on the HTTP response of the API call
func validateSlack(resp *http.Response) (ok bool, err error) {
	ok = resp.Header["X-Oauth-Scopes"] != nil
	return
}
