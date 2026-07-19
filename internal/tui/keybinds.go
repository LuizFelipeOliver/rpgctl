package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

type kmap []key.Binding

func (k kmap) ShortHelp() []key.Binding   { return k }
func (k kmap) FullHelp() [][]key.Binding   { return [][]key.Binding{k} }

func newHelpModel() help.Model {
	h := help.New()
	h.Styles.ShortKey = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#1d4ed8", Dark: "#89b4fa"}).
		Bold(true)
	h.Styles.ShortDesc = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#9ca3af", Dark: "#6c7086"})
	h.Styles.ShortSeparator = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#d1d5db", Dark: "#45475a"})
	return h
}

func (m Model) footerKeys() help.KeyMap {
	var km kmap

	km = append(km,
		key.NewBinding(key.WithKeys("ctrl+h"), key.WithHelp("ctrl+h", "aba ←")),
		key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("ctrl+l", "aba →")),
		key.NewBinding(key.WithKeys("ctrl+d"), key.WithHelp("ctrl+d", "dados")),
	)

	if m.popup != nil {
		km = append(km, key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "fechar")))
	} else {
		switch m.activeTab {
		case 0:
			km = append(km, key.NewBinding(key.WithKeys("g"), key.WithHelp("[N]g", "gerar")))
		case 1:
			km = append(km,
				key.NewBinding(key.WithKeys("↑/↓"), key.WithHelp("↑/↓", "navegar")),
				key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "detalhes")),
			)
		case 2:
			km = append(km,
				key.NewBinding(key.WithKeys("a"), key.WithHelp("a", "add")),
				key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "del")),
				key.NewBinding(key.WithKeys("n/p"), key.WithHelp("n/p", "turno")),
				key.NewBinding(key.WithKeys("+/-"), key.WithHelp("+/-", "cura/dano")),
				key.NewBinding(key.WithKeys("v"), key.WithHelp("v", "detalhes")),
				key.NewBinding(key.WithKeys("ctrl+r"), key.WithHelp("ctrl+r", "reset")),
			)
		case 3:
			km = append(km, key.NewBinding(key.WithKeys("g"), key.WithHelp("g", "gerar")))
		}
	}

	km = append(km, key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "sair")))
	return km
}
