package encounter

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"rpg-tui/internal/monster"
)

type Group struct {
	Monster  monster.Monster
	Quantity int
}

type Result struct {
	Groups     []Group
	Difficulty string
	TargetEL   int
	PartyLevel int
	PartyCount int
}

func CRToFloat(cr string) (float64, error) {
	cr = strings.TrimSpace(cr)
	if cr == "" {
		return 0, fmt.Errorf("CR vazio")
	}
	if idx := strings.IndexByte(cr, '/'); idx >= 0 {
		num, err := strconv.ParseFloat(cr[:idx], 64)
		if err != nil {
			return 0, fmt.Errorf("CR invalido: %s", cr)
		}
		den, err := strconv.ParseFloat(cr[idx+1:], 64)
		if err != nil {
			return 0, fmt.Errorf("CR invalido: %s", cr)
		}
		if den == 0 {
			return 0, fmt.Errorf("CR invalido: divisao por zero")
		}
		return num / den, nil
	}
	v, err := strconv.ParseFloat(cr, 64)
	if err != nil {
		return 0, fmt.Errorf("CR invalido: %s", cr)
	}
	return v, nil
}

func CRToInt(cr string) int {
	v, err := CRToFloat(cr)
	if err != nil {
		return 0
	}
	return int(v)
}

func groupEL(cr, qty int) int {
	switch {
	case qty <= 0:
		return 0
	case qty == 1:
		return cr
	case qty == 2:
		return cr + 2
	case qty <= 4:
		return cr + 3
	case qty <= 7:
		return cr + 4
	case qty <= 11:
		return cr + 5
	default:
		return cr + 6
	}
}

func combineEL(els []int) int {
	if len(els) == 0 {
		return 0
	}
	maxEL := els[0]
	for _, el := range els[1:] {
		switch {
		case maxEL == el:
			maxEL += 2
		case maxEL-el <= 2 && maxEL-el >= -2:
			maxEL++
		}
	}
	return maxEL
}

func pickAtCR(monsters []monster.Monster, target int) *monster.Monster {
	var pool []monster.Monster
	for _, m := range monsters {
		cr := CRToInt(m.ChallengeRating)
		if cr == target || (target < 0 && cr == 0) || (target == 0 && cr == 0) {
			pool = append(pool, m)
		}
	}
	if len(pool) == 0 {
		return nil
	}
	return &pool[rand.Intn(len(pool))]
}

func pickLowestCR(monsters []monster.Monster) *monster.Monster {
	if len(monsters) == 0 {
		return nil
	}
	best := &monsters[0]
	bestVal := CRToInt(best.ChallengeRating)
	for i := range monsters {
		if v := CRToInt(monsters[i].ChallengeRating); v < bestVal {
			best = &monsters[i]
			bestVal = v
		}
	}
	return best
}

var diffMap = map[string]int{
	"F": -2,
	"M": 0,
	"D": 1,
}

type GenerateOptions struct {
	TypeFilter string
	Quantity   int
}

func Generate(all []monster.Monster, players, level int, diff string) (*Result, error) {
	return GenerateWithOpts(all, players, level, diff, GenerateOptions{})
}

func GenerateWithOpts(all []monster.Monster, players, level int, diff string, opts GenerateOptions) (*Result, error) {
	diff = strings.ToUpper(strings.TrimSpace(diff))
	offset, ok := diffMap[diff]
	if !ok {
		return nil, fmt.Errorf("dificuldade invalida: %s (use: F (Facil), M (Medio), D (Dificil))", diff)
	}

	target := level + offset
	if target < 0 {
		target = 0
	}

	filtered := all
	if opts.TypeFilter != "" {
		ft := strings.ToLower(strings.TrimSpace(opts.TypeFilter))
		var f []monster.Monster
		for _, m := range filtered {
			if strings.Contains(strings.ToLower(m.Type), ft) {
				f = append(f, m)
			}
		}
		if len(f) == 0 {
			return nil, fmt.Errorf("nenhum monstro do tipo '%s' encontrado", opts.TypeFilter)
		}
		filtered = f
	}

	var candidates []monster.Monster
	for _, m := range filtered {
		if CRToInt(m.ChallengeRating) <= target+1 {
			candidates = append(candidates, m)
		}
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("nenhum monstro encontrado para EL ~%d", target)
	}

	if opts.Quantity > 0 {
		return generateFixed(candidates, target, opts.Quantity, diff, level, players), nil
	}

	leader := pickAtCR(candidates, target-1)
	if leader == nil {
		leader = pickAtCR(candidates, target)
	}
	if leader == nil {
		leader = pickAtCR(candidates, target-2)
	}
	if leader == nil {
		leader = &candidates[rand.Intn(len(candidates))]
	}

	groups := []Group{{Monster: *leader, Quantity: 1}}
	leaderEL := groupEL(CRToInt(leader.ChallengeRating), 1)

	minionCR := target - 3
	if minionCR < 0 {
		minionCR = 0
	}
	minion := pickAtCR(candidates, minionCR)
	if minion == nil {
		minion = pickLowestCR(candidates)
	}

	minionEL := 0
	if minion != nil && minion.ChallengeRating != leader.ChallengeRating {
		qty := rand.Intn(3) + 1
		minionEL = groupEL(CRToInt(minion.ChallengeRating), qty)
		groups = append(groups, Group{Monster: *minion, Quantity: qty})
	}

	totalEL := combineEL([]int{leaderEL, minionEL})

	return &Result{
		Groups:     groups,
		Difficulty: diff,
		TargetEL:   totalEL,
		PartyLevel: level,
		PartyCount: players,
	}, nil
}

func generateFixed(candidates []monster.Monster, target, qty int, diff string, level, players int) *Result {
	var pool []monster.Monster
	for _, m := range candidates {
		cr := CRToInt(m.ChallengeRating)
		if cr >= target-2 && cr <= target+1 {
			pool = append(pool, m)
		}
	}
	if len(pool) == 0 {
		pool = candidates
	}

	selected := make([]monster.Monster, qty)
	for i := 0; i < qty; i++ {
		selected[i] = pool[rand.Intn(len(pool))]
	}

	counts := make(map[string]int)
	monsterMap := make(map[string]monster.Monster)
	for _, m := range selected {
		counts[m.Name]++
		monsterMap[m.Name] = m
	}

	var groups []Group
	for _, m := range pool {
		if q := counts[m.Name]; q > 0 {
			groups = append(groups, Group{Monster: m, Quantity: q})
			delete(counts, m.Name)
		}
	}

	var els []int
	for _, g := range groups {
		els = append(els, groupEL(CRToInt(g.Monster.ChallengeRating), g.Quantity))
	}
	totalEL := combineEL(els)

	return &Result{
		Groups:     groups,
		Difficulty: diff,
		TargetEL:   totalEL,
		PartyLevel: level,
		PartyCount: players,
	}
}

func DisplayResult(r *Result) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Encontro para %d jogadores de nível %d (%s — EL %d)",
		r.PartyCount, r.PartyLevel, r.Difficulty, r.TargetEL))
	b.WriteString("\n")
	b.WriteString(strings.Repeat("─", 50))
	b.WriteString("\n")
	for _, g := range r.Groups {
		b.WriteString(fmt.Sprintf("  %s (CR %s) × %d\n", g.Monster.Name, g.Monster.ChallengeRating, g.Quantity))
	}
	return b.String()
}
