package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Clean branch and squash to one commit",
	Long:  "Resets all commits to one commit and pushes up to HEAD",
	RunE:  reset,
}

func reset(cmd *cobra.Command, args []string) error {
	// 1. Find the Parent's last commit
	fmt.Println("Fetching parent's last commit...")
	firstCommitCmd := exec.Command("git", "rev-list", "--max-parents=0", "HEAD")
	firstCommitOutput, err := firstCommitCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get first commit: %w", err)
	}
	firstCommit := strings.TrimSpace(string(firstCommitOutput))
	fmt.Printf("%s", firstCommit)

	// 2. Reset commits
	resetCmd := exec.Command("git", "reset", "--soft", firstCommit)
	resetOutput, err := resetCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Reset failed:\n%s\n", string(resetOutput))
		return fmt.Errorf("reset failed: %w", err)
	}
	fmt.Println("Reset complete.")

	// 3. Prompt the user for a commit message
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter a commit message: ")
	commitMsg, _ := reader.ReadString('\n')
	commitMsg = strings.TrimSpace(commitMsg)

	// 4. Prompt the user to add file(s) or path
	fmt.Print("Git file(s) or path to add. Default is '.': ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)
	if filePath == "" {
		filePath = "."
	}
	exec.Command("git", "add", filePath).Output()

	// 5. Commit the changes
	commitCmd := exec.Command("git", "commit", "-a", "-m", commitMsg)
	commitOutput, err := commitCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Commit failed:\n%s\n", string(commitOutput))
		return fmt.Errorf("commit failed: %w", err)
	}
	fmt.Println("Commit complete.")

	// 6. Push to the HEAD branch
	fmt.Println("Pushing...")
	pushCmd := exec.Command("git", "push", "--force-with-lease", "origin", "HEAD")
	pushOutput, err := pushCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Push failed:\n%s\n", string(pushOutput))
		return fmt.Errorf("push failed: %w", err)
	}
	fmt.Println("Branch is cleaned and updated.")

	return nil
}
