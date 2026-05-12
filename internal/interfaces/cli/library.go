package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"text/tabwriter"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/mangahub/pkg/models"
)

func init() {
	rootCmd.AddCommand(libraryCmd)
	libraryCmd.AddCommand(libraryListCmd)
	libraryCmd.AddCommand(libraryAddCmd)
	libraryCmd.AddCommand(libraryUpdateCmd)

	libraryAddCmd.Flags().String("status", "reading", "Reading status (reading, completed, plan_to_read)")
	libraryUpdateCmd.Flags().Int("chapter", 0, "Current chapter being read")
	libraryUpdateCmd.Flags().String("status", "", "Update status (optional)")
}

var libraryCmd = &cobra.Command{
	Use:   "library",
	Short: "Manage your personal manga library",
}

var libraryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all manga in your library",
	Run: func(cmd *cobra.Command, args []string) {
		token := getToken()

		req, _ := http.NewRequest("GET", serverURL+"/api/manga/library", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: Server unreachable")
			return
		}
		defer resp.Body.Close()

		var results []*models.UserProgress
		json.NewDecoder(resp.Body).Decode(&results)

		if len(results) == 0 {
			fmt.Println("Your library is empty. Use 'mangahub library add <id>' to track manga.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "MANGA ID\tCHAPTER\tSTATUS\tUPDATED AT")
		for _, p := range results {
			fmt.Fprintf(w, "%d\t%d\t%s\t%s\n", p.MangaID, p.CurrentChapter, p.Status, p.UpdatedAt.Format("2006-01-02 15:04"))
		}
		w.Flush()
	},
}

var libraryAddCmd = &cobra.Command{
	Use:   "add [manga_id]",
	Short: "Add a manga to your library",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mID, _ := strconv.Atoi(args[0])
		status, _ := cmd.Flags().GetString("status")
		token := getToken()

		payload := map[string]interface{}{
			"manga_id":        mID,
			"current_chapter": 0,
			"status":          status,
		}
		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", serverURL+"/api/manga/progress", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: Server unreachable")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error: Failed to add manga to library")
			return
		}

		fmt.Printf("Manga #%d successfully added to your library with status: %s\n", mID, status)
	},
}

var libraryUpdateCmd = &cobra.Command{
	Use:   "update [manga_id]",
	Short: "Update progress for a manga in your library",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		mID, _ := strconv.Atoi(args[0])
		chapter, _ := cmd.Flags().GetInt("chapter")
		status, _ := cmd.Flags().GetString("status")
		token := getToken()

		// First, get current progress to preserve status if not provided
		reqGet, _ := http.NewRequest("GET", serverURL+"/api/manga/library", nil)
		reqGet.Header.Set("Authorization", "Bearer "+token)
		client := &http.Client{}
		respGet, _ := client.Do(reqGet)
		var current []*models.UserProgress
		json.NewDecoder(respGet.Body).Decode(&current)
		respGet.Body.Close()

		if status == "" {
			for _, p := range current {
				if p.MangaID == mID {
					status = p.Status
					break
				}
			}
			if status == "" { status = "reading" }
		}

		payload := map[string]interface{}{
			"manga_id":        mID,
			"current_chapter": chapter,
			"status":          status,
		}
		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", serverURL+"/api/manga/progress", bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: Server unreachable")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error: Failed to update progress")
			return
		}

		fmt.Printf("Manga #%d progress updated to Chapter %d (%s)\n", mID, chapter, status)
	},
}
