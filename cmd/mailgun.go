package cli

func init() {
	mailgunCmd := newCommand("mailgun", "Checks a token against the Mailgun API")
	rootCmd.AddCommand(mailgunCmd)
}
