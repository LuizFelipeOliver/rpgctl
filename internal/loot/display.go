package loot

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type styledLine struct {
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
	var allLines []styledLine

	for _, item := range items {
		if len(allLines) > 0 {
			allLines = append(allLines, styledLine{plain: "─", styled: "─"})
		}
		switch v := item.(type) {
		case Weapon:
			allLines = append(allLines, weaponLines(v)...)
		case Armor:
			allLines = append(allLines, armorLines(v)...)
		case Accessory:
			allLines = append(allLines, accessoryLines(v)...)
		case PotionResult:
			allLines = append(allLines, potionLines(v)...)
		}
	}

	return borderBox(allLines)
}

func weaponLines(w Weapon) []styledLine {
	lines := []styledLine{
		{
			plain:  fmt.Sprintf("%s (%d %s)  %s %s", w.Name, w.Cost, w.Currency, w.Damage, w.DamageType),
			styled: fmt.Sprintf("%s %s  %s", bold.Render(w.Name), faint.Render(fmt.Sprintf("(%d %s)", w.Cost, w.Currency)), faint.Render(fmt.Sprintf("%s %s", w.Damage, w.DamageType))),
		},
	}
	lines = append(lines, enchantLines(w.Enchantment)...)
	return lines
}

func armorLines(a Armor) []styledLine {
	lines := []styledLine{
		{
			plain:  fmt.Sprintf("%s (%d %s)  CA %d", a.Name, a.Cost, a.Currency, a.ArmorBonus),
			styled: fmt.Sprintf("%s %s  %s", bold.Render(a.Name), faint.Render(fmt.Sprintf("(%d %s)", a.Cost, a.Currency)), faint.Render(fmt.Sprintf("CA %d", a.ArmorBonus))),
		},
	}
	lines = append(lines, enchantLines(a.Enchantment)...)
	return lines
}

func accessoryLines(a Accessory) []styledLine {
	lines := []styledLine{
		{
			plain:  fmt.Sprintf("%s (%d %s)  %s", a.Name, a.Cost, a.Currency, a.Type),
			styled: fmt.Sprintf("%s %s  %s", bold.Render(a.Name), faint.Render(fmt.Sprintf("(%d %s)", a.Cost, a.Currency)), faint.Render(a.Type)),
		},
	}
	lines = append(lines, enchantLines(a.Enchantment)...)
	return lines
}

func potionLines(p PotionResult) []styledLine {
	lines := []styledLine{
		{
			plain:  fmt.Sprintf("%s (%d %s)  %s", p.Potion.Name, p.Potion.Cost, p.Potion.Currency, p.Potion.Effect),
			styled: fmt.Sprintf("%s %s  %s", bold.Render(p.Potion.Name), faint.Render(fmt.Sprintf("(%d %s)", p.Potion.Cost, p.Potion.Currency)), faint.Render(p.Potion.Effect)),
		},
	}
	if p.ColateralEffect != nil {
		l := "colateral: " + p.ColateralEffect.Description
		lines = append(lines, styledLine{plain: "  " + l, styled: "  " + boldYellow.Render(l)})
	}
	return lines
}

func enchantLines(enchants []Enchantment) []styledLine {
	var lines []styledLine
	var blessings, curses []Enchantment

	for _, e := range enchants {
		if e.Type == Blessing {
			blessings = append(blessings, e)
		} else {
			curses = append(curses, e)
		}
	}

	if len(blessings) > 0 {
		lines = append(lines, styledLine{plain: "  benções", styled: "  " + boldGreen.Render("benções")})
		for _, e := range blessings {
			lines = append(lines, styledLine{plain: "    titulo: " + e.Title, styled: "    " + green.Render("titulo: "+e.Title)})
			lines = append(lines, styledLine{plain: "    descrição: " + e.Description, styled: "    " + green.Render("descrição: "+e.Description)})
		}
	}

	if len(curses) > 0 {
		lines = append(lines, styledLine{plain: "  maldições", styled: "  " + boldRed.Render("maldições")})
		for _, e := range curses {
			lines = append(lines, styledLine{plain: "    titulo: " + e.Title, styled: "    " + red.Render("titulo: "+e.Title)})
			lines = append(lines, styledLine{plain: "    descrição: " + e.Description, styled: "    " + red.Render("descrição: "+e.Description)})
		}
	}

	return lines
}

func borderBox(lines []styledLine) string {
	max := 0
	for _, l := range lines {
		w := runewidth.StringWidth(l.plain)
		if w > max {
			max = w
		}
	}

	var b strings.Builder
	b.WriteString("╭" + strings.Repeat("─", max+2) + "╮\n")

	for _, l := range lines {
		if l.plain == "─" {
			b.WriteString("├" + strings.Repeat("─", max+2) + "┤\n")
			continue
		}
		w := runewidth.StringWidth(l.plain)
		pad := max - w
		b.WriteString("│ " + l.styled + strings.Repeat(" ", pad) + " │\n")
	}

	b.WriteString("╰" + strings.Repeat("─", max+2) + "╯\n")
	return b.String()
}
