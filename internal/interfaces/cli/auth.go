package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(registerCmd)

	registerCmd.Flags().String("role", "user", "Role of the user (user/admin)")
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
}

var loginCmd = &cobra.Command{
	Use:   "login [username] [password]",
	Short: "Login to MangaHub",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		password := args[1]

		payload := map[string]string{"username": username, "password": password}
		jsonData, _ := json.Marshal(payload)
		resp, err := http.Post(serverURL+"/api/auth/login", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error: Server unreachable at %s\n", serverURL)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error: Invalid credentials")
			return
		}

		var res struct{ Token string }
		json.NewDecoder(resp.Body).Decode(&res)

		saveToken(res.Token)
		fmt.Printf("Successfully logged in as %s\n", username)
	},
}

var registerCmd = &cobra.Command{
	Use:   "register [username] [password]",
	Short: "Register a new account",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		password := args[1]
		role, _ := cmd.Flags().GetString("role")

		payload := map[string]string{"username": username, "password": password, "role": role}
		jsonData, _ := json.Marshal(payload)
		resp, err := http.Post(serverURL+"/api/auth/register", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error: Server unreachable at %s\n", serverURL)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			fmt.Println("Error: Registration failed")
			return
		}

		fmt.Printf("Successfully registered user: %s\n", username)
	},
}

func saveToken(token string) {
	config := map[string]string{"token": token}
	data, _ := json.MarshalIndent(config, "", "  ")

	dir := filepath.Dir(cfgFile)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	os.WriteFile(cfgFile, data, 0644)
}

func getToken() string {
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		return ""
	}
	var config struct{ Token string }
	json.Unmarshal(data, &config)
	return config.Token
}
