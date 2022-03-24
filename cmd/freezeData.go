/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (

	"github.com/spf13/cobra"
	"main/tools"
	"strconv"
)

// freezeDataCmd represents the freezeData command
var freezeDataCmd = &cobra.Command{
	Use:   "freezeData",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		blockNumber, _ := strconv.ParseUint(args[1], 10, 64)
		tools.GetBlockData(args[0], blockNumber)
	},
}

func init() {
	rootCmd.AddCommand(freezeDataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// freezeDataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// freezeDataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
