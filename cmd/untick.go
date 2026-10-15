/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

var untickCmd = &cobra.Command{
	Use:   "untick [name] ...",
	Short: "Undo ticking today's date for the given task(s)",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		var yamlTasks, _ = Unmarshal()
		var query = make(map[string]struct{})
		for _, arg := range args {
			query[arg] = struct{}{}
		}
		Untick(Query(query, yamlTasks))
		Marshal(yamlTasks)
	},
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// Return dynamic completions based on user's data
		_, validArgs := Unmarshal()
		return validArgs, cobra.ShellCompDirectiveNoFileComp
	},
}

func init() {
	rootCmd.AddCommand(untickCmd)
}
