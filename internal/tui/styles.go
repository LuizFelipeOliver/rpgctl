package tui

import "github.com/charmbracelet/lipgloss"

// ── Catppuccin Mocha Palette ──────────────────────────────────
var (
	catBase    = lipgloss.Color("#1e1e2e")
	catSurface = lipgloss.Color("#313244")
	catOverlay = lipgloss.Color("#45475a")
	catMuted   = lipgloss.Color("#6c7086")
	catText    = lipgloss.Color("#cdd6f4")
	catGreen   = lipgloss.Color("#a6e3a1")
	catRed     = lipgloss.Color("#f38ba8")
	catYellow  = lipgloss.Color("#f9e2af")
	catBlue    = lipgloss.Color("#89b4fa")
	catMauve   = lipgloss.Color("#cba6f7")
	catPeach   = lipgloss.Color("#fab387")
	catSky     = lipgloss.Color("#89dceb")
	catPink    = lipgloss.Color("#f5c2e7")
)

// ── Semantic Styles ──────────────────────────────────────────
var (
	green  = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#1a8a1a", Dark: "#a6e3a1"})
	red    = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#cc0000", Dark: "#f38ba8"})
	yellow = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#9a6700", Dark: "#f9e2af"})
	faint  = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#9ca3af", Dark: "#6c7086"})

	hpStyle    = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#1a8a1a", Dark: "#a6e3a1"})
	acStyle    = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#7c3aed", Dark: "#cba6f7"})
	crStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{Light: "#be185d", Dark: "#f5c2e7"})
	savesStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#0e7490", Dark: "#89dceb"})
	attackStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#c2410c", Dark: "#fab387"})

	// ── Layout ──────────────────────────────────────────────
	borderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(catOverlay).
			Padding(1, 2)

	tabBarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderBottom(true).
			BorderForeground(catOverlay)

	activeTabStyle = lipgloss.NewStyle().
			Background(catBlue).
			Foreground(catBase).
			Bold(true).
			Padding(0, 2)

	inactiveTabStyle = lipgloss.NewStyle().
			Foreground(catMuted).
			Padding(0, 2)

	footerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderForeground(catOverlay)

	popupStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(catBlue).
			Padding(1, 2)

	popupTitle = lipgloss.NewStyle().Bold(true).Foreground(catYellow)
	lootTitle  = popupTitle

	sectionStyle = lipgloss.NewStyle().Bold(true).Foreground(catBlue)
	sectionLine  = lipgloss.NewStyle().Bold(true).Foreground(catBlue)

	statLabel = lipgloss.NewStyle().
			Bold(true).
			Foreground(catMuted).
			Width(14).
			Align(lipgloss.Right)

	statValue = lipgloss.NewStyle().Foreground(catText)

	inputBorder = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(catOverlay).
			Padding(0, 1)

	searchBox = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(catBlue).
			Padding(0, 1).
			Foreground(catYellow).
			Width(30)

	diceResult = lipgloss.NewStyle().Bold(true).Foreground(catGreen)
)
