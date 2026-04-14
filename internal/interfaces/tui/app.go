package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.ActivePage != PageLogin {
				m.ActivePage = (m.ActivePage % 3) + 1 // Cycle Chat -> Events -> Create
			}
		}

		// Input handling based on page
		if m.ActivePage == PageLogin {
			if msg.String() == "enter" {
				return m, m.LoginCmd()
			}
			// Simplified input (for demo, just username for now)
			if len(msg.String()) == 1 {
				m.Username += msg.String()
			} else if msg.String() == "backspace" && len(m.Username) > 0 {
				m.Username = m.Username[:len(m.Username)-1]
			}
		} else if m.ActivePage == PageChat {
			if msg.String() == "enter" && len(m.ChatInput) > 0 {
				cmd := m.SendChatCmd(m.ChatInput)
				m.ChatInput = ""
				return m, cmd
			}
			if len(msg.String()) == 1 {
				m.ChatInput += msg.String()
			} else if msg.String() == "backspace" && len(m.ChatInput) > 0 {
				m.ChatInput = m.ChatInput[:len(m.ChatInput)-1]
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height

	case TickMsg:
		m.AsciiIndex = (m.AsciiIndex + 1) % len(AsciiIcons)
		return m, TickCommand()

	case LoginSuccessMsg:
		m.Token = string(msg)
		m.ActivePage = PageChat
		m.Status = "Logged in as Admin 🌸"
		return m, m.ListenWSCmd()

	case ChatMsg:
		m.Messages = append(m.Messages, string(msg))
		if len(m.Messages) > 10 {
			m.Messages = m.Messages[1:]
		}

	case EventMsg:
		m.Events = append(m.Events, string(msg))
		m.Status = "New Release! ✨"

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
		content = fmt.Sprintf("\n  🔐 Admin Login Required\n\n  Username: %s_\n\n  (Press Enter to login)", m.Username)
	case PageChat:
		chat := strings.Join(m.Messages, "\n")
		content = fmt.Sprintf("💬 CHAT ROOM\n\n%s\n\n> %s_", chat, m.ChatInput)
	case PageEvents:
		evts := strings.Join(m.Events, "\n")
		content = fmt.Sprintf("📡 SYSTEM EVENTS\n\n%s", evts)
	case PageCreate:
		content = "🏗️ CREATE NEW MANGA\n\n(Coming Soon in TUI...)"
	}

	// Tabs
	tabs := []string{"CHAT [1]", "EVENTS [2]", "CREATE [3]"}
	var renderedTabs []string
	for i, t := range tabs {
		if int(m.ActivePage) == i+1 {
			renderedTabs = append(renderedTabs, ActiveTabStyle.Render(t))
		} else {
			renderedTabs = append(renderedTabs, TabStyle.Render(t))
		}
	}
	tabRow := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	mainView := lipgloss.JoinVertical(lipgloss.Left, tabRow, BorderStyle.Width(m.Width-30).Height(m.Height-12).Render(content))
	
	fullView := lipgloss.JoinHorizontal(lipgloss.Center, 
		lipgloss.NewStyle().Width(25).Render(ascii),
		mainView,
	)

	return BaseStyle.Render(lipgloss.JoinVertical(lipgloss.Center, header, fullView, status))
}

// --- Commands (Non-blocking as per Mẹ Architect's instruction) ---

func (m Model) LoginCmd() tea.Cmd {
	return func() tea.Msg {
		// Mock/Simple Login for demo (admin/password)
		url := "http://localhost:8080/api/auth/login"
		body, _ := json.Marshal(map[string]string{"username": "admin", "password": "password"})
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
		if err != nil {
			return ErrorMsg("Connection Failed")
		}
		defer resp.Body.Close()
		
		var res struct{ Token string }
		json.NewDecoder(resp.Body).Decode(&res)
		if res.Token == "" {
			// Try Register first in case of clean DB
			regUrl := "http://localhost:8080/api/auth/register"
			regBody, _ := json.Marshal(map[string]string{"username": "admin", "password": "password", "role": "admin"})
			http.Post(regUrl, "application/json", bytes.NewBuffer(regBody))
			return m.LoginCmd()() // Retry once
		}
		return LoginSuccessMsg(res.Token)
	}
}

func (m Model) ListenWSCmd() tea.Cmd {
	return func() tea.Msg {
		if m.Token == "" {
			return nil
		}
		
		u := fmt.Sprintf("ws://localhost:8080/api/chat?token=%s&manga_id=1", m.Token)
		conn, _, err := websocket.DefaultDialer.Dial(u, nil)
		if err != nil {
			return ErrorMsg("WS Connection Failed")
		}
		
		// In a real TUI system, we might need a more complex state to keep the connection open.
		// For this demo, we'll start a goroutine that reads and returns only the first message
		go func() {
			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
				// Normally we'd use a channel to p.Send(msg).
				// Since we can't easily access p here without passing it in,
				// we'll stick to the "Non-blocking" principle by not blocking the Update() call.
				fmt.Printf("\a") // Audible beep on new message (TUI style)
			}
		}()

		return ChatMsg("Connected to Global Chat! 🌸")
	}
}

func (m Model) SendChatCmd(content string) tea.Cmd {
	return func() tea.Msg {
		// Mock send for the TUI demo
		return ChatMsg("You: " + content)
	}
}
