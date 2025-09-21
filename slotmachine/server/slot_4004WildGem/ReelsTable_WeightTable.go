package slot_4004WildGem

import (
	"slotserver/utils/random"

	"github.com/shopspring/decimal"
)

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
		{21, 11, 3, 13, 5, 5, 14, 1, 15, 2, 2, 11, 15, 3, 12, 11, 13, 12, 5, 15, 3, 3, 14, 12, 11, 4, 4, 11, 12, 15, 21, 14, 2, 12, 13, 4, 4, 13, 15, 5, 14, 11, 2, 2, 12, 21, 13, 15, 4, 14, 13, 2, 15, 11, 1, 14, 12, 3, 3, 13, 15, 5, 12, 14, 4, 4, 15, 13, 5, 5, 11, 14, 12},
		{21, 5, 11, 2, 2, 12, 15, 13, 3, 3, 14, 1, 11, 15, 14, 12, 2, 2, 11, 15, 14, 21, 12, 13, 3, 3, 11, 15, 4, 4, 2, 11, 15, 14, 13, 5, 5, 13, 11, 12, 5, 14, 14, 3, 12, 14, 15, 13, 1, 13, 5, 15, 12, 4, 4, 11, 13, 3, 15, 11, 21, 12, 2, 2, 3, 14, 12, 4, 4, 11, 13, 5, 14},
		{21, 15, 5, 5, 11, 1, 12, 14, 31, 12, 13, 2, 2, 14, 11, 3, 31, 11, 15, 12, 4, 4, 15, 21, 13, 3, 14, 31, 12, 2, 12, 4, 13, 12, 1, 15, 13, 31, 15, 14, 11, 3, 13, 15, 21, 14, 12, 4, 11, 14, 13, 31, 13, 14, 5, 5, 15, 12, 3, 3, 11, 15, 13, 11, 5, 15, 31, 2, 2, 11, 12, 14, 13},
		{21, 5, 11, 12, 31, 13, 12, 14, 2, 2, 11, 15, 11, 1, 15, 15, 31, 14, 13, 12, 2, 15, 14, 15, 14, 3, 11, 12, 13, 15, 3, 14, 5, 4, 11, 31, 13, 12, 2, 13, 14, 11, 1, 13, 31, 3, 3, 15, 14, 4, 13, 12, 15, 21, 5, 5, 11, 12, 31, 14, 5, 13, 15, 4, 12, 13, 3, 11, 14, 12, 4, 4, 11},
		{21, 4, 31, 11, 13, 1, 14, 15, 12, 3, 15, 5, 13, 15, 14, 5, 15, 14, 31, 12, 11, 21, 13, 3, 11, 13, 4, 12, 15, 11, 11, 31, 2, 13, 4, 12, 1, 14, 3, 3, 11, 31, 12, 3, 5, 4, 14, 21, 13, 15, 2, 12, 13, 11, 31, 15, 15, 12, 4, 4, 11, 15, 2, 5, 5, 14, 31, 14, 12, 2, 2, 11, 13},
	}

	fgReelStrips = ReelStrips{
		{21, 11, 12, 2, 2, 13, 14, 4, 11, 15, 5, 14, 5, 13, 1, 14, 11, 5, 14, 5, 12, 3, 15, 13, 2, 2, 14, 11, 3, 3, 15, 13, 4, 12, 15, 21, 14, 13, 5, 3, 15, 13, 4, 4, 12, 11, 2, 14, 15, 4, 11, 13, 1, 12, 15, 4, 5, 13, 11, 3, 3, 15, 12, 2, 4, 12, 14, 21, 11, 12, 5, 3},
		{21, 14, 5, 5, 13, 15, 2, 3, 11, 14, 4, 15, 4, 13, 5, 3, 11, 15, 4, 14, 12, 3, 3, 13, 12, 2, 13, 14, 1, 11, 21, 2, 2, 11, 13, 5, 4, 12, 11, 4, 4, 13, 14, 3, 12, 15, 5, 5, 11, 13, 2, 14, 21, 3, 3, 15, 14, 1, 12, 11, 2, 2, 15, 11, 5, 13, 12, 4, 4, 14, 15, 11, 3, 2},
		{21, 15, 12, 31, 31, 4, 15, 2, 15, 11, 31, 15, 13, 3, 12, 14, 31, 31, 11, 4, 13, 31, 5, 4, 13, 15, 2, 11, 31, 31, 31, 13, 11, 3, 5, 13, 1, 11, 31, 31, 12, 21, 11, 2, 4, 13, 14, 5, 12, 31, 31, 31, 14, 13, 2, 14, 21, 3, 12, 14, 31, 31, 31, 11, 15, 1, 14, 12, 31, 31, 15, 14, 31},
		{21, 12, 31, 31, 13, 4, 11, 31, 31, 31, 13, 12, 5, 21, 12, 2, 31, 11, 1, 12, 31, 31, 14, 13, 31, 11, 14, 2, 4, 11, 21, 12, 15, 4, 31, 31, 5, 15, 13, 3, 14, 5, 31, 31, 11, 14, 2, 15, 14, 31, 11, 12, 1, 11, 13, 31, 31, 15, 14, 4, 15, 31, 5, 12, 13, 31, 31, 31, 15, 3, 12, 14, 0},
		{21, 13, 15, 2, 14, 11, 31, 31, 13, 12, 4, 15, 13, 31, 31, 3, 11, 21, 5, 14, 31, 31, 31, 13, 1, 12, 31, 2, 11, 15, 31, 31, 13, 11, 4, 31, 12, 5, 14, 11, 2, 15, 31, 31, 14, 13, 3, 15, 1, 12, 31, 3, 21, 15, 11, 4, 14, 13, 5, 31, 31, 31, 12, 13, 4, 12, 31, 5, 15, 14, 3, 31, 31},
	}
)

func GetNGReelsStrips() ReelStrips {
	return ngReelStrips
}

func GetFGReelsStrips() ReelStrips {
	return fgReelStrips
}

// NormalGame Wild Weight
var (
	WWMultiObjectTable_R3 = []int{32, 33, 34, 35, 36, 37, 38}
	WWMultiObjectTable_R4 = []int{38, 39, 40, 41, 42, 43, 44}
	WWMultiObjectTable_R5 = []int{45, 46, 47, 48, 49}
	//NGwild
	//	@input_1	rtp
	//	@input_2	ReelNum
	NGwild = map[string]map[int]*WeightGames{
		"98": WWMulti_ng_98,
		"97": WWMulti_ng_97,
		"92": WWMulti_ng_92,
	}

	//92
	WWMultiWeightTable_R3_ng_92 = []int{90, 80, 50, 30, 15, 8, 1}
	WWMultiWeightTable_R4_ng_92 = []int{100, 90, 50, 25, 20, 9, 1}
	WWMultiWeightTable_R5_ng_92 = []int{280, 80, 20, 6, 1}

	WWMulti_R3_ng_92 = WeightNewGames(
		WWMultiWeightTable_R3_ng_92,
		WWMultiObjectTable_R3,
	)
	WWMulti_R4_ng_92 = WeightNewGames(
		WWMultiWeightTable_R4_ng_92,
		WWMultiObjectTable_R4,
	)
	WWMulti_R5_ng_92 = WeightNewGames(
		WWMultiWeightTable_R5_ng_92,
		WWMultiObjectTable_R5,
	)
	//input:reel number
	WWMulti_ng_92 = map[int]*WeightGames{
		2: WWMulti_R3_ng_92,
		3: WWMulti_R4_ng_92,
		4: WWMulti_R5_ng_92,
	}

	// 97
	WWMultiWeightTable_R3_ng_97 = []int{10, 70, 50, 30, 30, 20, 5} //{0, 0, 0, 1, 0, 0, 0} //
	WWMultiWeightTable_R4_ng_97 = []int{5, 80, 40, 40, 40, 8, 5}   //{0, 0, 1, 0, 0, 0, 0} //
	WWMultiWeightTable_R5_ng_97 = []int{200, 100, 20, 6, 1}        //{1, 0, 0, 0, 0}       //

	WWMulti_R3_ng_97 = WeightNewGames(
		WWMultiWeightTable_R3_ng_97,
		WWMultiObjectTable_R3,
	)
	WWMulti_R4_ng_97 = WeightNewGames(
		WWMultiWeightTable_R4_ng_97,
		WWMultiObjectTable_R4,
	)
	WWMulti_R5_ng_97 = WeightNewGames(
		WWMultiWeightTable_R5_ng_97,
		WWMultiObjectTable_R5,
	)

	//	@input	- Reel number
	WWMulti_ng_97 = map[int]*WeightGames{
		2: WWMulti_R3_ng_97,
		3: WWMulti_R4_ng_97,
		4: WWMulti_R5_ng_97,
	}
	// 98
	WWMultiWeightTable_R3_ng_98 = []int{10, 60, 50, 40, 30, 20, 5} //{0, 0, 0, 1, 0, 0, 0} //
	WWMultiWeightTable_R4_ng_98 = []int{5, 40, 40, 49, 49, 8, 5}   //{0, 0, 1, 0, 0, 0, 0} //
	WWMultiWeightTable_R5_ng_98 = []int{200, 120, 20, 6, 1}        //{1, 0, 0, 0, 0}       //

	WWMulti_R3_ng_98 = WeightNewGames(
		WWMultiWeightTable_R3_ng_98,
		WWMultiObjectTable_R3,
	)
	WWMulti_R4_ng_98 = WeightNewGames(
		WWMultiWeightTable_R4_ng_98,
		WWMultiObjectTable_R4,
	)
	WWMulti_R5_ng_98 = WeightNewGames(
		WWMultiWeightTable_R5_ng_98,
		WWMultiObjectTable_R5,
	)

	//	@param	reel number
	WWMulti_ng_98 = map[int]*WeightGames{
		2: WWMulti_R3_ng_98,
		3: WWMulti_R4_ng_98,
		4: WWMulti_R5_ng_98,
	}
)

// FreeGame Wild Weight
var (
	//input:rtp
	FGwild = map[string]map[int]*WeightGames{
		"98": WWMulti_fg,
		"97": WWMulti_fg,
		"92": WWMulti_fg,
	}
	//WWMulti_fg
	//	@input	reel number
	WWMulti_fg = map[int]*WeightGames{
		2: WWMulti_R3_fg,
		3: WWMulti_R4_fg,
		4: WWMulti_R5_fg,
	}
	WWMultiWeightTable_R3_fg = []int{10, 70, 200, 200, 100, 40, 10} //{0, 0, 0, 1, 0, 0, 0} //
	WWMultiWeightTable_R4_fg = []int{50, 60, 200, 200, 200, 30, 10} //{0, 0, 1, 0, 0, 0, 0} //
	WWMultiWeightTable_R5_fg = []int{400, 500, 350, 110, 30}        //{1, 0, 0, 0, 0}       //

	WWMulti_R3_fg = WeightNewGames(
		WWMultiWeightTable_R3_fg,
		WWMultiObjectTable_R3,
	)
	WWMulti_R4_fg = WeightNewGames(
		WWMultiWeightTable_R4_fg,
		WWMultiObjectTable_R4,
	)
	WWMulti_R5_fg = WeightNewGames(
		WWMultiWeightTable_R5_fg,
		WWMultiObjectTable_R5,
	)

	//Fgtimes
	//	@input	SFcount
	//	@output	FGtimes
	FGtimes = map[int]int{
		3: 7,
		4: 10,
		5: 15,
	}
)

// Wild Multi. Const
var (
	WW_number = map[int]decimal.Decimal{
		32: decimal.NewFromFloat(0.3),
		33: decimal.NewFromFloat(0.5),
		34: decimal.NewFromFloat(0.8),
		35: decimal.NewFromFloat(1),
		36: decimal.NewFromFloat(1.5),
		37: decimal.NewFromFloat(2),
		38: decimal.NewFromFloat(3),
		39: decimal.NewFromFloat(4),
		40: decimal.NewFromFloat(5),
		41: decimal.NewFromFloat(6),
		42: decimal.NewFromFloat(7),
		43: decimal.NewFromFloat(8),
		44: decimal.NewFromFloat(9),
		45: decimal.NewFromFloat(10),  //Mini
		46: decimal.NewFromFloat(20),  //Minor
		47: decimal.NewFromFloat(30),  //Major
		48: decimal.NewFromFloat(50),  //Super
		49: decimal.NewFromFloat(100), //Grand
	}
)

var (
	NgToFgPos = map[int][]int{ //隨不同遊戲不同
		0: {27, 28, 29, 42, 43, 44, 70, 71, 72},
		1: {18, 19, 20, 57, 58, 59, 70, 71, 72},
		2: {20, 21, 22, 41, 42, 43, 70, 71, 72},
		3: {50, 51, 52, 70, 71, 72},
		4: {18, 19, 20, 44, 45, 46, 70, 71, 72},
	}

	Rand_SF      = []int{3, 4, 5}
	Rand_SF_weig = []int{90000, 6000, 160}
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
