package tui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"rpg-tui/internal/monster"
)

type Combatant struct {
	Name       string
	Initiative int
	HP         int
	MaxHP      int
	AC         int
	Monster    *monster.Monster
}

type InitiativeModel struct {
	combatants []Combatant
	current    int
	round      int
	tbl        table.Model
}

func newInitTable() table.Model {
	columns := []table.Column{
		{Title: "Init", Width: 8},
		{Title: "Name", Width: 30},
		{Title: "HP", Width: 10},
		{Title: "AC", Width: 6},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true).
		Foreground(catBlue)
	s.Selected = s.Selected.
		Foreground(catGreen).
		Background(catSurface).
		Bold(true)
	s.Cell = s.Cell.
		Foreground(catText)
	t.SetStyles(s)
	return t
}

func NewInitiativeModel() InitiativeModel {
	return InitiativeModel{
		tbl: newInitTable(),
	}
}

func (m *InitiativeModel) SetTableSize(w, h int) {
	m.tbl.SetWidth(w)
	th := h - 7
	if th < 5 {
		th = 5
	}
	m.tbl.SetHeight(th)
}

func (m *InitiativeModel) sorted() {
	sort.SliceStable(m.combatants, func(i, j int) bool {
		return m.combatants[i].Initiative > m.combatants[j].Initiative
	})
}

func (m *InitiativeModel) syncTable() {
	rows := make([]table.Row, len(m.combatants))
	for i, c := range m.combatants {
		hpStr := fmt.Sprintf("%d/%d", c.HP, c.MaxHP)
		if c.MaxHP == 0 {
			hpStr = "-"
		}
		acStr := "-"
		if c.AC > 0 {
			acStr = strconv.Itoa(c.AC)
		}
		rows[i] = table.Row{
			strconv.Itoa(c.Initiative),
			c.Name,
			hpStr,
			acStr,
		}
	}
	m.tbl.SetRows(rows)
}

func (m *InitiativeModel) remove(idx int) {
	if idx < 0 || idx >= len(m.combatants) {
		return
	}
	m.combatants = append(m.combatants[:idx], m.combatants[idx+1:]...)
	if m.current >= len(m.combatants) && len(m.combatants) > 0 {
		m.current = 0
	}
}

func (m *InitiativeModel) reset() {
	m.combatants = nil
	m.current = 0
	m.round = 0
}

func (m *InitiativeModel) nextTurn() {
	if len(m.combatants) == 0 {
		return
	}
	if m.current < len(m.combatants)-1 {
		m.current++
	} else {
		m.current = 0
		m.round++
	}
}

func (m *InitiativeModel) prevTurn() {
	if len(m.combatants) == 0 {
		return
	}
	if m.current > 0 {
		m.current--
	} else {
		m.current = len(m.combatants) - 1
		if m.round > 0 {
			m.round--
		}
	}
}

func (m Model) updateInitiative(msg tea.KeyMsg) Model {
	s := msg.String()

	switch s {
	case "a":
		wiz := NewWizardPopupModel()
		m.popup = NewPopup(wiz, 60, 75)

	case "d":
		if len(m.init.combatants) > 0 {
			cursor := m.init.tbl.Cursor()
			m.init.remove(cursor)
			m.init.syncTable()
		}

	case "ctrl+r":
		m.popup = NewPopup(NewResetPopupModel(), 40, 40)

	case "n", "right":
		m.init.nextTurn()
		m.init.tbl.SetCursor(m.init.current)

	case "p", "left":
		m.init.prevTurn()
		m.init.tbl.SetCursor(m.init.current)

	case "-":
		if len(m.init.combatants) > 0 {
			cursor := m.init.tbl.Cursor()
			hp := NewHealPopupModel(cursor, true)
			m.popup = NewPopup(hp, 40, 40)
		}

	case "+":
		if len(m.init.combatants) > 0 {
			cursor := m.init.tbl.Cursor()
			hp := NewHealPopupModel(cursor, false)
			m.popup = NewPopup(hp, 40, 40)
		}

	case "enter", "v":
		if len(m.init.combatants) > 0 {
			cursor := m.init.tbl.Cursor()
			if cursor < 0 {
				break
			}
			c := &m.init.combatants[cursor]
			outerW := (m.width - 6) * 60 / 100
			cw := outerW - 6
			dp := NewInitDetailPopupModel(c, cw)
			m.popup = NewPopup(dp, 60, 60)
		}

	default:
		var cmd tea.Cmd
		m.init.tbl, cmd = m.init.tbl.Update(msg)
		_ = cmd
	}
	return m
}

func (m Model) renderInitiativeContent() string {
	return m.renderInitList()
}

func (m Model) renderInitList() string {
	var b strings.Builder
	title := "Initiative"
	if m.init.round > 0 || len(m.init.combatants) > 0 {
		title += fmt.Sprintf(" — Round %d", m.init.round+1)
	}
	if len(m.init.combatants) > 0 {
		title += fmt.Sprintf(" — %s", m.init.combatants[m.init.current].Name)
	}
	b.WriteString(lootTitle.Render(title))
	b.WriteString("\n")

	if len(m.init.combatants) == 0 {
		b.WriteString("\n")
		b.WriteString(faint.Render("No combatants. Press 'a' to add."))
		return b.String()
	}

	b.WriteString(m.init.tbl.View())
	return b.String()
}

func firstAC(s string) int {
	if s == "" {
		return 0
	}
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}
