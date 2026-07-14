package dice

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func RollDie(sides int) (int, error) {
	if sides < 1 {
		return 0, errors.New("numero de faces invalido")
	}
	return rand.Intn(sides) + 1, nil
}

type RolarResultado struct {
	Detalhes string
	Total    int
}

func Rolar(expressao string) (*RolarResultado, error) {
	expressao = strings.TrimSpace(expressao)
	if expressao == "" {
		return nil, errors.New("expressao vazia")
	}

	tokens := strings.Fields(expressao)
	if len(tokens) == 0 {
		return nil, errors.New("expressao vazia")
	}

	var detalhes []string
	total := 0
	sinal := 1

	for _, tok := range tokens {
		switch tok {
		case "+":
			sinal = 1
		case "-":
			sinal = -1
		default:
			var v int
			if strings.Contains(tok, "d") {
				var err error
				v, err = RollNotation(tok)
				if err != nil {
					return nil, err
				}
			} else {
				n, err := strconv.Atoi(tok)
				if err != nil {
					return nil, fmt.Errorf("token invalido: %s", tok)
				}
				v = n
			}

			if sinal < 0 {
				detalhes = append(detalhes, "- "+strconv.Itoa(v))
			} else if len(detalhes) > 0 {
				detalhes = append(detalhes, "+ "+strconv.Itoa(v))
			} else {
				detalhes = append(detalhes, strconv.Itoa(v))
			}
			total += sinal * v
			sinal = 1
		}
	}

	return &RolarResultado{
		Detalhes: strings.Join(detalhes, " "),
		Total:    total,
	}, nil
}

func RollNotation(notacao string) (int, error) {
	notacao = strings.TrimSpace(notacao)
	if notacao == "" {
		return 0, errors.New("notacao vazia")
	}

	parts := strings.SplitN(notacao, "d", 2)
	if len(parts) != 2 {
		return 0, errors.New("formato invalido, use NdS (ex: 2d6, d20)")
	}

	quantityDice := 1
	if parts[0] != "" {
		n, err := strconv.Atoi(parts[0])
		if err != nil || n < 1 {
			return 0, errors.New("numero de dados invalido")
		}
		quantityDice = n
	}

	sides, err := strconv.Atoi(parts[1])
	if err != nil || sides < 1 {
		return 0, errors.New("numero de faces invalido")
	}

	total := 0
	for range quantityDice {
		v, err := RollDie(sides)
		if err != nil {
			return 0, err
		}
		total += v
	}
	return total, nil
}
