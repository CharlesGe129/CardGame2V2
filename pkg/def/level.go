package def

var (
	MapLevelToName = map[uint8]string{
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
	}
	MapNameToLevel map[string]uint8
)

func init() {
	MapNameToLevel = make(map[string]uint8)
	for k, v := range MapLevelToName {
		MapNameToLevel[v] = k
	}
}
