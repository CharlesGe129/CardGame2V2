package player

import (
	"testing"

	"github.com/CharlesGe129/CardGame2V2/pkg/core"
	"github.com/CharlesGe129/CardGame2V2/pkg/def"
	"github.com/stretchr/testify/require"
)

func TestAI_NextShot_HasCard(t *testing.T) {
	player := NewAiPlayer("测试", 1)
	cardList := []core.Card{
		{Num: 10, Color: def.CardColorClub},
		{Num: 10, Color: def.CardColorClub},
	}
	for _, card := range cardList {
		player.AddCard(card)
	}

	cards := core.NewCards(def.CardColorClub, core.Card{Num: 11, Color: def.CardColorClub})
	shot := core.NewShot(cards, 1, "")
	round := core.NewRound(*shot)
	shotGot, err := player.NextShot(round)
	require.NoError(t, err)
	_, cardsGot := shotGot.Info()
	require.Equal(t, []core.Card{cardList[0]}, cardsGot.Cards)
	require.Equal(t, []core.Card{cardList[0]}, player.cardsByColor[def.CardColorClub])

	// double
	player = NewAiPlayer("测试", 1)
	cardList = []core.Card{
		{Num: 12, Color: def.CardColorClub},
		{Num: 11, Color: def.CardColorClub},
		{Num: 10, Color: def.CardColorClub},
		{Num: 10, Color: def.CardColorClub},
	}
	for _, card := range cardList {
		player.AddCard(card)
	}

	cards = core.NewCards(def.CardColorClub,
		core.Card{Num: 9, Color: def.CardColorClub},
		core.Card{Num: 9, Color: def.CardColorClub},
	)
	shot = core.NewShot(cards, 1, "")
	round = core.NewRound(*shot)
	shotGot, err = player.NextShot(round)
	require.NoError(t, err)
	_, cardsGot = shotGot.Info()
	expectCards := []core.Card{
		{Num: 10, Color: def.CardColorClub},
		{Num: 10, Color: def.CardColorClub},
	}
	require.Equal(t, expectCards, cardsGot.Cards)
}

func TestAI_NextShot_NotEnoughCard(t *testing.T) {
	// 1+1main
	player := NewAiPlayer("测试", 1)
	cardList := []core.Card{
		{Num: 10, Color: def.CardColorClub, IsMain: false},
		{Num: 9, Color: def.CardColorDiamond, IsMain: true},
	}
	for _, card := range cardList {
		player.AddCard(card)
	}
	pool := core.NewCardPool(def.CardColorDiamond)
	player.SetMainColor(*pool)

	cards := core.NewCards(def.CardColorClub,
		core.Card{Num: 11, Color: def.CardColorClub},
		core.Card{Num: 11, Color: def.CardColorClub})
	shot := core.NewShot(cards, 1, "")
	round := core.NewRound(*shot)
	shotGot, err := player.NextShot(round)
	require.NoError(t, err)
	_, cardsGot := shotGot.Info()
	require.Equal(t, cardList, cardsGot.Cards)
	require.Empty(t, player.cardsByColor[def.CardColorClub])
	require.Empty(t, player.cardsByColor[def.CardColorDiamond])

	// 0+2main
	player = NewAiPlayer("测试", 1)
	cardList = []core.Card{
		{Num: 10, Color: def.CardColorDiamond, IsMain: true},
		{Num: 9, Color: def.CardColorDiamond, IsMain: true},
	}
	for _, card := range cardList {
		player.AddCard(card)
	}
	pool = core.NewCardPool(def.CardColorDiamond)
	player.SetMainColor(*pool)

	cards = core.NewCards(def.CardColorClub,
		core.Card{Num: 11, Color: def.CardColorClub},
		core.Card{Num: 11, Color: def.CardColorClub})
	shot = core.NewShot(cards, 1, "")
	round = core.NewRound(*shot)
	shotGot, err = player.NextShot(round)
	require.NoError(t, err)
	_, cardsGot = shotGot.Info()
	require.Equal(t, cardList, cardsGot.Cards)
	require.Empty(t, player.cardsByColor[def.CardColorClub])
	require.Empty(t, player.cardsByColor[def.CardColorDiamond])

	// 0+0main+2
	player = NewAiPlayer("测试", 1)
	cardList = []core.Card{
		{Num: 11, Color: def.CardColorHeart, IsMain: false},
		{Num: 11, Color: def.CardColorHeart, IsMain: false},
	}
	for _, card := range cardList {
		player.AddCard(card)
	}
	pool = core.NewCardPool(def.CardColorDiamond)
	player.SetMainColor(*pool)

	cards = core.NewCards(def.CardColorClub,
		core.Card{Num: 11, Color: def.CardColorClub},
		core.Card{Num: 11, Color: def.CardColorClub})
	shot = core.NewShot(cards, 1, "")
	round = core.NewRound(*shot)
	shotGot, err = player.NextShot(round)
	require.NoError(t, err)
	_, cardsGot = shotGot.Info()
	require.Equal(t, cardList, cardsGot.Cards)
	require.Empty(t, player.cardsByColor[def.CardColorClub])
	require.Empty(t, player.cardsByColor[def.CardColorDiamond])
}

func TestAI_NextShot_MainCards(t *testing.T) {
	// 2main
	player := NewAiPlayer("测试", 1)
	cardList := []core.Card{
		{Num: 10, Color: def.CardColorClub, IsMain: true},
		{Num: 9, Color: def.CardColorClub, IsMain: true},
	}
	for _, card := range cardList {
		player.AddCard(card)
	}
	pool := core.NewCardPool(def.CardColorClub)
	player.SetMainColor(*pool)

	cards := core.NewCards(def.CardColorClub,
		core.Card{Num: 11, Color: def.CardColorClub},
		core.Card{Num: 11, Color: def.CardColorClub})
	shot := core.NewShot(cards, 1, "")
	round := core.NewRound(*shot)
	shotGot, err := player.NextShot(round)
	require.NoError(t, err)
	_, cardsGot := shotGot.Info()
	require.Equal(t, cardList, cardsGot.Cards)
	require.Empty(t, player.cardsByColor[def.CardColorClub])

	// 1main+1
	player = NewAiPlayer("测试", 1)
	cardList = []core.Card{
		{Num: 10, Color: def.CardColorClub, IsMain: true},
		{Num: 9, Color: def.CardColorHeart, IsMain: false},
	}
	for _, card := range cardList {
		player.AddCard(card)
	}
	pool = core.NewCardPool(def.CardColorClub)
	player.SetMainColor(*pool)

	cards = core.NewCards(def.CardColorClub,
		core.Card{Num: 11, Color: def.CardColorClub},
		core.Card{Num: 11, Color: def.CardColorClub})
	shot = core.NewShot(cards, 1, "")
	round = core.NewRound(*shot)
	shotGot, err = player.NextShot(round)
	require.NoError(t, err)
	_, cardsGot = shotGot.Info()
	require.Equal(t, cardList, cardsGot.Cards)
	require.Empty(t, player.cardsByColor[def.CardColorClub])
	require.Empty(t, player.cardsByColor[def.CardColorHeart])

	// 1大+1K
	player = NewAiPlayer("测试", 1)
	cardList = []core.Card{
		{Num: 22, Color: def.CardColorNil, IsMain: true},
		{Num: 15, Color: def.CardColorClub, IsMain: true},
	}
	for _, card := range cardList {
		player.AddCard(card)
	}
	pool = core.NewCardPool(def.CardColorClub)
	player.SetMainColor(*pool)

	cards = core.NewCards(def.CardColorClub,
		core.Card{Num: 11, Color: def.CardColorClub},
		core.Card{Num: 11, Color: def.CardColorClub})
	shot = core.NewShot(cards, 1, "")
	round = core.NewRound(*shot)
	shotGot, err = player.NextShot(round)
	require.NoError(t, err)
	_, cardsGot = shotGot.Info()

	for idx, card := range cardList {
		if card.Num == 22 {
			cardList[idx].Color = player.mainColor
		}
	}
	require.Equal(t, cardList, cardsGot.Cards)
	require.Empty(t, player.cardsByColor[def.CardColorClub])
	require.Empty(t, player.cardsByColor[def.CardColorHeart])
}
