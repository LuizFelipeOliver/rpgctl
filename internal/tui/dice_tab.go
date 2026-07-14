package tui

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"rpg-tui/internal/dice"
)

type DiceModel struct {
	showPopup bool
	input     textinput.Model
	result    string
	hasError  bool
}

func NewDiceModel() DiceModel {
	ti := textinput.New()
	ti.Placeholder = "ex: 2d20+3"
	ti.Validate = func(s string) error {
		for _, r := range s {
			if !unicode.IsDigit(r) && r != 'd' && r != '+' && r != '-' && r != ' ' {
				return errInvalidChar
			}
		}
		return nil
	}
	return DiceModel{input: ti}
}

var errInvalidChar = &invalidCharError{}

type invalidCharError struct{}

func (e *invalidCharError) Error() string { return "invalid character" }

func (d DiceModel) Update(msg tea.KeyMsg) DiceModel {
	if !d.showPopup {
		if msg.String() == "d" {
			d.showPopup = true
			d.input.SetValue("")
			d.input.Focus()
		}
		return d
	}

	switch msg.String() {
	case "esc":
		d.showPopup = false
		d.input.Blur()
	case "enter":
		d.roll()
	default:
		var cmd tea.Cmd
		d.input, cmd = d.input.Update(msg)
		_ = cmd
	}
	return d
}

func (d *DiceModel) roll() {
	r, err := dice.Rolar(d.input.Value())
	if err != nil {
		d.result = "Formato invalido. Ex: d20 + 1 + d3"
		d.hasError = true
		return
	}
	d.result = r.Detalhes + " = " + strconv.Itoa(r.Total)
	d.hasError = false
}

func (m Model) renderDiceContent() string {
	if m.dice.showPopup {
		return m.renderDicePopup()
	}

	var b strings.Builder
	b.WriteString(lootTitle.Render("Rolar Dados"))
	b.WriteString("\n\n")

	if m.dice.result != "" {
		b.WriteString("Ultimo resultado:\n")
		if m.dice.hasError {
			b.WriteString(red.Render(m.dice.result))
		} else {
			b.WriteString(green.Render(m.dice.result))
		}
		b.WriteString("\n\n")
	}

	b.WriteString("Pressione d para abrir o seletor de dados")
	return b.String()
}

func (m Model) renderDicePopup() string {
	var content strings.Builder
	content.WriteString(popupTitle.Render("Rolar Dados"))
	content.WriteString("\n\n")

	content.WriteString(inputBorder.Render(m.dice.input.View()))
	content.WriteString("\n\n")

	if m.dice.result != "" {
		content.WriteString("Resultado:\n")
		if m.dice.hasError {
			content.WriteString(red.Render(m.dice.result))
		} else {
			content.WriteString(diceResult.Render(m.dice.result))
		}
		content.WriteString("\n\n")
	}

	content.WriteString(faint.Render("Enter: rolar  •  esc: fechar"))

	popup := popupStyle.Render(content.String())

	ch := m.height - 4
	if ch < 1 {
		ch = 1
	}
	return lipgloss.Place(m.width, ch, lipgloss.Center, lipgloss.Center, popup)
}
