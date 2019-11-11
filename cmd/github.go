package cli

import (
	"net/http"

	"github.com/audibleblink/kh/pkg/keyhack"
)

// each subcommand's init function must add the subcommand to the root command
// and then add the validator function to the keyhack registry so that it know
// what a good http response looks like, and thus reports successful authentication
func init() {
	rootCmd.AddCommand(githubTokenCmd)
	keyhack.Registry["github-token"].Validator = validateGithub
}

// ensure the command name matches the entry in the YAML file
var githubTokenCmd = newCommand("github-token", "Checks a token against the GitHub API")

// define what a successful authentication means based on the
// of the API call issued by keyhacks
func validateGithub(resp *http.Response) (ok bool, err error) {
	ok = resp.StatusCode == 200
	return
}
