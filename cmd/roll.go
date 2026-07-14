package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"rpg-tui/internal/dice"
)

var rollCmd = &cobra.Command{
	Use:                "roll <expressao>",
	Short:              "Rola dados. Ex: \"d20 + 1 + d3\", \"2d6\"",
	DisableFlagParsing: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, a := range args {
			if a == "--help" || a == "-h" {
				return cmd.Help()
			}
		}

		expr := strings.Join(args, " ")
		result, err := dice.Rolar(expr)
		if err != nil {
			return fmt.Errorf("erro: %w", err)
		}
		fmt.Printf("%s = %d\n", result.Detalhes, result.Total)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(rollCmd)
}
