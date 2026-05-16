package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

func RunDice(args []string) {
	if len(args) == 0 {
		fmt.Println("use: rpgctl dice d20")
		return
	}

	result, err := dice(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}

func dice(arg string) (int, error) {
	if !strings.Contains(arg, "d") {
		return 0, fmt.Errorf("Formato de dado invalido: %s", arg)
	}

	parts := strings.Split(arg, "d")

	quantityText := parts[0]

	if quantityText == "" {
		quantityText = "1"
	}

	quantity, err := strconv.Atoi(quantityText)
	if err != nil {
		return 0, err
	}

	if quantity <= 0 {
		return 0, fmt.Errorf("quantidade de dados invalida: %d", quantity)
	}

	if parts[1] == "" {
		return 0, fmt.Errorf("Formato de dado invalido: %s", arg)
	}

	expr := strings.ReplaceAll(parts[1], "-", "+-")
	terms := strings.Split(expr, "+")

	sidesText := terms[0]

	modifierTexts := terms[1:]

	modifier := 0
	for _, text := range modifierTexts {
		value, err := strconv.Atoi(text)
		if err != nil {
			return 0, err
		}
		modifier += value
	}

	sides, err := strconv.Atoi(sidesText)
	if err != nil {
		return 0, err
	}

	switch sides {
	case 2, 4, 6, 8, 10, 12, 20, 100:
		total := 0
		for range quantity {
			total += rand.IntN(sides) + 1
		}
		return total + modifier, nil
	default:
		return 0, fmt.Errorf("formato de dado nao permitido: d%d", sides)
	}
}
