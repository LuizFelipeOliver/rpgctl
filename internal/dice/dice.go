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
	Expression    string
	Groups        []Group
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
		s := joinInts(g.Rolls, "+")
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
	pos := 0
	groups := []Group{}
	mod := 0

	sign := 1
	if pos < len(s) && s[pos] == '+' { pos++ }
	if pos < len(s) && s[pos] == '-' { sign = -1; pos++ }

	c, sides, ok := dice(s, &pos)
	if !ok {
		return nil, 0, errors.New("expressao deve comecar com um dado (ex: d20, 2d6)")
	}
	groups = append(groups, Group{Count: c, Sides: sides, Sign: sign})

	for pos < len(s) {
		switch s[pos] {
		case '+': sign = 1; pos++
		case '-': sign = -1; pos++
		default: return nil, 0, fmt.Errorf("caractere inesperado: %c", s[pos])
		}

		if c, sides, ok = dice(s, &pos); ok {
			groups = append(groups, Group{Count: c, Sides: sides, Sign: sign})
		} else {
			mod += sign * digits(s, &pos)
		}
	}
	return groups, mod, nil
}

func dice(s string, pos *int) (int, int, bool) {
	saved := *pos
	c := digits(s, pos)
	if *pos >= len(s) || s[*pos] != 'd' {
		*pos = saved
		return 0, 0, false
	}
	*pos++
	return c, max(digits(s, pos), 1), true
}

func digits(s string, pos *int) int {
	start := *pos
	for *pos < len(s) && unicode.IsDigit(rune(s[*pos])) {
		*pos++
	}
	if start == *pos {
		return 1
	}
	return max(1, atoi(s[start:*pos]))
}

func atoi(s string) int {
	n, _ := strconv.Atoi(s)
	return n
}

func joinInts(v []int, sep string) string {
	s := make([]string, len(v))
	for i, n := range v {
		s[i] = strconv.Itoa(n)
	}
	return strings.Join(s, sep)
}
