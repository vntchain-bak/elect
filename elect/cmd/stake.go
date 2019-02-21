package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vntchain/elect"
)

var stakeCmd = &cobra.Command{
	Use:     "stake vntCount",
	Short:   "Stake vnt token",
	Long:    "Stake provides pre-check before create stake a transaction, and sends the tx if it may success.",
	Example: "elect stake 1",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}
		stakeCnt, err := strconv.Atoi(args[0])
		if err != nil {
			log.Printf("vntCount is invalid: %s, got error: %s\n", args[0], err)
			return
		}

		e, err := elect.NewElection()
		if err != nil {
			panic(err)
		}
		if txhash, err := e.Stake(stakeCnt); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("stake transaction send success, tx hash: %s\n", txhash.String())
		}
	},
}
