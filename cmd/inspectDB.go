/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/spf13/cobra"
	"main/inspectDatabase"
)

// inspectDBCmd represents the inspectDB command
var inspectDBCmd = &cobra.Command{
	Use:   "inspectDB",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		path := "../.ethereum-testnet/goerli/geth/chaindata"
		// path := "../.ethereum-test/geth/chaindata"
		db, _ := rawdb.NewLevelDBDatabaseWithFreezer(path, 0, 0, path+"/ancient", "", false)
		inspectDatabase.InspectDatabase(db, nil, nil)
	},
}

func init() {
	rootCmd.AddCommand(inspectDBCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inspectDBCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inspectDBCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
