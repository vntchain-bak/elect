package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "elect",
	Short: "Election tools of VNT Chain",
	Long: `This tool is used to take part in VNT hubble network election.
You can do any operation of election contract, more information see help command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Using `elect -h` to see how to use elect.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Set flags of elect command
	// Add sub commands
	rootCmd.AddCommand(stakeCmd, unStakeCmd)
}
