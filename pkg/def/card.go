package def

const (
	CardColorNil     CardColor = ""
	CardColorSpade   CardColor = "spade"
	CardColorHeart   CardColor = "heart"
	CardColorClub    CardColor = "club"
	CardColorDiamond CardColor = "diamond"
)

var (
	MapCardName = map[uint8]string{
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
		15: "2",
		21: "小",
		22: "大",
	}
	MapCardColor = map[string]CardColor{
		"hei":  CardColorSpade,
		"hong": CardColorHeart,
		"cao":  CardColorClub,
		"fang": CardColorDiamond,
	}

	MapNameToCard  map[string]uint8
	MapColorToCard map[CardColor]string
)

type CardColor string

func init() {
	MapNameToCard = make(map[string]uint8)
	for k, v := range MapCardName {
		MapNameToCard[v] = k
	}
	MapColorToCard = make(map[CardColor]string)
	for k, v := range MapCardColor {
		MapColorToCard[v] = k
	}
}
