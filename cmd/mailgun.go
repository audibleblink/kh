package cli

func init() {
	mailgunCmd := newCommand("mailgun", "<token>")
	rootCmd.AddCommand(mailgunCmd)
}
