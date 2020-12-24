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

	defendTeam   uint8
	assaultScore uint8
}

func NewGame(playerList [4]player.Player) *Game {
	return &Game{
		Players: playerList,
	}
}

func (g *Game) Start() {
	mainColor, p0 := g.AssignCards()
	g.defendTeam = p0.GetTeam()
	for !g.Players[0].IsFinished() {
		// p0
		shot := p0.NewShot()
		round := core.NewRound(mainColor, *shot)

		// p1
		p1, err := g.nextPlayer(p0)
		if err != nil {
			log.Fatal(err)
		}
		shot = p1.NextShot(round)
		round.AddShot(*shot)

		// p2
		p2, err := g.nextPlayer(p0)
		if err != nil {
			log.Fatal(err)
		}
		shot = p2.NextShot(round)
		round.AddShot(*shot)

		// p3
		p3, err := g.nextPlayer(p0)
		if err != nil {
			log.Fatal(err)
		}
		shot = p3.NextShot(round)
		round.AddShot(*shot)

		team, score, playerName := round.GetResult()
		p0, err = g.getPlayerByName(playerName)
		if err != nil {
			log.Fatal(err)
		}
		if team != g.defendTeam {
			g.assaultScore += score
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
		for _, player := range g.Players {
			if player.GetTeam() != g.defendTeam {
				fmt.Printf("%s ", player.GetName())
			}
		}
	} else {
		fmt.Printf("防守方获胜，进攻方只获得%d分，获胜玩家：", g.assaultScore)
		for _, player := range g.Players {
			if player.GetTeam() == g.defendTeam {
				fmt.Printf("%s ", player.GetName())
			}
		}
	}
	fmt.Println()
}

func (g *Game) AssignCards() (def.CardColor, player.Player) {
	cards := initialCards()
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
	for _, player := range g.Players {
		player.SetMainColor(mainColor)
	}
	g.SetCoveredCards(origCoveredCards, startPlayer)
	return mainColor, startPlayer
}

func (g *Game) SetCoveredCards(cardList []core.Card, player player.Player) {
	g.CoveredCards = player.SetCoveredCards(cardList)
}

func (g *Game) getPlayerByName(name string) (player.Player, error) {
	for _, player := range g.Players {
		if player.GetName() == name {
			return player, nil
		}
	}
	return nil, fmt.Errorf("unable to find player %q", name)
}

func (g *Game) nextPlayer(curPlayer player.Player) (player.Player, error) {
	idx := 0
	for idx < len(g.Players) {
		idx += 1
		if g.Players[idx-1].GetName() == curPlayer.GetName() {
			return g.Players[idx], nil
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

func initialCards() []core.Card {
	var cardList []core.Card
	for num := 0; num < 2; num++ {
		for i := 3; i <= 15; i++ {
			cardList = append(cardList, initialCard(uint8(i))...)
		}
		cardList = append(cardList, core.Card{Num: 21})
		cardList = append(cardList, core.Card{Num: 22})
	}
	return cardList
}
