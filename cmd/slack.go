package cli

import (
	"net/http"

	"github.com/audibleblink/kh/pkg/registry"
)

// each subcommand's init function must add the subcommand to the root cli command
// and then add the validator function to the keyhack registry so that it knows
// what a good http response looks like
func init() {
	rootCmd.AddCommand(slackTokenCmd)
	registry.Registry["slack-token"].Validator.Fn = validateSlack
}

// ensure the command name matches the entry in the YAML file
var slackTokenCmd = newCommand("slack-token", "<token>")

// validator functions define what a successful authentication means
// based on the http response of the API call issued by keyhacks
func validateSlack(resp *http.Response) (ok bool, err error) {
	ok = resp.Header["X-Oauth-Scopes"] != nil
	return
}
