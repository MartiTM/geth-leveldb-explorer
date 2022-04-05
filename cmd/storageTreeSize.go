/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"main/tools"

	"github.com/spf13/cobra"
)

// storageTreeSizeCmd represents the storageTreeSize command
var storageTreeSizeCmd = &cobra.Command{
	Use:   "storageTreeSize <LevelDB path>",
	Short: "Displays the size of the storage trees of the last state tree present in levelDB",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if(len(args) == 0 ) {
			args = append(args, "./.ethereum/geth/chaindata")
		}
		tools.LatestStateTreeSize(args[0])
	},
}

// storageTreeCountCmd represents the storageTreeSize command
var storageTreeCountCmd = &cobra.Command{
	Use:   "storageTreeCount <LevelDB path>",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if(len(args) == 0 ) {
			args = append(args, "./.ethereum/geth/chaindata")
		}
		tools.CountingStorageTrees(args[0])
	},
}

// 
var snapshotReadCmd = &cobra.Command{
	Use:   "snapshotRead",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		tools.ReadSnapshot()
	},
}

func init() {
	rootCmd.AddCommand(storageTreeSizeCmd)
	rootCmd.AddCommand(storageTreeCountCmd)
	rootCmd.AddCommand(snapshotReadCmd)
}