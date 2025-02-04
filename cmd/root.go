package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitclean",
	Short: "Quick tool for rebasing, reset, and push to branch",
	Long:  "GitClean is a tool that allows the user to easily rebase and reset their branch to make it clean and easier to review.",
}

func Execute() error {
	rebaseCmd.Flags().StringP("branch", "b", "origin/master", "Branch to rebase")
	cleanCmd.Flags().StringP("branch", "b", "origin/master", "Branch to rebase")

	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(rebaseCmd)
	rootCmd.AddCommand(resetCmd)
	return rootCmd.Execute()
}