package loot

import (
	"math/rand"
	"testing"
)

func weaponPool() []Weapon {
	return []Weapon{
		{Name: "Espada", Damage: "1d8", DamageType: Slashing, Cost: 15, Currency: "gp", Category: Martial, Type: OneHanded},
		{Name: "Machado", Damage: "1d12", DamageType: Slashing, Cost: 30, Currency: "gp", Category: Martial, Type: TwoHanded},
		{Name: "Arco", Damage: "1d6", DamageType: Piercing, Cost: 25, Currency: "gp", Category: Simple, Type: Ranged},
	}
}

func armorPool() []Armor {
	return []Armor{
		{Name: "Couro", ArmorBonus: 11, Cost: 10, Currency: "gp", Type: LightArmor},
		{Name: "Malha", ArmorBonus: 16, Cost: 75, Currency: "gp", Type: HeavyArmor},
	}
}

func accessoryPool() []Accessory {
	return []Accessory{
		{Name: "Anel", Type: "ring", Cost: 5, Currency: "gp"},
		{Name: "Colar", Type: "necklace", Cost: 10, Currency: "gp"},
	}
}

func potionPool() []Potion {
	return []Potion{
		{Name: "Cura", Cost: 50, Currency: "gp", Effect: "2d4+2 PV"},
		{Name: "Força", Cost: 100, Currency: "gp", Effect: "Força +4"},
	}
}

func enchantPool() []Enchantment {
	return []Enchantment{
		{Title: "Fogo", Type: Blessing, Description: "dano de fogo", Chance: 0.20},
		{Title: "Gelo", Type: Blessing, Description: "dano de gelo", Chance: 0.15},
		{Title: "Sangramento", Type: Curse, Description: "sangramento", Chance: 0.25},
		{Title: "Ferrugem", Type: Curse, Description: "quebra facil", Chance: 0.15},
	}
}

func colateralPool() []ColateralEffect {
	return []ColateralEffect{
		{Description: "Enjoo", Chance: 0.30},
		{Description: "Tontura", Chance: 0.25},
	}
}

func TestFilterByKind(t *testing.T) {
	enchants := enchantPool()

	blessings := filterByKind(enchants, Blessing)
	if len(blessings) != 2 {
		t.Errorf("expected 2 blessings, got %d", len(blessings))
	}
	for _, e := range blessings {
		if e.Type != Blessing {
			t.Errorf("expected blessing, got %v", e.Type)
		}
	}

	curses := filterByKind(enchants, Curse)
	if len(curses) != 2 {
		t.Errorf("expected 2 curses, got %d", len(curses))
	}
	for _, e := range curses {
		if e.Type != Curse {
			t.Errorf("expected curse, got %v", e.Type)
		}
	}
}

func TestFilterByKindEmpty(t *testing.T) {
	result := filterByKind(nil, Blessing)
	if len(result) != 0 {
		t.Errorf("expected empty, got %d", len(result))
	}
}

func TestPickWeighted(t *testing.T) {
	enchants := []Enchantment{
		{Title: "A", Chance: 1.0},
		{Title: "B", Chance: 0.0},
	}

	result := pickWeighted(enchants)
	if result.Title != "A" {
		t.Errorf("expected A, got %s", result.Title)
	}
}

func TestPickWeightedSingle(t *testing.T) {
	enchants := []Enchantment{{Title: "Unico", Chance: 0.50}}
	result := pickWeighted(enchants)
	if result.Title != "Unico" {
		t.Errorf("expected Unico, got %s", result.Title)
	}
}

func TestRandomWeaponReturnsFromPool(t *testing.T) {
	weapons := weaponPool()
	enchants := enchantPool()

	for range 100 {
		w := RandomWeapon(weapons, enchants)
		if w.Name == "" {
			t.Fatal("weapon should not be empty")
		}
		found := false
		for _, wp := range weapons {
			if wp.Name == w.Name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("weapon %s not in pool", w.Name)
		}
	}
}

func TestRandomWeaponEnchantmentsFromPool(t *testing.T) {
	weapons := weaponPool()
	enchants := enchantPool()

	for range 200 {
		w := RandomWeapon(weapons, enchants)
		for _, e := range w.Enchantment {
			found := false
			for _, ep := range enchants {
				if ep.Title == e.Title {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("enchant %s not in pool", e.Title)
			}
		}
	}
}

func TestRandomWeaponMaxEnchantments(t *testing.T) {
	weapons := weaponPool()
	enchants := enchantPool()

	maxSeen := 0
	for range 500 {
		w := RandomWeapon(weapons, enchants)
		if len(w.Enchantment) > maxSeen {
			maxSeen = len(w.Enchantment)
		}
	}
	if maxSeen > maxEnchants {
		t.Errorf("max enchants exceeded: got %d, max %d", maxSeen, maxEnchants)
	}
}

func TestRandomArmorReturnsFromPool(t *testing.T) {
	armors := armorPool()
	enchants := enchantPool()

	for range 100 {
		a := RandomArmor(armors, enchants)
		if a.Name == "" {
			t.Fatal("armor should not be empty")
		}
		found := false
		for _, ap := range armors {
			if ap.Name == a.Name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("armor %s not in pool", a.Name)
		}
	}
}

func TestRandomArmorMaxOneEnchant(t *testing.T) {
	armors := armorPool()
	enchants := enchantPool()

	for range 200 {
		a := RandomArmor(armors, enchants)
		if len(a.Enchantment) > 1 {
			t.Errorf("armor should have at most 1 enchant, got %d", len(a.Enchantment))
		}
	}
}

func TestRandomAccessoryReturnsFromPool(t *testing.T) {
	accessories := accessoryPool()
	enchants := enchantPool()

	for range 100 {
		a := RandomAccessory(accessories, enchants)
		if a.Name == "" {
			t.Fatal("accessory should not be empty")
		}
		found := false
		for _, ap := range accessories {
			if ap.Name == a.Name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("accessory %s not in pool", a.Name)
		}
	}
}

func TestRandomAccessoryMaxOneEnchant(t *testing.T) {
	accessories := accessoryPool()
	enchants := enchantPool()

	for range 200 {
		a := RandomAccessory(accessories, enchants)
		if len(a.Enchantment) > 1 {
			t.Errorf("accessory should have at most 1 enchant, got %d", len(a.Enchantment))
		}
	}
}

func TestRandomPotionReturnsFromPool(t *testing.T) {
	potions := potionPool()
	effects := colateralPool()

	for range 100 {
		r := RandomPotion(potions, effects)
		if r.Potion.Name == "" {
			t.Fatal("potion should not be empty")
		}
		found := false
		for _, pp := range potions {
			if pp.Name == r.Potion.Name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("potion %s not in pool", r.Potion.Name)
		}
	}
}

func TestRandomPotionColateralChance(t *testing.T) {
	potions := potionPool()
	effects := colateralPool()

	withColateral := 0
	total := 10000
	for range total {
		r := RandomPotion(potions, effects)
		if r.ColateralEffect != nil {
			withColateral++
		}
	}

	pct := float64(withColateral) / float64(total)
	if pct < 0.20 || pct > 0.40 {
		t.Errorf("colateral chance outside expected range: got %.2f, expected ~0.30", pct)
	}
}

func TestGenerateReturnsCorrectCount(t *testing.T) {
	rand.Seed(42)

	lt := &LootTable{
		Weapons:          weaponPool(),
		Armors:           armorPool(),
		Accessories:      accessoryPool(),
		Potions:          potionPool(),
		WeaponEnchants:   enchantPool(),
		ArmorEnchants:    enchantPool(),
		AccessoryEnchants: enchantPool(),
		ColateralEffects: colateralPool(),
	}

	items := lt.Generate(10)
	if len(items) != 10 {
		t.Errorf("expected 10 items, got %d", len(items))
	}
}

func TestGenerateReturnsValidTypes(t *testing.T) {
	lt := &LootTable{
		Weapons:          weaponPool(),
		Armors:           armorPool(),
		Accessories:      accessoryPool(),
		Potions:          potionPool(),
		WeaponEnchants:   enchantPool(),
		ArmorEnchants:    enchantPool(),
		AccessoryEnchants: enchantPool(),
		ColateralEffects: colateralPool(),
	}

	items := lt.Generate(100)
	for _, item := range items {
		switch item.(type) {
		case Weapon, Armor, Accessory, PotionResult:
		default:
			t.Errorf("unexpected item type: %T", item)
		}
	}
}

func TestPickSingleEnchant(t *testing.T) {
	enchants := []Enchantment{
		{Title: "Bless", Type: Blessing, Chance: 1.0},
		{Title: "Curse", Type: Curse, Chance: 1.0},
	}

	gotBless := false
	gotCurse := false
	for range 100 {
		e, ok := pickSingleEnchant(enchants)
		if !ok {
			continue
		}
		if e.Title == "Bless" {
			gotBless = true
		}
		if e.Title == "Curse" {
			gotCurse = true
		}
		if gotBless && gotCurse {
			break
		}
	}
	if !gotBless {
		t.Error("never got a blessing")
	}
	if !gotCurse {
		t.Error("never got a curse")
	}
}

func TestPickSingleEnchantEmptyPool(t *testing.T) {
	_, ok := pickSingleEnchant(nil)
	if ok {
		t.Error("expected false for empty pool")
	}
}

func TestLoadJSONWithInvalidPath(t *testing.T) {
	_, err := LoadJSON[Weapon]("nonexistent.json")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestPickWeightedColateral(t *testing.T) {
	effects := []ColateralEffect{
		{Description: "A", Chance: 1.0},
		{Description: "B", Chance: 0.0},
	}

	result := pickWeightedColateral(effects)
	if result.Description != "A" {
		t.Errorf("expected A, got %s", result.Description)
	}
}


