package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type Player struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func RunInitiative(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("use: rpgctl init <add|list|remove|clear>")
	}

	switch args[0] {
	case "add":
		if len(args) < 3 {
			return fmt.Errorf("use: rpgctl init add <name> <value>")
		}

		name := args[1]

		value, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("iniciativa invalida: %s", args[2])
		}

		entries, err := loadInitiative()
		if err != nil {
			return err
		}

		entries = addInitiative(entries, name, value)

		if err := prepareInitiative(entries); err != nil {
			return err
		}

		if err := saveInitiative(entries); err != nil {
			return err
		}

	case "remove":
		if len(args) < 2 {
			return fmt.Errorf("use: rpgctl init remove <name>")
		}

		name := args[1]

		entries, err := loadInitiative()
		if err != nil {
			return err
		}

		entries = removeInitiative(entries, name)

		if err := saveInitiative(entries); err != nil {
			return err
		}

	case "list":
		entries, err := loadInitiative()
		if err != nil {
			return err
		}
		for index, entry := range entries {
			fmt.Printf("[%d] - %s\n", index+1, entry.Name)
		}
	default:
		return fmt.Errorf("comando invalido: %s", args[0])
	}
	return nil
}

func addInitiative(entries []Player, name string, value int) []Player {
	return append(entries, Player{Name: name, Value: value})
}

func removeInitiative(entries []Player, name string) []Player {
	var result []Player

	for _, entry := range entries {
		if entry.Name != name {
			result = append(result, entry)
		}
	}

	return result
}

func prepareInitiative(entries []Player) error {
	seen := make(map[int]string)

	for _, entry := range entries {
		if name, ok := seen[entry.Value]; ok {
			return fmt.Errorf(
				"empate na iniciativa %d entre %s e %s; rolem para desempatar",
				entry.Value,
				name,
				entry.Name,
			)
		}

		seen[entry.Value] = entry.Name
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Value > entries[j].Value
	})

	return nil
}

func loadInitiative() ([]Player, error) {
	path, err := initiativePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []Player{}, nil
		}
		return nil, err
	}

	var entries []Player
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func initiativePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "rpgctl", "initiative.json"), nil
}

func saveInitiative(entries []Player) error {
	path, err := initiativePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
