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
	{Name: "Monsters"},
	{Name: "Iniciativa"},
	{Name: "Encontro"},
}

func renderTabs(active int, width int) string {
	var b strings.Builder

	for i, t := range tabs {
		if i == active {
			b.WriteString(activeTabStyle.Render("● " + t.Name))
		} else {
			b.WriteString(inactiveTabStyle.Render("○ " + t.Name))
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
