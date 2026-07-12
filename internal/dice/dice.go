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
	p := &parser{raw: s}

	count, sides, ok := p.readDice()
	if !ok {
		return nil, 0, errors.New("expressao deve comecar com um dado (ex: d20, 2d6)")
	}

	groups := []Group{{Count: count, Sides: sides}}
	modifier := 0

	for p.next() {
		sign := p.readSign()
		if p.err != nil {
			return nil, 0, p.err
		}

		if count, sides, ok = p.readDice(); ok {
			if sign < 0 {
				return nil, 0, errors.New("dado negativo nao permitido")
			}
			groups = append(groups, Group{Count: count, Sides: sides})
		} else {
			modifier += sign * p.readNumber()
		}
	}

	return groups, modifier, nil
}

type parser struct {
	raw string
	pos int
	err error
}

func (p *parser) next() bool {
	return p.err == nil && p.pos < len(p.raw)
}

func (p *parser) readSign() int {
	switch p.raw[p.pos] {
	case '+':
		p.pos++
		return 1
	case '-':
		p.pos++
		return -1
	default:
		p.err = fmt.Errorf("caractere inesperado: %c", p.raw[p.pos])
		return 0
	}
}

func (p *parser) readNumber() int {
	start := p.pos
	for p.pos < len(p.raw) && unicode.IsDigit(rune(p.raw[p.pos])) {
		p.pos++
	}
	if start == p.pos {
		return 1
	}
	n, _ := strconv.Atoi(p.raw[start:p.pos])
	if n < 1 {
		return 1
	}
	return n
}

func (p *parser) readDice() (count, sides int, ok bool) {
	saved := p.pos
	count = p.readNumber()
	if !p.next() || p.raw[p.pos] != 'd' {
		p.pos = saved
		return 0, 0, false
	}
	p.pos++
	sides = p.readNumber()
	if sides < 2 {
		sides = 2
	}
	return count, sides, true
}
