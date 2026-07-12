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
	Sign  int
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

	p := &parser{expr: expr}
	groups, modifier := p.parse()

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
			result.Total += g.Sign * roll
		}
	}
	result.Total += modifier

	return result, nil
}

func (r *RollResult) String() string {
	parts := make([]string, 0, len(r.Groups))

	for _, g := range r.Groups {
		rollStrs := make([]string, len(g.Rolls))
		for i, v := range g.Rolls {
			rollStrs[i] = strconv.Itoa(v)
		}
		groupText := strings.Join(rollStrs, "+")

		if g.Sign < 0 || g.Count > 1 || len(r.Groups) > 1 || r.Modifier != 0 {
			groupText = "(" + groupText + ")"
		}

		if g.Sign < 0 {
			if len(parts) == 0 {
				parts = append(parts, "-"+groupText)
			} else {
				parts = append(parts, "- "+groupText)
			}
		} else if len(parts) > 0 {
			parts = append(parts, "+ "+groupText)
		} else {
			parts = append(parts, groupText)
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

	if len(parts) == 0 {
		return fmt.Sprintf("%d = %d", r.Modifier, r.Total)
	}

	return strings.Join(parts, " ") + " = " + strconv.Itoa(r.Total)
}

type parser struct {
	expr string
	pos  int
}

func (p *parser) parse() (groups []Group, modifier int) {
	for p.pos < len(p.expr) {
		sign := 1
		if p.expr[p.pos] == '+' {
			sign = 1
			p.pos++
		} else if p.expr[p.pos] == '-' {
			sign = -1
			p.pos++
		}

		start := p.pos
		num := p.parseNumber()

		if p.pos < len(p.expr) && p.expr[p.pos] == 'd' {
			if start == p.pos {
				num = 1
			}
			p.pos++
			sides := p.parseNumber()
			if sides < 2 {
				sides = 2
			}
			groups = append(groups, Group{
				Count: num,
				Sides: sides,
				Sign:  sign,
			})
		} else {
			modifier += sign * num
		}
	}
	return
}

func (p *parser) parseNumber() int {
	start := p.pos
	for p.pos < len(p.expr) && unicode.IsDigit(rune(p.expr[p.pos])) {
		p.pos++
	}
	if start == p.pos {
		return 1
	}
	n, _ := strconv.Atoi(p.expr[start:p.pos])
	if n < 1 {
		return 1
	}
	return n
}
