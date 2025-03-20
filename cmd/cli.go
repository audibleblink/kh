package cli

import (
	"bufio"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/audibleblink/kh/pkg/keyhack"
)

var rootCmd = &cobra.Command{
	Use:   "kh",
	Short: "Validate API tokens/webhooks for various services",
}

// serviceCommands tracks all registered service commands
var serviceCommands []*cobra.Command

// Execute runs the CLI application
func Execute() {
	// Add all registered service commands to the root command
	for _, cmd := range serviceCommands {
		rootCmd.AddCommand(cmd)
	}

	// Let Cobra handle command execution and errors
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// NewServiceCommand creates a new cobra command for a service
func NewServiceCommand(name, desc string) *cobra.Command {
	usage := strings.Join([]string{name, desc}, " ")
	cmd := &cobra.Command{
		DisableFlagsInUseLine: true,
		DisableFlagParsing:    true,
		Use:                   usage,
		Short:                 desc,
		Args:                  cobra.MinimumNArgs(1),
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

	// Register the command in the registry
	serviceCommands = append(serviceCommands, cmd)

	return cmd
}
