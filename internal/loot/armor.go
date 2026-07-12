package loot

import "math/rand"

type Armor struct {
	Name        string        `json:"name"`
	ArmorBonus  int           `json:"armor_bonus"`
	Cost        int           `json:"cost"`
	Currency    string        `json:"currency"`
	Type        Types         `json:"type"`
	Enchantment []Enchantment `json:"enchantment"`
}

type Types string

const (
	LightArmor  Types = "Armadura Leve"
	MediumArmor Types = "Armadura Média"
	HeavyArmor  Types = "Armadura Pesada"
	Shield      Types = "Escudo"
	Wooden      Types = "Madeira"
	Extras      Types = "Extras"
)

const armorEnchantChance = 0.40

func RandomArmor(armors []Armor, enchants []Enchantment) Armor {
	a := armors[rand.Intn(len(armors))]

	if rand.Float64() >= armorEnchantChance {
		return a
	}

	e, ok := pickSingleEnchant(enchants)
	if ok {
		a.Enchantment = append(a.Enchantment, e)
	}

	return a
}
