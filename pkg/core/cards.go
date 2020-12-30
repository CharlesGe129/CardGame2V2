package core

import (
	"sort"
	"strings"

	"github.com/CharlesGe129/CardGame2V2/pkg/def"
)

const (
	CardTypeNil               uint8 = 0
	CardTypeSingle            uint8 = 1
	CardTypeDouble            uint8 = 2
	CardTypeConsecutiveDouble uint8 = 3
)

type Cards struct {
	Cards     []Card
	mainColor def.CardColor
}

func NewCards(mainColor def.CardColor, cardList ...Card) Cards {
	cards := Cards{
		mainColor: mainColor,
		Cards:     make([]Card, len(cardList)),
	}
	sort.Slice(cardList, func(i, j int) bool {
		return cardList[i].Num < cardList[j].Num
	})
	l := len(cardList)
	for idx, card := range cardList {
		cards.Cards[l-1-idx] = card
	}
	return cards
}

func ParseCards(pool CardPool, rawStr string) (*Cards, error) {
	var cardList []Card
	cardStrList := strings.Split(rawStr, " ")
	curColor := def.CardColorNil
	for _, cardStr := range cardStrList {
		card, err := pool.ParseCard(cardStr, curColor)
		if err != nil {
			return nil, err
		}
		cardList = append(cardList, *card)
		curColor = card.Color
	}
	cards := NewCards(pool.MainColor, cardList...)
	return &cards, nil
}

func (cards Cards) ValidateSameColor() bool {
	if len(cards.Cards) == 0 {
		return false
	}
	color := def.CardColorNil
	for _, card := range cards.Cards {
		if card.Num == 21 || card.Num == 22 {
			color = cards.mainColor
		} else if color == def.CardColorNil {
			color = card.Color
		} else {
			if card.Color != color {
				return false
			}
		}
	}
	return true
}

// assume same color
func (cards Cards) ParseBiggest() (_ Card, cardType, num uint8) {
	switch len(cards.Cards) {
	case 1:
		return cards.Cards[0], CardTypeSingle, 1
	case 2:
		big := cards.Cards[0]
		if big.Num == cards.Cards[1].Num {
			return big, CardTypeDouble, 1
		} else {
			return big, CardTypeSingle, 1
		}
	case 3:
		big := cards.Cards[1]
		if cards.Cards[0].Num == big.Num {
			return big, CardTypeDouble, 1
		} else if big.Num == cards.Cards[2].Num {
			return big, CardTypeDouble, 1
		} else {
			return cards.Cards[0], CardTypeSingle, 1
		}
	default:
		var doubleList []Card
		// find doubles
		prevCard := cards.Cards[0]
		for _, card := range cards.Cards[1:] {
			if card.Num == prevCard.Num {
				doubleList = append(doubleList, card)
			}
			prevCard = card
		}
		if len(doubleList) == 0 {
			return cards.Cards[0], CardTypeSingle, 1
		}
		var (
			found            bool
			curIsConsecutive bool
			big              = doubleList[0]
			curBig           = doubleList[0]
			curNum           uint8
			num              uint8
		)
		prevCard = doubleList[0]
		// find consecutive double
		for _, card := range doubleList[1:] {
			if prevCard.Num-card.Num == 1 {
				if !found {
					found = true
				}
				if !curIsConsecutive {
					curIsConsecutive = true
					curBig = prevCard
					curNum = 2
				} else {
					curNum++
				}
			} else {
				curIsConsecutive = false
				if curNum > num {
					num = curNum
					big = curBig
				}
			}
			prevCard = card
		}
		if curNum > num {
			num = curNum
			big = curBig
		}
		if found {
			return big, CardTypeConsecutiveDouble, num
		} else {
			return doubleList[0], CardTypeDouble, uint8(len(doubleList))
		}
	}
}

func (cards *Cards) LargerOrEqualTo(others *Cards) bool {
	if !others.ValidateSameColor() {
		return true
	}
	curBig, curType, curNum := cards.ParseBiggest()
	otherBig, otherType, otherNum := others.ParseBiggest()
	switch {
	case curBig.IsMain && !otherBig.IsMain:
		return true
	case !curBig.IsMain && otherBig.IsMain:
		if curType < otherType {
			return false
		} else if curType == otherType && curNum <= otherNum {
			return false
		} else {
			return true
		}
	default:
		if curType > otherType {
			return true
		} else if curType == otherType {
			if curNum > otherNum {
				return true
			} else {
				return curBig.LargerOrEqualTo(otherBig)
			}
		} else {
			return curBig.LargerOrEqualTo(otherBig)
		}
	}
}

func (cards *Cards) String() string {
	str := ""
	for _, card := range cards.Cards {
		str += card.String() + " "
	}
	return str
}

func RemoveCardsFromList(removeCards, cardList []Card) []Card {
	for _, rmCard := range removeCards {
		for idx, curCard := range cardList {
			if rmCard.Num == curCard.Num && rmCard.Color == curCard.Color {
				cardList = append(cardList[:idx], cardList[idx+1:]...)
			}
		}
	}
	return cardList
}
