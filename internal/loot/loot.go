package loot

import (
	"encoding/json"
	"math/rand"
	"os"
)

type EnchantmentType bool

const (
	Blessing EnchantmentType = true
	Curse    EnchantmentType = false
)

type Enchantment struct {
	Title       string          `json:"title"`
	Type        EnchantmentType `json:"type"`
	Description string          `json:"description"`
	Chance      float64         `json:"chance"`
}

func LoadJSON[T any](path string) ([]T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	items := []T{}
	err = json.Unmarshal(data, &items)
	return items, err
}

type LootTable struct {
	Weapons           []Weapon
	Armors            []Armor
	Accessories       []Accessory
	Potions           []Potion
	WeaponEnchants    []Enchantment
	ArmorEnchants     []Enchantment
	AccessoryEnchants []Enchantment
	ColateralEffects  []ColateralEffect
}

func (lt *LootTable) Generate(n int) []any {
	items := make([]any, 0, n)
	for range n {
		items = append(items, lt.randomItem())
	}
	return items
}

func (lt *LootTable) randomItem() any {
	switch rand.Intn(4) {
	case 0:
		return RandomWeapon(lt.Weapons, lt.WeaponEnchants)
	case 1:
		return RandomArmor(lt.Armors, lt.ArmorEnchants)
	case 2:
		return RandomAccessory(lt.Accessories, lt.AccessoryEnchants)
	case 3:
		return RandomPotion(lt.Potions, lt.ColateralEffects)
	}

	return nil
}
