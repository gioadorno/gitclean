package cmd

import (
	"fmt"
	"os/exec"
	"strings"

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
		RunE:  rebase,
	}
	rebaseCmd.Flags().StringP("branch", "b", "origin/master", "Branch to rebase")
	rootCmd.AddCommand(rebaseCmd)
	return rootCmd.Execute()
}

// TODO: Move to a different file and extract each command/error to their own function for cleaner code
func rebase(cmd *cobra.Command, args []string) error {
	branch, _ := cmd.Flags().GetString("branch")

	fmt.Printf("Starting rebase onto %s...\n", branch)
	rebaseCmd := exec.Command("git", "rebase", branch)
	rebaseOutput, err := rebaseCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Rebase failed:\n%s\n", string(rebaseOutput))

		fmt.Println("Checking for conflicts...")
		diffCmd := exec.Command("git", "diff", "--name-only", "--conflicts")
		conflictFiles, err := diffCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to check for conflicts: %w", err)
		}

		if len(conflictFiles) > 0 {
			fmt.Printf("Conflicts found in the following files:\n%s\n", string(conflictFiles))
			return fmt.Errorf("rebase failed due to conflicts")
		} else {
			return fmt.Errorf("rebase failed (unknown reason)")
		}

	}

	fmt.Println("Rebase successful.")

	fmt.Println("Resetting to prepare for single commit...")
	resetCmd := exec.Command("git", "reset", "--soft", "HEAD~1")
	resetOutput, err := resetCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Reset failed:\n%s\n", string(resetOutput))
		return fmt.Errorf("reset failed: %w", err)
	}
	fmt.Println("Reset complete.")

	fmt.Println("Checking for single commit...")
	commitCountCmd := exec.Command("git", "rev-list", "--count", "HEAD")
	commitCountOutput, err := commitCountCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to count commits: %w", err)
	}

	commitCount := strings.TrimSpace(string(commitCountOutput))
	if commitCount != "1" {
		fmt.Printf("Branch does not have a single commit after reset. Commit count: %s\n", commitCount)
		return fmt.Errorf("branch is not in a single commit state. Please squash your commits before pushing")
	}
	fmt.Println("Single commit check passed.")

	fmt.Println("Force pushing changes...")
	pushCmd := exec.Command("git", "push", "--force-with-lease", "origin", "HEAD")
	pushOutput, err := pushCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Force push failed:\n%s\n", string(pushOutput))
		return fmt.Errorf("force push failed: %w", err)
	}

	fmt.Println("Force push successful.")

	return nil
}
