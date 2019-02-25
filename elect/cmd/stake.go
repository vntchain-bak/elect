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
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("stake transaction send success, tx hash: %s\n", txhash.String())
		}
	},
}

var unStakeCmd = &cobra.Command{
	Use:     "unstake",
	Short:   "Unstake vnt token",
	Long:    "Unstake provides pre-check before create a unstake transaction, and sends the tx if it may success.",
	Example: "elect unstake",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmd.Help()
			return
		}

		e, err := elect.NewElection()
		if err != nil {
			panic(err)
		}
		if txhash, err := e.Unstake(); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("unstake transaction send success, tx hash: %s\n", txhash.String())
		}
	},
}

var registerWitnessCmd = &cobra.Command{
	Use:     "register",
	Short:   "Register witness candidate",
	Long:    "Register provides pre-check before create a register witness transaction, and sends the tx if it may success.",
	Example: "elect register nodeName nodeUrl website",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			cmd.Help()
			return
		}

		e, err := elect.NewElection()
		if err != nil {
			panic(err)
		}
		if txhash, err := e.RegisterWitness(args[0], args[1], args[2]); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("register witness transaction send success, tx hash: %s\n", txhash.String())
		}
	},
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query election data",
	Long: `Query supports getting the stake or vote information of account in config, 
getting witness candidates list and rest bounty.`, // TODO 修改描述
	Example: `elect query stake/vote/candidates/rest`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}

		var (
			ret []byte
			err error
			e   *elect.Election
		)

		e, err = elect.NewElection()
		if err != nil {
			panic(err)
		}

		switch args[0] {
		case "stake":
			ret, err = e.QueryStake()
		case "vote":
			ret, err = e.QueryVote()
		case "candidates":
			ret, err = e.QueryCandidates()
		case "rest":
			ret, err = e.QueryRestVNTBounty()
		default:
			fmt.Printf("error: query not support %s\n", args[0])
			fmt.Printf("\nQuery help:\n")
			cmd.Help()
			return
		}

		if err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("Result:\n%s\n", string(ret))
		}
	},
}
