package core

import (
	"fmt"

	"github.com/CharlesGe129/CardGame2V2/pkg/def"
)

type Card struct {
	Num    uint8
	Color  def.CardColor
	IsMain bool
}

func (c Card) Name() string {
	return def.MapCardName[c.Num]
}

func (c Card) LargerOrEqualTo(other Card) bool {
	if c.IsMain && other.IsMain {
		return c.Num >= other.Num
	} else if c.IsMain && !other.IsMain {
		return true
	} else if !c.IsMain && other.IsMain {
		return false
	}
	if c.Color == other.Color {
		return c.Num >= other.Num
	} else {
		return true
	}
}

func (c Card) String() string {
	if _, ok := def.MapCardName[c.Num]; !ok {
		fmt.Printf("unable to find card name: %d\n", c.Num)
	}
	return def.MapColorToCard[c.Color] + def.MapCardName[c.Num]
}
