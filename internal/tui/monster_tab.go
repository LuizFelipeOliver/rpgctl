package tui

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"

	"rpg-tui/internal/monster"
)

type MonsterModel struct {
	monsters []monster.Monster
	filtered []monster.Monster
	search   string
	err      error
	tbl      table.Model
	fullTblH int
	fullTblW int
}

func newMonsterTable(h int) table.Model {
	columns := []table.Column{
		{Title: "Name", Width: 40},
		{Title: "CR", Width: 8},
		{Title: "Type", Width: 30},
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
		Bold(true).
		Foreground(catBlue)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#1e1e2e")).
		Background(catBlue).
		Bold(true)
	s.Cell = s.Cell.
		Foreground(catText)
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

func (m *MonsterModel) SetTableSize(w, h int) {
	m.fullTblW = w
	m.fullTblH = h - 8
	if m.fullTblH < 5 {
		m.fullTblH = 5
	}
	m.tbl.SetWidth(w)
	m.tbl.SetHeight(m.fullTblH)
}

type MonsterPopupModel struct {
	monster  monster.Monster
	descMode bool
	detailVP viewport.Model
	width    int
}

func NewMonsterPopupModel(mon monster.Monster, vpW int) *MonsterPopupModel {
	var content strings.Builder
	writeMonsterStats(&content, mon, vpW)
	vp := viewport.New(vpW, 12)
	vp.SetContent(content.String())
	vp.GotoTop()
	return &MonsterPopupModel{
		monster:  mon,
		descMode: false,
		detailVP: vp,
		width:    vpW,
	}
}

func (m *MonsterPopupModel) SetContentWidth(w int) {
	if w == m.width {
		return
	}
	m.width = w
	m.detailVP.Width = w
	if m.descMode {
		desc := stripHTML(m.monster.FullText)
		m.detailVP.SetContent(wordwrap.String(desc, w))
	} else {
		var content strings.Builder
		writeMonsterStats(&content, m.monster, w)
		m.detailVP.SetContent(content.String())
	}
	m.detailVP.GotoTop()
}

func (m *MonsterPopupModel) Init() tea.Cmd { return nil }

func (m *MonsterPopupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.descMode {
				m.descMode = false
				var content strings.Builder
				writeMonsterStats(&content, m.monster, m.width)
				m.detailVP.SetContent(content.String())
				m.detailVP.GotoTop()
				return m, nil
			}
			return m, func() tea.Msg { return PopupCloseMsg{} }

		case "enter":
			if m.descMode {
				m.descMode = false
				var content strings.Builder
				writeMonsterStats(&content, m.monster, m.width)
				m.detailVP.SetContent(content.String())
				m.detailVP.GotoTop()
			} else if m.monster.FullText != "" {
				m.descMode = true
				desc := stripHTML(m.monster.FullText)
				m.detailVP.SetContent(wordwrap.String(desc, m.width))
				m.detailVP.GotoTop()
			}
			return m, nil

		case "i":
			return m, func() tea.Msg { return AddToInitMsg{Monster: &m.monster} }
		}

		var cmd tea.Cmd
		m.detailVP, cmd = m.detailVP.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *MonsterPopupModel) View() string {
	var content strings.Builder
	if m.descMode {
		content.WriteString(popupTitle.Render("Descrição"))
		content.WriteString("\n\n")
		content.WriteString(m.detailVP.View())
		content.WriteString("\n\n")
		content.WriteString(faint.Render("[enter] voltar  •  [esc] fechar"))
	} else {
		crBadge := fmt.Sprintf(" CR %s ", m.monster.ChallengeRating)
		content.WriteString(popupTitle.Render(m.monster.Name))
		content.WriteString("  ")
		content.WriteString(crStyle.Render(crBadge))
		content.WriteString("\n\n")
		content.WriteString(m.detailVP.View())
		content.WriteString("\n\n")
		if m.monster.FullText != "" {
			content.WriteString(green.Render("▶ Descrição  "))
			content.WriteString(faint.Render("[enter]"))
			content.WriteString("\n\n")
		}
		content.WriteString(faint.Render("[i] iniciativa  •  [esc] fechar"))
	}
	return content.String()
}

func (m Model) updateMonster(msg tea.KeyMsg) Model {
	if m.monster.monsters == nil {
		return m
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
				outerW := (m.width - 6) * 60 / 100
				cw := outerW - 6
				pm := NewMonsterPopupModel(mon, cw)
				m.popup = NewPopup(pm, 60, 60)
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

func (m Model) renderMonsterContent() string {
	if m.monster.err != nil {
		return red.Render("Erro: " + m.monster.err.Error())
	}
	if m.monster.monsters == nil {
		return "Loading..."
	}
	return m.renderMonsterList()
}

func (m Model) renderMonsterList() string {
	var b strings.Builder

	total := len(m.monster.monsters)
	count := len(m.monster.filtered)
	title := fmt.Sprintf("Monsters  •  %d/%d", count, total)
	b.WriteString(lootTitle.Render(title))
	b.WriteString("\n\n")

	searchDisplay := m.monster.search
	if searchDisplay == "" {
		searchDisplay = "type to search..."
	}
	b.WriteString(faint.Render("Search: "))
	if m.monster.search != "" {
		b.WriteString(searchBox.Render(m.monster.search))
	} else {
		b.WriteString(faint.Render(searchDisplay))
	}
	b.WriteString("\n")

	if count == 0 {
		b.WriteString("\n")
		b.WriteString(faint.Render("No monsters found"))
		return b.String()
	}

	b.WriteString(m.monster.tbl.View())
	return b.String()
}

func writeField(b *strings.Builder, label, value string, style ...lipgloss.Style) {
	if value == "" || value == "0" {
		return
	}
	b.WriteString(statLabel.Render(label))
	b.WriteString("  ")
	if len(style) > 0 {
		b.WriteString(style[0].Render(value))
	} else {
		b.WriteString(value)
	}
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

func parseAbilities(s string) [][2]string {
	var result [][2]string
	if s == "" {
		return result
	}
	parts := strings.Split(s, ",")
	abbr := map[string]string{
		"Str": "FOR", "Dex": "DES", "Con": "CON",
		"Int": "INT", "Wis": "SAB", "Cha": "CAR",
	}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		fields := strings.Fields(p)
		if len(fields) >= 2 {
			name := fields[0]
			val := fields[len(fields)-1]
			if pt, ok := abbr[name]; ok {
				name = pt
			}
			result = append(result, [2]string{name, val})
		}
	}
	return result
}

func renderAbilityGrid(abilities string) string {
	pairs := parseAbilities(abilities)
	if len(pairs) == 0 {
		return ""
	}

	var row1, row2 []string
	for i, p := range pairs {
		card := lipgloss.JoinVertical(lipgloss.Center,
			lipgloss.NewStyle().Bold(true).Foreground(catBlue).Width(6).Align(lipgloss.Center).Render(p[0]),
			lipgloss.NewStyle().Bold(true).Foreground(catText).Width(6).Align(lipgloss.Center).Render(p[1]),
		)
		if i < 3 {
			row1 = append(row1, card)
		} else {
			row2 = append(row2, card)
		}
	}

	var b strings.Builder
	if len(row1) > 0 {
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, row1...))
	}
	b.WriteString("\n")
	if len(row2) > 0 {
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, row2...))
	}
	b.WriteString("\n")
	return b.String()
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

func writeMonsterStats(b *strings.Builder, mon monster.Monster, maxW int) {
	valW := maxW - 16
	if valW < 20 {
		valW = 20
	}
	txtW := maxW - 2
	if txtW < 20 {
		txtW = 20
	}

	sec := func(title string) {
		line := strings.Repeat("─", 8)
		b.WriteString(sectionLine.Render(line + " " + title + " " + line))
		b.WriteString("\n")
	}

	sec("Identidade")
	writeField(b, "Size/Type", formatTypeLine(mon))
	writeField(b, "Hit Dice", mon.HitDice)
	writeField(b, "Alignment", mon.Alignment)
	b.WriteString("\n")

	sec("Defesa")
	writeField(b, "HP", fmt.Sprintf("%d", mon.HP), hpStyle)
	writeField(b, "AC", mon.ArmorClass, acStyle)
	writeField(b, "Saves", mon.Saves, savesStyle)
	b.WriteString("\n")

	sec("Combate")
	writeField(b, "Attack", mon.Attack, attackStyle)
	writeField(b, "Full Attack", wordwrap.String(mon.FullAttack, valW))
	writeField(b, "Base Atk", mon.BaseAttack)
	writeField(b, "Grapple", mon.Grapple)
	writeField(b, "Space/Reach", formatSpaceReach(mon))
	b.WriteString("\n")

	sec("Deslocamento")
	writeField(b, "Speed", mon.Speed)
	writeField(b, "Initiative", mon.Initiative)
	b.WriteString("\n")

	sec("Atributos")
	if g := renderAbilityGrid(mon.Abilities); g != "" {
		b.WriteString("  ")
		b.WriteString(g)
	}
	writeField(b, "Skills", wordwrap.String(mon.Skills, valW))
	writeField(b, "Feats", wordwrap.String(mon.Feats, valW))
	b.WriteString("\n")

	sec("Ambiente")
	writeField(b, "Environment", mon.Environment)
	writeField(b, "Organization", wordwrap.String(mon.Organization, valW))
	writeField(b, "Treasure", mon.Treasure)

	if mon.SpecialAttacks != "" {
		b.WriteString("\n\n")
		line := strings.Repeat("─", 6)
		b.WriteString(sectionLine.Render(line + " Ataques Especiais " + line))
		b.WriteString("\n")
		b.WriteString(wordwrap.String(mon.SpecialAttacks, txtW))
	}

	if mon.SpecialQualities != "" {
		b.WriteString("\n\n")
		line := strings.Repeat("─", 5)
		b.WriteString(sectionLine.Render(line + " Qualidades Especiais " + line))
		b.WriteString("\n")
		b.WriteString(wordwrap.String(mon.SpecialQualities, txtW))
	}
}
