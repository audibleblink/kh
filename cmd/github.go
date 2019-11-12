package cli

func init() {
	githubTokenCmd := newCommand("github-token", "<token>")
	rootCmd.AddCommand(githubTokenCmd)
}
