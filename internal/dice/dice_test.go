package dice

import (
	"testing"
)

func TestRoll(t *testing.T) {
	tests := []struct {
		expr      string
		err       bool
		numGroups int
		modifier  int
	}{
		{"d20", false, 1, 0},
		{"2d6", false, 1, 0},
		{"d20+1d4", false, 2, 0},
		{"2d6+3", false, 1, 3},
		{"d20+1d4+5", false, 2, 5},
		{"3d8", false, 1, 0},
		{"d6-2", false, 1, -2},
		{"d20-3", false, 1, -3},
		{"d20+1d4-2", false, 2, -2},

		{"d20-1d4", true, 0, 0},
		{"-1d4", true, 0, 0},
		{"3+d20", true, 0, 0},
		{"3", true, 0, 0},
		{"", true, 0, 0},
	}

	for _, tt := range tests {
		result, err := Roll(tt.expr)
		if tt.err {
			if err == nil {
				t.Errorf("Roll(%q): expected error", tt.expr)
			}
			continue
		}
		if err != nil {
			t.Errorf("Roll(%q): unexpected error: %v", tt.expr, err)
			continue
		}
		if len(result.Groups) != tt.numGroups {
			t.Errorf("Roll(%q): got %d groups, want %d", tt.expr, len(result.Groups), tt.numGroups)
		}
		if result.Modifier != tt.modifier {
			t.Errorf("Roll(%q): got modifier %d, want %d", tt.expr, result.Modifier, tt.modifier)
		}

		for _, g := range result.Groups {
			if len(g.Rolls) != g.Count {
				t.Errorf("Roll(%q): group got %d rolls, want %d", tt.expr, len(g.Rolls), g.Count)
			}
			for _, r := range g.Rolls {
				if r < 1 || r > g.Sides {
					t.Errorf("Roll(%q): roll %d out of range [1,%d]", tt.expr, r, g.Sides)
				}
			}
		}

		expectedTotal := result.Modifier
		for _, g := range result.Groups {
			for _, r := range g.Rolls {
				expectedTotal += r
			}
		}
		if result.Total != expectedTotal {
			t.Errorf("Roll(%q): total %d != expected %d", tt.expr, result.Total, expectedTotal)
		}
	}
}

func TestRoll_Range(t *testing.T) {
	const iterations = 100
	for range iterations {
		result, err := Roll("2d6+3")
		if err != nil {
			t.Fatal(err)
		}
		if result.Total < 5 || result.Total > 15 {
			t.Errorf("2d6+3 total %d out of range [5,15]", result.Total)
		}
	}
}

func TestDiceConstants(t *testing.T) {
	if D4 != 4 || D6 != 6 || D8 != 8 || D10 != 10 || D12 != 12 || D20 != 20 || D100 != 100 {
		t.Error("dice constants mismatch")
	}
}

func TestStringFormat(t *testing.T) {
	result, err := Roll("d20")
	if err != nil {
		t.Fatal(err)
	}
	s := result.String()
	if result.Total < 1 || result.Total > 20 {
		t.Errorf("d20 total %d out of range", result.Total)
	}
	_ = s
}
