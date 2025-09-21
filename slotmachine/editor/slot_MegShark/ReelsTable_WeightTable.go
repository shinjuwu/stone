package slot_MegShark

import "slotEditor/utils/random"

type (
	RTPs                 string
	Symbol               int
	Reel                 []Symbol //每輪獎圖
	ReelStrips           []Reel   //盤面
	ReelStripsDef        []int
	ReelStripList        map[RTPs]*ReelStrips
	ReelStripLengthTable map[RTPs][]int
	PayTable             [][]int
)

var (
	ngReelStrips = ReelStrips{
		{11, 21, 16, 1, 13, 15, 4, 13, 12, 4, 15, 12, 4, 14, 12, 16, 3, 14, 11, 21, 2, 13, 14, 3, 2, 16, 16, 21, 4, 13, 1, 15, 14, 3, 14, 15, 2, 13, 11, 3, 15, 3, 14, 15, 15, 14},
		{13, 31, 12, 16, 3, 15, 14, 2, 4, 15, 14, 13, 4, 16, 15, 4, 15, 11, 15, 14, 4, 12, 1, 15, 14, 16, 3, 13, 16, 13, 31, 15, 14, 12, 2, 11, 16, 16, 3, 16, 11, 12, 2, 11, 15, 14},
		{14, 21, 16, 14, 2, 13, 4, 15, 15, 31, 11, 12, 21, 14, 13, 2, 1, 16, 14, 15, 4, 16, 12, 21, 3, 13, 14, 1, 13, 11, 3, 16, 11, 31, 4, 15, 15, 2, 14, 12, 14, 11, 4, 15, 16, 11},
		{15, 31, 12, 15, 14, 1, 16, 12, 2, 12, 11, 15, 4, 3, 16, 13, 11, 4, 15, 14, 13, 3, 15, 12, 31, 2, 11, 12, 3, 15, 11, 4, 15, 16, 15, 1, 14, 13, 2, 14, 14, 16, 15, 1, 4, 15},
		{16, 21, 16, 12, 4, 13, 12, 2, 13, 11, 4, 21, 11, 16, 15, 3, 14, 12, 31, 12, 12, 16, 1, 21, 13, 12, 4, 15, 15, 2, 16, 11, 14, 3, 16, 31, 13, 3, 13, 14, 11, 21, 16, 13, 11, 14},
	}
	fgReelStrips = ReelStrips{
		{16, 21, 3, 11, 1, 11, 16, 2, 1, 12, 4, 2, 14, 14, 1, 4, 15, 14, 2, 12, 2, 3, 21, 2, 11, 4, 3, 13, 15, 4, 1, 15, 13, 2, 3, 14, 13, 4, 16},
		{15, 12, 2, 11, 1, 16, 4, 16, 31, 15, 3, 1, 14, 12, 14, 2, 4, 2, 11, 12, 1, 16, 16, 4, 13, 13, 1, 15, 16, 4, 12, 31, 13, 3, 3, 14, 15, 11, 3},
		{12, 21, 14, 2, 12, 4, 16, 15, 11, 1, 21, 4, 16, 13, 4, 13, 15, 14, 4, 31, 15, 11, 2, 15, 1, 16, 13, 3, 21, 14, 16, 11, 1, 12, 15, 3, 16, 14, 31},
		{2, 1, 13, 16, 1, 14, 15, 4, 14, 13, 31, 3, 12, 13, 4, 3, 15, 14, 15, 2, 11, 14, 3, 16, 12, 3, 12, 31, 14, 12, 4, 16, 2, 15, 16, 1, 11, 15, 4},
		{1, 21, 16, 14, 2, 12, 16, 11, 31, 1, 13, 4, 16, 16, 3, 15, 13, 21, 13, 14, 3, 12, 2, 14, 15, 4, 11, 21, 14, 12, 2, 11, 15, 31, 4, 15, 13, 1, 15},
	}

	FG_randomSym = map[string]*WeightGames{
		"98": randomsym_98,
		"97": randomsym_97,
		"92": randomsym_92,
	}
	//RTP:98

	// fgWeightTable_98 = []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	fgWeightTable_98 = []int{160, 150, 150, 150, 150, 150, 150, 150, 355, 360}
	fgWeightTable_97 = []int{160, 155, 155, 150, 150, 160, 150, 150, 330, 330}
	fgWeightTable_92 = []int{240, 250, 200, 190, 180, 200, 180, 180, 170, 170}

	fgObjectTable = []int{1, 2, 3, 4, 11, 12, 13, 14, 15, 16}

	randomsym_98 = WeightNewGames(
		fgWeightTable_98,
		fgObjectTable,
	)
	randomsym_97 = WeightNewGames(
		fgWeightTable_97,
		fgObjectTable,
	)
	randomsym_92 = WeightNewGames(
		fgWeightTable_92,
		fgObjectTable,
	)
)

var (
	NgToFgPos = map[int][]int{ //隨不同遊戲不同
		0: {44, 45, 0, 16, 17, 18, 24, 25, 26},
		2: {44, 45, 0, 9, 10, 11, 20, 21, 22},
		4: {44, 45, 0, 8, 9, 10, 20, 21, 22, 38, 39, 40},
	}

	Rand_SF      = []int{3, 4, 5}
	Rand_SF_weig = []int{93830, 6017, 152}
	Rand_SFcount = WeightNewGames(
		Rand_SF_weig,
		Rand_SF,
	)
)

func GetSFpos(randsf int) []int {
	var pos []int
	pos = []int{1, 0, 1, 0, 1}
	return pos
}
func GetFGpos(pos []int) []int {
	var ptpos []int
	for i, v := range pos {
		if v == 1 {
			id := random.Intn(len(NgToFgPos[i]))
			ptpos = append(ptpos, NgToFgPos[i][id])
		} else {
			id := random.Intn(len(ngReelStrips[i]))
			ptpos = append(ptpos, id)
		}
	}
	return ptpos
}
