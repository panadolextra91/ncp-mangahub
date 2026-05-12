package tui

import (
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

type Page int

const (
	PageLogin Page = iota
	PageChat
	PageEvents
	PageCreate
	PageProgress
	PageSearch
)

type Model struct {
	ActivePage  Page
	AsciiIndex  int
	Token       string
	Username    string
	Password    string
	Role        string
	ChatInput   string
	MangaTitle   string
	MangaAuthor  string
	MangaGenres  string
	MangaStatus  string
	MangaChapters string
	MangaDesc    string
	MangaIDInput string
	ChapterInput string
	StatusInput  string // Reading, Completed, etc.
	SearchInput  string
	SearchResults []string
	Messages    []string
	Events      []string
	Status      string
	Width       int
	Height      int
	Error       string
	FocusIndex  int // 0, 1, 2 for inputs
	SearchIndex int // Scrolling index for search results
	WS          *websocket.Conn
	Stream      EventStream
}

type TickMsg time.Time
type ChatMsg string
type EventMsg string
type LoginSuccessMsg struct {
	Token string
	Role  string
}
type ErrorMsg string

func NewModel() Model {
	return Model{
		ActivePage: PageLogin,
		AsciiIndex: 0,
		Status:     "Pink Hub Ready 🌸",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		TickCommand(),
	)
}

func TickCommand() tea.Cmd {
	return tea.Tick(60*time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
