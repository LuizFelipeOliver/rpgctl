package tui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"

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
	combatants     []Combatant
	current        int
	round          int
	addMode        int
	showDetail     bool
	detail         *monster.Monster
	detailCombat   *Combatant
	pendingMonster *monster.Monster
	inputs         [4]textinput.Model
	tbl            table.Model
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
		Bold(true)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#00FF00")).
		Bold(true)
	t.SetStyles(s)
	return t
}

func newWizardInput(placeholder string, validate func(string) error) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Validate = validate
	return ti
}

func initValidator(s string) error {
	if s == "" {
		return nil
	}
	for i, r := range s {
		if r >= '0' && r <= '9' {
			continue
		}
		if i == 0 && (r == '-' || r == '+') {
			continue
		}
		return errInvalidChar
	}
	return nil
}

func digitOptional(s string) error {
	if s == "" {
		return nil
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return errInvalidChar
		}
	}
	return nil
}

func NewInitiativeModel() InitiativeModel {
	return InitiativeModel{
		inputs: [4]textinput.Model{
			newWizardInput("Name...", nil),
			newWizardInput("Initiative...", initValidator),
			newWizardInput("HP (optional)...", digitOptional),
			newWizardInput("AC (optional)...", digitOptional),
		},
		tbl: newInitTable(),
	}
}

func (m *InitiativeModel) SetTableHeight(h int) {
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

func (m *InitiativeModel) addFromWizard() {
	name := strings.TrimSpace(m.inputs[0].Value())
	if name == "" {
		return
	}
	init, _ := strconv.Atoi(strings.TrimSpace(m.inputs[1].Value()))
	hp, _ := strconv.Atoi(strings.TrimSpace(m.inputs[2].Value()))
	ac, _ := strconv.Atoi(strings.TrimSpace(m.inputs[3].Value()))

	c := Combatant{
		Name:       name,
		Initiative: init,
		HP:         hp,
		MaxHP:      hp,
		AC:         ac,
		Monster:    m.pendingMonster,
	}
	m.pendingMonster = nil
	m.combatants = append(m.combatants, c)
	m.sorted()
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
	m.addMode = 0
	m.showDetail = false
	m.detail = nil
	m.detailCombat = nil
	m.pendingMonster = nil
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
	if m.init.addMode > 0 {
		return m.updateInitWizard(msg)
	}
	if m.init.showDetail {
		return m.updateInitDetail(msg)
	}
	return m.updateInitList(msg)
}

func (m Model) updateInitWizard(msg tea.KeyMsg) Model {
	switch msg.String() {
	case "esc":
		for i := range m.init.inputs {
			m.init.inputs[i].Blur()
			m.init.inputs[i].SetValue("")
		}
		m.init.addMode = 0
	case "enter":
		switch m.init.addMode {
		case 1:
			if strings.TrimSpace(m.init.inputs[0].Value()) != "" {
				m.init.inputs[0].Blur()
				m.init.inputs[1].Focus()
				m.init.addMode = 2
			}
		case 2:
			if strings.TrimSpace(m.init.inputs[1].Value()) != "" {
				m.init.inputs[1].Blur()
				m.init.inputs[2].Focus()
				m.init.addMode = 3
			}
		case 3:
			m.init.inputs[2].Blur()
			m.init.inputs[3].Focus()
			m.init.addMode = 4
		case 4:
			m.init.addFromWizard()
			m.init.syncTable()
			for i := range m.init.inputs {
				m.init.inputs[i].SetValue("")
				m.init.inputs[i].Blur()
			}
			m.init.addMode = 0
		}
	default:
		focused := m.init.addMode - 1
		if focused >= 0 && focused < 4 {
			var cmd tea.Cmd
			m.init.inputs[focused], cmd = m.init.inputs[focused].Update(msg)
			_ = cmd
		}
	}
	return m
}

func (m Model) updateInitDetail(msg tea.KeyMsg) Model {
	switch msg.String() {
	case "esc", "enter":
		m.init.showDetail = false
		m.init.detail = nil
		m.init.detailCombat = nil
	case "up":
		if m.monster.scroll > 0 {
			m.monster.scroll--
		}
	case "down":
		m.monster.scroll++
	}
	return m
}

func (m Model) updateInitList(msg tea.KeyMsg) Model {
	switch msg.String() {
	case "a":
		m.init.addMode = 1
		m.init.inputs[0].Focus()
		m.init.inputs[0].SetValue("")
	case "d":
		if len(m.init.combatants) > 0 {
			cursor := m.init.tbl.Cursor()
			m.init.remove(cursor)
			m.init.syncTable()
		}
	case "r":
		m.init.reset()
		m.init.syncTable()
	case "n", "right":
		m.init.nextTurn()
		m.init.tbl.SetCursor(m.init.current)
	case "p", "left":
		m.init.prevTurn()
		m.init.tbl.SetCursor(m.init.current)
	case "enter", "v":
		if len(m.init.combatants) > 0 {
			cursor := m.init.tbl.Cursor()
			m.init.detailCombat = &m.init.combatants[cursor]
			m.init.detail = m.init.combatants[cursor].Monster
			m.init.showDetail = true
			m.monster.scroll = 0
		}
	default:
		var cmd tea.Cmd
		m.init.tbl, cmd = m.init.tbl.Update(msg)
		_ = cmd
	}
	return m
}

func (m Model) renderInitiativeContent() string {
	if m.init.addMode > 0 {
		return m.renderInitWizard()
	}
	if m.init.showDetail {
		return m.renderInitDetail()
	}
	return m.renderInitList()
}

func (m Model) renderInitWizard() string {
	var b strings.Builder
	b.WriteString(popupTitle.Render("Add Combatant"))
	b.WriteString("\n\n")

	labels := []string{"Name", "Initiative", "HP (optional)", "AC (optional)"}

	for i := 0; i < 4; i++ {
		step := i + 1
		if step == m.init.addMode {
			b.WriteString(green.Render("> " + m.init.inputs[i].View()))
		} else if step < m.init.addMode {
			b.WriteString(faint.Render("  " + m.init.inputs[i].Value()))
		} else {
			b.WriteString(faint.Render("  " + labels[i]))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(faint.Render("Enter: next  •  esc: cancel"))

	popup := popupStyle.Render(b.String())
	ch := m.height - 2
	if ch < 1 {
		ch = 1
	}
	return lipgloss.Place(m.width, ch, lipgloss.Center, lipgloss.Center, popup)
}

func (m Model) renderInitDetail() string {
	var b strings.Builder

	if m.init.detailCombat != nil {
		c := m.init.detailCombat
		b.WriteString(lootTitle.Render(c.Name))
		b.WriteString("\n")
		b.WriteString(faint.Render(strings.Repeat("─", 50)))
		b.WriteString("\n\n")

		if m.init.detail != nil {
			writeMonsterStats(&b, *m.init.detail, m.width-6)
		} else {
			writeField(&b, "Initiative", strconv.Itoa(c.Initiative))
			hpStr := fmt.Sprintf("%d/%d", c.HP, c.MaxHP)
			if c.MaxHP == 0 {
				hpStr = "-"
			}
			writeField(&b, "HP", hpStr)
			if c.AC > 0 {
				writeField(&b, "AC", strconv.Itoa(c.AC))
			}
		}
	} else if m.init.detail != nil {
		mon := *m.init.detail
		b.WriteString(lootTitle.Render(mon.Name))
		b.WriteString("\n")
		b.WriteString(faint.Render(strings.Repeat("─", 50)))
		b.WriteString("\n\n")
		writeMonsterStats(&b, mon, m.width-6)
	}

	b.WriteString("\n")
	b.WriteString(faint.Render("[esc] back"))

	popup := popupStyle.Render(b.String())
	ch := m.height - 2
	if ch < 1 {
		ch = 1
	}
	return lipgloss.Place(m.width, ch, lipgloss.Center, lipgloss.Center, popup)
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

func writeMonsterStats(b *strings.Builder, mon monster.Monster, maxW int) {
	valW := maxW - 15
	if valW < 20 {
		valW = 20
	}
	txtW := maxW - 2
	if txtW < 20 {
		txtW = 20
	}

	writeField(b, "Size/Type", formatTypeLine(mon))
	writeField(b, "Hit Dice", mon.HitDice)
	writeField(b, "HP", fmt.Sprintf("%d", mon.HP))
	writeField(b, "Initiative", mon.Initiative)
	writeField(b, "Speed", mon.Speed)
	writeField(b, "AC", mon.ArmorClass)
	writeField(b, "Base Atk", mon.BaseAttack)
	writeField(b, "Grapple", mon.Grapple)
	writeField(b, "Attack", mon.Attack)
	writeField(b, "Full Attack", wordwrap.String(mon.FullAttack, valW))
	writeField(b, "Space/Reach", formatSpaceReach(mon))
	writeField(b, "Saves", mon.Saves)
	writeField(b, "Abilities", mon.Abilities)
	writeField(b, "Skills", wordwrap.String(mon.Skills, valW))
	writeField(b, "Feats", wordwrap.String(mon.Feats, valW))
	writeField(b, "Environment", mon.Environment)
	writeField(b, "Organization", wordwrap.String(mon.Organization, valW))
	writeField(b, "Treasure", mon.Treasure)
	writeField(b, "Alignment", mon.Alignment)

	if mon.SpecialAttacks != "" {
		b.WriteString("\n")
		b.WriteString(yellow.Render("Special Attacks:"))
		b.WriteString("\n")
		b.WriteString(wordwrap.String(mon.SpecialAttacks, txtW))
		b.WriteString("\n")
	}

	if mon.SpecialQualities != "" {
		b.WriteString("\n")
		b.WriteString(yellow.Render("Special Qualities:"))
		b.WriteString("\n")
		b.WriteString(wordwrap.String(mon.SpecialQualities, txtW))
		b.WriteString("\n")
	}

	if mon.FullText != "" {
		b.WriteString("\n")
		b.WriteString(faint.Render("Description:"))
		b.WriteString("\n")
		b.WriteString(stripHTML(mon.FullText))
		b.WriteString("\n")
	}
}
