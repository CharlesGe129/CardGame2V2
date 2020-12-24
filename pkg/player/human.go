package player

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/CharlesGe129/CardGame2V2/pkg/core"
	"github.com/CharlesGe129/CardGame2V2/pkg/def"
)

type HumanPlayer struct {
	Name string
	Team uint8

	mainColor    def.CardColor
	cardsByColor map[def.CardColor][]core.Card
	scanner      *bufio.Scanner
}

func NewPlayer(name string, team uint8) *HumanPlayer {
	return &HumanPlayer{
		Name:         name,
		Team:         team,
		cardsByColor: make(map[def.CardColor][]core.Card),
		scanner:      bufio.NewScanner(os.Stdin),
	}
}

func (p *HumanPlayer) SetMainColor(color def.CardColor) {
	p.mainColor = color
	for color, cardList := range p.cardsByColor {
		if color == p.mainColor {
			cardList = append(cardList, p.cardsByColor[def.CardColorNil]...)
		}
		p.cardsByColor[color] = cardList
	}
	delete(p.cardsByColor, def.CardColorNil)
	p.sortCards()
}

func (p *HumanPlayer) sortCards() {
	for color, cardList := range p.cardsByColor {
		sort.Slice(cardList, func(i, j int) bool {
			return cardList[i].Num < cardList[j].Num
		})
		p.cardsByColor[color] = cardList
	}
}

func (p *HumanPlayer) BidMainColor() def.CardColor {
	p.ShowCards()
	for {
		rawStr := p.getInput("请输入王花色: ")
		if color, ok := def.MapCardColor[rawStr]; ok {
			return color
		}
	}
}

func (p *HumanPlayer) ShowCards() {
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

func (p *HumanPlayer) AddCard(card core.Card) {
	cardList, ok := p.cardsByColor[card.Color]
	if ok {
		p.cardsByColor[card.Color] = append(cardList, card)
	} else {
		p.cardsByColor[card.Color] = []core.Card{card}
	}
}

func (p *HumanPlayer) NextShot(r *core.Round) *core.Shot {
	for {
		rawStr := p.getInput("请出牌:")
		switch rawStr {
		case "show":
			if r != nil {
				r.ShowShots()
			}
		default:
			cards, err := core.ParseCards(p.mainColor, rawStr)
			if err != nil {
				fmt.Printf("出错了: \n%s\n\n", err)
				continue
			}
			return core.NewShot(*cards, p.Team, p.Name)
		}
	}
}

func (p *HumanPlayer) getInput(msg string) string {
	for {
		fmt.Println(msg)
		p.scanner.Scan()
		if err := p.scanner.Err(); err != nil {
			fmt.Printf("出错了: \n%s\n\n", err)
			continue
		}
		return p.scanner.Text()
	}
}

func (p *HumanPlayer) NewShot() *core.Shot {
	return p.NextShot(nil)
}

func (p *HumanPlayer) RemoveCards(cardList []core.Card) error {
	cardsByColor := make(map[def.CardColor][]core.Card)
	for color, cardList := range p.cardsByColor {
		cardsByColor[color] = append([]core.Card{}, cardList...)
	}
	for _, card := range cardList {
		curCardList, ok := cardsByColor[card.Color]
		if !ok {
			return fmt.Errorf("player %s doesn't have any cards of color %q; cards to be removed: %s",
				p.Name, card.Color, cardList)
		}
		for idx, curCard := range curCardList {
			if curCard.Num == card.Num {
				if idx+1 <= len(curCardList) {
					curCardList = append(curCardList[:idx], curCardList[idx+1:]...)
				} else {
					curCardList = curCardList[:idx]
				}
			}
		}
	}
	p.cardsByColor = cardsByColor
	return nil
}

func (p *HumanPlayer) SetCoveredCards(origCoveredCards []core.Card) []core.Card {
	for _, card := range origCoveredCards {
		p.AddCard(card)
	}
	p.sortCards()
	p.ShowCards()
	for {
		rawStr := p.getInput("请扣底牌")
		newCoveredCards, err := core.ParseCards(p.mainColor, rawStr)
		if err != nil {
			fmt.Printf("出错了: \n%s\n\n", err)
			continue
		}
		if err := p.RemoveCards(newCoveredCards.Cards); err != nil {
			fmt.Printf("出错了: \n%s\n\n", err)
			continue
		}
		return newCoveredCards.Cards
	}
}

func (p *HumanPlayer) IsFinished() bool {
	for _, cardList := range p.cardsByColor {
		if len(cardList) > 0 {
			return false
		}
	}
	return true
}

func (p *HumanPlayer) GetName() string {
	return p.Name
}

func (p *HumanPlayer) GetTeam() uint8 {
	return p.Team
}
