package slot_Midas

import "slotserver/utils/random"

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
		{12, 21, 15, 14, 3, 13, 5, 13, 12, 15, 12, 12, 11, 5, 15, 12, 14, 1, 11, 13, 14, 13, 14, 4, 15, 13, 11, 13, 11, 2, 11, 14, 13, 11},
		{14, 21, 13, 3, 31, 14, 14, 4, 15, 14, 12, 2, 14, 12, 4, 13, 13, 11, 15, 31, 15, 5, 13, 12, 13, 12, 13, 14, 5, 13, 11, 1, 14, 5},
		{15, 21, 13, 14, 11, 13, 13, 14, 1, 14, 12, 14, 31, 13, 11, 4, 13, 14, 11, 3, 14, 2, 12, 15, 15, 3, 15, 11, 5, 12, 31, 12, 11, 5},
		{2, 21, 14, 5, 13, 12, 1, 13, 15, 31, 13, 14, 11, 5, 14, 12, 14, 5, 15, 13, 2, 13, 15, 12, 31, 13, 14, 11, 3, 14, 13, 4, 11, 3},
		{12, 21, 13, 4, 14, 2, 15, 11, 15, 5, 13, 11, 12, 4, 13, 14, 11, 3, 11, 12, 21, 13, 14, 5, 11, 2, 14, 12, 1, 14, 13, 3, 12, 1},
	}

	fgReelStrips = ReelStrips{
		{5, 15, 12, 2, 14, 13, 3, 15, 14, 12, 4, 11, 12, 13, 1, 15, 12, 5, 13, 13, 14, 5, 12, 13, 14, 4, 11, 15, 2, 14, 12, 4, 11, 13, 12, 14, 1, 11, 15, 3, 11, 15, 4, 14, 13, 5, 11},
		{2, 31, 12, 5, 13, 11, 13, 15, 3, 11, 14, 4, 13, 13, 5, 15, 14, 11, 12, 3, 14, 5, 15, 12, 1, 14, 11, 5, 12, 14, 13, 11, 3, 14, 12, 15, 4, 14, 13, 1, 15, 14, 4, 12, 15, 2, 12},
		{2, 12, 15, 13, 5, 14, 11, 1, 31, 4, 12, 13, 11, 2, 15, 14, 3, 15, 11, 4, 14, 12, 1, 13, 13, 14, 3, 11, 15, 12, 2, 15, 14, 11, 5, 12, 11, 5, 14, 13, 11, 3, 15, 12, 14, 3, 13},
		{15, 1, 11, 12, 3, 2, 13, 11, 3, 4, 15, 14, 1, 12, 15, 14, 31, 13, 14, 15, 2, 13, 11, 5, 12, 5, 11, 12, 31, 1, 13, 14, 3, 12, 15, 13, 4, 14, 13, 11, 5, 12, 15, 13, 14, 14, 2},
		{11, 14, 13, 3, 11, 14, 15, 12, 2, 12, 13, 15, 14, 12, 1, 11, 13, 2, 3, 12, 15, 4, 11, 15, 13, 14, 4, 12, 3, 15, 5, 11, 13, 14, 2, 12, 13, 14, 12, 15, 11, 1, 14, 13, 5, 11, 4},
	}

	//ng_def
	//ng_WW倍數權重
	WWMulti = map[string]*WeightGames{
		"98": WWMulti_all,
		"97": WWMulti_all,
		"92": WWMulti_all,
	}
	WWMultiWeightTable_all = []int{150, 10, 1}

	WWMultiObjectTable = []int{1, 2, 3}

	WWMulti_all = WeightNewGames(
		WWMultiWeightTable_all,
		WWMultiObjectTable,
	)

	//fg_def
	//fg_抽取轉數權重
	fgTimes = map[string]*WeightGames{
		"98": fgTimes_98,
		"97": fgTimes_97,
		"92": fgTimes_92,
	}

	fgTimesWeightTable_98 = []int{140, 23, 3}
	fgTimesWeightTable_97 = []int{150, 21, 3}
	fgTimesWeightTable_92 = []int{190, 14, 2}
	fgTimesObjectTable    = []int{1, 2, 3}

	fgTimes_98 = WeightNewGames(
		fgTimesWeightTable_98,
		fgTimesObjectTable,
	)
	fgTimes_97 = WeightNewGames(
		fgTimesWeightTable_97,
		fgTimesObjectTable,
	)
	fgTimes_92 = WeightNewGames(
		fgTimesWeightTable_92,
		fgTimesObjectTable,
	)
	//fg_WW倍數權重
	WWMulti_fg = map[string]*WeightGames{
		"98": WWMulti_fg_all,
		"97": WWMulti_fg_all,
		"92": WWMulti_fg_all,
	}
	WWMultiWeightTable_fg = []int{150, 10, 1}

	WWMulti_fg_all = WeightNewGames(
		WWMultiWeightTable_fg,
		WWMultiObjectTable,
	)
)

func GetNGReelsStrips() ReelStrips {
	return ngReelStrips
}

func GetFGReelsStrips() ReelStrips {
	return fgReelStrips
}

var (
	NgToFgPos = map[int][]int{ //隨不同遊戲不同
		0: {0, 32, 33},
		1: {0, 32, 33},
		2: {0, 32, 33},
		3: {0, 32, 33},
		4: {0, 17, 18, 19, 32, 33},
	}

	Rand_SF      = []int{3, 4, 5}
	Rand_SF_weig = []int{94000, 6100, 150}
	Rand_SFcount = WeightNewGames(
		Rand_SF_weig,
		Rand_SF,
	)
	SF_3 = map[int][]int{
		1: {1, 1, 1, 0, 0},
		2: {1, 1, 0, 1, 0},
		3: {1, 1, 0, 0, 1},
		4: {1, 0, 1, 1, 0},
		5: {1, 0, 1, 0, 1},
		6: {1, 0, 0, 1, 1},
		7: {0, 1, 1, 1, 0},
		8: {0, 1, 1, 0, 1},
		9: {0, 1, 0, 1, 1},
		0: {0, 0, 1, 1, 1},
	}
	SF_4 = map[int][]int{
		1: {1, 1, 1, 1, 0},
		2: {1, 1, 1, 0, 1},
		3: {1, 1, 0, 1, 1},
		4: {1, 0, 1, 1, 1},
		0: {0, 1, 1, 1, 1},
	}
)

func GetSFpos(randsf int) []int {
	var pos []int
	switch randsf {
	case 3:
		id := random.Intn(10)
		pos = SF_3[id]
	case 4:
		id := random.Intn(5)
		pos = SF_4[id]
	case 5:
		pos = []int{1, 1, 1, 1, 1}
	}
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
