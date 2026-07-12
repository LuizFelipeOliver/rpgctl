package tui

import (
	"strings"

	"github.com/mattn/go-runewidth"
)

type Tab struct {
	Name string
}

var tabs = []Tab{
	{Name: "Loot"},
	{Name: "Dados"},
	{Name: "Iniciativa"},
}

func renderTabs(active int, width int) string {
	var b strings.Builder

	for i, t := range tabs {
		name := t.Name
		if i == active {
			name = "[" + name + "]"
			b.WriteString(activeTabStyle.Render(name))
		} else {
			b.WriteString(inactiveTabStyle.Render(name))
		}
		if i < len(tabs)-1 {
			b.WriteString("  ")
		}
	}

	text := b.String()
	tw := runewidth.StringWidth(text)
	if pad := width - tw; pad > 0 {
		text += strings.Repeat(" ", pad)
	}

	return tabBarStyle.Render(text)
}
