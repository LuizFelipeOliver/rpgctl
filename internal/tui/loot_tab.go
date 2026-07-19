package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"rpg-tui/internal/loot"
)

type LootModel struct {
	table  *loot.LootTable
	items  []any
	err    error
	tbl    table.Model
	fullW  int
	fullTH int
}

func newLootTable(h int) table.Model {
	columns := []table.Column{
		{Title: "Item", Width: 26},
		{Title: "Custo", Width: 12},
		{Title: "Detalhe", Width: 36},
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows([]table.Row{}),
		table.WithFocused(true),
		table.WithHeight(h),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true).
		Foreground(catBlue)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("#1e1e2e")).
		Background(catBlue).
		Bold(true)
	s.Cell = s.Cell.
		Foreground(catText)
	t.SetStyles(s)
	return t
}

func NewLootModel() LootModel {
	return LootModel{
		tbl: newLootTable(20),
	}
}

func (m *LootModel) SetTableSize(w, h int) {
	m.fullW = w
	m.fullTH = h - 5
	if m.fullTH < 5 {
		m.fullTH = 5
	}
	detailW := w - 26 - 12 - 4
	if detailW < 20 {
		detailW = 20
	}
	m.tbl.SetColumns([]table.Column{
		{Title: "Item", Width: 26},
		{Title: "Custo", Width: 12},
		{Title: "Detalhe", Width: detailW},
	})
	m.tbl.SetWidth(w)
	m.tbl.SetHeight(m.fullTH)
}

func formatEnchants(enchants []loot.Enchantment) string {
	if len(enchants) == 0 {
		return ""
	}
	var blessings, curses []string
	for _, e := range enchants {
		if e.Type == loot.Blessing {
			blessings = append(blessings, e.Title)
		} else {
			curses = append(curses, e.Title)
		}
	}
	var parts []string
	if len(blessings) > 0 {
		s := "✦ " + blessings[0]
		if len(blessings) > 1 {
			s += " +" + strconv.Itoa(len(blessings)-1)
		}
		parts = append(parts, s)
	}
	if len(curses) > 0 {
		s := "⚠ " + curses[0]
		if len(curses) > 1 {
			s += " +" + strconv.Itoa(len(curses)-1)
		}
		parts = append(parts, s)
	}
	return " · " + strings.Join(parts, " ")
}

func (m *LootModel) syncTable() {
	rows := make([]table.Row, len(m.items))
	for i, item := range m.items {
		switch v := item.(type) {
		case loot.Weapon:
			detail := fmt.Sprintf("%s %s · %s %s", v.Damage, v.DamageType, v.Category, v.Type)
			detail += formatEnchants(v.Enchantment)
			rows[i] = table.Row{v.Name, fmt.Sprintf("%d %s", v.Cost, v.Currency), detail}
		case loot.Armor:
			detail := fmt.Sprintf("CA +%d · %s", v.ArmorBonus, v.Type)
			detail += formatEnchants(v.Enchantment)
			rows[i] = table.Row{v.Name, fmt.Sprintf("%d %s", v.Cost, v.Currency), detail}
		case loot.Accessory:
			detail := v.Type
			detail += formatEnchants(v.Enchantment)
			rows[i] = table.Row{v.Name, fmt.Sprintf("%d %s", v.Cost, v.Currency), detail}
		case loot.PotionResult:
			detail := "Poção: " + v.Potion.Effect
			if v.ColateralEffect != nil {
				detail += " | ⚠ " + v.ColateralEffect.Description
			}
			rows[i] = table.Row{v.Potion.Name, fmt.Sprintf("%d %s", v.Potion.Cost, v.Potion.Currency), detail}
		}
	}
	m.tbl.SetRows(rows)
}

func (m *LootModel) Generate(n int) {
	if m.table == nil {
		t, err := loot.NewLootTable("data")
		if err != nil {
			m.err = err
			return
		}
		m.table = t
	}
	m.err = nil
	m.items = m.table.Generate(n)
	m.syncTable()
}

type LootItemPopupModel struct {
	item  any
	width int
}

func NewLootItemPopupModel(item any, w int) *LootItemPopupModel {
	return &LootItemPopupModel{item: item, width: w}
}

func (p *LootItemPopupModel) Init() tea.Cmd { return nil }

func (p *LootItemPopupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "enter":
			return p, func() tea.Msg { return PopupCloseMsg{} }
		}
	}
	return p, nil
}

func (p *LootItemPopupModel) View() string {
	var b strings.Builder

	switch v := p.item.(type) {
	case loot.Weapon:
		b.WriteString(popupTitle.Render(v.Name))
		b.WriteString("  ")
		b.WriteString(faint.Render(fmt.Sprintf("(%d %s)", v.Cost, v.Currency)))
		b.WriteString("\n\n")
		b.WriteString(statLabel.Render("Dano"))
		b.WriteString("  " + v.Damage + " " + string(v.DamageType) + "\n")
		b.WriteString(statLabel.Render("Categoria"))
		b.WriteString("  " + string(v.Category) + " · " + string(v.Type) + "\n")
		b.WriteString(renderEnchantDetail(v.Enchantment))

	case loot.Armor:
		b.WriteString(popupTitle.Render(v.Name))
		b.WriteString("  ")
		b.WriteString(faint.Render(fmt.Sprintf("(%d %s)", v.Cost, v.Currency)))
		b.WriteString("\n\n")
		b.WriteString(statLabel.Render("CA"))
		b.WriteString("  +" + strconv.Itoa(v.ArmorBonus) + "\n")
		b.WriteString(statLabel.Render("Tipo"))
		b.WriteString("  " + string(v.Type) + "\n")
		b.WriteString(renderEnchantDetail(v.Enchantment))

	case loot.Accessory:
		b.WriteString(popupTitle.Render(v.Name))
		b.WriteString("  ")
		b.WriteString(faint.Render(fmt.Sprintf("(%d %s)", v.Cost, v.Currency)))
		b.WriteString("\n\n")
		b.WriteString(statLabel.Render("Tipo"))
		b.WriteString("  " + v.Type + "\n")
		b.WriteString(renderEnchantDetail(v.Enchantment))

	case loot.PotionResult:
		b.WriteString(popupTitle.Render(v.Potion.Name))
		b.WriteString("  ")
		b.WriteString(faint.Render(fmt.Sprintf("(%d %s)", v.Potion.Cost, v.Potion.Currency)))
		b.WriteString("\n\n")
		b.WriteString(statLabel.Render("Efeito"))
		b.WriteString("  " + v.Potion.Effect + "\n")
		if v.ColateralEffect != nil {
			b.WriteString("\n")
			b.WriteString(red.Render("⚠ Efeito Colateral"))
			b.WriteString("\n  " + v.ColateralEffect.Description + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(faint.Render("[esc] fechar"))
	return b.String()
}

func renderEnchantDetail(enchants []loot.Enchantment) string {
	var b strings.Builder
	if len(enchants) == 0 {
		return ""
	}
	b.WriteString("\n")
	for _, e := range enchants {
		if e.Type == loot.Blessing {
			b.WriteString(green.Render("✨ Bênção: " + e.Title) + "\n")
		} else {
			b.WriteString(red.Render("☠️ Maldição: " + e.Title) + "\n")
		}
		b.WriteString("  " + e.Description + "\n")
	}
	b.WriteString("\n")
	return b.String()
}

func (m Model) updateLoot(msg tea.KeyMsg) Model {
	s := msg.String()

	switch s {
	case "0", "1", "2", "3", "4", "5", "6", "7", "8", "9":
		m.countBuf += s

	case "g":
		n := 5
		if m.countBuf != "" {
			if v, err := strconv.Atoi(m.countBuf); err == nil && v > 0 {
				n = v
			}
			if n > 100 {
				n = 100
			}
		}
		m.countBuf = ""
		m.loot.Generate(n)

	case "enter":
		if len(m.loot.items) > 0 {
			cursor := m.loot.tbl.Cursor()
			if cursor >= 0 && cursor < len(m.loot.items) {
				item := m.loot.items[cursor]
				outerW := (m.width - 6) * 60 / 100
				cw := outerW - 6
				pm := NewLootItemPopupModel(item, cw)
				m.popup = NewPopup(pm, 60, 60)
			}
		}

	default:
		if len(m.loot.items) > 0 {
			var cmd tea.Cmd
			m.loot.tbl, cmd = m.loot.tbl.Update(msg)
			_ = cmd
		}
		m.countBuf = ""
	}

	return m
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
		b.WriteString(m.loot.tbl.View())
	}

	return b.String()
}
