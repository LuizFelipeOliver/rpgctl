package main

import "testing"

func TestDiceCurrentBehavior(t *testing.T) {
	t.Run("valid expressions", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
			min   int
			max   int
		}{
			{name: "single die without quantity", input: "d20", min: 1, max: 20},
			{name: "single die with quantity", input: "1d20", min: 1, max: 20},
			{name: "multiple dice", input: "2d6", min: 2, max: 12},
			{name: "die with positive modifier", input: "1d4+3", min: 4, max: 7},
			{name: "die with negative and positive modifiers", input: "1d20-2+3", min: 2, max: 21},
			{name: "die without quantity and mixed modifiers", input: "d12+3-1", min: 3, max: 14},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got, err := dice(test.input)

				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}

				if got < test.min || got > test.max {
					t.Fatalf("expected result between %d and %d, got %d", test.min, test.max, got)
				}
			})
		}
	})

	t.Run("invalid expressions", func(t *testing.T) {
		tests := []struct {
			name  string
			input string
		}{
			{name: "invalid quantity zero", input: "0d6"},
			{name: "invalid negative quantity", input: "-1d6"},
			{name: "missing sides", input: "2d"},
			{name: "too many separators", input: "2dd6"},
			{name: "unsupported die", input: "1d7"},
			{name: "missing d", input: "20"},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				got, err := dice(test.input)

				if err == nil {
					t.Fatalf("expected error, got nil and result %d", got)
				}
			})
		}
	})
}
