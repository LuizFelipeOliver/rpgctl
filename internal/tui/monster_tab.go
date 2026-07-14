package tui

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"rpg-tui/internal/monster"
)

type MonsterModel struct {
	monsters []monster.Monster
	filtered []monster.Monster
	search   string
	scroll   int
	detail   *monster.Monster
	err      error
	tbl      table.Model
}

func newMonsterTable(h int) table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 35},
		{Title: "CR", Width: 8},
		{Title: "Type", Width: 25},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(h),
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

func NewMonsterModel() MonsterModel {
	m := MonsterModel{
		tbl: newMonsterTable(20),
	}
	m.load()
	return m
}

func (m *MonsterModel) load() {
	monsters, err := monster.LoadMonsters("data/monster/monsters-parsed.json")
	if err != nil {
		m.err = err
		return
	}
	m.monsters = monsters
	m.filtered = monsters
	m.syncTable()
}

func (m *MonsterModel) syncTable() {
	rows := make([]table.Row, len(m.filtered))
	for i, mon := range m.filtered {
		rows[i] = table.Row{mon.Name, mon.ChallengeRating, mon.Type}
	}
	m.tbl.SetRows(rows)
}

func (m *MonsterModel) doSearch() {
	m.filtered = monster.Search(m.monsters, m.search)
	m.syncTable()
}

func (m *MonsterModel) SetTableHeight(h int) {
	th := h - 8
	if th < 5 {
		th = 5
	}
	m.tbl.SetHeight(th)
}

func (m Model) updateMonster(msg tea.KeyMsg) Model {
	if m.monster.monsters == nil {
		return m
	}
	if m.monster.detail != nil {
		return m.updateMonsterDetail(msg)
	}
	return m.updateMonsterList(msg)
}

func (m Model) updateMonsterList(msg tea.KeyMsg) Model {
	s := msg.String()

	switch s {
	case "enter":
		if len(m.monster.filtered) > 0 {
			cursor := m.monster.tbl.Cursor()
			if cursor < len(m.monster.filtered) {
				mon := m.monster.filtered[cursor]
				m.monster.detail = &mon
				m.monster.scroll = 0
				m.monster.tbl.Blur()
			}
		}
	case "backspace":
		if len(m.monster.search) > 0 {
			m.monster.search = m.monster.search[:len(m.monster.search)-1]
			m.monster.doSearch()
		}
	case "esc":
		m.monster.search = ""
		m.monster.doSearch()
	case "i":
		if len(m.monster.filtered) > 0 {
			cursor := m.monster.tbl.Cursor()
			if cursor < len(m.monster.filtered) {
				mon := m.monster.filtered[cursor]
				m.init.pendingMonster = &mon
				m.init.inputs[0].SetValue(mon.Name)
				m.init.inputs[1].SetValue("")
				m.init.inputs[2].SetValue(strconv.Itoa(mon.HP))
				m.init.inputs[3].SetValue(strconv.Itoa(firstAC(mon.ArmorClass)))
				m.init.inputs[0].Blur()
				m.init.inputs[1].Focus()
				m.init.addMode = 2
				m.activeTab = 3
			}
		}
	default:
		if len(s) == 1 && s[0] >= 32 && s[0] <= 126 {
			m.monster.search += s
			m.monster.doSearch()
		} else {
			var cmd tea.Cmd
			m.monster.tbl, cmd = m.monster.tbl.Update(msg)
			_ = cmd
		}
	}
	return m
}

func (m Model) updateMonsterDetail(msg tea.KeyMsg) Model {
	switch msg.String() {
	case "esc", "enter":
		m.monster.detail = nil
		m.monster.tbl.Focus()
	case "up":
		if m.monster.scroll > 0 {
			m.monster.scroll--
		}
	case "down":
		m.monster.scroll++
	case "i":
		m.init.pendingMonster = m.monster.detail
		m.init.inputs[0].SetValue(m.monster.detail.Name)
		m.init.inputs[1].SetValue("")
		m.init.inputs[2].SetValue(strconv.Itoa(m.monster.detail.HP))
		m.init.inputs[3].SetValue(strconv.Itoa(firstAC(m.monster.detail.ArmorClass)))
		m.init.inputs[0].Blur()
		m.init.inputs[1].Focus()
		m.init.addMode = 2
		m.init.detail = nil
		m.monster.detail = nil
		m.monster.tbl.Focus()
		m.activeTab = 3
	}
	return m
}

func (m Model) renderMonsterContent() string {
	if m.monster.err != nil {
		return red.Render("Erro: " + m.monster.err.Error())
	}
	if m.monster.monsters == nil {
		return "Loading..."
	}
	if m.monster.detail != nil {
		return m.renderMonsterDetail()
	}
	return m.renderMonsterList()
}

func (m Model) renderMonsterList() string {
	var b strings.Builder
	b.WriteString(lootTitle.Render("Monsters"))
	b.WriteString("\n\n")

	searchDisplay := m.monster.search
	if searchDisplay == "" {
		searchDisplay = "type to search..."
	}
	b.WriteString(faint.Render("Search: "))
	b.WriteString(yellow.Render(searchDisplay))
	b.WriteString("\n")

	if len(m.monster.filtered) == 0 {
		b.WriteString(faint.Render("No monsters found"))
		return b.String()
	}

	b.WriteString(m.monster.tbl.View())
	return b.String()
}

func (m Model) renderMonsterDetail() string {
	mon := *m.monster.detail
	var b strings.Builder

	b.WriteString(lootTitle.Render(mon.Name))
	b.WriteString("  ")
	b.WriteString(faint.Render("CR: " + mon.ChallengeRating))
	b.WriteString("\n")
	b.WriteString(faint.Render(strings.Repeat("─", 50)))
	b.WriteString("\n\n")

	writeMonsterStats(&b, mon, m.width-2)

	b.WriteString("\n")
	b.WriteString(faint.Render("[i] add to initiative  •  [esc] back"))
	return b.String()
}

func writeField(b *strings.Builder, label, value string) {
	if value == "" || value == "0" {
		return
	}
	b.WriteString(bold.Render(label + ":"))
	b.WriteString(strings.Repeat(" ", 13-len(label)))
	if len(label) >= 13 {
		b.WriteString(" ")
	}
	b.WriteString(value)
	b.WriteString("\n")
}

func formatTypeLine(mon monster.Monster) string {
	var parts []string
	if mon.Size != "" {
		parts = append(parts, mon.Size)
	}
	if mon.Type != "" {
		parts = append(parts, mon.Type)
	}
	if mon.Descriptor != "" {
		parts = append(parts, "("+mon.Descriptor+")")
	}
	return strings.Join(parts, " ")
}

func formatSpaceReach(mon monster.Monster) string {
	var parts []string
	if mon.Space != "" {
		parts = append(parts, mon.Space)
	}
	if mon.Reach != "" {
		parts = append(parts, mon.Reach)
	}
	return strings.Join(parts, " / ")
}

func truncate(s string, n int) string {
	if n <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n-3]) + "..."
}

func stripHTML(s string) string {
	var b strings.Builder
	inTag := false
	for _, r := range s {
		if r == '<' {
			inTag = true
			continue
		}
		if r == '>' {
			inTag = false
			continue
		}
		if !inTag {
			b.WriteRune(r)
		}
	}
	result := strings.TrimSpace(b.String())
	lines := strings.Split(result, "\n")
	var cleaned []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" || unicode.IsLetter(rune(lastChar(cleaned))) {
			cleaned = append(cleaned, line)
		}
	}
	return strings.Join(cleaned, "\n")
}

func lastChar(s []string) byte {
	if len(s) == 0 {
		return 0
	}
	last := s[len(s)-1]
	if len(last) == 0 {
		return 0
	}
	return last[len(last)-1]
}
