package cmd

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// TODO Make the first positional argument obligatory. Name it dirname.
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a new entry to your second brain.",
	Long: `Creates a new directory in $B2_DIR with metadata.yaml.

You will provide a human readable description of what does the entry contain.
It is later used during lookups.`,
	Run: func(cmd *cobra.Command, args []string) {
		var description string

		huh.NewText().
			Title("Describe the entry.").
			Value(&description).
			Run()

		// TODO
		// Create new directory in user's home named $dirname (as the positional argument passwed to the command)
		// Create empty metadata.yaml file in that folder.
		// Serialize struct with two field to the metadata.yaml.
		//  one description (taken from the variable description),
		//  second createdAt with current date in YYYY-MM-DD HH:MM:SS format.
		// Change user's current directory to the newly created dir.

		fmt.Printf("Your description: %s\n", description)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
