package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vntchain/elect"
)

var stakeCmd = &cobra.Command{
	Use:   "stake vntCount",
	Short: "Stake vnt token",
	Long: `Stake provides checks before creating transaction to stake, 
and sends the transaction if it may execute success.`,
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

		e, err := elect.NewElection("./cofig.json")
		if err != nil {
			panic(err)
		}
		if txhash, err := e.Stake(stakeCnt); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("stake transaction send success, transaction hash: %s\n", txhash.String())
		}
	},
}

var unStakeCmd = &cobra.Command{
	Use:   "unstake",
	Short: "Unstake vnt token",
	Long: `Unstake provides checks before creating a transaction to unstake,
and sends the transaction if it may execute success.`,
	Example: "elect unstake",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmd.Help()
			return
		}

		e, err := elect.NewElection("./cofig.json")
		if err != nil {
			panic(err)
		}
		if txhash, err := e.Unstake(); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("unstake transaction send success, transaction hash: %s\n", txhash.String())
		}
	},
}

var registerWitnessCmd = &cobra.Command{
	Use:   "register",
	Short: "Register witness candidate",
	Long: `Register provides checks before creating a transaction to register
witness, and sends the transaction if it may execute success.`,
	Example: "elect register nodeName nodeUrl website",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 3 {
			cmd.Help()
			return
		}

		e, err := elect.NewElection("./config.json")
		if err != nil {
			panic(err)
		}
		if txhash, err := e.RegisterWitness(args[0], args[1], args[2]); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("register witness transaction send success, transaction hash: %s\n", txhash.String())
		}
	},
}

var unregisterWitnessCmd = &cobra.Command{
	Use:   "unregister",
	Short: "Unregister witness candidate",
	Long: `Unregister provides checks before creating a transaction to unregister
witness, and sends the transaction if it may execute success.`,
	Example: "elect unregister",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmd.Help()
			return
		}

		e, err := elect.NewElection("./config.json")
		if err != nil {
			panic(err)
		}
		if txhash, err := e.UnregisterWitness(); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("unregister witness transaction send success, transaction hash: %s\n", txhash.String())
		}
	},
}

var voteCmd = &cobra.Command{
	Use:   "vote",
	Short: "Vote witness candidate, up to 30 witnesses",
	Long: `Vote provides checks before creating a transaction to vote witness,
and sends the transaction if it may execute success.`,
	Example: `elect vote "0x123....456" "0x789...123"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			cmd.Help()
			return
		}

		e, err := elect.NewElection("./config.json")
		if err != nil {
			panic(err)
		}
		if txhash, err := e.VoteWitness(args); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("vote witness transaction send success, transaction hash: %s\n", txhash.String())
		}
	},
}

var cancelVoteCmd = &cobra.Command{
	Use:   "cancelVote",
	Short: "Cancel the vote for witness candidate",
	Long: `Cancel vote provides checks before creating a transaction to cancel
the vote, and sends the transaction if it may execute success.`,
	Example: `elect cancelVote`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmd.Help()
			return
		}

		e, err := elect.NewElection("./config.json")
		if err != nil {
			panic(err)
		}
		if txhash, err := e.CancelVote(); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("cancel vote witness transaction send success, transaction hash: %s\n", txhash.String())
		}
	},
}

var startProxyCmd = &cobra.Command{
	Use:   "startProxy",
	Short: "Become a vote proxy",
	Long: `Start proxy provides checks before creating a transaction to become
a proxy, and sends the transaction if it may execute success.`,
	Example: `elect startProxy`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmd.Help()
			return
		}

		e, err := elect.NewElection("./config.json")
		if err != nil {
			panic(err)
		}
		if txhash, err := e.StartProxy(); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("start proxy transaction send success, transaction hash: %s\n", txhash.String())
		}
	},
}

var stopProxyCmd = &cobra.Command{
	Use:   "stopProxy",
	Short: "Back to a normal voter",
	Long: `Stop proxy provides checks before creating transaction to become a 
normal voter, and sends the transaction if it may execute success.`,
	Example: `elect stopProxy`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmd.Help()
			return
		}

		e, err := elect.NewElection("./config.json")
		if err != nil {
			panic(err)
		}
		if txhash, err := e.StopProxy(); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("stop proxy transaction send success, transaction hash: %s\n", txhash.String())
		}
	},
}

var setProxyCmd = &cobra.Command{
	Use:   "setProxy",
	Short: "Vote by proxy ",
	Long: `Set proxy provides checks before creating a transaction to set a account
as your proxy, and sends the transaction if it may execute success.`,
	Example: `elect setProxy proxyAccountAddr`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Help()
			return
		}

		e, err := elect.NewElection("./config.json")
		if err != nil {
			panic(err)
		}
		if txhash, err := e.SetProxy(args[0]); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("stop proxy transaction send success, transaction hash: %s\n", txhash.String())
		}
	},
}

var cancelProxyCmd = &cobra.Command{
	Use:   "cancelProxy",
	Short: "Cancel vote by proxy ",
	Long: `Cancel proxy provides checks before creating a transaction to cancel
your proxy, and sends the transaction if it may execute success.`,
	Example: `elect cancelProxy`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmd.Help()
			return
		}

		e, err := elect.NewElection("./config.json")
		if err != nil {
			panic(err)
		}
		if txhash, err := e.CancelProxy(); err != nil {
			fmt.Printf("error: %s\n", err)
		} else {
			fmt.Printf("cancel proxy transaction send success, transaction hash: %s\n", txhash.String())
		}
	},
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query election data",
	Long: `Query supports getting the stake or vote information of account, 
and getting witness candidates list and rest bounty.`,
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

		e, err = elect.NewElection("./config.json")
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
