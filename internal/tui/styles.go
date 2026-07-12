package tui

import "github.com/charmbracelet/lipgloss"

var (
	green      = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	red        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	yellow     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
	faint      = lipgloss.NewStyle().Faint(true)
	menuCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).SetString("▸ ")
	menuTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFF00"))
	selected   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00"))
	lootTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFF00"))
	helpStyle  = lipgloss.NewStyle().Faint(true).PaddingTop(1)
)
