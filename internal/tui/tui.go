package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	activeTab int
	width     int
	height    int
	countBuf  string
	popup     *Popup
	help      help.Model
	loot      LootModel
	dice      DiceModel
	monster   MonsterModel
	init      InitiativeModel
	encounter EncounterModel
}

func New() Model {
	return Model{
		activeTab: 0,
		width:     80,
		height:    24,
		help:      newHelpModel(),
		loot:      NewLootModel(),
		dice:      NewDiceModel(),
		monster:   NewMonsterModel(),
		init:      NewInitiativeModel(),
		encounter: NewEncounterModel(),
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
		m.loot.SetTableSize(msg.Width-6, msg.Height-10)
		m.monster.SetTableSize(msg.Width-6, msg.Height-10)
		m.init.SetTableSize(msg.Width-6, msg.Height-9)
		m.encounter.SetTableSize(msg.Width-6, msg.Height-10)
		if m.popup != nil {
			m.popup.SetSize(msg.Width-6, msg.Height-12)
		}
		return m, nil

	case PopupCloseMsg:
		m.popup = nil
		return m, nil

	case AddToInitMsg:
		m.popup = nil
		wiz := NewWizardPopupModel()
		wiz.Prefill(msg.Monster)
		m.popup = NewPopup(wiz, 60, 75)
		m.activeTab = 3
		return m, nil

	case ResetConfirmMsg:
		m.popup = nil
		m.init.reset()
		m.init.syncTable()
		return m, nil

	case HealApplyMsg:
		m.popup = nil
		if msg.Cursor >= 0 && msg.Cursor < len(m.init.combatants) {
			if msg.IsDamage {
				m.init.combatants[msg.Cursor].HP -= msg.Amount
				if m.init.combatants[msg.Cursor].HP < 0 {
					m.init.combatants[msg.Cursor].HP = 0
				}
			} else {
				m.init.combatants[msg.Cursor].HP += msg.Amount
				max := m.init.combatants[msg.Cursor].MaxHP
				if m.init.combatants[msg.Cursor].HP > max {
					m.init.combatants[msg.Cursor].HP = max
				}
			}
			m.init.syncTable()
		}
		return m, nil

	case WizardCompleteMsg:
		m.popup = nil
		name := strings.TrimSpace(msg.Name)
		if name != "" {
			c := Combatant{
				Name:       name,
				Initiative: msg.Initiative,
				HP:         msg.HP,
				MaxHP:      msg.HP,
				AC:         msg.AC,
				Monster:    msg.Monster,
			}
			m.init.combatants = append(m.init.combatants, c)
			m.init.sorted()
			m.init.syncTable()
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+h":
			if m.activeTab > 0 {
				m.activeTab--
				m.countBuf = ""
			}
			return m, nil

		case "ctrl+l":
			if m.activeTab < len(tabs)-1 {
				m.activeTab++
				m.countBuf = ""
			}
			return m, nil

		case "ctrl+d":
			m.dice = DiceModel{input: m.dice.input}
			m.dice.input.Focus()
			m.popup = NewPopup(&m.dice, 50, 50)
			return m, nil
		}

		if m.popup != nil {
			var cmd tea.Cmd
			var updated tea.Model
			updated, cmd = m.popup.Update(msg)
			m.popup = updated.(*Popup)
			return m, cmd
		}

		switch m.activeTab {
		case 0:
			return m.updateLoot(msg), nil
		case 1:
			return m.updateMonster(msg), nil
		case 2:
			return m.updateInitiative(msg), nil
		case 3:
			return m.updateEncounter(msg), nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	innerW := m.width - 6
	tabBar := renderTabs(m.activeTab, innerW)
	footer := m.renderFooter()

	tabLines := lipgloss.Height(tabBar)
	footerLines := lipgloss.Height(footer)
	availH := m.height - 4 - tabLines - footerLines
	if availH < 1 {
		availH = 1
	}

	body := lipgloss.Place(innerW, availH, lipgloss.Top, lipgloss.Left, m.renderContent())
	if m.popup != nil {
		m.popup.SetSize(innerW, availH)
		body = m.popup.View()
	}

	content := lipgloss.JoinVertical(lipgloss.Top, tabBar, body, footer)
	return borderStyle.Render(content)
}

func (m Model) renderContent() string {
	switch m.activeTab {
	case 0:
		return m.renderLootContent()
	case 1:
		return m.renderMonsterContent()
	case 2:
		return m.renderInitiativeContent()
	case 3:
		return m.renderEncounterContent()
	}
	return ""
}

func (m Model) renderFooter() string {
	m.help.Width = m.width - 6
	return footerStyle.Render(m.help.View(m.footerKeys()))
}
