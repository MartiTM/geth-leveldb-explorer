package cmd

import (

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/spf13/cobra"
	"main/tools"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect <Chaindata path>",
	Short: "Same function as geth inspect",
	Long: ``,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		db, _ := rawdb.NewLevelDBDatabaseWithFreezer(path, 0, 0, path+"/ancient", "", false)
		tools.InspectDatabase(db, nil, nil)
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectDBCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inspectDBCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
