# rpgctl

Gerenciador de campanhas de RPG via terminal com interface TUI.

## Funcionalidades

- **🎲 Rolar dados** — expressões como `3d6+2`, `2d20`, `1d8+4`
- **⚔️ Gerar loot** — armas, armaduras, acessórios e poções com:
  - Enchants aleatórios (bênçãos e maldições com chances configuráveis)
  - Até 3 enchants por arma, 1 por armadura/acessório
  - Efeitos colaterais em poções (30% de chance)
- **📋 Iniciativa** — gerenciamento de ordem de combate *(em breve)*
- **🖥️ TUI** — interface interativa com Bubble Tea *(em breve)*

## Instalação

```bash
go install github.com/LuizFelipeOliver/rpgctl@latest
```

## Uso

```bash
# Gerar 5 itens aleatórios
rpgctl loot 5

# Rolar dados
rpgctl roll 3d6+2

# Entrar na interface TUI
rpgctl tui
```

## Estrutura do projeto

```
.
├── main.go                 # Ponto de entrada
├── cmd/                    # Comandos CLI (Cobra)
├── data/                   # Dados em JSON (editáveis)
│   ├── weapon/
│   ├── armor/
│   ├── accessory/
│   └── potion/
├── internal/
│   ├── loot/               # Geração de loot
│   ├── dice/               # Rolagem de dados
│   ├── initiative/         # Gerenciamento de iniciativa
│   └── tui/                # Interface Bubble Tea
└── LICENSE
```

## Configuração

Os dados de loot ficam em `data/` como JSON. Você pode editar diretamente para:

- Adicionar novas armas, armaduras, acessórios e poções
- Criar novos enchants (bênçãos e maldições)
- Ajustar chances de drop

## Sistema de loot

### Chances

| Categoria | Com enchant | Blessing | Curse | Max |
|-----------|:-----------:|:--------:|:-----:|:---:|
| Arma | 50% | 25% | 35% | 3 |
| Armadura | 40% | 50% | 50% | 1 |
| Acessório | 60% | 50% | 50% | 1 |
| Poção | — | — | — | Efeito colateral 30% |

Ao chamar `Generate(n)`, os itens são sorteados igualmente entre as 4 categorias.

## Tecnologias

- [Go](https://go.dev)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — framework TUI
- [Cobra](https://github.com/spf13/cobra) — CLI
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — estilização
- Dados em JSON editáveis sem recompilar

## Licença

[![CC BY-NC-SA 4.0][cc-by-nc-sa-shield]][cc-by-nc-sa]

Esta obra tem a licença [Creative Commons Atribuição-NãoComercial-CompartilhaIgual 4.0 Internacional][cc-by-nc-sa].

[![CC BY-NC-SA 4.0][cc-by-nc-sa-image]][cc-by-nc-sa]

[cc-by-nc-sa]: https://creativecommons.org/licenses/by-nc-sa/4.0/deed.pt_BR
[cc-by-nc-sa-image]: https://licensebuttons.net/l/by-nc-sa/4.0/88x31.png
[cc-by-nc-sa-shield]: https://img.shields.io/badge/License-CC%20BY--NC--SA%204.0-lightgrey.svg
