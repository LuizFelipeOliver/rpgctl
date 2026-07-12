package tui

import (
	"strings"

	"rpg-tui/internal/loot"
)

type LootModel struct {
	table *loot.LootTable
	items []any
	err   error
}

func NewLootModel() LootModel {
	return LootModel{}
}

func (m *LootModel) Generate(n int) {
	if m.table == nil {
		table, err := loot.NewLootTable("data")
		if err != nil {
			m.err = err
			return
		}
		m.table = table
	}
	m.err = nil
	m.items = m.table.Generate(n)
}

func (m Model) renderLootContent() string {
	var b strings.Builder

	b.WriteString(lootTitle.Render("Gerar Itens"))
	b.WriteString("\n\n")

	if m.countBuf != "" {
		b.WriteString(yellow.Render("Quantidade: " + m.countBuf))
		b.WriteString("\n\n")
	}

	if m.loot.err != nil {
		b.WriteString(red.Render("Erro: " + m.loot.err.Error()))
	} else if len(m.loot.items) == 0 {
		b.WriteString("Pressione g para gerar itens")
	} else {
		b.WriteString(loot.DisplayItems(m.loot.items))
	}

	return b.String()
}
