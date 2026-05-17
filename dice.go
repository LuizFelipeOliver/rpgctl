package main

import (
	"fmt"
	"math/rand/v2"
	"strconv"
	"strings"
)

type Dice struct {
	Quantity int
	Sides    int
	Modifier int
}

func RunDice(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("use: rpgctl dice d20")

	}

	result, err := dice(args[0])
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}

func dice(arg string) (int, error) {
	d, err := parseDice(arg)
	if err != nil {
		return 0, err
	}

	return rollDice(d), nil
}

const maxDiceQuantity = 100

func parseDice(arg string) (Dice, error) {
	if !strings.Contains(arg, "d") {
		return Dice{}, fmt.Errorf("formato de dado invalido: %s", arg)
	}

	if strings.Count(arg, "d") != 1 {
		return Dice{}, fmt.Errorf("formato de dado invalido: %s", arg)
	}

	parts := strings.Split(arg, "d")

	if parts[1] == "" {
		return Dice{}, fmt.Errorf("formato de dado invalido: %s", arg)
	}

	quantityText := parts[0]

	if quantityText == "" {
		quantityText = "1"
	}

	quantity, err := strconv.Atoi(quantityText)
	if err != nil {
		return Dice{}, err
	}

	if quantity <= 0 {
		return Dice{}, fmt.Errorf("quantidade de dados invalida: %d", quantity)
	}

	if quantity > maxDiceQuantity {
		return Dice{}, fmt.Errorf("exedido quantidade de dados: %d (max:%d)", quantity, maxDiceQuantity)
	}

	expr := strings.ReplaceAll(parts[1], "-", "+-")
	terms := strings.Split(expr, "+")

	sidesText := terms[0]

	sides, err := parseSides(sidesText)
	if err != nil {
		return Dice{}, err
	}

	modifier, err := parseModifier(terms[1:])
	if err != nil {
		return Dice{}, err
	}

	return Dice{
		Quantity: quantity,
		Sides:    sides,
		Modifier: modifier,
	}, nil
}

func parseModifier(arg []string) (int, error) {
	modifier := 0
	for _, text := range arg {
		value, err := strconv.Atoi(text)
		if err != nil {
			return 0, err
		}
		modifier += value
	}
	return modifier, nil
}

func parseSides(sidesText string) (int, error) {
	sides, err := strconv.Atoi(sidesText)
	if err != nil {
		return 0, err
	}

	switch sides {
	case 2, 4, 6, 8, 10, 12, 20, 100:
		return sides, nil
	default:
		return 0, fmt.Errorf("quantidade de faces invalida: %d", sides)
	}
}

func rollDice(d Dice) int {
	total := 0
	for range d.Quantity {
		total += rand.IntN(d.Sides) + 1
	}
	return total + d.Modifier
}
