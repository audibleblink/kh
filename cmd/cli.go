package cli

import (
	"bufio"
	"os"

	"github.com/audibleblink/kh/pkg/keyhack"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kh",
	Short: "Validate API tokens/webhooks for various services",
}

func Execute() {
	rootCmd.Execute()
}

func newCommand(name, desc string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: desc,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			var (
				ok  bool
				err error
			)

			// accept multiple lines through stdin
			if args[0] == "-" {
				scanner := bufio.NewScanner(os.Stdin)
				for scanner.Scan() {
					token := scanner.Text()
					ok, err = keyhack.Check(name, token)
					if err != nil {
						cmd.PrintErr(err)
						continue
					}
					if ok {
						cmd.Println(token)
					}
				}
				return
			}

			// allow multiple args so commands like xargs also work
			if len(args) > 1 {
				for _, token := range args {
					ok, err = keyhack.Check(name, token)
					if err != nil {
						cmd.PrintErr(err)
						continue
					}
					if ok {
						cmd.Println(token)
					}
				}
				return
			}

			token := args[0]
			ok, err = keyhack.Check(name, token)
			if err != nil {
				cmd.PrintErr(err)
				os.Exit(1)
			}
			if ok {
				cmd.Println(token)
				return
			}
			os.Exit(1)
		},
	}
}
