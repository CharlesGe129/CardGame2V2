package test

import (
	"testing"

	"github.com/CharlesGe129/CardGame2V2/pkg/core"
	"github.com/CharlesGe129/CardGame2V2/pkg/def"
	"github.com/stretchr/testify/require"
)

func TestCards_ParseBiggest_Single(t *testing.T) {
	big := core.Card{
		Num:   10,
		Color: def.CardColorClub,
	}
	cards := core.NewCards(def.CardColorClub, big)
	bigGot, typ, num := cards.ParseBiggest()
	require.Equal(t, big, bigGot)
	require.Equal(t, core.CardTypeSingle, typ)
	require.Equal(t, uint8(1), num)
}

func TestCards_ParseBiggest_Double(t *testing.T) {
	big := core.Card{
		Num:   10,
		Color: def.CardColorClub,
	}
	cards := core.NewCards(def.CardColorClub, big, big)
	bigGot, typ, num := cards.ParseBiggest()
	require.Equal(t, big, bigGot)
	require.Equal(t, core.CardTypeDouble, typ)
	require.Equal(t, uint8(1), num)

	small := core.Card{
		Num:   9,
		Color: def.CardColorClub,
	}
	cards = core.NewCards(def.CardColorClub, big, small)
	bigGot, typ, num = cards.ParseBiggest()
	require.Equal(t, big, bigGot)
	require.Equal(t, core.CardTypeSingle, typ)
	require.Equal(t, uint8(1), num)

	cards = core.NewCards(def.CardColorClub, small, big)
	bigGot, typ, num = cards.ParseBiggest()
	require.Equal(t, big, bigGot)
	require.Equal(t, core.CardTypeSingle, typ)
	require.Equal(t, uint8(1), num)
}

func TestCards_ParseBiggest_3Cards(t *testing.T) {
	big := core.Card{
		Num:   10,
		Color: def.CardColorClub,
	}
	next := core.Card{
		Num:   big.Num - 1,
		Color: big.Color,
	}
	cards := core.NewCards(def.CardColorClub, big, big, next)
	bigGot, typ, num := cards.ParseBiggest()
	require.Equal(t, big, bigGot)
	require.Equal(t, core.CardTypeDouble, typ)
	require.Equal(t, uint8(1), num)

	cards = core.NewCards(def.CardColorClub, big, next, next)
	bigGot, typ, num = cards.ParseBiggest()
	require.Equal(t, next, bigGot)
	require.Equal(t, core.CardTypeDouble, typ)
	require.Equal(t, uint8(1), num)

	cards = core.NewCards(def.CardColorClub, big, next, core.Card{
		Num:   next.Num - 1,
		Color: next.Color,
	})
	bigGot, typ, num = cards.ParseBiggest()
	require.Equal(t, big, bigGot)
	require.Equal(t, core.CardTypeSingle, typ)
	require.Equal(t, uint8(1), num)
}

func TestCards_ParseBiggest_ConsecutiveDouble(t *testing.T) {
	// 9988
	big := core.Card{
		Num:   9,
		Color: def.CardColorClub,
	}
	next := core.Card{
		Num:   big.Num - 1,
		Color: big.Color,
	}
	cards := core.NewCards(def.CardColorClub, big, big, next, next)
	bigGot, typ, num := cards.ParseBiggest()
	require.Equal(t, big, bigGot)
	require.Equal(t, core.CardTypeConsecutiveDouble, typ)
	require.Equal(t, uint8(2), num)

	// 998877
	next3 := core.Card{
		Num:   next.Num - 1,
		Color: next.Color,
	}
	cards = core.NewCards(def.CardColorClub, big, big, next, next, next3, next3)
	bigGot, typ, num = cards.ParseBiggest()
	require.Equal(t, big, bigGot)
	require.Equal(t, core.CardTypeConsecutiveDouble, typ)
	require.Equal(t, uint8(3), num)

	// A9988
	cardA := core.Card{
		Num:   14,
		Color: big.Color,
	}
	cards = core.NewCards(def.CardColorClub, big, big, next, next, cardA)
	bigGot, typ, num = cards.ParseBiggest()
	require.Equal(t, big, bigGot)
	require.Equal(t, core.CardTypeConsecutiveDouble, typ)
	require.Equal(t, uint8(2), num)

	// 99887554433
	card7 := core.Card{
		Num:   7,
		Color: big.Color,
	}
	card5 := core.Card{
		Num:   5,
		Color: big.Color,
	}
	card4 := core.Card{
		Num:   4,
		Color: big.Color,
	}
	card3 := core.Card{
		Num:   3,
		Color: big.Color,
	}
	cards = core.NewCards(def.CardColorClub, big, big, next, next, card7, card5, card5, card4, card4, card3, card3)
	bigGot, typ, num = cards.ParseBiggest()
	require.Equal(t, card5, bigGot)
	require.Equal(t, core.CardTypeConsecutiveDouble, typ)
	require.Equal(t, uint8(3), num)
}
