package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"text/tabwriter"
	"os"

	"github.com/spf13/cobra"
	"github.com/user/mangahub/pkg/models"
)

func init() {
	rootCmd.AddCommand(mangaCmd)
	mangaCmd.AddCommand(searchMangaCmd)
	mangaCmd.AddCommand(infoMangaCmd)
}

var mangaCmd = &cobra.Command{
	Use:   "manga",
	Short: "Manga discovery commands",
}

var searchMangaCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for manga by title, author, or genre",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		q := url.QueryEscape(args[0])
		token := getToken()

		req, _ := http.NewRequest("GET", serverURL+"/api/manga?q="+q, nil)
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: Server unreachable")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error: Unauthorized or server error")
			return
		}

		var results []*models.Manga
		json.NewDecoder(resp.Body).Decode(&results)

		if len(results) == 0 {
			fmt.Println("No manga found matching your query.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tAUTHOR\tGENRES\tSTATUS")
		for _, m := range results {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", m.ID, m.Title, m.Author, m.Genres, m.Status)
		}
		w.Flush()
	},
}

var infoMangaCmd = &cobra.Command{
	Use:   "info [id]",
	Short: "Get detailed information about a manga",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		token := getToken()

		req, _ := http.NewRequest("GET", serverURL+"/api/manga/"+id, nil) // We need a single manga endpoint
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error: Server unreachable")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error: Manga not found")
			return
		}

		var m models.Manga
		json.NewDecoder(resp.Body).Decode(&m)

		fmt.Printf("--- %s ---\n", m.Title)
		fmt.Printf("ID: %d\n", m.ID)
		fmt.Printf("Author: %s\n", m.Author)
		fmt.Printf("Genres: %s\n", m.Genres)
		fmt.Printf("Status: %s\n", m.Status)
		fmt.Printf("Total Chapters: %d\n", m.TotalChapters)
		fmt.Printf("\nDescription:\n%s\n", m.Description)
	},
}
