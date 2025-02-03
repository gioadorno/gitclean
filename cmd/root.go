package cmd

import (
	"fmt"
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
		RunE:  rebase,
	}
	rebaseCmd.Flags().StringP("branch", "b", "origin/master", "Branch to rebase")
	rootCmd.AddCommand(rebaseCmd)
	return rootCmd.Execute()
}

// TODO: Move to a different file and extract each command/error to their own function for cleaner code
func rebase(cmd *cobra.Command, args []string) error {
	branch, _ := cmd.Flags().GetString("branch")

	// Check if there are changes in Master
	fmt.Printf("Fetching latest changes from origin...\n")
	_, err := exec.Command("git", "fetch", branch).Output()
	if err != nil {
		return err
	}
	fmt.Println("Fetch complete.")

	fmt.Printf("Checking for changes on %s...\n", branch)
	// 2. Check for changes on the target branch
	logCmd := exec.Command("git", "log", fmt.Sprintf("HEAD..%s", branch))
	logOutput, err := logCmd.Output()
	if err != nil {
		// Handle the case where git log fails (e.g., branch doesn't exist)
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

	fmt.Printf("Starting rebase onto %s...\n", branch)
	// 3. Perform the rebase (if changes exist)
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

	fmt.Println("Resetting to prepare for single commit...")
	// 4. Reset soft and check for single commit
	resetCmd := exec.Command("git", "reset", "--soft", "HEAD~1") // Reset to the commit before the rebase
	resetOutput, err := resetCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Reset failed:\n%s\n", string(resetOutput))
		return fmt.Errorf("reset failed: %w", err)
	}
	fmt.Printf("Reset successful:\n%s\n", string(resetOutput))

	fmt.Println("Force pushing changes...")
	// 5. Force push if no conflicts
	pushCmd := exec.Command("git", "push", "--force-with-lease", "origin", "HEAD")
	pushOutput, err := pushCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Force push failed:\n%s\n", string(pushOutput))
		return fmt.Errorf("force push failed: %w", err)
	}

	fmt.Printf("Force push successful:\n%s\n", string(pushOutput))

	return nil
}
