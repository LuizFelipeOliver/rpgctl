package loot

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type styledLine struct {
	sep    bool
	plain  string
	styled string
}

var (
	green      = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	red        = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	yellow     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
	bold       = lipgloss.NewStyle().Bold(true)
	faint      = lipgloss.NewStyle().Faint(true)
	boldGreen  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00"))
	boldRed    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
	boldYellow = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFF00"))
)

func DisplayItems(items []any) string {
	var lines []styledLine

	for i, item := range items {
		if i > 0 {
			lines = append(lines, styledLine{sep: true})
		}
		switch v := item.(type) {
		case Weapon:
			lines = appendItemLine(lines, v.Name, v.Cost, v.Currency, fmt.Sprintf("%s %s", v.Damage, v.DamageType))
			lines = appendEnchants(lines, v.Enchantment)
		case Armor:
			lines = appendItemLine(lines, v.Name, v.Cost, v.Currency, fmt.Sprintf("CA %d", v.ArmorBonus))
			lines = appendEnchants(lines, v.Enchantment)
		case Accessory:
			lines = appendItemLine(lines, v.Name, v.Cost, v.Currency, v.Type)
			lines = appendEnchants(lines, v.Enchantment)
		case PotionResult:
			lines = appendItemLine(lines, v.Potion.Name, v.Potion.Cost, v.Potion.Currency, v.Potion.Effect)
			if v.ColateralEffect != nil {
				lines = appendStyledLine(lines, "  colateral: "+v.ColateralEffect.Description, boldYellow)
			}
		}
	}

	return borderBox(lines)
}

func appendItemLine(lines []styledLine, name string, cost int, currency, detail string) []styledLine {
	plain := fmt.Sprintf("%s (%d %s)  %s", name, cost, currency, detail)
	styled := fmt.Sprintf("%s %s  %s",
		bold.Render(name),
		faint.Render(fmt.Sprintf("(%d %s)", cost, currency)),
		faint.Render(detail))
	return append(lines, styledLine{plain: plain, styled: styled})
}

func appendStyledLine(lines []styledLine, text string, style lipgloss.Style) []styledLine {
	return append(lines, styledLine{plain: text, styled: style.Render(text)})
}

func appendEnchants(lines []styledLine, enchants []Enchantment) []styledLine {
	var blessings, curses []Enchantment

	for _, e := range enchants {
		if e.Type == Blessing {
			blessings = append(blessings, e)
		} else {
			curses = append(curses, e)
		}
	}

	if len(blessings) > 0 {
		lines = appendStyledLine(lines, "  benções", boldGreen)
		for _, e := range blessings {
			lines = appendStyledLine(lines, "    titulo: "+e.Title, green)
			lines = appendStyledLine(lines, "    descrição: "+e.Description, green)
		}
	}

	if len(curses) > 0 {
		lines = appendStyledLine(lines, "  maldições", boldRed)
		for _, e := range curses {
			lines = appendStyledLine(lines, "    titulo: "+e.Title, red)
			lines = appendStyledLine(lines, "    descrição: "+e.Description, red)
		}
	}

	return lines
}

func borderBox(lines []styledLine) string {
	if len(lines) == 0 {
		return ""
	}

	max := 0
	for _, l := range lines {
		if w := runewidth.StringWidth(l.plain); w > max {
			max = w
		}
	}

	var b strings.Builder
	b.Grow((max + 6) * (len(lines) + 2))

	b.WriteString("╭" + strings.Repeat("─", max+2) + "╮\n")

	for _, l := range lines {
		if l.sep {
			b.WriteString("├" + strings.Repeat("─", max+2) + "┤\n")
			continue
		}
		w := runewidth.StringWidth(l.plain)
		b.WriteString("│ " + l.styled + strings.Repeat(" ", max-w) + " │\n")
	}

	b.WriteString("╰" + strings.Repeat("─", max+2) + "╯\n")
	return b.String()
}
