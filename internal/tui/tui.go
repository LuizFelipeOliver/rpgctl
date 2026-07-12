package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type State int

const (
	StateMenu State = iota
	StateLoot
	StateDice
	StateInitiative
	StateQuit
)

type Model struct {
	state State
	menu  MenuModel
	loot  LootModel
}

func New() Model {
	return Model{
		state: StateMenu,
		menu:  NewMenu(),
		loot:  NewLootModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}

		switch m.state {
		case StateMenu:
			return m.updateMenu(msg)
		case StateLoot:
			return m.updateLoot(msg)
		case StateDice, StateInitiative:
			return m.updatePlaceholder(msg)
		}
	}
	return m, nil
}

func (m Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.menu.cursor < len(m.menu.items)-1 {
			m.menu.cursor++
		}
	case "k", "up":
		if m.menu.cursor > 0 {
			m.menu.cursor--
		}
	case "enter", " ":
		m.state = m.menu.items[m.menu.cursor].state
		if m.state == StateLoot {
			m.loot.Init()
		}
		if m.state == StateQuit {
			return m, tea.Quit
		}
	case "1":
		m.state = StateLoot
		m.menu.cursor = 0
		m.loot.Init()
	case "2":
		m.state = StateDice
		m.menu.cursor = 1
	case "3":
		m.state = StateInitiative
		m.menu.cursor = 2
	case "4", "q":
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) updateLoot(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "g":
		m.loot.Generate()
	case "esc", "q":
		m.state = StateMenu
	}
	return m, nil
}

func (m Model) updatePlaceholder(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q":
		m.state = StateMenu
	}
	return m, nil
}

func (m Model) View() string {
	var s string

	switch m.state {
	case StateMenu:
		s = m.menu.View()
	case StateLoot:
		s = m.loot.View()
	case StateDice:
		s = m.placeholderView("Rolar Dados", "Em breve...")
	case StateInitiative:
		s = m.placeholderView("Iniciativa", "Em breve...")
	case StateQuit:
		return ""
	}

	return s + "\n"
}

func (m Model) placeholderView(title, text string) string {
	var b strings.Builder
	b.WriteString(lootTitle.Render(title))
	b.WriteString("\n\n")
	b.WriteString(text)
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("esc/q: voltar"))
	return b.String()
}
