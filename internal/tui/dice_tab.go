package tui

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"rpg-tui/internal/dice"
)

type DiceModel struct {
	input    textinput.Model
	result   string
	hasError bool
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
	ti.Focus()
	return DiceModel{input: ti}
}

func (d *DiceModel) Init() tea.Cmd {
	return textinput.Blink
}

func (d *DiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			d.input.Blur()
			return d, func() tea.Msg { return PopupCloseMsg{} }

		case "enter":
			r, err := dice.Rolar(d.input.Value())
			if err != nil {
				d.result = "Formato invalido. Ex: d20 + 1 + d3"
				d.hasError = true
			} else {
				d.result = r.Detalhes + " = " + strconv.Itoa(r.Total)
				d.hasError = false
			}
			return d, nil
		}
	}

	var cmd tea.Cmd
	d.input, cmd = d.input.Update(msg)
	return d, cmd
}

func (d DiceModel) View() string {
	var content strings.Builder
	content.WriteString(popupTitle.Render("Rolar Dados"))
	content.WriteString("\n\n")

	content.WriteString(inputBorder.Render(d.input.View()))
	content.WriteString("\n\n")

	if d.result != "" {
		content.WriteString("Resultado:\n")
		if d.hasError {
			content.WriteString(red.Render(d.result))
		} else {
			content.WriteString(diceResult.Render(d.result))
		}
		content.WriteString("\n\n")
	}

	content.WriteString(faint.Render("Enter: rolar  •  esc: fechar"))
	return content.String()
}
