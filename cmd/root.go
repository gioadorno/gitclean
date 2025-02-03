package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitclean",
	Short: "Quick tool for rebasing, reset, and push to branch",
	Long:  "GitClean is a tool that allows the user to easily rebase and reset their branch to make it clean and easier to review.",
}

func Execute() error {
	rebaseCmd := &cobra.Command{
		Use:   "rebase",
		Short: "Rebase with Master",
		Long:  "Checks if Master branch has any changes, then pulls those changes to Master, and rebase Master with current branch. Outputs if there are any conflicts. If there are no conflicts then asks the user if they want to push to current branch.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if there are changes in Master
			_, err := exec.Command("git status --porcelain origin/master").Output()
			if err != nil {
				return err
			}
			return nil
		},
	}
	rootCmd.AddCommand(rebaseCmd)
	return rootCmd.Execute()
}
