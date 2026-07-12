package tui

import "github.com/charmbracelet/lipgloss"

var (
	green      = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	red        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	yellow     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
	faint      = lipgloss.NewStyle().Faint(true)
	bold       = lipgloss.NewStyle().Bold(true)

	activeTabStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00")).Padding(0, 1)
	inactiveTabStyle = lipgloss.NewStyle().Faint(true).Padding(0, 1)

	tabBarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true)

	footerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true)

	popupStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true).
			BorderForeground(lipgloss.Color("#00FF00")).
			Padding(1, 2)

	popupTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFF00"))
	lootTitle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFF00"))
	diceResult  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00"))
	inputBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true).
			Padding(0, 1).
			Width(20)
)
