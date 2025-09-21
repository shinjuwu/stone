package slot_4006PryTreasure

import (
	"github.com/shopspring/decimal"
)

const (
	bg_1 = Symbol(100) + iota
	bg_2
	bg_3
	bg_4
	bg_5
	bg_7o5
	bg_10
	bg_8
	bg_12
	bg_15
	bg_20
	bg_X2
	bg_X3
	bg_X5
)
const (
	bg_GG = Symbol(300) + iota
	bg_up
)

var (
	BgPayTable = map[int]decimal.Decimal{
		100: decimal.NewFromFloat(1),
		101: decimal.NewFromFloat(2),
		102: decimal.NewFromFloat(3),
		103: decimal.NewFromFloat(4),
		104: decimal.NewFromFloat(5),
		105: decimal.NewFromFloat(7.5),
		106: decimal.NewFromFloat(10),
		107: decimal.NewFromFloat(8),
		108: decimal.NewFromFloat(12),
		109: decimal.NewFromFloat(15),
		110: decimal.NewFromFloat(20),
		111: decimal.NewFromInt(2),
		112: decimal.NewFromInt(3),
		113: decimal.NewFromInt(5),
	}

	//每層中獎球數權重
	bgNumObjectTable    = []int{0, 1, 2, 3, 4, 5, 6}
	bgNumWeightTable_92 = map[int][]int{
		1: {0, 400, 1000, 400, 300, 150, 20},
		2: {599, 804, 300, 200, 80, 0, 0},
		3: {170, 600, 350, 120, 0, 0, 0},
		4: {200, 670, 405, 0, 0, 0, 0},
		5: {0, 1, 0, 0, 0, 0, 0},
	}

	bgNumWeightTable_97 = map[int][]int{
		1: {0, 900, 900, 600, 300, 150, 20},
		2: {140, 524, 400, 250, 120, 0, 0},
		3: {150, 600, 350, 230, 0, 0, 0},
		4: {170, 500, 402, 0, 0, 0, 0},
		5: {0, 1, 0, 0, 0, 0, 0},
	}

	bgNumWeightTable_98 = map[int][]int{
		1: {0, 390, 800, 600, 300, 150, 20},
		2: {140, 600, 400, 250, 120, 0, 0},
		3: {100, 600, 350, 230, 0, 0, 0},
		4: {200, 500, 402, 0, 0, 0, 0},
		5: {0, 1, 0, 0, 0, 0, 0},
	}

	BgBallNumber = map[string]map[int]*WeightGames{
		"92": BgBallNumber_92,
		"97": BgBallNumber_97,
		"98": BgBallNumber_98,
	}
	BgBallNumber_92 = map[int]*WeightGames{
		1: WeightNewGames(
			bgNumWeightTable_92[1],
			bgNumObjectTable,
		),
		2: WeightNewGames(
			bgNumWeightTable_92[2],
			bgNumObjectTable,
		),
		3: WeightNewGames(
			bgNumWeightTable_92[3],
			bgNumObjectTable,
		),
		4: WeightNewGames(
			bgNumWeightTable_92[4],
			bgNumObjectTable,
		),
		5: WeightNewGames(
			bgNumWeightTable_92[5],
			bgNumObjectTable,
		),
	}

	BgBallNumber_97 = map[int]*WeightGames{
		1: WeightNewGames(
			bgNumWeightTable_97[1],
			bgNumObjectTable,
		),
		2: WeightNewGames(
			bgNumWeightTable_97[2],
			bgNumObjectTable,
		),
		3: WeightNewGames(
			bgNumWeightTable_97[3],
			bgNumObjectTable,
		),
		4: WeightNewGames(
			bgNumWeightTable_97[4],
			bgNumObjectTable,
		),
		5: WeightNewGames(
			bgNumWeightTable_97[5],
			bgNumObjectTable,
		),
	}
	BgBallNumber_98 = map[int]*WeightGames{
		1: WeightNewGames(
			bgNumWeightTable_98[1],
			bgNumObjectTable,
		),
		2: WeightNewGames(
			bgNumWeightTable_98[2],
			bgNumObjectTable,
		),
		3: WeightNewGames(
			bgNumWeightTable_98[3],
			bgNumObjectTable,
		),
		4: WeightNewGames(
			bgNumWeightTable_98[4],
			bgNumObjectTable,
		),
		5: WeightNewGames(
			bgNumWeightTable_98[5],
			bgNumObjectTable,
		),
	}
	//每層是否結束的權重
	bgGUTable       = []int{int(bg_up), int(bg_GG)}
	bgGUWeightTable = map[int][]int{
		1: {1, 0},
		2: {70, 30},
		3: {60, 40},
		4: {10, 100},
		5: {0, 1},
	}
	BgOverUp = map[int]*WeightGames{
		1: WeightNewGames(
			bgGUWeightTable[1],
			bgGUTable,
		),
		2: WeightNewGames(
			bgGUWeightTable[2],
			bgGUTable,
		),
		3: WeightNewGames(
			bgGUWeightTable[3],
			bgGUTable,
		),
		4: WeightNewGames(
			bgGUWeightTable[4],
			bgGUTable,
		),
		5: WeightNewGames(
			bgGUWeightTable[5],
			bgGUTable,
		),
	}
	//每層球上分數的權重
	bgPtTable = map[int][]int{
		1: {int(bg_1), int(bg_2), int(bg_3), int(bg_4)},
		2: {int(bg_5), int(bg_7o5), int(bg_10)},
		3: {int(bg_8), int(bg_12)},
		4: {int(bg_15), int(bg_20)},
		5: {int(bg_X2), int(bg_X3), int(bg_X5)},
	}
	bgPtWeightTable_92 = map[int][]int{
		1: {3450, 5000, 3000, 800},
		2: {600, 800, 600},
		3: {1200, 500},
		4: {1300, 450},
		5: {8000, 2500, 1000},
	}

	bgPtWeightTable_97 = map[int][]int{
		1: {3000, 4000, 1500, 800},
		2: {8000, 8000, 6000},
		3: {10000, 4000},
		4: {1200, 450},
		5: {8000, 2500, 1000},
	}

	bgPtWeightTable_98 = map[int][]int{
		1: {3000, 5000, 3000, 600},
		2: {6000, 8000, 6000},
		3: {12000, 5910},
		4: {1000, 500},
		5: {8000, 2500, 1000},
	}

	BgPt = map[string]map[int]*WeightGames{
		"92": BgPt_92,
		"97": BgPt_97,
		"98": BgPt_98,
	}
	BgPt_92 = map[int]*WeightGames{
		1: WeightNewGames(
			bgPtWeightTable_92[1],
			bgPtTable[1],
		),
		2: WeightNewGames(
			bgPtWeightTable_92[2],
			bgPtTable[2],
		),
		3: WeightNewGames(
			bgPtWeightTable_92[3],
			bgPtTable[3],
		),
		4: WeightNewGames(
			bgPtWeightTable_92[4],
			bgPtTable[4],
		),
		5: WeightNewGames(
			bgPtWeightTable_92[5],
			bgPtTable[5],
		),
	}
	BgPt_97 = map[int]*WeightGames{
		1: WeightNewGames(
			bgPtWeightTable_97[1],
			bgPtTable[1],
		),
		2: WeightNewGames(
			bgPtWeightTable_97[2],
			bgPtTable[2],
		),
		3: WeightNewGames(
			bgPtWeightTable_97[3],
			bgPtTable[3],
		),
		4: WeightNewGames(
			bgPtWeightTable_97[4],
			bgPtTable[4],
		),
		5: WeightNewGames(
			bgPtWeightTable_97[5],
			bgPtTable[5],
		),
	}
	BgPt_98 = map[int]*WeightGames{
		1: WeightNewGames(
			bgPtWeightTable_98[1],
			bgPtTable[1],
		),
		2: WeightNewGames(
			bgPtWeightTable_98[2],
			bgPtTable[2],
		),
		3: WeightNewGames(
			bgPtWeightTable_98[3],
			bgPtTable[3],
		),
		4: WeightNewGames(
			bgPtWeightTable_98[4],
			bgPtTable[4],
		),
		5: WeightNewGames(
			bgPtWeightTable_98[5],
			bgPtTable[5],
		),
	}
)
