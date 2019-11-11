package cli

import (
	"log"
	"net/http"

	"github.com/audibleblink/kh/pkg/keyhack"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(githubTokenCmd)
	keyhack.Registry["github-token"].Validator = validateGithub
}

var githubTokenCmd = &cobra.Command{
	Use:   "github-token",
	Short: "Checks a token against the GitHub API",
	Run: func(cmd *cobra.Command, args []string) {

		ok, err := keyhack.Check("github-token", args[0])
		if err != nil {
			log.Fatal(err)
		}

		cmd.Println(ok)
	},
}

func validateGithub(resp *http.Response) (ok bool, err error) {
	ok = resp.StatusCode == 200
	return
}
