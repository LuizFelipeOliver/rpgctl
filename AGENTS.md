# rpgctl

## Build & verify
- `go build ./... && go vet ./... && go test ./...`

## CLI
- `rpgctl loot [n]` — generates n items (default 1)
- `rpgctl roll "d20+1d4"` — dice expression
- `rpgctl tui` — Bubble Tea interface

## Architecture
- `main.go` → `cmd.Execute()` (Cobra)
- `cmd/` — 4 commands; each `init()` registers to `rootCmd`
- `internal/loot/` — `Generate(n) []any`, caller type-switches
- `internal/dice/` — `Roll(expr) (*RollResult, error)`, parser with save/restore
- `internal/tui/` — Bubble Tea; tabs[] hardcoded in tabs.go
- `data/` — editable JSON, one dir per category

## Code patterns
- **Guard clauses + early returns** throughout (validate → bail → proceed)
- **Speculative parsing** in dice/parser.go: `readDice()` saves position, restores on failure
- roll command uses `DisableFlagParsing: true` (dice `-1d4` would be parsed as flags)

## Conventions
- All labels, types, errors in **Portuguese**
- JSON: snake_case fields; enchantment `type: true` = blessing, `false` = curse
- Dice constants (D4..D100) are convenience — any side count works (d1, d3, d24, d1000)
- Expression MUST start with a dice group (`d20`, `-d24`, `2d6`). Stray modifier-only (`3`) or number-first (`3+d20`) is rejected.
- Negative dice groups allowed (`d20-1d4`, `-d24+d4`). Use `g.Sign` to track (±1).
- All data types, labels, and identifiers in Portuguese throughout the codebase.

## Gotchas
- loot `Generate()` returns `[]any` — must switch on concrete type
- `DisplayItems()` returns ANSI-colored string with `╭─╮` box-drawing borders
- TUI footer and count buffer ([N]g) are tab-specific (tui.go switch)
- roll command: `DisableFlagParsing: true` — manual `--help` detection in RunE
