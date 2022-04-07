package cmd

import (
	"github.com/spf13/cobra"
	"main/tools"
	"strconv"
)

var freezeDataCmd = &cobra.Command{
	Use:   "freezeBlock <FreezeDB path> <block number>",
	Short: "Read all the data from a block on FreezeDB",
	Long: ``,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		blockNumber, _ := strconv.ParseUint(args[1], 10, 64)
		tools.FreezerBlockData(args[0], blockNumber)
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
