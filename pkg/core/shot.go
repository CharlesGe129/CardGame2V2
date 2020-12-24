package core

import (
	"fmt"

	"github.com/CharlesGe129/CardGame2V2/pkg/def"
)

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
		if card.Num == 5 {
			shot.score += 5
		} else if card.Num == 10 {
			shot.score += 10
		} else if card.Num == 13 {
			shot.score += 10
		}
	}
	return &shot
}

type Round struct {
	shotList []Shot

	biggest Shot
	num     int
}

func NewRound(mainColor def.CardColor, shot Shot) *Round {
	return &Round{
		shotList: []Shot{shot},
		biggest:  shot,
		num:      len(shot.cards.Cards),
	}
}

func (r *Round) AddShot(shot Shot) {
	r.shotList = append(r.shotList, shot)
	if !r.biggest.cards.LargerOrEqualTo(&shot.cards) {
		r.biggest = shot
	}
	return
}

func (r *Round) GetResult() (team, score uint8, player string) {
	team = r.biggest.team
	for _, shot := range r.shotList {
		if shot.team == team {
			score += shot.score
		}
	}
	return team, score, r.biggest.player
}

func (r *Round) ShowShots() {
	for idx, shot := range r.shotList {
		if shot.player == r.biggest.player {
			fmt.Printf("[最大]")
		}
		fmt.Printf("第%d位玩家出的牌: %s\n", idx+1, shot.cards)
	}
	fmt.Printf("您是第%d位\n", len(r.shotList))
}
