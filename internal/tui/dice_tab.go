package tui

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DiceModel struct {
	showPopup bool
	input     string
	result    string
}

func NewDiceModel() DiceModel {
	return DiceModel{}
}

func (d DiceModel) Update(msg tea.KeyMsg) DiceModel {
	if !d.showPopup {
		if msg.String() == "d" {
			d.showPopup = true
			d.input = ""
			d.result = ""
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

var dicePattern = regexp.MustCompile(`^(\d*)d(\d+)([+-]\d+)?$`)

func (d *DiceModel) roll() {
	matches := dicePattern.FindStringSubmatch(d.input)
	if matches == nil {
		d.result = "Formato invalido. Use NdS ou NdS+M"
		return
	}

	numDice := 1
	if matches[1] != "" {
		n, err := strconv.Atoi(matches[1])
		if err == nil && n > 0 {
			numDice = n
		}
	}

	sides, _ := strconv.Atoi(matches[2])
	if sides < 2 {
		sides = 2
	}
	if sides > 1000 {
		sides = 1000
	}
	if numDice > 100 {
		numDice = 100
	}

	modifier := 0
	if matches[3] != "" {
		modifier, _ = strconv.Atoi(matches[3])
	}

	rolls := make([]string, numDice)
	total := 0
	for i := range numDice {
		r := rand.Intn(sides) + 1
		rolls[i] = strconv.Itoa(r)
		total += r
	}

	parts := strings.Join(rolls, " + ")
	if modifier != 0 {
		parts = fmt.Sprintf("%s %+d", parts, modifier)
		total += modifier
	}
	d.result = fmt.Sprintf("%s = %d", parts, total)
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
		b.WriteString(green.Render(m.dice.result))
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
		content.WriteString(diceResult.Render(m.dice.result))
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
