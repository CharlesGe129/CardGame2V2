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

type player struct {
	Name string
	Team uint8

	mainColor    def.CardColor
	pool         core.CardPool
	cardsByColor map[def.CardColor][]core.Card
}

func (p *player) SetMainColor(pool core.CardPool) {
	p.mainColor = pool.MainColor
	p.pool = pool

	// arrange main cards
	var mainCardList []core.Card
	for color, curCardList := range p.cardsByColor {
		switch color {
		case def.CardColorNil:
			for _, card := range curCardList {
				if card.Color == def.CardColorNil {
					card.Color = p.mainColor
				}
				mainCardList = append(mainCardList, card)
			}
		case p.mainColor:
			mainCardList = append(mainCardList, curCardList...)
		default:
			cardList := make([]core.Card, 0, len(curCardList))
			for _, card := range curCardList {
				if card.Num == 15 {
					mainCardList = append(mainCardList, card)
				} else {
					cardList = append(cardList, card)
				}
			}
			p.cardsByColor[color] = cardList
		}
	}
	for idx := range mainCardList {
		mainCardList[idx].IsMain = true
	}
	p.cardsByColor[p.mainColor] = mainCardList
	delete(p.cardsByColor, def.CardColorNil)

	p.sortCards()
}

func (p *player) sortCards() {
	for color, origCardList := range p.cardsByColor {
		if color == p.mainColor {
			continue
		}
		var cardList []core.Card
		for _, card := range origCardList {
			if card.Num == 15 {
				p.cardsByColor[p.mainColor] = append(p.cardsByColor[p.mainColor], card)
			} else {
				cardList = append(cardList, card)
			}
		}
		p.cardsByColor[color] = SortCards(cardList)
	}
	p.cardsByColor[p.mainColor] = SortCards(p.cardsByColor[p.mainColor])
}

func (p *player) IsFinished() bool {
	for _, cardList := range p.cardsByColor {
		if len(cardList) > 0 {
			return false
		}
	}
	return true
}

func (p *player) GetName() string {
	return p.Name
}

func (p *player) GetTeam() uint8 {
	return p.Team
}

func (p *player) findCardsByType(curCardList []core.Card, cardType uint8) [][]core.Card {
	if len(curCardList) == 0 {
		return nil
	}
	var resultList [][]core.Card
	switch cardType {
	case core.CardTypeSingle:
		resultList = append(resultList, []core.Card{curCardList[0]})
		return append(resultList, p.findCardsByType(curCardList[1:], cardType)...)
	case core.CardTypeDouble:
		if len(curCardList) < 2 {
			break
		}
		for idx := 0; idx < len(curCardList)-1; idx++ {
			if curCardList[idx].Num == curCardList[idx+1].Num {
				resultList = append(resultList, []core.Card{
					curCardList[idx],
					curCardList[idx+1],
				})
				curCardList = append(curCardList[:idx], curCardList[idx+2:]...)
				return append(resultList, p.findCardsByType(curCardList, cardType)...)
			}
		}
	case core.CardTypeConsecutiveDouble:
		if len(curCardList) < 4 {
			break
		}
		for idx := 0; idx < len(curCardList)-3; idx++ {
			if curCardList[idx].Num == curCardList[idx+1].Num &&
				curCardList[idx+1].Num == curCardList[idx+2].Num &&
				curCardList[idx+2].Num == curCardList[idx+3].Num {
				resultList = append(resultList, []core.Card{
					curCardList[idx],
					curCardList[idx+1],
					curCardList[idx+2],
					curCardList[idx+3],
				})
				curCardList = append(curCardList[:idx], curCardList[idx+4:]...)
				return append(resultList, p.findCardsByType(curCardList, cardType)...)
			}
		}
	}
	return resultList
}
