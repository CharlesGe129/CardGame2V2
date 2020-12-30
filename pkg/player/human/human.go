package human

import (
	"bufio"
	"fmt"
	"os"

	"github.com/CharlesGe129/CardGame2V2/pkg/core"
	"github.com/CharlesGe129/CardGame2V2/pkg/def"
	"github.com/CharlesGe129/CardGame2V2/pkg/player/common"
)

type HumanPlayer struct {
	common.Player
	scanner *bufio.Scanner
}

func NewHumanPlayer(name string, team uint8) *HumanPlayer {
	return &HumanPlayer{
		Player: common.Player{
			Name:         name,
			Team:         team,
			CardsByColor: make(map[def.CardColor][]core.Card),
		},
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (p *HumanPlayer) BidMainColor() def.CardColor {
	p.ShowCards()
	for {
		rawStr := p.getInput("请输入王花色(hei, hong, fang, cao): ")
		if color, ok := def.MapCardColor[rawStr]; ok {
			return color
		}
	}
}

func (p *HumanPlayer) ShowCards() {
	p.SortCards()
	mainCards := p.CardsByColor[p.MainColor]
	fmt.Printf("[王]" + def.MapColorZnCh[p.MainColor] + ": ")
	for _, card := range mainCards {
		if card.Num != 15 {
			fmt.Printf(card.Name())
		} else {
			fmt.Printf(card.Name() + "(" + def.MapColorZnCh[card.Color] + ")")
		}
	}
	fmt.Println()
	for color, cards := range p.CardsByColor {
		if color == p.MainColor {
			continue
		}
		fmt.Printf(def.MapColorZnCh[color] + ": ")
		for _, card := range cards {
			fmt.Printf(card.Name())
		}
		fmt.Println("")
	}
}

func (p *HumanPlayer) AddCard(card core.Card) {
	cardList, ok := p.CardsByColor[card.Color]
	if ok {
		p.CardsByColor[card.Color] = append(cardList, card)
	} else if p.MainColor != def.CardColorNil {
		p.CardsByColor[p.MainColor] = append(p.CardsByColor[p.MainColor], card)
	} else {
		p.CardsByColor[card.Color] = []core.Card{card}
	}
}

func (p *HumanPlayer) NextShot(r *core.Round) (*core.Shot, error) {
	for {
		rawStr := p.getInput("请出牌:")
		switch rawStr {
		case "round", "r":
			if r != nil {
				r.ShowShots()
			}
		case "show", "s":
			p.ShowCards()
		case "help", "h":
			fmt.Println("输入`round`查看本轮其他玩家已打的牌，输入`show`查看自己手牌")
		default:
			cards, err := core.ParseCards(p.Pool, rawStr)
			if err != nil {
				fmt.Printf("出错了: \n%s\n\n", err)
				continue
			}
			realCardList, err := p.RemoveCards(cards.Cards)
			if err != nil {
				fmt.Printf("出错了: \n%s\n\n", err)
				continue
			}
			realCards := core.NewCards(p.MainColor, realCardList...)
			return core.NewShot(realCards, p.Team, p.Name), nil
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
	shot, _ := p.NextShot(nil)
	return shot
}

func (p *HumanPlayer) RemoveCards(rawCardList []core.Card) ([]core.Card, error) {
	cardsByColor := make(map[def.CardColor][]core.Card)
	for color, cardList := range p.CardsByColor {
		cardsByColor[color] = append([]core.Card{}, cardList...)
	}
	realCardList := make([]core.Card, 0, len(rawCardList))
	for _, card := range rawCardList {
		color := card.Color
		if card.IsMain {
			color = p.MainColor
		}
		curCardList := cardsByColor[color]
		for idx, curCard := range curCardList {
			if curCard.Num == card.Num {
				realCardList = append(realCardList, curCard)
				curCardList = append(curCardList[:idx], curCardList[idx+1:]...)
				break
			}
		}
		cardsByColor[color] = curCardList
	}
	if len(realCardList) != len(rawCardList) {
		return nil, fmt.Errorf("unable to find cards: %v", rawCardList)
	}
	p.CardsByColor = cardsByColor
	return realCardList, nil
}

func (p *HumanPlayer) SetCoveredCards(origCoveredCards []core.Card) ([]core.Card, error) {
	for _, card := range origCoveredCards {
		p.AddCard(card)
	}
	p.SortCards()
	p.ShowCards()
	for {
		rawStr := p.getInput("请扣底牌")
		newCoveredCards, err := core.ParseCards(p.Pool, rawStr)
		if err != nil {
			fmt.Printf("出错了: \n%s\n\n", err)
			continue
		}
		if len(newCoveredCards.Cards) != 8 {
			fmt.Printf("底牌数量不对: %v\n", newCoveredCards)
			continue
		}
		realCardList, err := p.RemoveCards(newCoveredCards.Cards)
		if err != nil {
			fmt.Printf("出错了: \n%s\n\n", err)
			continue
		}
		return realCardList, nil
	}
}
