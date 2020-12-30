package player

import (
	"github.com/CharlesGe129/CardGame2V2/pkg/core"
	"github.com/CharlesGe129/CardGame2V2/pkg/def"
)

type Player interface {
	NextShot(r *core.Round) (*core.Shot, error)
	NewShot() *core.Shot
	SetMainColor(pool core.CardPool)
	BidMainColor() def.CardColor
	ShowCards()
	AddCard(card core.Card)
	RemoveCards(cardList []core.Card) ([]core.Card, error)
	SetCoveredCards(origCoveredCards []core.Card) ([]core.Card, error)
	IsFinished() bool
	GetName() string
	GetTeam() uint8
}