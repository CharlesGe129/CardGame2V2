package core

import (
	"fmt"
	"strings"

	"github.com/CharlesGe129/CardGame2V2/pkg/def"
)

type CardPool struct {
	MainColor def.CardColor
}

func NewCardPool(mainColor def.CardColor) *CardPool {
	return &CardPool{
		MainColor: mainColor,
	}
}

func (pool *CardPool) ParseCard(rawStr string, curColor def.CardColor) (*Card, error) {
	var card Card
	switch rawStr {
	case "大":
		card.Num = 22
		return pool.getRealCard(card), nil
	case "小":
		card.Num = 21
		card.IsMain = true
		return pool.getRealCard(card), nil
	}
	found := false
	for prefix, color := range def.MapCardColor {
		if strings.HasPrefix(rawStr, prefix) {
			card.Color = color
			rawStr = strings.Split(rawStr, prefix)[1]
			found = true
			break
		}
	}
	if !found {
		if curColor != def.CardColorNil {
			card.Color = curColor
		} else {
			return nil, fmt.Errorf("unable to parse card %q: color not found", rawStr)
		}
	}
	card.Num, found = def.MapNameToCard[rawStr]
	if !found {
		return nil, fmt.Errorf("unable to parse card %q: number not found", rawStr)
	}
	return pool.getRealCard(card), nil
}

func (pool *CardPool) GetRealCardList(cardList ...Card) []Card {
	realCardList := make([]Card, len(cardList))
	for idx, card := range cardList {
		realCardList[idx] = *pool.getRealCard(card)
	}
	return realCardList
}

func (pool *CardPool) getRealCard(card Card) *Card {
	if card.Num > 20 {
		return &Card{
			Num:    card.Num,
			Color:  pool.MainColor,
			IsMain: true,
		}
	} else if card.Num == 15 {
		return &Card{
			Num:    card.Num,
			Color:  card.Color,
			IsMain: true,
		}
	} else {
		return &Card{
			Num:    card.Num,
			Color:  card.Color,
			IsMain: card.Color == pool.MainColor,
		}
	}
}
