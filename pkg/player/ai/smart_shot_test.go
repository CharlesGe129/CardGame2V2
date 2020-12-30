package ai

import (
	"testing"

	"github.com/CharlesGe129/CardGame2V2/pkg/core"
	"github.com/CharlesGe129/CardGame2V2/pkg/def"
	"github.com/stretchr/testify/require"
)

func TestSmartShot_NextShot_Single_HasCard(t *testing.T) {
	// bigger
	mainColor := def.CardColorClub
	player := NewSmartShot(mainColor, map[def.CardColor][]core.Card{
		mainColor: {
			{Num: 10, Color: mainColor, IsMain: true},
			{Num: 5, Color: mainColor, IsMain: true},
		},
	})
	shot := core.NewShot(core.NewCards(mainColor, []core.Card{
		{Num: 9, Color: mainColor, IsMain: true},
	}...), 1, "test")
	cardList, err := player.NextShot(shot)
	require.NoError(t, err)

	expectCardList := []core.Card{
		{Num: 10, Color: mainColor, IsMain: true},
	}
	require.Equal(t, expectCardList, cardList)

	// bigger
	player = NewSmartShot(mainColor, map[def.CardColor][]core.Card{
		mainColor: {
			{Num: 22, Color: mainColor, IsMain: true},
			{Num: 5, Color: mainColor, IsMain: true},
		},
	})
	shot = core.NewShot(core.NewCards(mainColor, []core.Card{
		{Num: 9, Color: mainColor, IsMain: true},
	}...), 1, "test")
	cardList, err = player.NextShot(shot)
	require.NoError(t, err)

	expectCardList = []core.Card{
		{Num: 22, Color: mainColor, IsMain: true},
	}
	require.Equal(t, expectCardList, cardList)

	// smaller
	player = NewSmartShot(mainColor, map[def.CardColor][]core.Card{
		mainColor: {
			{Num: 10, Color: mainColor, IsMain: true},
			{Num: 5, Color: mainColor, IsMain: true},
			{Num: 12, Color: mainColor, IsMain: true},
		},
	})
	shot = core.NewShot(core.NewCards(mainColor, []core.Card{
		{Num: 15, Color: mainColor, IsMain: true},
	}...), 1, "test")
	cardList, err = player.NextShot(shot)
	require.NoError(t, err)

	expectCardList = []core.Card{
		{Num: 12, Color: mainColor, IsMain: true},
	}
	require.Equal(t, expectCardList, cardList)
}

func TestSmartShot_NextShot_Single_NoCard(t *testing.T) {
	// has main
	mainColor := def.CardColorClub
	player := NewSmartShot(mainColor, map[def.CardColor][]core.Card{
		mainColor: {
			{Num: 10, Color: mainColor, IsMain: true},
			{Num: 5, Color: mainColor, IsMain: true},
		},
	})
	shot := core.NewShot(core.NewCards(mainColor, []core.Card{
		{Num: 9, Color: def.CardColorHeart, IsMain: false},
	}...), 1, "test")
	cardList, err := player.NextShot(shot)
	require.NoError(t, err)

	expectCardList := []core.Card{
		{Num: 5, Color: mainColor, IsMain: true},
	}
	require.Equal(t, expectCardList, cardList)

	// no main
	player = NewSmartShot(mainColor, map[def.CardColor][]core.Card{
		def.CardColorDiamond: {
			{Num: 12, Color: mainColor, IsMain: true},
			{Num: 2, Color: mainColor, IsMain: true},
		},
	})
	shot = core.NewShot(core.NewCards(mainColor, []core.Card{
		{Num: 9, Color: def.CardColorHeart, IsMain: false},
	}...), 1, "test")
	cardList, err = player.NextShot(shot)
	require.NoError(t, err)

	expectCardList = []core.Card{
		{Num: 2, Color: mainColor, IsMain: true},
	}
	require.Equal(t, expectCardList, cardList)
}

func TestSmartShot_NextShot_Double_HasCard(t *testing.T) {
	// bigger
	mainColor := def.CardColorClub
	player := NewSmartShot(mainColor, map[def.CardColor][]core.Card{
		mainColor: {
			{Num: 10, Color: mainColor, IsMain: true},
			{Num: 10, Color: mainColor, IsMain: true},
			{Num: 5, Color: mainColor, IsMain: true},
			{Num: 5, Color: mainColor, IsMain: true},
			{Num: 3, Color: mainColor, IsMain: true},
		},
	})
	shot := core.NewShot(core.NewCards(mainColor, []core.Card{
		{Num: 9, Color: mainColor, IsMain: true},
		{Num: 9, Color: mainColor, IsMain: true},
		{Num: 7, Color: mainColor, IsMain: true},
	}...), 1, "test")
	cardList, err := player.NextShot(shot)
	require.NoError(t, err)

	expectCardList := []core.Card{
		{Num: 10, Color: mainColor, IsMain: true},
		{Num: 10, Color: mainColor, IsMain: true},
		{Num: 3, Color: mainColor, IsMain: true},
	}
	require.Equal(t, expectCardList, cardList)

	// smaller
	player = NewSmartShot(mainColor, map[def.CardColor][]core.Card{
		mainColor: {
			{Num: 12, Color: mainColor, IsMain: true},
			{Num: 12, Color: mainColor, IsMain: true},
			{Num: 3, Color: mainColor, IsMain: true},
			{Num: 3, Color: mainColor, IsMain: true},
			{Num: 2, Color: mainColor, IsMain: true},
		},
	})
	shot = core.NewShot(core.NewCards(mainColor, []core.Card{
		{Num: 14, Color: mainColor, IsMain: true},
		{Num: 14, Color: mainColor, IsMain: true},
	}...), 1, "test")
	cardList, err = player.NextShot(shot)
	require.NoError(t, err)

	expectCardList = []core.Card{
		{Num: 3, Color: mainColor, IsMain: true},
		{Num: 3, Color: mainColor, IsMain: true},
	}
	require.Equal(t, expectCardList, cardList)
}

func TestSmartShot_NextShot_Double_NoCard(t *testing.T) {
	// has main
	mainColor := def.CardColorClub
	player := NewSmartShot(mainColor, map[def.CardColor][]core.Card{
		mainColor: {
			{Num: 10, Color: mainColor, IsMain: true},
			{Num: 10, Color: mainColor, IsMain: true},
			{Num: 5, Color: mainColor, IsMain: true},
			{Num: 5, Color: mainColor, IsMain: true},
		},
	})
	shot := core.NewShot(core.NewCards(mainColor, []core.Card{
		{Num: 9, Color: def.CardColorHeart, IsMain: false},
		{Num: 9, Color: def.CardColorHeart, IsMain: false},
	}...), 1, "test")
	cardList, err := player.NextShot(shot)
	require.NoError(t, err)

	expectCardList := []core.Card{
		{Num: 5, Color: mainColor, IsMain: true},
		{Num: 5, Color: mainColor, IsMain: true},
	}
	require.Equal(t, expectCardList, cardList)

	// has main, no double
	player = NewSmartShot(mainColor, map[def.CardColor][]core.Card{
		mainColor: {
			{Num: 12, Color: mainColor, IsMain: true},
			{Num: 3, Color: mainColor, IsMain: true},
			{Num: 2, Color: mainColor, IsMain: true},
		},
		def.CardColorDiamond: {
			{Num: 4, Color: def.CardColorDiamond, IsMain: false},
			{Num: 6, Color: def.CardColorDiamond, IsMain: false},
		},
	})
	shot = core.NewShot(core.NewCards(mainColor, []core.Card{
		{Num: 9, Color: def.CardColorHeart, IsMain: false},
		{Num: 9, Color: def.CardColorHeart, IsMain: false},
	}...), 1, "test")
	cardList, err = player.NextShot(shot)
	require.NoError(t, err)

	expectCardList = []core.Card{
		{Num: 4, Color: def.CardColorDiamond, IsMain: false},
		{Num: 6, Color: def.CardColorDiamond, IsMain: false},
	}
	require.Equal(t, expectCardList, cardList)

	// no main
	player = NewSmartShot(mainColor, map[def.CardColor][]core.Card{
		def.CardColorDiamond: {
			{Num: 12, Color: mainColor, IsMain: false},
			{Num: 12, Color: mainColor, IsMain: false},
			{Num: 3, Color: mainColor, IsMain: false},
			{Num: 2, Color: mainColor, IsMain: false},
		},
	})
	shot = core.NewShot(core.NewCards(mainColor, []core.Card{
		{Num: 9, Color: def.CardColorHeart, IsMain: false},
		{Num: 9, Color: def.CardColorHeart, IsMain: false},
	}...), 1, "test")
	cardList, err = player.NextShot(shot)
	require.NoError(t, err)

	expectCardList = []core.Card{
		{Num: 2, Color: mainColor, IsMain: false},
		{Num: 3, Color: mainColor, IsMain: false},
	}
	require.Equal(t, expectCardList, cardList)
}
