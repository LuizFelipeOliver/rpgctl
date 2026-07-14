package monster

import (
	"encoding/json"
	"os"
	"strings"
)

type Monster struct {
	ID               int    `json:"id"`
	Family           string `json:"family"`
	Name             string `json:"name"`
	AltName          string `json:"altname"`
	Size             string `json:"size"`
	Type             string `json:"type"`
	Descriptor       string `json:"descriptor"`
	HitDice          string `json:"hit_dice"`
	HD               int    `json:"hd"`
	Die              int    `json:"die"`
	HP               int    `json:"hp"`
	Initiative       string `json:"initiative"`
	Speed            string `json:"speed"`
	ArmorClass       string `json:"armor_class"`
	BaseAttack       string `json:"base_attack"`
	Grapple          string `json:"grapple"`
	Attack           string `json:"attack"`
	FullAttack       string `json:"full_attack"`
	Space            string `json:"space"`
	Reach            string `json:"reach"`
	SpecialAttacks   string `json:"special_attacks"`
	SpecialQualities string `json:"special_qualities"`
	Saves            string `json:"saves"`
	Abilities        string `json:"abilities"`
	Skills           string `json:"skills"`
	Feats            string `json:"feats"`
	EpicFeats        string `json:"epic_feats"`
	Environment      string `json:"environment"`
	Organization     string `json:"organization"`
	ChallengeRating  string `json:"challenge_rating"`
	Treasure         string `json:"treasure"`
	Alignment        string `json:"alignment"`
	Advancement      string `json:"advancement"`
	LevelAdjustment  string `json:"level_adjustment"`
	SpecialAbilities string `json:"special_abilities"`
	StatBlock        string `json:"stat_block"`
	FullText         string `json:"full_text"`
}

func LoadMonsters(path string) ([]Monster, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var monsters []Monster
	err = json.Unmarshal(data, &monsters)
	return monsters, err
}

func Search(monsters []Monster, query string) []Monster {
	query = strings.TrimSpace(query)
	if query == "" {
		return monsters
	}
	q := strings.ToLower(query)
	var result []Monster
	for _, m := range monsters {
		if strings.Contains(strings.ToLower(m.Name), q) {
			result = append(result, m)
		}
	}
	return result
}
