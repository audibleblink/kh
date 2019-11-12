package cli

func init() {
	twitterCmd := newCommand("twitter", "<token:secret>")
	rootCmd.AddCommand(twitterCmd)
}
