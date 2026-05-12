package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	serverURL string
)

var rootCmd = &cobra.Command{
	Use:   "mangahub",
	Short: "MangaHub CLI - Manage your manga library from the terminal",
	Long: `MangaHub CLI is a standard interface for interacting with the MangaHub server.
It allows users to search manga, manage their reading library, and chat with others.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mangahub/config.json)")
	rootCmd.PersistentFlags().StringVar(&serverURL, "server", "http://localhost:8080", "MangaHub server URL")
}

func initConfig() {
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cfgFile = filepath.Join(home, ".mangahub", "config.json")
	}
}
