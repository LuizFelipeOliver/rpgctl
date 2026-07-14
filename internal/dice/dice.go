package dice

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"unicode"
)

const (D4, D6, D8, D10, D12, D20, D100 = 4, 6, 8, 10, 12, 20, 100)

type Group struct {
	Count, Sides int
	Rolls        []int
	Sign         int
}

type RollResult struct {
	Expression      string
	Groups          []Group
	Modifier, Total int
}

func Roll(expr string) (*RollResult, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return nil, errors.New("expressao vazia")
	}

	groups, mod, err := parse(expr)
	if err != nil {
		return nil, err
	}

	r := &RollResult{Expression: expr, Groups: groups, Modifier: mod}
	for i := range r.Groups {
		g := &r.Groups[i]
		for range g.Count {
			v := rand.Intn(g.Sides) + 1
			g.Rolls = append(g.Rolls, v)
			r.Total += g.Sign * v
		}
	}
	r.Total += mod
	return r, nil
}

func (r *RollResult) String() string {
	parts := make([]string, 0, len(r.Groups))
	for _, g := range r.Groups {
		rolls := make([]string, len(g.Rolls))
		for i, v := range g.Rolls {
			rolls[i] = strconv.Itoa(v)
		}
		s := strings.Join(rolls, "+")

		if g.Sign < 0 || g.Count > 1 || len(r.Groups) > 1 || r.Modifier != 0 {
			s = "(" + s + ")"
		}
		if g.Sign < 0 {
			if len(parts) == 0 {
				parts = append(parts, "-"+s)
			} else {
				parts = append(parts, "- "+s)
			}
		} else if len(parts) > 0 {
			parts = append(parts, "+ "+s)
		} else {
			parts = append(parts, s)
		}
	}
	if r.Modifier != 0 {
		if r.Modifier > 0 {
			parts = append(parts, "+ "+strconv.Itoa(r.Modifier))
		} else {
			parts = append(parts, "- "+strconv.Itoa(-r.Modifier))
		}
	}
	return strings.Join(parts, " ") + " = " + strconv.Itoa(r.Total)
}

func parse(s string) ([]Group, int, error) {
	sign := 1
	pos := 0
	first := true

	var groups []Group
	mod := 0

	for pos < len(s) {
		if first {
			if s[pos] == '+' {
				pos++
			} else if s[pos] == '-' {
				sign = -1
				pos++
			}
		} else {
			switch s[pos] {
			case '+':
				sign = 1
				pos++
			case '-':
				sign = -1
				pos++
			default:
				return nil, 0, fmt.Errorf("caractere inesperado: %c", s[pos])
			}
		}

		start := pos
		for pos < len(s) && unicode.IsDigit(rune(s[pos])) {
			pos++
		}
		n := 1
		if pos > start {
			n, _ = strconv.Atoi(s[start:pos])
		}

		if pos < len(s) && s[pos] == 'd' {
			pos++
			start = pos
			for pos < len(s) && unicode.IsDigit(rune(s[pos])) {
				pos++
			}
			sides := 1
			if pos > start {
				sides, _ = strconv.Atoi(s[start:pos])
			}
			groups = append(groups, Group{Count: n, Sides: sides, Sign: sign})
			first = false
		} else if first {
			return nil, 0, errors.New("expressao deve comecar com um dado (ex: d20, 2d6)")
		} else {
			mod += sign * n
		}
	}

	if first {
		return nil, 0, errors.New("expressao vazia")
	}
	return groups, mod, nil
}
