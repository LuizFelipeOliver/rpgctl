package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"rpg-tui/internal/encounter"
	"rpg-tui/internal/monster"
)

var encounterCmd = &cobra.Command{
	Use:   "encounter <jogadores> <nivel> [dificuldade]",
	Short: "Gera um encontro aleatório",
	Long: `Gera um encontro balanceado para um grupo de jogadores.

Dificuldades: F (Fácil), M (Médio, padrão), D (Difícil)

Exemplos:
  rpgctl encounter 4 5
  rpgctl encounter 4 5 D
  rpgctl encounter 3 1 F`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		players, err := parsePosInt(args[0])
		if err != nil {
			return fmt.Errorf("numero de jogadores invalido: %w", err)
		}
		level, err := parsePosInt(args[1])
		if err != nil {
			return fmt.Errorf("nivel invalido: %w", err)
		}

		diff := "medio"
		if len(args) > 2 {
			diff = args[2]
		}

		monsters, err := monster.LoadMonsters("data/monster/monsters-parsed.json")
		if err != nil {
			return fmt.Errorf("erro carregando monstros: %w", err)
		}

		r, err := encounter.Generate(monsters, players, level, diff)
		if err != nil {
			return fmt.Errorf("erro gerando encontro: %w", err)
		}

		fmt.Print(encounter.DisplayResult(r))
		return nil
	},
}

func parsePosInt(s string) (int, error) {
	var n int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("'%s' nao e um numero valido", s)
		}
		n = n*10 + int(c-'0')
	}
	if n <= 0 {
		return 0, fmt.Errorf("'%s' precisa ser positivo", s)
	}
	return n, nil
}

func init() {
	rootCmd.AddCommand(encounterCmd)
}
