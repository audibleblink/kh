package cli

func init() {
	twitterBearerCmd := newCommand("twitter-bearer", "<token>")
	rootCmd.AddCommand(twitterBearerCmd)
}
