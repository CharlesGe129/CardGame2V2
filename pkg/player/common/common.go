package common

import (
	"github.com/CharlesGe129/CardGame2V2/pkg/core"
	"github.com/CharlesGe129/CardGame2V2/pkg/def"
	player2 "github.com/CharlesGe129/CardGame2V2/pkg/player"
)

type Player struct {
	Name string
	Team uint8

	MainColor    def.CardColor
	Pool         core.CardPool
	CardsByColor map[def.CardColor][]core.Card
}

func (p *Player) SetMainColor(pool core.CardPool) {
	p.MainColor = pool.MainColor
	p.Pool = pool

	// arrange main cards
	var mainCardList []core.Card
	for color, curCardList := range p.CardsByColor {
		switch color {
		case def.CardColorNil:
			for _, card := range curCardList {
				if card.Color == def.CardColorNil {
					card.Color = p.MainColor
				}
				mainCardList = append(mainCardList, card)
			}
		case p.MainColor:
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
			p.CardsByColor[color] = cardList
		}
	}
	for idx := range mainCardList {
		mainCardList[idx].IsMain = true
	}
	p.CardsByColor[p.MainColor] = mainCardList
	delete(p.CardsByColor, def.CardColorNil)

	p.SortCards()
}

func (p *Player) SortCards() {
	for color, origCardList := range p.CardsByColor {
		if color == p.MainColor {
			continue
		}
		var cardList []core.Card
		for _, card := range origCardList {
			if card.Num == 15 {
				p.CardsByColor[p.MainColor] = append(p.CardsByColor[p.MainColor], card)
			} else {
				cardList = append(cardList, card)
			}
		}
		p.CardsByColor[color] = player2.SortCards(cardList)
	}
	p.CardsByColor[p.MainColor] = player2.SortCards(p.CardsByColor[p.MainColor])
}

func (p *Player) IsFinished() bool {
	for _, cardList := range p.CardsByColor {
		if len(cardList) > 0 {
			return false
		}
	}
	return true
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) GetTeam() uint8 {
	return p.Team
}

func FindCardsByType(curCardList []core.Card, cardType uint8) [][]core.Card {
	if len(curCardList) == 0 {
		return nil
	}
	var resultList [][]core.Card
	switch cardType {
	case core.CardTypeSingle:
		resultList = append(resultList, []core.Card{curCardList[0]})
		return append(resultList, FindCardsByType(curCardList[1:], cardType)...)
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
				return append(resultList, FindCardsByType(curCardList, cardType)...)
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
				return append(resultList, FindCardsByType(curCardList, cardType)...)
			}
		}
	}
	return resultList
}

func (p *Player) AddCard(card core.Card) {
	cardList, ok := p.CardsByColor[card.Color]
	if ok {
		p.CardsByColor[card.Color] = append(cardList, card)
	} else {
		p.CardsByColor[card.Color] = []core.Card{card}
	}
}
