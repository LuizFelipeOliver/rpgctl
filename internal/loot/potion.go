package loot

import "math/rand"

type Potion struct {
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Currency string `json:"currency"`
	Effect   string `json:"effect"`
}

type ColateralEffect struct {
	Description string  `json:"description"`
	Chance      float64 `json:"chance"`
}

type PotionResult struct {
	Potion          Potion
	ColateralEffect *ColateralEffect
}

const collateralChance = 0.30

func RandomPotion(potions []Potion, effects []ColateralEffect) PotionResult {
	p := potions[rand.Intn(len(potions))]

	if rand.Float64() >= collateralChance {
		return PotionResult{Potion: p}
	}

	e := pickWeightedColateral(effects)
	return PotionResult{Potion: p, ColateralEffect: &e}
}

func pickWeightedColateral(effects []ColateralEffect) ColateralEffect {
	total := 0.0
	for _, e := range effects {
		total += e.Chance
	}

	roll := rand.Float64() * total
	cum := 0.0
	for _, e := range effects {
		cum += e.Chance
		if roll < cum {
			return e
		}
	}

	return effects[0]
}
