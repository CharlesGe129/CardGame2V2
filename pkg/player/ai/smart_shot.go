package ai

import (
	"fmt"
	"sort"

	"github.com/CharlesGe129/CardGame2V2/pkg/core"
	"github.com/CharlesGe129/CardGame2V2/pkg/def"
	"github.com/CharlesGe129/CardGame2V2/pkg/player/common"
	"github.com/pkg/errors"
)

type SmartShot struct {
	mainColor    def.CardColor
	cardsByColor map[def.CardColor][]core.Card
}

func NewSmartShot(mainColor def.CardColor, cardsByColor map[def.CardColor][]core.Card) *SmartShot {
	shot := SmartShot{
		mainColor:    mainColor,
		cardsByColor: make(map[def.CardColor][]core.Card),
	}
	for color, cardList := range cardsByColor {
		shot.cardsByColor[color] = append([]core.Card{}, cardList...)
	}
	return &shot
}

func (p *SmartShot) HasCard(origShot core.Shot) bool {
	_, cards := origShot.Info()
	card := cards.Cards[0]
	if card.IsMain {
		return len(p.cardsByColor[p.mainColor]) > 0
	} else {
		return len(p.cardsByColor[card.Color]) > 0
	}
}

func (p *SmartShot) NextShot(origShot, bigShot core.Shot) ([]core.Card, error) {
	// has cards
	if p.HasCard(origShot) {
		bigShot = origShot
	}

	_, bigCards := bigShot.Info()
	bigCard, bigCardType, bigTypeNum := bigCards.ParseBiggest()
	bigCount := len(bigCards.Cards)
	shotColor := bigCard.Color
	if bigCard.IsMain {
		shotColor = p.mainColor
	}
	if shotColor != p.mainColor && len(p.cardsByColor[shotColor]) == 0 {
		fromCardList := p.cardsByColor[p.mainColor]
		cardList := p.nextBiggerShotFromMain(bigCard, bigCardType, bigTypeNum, bigCount, fromCardList, nil)
		if len(cardList) == bigCount {
			return cardList, nil
		}
	}
	return p.nextShotFromNormal(bigShot)
}

func (p *SmartShot) nextBiggerShotFromMain(bigCard core.Card, bigType, bigTypeNum uint8, bigCount int,
	fromCardList, cardList []core.Card) []core.Card {
	if bigCount == 0 {
		return nil
	}
	if len(fromCardList) < bigCount {
		return nil
	}
	copiedCardList := append([]core.Card{}, fromCardList...)
	candidateCardLists := common.FindCardsByType(copiedCardList, bigType)
	if len(candidateCardLists) < int(bigTypeNum) {
		return nil
	}
	myBigCards := core.NewCards(p.mainColor, candidateCardLists[0]...)
	myBig, _, _ := myBigCards.ParseBiggest()
	if bigCard.LargerOrEqualTo(myBig) {
		return nil
	}
	// add main cards
	for _, newCardList := range candidateCardLists {
		bigTypeNum -= 1
		bigCount -= len(newCardList)
		fromCardList = p.RemoveCards(fromCardList, newCardList)
		cardList = append(cardList, newCardList...)
		if bigTypeNum == 0 {
			break
		}
	}
	return p.nextBiggerShotFromMain(bigCard, bigType, bigTypeNum, bigCount, fromCardList, cardList)
}

func (p *SmartShot) nextShotFromNormal(shot core.Shot) ([]core.Card, error) {
	_, bigCards := shot.Info()
	bigCard, bigType, bigTypeNum := bigCards.ParseBiggest()
	bigCount := len(bigCards.Cards)
	shotColor := bigCard.Color
	if bigCard.IsMain {
		shotColor = p.mainColor
	}
	var cardList []core.Card

	// shot from color
	cardList = p.nextShotFromCardList(bigCard, bigType, bigTypeNum, bigCount,
		p.cardsByColor[shotColor], cardList, true)
	if len(cardList) == bigCount {
		return cardList, nil
	}
	// add cards from other colors
	for color, fromCardList := range p.cardsByColor {
		if color == shotColor || color == p.mainColor {
			continue
		}
		cardList = p.nextShotFromCardList(bigCard, bigType, bigTypeNum, bigCount,
			fromCardList, cardList, false)
		if len(cardList) == bigCount {
			return cardList, nil
		}
	}
	// add cards from main
	if shotColor != p.mainColor {
		cardList = p.nextShotFromCardList(bigCard, bigType, bigTypeNum, bigCount,
			p.cardsByColor[p.mainColor], cardList, false)
	}
	if bigCount != len(cardList) {
		return nil, errors.Errorf("unable to nextShotFromNormal(); shot=%+v\n, cardList=%+v", shot, cardList)
	}
	return cardList, nil
}

func (p *SmartShot) nextShotFromCardList(bigCard core.Card, bigType, bigTypeNum uint8, bigCount int,
	fromCardList, cardList []core.Card, mustDouble bool) []core.Card {
	if len(fromCardList) == 0 {
		return cardList
	}
	if (bigCard.IsMain && !fromCardList[0].IsMain) || (bigCard.Color != fromCardList[0].Color) || !mustDouble {
		return p.nextShotSmallest(bigCount, fromCardList, cardList)
	} else {
		// try to be bigger
		copiedCardList := append([]core.Card{}, fromCardList...)
		candidateCardLists := common.FindCardsByType(copiedCardList, bigType)
		if len(candidateCardLists) == 0 {
			return p.nextShotSmallest(bigCount, fromCardList, cardList)
		}
		myBigCards := core.NewCards(p.mainColor, candidateCardLists[0]...)
		myBig, _, _ := myBigCards.ParseBiggest()
		if !bigCard.LargerOrEqualTo(myBig) && len(candidateCardLists) >= int(bigTypeNum) {
			// I'm bigger
			for _, newCardList := range candidateCardLists {
				bigTypeNum -= 1
				fromCardList = p.RemoveCards(fromCardList, newCardList)
				cardList = append(cardList, newCardList...)
				if bigTypeNum == 0 {
					break
				}
			}
			return p.nextShotSmallest(bigCount, fromCardList, cardList)
		} else if mustDouble {
			// must double
			for idx := range candidateCardLists {
				newCardList := candidateCardLists[len(candidateCardLists)-1-idx]
				bigTypeNum -= 1
				fromCardList = p.RemoveCards(fromCardList, newCardList)
				cardList = append(cardList, newCardList...)
				if bigTypeNum == 0 {
					break
				}
			}
			return p.nextShotSmallest(bigCount, fromCardList, cardList)
		} else {
			return p.nextShotSmallest(bigCount, fromCardList, cardList)
		}
	}
}

func (p *SmartShot) nextShotSmallest(count int, fromCardList, cardList []core.Card) []core.Card {
	if count == 0 || count == len(cardList) {
		return cardList
	}
	fromCardList = p.sortCardListFromSmallToScore(fromCardList)
	for _, card := range fromCardList {
		cardList = append(cardList, card)
		if len(cardList) == count {
			break
		}
	}
	return cardList
}

func (p *SmartShot) sortCardListFromSmallToScore(cardList []core.Card) []core.Card {
	sort.Slice(cardList, func(i, j int) bool {
		card1 := cardList[i]
		card2 := cardList[j]
		switch def.MapCardName[card1.Num] {
		case "k":
			return false
		case "0":
			return def.MapCardName[card2.Num] == "k"
		case "5":
			return def.MapCardName[card2.Num] == "k" || def.MapCardName[card2.Num] == "0"
		default:
			name2 := def.MapCardName[card2.Num]
			if name2 == "k" || name2 == "0" || name2 == "5" {
				return true
			} else {
				return card1.Num <= card2.Num
			}
		}
	})
	return cardList
}

func (p *SmartShot) RemoveCards(curCardList, rmCardList []core.Card) []core.Card {
	fmt.Printf("remove %+v\n from %+v\n", rmCardList, curCardList)
	for _, rmCard := range rmCardList {
		for idx, card := range curCardList {
			if rmCard.Num == card.Num && rmCard.Color == card.Color {
				curCardList = append(curCardList[:idx], curCardList[idx+1:]...)
				break
			}
		}
	}
	return curCardList
}
