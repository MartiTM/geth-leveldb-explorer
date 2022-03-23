/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"main/levelDbTools"

	"github.com/spf13/cobra"
)

// storageTreeSizeCmd represents the storageTreeSize command
var storageTreeSizeCmd = &cobra.Command{
	Use:   "storageTreeSize <LevelDB path>",
	Short: "",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if(len(args) == 0 ) {
			args = append(args, "./.ethereum/geth/chaindata")
		}
		levelDbTools.GetStorageSize(args[0])
	},
}

func init() {
	rootCmd.AddCommand(storageTreeSizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// storageTreeSizeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// storageTreeSizeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}