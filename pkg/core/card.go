package core

import (
	"fmt"
	"strings"

	"github.com/CharlesGe129/CardGame2V2/pkg/def"
)

type Card struct {
	Num   uint8
	Color def.CardColor
}

func ParseCard(rawStr string) (*Card, error) {
	var card Card
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
		return nil, fmt.Errorf("unable to parse card %q: color not found", rawStr)
	}
	card.Num, found = def.MapNameToCard[rawStr]
	if !found {
		return nil, fmt.Errorf("unable to parse card %q: number not found", rawStr)
	}
	return &card, nil
}

func (c Card) Name() string {
	return def.MapCardName[c.Num]
}

func (c Card) Larger(other Card, mainColor def.CardColor) bool {
	if c.Color == other.Color {
		return c.Num > other.Num
	} else if c.Color == mainColor {
		return true
	} else if other.Color == mainColor {
		return false
	} else {
		return false
	}
}

func (c Card) String() string {
	return def.MapColorToCard[c.Color] + def.MapCardName[c.Num]
}
