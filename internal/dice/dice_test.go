package dice

import (
	"strconv"
	"strings"
	"testing"
)

func TestRollDie(t *testing.T) {
	v, err := RollDie(20)
	if err != nil {
		t.Fatal(err)
	}
	if v < 1 || v > 20 {
		t.Errorf("RollDie(20) = %d, fora do intervalo [1,20]", v)
	}
}

func TestRollDie_Invalid(t *testing.T) {
	_, err := RollDie(0)
	if err == nil {
		t.Error("RollDie(0): expected error")
	}
	_, err = RollDie(-1)
	if err == nil {
		t.Error("RollDie(-1): expected error")
	}
}

func TestRollNotation(t *testing.T) {
	tests := []struct {
		expr    string
		err     bool
		min, max int
	}{
		{"d20", false, 1, 20},
		{"2d6", false, 2, 12},
		{"1d4", false, 1, 4},
		{"d1", false, 1, 1},
		{"", true, 0, 0},
		{"abc", true, 0, 0},
		{"0d6", true, 0, 0},
		{"d0", true, 0, 0},
		{"-d6", true, 0, 0},
	}

	for _, tt := range tests {
		result, err := RollNotation(tt.expr)
		if tt.err {
			if err == nil {
				t.Errorf("RollNotation(%q): expected error, got %d", tt.expr, result)
			}
			continue
		}
		if err != nil {
			t.Errorf("RollNotation(%q): unexpected error: %v", tt.expr, err)
			continue
		}
		if result < tt.min || result > tt.max {
			t.Errorf("RollNotation(%q) = %d, fora do intervalo [%d,%d]", tt.expr, result, tt.min, tt.max)
		}
	}
}

func TestRollNotation_Range(t *testing.T) {
	const iterations = 100
	for range iterations {
		result, err := RollNotation("2d6")
		if err != nil {
			t.Fatal(err)
		}
		if result < 2 || result > 12 {
			t.Errorf("2d6 = %d, fora do intervalo [2,12]", result)
		}
	}
}

func TestD1(t *testing.T) {
	for range 10 {
		result, err := RollNotation("d1")
		if err != nil {
			t.Fatal(err)
		}
		if result != 1 {
			t.Errorf("d1 = %d, want 1", result)
		}
	}
}

func TestRolar_DiceOnly(t *testing.T) {
	for range 10 {
		r, err := Rolar("d20")
		if err != nil {
			t.Fatal(err)
		}
		if r.Total < 1 || r.Total > 20 {
			t.Errorf("d20 total=%d, fora [1,20]", r.Total)
		}
		if !strings.Contains(r.Detalhes, strconv.Itoa(r.Total)) {
			t.Errorf("detalhes=%q nao contem total=%d", r.Detalhes, r.Total)
		}
	}
}

func TestRolar_DicePlusMod(t *testing.T) {
	for range 10 {
		r, err := Rolar("d20 + 5")
		if err != nil {
			t.Fatal(err)
		}
		if r.Total < 6 || r.Total > 25 {
			t.Errorf("d20+5 total=%d, fora [6,25]", r.Total)
		}
		if !strings.Contains(r.Detalhes, " + 5") {
			t.Errorf("detalhes=%q deve conter \" + 5\"", r.Detalhes)
		}
	}
}

func TestRolar_DiceMinusMod(t *testing.T) {
	for range 10 {
		r, err := Rolar("d20 - 3")
		if err != nil {
			t.Fatal(err)
		}
		if r.Total < -2 || r.Total > 17 {
			t.Errorf("d20-3 total=%d, fora [-2,17]", r.Total)
		}
		if !strings.Contains(r.Detalhes, " - 3") {
			t.Errorf("detalhes=%q deve conter \" - 3\"", r.Detalhes)
		}
	}
}

func TestRolar_MultipleTerms(t *testing.T) {
	for range 10 {
		r, err := Rolar("d20 + 1 + d3")
		if err != nil {
			t.Fatal(err)
		}
		parts := strings.Split(r.Detalhes, " + ")
		if len(parts) != 3 {
			t.Errorf("detalhes=%q deve ter 3 partes separadas por \" + \"", r.Detalhes)
		}
	}
}

func TestRolar_PureNumbers(t *testing.T) {
	r, err := Rolar("5 + 3")
	if err != nil {
		t.Fatal(err)
	}
	if r.Total != 8 {
		t.Errorf("5+3 total=%d, want 8", r.Total)
	}
	if r.Detalhes != "5 + 3" {
		t.Errorf("detalhes=%q, want \"5 + 3\"", r.Detalhes)
	}
}

func TestRolar_Subtrai(t *testing.T) {
	r, err := Rolar("10 - 3")
	if err != nil {
		t.Fatal(err)
	}
	if r.Total != 7 {
		t.Errorf("10-3 total=%d, want 7", r.Total)
	}
	if r.Detalhes != "10 - 3" {
		t.Errorf("detalhes=%q, want \"10 - 3\"", r.Detalhes)
	}
}

func TestRolar_Error(t *testing.T) {
	_, err := Rolar("")
	if err == nil {
		t.Error("Rolar(\"\"): expected error")
	}
	_, err = Rolar("d20 + abc")
	if err == nil {
		t.Error("Rolar(\"d20 + abc\"): expected error")
	}
}
