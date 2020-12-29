package player

import (
	"sort"

	"github.com/CharlesGe129/CardGame2V2/pkg/core"
)

func SortCards(cardList []core.Card) []core.Card {
	sort.Slice(cardList, func(i, j int) bool {
		return cardList[i].Num < cardList[j].Num
	})
	revCardList := make([]core.Card, len(cardList))
	for idx, card := range cardList {
		revCardList[len(cardList)-1-idx] = card
	}
	return revCardList
}

func HasDouble(cardList []core.Card, num uint8) bool {
	var count int
	for _, card := range cardList {
		if card.Num == num {
			count++
		}
	}
	return count >= 2
}
