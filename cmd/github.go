package cli

func init() {
	githubTokenCmd := newCommand("github-token", "Checks a token against the GitHub API")
	rootCmd.AddCommand(githubTokenCmd)
}
