package def

const (
	CardColorNil     CardColor = ""
	CardColorSpade   CardColor = "spade"
	CardColorHeart   CardColor = "heart"
	CardColorClub    CardColor = "club"
	CardColorDiamond CardColor = "diamond"
)

var (
	mapCardName = map[uint8]string{
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "0",
		11: "J",
		12: "Q",
		13: "K",
		14: "A",
		// 15: "main card, should be set by game.NewGame()"
		21: "小",
		22: "大",
	}
	MapCardName  map[uint8]string
	MapCardColor = map[string]CardColor{
		"hei":  CardColorSpade,
		"hong": CardColorHeart,
		"cao":  CardColorClub,
		"fang": CardColorDiamond,
	}
	MapColorZnCh = map[CardColor]string{
		CardColorSpade:   "黑桃",
		CardColorHeart:   "红桃",
		CardColorClub:    "草花",
		CardColorDiamond: "方块",
	}

	MapNameToCard  map[string]uint8
	MapColorToCard map[CardColor]string
)

type CardColor string

func Init(mainCard string) (origMainNum uint8) {
	MapCardName = make(map[uint8]string)
	for k, v := range mapCardName {
		MapCardName[k] = v
	}
	for num, name := range MapCardName {
		if name == mainCard {
			origMainNum = num
			delete(MapCardName, num)
			break
		}
	}
	MapCardName[15] = mainCard

	MapNameToCard = make(map[string]uint8)
	for k, v := range MapCardName {
		MapNameToCard[v] = k
	}
	MapColorToCard = make(map[CardColor]string)
	for k, v := range MapCardColor {
		MapColorToCard[v] = k
	}
	return
}
