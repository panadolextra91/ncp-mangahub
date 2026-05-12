package tui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
	"net/url"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
	"github.com/user/mangahub/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"github.com/user/mangahub/internal/application"
	"github.com/user/mangahub/pkg/models"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		// Input handling based on page
		if m.ActivePage == PageLogin {
			if msg.String() == "enter" {
				m.Status = "⏳ Logging in..."
				return m, m.LoginCmd()
			}
			if msg.String() == "up" || msg.String() == "down" || msg.String() == "tab" {
				m.FocusIndex = 1 - m.FocusIndex
				return m, nil
			}

			var target *string
			if m.FocusIndex == 0 {
				target = &m.Username
			} else {
				target = &m.Password
			}

			if msg.String() == "space" {
				*target += " "
			} else if len(msg.String()) == 1 {
				*target += msg.String()
			} else if msg.String() == "backspace" && len(*target) > 0 {
				*target = (*target)[:len(*target)-1]
			}
		} else if m.ActivePage == PageChat {
			if msg.String() == "enter" && len(m.ChatInput) > 0 {
				cmd := m.SendChatCmd(m.ChatInput)
				m.ChatInput = ""
				return m, cmd
			}
			if msg.String() == "space" {
				m.ChatInput += " "
			} else if len(msg.String()) == 1 {
				m.ChatInput += msg.String()
			} else if msg.String() == "backspace" && len(m.ChatInput) > 0 {
				m.ChatInput = m.ChatInput[:len(m.ChatInput)-1]
			}
		} else if m.ActivePage == PageCreate {
			isExempt := strings.ToLower(m.Role) == "admin"
			if !isExempt && msg.String() != "tab" {
				return m, nil
			}
			if msg.String() == "enter" {
				return m, m.CreateMangaCmd()
			}
			if msg.String() == "up" {
				m.FocusIndex = (m.FocusIndex - 1 + 6) % 6
			} else if msg.String() == "down" {
				m.FocusIndex = (m.FocusIndex + 1) % 6
			}

			// Handle input for focused field
			var target *string
			switch m.FocusIndex {
			case 0: target = &m.MangaTitle
			case 1: target = &m.MangaAuthor
			case 2: target = &m.MangaGenres
			case 3: target = &m.MangaStatus
			case 4: target = &m.MangaChapters
			case 5: target = &m.MangaDesc
			}

			if target == nil { return m, nil }
			if msg.String() == "space" {
				*target += " "
			} else if len(msg.String()) == 1 {
				*target += msg.String()
			} else if msg.String() == "backspace" && len(*target) > 0 {
				*target = (*target)[:len(*target)-1]
			}
		} else if m.ActivePage == PageProgress {
			if msg.String() == "enter" {
				return m, m.UpdateProgressCmd()
			}
			if msg.String() == "up" {
				m.FocusIndex = (m.FocusIndex - 1 + 3) % 3
			} else if msg.String() == "down" {
				m.FocusIndex = (m.FocusIndex + 1) % 3
			}
			var target *string
			switch m.FocusIndex {
			case 0: target = &m.MangaIDInput
			case 1: target = &m.ChapterInput
			case 2: target = &m.StatusInput
			}

			if target == nil { return m, nil }
			if msg.String() == "space" {
				*target += " "
			} else if len(msg.String()) == 1 {
				*target += msg.String()
			} else if msg.String() == "backspace" && len(*target) > 0 {
				*target = (*target)[:len(*target)-1]
			}
		} else if m.ActivePage == PageSearch {
			if msg.String() == "enter" {
				m.SearchIndex = 0 // Reset scroll on new search
				return m, m.SearchMangaCmd()
			}
			if msg.String() == "right" || msg.String() == "enter" && len(m.SearchResults) > 0 {
				// Go to Detail view for selected manga
				selectedLine := m.SearchResults[m.SearchIndex]
				// Parse ID from line (starts with ID)
				idStr := strings.Fields(selectedLine)[0]
				id, _ := strconv.Atoi(idStr)
				m.ActivePage = PageDetail
				m.Status = fmt.Sprintf("Loading Manga #%d...", id)
				return m, m.GetMangaDetailCmd(id)
			}
			if msg.String() == "up" && m.SearchIndex > 0 {
				m.SearchIndex--
			} else if msg.String() == "down" && m.SearchIndex < len(m.SearchResults)-1 {
				m.SearchIndex++
			}

			if msg.String() == "space" {
				m.SearchInput += " "
			} else if len(msg.String()) == 1 {
				m.SearchInput += msg.String()
			} else if msg.String() == "backspace" && len(m.SearchInput) > 0 {
				m.SearchInput = m.SearchInput[:len(m.SearchInput)-1]
			}
		} else if m.ActivePage == PageDiscover {
			if msg.String() == "up" && m.CurrentQuoteIndex > 0 {
				m.CurrentQuoteIndex--
			} else if msg.String() == "down" && m.CurrentQuoteIndex < len(m.Quotes)-1 {
				m.CurrentQuoteIndex++
			}
		} else if m.ActivePage == PageDetail {
			if msg.String() == "a" && m.SelectedManga != nil {
				m.Status = "➕ Adding to Library..."
				return m, m.AddLibraryCmd(m.SelectedManga.ID)
			}
			if msg.String() == "esc" || msg.String() == "backspace" {
				m.ActivePage = PageSearch
				return m, nil
			}
		} else if m.ActivePage == PageLibrary {
			if msg.String() == "up" && m.LibraryIndex > 0 {
				m.LibraryIndex--
			} else if msg.String() == "down" && m.LibraryIndex < len(m.LibraryResults)-1 {
				m.LibraryIndex++
			}
			if msg.String() == "r" {
				return m, m.FetchLibraryCmd()
			}
		}

		// Global keys
		switch msg.String() {
		case "r":
			if m.ActivePage == PageDiscover {
				m.Status = "🔄 Refreshing Quotes..."
				return m, m.ScrapeQuotesCmd()
			}
			if m.ActivePage == PageLibrary {
				m.Status = "🔄 Refreshing Library..."
				return m, m.FetchLibraryCmd()
			}
		}

		// Tab switching
		switch msg.String() {
		case "tab":
			if m.ActivePage != PageLogin {
				// Cycle through 7 main tabs (1 to 7)
				// PageDetail (8) is excluded from tab cycling
				if m.ActivePage >= PageLibrary {
					m.ActivePage = PageChat
				} else {
					m.ActivePage++
				}
				m.FocusIndex = 0
				if m.ActivePage == PageLibrary {
					return m, m.FetchLibraryCmd()
				}
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case TickMsg:
		m.AsciiIndex = (m.AsciiIndex + 1) % len(AsciiIcons)
		return m, TickCommand()

	case LoginSuccessMsg:
		m.Token = msg.Token
		m.Role = strings.TrimSpace(msg.Role) // Ensure no hidden spaces
		m.ActivePage = PageChat
		m.Status = fmt.Sprintf("Logged in as %s 🌸", m.Username)
		log.Printf("🕵️ [DEBUG] Login Successful! User: %s, Role: '%s'", m.Username, m.Role)
		return m, tea.Batch(m.ListenWSCmd(), m.ListenEventsCmd(), m.ListenTCPCmd())

	case TCPSyncMsg:
		m.Status = "🔄 Sync: " + string(msg)
		if m.ActivePage == PageLibrary {
			return m, m.FetchLibraryCmd()
		}
		return m, nil

	case *websocket.Conn:
		m.WS = msg
		return m, m.ReceiveWSCmd()

	case ChatMsg:
		m.Messages = append(m.Messages, string(msg))
		if len(m.Messages) > 10 {
			m.Messages = m.Messages[1:]
		}
		return m, m.ReceiveWSCmd() // Listen for next message

	case EventStream:
		m.Stream = msg
		return m, m.ReceiveEventCmd(m.Stream)

	case EventMsg:
		m.Events = append(m.Events, string(msg))
		m.Status = "New Release! ✨"
		return m, m.ReceiveEventCmd(m.Stream)

	case SearchSuccessMsg:
		m.SearchResults = msg
		m.Status = fmt.Sprintf("Found %d results 📚", len(msg))
		return m, nil

	case []models.Quote:
		m.Quotes = msg
		m.Status = "Quotes Updated! ✨"
		return m, nil

	case *models.Manga:
		m.SelectedManga = msg
		m.Status = "Viewing: " + msg.Title
		return m, nil

	case []*models.UserProgress:
		m.LibraryResults = msg
		m.Status = fmt.Sprintf("Library Updated (%d items) 📚", len(m.LibraryResults))
		return m, nil

	case ErrorMsg:
		m.Error = string(msg)
	}

	return m, nil
}

func (m Model) View() string {
	if m.Width == 0 {
		return "Initializing Pink Hub..."
	}

	header := HeaderStyle.Render(" MANGAHUB PREMIUM HUB ")
	status := lipgloss.NewStyle().Foreground(PinkGlow).Render(" Status: " + m.Status)
	
	ascii := AsciiStyle.Render(AsciiIcons[m.AsciiIndex])
	
	var content string
	switch m.ActivePage {
	case PageLogin:
		f := func(i int, label, val string, isPassword bool) string {
			displayVal := val
			if isPassword {
				displayVal = strings.Repeat("*", len(val))
			}
			style := lipgloss.NewStyle().Foreground(Slate400)
			if m.FocusIndex == i {
				style = lipgloss.NewStyle().Foreground(PinkPastel).Bold(true)
				return style.Render(fmt.Sprintf("> %s: %s_", label, displayVal))
			}
			return style.Render(fmt.Sprintf("  %s: %s", label, displayVal))
		}
		content = fmt.Sprintf("\n  🔐 Admin Login Required\n\n%s\n%s\n\n  (Press Enter to login)", 
			f(0, "Username", m.Username, false),
			f(1, "Password", m.Password, true))
	case PageChat:
		chat := strings.Join(m.Messages, "\n")
		content = fmt.Sprintf("💬 CHAT ROOM\n\n%s\n\n> %s_", chat, m.ChatInput)
	case PageEvents:
		evts := strings.Join(m.Events, "\n")
		content = fmt.Sprintf("📡 SYSTEM EVENTS\n\n%s", evts)
	case PageCreate:
		isAdmin := strings.ToLower(m.Role) == "admin"
		if !isAdmin {
			content = "\n\n  🚫 ACCESS DENIED\n\n  You are not an admin!\n  Only administrators can create new manga.\n\n  (Please switch to another tab)"
		} else {
			f := func(i int, label, val string) string {
				style := lipgloss.NewStyle().Foreground(Slate400)
				if m.FocusIndex == i {
					style = lipgloss.NewStyle().Foreground(PinkPastel).Bold(true)
					return style.Render(fmt.Sprintf("> %s: %s_", label, val))
				}
				return style.Render(fmt.Sprintf("  %s: %s", label, val))
			}
			content = fmt.Sprintf("🏗️ CREATE NEW MANGA\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n%s\n\n(Press Enter to Broadcast)", 
				f(0, "Title   ", m.MangaTitle),
				f(1, "Author  ", m.MangaAuthor),
				f(2, "Genres  ", m.MangaGenres),
				f(3, "Status  ", m.MangaStatus),
				f(4, "Chapters", m.MangaChapters),
				f(5, "Desc    ", m.MangaDesc))
		}
	case PageProgress:
		f := func(i int, label, val string) string {
			style := lipgloss.NewStyle().Foreground(Slate400)
			if m.FocusIndex == i {
				style = lipgloss.NewStyle().Foreground(PinkPastel).Bold(true)
				return style.Render(fmt.Sprintf("> %s: %s_", label, val))
			}
			return style.Render(fmt.Sprintf("  %s: %s", label, val))
		}
		content = fmt.Sprintf("📖 UPDATE PROGRESS\n\n%s\n\n%s\n\n%s\n\n(Press Enter to Sync)", 
			f(0, "Manga ID", m.MangaIDInput),
			f(1, "Chapter ", m.ChapterInput),
			f(2, "Status  ", m.StatusInput))
	case PageSearch:
		windowSize := 12
		start := m.SearchIndex - windowSize/2
		if start < 0 { start = 0 }
		end := start + windowSize
		if end > len(m.SearchResults) {
			end = len(m.SearchResults)
			start = end - windowSize
			if start < 0 { start = 0 }
		}

		var displayedResults []string
		for i := start; i < end; i++ {
			row := m.SearchResults[i]
			if i == m.SearchIndex {
				displayedResults = append(displayedResults, lipgloss.NewStyle().Foreground(PinkPastel).Bold(true).Render("> " + row))
			} else {
				displayedResults = append(displayedResults, "  " + row)
			}
		}

		results := strings.Join(displayedResults, "\n")
		scrollInfo := ""
		if len(m.SearchResults) > 0 {
			scrollInfo = fmt.Sprintf("\n\n(Showing %d-%d of %d | Up/Down to Scroll)", start+1, end, len(m.SearchResults))
		}

		content = fmt.Sprintf("🔍 SEARCH MANGA\n\nQuery: %s_\n\n(Press Enter to Search)\n\nRESULTS:\n  ID\tTITLE\t\t\t\t\tCHAPS\tSTATUS\n%s%s", 
			m.SearchInput, results, scrollInfo)
	case PageDiscover:
		var lines []string
		for i, q := range m.Quotes {
			prefix := "  "
			if i == m.CurrentQuoteIndex {
				prefix = "> "
			}
			text := q.Text
			if len(text) > 60 { text = text[:57] + "..." }
			lines = append(lines, fmt.Sprintf("%s\"%s\" - %s", prefix, text, q.Author))
		}
		content = fmt.Sprintf("🌟 DISCOVER QUOTES (Scraped from httpbin/quotes)\n\n%s\n\n(Press 'r' to re-scrape | Up/Down to scroll)", strings.Join(lines, "\n"))
	case PageDetail:
		if m.SelectedManga == nil {
			content = "Loading details..."
		} else {
			m := m.SelectedManga
			content = fmt.Sprintf("📖 MANGA DETAILS\n\n"+
				"Title:    %s\n"+
				"Author:   %s\n"+
				"Status:   %s\n"+
				"Chapters: %d\n"+
				"Genres:   %s\n\n"+
				"Description:\n%s\n\n"+
				"(Press 'a' to Add to Library | Backspace to return)", 
				m.Title, m.Author, m.Status, m.TotalChapters, m.Genres, m.Description)
		}
	case PageLibrary:
		var lines []string
		for i, p := range m.LibraryResults {
			prefix := "  "
			if i == m.LibraryIndex {
				prefix = "> "
			}
			line := fmt.Sprintf("%sManga #%-4d | Chap %-4d | %s", prefix, p.MangaID, p.CurrentChapter, p.Status)
			lines = append(lines, line)
		}
		if len(lines) == 0 {
			lines = append(lines, "  Your library is empty. Go search and add some manga! 🌸")
		}
		content = fmt.Sprintf("📚 MY PERSONAL LIBRARY\n\n%s\n\n(Press 'r' to Refresh | Up/Down to Scroll)", strings.Join(lines, "\n"))
	}

	// Tabs
	tabs := []string{"CHAT [1]", "EVENTS [2]", "CREATE [3]", "PROGRESS [4]", "SEARCH [5]", "DISCOVER [6]", "LIBRARY [7]"}
	var renderedTabs []string
	for i, t := range tabs {
		if int(m.ActivePage) == i+1 {
			renderedTabs = append(renderedTabs, ActiveTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, TabStyle.Render(t))
		}
	}
	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	// Dashboard Enhancement: Add a quote below ASCII if available
	quoteView := ""
	if len(m.Quotes) > 0 {
		q := m.Quotes[time.Now().UnixNano()%int64(len(m.Quotes))]
		quoteView = lipgloss.NewStyle().
			Foreground(PinkPastel).
			Italic(true).
			Width(25).
			Align(lipgloss.Center).
			Render(fmt.Sprintf("\n\"%s\"\n— %s", q.Text, q.Author))
	}

	mainView := lipgloss.JoinVertical(lipgloss.Left, tabRow, BorderStyle.Width(m.Width-30).Height(m.Height-12).Render(content))
	
	fullView := lipgloss.JoinHorizontal(lipgloss.Center, 
		lipgloss.JoinVertical(lipgloss.Center, ascii, quoteView),
		mainView,
	)

	return BaseStyle.Render(lipgloss.JoinVertical(lipgloss.Center, header, fullView, status + fmt.Sprintf(" [%s]", m.Role)))
}

// --- Commands (Non-blocking as per Mẹ Architect's instruction) ---

func (m Model) LoginCmd() tea.Cmd {
	return func() tea.Msg {
		url := "http://127.0.0.1:8080/api/auth/login"
		body, _ := json.Marshal(map[string]string{"username": m.Username, "password": m.Password})
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Post(url, "application/json", bytes.NewBuffer(body))
		if err != nil {
			return ErrorMsg("❌ Server Unreachable")
		}
		defer resp.Body.Close()
		
		if resp.StatusCode == http.StatusUnauthorized {
			// User likely doesn't exist, try Register with default 'user' role
			regUrl := "http://127.0.0.1:8080/api/auth/register"
			regBody, _ := json.Marshal(map[string]string{"username": m.Username, "password": m.Password, "role": "user"})
			regResp, err := client.Post(regUrl, "application/json", bytes.NewBuffer(regBody))
			if err != nil || (regResp.StatusCode != http.StatusCreated && regResp.StatusCode != http.StatusConflict) {
				return ErrorMsg("❌ Registration Failed")
			}
			regResp.Body.Close()
			return m.LoginCmd()() // Retry login after registration
		}

		var res struct {
			Token string `json:"token"`
			User struct {
				Role string `json:"role"`
			} `json:"user"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return ErrorMsg("❌ Login Failed")
		}

		if res.Token == "" {
			return ErrorMsg("❌ Invalid Token")
		}
		return LoginSuccessMsg{Token: res.Token, Role: res.User.Role}
	}
}

func (m Model) ListenWSCmd() tea.Cmd {
	return func() tea.Msg {
		if m.Token == "" {
			return nil
		}
		u := fmt.Sprintf("ws://127.0.0.1:8080/api/chat?token=%s&manga_id=1", m.Token)
		conn, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			return ErrorMsg("WS Connection Failed")
		}
		return conn
	}
}

func (m Model) ReceiveWSCmd() tea.Cmd {
	return func() tea.Msg {
		if m.WS == nil {
			return nil
		}
		_, p, err := m.WS.ReadMessage()
		if err != nil {
			return ErrorMsg("WS Read Error")
		}
		var chat struct {
			SenderName string `json:"sender_name"`
			Content    string `json:"content"`
		}
		json.Unmarshal(p, &chat)
		return ChatMsg(fmt.Sprintf("%s: %s", chat.SenderName, chat.Content))
	}
}

func (m Model) SendChatCmd(content string) tea.Cmd {
	return func() tea.Msg {
		if m.WS == nil {
			return ErrorMsg("Not connected to Chat")
		}
		m.WS.WriteMessage(websocket.TextMessage, []byte(content))
		return nil
	}
}

// gRPC Event Streaming
type EventStream pb.MangaService_SubscribeEventsClient

func (m Model) ListenEventsCmd() tea.Cmd {
	return func() tea.Msg {
		conn, err := grpc.Dial("127.0.0.1:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return ErrorMsg("gRPC Connection Failed")
		}
		client := pb.NewMangaServiceClient(conn)
		ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+m.Token)
		stream, err := client.SubscribeEvents(ctx, &emptypb.Empty{})
		if err != nil {
			return ErrorMsg("gRPC Stream Failed")
		}
		return EventStream(stream)
	}
}

func (m Model) ReceiveEventCmd(stream EventStream) tea.Cmd {
	return func() tea.Msg {
		resp, err := stream.Recv()
		if err != nil {
			return ErrorMsg("Event Stream Error")
		}
		return EventMsg(fmt.Sprintf("[%s] %s", resp.Topic, resp.Message))
	}
}

func (m Model) CreateMangaCmd() tea.Cmd {
	return func() tea.Msg {
		url := "http://127.0.0.1:8080/api/manga"
		chaps, _ := strconv.Atoi(m.MangaChapters)
		payload := map[string]interface{}{
			"title":          m.MangaTitle,
			"author":         m.MangaAuthor,
			"genres":         m.MangaGenres,
			"status":         m.MangaStatus,
			"total_chapters": chaps,
			"description":    m.MangaDesc,
		}
		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+m.Token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusCreated {
			return ErrorMsg("Failed to create Manga")
		}
		return EventMsg("You created: " + m.MangaTitle)
	}
}

func (m Model) UpdateProgressCmd() tea.Cmd {
	return func() tea.Msg {
		url := "http://127.0.0.1:8080/api/manga/progress"
		mID, _ := strconv.Atoi(m.MangaIDInput)
		chap, _ := strconv.Atoi(m.ChapterInput)
		payload := map[string]interface{}{
			"manga_id":        mID,
			"current_chapter": chap,
			"status":          m.StatusInput,
		}
		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+m.Token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return ErrorMsg("Failed to update progress")
		}
		return EventMsg(fmt.Sprintf("Progress Updated: Manga #%d to Chap %d (%s)", mID, chap, m.StatusInput))
	}
}

type SearchSuccessMsg []string

func (m Model) SearchMangaCmd() tea.Cmd {
	return func() tea.Msg {
		q := url.QueryEscape(m.SearchInput)
		url := "http://127.0.0.1:8080/api/manga?q=" + q
		
		client := &http.Client{Timeout: 5 * time.Second}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+m.Token)
		
		resp, err := client.Do(req)
		if err != nil {
			return ErrorMsg("❌ Search Failed: " + err.Error())
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return ErrorMsg(fmt.Sprintf("❌ Search Error: %d", resp.StatusCode))
		}

		var results []struct {
			ID            int    `json:"id"`
			Title         string `json:"title"`
			TotalChapters int    `json:"total_chapters"`
			Status        string `json:"status"`
		}
		json.NewDecoder(resp.Body).Decode(&results)

		var lines []string
		for _, res := range results {
			title := res.Title
			if len(title) > 25 { title = title[:22] + "..." }
			// Use fixed width spaces instead of tabs for better TUI alignment
			line := fmt.Sprintf("%-4d %-25s %-6d %s", res.ID, title, res.TotalChapters, res.Status)
			lines = append(lines, line)
		}
		return SearchSuccessMsg(lines)
	}
}

func (m Model) ScrapeQuotesCmd() tea.Cmd {
	return func() tea.Msg {
		scraper := application.NewScraperService()
		quotes, err := scraper.FetchQuotes()
		if err != nil {
			return ErrorMsg("Scrape Failed: " + err.Error())
		}
		return quotes
	}
}

func (m Model) GetMangaDetailCmd(id int) tea.Cmd {
	return func() tea.Msg {
		url := fmt.Sprintf("http://127.0.0.1:8080/api/manga/%d", id)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+m.Token)
		
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return ErrorMsg("Failed to fetch detail")
		}
		defer resp.Body.Close()

		var manga models.Manga
		json.NewDecoder(resp.Body).Decode(&manga)
		return &manga
	}
}

func (m Model) AddLibraryCmd(mangaID int) tea.Cmd {
	return func() tea.Msg {
		url := "http://127.0.0.1:8080/api/manga/progress"
		payload := map[string]interface{}{
			"manga_id":        mangaID,
			"current_chapter": 0,
			"status":          "plan_to_read",
		}
		jsonData, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
		req.Header.Set("Authorization", "Bearer "+m.Token)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return ErrorMsg("Failed to add to library")
		}
		return EventMsg("Added to Library! ✨")
	}
}

func (m Model) FetchLibraryCmd() tea.Cmd {
	return func() tea.Msg {
		url := "http://127.0.0.1:8080/api/manga/library"
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+m.Token)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return ErrorMsg("Failed to fetch library")
		}
		defer resp.Body.Close()

		var results []*models.UserProgress
		json.NewDecoder(resp.Body).Decode(&results)
		return results
	}
}

func (m Model) ListenTCPCmd() tea.Cmd {
	return func() tea.Msg {
		conn, err := net.Dial("tcp", "127.0.0.1:9090")
		if err != nil {
			return ErrorMsg("TCP Sync Failed")
		}
		// In a real app we'd send an AUTH message here, but for demo we just listen
		return m.ReceiveTCPCmd(conn)()
	}
}

func (m Model) ReceiveTCPCmd(conn net.Conn) tea.Cmd {
	return func() tea.Msg {
		buf := make(chan tea.Msg)
		go func() {
			for {
				data := make([]byte, 1024)
				n, err := conn.Read(data)
				if err != nil {
					buf <- ErrorMsg("TCP Read Error")
					return
				}
				buf <- TCPSyncMsg(string(data[:n]))
			}
		}()
		return <-buf
	}
}
