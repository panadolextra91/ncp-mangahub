package tui

import (
	"time"

	"github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"github.com/user/mangahub/pkg/models"
)

type Page int

const (
	PageLogin Page = iota
	PageChat
	PageEvents
	PageCreate
	PageProgress
	PageSearch
	PageDiscover
	PageLibrary
	PageDetail
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
	Quotes      []models.Quote
	CurrentQuoteIndex int
	LibraryResults    []*models.UserProgress
	SelectedManga     *models.Manga
	LibraryIndex      int
}

type TickMsg time.Time
type ChatMsg string
type EventMsg string
type LoginSuccessMsg struct {
	Token string
	Role  string
}
type ErrorMsg string
type TCPSyncMsg string

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
		m.ScrapeQuotesCmd(),
	)
}

func TickCommand() tea.Cmd {
	return tea.Tick(60*time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
