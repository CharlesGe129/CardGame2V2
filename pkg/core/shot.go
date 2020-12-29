package core

import "github.com/CharlesGe129/CardGame2V2/pkg/def"

type Shot struct {
	cards  Cards
	team   uint8
	score  uint8
	player string
}

func NewShot(cards Cards, team uint8, player string) *Shot {
	shot := Shot{
		cards:  cards,
		team:   team,
		player: player,
	}
	for _, card := range cards.Cards {
		if card.Num == 5 || def.MapCardName[card.Num] == "5" {
			shot.score += 5
		} else if card.Num == 10 || def.MapCardName[card.Num] == "0" {
			shot.score += 10
		} else if card.Num == 13 || def.MapCardName[card.Num] == "K" {
			shot.score += 10
		}
	}
	return &shot
}

func (s *Shot) Info() (uint8, Cards) {
	return s.team, s.cards
}
