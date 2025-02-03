package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var rebaseCmd = &cobra.Command{
	Use:   "rebase",
	Short: "Rebase with input branch",
	Long:  "Checks if input branch has any changes, then pulls those changes to HEAD, and rebase Master with HEAD. Outputs if there are any conflicts. If there are no conflicts then asks the user if they want to push to current branch.",
	RunE:  rebase,
}

func rebase(cmd *cobra.Command, args []string) error {
	branch, _ := cmd.Flags().GetString("branch")

	// 1. Fetch input branch or default origin/master
	fmt.Printf("Fetching latest changes from origin...\n")
	_, err := exec.Command("git", "fetch", branch).Output()
	if err != nil {
		return err
	}
	fmt.Println("Fetch complete.")

	// 2. Check for changes on the target branch
	fmt.Printf("Checking for changes on %s...\n", branch)
	logCmd := exec.Command("git", "log", fmt.Sprintf("HEAD..%s", branch))
	logOutput, err := logCmd.Output()
	if err != nil {
		if strings.Contains(err.Error(), "unknown revision") {
			return fmt.Errorf("branch %s not found", branch)
		}
		return fmt.Errorf("failed to check for changes: %w", err)
	}

	if len(logOutput) == 0 {
		fmt.Printf("There are no new changes on %s branch\n", branch)
		return nil
	}
	fmt.Printf("Changes found on %s:\n%s\n", branch, string(logOutput))

	// 3. Perform the rebase (if changes exist)
	fmt.Printf("Starting rebase onto %s...\n", branch)
	rebaseCmd := exec.Command("git", "rebase", branch)
	rebaseOutput, err := rebaseCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Rebase failed:\n%s\n", string(rebaseOutput))

		// Check for conflicts:
		diffCmd := exec.Command("git", "diff", "--name-only", "--conflicts")
		conflictFiles, err := diffCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to check for conflicts: %w", err)
		}

		if len(conflictFiles) > 0 {
			fmt.Printf("Conflicts found in the following files:\n%s\n", string(conflictFiles))
			return fmt.Errorf("rebase failed due to conflicts")
		} else {
			return fmt.Errorf("rebase failed (unknown reason)") // Rebase failed, but no conflicts detected.
		}

	}
	fmt.Printf("Rebase successful:\n%s\n", string(rebaseOutput))

	// 4. Force Push
	fmt.Println("Force pushing changes...")
	pushCmd := exec.Command("git", "push", "--force-with-lease", "origin", "HEAD")
	pushOutput, err := pushCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Force push failed:\n%s\n", string(pushOutput))
		return fmt.Errorf("force push failed: %w", err)
	}
	fmt.Printf("Force push successful:\n%s\n", string(pushOutput))

	return nil
}
