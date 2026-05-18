package main

import "testing"

func TestAddInitiative(t *testing.T) {
	entries := []Player{
		{Name: "Aragorn", Value: 18},
	}

	got := addInitiative(entries, "Orc", 12)

	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}

	if got[1].Name != "Orc" || got[1].Value != 12 {
		t.Fatalf("expected Orc with value 12, got %+v", got[1])
	}
}

func TestRemoveInitiative(t *testing.T) {
	entries := []Player{
		{Name: "Aragorn", Value: 18},
		{Name: "Orc", Value: 12},
		{Name: "Goblin", Value: 8},
	}

	got := removeInitiative(entries, "Orc")

	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}

	for _, entry := range got {
		if entry.Name == "Orc" {
			t.Fatalf("expected Orc to be removed, got %+v", got)
		}
	}
}

func TestPrepareInitiativeSortsDescending(t *testing.T) {
	entries := []Player{
		{Name: "Goblin", Value: 8},
		{Name: "Aragorn", Value: 18},
		{Name: "Orc", Value: 12},
	}

	if err := prepareInitiative(entries); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	want := []string{"Aragorn", "Orc", "Goblin"}
	for i, name := range want {
		if entries[i].Name != name {
			t.Fatalf("expected entry %d to be %s, got %s", i, name, entries[i].Name)
		}
	}
}

func TestPrepareInitiativeRejectsTies(t *testing.T) {
	entries := []Player{
		{Name: "Aragorn", Value: 18},
		{Name: "Legolas", Value: 18},
	}

	if err := prepareInitiative(entries); err == nil {
		t.Fatal("expected tie error, got nil")
	}
}
