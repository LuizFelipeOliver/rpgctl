package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
)

type Model struct {
	activeTab int
	width     int
	height    int
	countBuf  string
	loot      LootModel
	dice      DiceModel
	monster   MonsterModel
	init      InitiativeModel
}

func New() Model {
	return Model{
		activeTab: 0,
		width:     80,
		height:    24,
		loot:      NewLootModel(),
		dice:      NewDiceModel(),
		monster:   NewMonsterModel(),
		init:      NewInitiativeModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.monster.SetTableHeight(msg.Height)
		m.init.SetTableHeight(msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "q":
			return m, tea.Quit

		case "h", "left":
			if m.activeTab > 0 {
				m.activeTab--
				m.countBuf = ""
			}
			return m, nil

		case "l", "right":
			if m.activeTab < len(tabs)-1 {
				m.activeTab++
				m.countBuf = ""
			}
			return m, nil
		}

		switch m.activeTab {
		case 0:
			return m.updateLoot(msg), nil
		case 1:
			return m.updateDice(msg), nil
		case 2:
			return m.updateMonster(msg), nil
		case 3:
			return m.updateInitiative(msg), nil
		}
	}

	return m, nil
}

func (m Model) updateLoot(msg tea.KeyMsg) Model {
	switch msg.String() {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		m.countBuf += msg.String()
	case "g":
		n := 5
		if m.countBuf != "" {
			if v, err := strconv.Atoi(m.countBuf); err == nil && v > 0 {
				n = v
			}
			if n > 100 {
				n = 100
			}
		}
		m.countBuf = ""
		m.loot.Generate(n)
	default:
		m.countBuf = ""
	}
	return m
}

func (m Model) updateDice(msg tea.KeyMsg) Model {
	m.dice = m.dice.Update(msg)
	return m
}

func (m Model) View() string {
	tabBar := renderTabs(m.activeTab, m.width)
	content := m.renderContent()
	footer := m.renderFooter()

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lw := runewidth.StringWidth(line)
		if pad := m.width - lw; pad > 0 {
			lines[i] = line + strings.Repeat(" ", pad)
		}
	}
	content = strings.Join(lines, "\n")

	tabLines := strings.Count(tabBar, "\n") + 1
	contentLines := strings.Count(content, "\n") + 1
	footerLines := strings.Count(footer, "\n") + 1

	padding := m.height - tabLines - contentLines - footerLines
	if padding < 0 {
		padding = 0
	}

	return tabBar + "\n" + content + strings.Repeat("\n", padding) + "\n" + footer
}

func (m Model) renderContent() string {
	switch m.activeTab {
	case 0:
		return m.renderLootContent()
	case 1:
		return m.renderDiceContent()
	case 2:
		return m.renderMonsterContent()
	case 3:
		return m.renderInitiativeContent()
	}
	return ""
}

func (m Model) renderFooter() string {
	var text string
	switch m.activeTab {
	case 0:
		text = "h/l: abas  •  [N]g: gerar  •  q: sair"
	case 1:
		text = "h/l: abas  •  d: dados  •  q: sair"
	case 2:
		if m.monster.detail != nil {
			text = "h/l: abas  •  esc: voltar  •  ↑↓: rolar  •  i: add iniciativa  •  q: sair"
		} else {
			text = "h/l: abas  •  digitar: buscar  •  ↑↓: navegar  •  enter: detalhes  •  q: sair"
		}
	case 3:
		if m.init.addMode > 0 || m.init.showDetail {
			text = "h/l: abas  •  esc: voltar  •  q: sair"
		} else {
			text = "h/l: abas  •  a: add  •  d: delete  •  n/p: turno  •  enter: detalhes  •  r: reset  •  q: sair"
		}
	}
	tw := runewidth.StringWidth(text)
	if pad := m.width - tw; pad > 0 {
		text += strings.Repeat(" ", pad)
	}
	return footerStyle.Render(text)
}
