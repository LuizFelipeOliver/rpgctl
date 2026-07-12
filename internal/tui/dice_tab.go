package tui

import (
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"rpg-tui/internal/dice"
)

type DiceModel struct {
	showPopup bool
	input     string
	result    string
	hasError  bool
}

func NewDiceModel() DiceModel {
	return DiceModel{}
}

func (d DiceModel) Update(msg tea.KeyMsg) DiceModel {
	if !d.showPopup {
		if msg.String() == "d" {
			d.showPopup = true
			d.input = ""
		}
		return d
	}

	switch msg.String() {
	case "esc":
		d.showPopup = false
	case "enter":
		d.roll()
	case "backspace":
		if len(d.input) > 0 {
			d.input = d.input[:len(d.input)-1]
		}
	default:
		if isDiceChar(msg.String()) {
			d.input += msg.String()
		}
	}
	return d
}

func isDiceChar(s string) bool {
	if len(s) != 1 {
		return false
	}
	c := rune(s[0])
	return unicode.IsDigit(c) || c == 'd' || c == '+' || c == '-'
}

func (d *DiceModel) roll() {
	r, err := dice.Roll(d.input)
	if err != nil {
		d.result = "Formato invalido. Use NdS ou NdS+M"
		d.hasError = true
		return
	}
	d.result = r.String()
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

	inputDisplay := m.dice.input
	if inputDisplay == "" {
		inputDisplay = "ex: 2d20+3"
	}
	content.WriteString(inputBorder.Render(inputDisplay))
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
