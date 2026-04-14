package tui

import (
	"time"

	"github.com/charmbracelet/bubbletea"
)

type Page int

const (
	PageLogin Page = iota
	PageChat
	PageEvents
	PageCreate
)

type Model struct {
	ActivePage  Page
	AsciiIndex  int
	Token       string
	Username    string
	Password    string
	ChatInput   string
	MangaTitle  string
	MangaAuthor string
	MangaDesc   string
	Messages    []string
	Events      []string
	Status      string
	Width       int
	Height      int
	Error       string
}

type TickMsg time.Time
type ChatMsg string
type EventMsg string
type LoginSuccessMsg string
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
