package tui

import "strings"

type menuItem struct {
	title string
	desc  string
	state State
}

type MenuModel struct {
	items  []menuItem
	cursor int
}

func NewMenu() MenuModel {
	return MenuModel{
		items: []menuItem{
			{title: "Gerar Itens", desc: "Gera itens aleatorios (arma, armadura, acessorio, pocao)", state: StateLoot},
			{title: "Rolar Dados", desc: "Simula rolagens de dados", state: StateDice},
			{title: "Iniciativa", desc: "Gerenciar ordem de iniciativa", state: StateInitiative},
			{title: "Sair", desc: "Encerrar o programa", state: StateQuit},
		},
		cursor: 0,
	}
}

func (m MenuModel) View() string {
	var b strings.Builder

	b.WriteString(menuTitle.Render("rpgctl - Gerenciador de Campanhas"))
	b.WriteString("\n\n")

	for i, item := range m.items {
		prefix := "  "
		if i == m.cursor {
			prefix = menuCursor.String()
		}

		title := item.title
		if i == m.cursor {
			title = selected.Render(item.title)
		}

		b.WriteString(prefix + title + "\n")

		if i == m.cursor {
			b.WriteString("   " + faint.Render(item.desc) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("j/k: navegar  •  Enter: selecionar  •  1-4: atalho  •  q: sair"))

	return b.String()
}
