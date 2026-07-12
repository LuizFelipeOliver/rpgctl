package loot

import "math/rand"

type Weapon struct {
	Name        string        `json:"name"`
	Damage      string        `json:"damage"`
	DamageType  DamageTypes   `json:"damage_type"`
	Cost        int           `json:"cost"`
	Currency    string        `json:"currency"`
	Category    Category      `json:"category"`
	Type        Type          `json:"type"`
	Enchantment []Enchantment `json:"enchantment"`
}

type DamageTypes string

const (
	Bludgeoning DamageTypes = "contusão"
	Piercing    DamageTypes = "perfurante"
	Slashing    DamageTypes = "cortante"
)

type Category string

const (
	Exotic  Category = "exótica"
	Martial Category = "marcial"
	Simple  Category = "simples"
)

type Type string

const (
	OneHanded Type = "uma-mão"
	TwoHanded Type = "duas-mãos"
	Light     Type = "leve"
	Ranged    Type = "distância"
	Unarmed   Type = "desarmado"
)

const (
	enchantChance  = 0.50
	blessingChance = 0.25
	curseChance    = 0.35
	maxEnchants    = 3
)

func RandomWeapon(weapons []Weapon, enchants []Enchantment) Weapon {
	weapon := weapons[rand.Intn(len(weapons))]

	if rand.Float64() >= enchantChance {
		return weapon
	}

	n := rand.Intn(maxEnchants) + 1
	for range n {
		e, ok := pickEnchantment(enchants)
		if !ok {
			continue
		}
		weapon.Enchantment = append(weapon.Enchantment, e)
	}

	return weapon
}

func pickEnchantment(enchants []Enchantment) (Enchantment, bool) {
	kind, ok := rollEnchantType()
	if !ok {
		return Enchantment{}, false
	}

	pool := filterByKind(enchants, kind)
	if len(pool) == 0 {
		return Enchantment{}, false
	}

	return pickWeighted(pool), true
}

func rollEnchantType() (EnchantmentType, bool) {
	roll := rand.Float64()
	switch {
	case roll < blessingChance:
		return Blessing, true
	case roll < blessingChance+curseChance:
		return Curse, true
	default:
		return Blessing, false
	}
}

func filterByKind(enchants []Enchantment, kind EnchantmentType) []Enchantment {
	var result []Enchantment
	for _, e := range enchants {
		if e.Type == kind {
			result = append(result, e)
		}
	}
	return result
}

func pickSingleEnchant(enchants []Enchantment) (Enchantment, bool) {
	kind := Blessing
	if rand.Float64() < 0.5 {
		kind = Curse
	}

	pool := filterByKind(enchants, kind)
	if len(pool) == 0 {
		return Enchantment{}, false
	}

	return pickWeighted(pool), true
}

func pickWeighted(pool []Enchantment) Enchantment {
	total := 0.0
	for _, e := range pool {
		total += e.Chance
	}

	roll := rand.Float64() * total
	cum := 0.0
	for _, e := range pool {
		cum += e.Chance
		if roll < cum {
			return e
		}
	}

	return pool[0]
}
