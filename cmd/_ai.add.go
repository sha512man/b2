Implement the missing code described in TODO comments

package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

type Metadata struct {
	Description string `yaml:"description"`
	CreatedAt   string `yaml:"createdAt"`
}

var addCmd = &cobra.Command{
	Use:   "add [dirname]",
	Short: "Adds a new entry to your second brain.",
	Long: `Creates a new directory in $B2_DIR with metadata.yaml.

You will provide a human readable description of what does the entry contain.
It is later used during lookups.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var description string

		huh.NewText().
			Title("Describe the entry.").
			Value(&description).
			Run()

		dirname := args[0]
		usr, _ := user.Current()
		homeDir := usr.HomeDir

		entryDir := filepath.Join(homeDir, dirname)
		err := os.MkdirAll(entryDir, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return
		}

		metadata := Metadata{
			Description: description,
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		}

		file, err := os.Create(filepath.Join(entryDir, "metadata.yaml"))
		if err != nil {
			fmt.Println("Error creating metadata file:", err)
			return
		}
		defer file.Close()

		encoder := yaml.NewEncoder(file)
		err = encoder.Encode(&metadata)
		if err != nil {
			fmt.Println("Error writing to metadata file:", err)
			return
		}
		encoder.Close()

		os.Chdir(entryDir)

		fmt.Printf("Your description: %s\n", description)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}

