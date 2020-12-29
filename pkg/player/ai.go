package player

import (
	"fmt"

	"github.com/CharlesGe129/CardGame2V2/pkg/core"
	"github.com/CharlesGe129/CardGame2V2/pkg/def"
	zlog "github.com/rs/zerolog/log"
)

type AiPlayer struct {
	player
}

func NewAiPlayer(name string, team uint8) *AiPlayer {
	return &AiPlayer{
		player: player{
			Name:         name,
			Team:         team,
			cardsByColor: make(map[def.CardColor][]core.Card),
		},
	}
}

func (p *AiPlayer) NextShot(r *core.Round) (*core.Shot, error) {
	big := r.GetBiggest()
	_, bigCards := big.Info()
	num := len(bigCards.Cards)
	color := r.Color

	var cardList []core.Card
	// color
	logger := zlog.Logger.With().Str("player", "[ai][nextShot]").Logger()
	logger.Debug().Msgf("current round biggest: %+v\ncolor=%s, current cards%+v", big, color, p.cardsByColor[color])
	bigCard, bigCardType, _ := bigCards.ParseBiggest()
	cardList = p.GetCards(p.cardsByColor[color], bigCard, bigCardType)
	logger.Debug().Msgf("after adding `get cards`: %+v", cardList)
	if len(cardList) == 0 {
		for _, card := range p.cardsByColor[color] {
			cardList = append(cardList, card)
			if len(cardList) >= num {
				break
			}
		}
	}
	logger.Debug().Msgf("after adding `current color`: %+v", cardList)

	// get from main color
	if len(cardList) < num && color != p.mainColor {
		for _, card := range p.cardsByColor[p.mainColor] {
			cardList = append(cardList, card)
			if len(cardList) >= num {
				break
			}
		}
	}
	logger.Debug().Msgf("after adding `main color`: %+v", cardList)

	// get from other color
	if len(cardList) < num {
	NextColor:
		for curColor, curCardList := range p.cardsByColor {
			if curColor == color || curColor == p.mainColor {
				continue
			}
			for _, card := range curCardList {
				cardList = append(cardList, card)
				if len(cardList) >= num {
					break NextColor
				}
			}
		}
	}
	logger.Debug().Msgf("after adding `other color`: %+v", cardList)

	realCardList, err := p.RemoveCards(cardList)
	if err != nil {
		return nil, err
	}
	cards := core.NewCards(p.mainColor, realCardList...)
	return core.NewShot(cards, p.Team, p.Name), nil
}

func (p *AiPlayer) GetCards(origCardList []core.Card, bigCard core.Card, cardType uint8) []core.Card {
	curCardList := append([]core.Card{}, origCardList...)
	resultList := p.findCardsByType(curCardList, cardType)
	if len(resultList) == 0 {
		return nil
	}
	myBigCards := core.NewCards(p.mainColor, resultList[0]...)
	myBig, _, _ := myBigCards.ParseBiggest()
	if !bigCard.LargerOrEqualTo(myBig) {
		return resultList[0]
	} else {
		return resultList[len(resultList)-1]
	}
}

func (p *AiPlayer) NewShot() *core.Shot {
	for _, cardList := range p.cardsByColor {
		if len(cardList) > 0 {
			realCardList, _ := p.RemoveCards([]core.Card{cardList[0]})
			cards := core.NewCards(p.mainColor, realCardList...)
			return core.NewShot(cards, p.Team, p.Name)
		}
	}
	return nil
}

func (p *AiPlayer) sortCards() {
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

func (p *AiPlayer) BidMainColor() def.CardColor {
	// TODO: bid main color
	return def.CardColorSpade
}

func (p *AiPlayer) ShowCards() {
	p.sortCards()
	mainCards := p.cardsByColor[p.mainColor]
	fmt.Printf("[王]" + string(p.mainColor) + ": ")
	for _, card := range mainCards {
		fmt.Printf(card.Name())
	}
	fmt.Println()
	for color, cards := range p.cardsByColor {
		if color == p.mainColor {
			continue
		}
		fmt.Printf(string(color) + ": ")
		for _, card := range cards {
			fmt.Printf(card.Name())
		}
		fmt.Println("")
	}
}

func (p *AiPlayer) AddCard(card core.Card) {
	cardList, ok := p.cardsByColor[card.Color]
	if ok {
		p.cardsByColor[card.Color] = append(cardList, card)
	} else {
		p.cardsByColor[card.Color] = []core.Card{card}
	}
}

func (p *AiPlayer) RemoveCards(rawCardList []core.Card) ([]core.Card, error) {
	cardsByColor := make(map[def.CardColor][]core.Card)
	for color, cards := range p.cardsByColor {
		cardsByColor[color] = append([]core.Card{}, cards...)
	}
	realCardList := make([]core.Card, 0, len(rawCardList))
	for _, card := range rawCardList {
		var curColor def.CardColor
		if card.IsMain {
			curColor = p.mainColor
		} else {
			curColor = card.Color
		}
		curCardList, ok := cardsByColor[curColor]
		if !ok {
			return nil, fmt.Errorf("player %s doesn't have any cards of color %q; cards to be removed: %s",
				p.Name, card.Color, rawCardList)
		}
		for idx, curCard := range curCardList {
			if curCard.Num == card.Num {
				realCardList = append(realCardList, curCard)
				curCardList = append(curCardList[:idx], curCardList[idx+1:]...)
				break
			}
		}
		cardsByColor[curColor] = curCardList
	}
	if len(realCardList) != len(rawCardList) {
		fmt.Printf("%+v\n", p.cardsByColor)
		return nil, fmt.Errorf("unable to find cards: %v", rawCardList)
	}
	p.cardsByColor = cardsByColor
	return realCardList, nil
}

func (p *AiPlayer) SetCoveredCards(origCoveredCards []core.Card) ([]core.Card, error) {
	for _, card := range origCoveredCards {
		p.AddCard(card)
	}
	p.sortCards()
	var coveredCards []core.Card
NextColor:
	for color, cardList := range p.cardsByColor {
		if color == p.mainColor {
			continue
		}
		tmpCardList := append([]core.Card{}, cardList...)
		for _, tmpCard := range tmpCardList {
			if HasDouble(tmpCardList, tmpCard.Num) {
				continue
			}
			if tmpCard.Num == 14 { // A
				continue
			}
			coveredCards = append(coveredCards, tmpCard)
			if len(coveredCards) == 8 {
				break NextColor
			}
		}
	}
	if len(coveredCards) < 8 {
		mainCardList := p.cardsByColor[p.mainColor]
		idx := len(mainCardList) - 1
		for len(coveredCards) < 8 {
			coveredCards = append(coveredCards, mainCardList[idx])
			idx--
		}
	}
	realCardList, err := p.RemoveCards(coveredCards)
	if err != nil {
		return nil, err
	}
	return realCardList, nil
}
