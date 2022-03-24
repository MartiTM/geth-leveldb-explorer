/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"main/levelDBTools"

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
		levelDBTools.GetStorageTreeSize(args[0])
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
		levelDBTools.CountStorageTree(args[0])
	},
}

func init() {
	rootCmd.AddCommand(storageTreeSizeCmd)
	rootCmd.AddCommand(storageTreeCountCmd)
}