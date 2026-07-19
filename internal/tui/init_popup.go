package tui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"rpg-tui/internal/monster"
)

var errInvalidChar = &invalidCharError{}

type invalidCharError struct{}

func (e *invalidCharError) Error() string { return "invalid character" }

func newWizardInput(placeholder string, validate func(string) error) textinput.Model {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Validate = validate
	return ti
}

func initValidator(s string) error {
	if s == "" {
		return nil
	}
	for i, r := range s {
		if r >= '0' && r <= '9' {
			continue
		}
		if i == 0 && (r == '-' || r == '+') {
			continue
		}
		return errInvalidChar
	}
	return nil
}

func digitOptional(s string) error {
	if s == "" {
		return nil
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return errInvalidChar
		}
	}
	return nil
}

type WizardPopupModel struct {
	inputs  [4]textinput.Model
	step    int
	monster *monster.Monster
}

func NewWizardPopupModel() *WizardPopupModel {
	return &WizardPopupModel{
		inputs: [4]textinput.Model{
			newWizardInput("Name...", nil),
			newWizardInput("Initiative...", initValidator),
			newWizardInput("HP (optional)...", digitOptional),
			newWizardInput("AC (optional)...", digitOptional),
		},
		step: 1,
	}
}

func (w *WizardPopupModel) Prefill(mon *monster.Monster) {
	w.monster = mon
	w.inputs[0].SetValue(mon.Name)
	w.inputs[1].SetValue("")
	w.inputs[2].SetValue(strconv.Itoa(mon.HP))
	w.inputs[3].SetValue(strconv.Itoa(firstAC(mon.ArmorClass)))
	w.step = 2
	w.inputs[0].Blur()
	w.inputs[1].Focus()
}

func (w *WizardPopupModel) Init() tea.Cmd { return nil }

func (w *WizardPopupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return w, func() tea.Msg { return PopupCloseMsg{} }

		case "enter":
			switch w.step {
			case 1:
				if strings.TrimSpace(w.inputs[0].Value()) != "" {
					w.inputs[0].Blur()
					w.inputs[1].Focus()
					w.step = 2
				}
			case 2:
				if strings.TrimSpace(w.inputs[1].Value()) != "" {
					w.inputs[1].Blur()
					w.inputs[2].Focus()
					w.step = 3
				}
			case 3:
				w.inputs[2].Blur()
				w.inputs[3].Focus()
				w.step = 4
			case 4:
				name := strings.TrimSpace(w.inputs[0].Value())
				if name == "" {
					return w, nil
				}
				init, _ := strconv.Atoi(strings.TrimSpace(w.inputs[1].Value()))
				hp, _ := strconv.Atoi(strings.TrimSpace(w.inputs[2].Value()))
				ac, _ := strconv.Atoi(strings.TrimSpace(w.inputs[3].Value()))
				return w, func() tea.Msg {
					return WizardCompleteMsg{
						Name: name, Initiative: init, HP: hp, AC: ac, Monster: w.monster,
					}
				}
			}
			return w, nil

		default:
			focused := w.step - 1
			if focused >= 0 && focused < 4 {
				var cmd tea.Cmd
				w.inputs[focused], cmd = w.inputs[focused].Update(msg)
				_ = cmd
			}
		}
	}
	return w, nil
}

func (w *WizardPopupModel) View() string {
	var b strings.Builder
	b.WriteString(popupTitle.Render("Add Combatant"))
	b.WriteString("\n\n")

	labels := []string{"Name", "Initiative", "HP (optional)", "AC (optional)"}

	for i := 0; i < 4; i++ {
		step := i + 1
		if step == w.step {
			b.WriteString(green.Render("> " + w.inputs[i].View()))
		} else if step < w.step {
			b.WriteString(faint.Render("  " + w.inputs[i].Value()))
		} else {
			b.WriteString(faint.Render("  " + labels[i]))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(faint.Render("Enter: next  •  esc: cancel"))
	return b.String()
}

type HealPopupModel struct {
	healInput textinput.Model
	isDamage  bool
	cursor    int
}

func NewHealPopupModel(cursor int, isDamage bool) *HealPopupModel {
	hi := textinput.New()
	hi.Placeholder = "0"
	hi.Validate = digitOptional
	hi.Focus()
	return &HealPopupModel{
		healInput: hi,
		isDamage:  isDamage,
		cursor:    cursor,
	}
}

func (h *HealPopupModel) Init() tea.Cmd {
	return textinput.Blink
}

func (h *HealPopupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			h.healInput.Blur()
			return h, func() tea.Msg { return PopupCloseMsg{} }

		case "enter":
			val, err := strconv.Atoi(h.healInput.Value())
			if err != nil || val <= 0 {
				return h, nil
			}
			h.healInput.Blur()
			return h, func() tea.Msg {
				return HealApplyMsg{Cursor: h.cursor, Amount: val, IsDamage: h.isDamage}
			}
		}
	}

	var cmd tea.Cmd
	h.healInput, cmd = h.healInput.Update(msg)
	return h, cmd
}

func (h *HealPopupModel) View() string {
	label := "Dano"
	style := red
	if !h.isDamage {
		label = "Cura"
		style = green
	}

	var content strings.Builder
	content.WriteString(popupTitle.Render(label))
	content.WriteString("\n\n")
	content.WriteString(style.Render(h.healInput.View()))
	content.WriteString("\n\n")
	content.WriteString(faint.Render("Enter: aplicar  •  esc: cancelar"))
	return content.String()
}

type ResetPopupModel struct{}

func NewResetPopupModel() *ResetPopupModel {
	return &ResetPopupModel{}
}

func (r *ResetPopupModel) Init() tea.Cmd { return nil }

func (r *ResetPopupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "n", "N":
			return r, func() tea.Msg { return PopupCloseMsg{} }
		case "enter", "s", "S", "ctrl+r":
			return r, func() tea.Msg { return ResetConfirmMsg{} }
		}
	}
	return r, nil
}

func (r *ResetPopupModel) View() string {
	var b strings.Builder
	b.WriteString(popupTitle.Render("Resetar Iniciativa"))
	b.WriteString("\n\n")
	b.WriteString("Tem certeza que deseja\nremover todos os combatentes?\n\n")
	b.WriteString(green.Render("[Enter/S] Confirmar"))
	b.WriteString("  ")
	b.WriteString(faint.Render("[Esc/N] Cancelar"))
	return b.String()
}

type InitDetailPopupModel struct {
	combatant *Combatant
	monster   *monster.Monster
	width     int
}

func NewInitDetailPopupModel(combatant *Combatant, w int) *InitDetailPopupModel {
	return &InitDetailPopupModel{
		combatant: combatant,
		monster:   combatant.Monster,
		width:     w,
	}
}

func (d *InitDetailPopupModel) Init() tea.Cmd { return nil }

func (d *InitDetailPopupModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "enter":
			return d, func() tea.Msg { return PopupCloseMsg{} }
		}
	}
	return d, nil
}

func (d *InitDetailPopupModel) View() string {
	var b strings.Builder

	if d.combatant != nil {
		c := d.combatant
		b.WriteString(lootTitle.Render(c.Name))
		b.WriteString("\n")
		b.WriteString(faint.Render(strings.Repeat("─", 50)))
		b.WriteString("\n\n")

		if d.monster != nil {
			writeMonsterStats(&b, *d.monster, d.width)
		} else {
			writeField(&b, "Initiative", strconv.Itoa(c.Initiative))
			hpStr := fmt.Sprintf("%d/%d", c.HP, c.MaxHP)
			if c.MaxHP == 0 {
				hpStr = "-"
			}
			writeField(&b, "HP", hpStr)
			if c.AC > 0 {
				writeField(&b, "AC", strconv.Itoa(c.AC))
			}
		}
	} else if d.monster != nil {
		mon := *d.monster
		b.WriteString(lootTitle.Render(mon.Name))
		b.WriteString("\n")
		b.WriteString(faint.Render(strings.Repeat("─", 50)))
		b.WriteString("\n\n")
		writeMonsterStats(&b, mon, d.width)
	}

	b.WriteString("\n")
	b.WriteString(faint.Render("[esc] back"))
	return b.String()
}
