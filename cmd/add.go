package cmd

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// TODO extract into separate file
type Metadata struct {
	Description     string   `yaml:"description"`
	CreatedAt       string   `yaml:"created_at"`
	MatchingQueries []string `yaml:"matching_queries"`
}

var addCmd = &cobra.Command{
	Use:   "add [dirname]",
	Short: "Add a new entry to your second brain.",
	Long: `Creates a new directory in $B2_DIR with metadata.yaml.

You will provide a human readable description of what does the entry contain.
It is later used during lookups.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a directory name.")
			return
		}

		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("$HOME/.config/b2/")

		if err := viper.ReadInConfig(); err != nil {
			fmt.Printf("Failed to read config file: %s\n", err.Error())
			return
		}

		b2Dir := viper.GetString("B2_DIR")

		if b2Dir == "" {
			homeDir, err := os.UserHomeDir()

			if err != nil {
				fmt.Printf("Failed to get user home directory: %s\n", err.Error())
				return
			}

			b2Dir = filepath.Join(homeDir, "b2")
		}

		newDir := filepath.Join(b2Dir, args[0])

		var description string

		huh.NewText().
			Title("Describe the entry.").
			Value(&description).
			Run()

		if err := os.MkdirAll(newDir, 0755); err != nil {
			fmt.Printf("Failed to create directory: %s\n", err.Error())
			return
		}

		metadata := Metadata{
			Description:     description,
			CreatedAt:       time.Now().Format("2006-01-02 15:04:05"),
			MatchingQueries: []string{description},
		}

		if file, err := os.Create(filepath.Join(newDir, "metadata.yaml")); err != nil {
			fmt.Printf("Failed to create metadata file: %s\n", err.Error())
			return
		} else {
			defer file.Close()

			encoder := yaml.NewEncoder(file)
			defer encoder.Close()

			if err := encoder.Encode(metadata); err != nil {
				fmt.Printf("Failed to write metadata: %s\n", err.Error())
				return
			}
		}

		cmdToRun := exec.Command("sh", "-c", "cd "+newDir+"; exec $SHELL")
		cmdToRun.Stdin = os.Stdin
		cmdToRun.Stdout = os.Stdout
		cmdToRun.Stderr = os.Stderr

		if err := cmdToRun.Run(); err != nil {
			fmt.Printf("Failed to change directory: %s\n", err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
