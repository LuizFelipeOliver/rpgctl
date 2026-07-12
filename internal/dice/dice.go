package dice

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
)

const (
	D4   = 4
	D6   = 6
	D8   = 8
	D10  = 10
	D12  = 12
	D20  = 20
	D100 = 100
)

type Group struct {
	Count int
	Sides int
	Rolls []int
}

type RollResult struct {
	Expression string
	Groups     []Group
	Modifier   int
	Total      int
}

func Roll(expr string) (*RollResult, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, errors.New("expressao vazia")
	}

	groups, modifier, err := parseExpr(expr)
	if err != nil {
		return nil, err
	}

	result := &RollResult{
		Expression: expr,
		Groups:     groups,
		Modifier:   modifier,
	}

	for i := range result.Groups {
		g := &result.Groups[i]
		for range g.Count {
			roll := rand.Intn(g.Sides) + 1
			g.Rolls = append(g.Rolls, roll)
			result.Total += roll
		}
	}
	result.Total += modifier

	return result, nil
}

func (r *RollResult) String() string {
	parts := make([]string, 0, len(r.Groups))

	for _, g := range r.Groups {
		rolls := make([]string, len(g.Rolls))
		for i, v := range g.Rolls {
			rolls[i] = strconv.Itoa(v)
		}
		text := strings.Join(rolls, "+")

		if g.Count > 1 || len(r.Groups) > 1 || r.Modifier != 0 {
			text = "(" + text + ")"
		}

		if len(parts) > 0 {
			parts = append(parts, "+ "+text)
		} else {
			parts = append(parts, text)
		}
	}

	if r.Modifier != 0 {
		if len(parts) > 0 {
			if r.Modifier > 0 {
				parts = append(parts, "+ "+strconv.Itoa(r.Modifier))
			} else {
				parts = append(parts, "- "+strconv.Itoa(-r.Modifier))
			}
		} else {
			parts = append(parts, strconv.Itoa(r.Modifier))
		}
	}

	return strings.Join(parts, " ") + " = " + strconv.Itoa(r.Total)
}

func parseExpr(s string) ([]Group, int, error) {
	pos := 0

	count, sides, newPos, ok := parseCountSides(s, pos)
	if !ok {
		return nil, 0, errors.New("expressao deve comecar com um dado (ex: d20, 2d6)")
	}
	groups := []Group{{Count: count, Sides: sides}}
	modifier := 0
	pos = newPos

	for pos < len(s) {
		var err error
		var sign int
		sign, pos, err = parseSign(s, pos)
		if err != nil {
			return nil, 0, err
		}

		count, sides, newPos, ok = parseCountSides(s, pos)
		if ok {
			if sign < 0 {
				return nil, 0, errors.New("dado negativo nao permitido")
			}
			groups = append(groups, Group{Count: count, Sides: sides})
			pos = newPos
		} else {
			var num int
			num, pos = parseNumber(s, pos)
			modifier += sign * num
		}
	}

	return groups, modifier, nil
}

func parseCountSides(s string, pos int) (count, sides, newPos int, ok bool) {
	count, newPos = parseNumber(s, pos)
	if newPos >= len(s) || s[newPos] != 'd' {
		return 0, 0, pos, false
	}
	newPos++
	sides, newPos = parseNumber(s, newPos)
	if sides < 2 {
		sides = 2
	}
	return count, sides, newPos, true
}

func parseNumber(s string, pos int) (val, newPos int) {
	if pos >= len(s) || !unicode.IsDigit(rune(s[pos])) {
		return 1, pos
	}
	end := pos
	for end < len(s) && unicode.IsDigit(rune(s[end])) {
		end++
	}
	n, _ := strconv.Atoi(s[pos:end])
	if n < 1 {
		n = 1
	}
	return n, end
}

func parseSign(s string, pos int) (sign int, newPos int, err error) {
	if s[pos] == '+' {
		return 1, pos + 1, nil
	}
	if s[pos] == '-' {
		return -1, pos + 1, nil
	}
	return 0, pos, fmt.Errorf("caractere inesperado: %c", s[pos])
}
