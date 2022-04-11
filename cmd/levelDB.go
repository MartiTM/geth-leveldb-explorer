package cmd

import (
	"main/tools"

	"github.com/spf13/cobra"
)

var trieDetailsCmd = &cobra.Command{
	Use:   "trieDetails <LevelDB path>",
	Short: "Search in levelDB the merkle-patricia trees and detail the last one",
	Long: `Search in levelDB the merkle-patricia trees and detail the last one.
	Returns:

   - Total number of state trees (for blocks present in levelDB).
   - Gives the block number and the root of the most recent state tree
   - Total number of accounts (including smartcontract) in the tree
   - Total number of smartcontract in the tree
   - Size of the most recent state tree with leaf details
   - Size of most recent storage tree with leaf details
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tools.StateAndStorageTrees(args[0])
	},
}

var countStateTreesCmd = &cobra.Command{
	Use:   "countStateTrees <LevelDB path>",
	Short: "Count in levelDB the merkle-patricia trees.",
	Long: `Count in levelDB the merkle-patricia trees.
	Return the total number of state trees (for blocks present in levelDB).
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tools.CountStateTrees(args[0])
	},
}

var snapshotAccountCmd = &cobra.Command{
	Use:   "snapshotAccount <LevelDB path> <account address>",
	Short: "Search for an account in the snapshot part of LevelDB",
	Long: `Search for an account in the snapshot part of LevelDB.
	Return raw and decoded informations about the account`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tools.SnapshotAccount(args[0], args[1])
	},
}

var treeAccountCmd = &cobra.Command{
	Use:   "treeAccount <LevelDB path> <account address>",
	Short: "Search for an account in the merkle-patricia tree part of LevelDB",
	Long: `Return raw and decoded informations about the account`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tools.TreeAccount(args[0], args[1])
	},
}

var compareAccountCmd = &cobra.Command{
	Use:   "compareAccount <LevelDB path> <account address>",
	Short: "Search for an account in the merkle-patricia tree and snapshot in LevelDB",
	Long: `Search for an account in the merkle-patricia tree and snapshot in LevelDB.
	Return raw and decoded informations about the account for both part`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tools.CompareAccount(args[0], args[1])
	},
}

var readCmd = &cobra.Command{
	Use:   "read <LevelDB path> <key in hex>",
	Short: "",
	Long: ``,
	// Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tools.Read(args[0])
	},
}

func init() {
	rootCmd.AddCommand(trieDetailsCmd)
	rootCmd.AddCommand(countStateTreesCmd)
	rootCmd.AddCommand(snapshotAccountCmd)
	rootCmd.AddCommand(treeAccountCmd)
	rootCmd.AddCommand(compareAccountCmd)
	rootCmd.AddCommand(readCmd)
}