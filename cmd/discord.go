package cli

func init() {
	mailgunCmd := newCommand("discord", "<token>")
	rootCmd.AddCommand(mailgunCmd)
}
