package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
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

		for {
			fmt.Print("Please resolve your conflicts. Press enter to continue or q to abort process: ")

			oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
			if err != nil {
				panic(err)
			}
			defer term.Restore(int(os.Stdin.Fd()), oldState)

			var input []byte = make([]byte, 1)
			_, err = os.Stdin.Read(input)
			if err != nil {
				panic(err)
			}

			switch input[0] {
			case '\r', '\n': // Enter key
				fmt.Println("Continuing rebase...")
				// continueCmd := exec.Command("git", "rebase", "--continue")
				// continueOutput, err := continueCmd.CombinedOutput()
				// 	fmt.Printf("Failed to continue rebase: %v\nOutput: %s\n", err, continueOutput)
				// }
				continue
			case 'q', 'Q':
				fmt.Println()
				fmt.Println("Aborting rebase...")
				abortCmd := exec.Command("git", "rebase", "--abort")
				abortOutput, err := abortCmd.CombinedOutput()
				if err != nil {
					return fmt.Errorf("failed to abort rebase: %w", err)
				}
				fmt.Printf("Rebase abort output: %s\n", abortOutput)
				return nil
			default:
				fmt.Print("Invalid input. Please press Enter or q: ")
			}
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

	// test
	return nil
}
