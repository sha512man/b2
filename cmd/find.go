package cmd

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"time"
)

type Query struct {
	Content    string `yaml:"content"`
	Target     string `yaml:"target"`
	LastUsedAt string `yaml:"last_used_at"`
}

type SearchHistory struct {
	Queries map[string]Query `yaml:"queries"`
}

var findCmd = &cobra.Command{
	Use:   "find",
	Short: "Find an entry in your second brain.",
	Long:  `Takes query as an input and uses LLM to find the best match.`,
	Run: func(cmd *cobra.Command, args []string) {
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

		var query string

		huh.NewText().
			Title("What are you looking for?").
			Value(&query).
			Run()

		hash := md5.Sum([]byte(query))
		hashString := hex.EncodeToString(hash[:])

		historyFilePath := filepath.Join(b2Dir, ".history.yaml")
		historyContent, err := os.ReadFile(historyFilePath)

		if err != nil {
			fmt.Printf("Failed to read .history.yaml file: %s\n", err.Error())
			return
		}

		var history SearchHistory

		if err := yaml.Unmarshal(historyContent, &history); err != nil {
			fmt.Printf("Failed to deserialize search history: %s\n", err.Error())
			return
		}

		now := time.Now().Format("2006-01-02 15:04:05")

		if value, exists := history.Queries[hashString]; exists {
			fmt.Println(value)
			value.LastUsedAt = now
		} else {
			fmt.Println("No match found for the query.")
			fmt.Println("Using AI...")

			// Now we need to perform AI search through the b2 entries.
			// Offer the user a list of possible matches and once they select the dir.
			// Store the query and target dir.
			// Also add the query into 'matching_queries' field to the metadata.yaml file in the target dir.

			targetDir := "foo"

			history.Queries[hashString] = Query{
				Content:    query,
				Target:     targetDir,
				LastUsedAt: now,
			}
		}

		file, err := os.OpenFile(historyFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)

		if err != nil {
			fmt.Printf("Failed to open history file for writing: %s\n", err.Error())
			return
		}

		encoder := yaml.NewEncoder(file)
		defer encoder.Close()

		if err := encoder.Encode(history); err != nil {
			fmt.Printf("Failed to write history: %s\n", err.Error())
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}
