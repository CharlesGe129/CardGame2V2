package pkg

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/CharlesGe129/CardGame2V2/pkg/core"
	"github.com/CharlesGe129/CardGame2V2/pkg/def"
	"github.com/CharlesGe129/CardGame2V2/pkg/player"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Game struct {
	Players      [4]player.Player
	CoveredCards []core.Card
	pool         *core.CardPool

	mainCardName    string
	mainCardOrigNum uint8
	defendTeam      uint8
	assaultScore    uint8
}

func NewGame(playerList [4]player.Player, mainCard string) *Game {
	game := Game{
		Players: playerList,
	}
	// set main card
	if _, ok := def.MapNameToCard[mainCard]; ok {
		game.mainCardName = mainCard
	} else {
		game.mainCardName = "2"
	}
	for num, name := range def.MapCardName {
		if name == mainCard {
			game.mainCardOrigNum = num
			delete(def.MapCardName, num)
			break
		}
	}
	def.MapCardName[15] = mainCard

	def.Init()
	fmt.Printf("游戏准备完毕，本轮牌库: %+v\n\n", def.MapCardName)
	return &game
}

func (g *Game) Start() {
	p0, err := g.AssignCards()
	if err != nil {
		log.Fatal(err)
	}
	g.defendTeam = g.Players[1].GetTeam()
	for !g.Players[0].IsFinished() {
		// p0
		shot := p0.NewShot()
		round := core.NewRound(*shot)

		// p1
		p1, err := g.nextPlayer(p0)
		if err != nil {
			log.Fatal(err)
		}
		shot, err = p1.NextShot(round)
		if err != nil {
			log.Fatal(err)
		}
		round.AddShot(*shot)

		// p2
		p2, err := g.nextPlayer(p1)
		if err != nil {
			log.Fatal(err)
		}
		shot, err = p2.NextShot(round)
		if err != nil {
			log.Fatal(err)
		}
		round.AddShot(*shot)

		// p3
		p3, err := g.nextPlayer(p2)
		if err != nil {
			log.Fatal(err)
		}
		shot, err = p3.NextShot(round)
		if err != nil {
			log.Fatal(err)
		}
		round.AddShot(*shot)

		team, score, playerName := round.GetResult()
		p0, err = g.getPlayerByName(playerName)
		if err != nil {
			log.Fatal(err)
		}
		if team != g.defendTeam {
			g.assaultScore += score
			fmt.Printf("进攻方得到%d分，目前共有%d分\n", score, g.assaultScore)
		} else {
			fmt.Printf("防守方逃脱%d分，目前进攻方共有%d分\n", score, g.assaultScore)
		}
	}
	// get result
	if p0.GetTeam() != g.defendTeam {
		var coveredScores uint8
		for _, card := range g.CoveredCards {
			switch card.Num {
			case 5:
				coveredScores += 10
			case 10, 13:
				coveredScores += 20
			}
		}
		fmt.Printf("进攻方挖底成功，原有%d分，底牌中获得%d分\n", g.assaultScore, coveredScores)
		g.assaultScore += coveredScores
	}

	if g.assaultScore >= 80 {
		fmt.Printf("进攻方获胜，共得到%d分，获胜玩家：", g.assaultScore)
		for _, p := range g.Players {
			if p.GetTeam() != g.defendTeam {
				fmt.Printf("%s ", p.GetName())
			}
		}
	} else {
		fmt.Printf("防守方获胜，进攻方只获得%d分，获胜玩家：", g.assaultScore)
		for _, p := range g.Players {
			if p.GetTeam() == g.defendTeam {
				fmt.Printf("%s ", p.GetName())
			}
		}
	}
	fmt.Println()
}

func (g *Game) AssignCards() (player.Player, error) {
	cards := initialCards(g.mainCardOrigNum)
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	for i := 0; i < len(cards)-8; i++ {
		g.Players[i%4].AddCard(cards[i])
		time.Sleep(time.Millisecond * 1)
	}
	origCoveredCards := cards[len(cards)-8:]
	// TODO: bid for the main color
	startPlayer := g.Players[0]
	mainColor := startPlayer.BidMainColor()
	g.pool = core.NewCardPool(mainColor)
	for _, p := range g.Players {
		p.SetMainColor(*g.pool)
	}
	if err := g.SetCoveredCards(origCoveredCards, startPlayer); err != nil {
		return nil, err
	}
	return startPlayer, nil
}

func (g *Game) SetCoveredCards(cardList []core.Card, player player.Player) error {
	coveredCards, err := player.SetCoveredCards(cardList)
	if err != nil {
		return err
	}
	g.CoveredCards = coveredCards
	return nil
}

func (g *Game) getPlayerByName(name string) (player.Player, error) {
	for _, p := range g.Players {
		if p.GetName() == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("unable to find player %q", name)
}

func (g *Game) nextPlayer(curPlayer player.Player) (player.Player, error) {
	for idx := 0; idx < len(g.Players); idx++ {
		if g.Players[idx].GetName() == curPlayer.GetName() {
			if idx == len(g.Players)-1 {
				return g.Players[0], nil
			} else {
				return g.Players[idx+1], nil
			}
		}
	}
	return nil, errors.New("unable to find next player")
}

func initialCard(n uint8) []core.Card {
	return []core.Card{
		{
			Num:   n,
			Color: def.CardColorSpade,
		},
		{
			Num:   n,
			Color: def.CardColorHeart,
		},
		{
			Num:   n,
			Color: def.CardColorClub,
		},
		{
			Num:   n,
			Color: def.CardColorDiamond,
		},
	}
}

func initialCards(mainCardNum uint8) []core.Card {
	var cardList []core.Card
	for num := 0; num < 2; num++ {
		for i := 2; i <= 15; i++ {
			if uint8(i) == mainCardNum {
				continue
			}
			cardList = append(cardList, initialCard(uint8(i))...)
		}
		cardList = append(cardList, core.Card{Num: 21, IsMain: true})
		cardList = append(cardList, core.Card{Num: 22, IsMain: true})
	}
	return cardList
}
