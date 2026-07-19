package encounter

import (
	"testing"

	"rpg-tui/internal/monster"
)

func TestCRToFloat(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"1/3", 0.3333333333333333},
		{"1/2", 0.5},
		{"1", 1},
		{"10", 10},
		{"0", 0},
	}
	for _, tt := range tests {
		got, err := CRToFloat(tt.input)
		if err != nil {
			t.Errorf("CRToFloat(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("CRToFloat(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestCRToFloatError(t *testing.T) {
	_, err := CRToFloat("")
	if err == nil {
		t.Error("expected error for empty CR")
	}
	_, err = CRToFloat("abc")
	if err == nil {
		t.Error("expected error for invalid CR")
	}
}

func TestCRToInt(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"1/3", 0},
		{"1/2", 0},
		{"1", 1},
		{"3", 3},
		{"10", 10},
	}
	for _, tt := range tests {
		got := CRToInt(tt.input)
		if got != tt.want {
			t.Errorf("CRToInt(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestGroupEL(t *testing.T) {
	tests := []struct {
		cr, qty, want int
	}{
		{3, 1, 3},
		{3, 2, 5},
		{3, 3, 6},
		{3, 4, 6},
		{3, 5, 7},
		{3, 8, 8},
		{3, 12, 9},
		{0, 0, 0},
	}
	for _, tt := range tests {
		got := groupEL(tt.cr, tt.qty)
		if got != tt.want {
			t.Errorf("groupEL(%d, %d) = %d, want %d", tt.cr, tt.qty, got, tt.want)
		}
	}
}

func TestCombineEL(t *testing.T) {
	tests := []struct {
		els  []int
		want int
	}{
		{[]int{5, 5}, 7},
		{[]int{5, 3}, 6},
		{[]int{5, 1}, 5},
		{[]int{5}, 5},
		{[]int{}, 0},
	}
	for _, tt := range tests {
		got := combineEL(tt.els)
		if got != tt.want {
			t.Errorf("combineEL(%v) = %d, want %d", tt.els, got, tt.want)
		}
	}
}

func TestGenerate(t *testing.T) {
	monsters := []monster.Monster{
		{Name: "Goblin", ChallengeRating: "1/3"},
		{Name: "Orc", ChallengeRating: "1/2"},
		{Name: "Bugbear", ChallengeRating: "2"},
		{Name: "Ogro", ChallengeRating: "3"},
		{Name: "Troll", ChallengeRating: "5"},
		{Name: "Bebith", ChallengeRating: "10"},
	}

	r, err := Generate(monsters, 4, 5, "D")
	if err != nil {
		t.Fatalf("Generate() unexpected error: %v", err)
	}

	if r.PartyCount != 4 {
		t.Errorf("expected PartyCount=4, got %d", r.PartyCount)
	}
	if r.PartyLevel != 5 {
		t.Errorf("expected PartyLevel=5, got %d", r.PartyLevel)
	}
	if r.Difficulty != "D" {
		t.Errorf("expected Difficulty=D, got %s", r.Difficulty)
	}
	if r.TargetEL < 4 || r.TargetEL > 10 {
		t.Errorf("TargetEL out of range [4,10]: got %d", r.TargetEL)
	}
	if len(r.Groups) == 0 {
		t.Fatal("expected at least 1 group")
	}
}

func TestGenerateInvalidDifficulty(t *testing.T) {
	_, err := Generate(nil, 4, 5, "X")
	if err == nil {
		t.Error("expected error for invalid difficulty")
	}
}

func TestGenerateEmpty(t *testing.T) {
	_, err := Generate(nil, 4, 1, "M")
	if err == nil {
		t.Error("expected error for empty monster list")
	}
}

func TestDisplayResult(t *testing.T) {
	r := &Result{
		Groups: []Group{
			{Monster: monster.Monster{Name: "Ogro", ChallengeRating: "3"}, Quantity: 1},
			{Monster: monster.Monster{Name: "Goblin", ChallengeRating: "1/3"}, Quantity: 3},
		},
		Difficulty: "D",
		TargetEL:   6,
		PartyLevel: 5,
		PartyCount: 4,
	}
	s := DisplayResult(r)
	if s == "" {
		t.Error("expected non-empty display string")
	}
}
