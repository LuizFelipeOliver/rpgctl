package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"rpg-tui/internal/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Inicia a interface interativa no terminal",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := tea.NewProgram(tui.New())
		_, err := p.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
