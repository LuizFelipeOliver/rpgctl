package loot

import "math/rand"

type Accessory struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Cost        int           `json:"cost"`
	Currency    string        `json:"currency"`
	Enchantment []Enchantment `json:"enchantment"`
}

const accessoryEnchantChance = 0.60

func RandomAccessory(accessories []Accessory, enchants []Enchantment) Accessory {
	a := accessories[rand.Intn(len(accessories))]

	if rand.Float64() >= accessoryEnchantChance {
		return a
	}

	e, ok := pickSingleEnchant(enchants)
	if ok {
		a.Enchantment = append(a.Enchantment, e)
	}

	return a
}
