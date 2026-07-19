package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"rpg-tui/internal/monster"
)

type PopupCloseMsg struct{}

type AddToInitMsg struct {
	Monster *monster.Monster
}

type HealApplyMsg struct {
	Cursor   int
	Amount   int
	IsDamage bool
}

type ResetConfirmMsg struct{}

type WizardCompleteMsg struct {
	Name       string
	Initiative int
	HP         int
	AC         int
	Monster    *monster.Monster
}

type contentWidthSetter interface {
	SetContentWidth(int)
}

type Popup struct {
	Inner     tea.Model
	WidthPct  int
	HeightPct int
	w, h      int
}

func NewPopup(inner tea.Model, wp, hp int) *Popup {
	return &Popup{Inner: inner, WidthPct: wp, HeightPct: hp}
}

func (p *Popup) SetSize(w, h int) {
	p.w = w
	p.h = h
}

func (p *Popup) Init() tea.Cmd {
	if p.Inner == nil {
		return nil
	}
	return p.Inner.Init()
}

func (p *Popup) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if p.Inner == nil {
		return p, nil
	}

	if _, ok := msg.(tea.WindowSizeMsg); ok {
		pw := p.w * p.WidthPct / 100
		cw := pw - 6
		if sw, ok := p.Inner.(contentWidthSetter); ok {
			sw.SetContentWidth(cw)
		}
	}

	var cmd tea.Cmd
	p.Inner, cmd = p.Inner.Update(msg)
	return p, cmd
}

func (p *Popup) View() string {
	if p.Inner == nil {
		return ""
	}
	pw := p.w * p.WidthPct / 100
	cw := pw - 6
	if sw, ok := p.Inner.(contentWidthSetter); ok {
		sw.SetContentWidth(cw)
	}
	return lipgloss.Place(p.w, p.h, lipgloss.Center, lipgloss.Center, popupStyle.Render(p.Inner.View()))
}
