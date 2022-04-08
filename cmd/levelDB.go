package cmd

import (
	"main/tools"

	"github.com/spf13/cobra"
)

var trieDetailsCmd = &cobra.Command{
	Use:   "trieDetails <LevelDB path>",
	Short: "Display informations on state and storage trees of levelDB",
	Long: ``,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tools.StateAndStorageTrees(args[0])
	},
}

var countStateTreesCmd = &cobra.Command{
	Use:   "countStateTrees <LevelDB path>",
	Short: "Count the state trees in LevelDB",
	Long: ``,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		tools.CountStateTrees(args[0])
	},
}

var snapshotAccountCmd = &cobra.Command{
	Use:   "snapshotAccount <LevelDB path> <account address>",
	Short: "Display the account informations store in the snapshot part of LevelDB",
	Long: ``,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tools.ReadSnapshot(args[0], args[1])
	},
}

var treeAccountCmd = &cobra.Command{
	Use:   "treeAccount <LevelDB path> <account address>",
	Short: "Display the account informations store in the merkle tree part of LevelDB",
	Long: ``,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tools.TreeAccount(args[0], args[1])
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
	rootCmd.AddCommand(readCmd)
}