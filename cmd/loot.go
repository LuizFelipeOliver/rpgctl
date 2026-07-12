package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"rpg-tui/internal/loot"
)

var lootCmd = &cobra.Command{
	Use:   "loot [n]",
	Short: "Gera itens aleatórios",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		n := 1
		if len(args) > 0 {
			var err error
			n, err = strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("número inválido: %s", args[0])
			}
		}

		lt, err := loot.NewLootTable("data")
		if err != nil {
			return fmt.Errorf("erro carregando dados: %w", err)
		}

		items := lt.Generate(n)
		fmt.Print(loot.DisplayItems(items))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(lootCmd)
}
