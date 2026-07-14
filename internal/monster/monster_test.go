package monster

import (
	"strings"
	"testing"
)

func TestLoadMonsters(t *testing.T) {
	monsters, err := LoadMonsters("../../data/monster/monsters-parsed.json")
	if err != nil {
		t.Fatal(err)
	}
	if len(monsters) == 0 {
		t.Fatal("no monsters loaded")
	}
	if monsters[0].Name == "" {
		t.Error("first monster has empty name")
	}
}

func TestLoadMonsters_InvalidPath(t *testing.T) {
	_, err := LoadMonsters("nonexistent.json")
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestSearch_FindsMatch(t *testing.T) {
	monsters, err := LoadMonsters("../../data/monster/monsters-parsed.json")
	if err != nil {
		t.Fatal(err)
	}
	results := Search(monsters, "goblin")
	if len(results) == 0 {
		t.Fatal("expected at least one goblin match")
	}
	found := false
	for _, m := range results {
		if containsIgnoreCase(m.Name, "goblin") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Search results don't contain 'goblin' in name")
	}
}

func TestSearch_EmptyQueryReturnsAll(t *testing.T) {
	monsters, err := LoadMonsters("../../data/monster/monsters-parsed.json")
	if err != nil {
		t.Fatal(err)
	}
	results := Search(monsters, "")
	if len(results) != len(monsters) {
		t.Errorf("empty query: got %d, want %d", len(results), len(monsters))
	}
}

func TestSearch_NoMatch(t *testing.T) {
	monsters, err := LoadMonsters("../../data/monster/monsters-parsed.json")
	if err != nil {
		t.Fatal(err)
	}
	results := Search(monsters, "zzzznonexistentzzzz")
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	monsters, err := LoadMonsters("../../data/monster/monsters-parsed.json")
	if err != nil {
		t.Fatal(err)
	}
	r1 := Search(monsters, "Dragon")
	r2 := Search(monsters, "dragon")
	if len(r1) != len(r2) {
		t.Errorf("case insensitive mismatch: %d vs %d", len(r1), len(r2))
	}
}

func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
