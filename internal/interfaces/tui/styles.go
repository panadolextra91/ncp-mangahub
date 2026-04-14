package tui

import "github.com/charmbracelet/lipgloss"

// Pink Palette
const (
    PinkPastel = lipgloss.Color("#FBCFE8")
    PinkGlow   = lipgloss.Color("#F9A8D4")
    SubtleBlack = lipgloss.Color("#0B0E11")
    Slate500    = lipgloss.Color("#64748B")
    Slate400    = lipgloss.Color("#94A3B8")
    Slate200    = lipgloss.Color("#E2E8F0")
)

var (
    BaseStyle = lipgloss.NewStyle().
        Foreground(Slate200).
        Background(SubtleBlack)

    HeaderStyle = lipgloss.NewStyle().
        Foreground(PinkPastel).
        Bold(true).
        Padding(0, 1).
        Border(lipgloss.RoundedBorder()).
        BorderForeground(PinkPastel)

    TabStyle = lipgloss.NewStyle().
        Foreground(Slate400).
        Padding(0, 1)

    ActiveTabStyle = lipgloss.NewStyle().
        Foreground(PinkPastel).
        Bold(true).
        Underline(true).
        Padding(0, 1)

    BorderStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("205")) // Pinkish

    GlowStyle = lipgloss.NewStyle().
        Foreground(PinkPastel).
        Bold(true)

    AsciiStyle = lipgloss.NewStyle().
        Foreground(PinkPastel).
        Italic(true)
)
