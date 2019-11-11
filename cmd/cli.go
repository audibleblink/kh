package cli

import (
	"log"
	"os"

	"github.com/audibleblink/kh/pkg/keyhack"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kh",
	Short: "KeyHacks validates tokens for services",
	// Run: func(cmd *cobra.Command, args []string) {
	// 	// Do Stuff Here
	// },
}

func Execute() {
	rootCmd.Execute()
}

func newCommand(name, desc string) *cobra.Command {

	return &cobra.Command{
		Use:   name,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {

			ok, err := keyhack.Check(name, args[0])
			if err != nil {
				log.Fatal(err)
			}

			if ok {
				cmd.Println(args[0])
			} else {
				os.Exit(1)
			}
		},
	}
}
