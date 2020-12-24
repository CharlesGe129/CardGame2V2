package player

import (
	"github.com/CharlesGe129/CardGame2V2/pkg/core"
	"github.com/CharlesGe129/CardGame2V2/pkg/def"
)

type Player interface {
	NextShot(r *core.Round) *core.Shot
	NewShot() *core.Shot
	SetMainColor(color def.CardColor)
	BidMainColor() def.CardColor
	ShowCards()
	AddCard(card core.Card)
	RemoveCards(cardList []core.Card) error
	SetCoveredCards(cardList []core.Card) []core.Card
	IsFinished() bool
	GetName() string
	GetTeam() uint8
}
