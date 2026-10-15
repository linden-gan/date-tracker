/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [name] ...",
	Short: "Show tasks' date tracker in  . . # . # . format",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var yamlTasks, _ = Unmarshal()
		var query = make(map[string]struct{})
		for _, arg := range args {
			query[arg] = struct{}{}
		}
		urgency, _ := cmd.Flags().GetBool("urgency")
		frequency, _ := cmd.Flags().GetBool("frequency")
		if urgency && frequency {
			fmt.Println("Warning: cannot sort by both urgency and frequency. Urgency will be used.")
		}
		Print(Query(query, yamlTasks), urgency, frequency)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Return dynamic completions based on user's data
		_, allStrings := Unmarshal()
		return allStrings, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolP("urgency", "u", false, "Sort output tasks by urgency (countdown)")
	showCmd.Flags().BoolP("frequency", "f", false, "Sort output tasks by frequency")
}
