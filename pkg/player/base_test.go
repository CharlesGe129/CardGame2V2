package player

import (
	"testing"

	"github.com/CharlesGe129/CardGame2V2/pkg/core"
	"github.com/CharlesGe129/CardGame2V2/pkg/def"
	"github.com/stretchr/testify/require"
)

func TestBasePlayer_GetCards_Double(t *testing.T) {
	// no doubles
	player := NewAiPlayer("测试", 1)
	curCardList := []core.Card{
		{Num: 10, Color: def.CardColorClub, IsMain: false},
		{Num: 9, Color: def.CardColorClub, IsMain: false},
		{Num: 8, Color: def.CardColorClub, IsMain: false},
	}

	cards := core.NewCards(def.CardColorHeart,
		core.Card{Num: 11, Color: def.CardColorClub},
		core.Card{Num: 11, Color: def.CardColorClub})
	bigCards, cardType, _ := cards.ParseBiggest()
	cardList := player.GetCards(curCardList, bigCards, cardType)
	require.Empty(t, cardList)

	// has 1 doubles
	player = NewAiPlayer("测试", 1)
	curCardList = []core.Card{
		{Num: 10, Color: def.CardColorClub, IsMain: false},
		{Num: 9, Color: def.CardColorClub, IsMain: false},
		{Num: 9, Color: def.CardColorClub, IsMain: false},
		{Num: 8, Color: def.CardColorClub, IsMain: false},
	}

	cards = core.NewCards(def.CardColorHeart,
		core.Card{Num: 11, Color: def.CardColorClub},
		core.Card{Num: 11, Color: def.CardColorClub})
	bigCards, cardType, _ = cards.ParseBiggest()
	cardList = player.GetCards(curCardList, bigCards, cardType)
	require.Equal(t, []core.Card{
		{Num: 9, Color: def.CardColorClub, IsMain: false},
		{Num: 9, Color: def.CardColorClub, IsMain: false},
	}, cardList)

	// has big doubles
	player = NewAiPlayer("测试", 1)
	curCardList = []core.Card{
		{Num: 14, Color: def.CardColorClub, IsMain: false},
		{Num: 13, Color: def.CardColorClub, IsMain: false},
		{Num: 13, Color: def.CardColorClub, IsMain: false},
		{Num: 9, Color: def.CardColorClub, IsMain: false},
		{Num: 9, Color: def.CardColorClub, IsMain: false},
		{Num: 8, Color: def.CardColorClub, IsMain: false},
	}

	cards = core.NewCards(def.CardColorHeart,
		core.Card{Num: 11, Color: def.CardColorClub},
		core.Card{Num: 11, Color: def.CardColorClub})
	bigCards, cardType, _ = cards.ParseBiggest()
	cardList = player.GetCards(curCardList, bigCards, cardType)
	require.Equal(t, []core.Card{
		{Num: 13, Color: def.CardColorClub, IsMain: false},
		{Num: 13, Color: def.CardColorClub, IsMain: false},
	}, cardList)

	// has small doubles
	player = NewAiPlayer("测试", 1)
	curCardList = []core.Card{
		{Num: 14, Color: def.CardColorClub, IsMain: false},
		{Num: 9, Color: def.CardColorClub, IsMain: false},
		{Num: 9, Color: def.CardColorClub, IsMain: false},
		{Num: 8, Color: def.CardColorClub, IsMain: false},
		{Num: 3, Color: def.CardColorClub, IsMain: false},
		{Num: 3, Color: def.CardColorClub, IsMain: false},
	}

	cards = core.NewCards(def.CardColorHeart,
		core.Card{Num: 11, Color: def.CardColorClub},
		core.Card{Num: 11, Color: def.CardColorClub})
	bigCards, cardType, _ = cards.ParseBiggest()
	cardList = player.GetCards(curCardList, bigCards, cardType)
	require.Equal(t, []core.Card{
		{Num: 3, Color: def.CardColorClub, IsMain: false},
		{Num: 3, Color: def.CardColorClub, IsMain: false},
	}, cardList)
}
