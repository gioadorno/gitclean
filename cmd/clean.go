package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Rebase and clean",
	Long:  "Rebases pointer branch and resets all commits to one commit and pushes up to HEAD",
	RunE:  clean,
}

func clean(cmd *cobra.Command, args []string) error {
	err := rebase(cmd, args)
	if err != nil {
		fmt.Printf("There was an error attempting to rebase: %v", err)
		os.Exit(1)
	}

	resetErr := reset(cmd, args)
	if resetErr != nil {
		fmt.Printf("There was an error attempting to reset branch: %v", err)
		os.Exit(1)
	}
	fmt.Println("Clean successful.")
	return nil
}
