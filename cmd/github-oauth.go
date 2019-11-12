package cli

func init() {
	githubOAuthCmd := newCommand("github-oauth", "<client_id:client_secret>")
	rootCmd.AddCommand(githubOAuthCmd)
}
