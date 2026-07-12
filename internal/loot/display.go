package loot

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	blessStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Italic(true)
	curseStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Italic(true)
	colateralStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
	itemStyle      = lipgloss.NewStyle().Bold(true)
	metaStyle      = lipgloss.NewStyle().Faint(true)
)

func DisplayItems(items []any) string {
	var content strings.Builder

	for i, item := range items {
		if i > 0 {
			content.WriteString("\n")
		}
		switch v := item.(type) {
		case Weapon:
			content.WriteString(displayWeapon(v))
		case Armor:
			content.WriteString(displayArmor(v))
		case Accessory:
			content.WriteString(displayAccessory(v))
		case PotionResult:
			content.WriteString(displayPotion(v))
		}
	}

	return borderBox(content.String())
}

func borderBox(content string) string {
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	max := 0
	for _, line := range lines {
		w := lipgloss.Width(line)
		if w > max {
			max = w
		}
	}

	var b strings.Builder
	b.WriteString("╭" + strings.Repeat("─", max+2) + "╮\n")
	for _, line := range lines {
		w := lipgloss.Width(line)
		padding := max - w
		// Use strings.Builder for correct Unicode padding
		padded := line + strings.Repeat(" ", padding)
		b.WriteString("│ " + padded + " │\n")
	}
	b.WriteString("╰" + strings.Repeat("─", max+2) + "╯\n")
	return b.String()
}

func displayWeapon(w Weapon) string {
	var b strings.Builder
	b.WriteString(itemStyle.Render(w.Name) + " ")
	b.WriteString(metaStyle.Render(fmt.Sprintf("%d %s", w.Cost, w.Currency)) + " ")
	b.WriteString(metaStyle.Render(fmt.Sprintf("%s %s", w.Damage, w.DamageType)) + "\n")
	for _, e := range w.Enchantment {
		b.WriteString("  ")
		b.WriteString(formatEnchant(e))
	}
	return b.String()
}

func displayArmor(a Armor) string {
	var b strings.Builder
	b.WriteString(itemStyle.Render(a.Name) + " ")
	b.WriteString(metaStyle.Render(fmt.Sprintf("%d %s", a.Cost, a.Currency)) + " ")
	b.WriteString(metaStyle.Render(fmt.Sprintf("CA %d", a.ArmorBonus)) + "\n")
	for _, e := range a.Enchantment {
		b.WriteString("  ")
		b.WriteString(formatEnchant(e))
	}
	return b.String()
}

func displayAccessory(a Accessory) string {
	var b strings.Builder
	b.WriteString(itemStyle.Render(a.Name) + " ")
	b.WriteString(metaStyle.Render(fmt.Sprintf("%d %s", a.Cost, a.Currency)) + " ")
	b.WriteString(metaStyle.Render(a.Type) + "\n")
	for _, e := range a.Enchantment {
		b.WriteString("  ")
		b.WriteString(formatEnchant(e))
	}
	return b.String()
}

func displayPotion(p PotionResult) string {
	var b strings.Builder
	b.WriteString(itemStyle.Render(p.Potion.Name) + " ")
	b.WriteString(metaStyle.Render(fmt.Sprintf("%d %s", p.Potion.Cost, p.Potion.Currency)) + " ")
	b.WriteString(metaStyle.Render(p.Potion.Effect) + "\n")
	if p.ColateralEffect != nil {
		b.WriteString("  ")
		b.WriteString(colateralStyle.Render(p.ColateralEffect.Description) + "\n")
	}
	return b.String()
}

func formatEnchant(e Enchantment) string {
	txt := fmt.Sprintf("%s: %s", e.Title, e.Description)
	if e.Type == Blessing {
		return blessStyle.Render(txt) + "\n"
	}
	return curseStyle.Render(txt) + "\n"
}
