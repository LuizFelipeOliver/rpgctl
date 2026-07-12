package tui

import "rpg-tui/internal/loot"

type LootModel struct {
	table *loot.LootTable
	items []any
	err   error
}

func NewLootModel() LootModel {
	return LootModel{}
}

func (m *LootModel) Init() {
	table, err := loot.NewLootTable("data")
	if err != nil {
		m.err = err
		return
	}
	m.table = table
	m.err = nil
	m.Generate()
}

func (m *LootModel) Generate() {
	m.items = m.table.Generate(5)
}

func (m LootModel) View() string {
	var s string
	s += lootTitle.Render("Gerar Itens") + "\n\n"

	if m.err != nil {
		s += red.Render("Erro: " + m.err.Error())
		s += helpStyle.Render("\nesc/q: voltar")
		return s
	}

	if len(m.items) == 0 {
		s += "Pressione g para gerar itens\n"
	} else {
		s += loot.DisplayItems(m.items)
	}

	s += helpStyle.Render("\ng: gerar  •  esc/q: voltar")
	return s
}
