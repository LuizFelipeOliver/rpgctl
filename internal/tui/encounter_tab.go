package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"rpg-tui/internal/dice"
	"rpg-tui/internal/encounter"
	"rpg-tui/internal/monster"
)

var difficultyOptions = []string{"F", "M", "D"}

type EncounterModel struct {
	inputs        [4]textinput.Model // jogadores, nivel, tipo, qtd
	focused       int                // 0-3 inputs, 4 difficulty
	difficultyIdx int
	result        *encounter.Result
	hasError      bool
	errMsg        string
	tbl           table.Model
	fullW         int
	fullTH        int
}

func NewEncounterModel() EncounterModel {
	players := textinput.New()
	players.Placeholder = "4"
	players.Validate = digitOptional

	level := textinput.New()
	level.Placeholder = "5"
	level.Validate = digitOptional

	tipo := textinput.New()
	tipo.Placeholder = "todos"

	qtd := textinput.New()
	qtd.Placeholder = "auto"
	qtd.Validate = digitOptional

	players.Focus()

	cols := []table.Column{
		{Title: "Monstro", Width: 28},
		{Title: "ND", Width: 8},
		{Title: "Qtd", Width: 6},
	}
	t := table.New(
		table.WithColumns(cols),
		table.WithRows([]table.Row{}),
		table.WithFocused(false),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true).
		Foreground(catBlue)
	s.Cell = s.Cell.Foreground(catText)
	s.Selected = s.Selected.Foreground(catText)
	t.SetStyles(s)

	return EncounterModel{
		inputs: [4]textinput.Model{players, level, tipo, qtd},
		tbl:    t,
	}
}

func parseInitMod(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	s = strings.TrimLeft(s, "+")
	n, err := strconv.Atoi(strings.Fields(s)[0])
	if err != nil {
		return 0
	}
	return n
}

func (m Model) updateEncounter(msg tea.KeyMsg) Model {
	if msg.String() == "g" {
		m.encounter.generate()
		return m
	}

	switch msg.String() {
	case "tab", "down":
		if m.encounter.focused < 4 {
			m.encounter.inputs[m.encounter.focused].Blur()
		}
		m.encounter.focused = (m.encounter.focused + 1) % 5
		if m.encounter.focused < 4 {
			m.encounter.inputs[m.encounter.focused].Focus()
		}
	case "shift+tab", "up":
		if m.encounter.focused < 4 {
			m.encounter.inputs[m.encounter.focused].Blur()
		}
		m.encounter.focused = (m.encounter.focused + 4) % 5
		if m.encounter.focused < 4 {
			m.encounter.inputs[m.encounter.focused].Focus()
		}
	case "enter":
		m.encounter.generate()
	case "esc":
		m.encounter.result = nil
		m.encounter.hasError = false
		m.encounter.syncTable()
	case "i":
		if m.encounter.result != nil {
			return m.addEncounterToInitiative()
		}
	case "left":
		if m.encounter.focused == 4 {
			m.encounter.difficultyIdx--
			if m.encounter.difficultyIdx < 0 {
				m.encounter.difficultyIdx = len(difficultyOptions) - 1
			}
			return m
		}
	case "right":
		if m.encounter.focused == 4 {
			m.encounter.difficultyIdx++
			if m.encounter.difficultyIdx >= len(difficultyOptions) {
				m.encounter.difficultyIdx = 0
			}
			return m
		}
	}

	if m.encounter.focused < 4 {
		var cmd tea.Cmd
		m.encounter.inputs[m.encounter.focused], cmd = m.encounter.inputs[m.encounter.focused].Update(msg)
		_ = cmd
	}

	return m
}

func (m *EncounterModel) SetTableSize(w, h int) {
	m.fullW = w
	m.fullTH = h - 5
	if m.fullTH < 5 {
		m.fullTH = 5
	}

	leftW := w * 40 / 100
	inputW := leftW - 16
	if inputW < 10 {
		inputW = 10
	}
	for i := 0; i < 4; i++ {
		m.inputs[i].Width = inputW
	}

	rightW := w*60/100 - 2
	if rightW < 30 {
		rightW = 30
	}
	monsterW := rightW - 8 - 6 - 4
	if monsterW < 10 {
		monsterW = 10
	}
	m.tbl.SetColumns([]table.Column{
		{Title: "Monstro", Width: monsterW},
		{Title: "ND", Width: 8},
		{Title: "Qtd", Width: 6},
	})
	m.tbl.SetWidth(rightW)
	m.tbl.SetHeight(m.fullTH)
}

func (e *EncounterModel) generate() {
	playersStr := strings.TrimSpace(e.inputs[0].Value())
	levelStr := strings.TrimSpace(e.inputs[1].Value())
	tipoStr := strings.TrimSpace(e.inputs[2].Value())
	qtdStr := strings.TrimSpace(e.inputs[3].Value())

	if playersStr == "" || levelStr == "" {
		e.result = nil
		e.hasError = true
		e.errMsg = "Preencha Jogadores e Nível"
		e.syncTable()
		return
	}

	players, err := strconv.Atoi(playersStr)
	if err != nil || players <= 0 {
		e.result = nil
		e.hasError = true
		e.errMsg = "Número de jogadores inválido"
		e.syncTable()
		return
	}

	level, err := strconv.Atoi(levelStr)
	if err != nil || level <= 0 {
		e.result = nil
		e.hasError = true
		e.errMsg = "Nível inválido"
		e.syncTable()
		return
	}

	var qty int
	if qtdStr != "" {
		qty, err = strconv.Atoi(qtdStr)
		if err != nil || qty <= 0 {
			e.result = nil
			e.hasError = true
			e.errMsg = "Quantidade inválida"
			e.syncTable()
			return
		}
	}

	diff := difficultyOptions[e.difficultyIdx]

	monsters, err := monster.LoadMonsters("data/monster/monsters-parsed.json")
	if err != nil {
		e.result = nil
		e.hasError = true
		e.errMsg = "Erro carregando monstros: " + err.Error()
		e.syncTable()
		return
	}

	opts := encounter.GenerateOptions{TypeFilter: tipoStr, Quantity: qty}
	r, err := encounter.GenerateWithOpts(monsters, players, level, diff, opts)
	if err != nil {
		e.result = nil
		e.hasError = true
		e.errMsg = err.Error()
		e.syncTable()
		return
	}

	e.result = r
	e.hasError = false
	e.syncTable()
}

func (e *EncounterModel) syncTable() {
	if e.result == nil {
		e.tbl.SetRows([]table.Row{})
		return
	}
	rows := make([]table.Row, len(e.result.Groups))
	for i, g := range e.result.Groups {
		rows[i] = table.Row{g.Monster.Name, g.Monster.ChallengeRating, strconv.Itoa(g.Quantity)}
	}
	e.tbl.SetRows(rows)
}

func (m Model) addEncounterToInitiative() Model {
	if m.encounter.result == nil {
		return m
	}

	for _, group := range m.encounter.result.Groups {
		for j := 0; j < group.Quantity; j++ {
			name := group.Monster.Name
			if group.Quantity > 1 {
				name += fmt.Sprintf(" #%d", j+1)
			}

			initRoll, err := dice.RollNotation("1d20")
			if err != nil {
				initRoll = 0
			}
			initiative := initRoll + parseInitMod(group.Monster.Initiative)

			c := Combatant{
				Name:       name,
				Initiative: initiative,
				HP:         group.Monster.HP,
				MaxHP:      group.Monster.HP,
				AC:         firstAC(group.Monster.ArmorClass),
				Monster:    &group.Monster,
			}
			m.init.combatants = append(m.init.combatants, c)
		}
	}

	m.init.sorted()
	m.init.syncTable()

	m.activeTab = 3
	return m
}

func (m Model) renderEncounterContent() string {
	innerW := m.width - 6

	var b strings.Builder
	b.WriteString(lootTitle.Render("Gerar Encontro"))
	b.WriteString("\n\n")

	leftW := innerW * 40 / 100
	rightW := innerW - leftW - 2

	left := m.renderEncounterControls(leftW)
	right := m.renderEncounterResults(rightW)

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, left, "  ", right))
	return b.String()
}

func (m Model) renderEncounterControls(w int) string {
	var b strings.Builder
	labels := []string{"Jogadores", "Nível", "Tipo", "Qtd"}

	for i, label := range labels {
		if i == m.encounter.focused {
			b.WriteString(green.Render("> "))
		} else {
			b.WriteString("  ")
		}
		b.WriteString(label)
		b.WriteString(" [")
		if i == m.encounter.focused {
			b.WriteString(m.encounter.inputs[i].View())
		} else {
			v := m.encounter.inputs[i].Value()
			if v == "" {
				b.WriteString(faint.Render(m.encounter.inputs[i].Placeholder))
			} else {
				b.WriteString(v)
			}
		}
		b.WriteString("]\n")
	}

	if m.encounter.focused == 4 {
		b.WriteString(green.Render("> "))
	} else {
		b.WriteString("  ")
	}
	b.WriteString("Dificuldade ")
	if m.encounter.focused == 4 {
		b.WriteString(difficultySelector(m.encounter.difficultyIdx))
	} else {
		b.WriteString(difficultyOptions[m.encounter.difficultyIdx])
	}
	b.WriteString("\n\n")

	if m.encounter.hasError {
		b.WriteString(red.Render(m.encounter.errMsg))
		b.WriteString("\n\n")
	}

	b.WriteString(faint.Render("g/Enter: gerar"))
	if m.encounter.result != nil {
		b.WriteString(faint.Render("  •  i: add iniciativa"))
	}
	b.WriteString(faint.Render("  •  esc: limpar"))

	return lipgloss.NewStyle().Width(w).Render(b.String())
}

func (m Model) renderEncounterResults(w int) string {
	if m.encounter.result == nil {
		return ""
	}

	var b strings.Builder
	b.WriteString(m.encounter.tbl.View())
	b.WriteString("\n")
	b.WriteString(faint.Render(fmt.Sprintf("EL %d (%d×%d, %s)",
		m.encounter.result.TargetEL,
		m.encounter.result.PartyCount,
		m.encounter.result.PartyLevel,
		m.encounter.result.Difficulty,
	)))

	return lipgloss.NewStyle().Width(w).Render(b.String())
}

func difficultySelector(idx int) string {
	prev := difficultyOptions[(idx+len(difficultyOptions)-1)%len(difficultyOptions)]
	curr := difficultyOptions[idx]
	next := difficultyOptions[(idx+1)%len(difficultyOptions)]
	return fmt.Sprintf("%s  %s  %s", faint.Render("< "+prev), green.Render(curr), faint.Render(next+" >"))
}
