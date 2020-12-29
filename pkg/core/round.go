package core

import (
	"fmt"

	"github.com/CharlesGe129/CardGame2V2/pkg/def"
)

type Round struct {
	shotList []Shot

	biggest Shot
	num     int
	Color   def.CardColor
}

func NewRound(shot Shot) *Round {
	fmt.Printf("%q出牌：%s\n", shot.player, shot.cards.String())
	round := Round{
		shotList: []Shot{shot},
		biggest:  shot,
		num:      len(shot.cards.Cards),
		Color:    shot.cards.Cards[0].Color,
	}
	if round.Color == def.CardColorNil || shot.cards.Cards[0].IsMain {
		round.Color = shot.cards.mainColor
	}
	return &round
}

func (r *Round) AddShot(shot Shot) {
	fmt.Printf("%q出牌：%s\n%+v\n", shot.player, shot.cards.String(), shot.cards)
	r.shotList = append(r.shotList, shot)
	if !r.biggest.cards.LargerOrEqualTo(&shot.cards) {
		r.biggest = shot
	}
	return
}

func (r *Round) GetResult() (team, score uint8, player string) {
	team = r.biggest.team
	for _, shot := range r.shotList {
		score += shot.score
	}
	fmt.Printf("%q最大\n\n", r.biggest.player)
	return team, score, r.biggest.player
}

func (r *Round) ShowShots() {
	for idx, shot := range r.shotList {
		if shot.player == r.biggest.player {
			fmt.Printf("[最大]")
		}
		fmt.Printf("第%d位玩家(%q)出的牌: %s\n", idx+1, shot.player, shot.cards)
	}
	fmt.Printf("您是第%d位\n", len(r.shotList)+1)
}

func (r *Round) GetBiggest() Shot {
	return r.biggest
}
