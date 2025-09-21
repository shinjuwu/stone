package slot_Midas

var (
	colNum     int //去檢查reelDef
	rowNum     int
	payLineNum int
	symbolNum  int
)

const (
	fg_sym_def = 3
)

// payline col
var PayLines = [][]int{
	{1, 1, 1, 1, 1},
	{0, 0, 0, 0, 0},
	{2, 2, 2, 2, 2},
	{0, 1, 2, 1, 0},
	{2, 1, 0, 1, 2},
	{1, 0, 0, 0, 1},
	{1, 2, 2, 2, 1},
	{0, 0, 1, 2, 2},
	{2, 2, 1, 0, 0},
	{1, 2, 1, 0, 1},
	{1, 0, 1, 2, 1},
	{0, 1, 1, 1, 0},
	{2, 1, 1, 1, 2},
	{0, 1, 0, 1, 0},
	{2, 1, 2, 1, 2},
	{1, 1, 0, 1, 1},
	{1, 1, 2, 1, 1},
	{0, 0, 2, 0, 0},
	{2, 2, 0, 2, 2},
	{0, 2, 2, 2, 0},
}

type Reels [][]int

// symbol col
var symbolPayout = [][]int{
	{0, 0, 0, 0, 0},     //default
	{0, 0, 40, 80, 400}, //H1
	{0, 0, 30, 60, 300}, //H2
	{0, 0, 25, 50, 250}, //H3
	{0, 0, 20, 40, 200}, //H4
	{0, 0, 15, 30, 150}, //H5
	{0, 0, 0, 0, 0},     //h6
	{0, 0, 0, 0, 0},     //h7
	{0, 0, 0, 0, 0},     //h8
	{0, 0, 0, 0, 0},     //h9
	{0, 0, 0, 0, 0},     //h10
	{0, 0, 5, 10, 50},   //L1
	{0, 0, 5, 10, 50},   //L2
	{0, 0, 2, 5, 25},    //L3
	{0, 0, 2, 5, 25},    //L4
	{0, 0, 2, 5, 25},    //L5
	{0, 0, 0, 0, 0},     //L6
	{0, 0, 0, 0, 0},     //L7
	{0, 0, 0, 0, 0},     //L8
	{0, 0, 0, 0, 0},     //L9
	{0, 0, 0, 0, 0},     //L10
	{0, 0, 0, 0, 0},     //SF
	{0, 0, 0, 0, 0},     //SB
	{0, 0, 0, 0, 0},     //FS3
	{0, 0, 0, 0, 0},     //FS4
	{0, 0, 0, 0, 0},     //FS5
	{0, 0, 0, 0, 0},     //FS6
	{0, 0, 0, 0, 0},     //FS7
	{0, 0, 0, 0, 0},     //FS8
	{0, 0, 0, 0, 0},     //FS9
	{0, 0, 0, 0, 0},     //FS10
	{0, 0, 0, 0, 0},     //WW
	{0, 0, 0, 0, 0},     //WW_2
	{0, 0, 0, 0, 0},     //WW_3
	{0, 0, 0, 0, 0},     //WW_4
	{0, 0, 0, 0, 0},     //WW_5
	{0, 0, 0, 0, 0},     //WW_6
	{0, 0, 0, 0, 0},     //WW_7
	{0, 0, 0, 0, 0},     //WW_8
	{0, 0, 0, 0, 0},     //WW_9
	{0, 0, 0, 0, 0},     //WW_10

}

var symbolTypeMap = []int{
	0:  SYMBOL_TYPE_WILD,
	1:  SYMBOL_TYPE_NORMAL,
	2:  SYMBOL_TYPE_NORMAL,
	3:  SYMBOL_TYPE_NORMAL,
	4:  SYMBOL_TYPE_NORMAL,
	5:  SYMBOL_TYPE_NORMAL,
	6:  SYMBOL_TYPE_NORMAL,
	7:  SYMBOL_TYPE_NORMAL,
	8:  SYMBOL_TYPE_NORMAL,
	9:  SYMBOL_TYPE_NORMAL,
	10: SYMBOL_TYPE_NORMAL,
	11: SYMBOL_TYPE_NORMAL,
	12: SYMBOL_TYPE_NORMAL,
	13: SYMBOL_TYPE_NORMAL,
	14: SYMBOL_TYPE_NORMAL,
	15: SYMBOL_TYPE_NORMAL,
	16: SYMBOL_TYPE_NORMAL,
	17: SYMBOL_TYPE_NORMAL,
	18: SYMBOL_TYPE_NORMAL,
	19: SYMBOL_TYPE_NORMAL,
	20: SYMBOL_TYPE_NORMAL,
	21: SYMBOL_TYPE_BONUS,
	22: SYMBOL_TYPE_BONUS,
	23: SYMBOL_TYPE_BONUS,
	24: SYMBOL_TYPE_BONUS,
	25: SYMBOL_TYPE_BONUS,
	26: SYMBOL_TYPE_BONUS,
	27: SYMBOL_TYPE_BONUS,
	28: SYMBOL_TYPE_BONUS,
	29: SYMBOL_TYPE_BONUS,
	30: SYMBOL_TYPE_BONUS,
	31: SYMBOL_TYPE_WILD,
	32: SYMBOL_TYPE_WILD,
	33: SYMBOL_TYPE_WILD,
	34: SYMBOL_TYPE_WILD,
	35: SYMBOL_TYPE_WILD,
	36: SYMBOL_TYPE_WILD,
	37: SYMBOL_TYPE_WILD,
	38: SYMBOL_TYPE_WILD,
	39: SYMBOL_TYPE_WILD,
	40: SYMBOL_TYPE_WILD,
}

var symbolTypeStr = []string{
	SYMBOL_TYPE_NORMAL: "normal",
	SYMBOL_TYPE_WILD:   "wild",
	SYMBOL_TYPE_BONUS:  "bonus",
}

var stringToTypeMap map[string]int = map[string]int{
	"normal": SYMBOL_TYPE_NORMAL,
	"wild":   SYMBOL_TYPE_WILD,
	"bonus":  SYMBOL_TYPE_BONUS,
}

var symbolNameList = []string{
	"default",
	"H1",
	"H2",
	"H3",
	"H4",
	"H5",
	"H6",
	"H7",
	"H8",
	"H9",
	"H10",
	"L1",
	"L2",
	"L3",
	"L4",
	"L5",
	"L6",
	"L7",
	"L8",
	"L9",
	"L10",
	"SF",
	"SB",
	"FS3",
	"FS4",
	"FS5",
	"FS6",
	"FS7",
	"FS8",
	"FS9",
	"FS10",
	"WW",
	"WW_2",
	"WW_3",
	"WW_4",
	"WW_5",
	"WW_6",
	"WW_7",
	"WW_8",
	"WW_9",
	"WW_10",
}

var stringToSymbol map[string]int = map[string]int{
	"default": 0,
	"H1":      1,
	"H2":      2,
	"H3":      3,
	"H4":      4,
	"H5":      5,
	"H6":      6,
	"H7":      7,
	"H8":      8,
	"H9":      9,
	"H10":     10,
	"L1":      11,
	"L2":      12,
	"L3":      13,
	"L4":      14,
	"L5":      15,
	"L6":      16,
	"L7":      17,
	"L8":      18,
	"L9":      19,
	"L10":     20,
	"SF":      21,
	"SB":      22,
	"FS3":     23,
	"FS4":     24,
	"FS5":     25,
	"FS6":     26,
	"FS7":     27,
	"FS8":     28,
	"FS9":     29,
	"FS10":    30,
	"WW":      31,
	"WW_2":    32,
	"WW_3":    33,
	"WW_4":    34,
	"WW_5":    35,
	"WW_6":    36,
	"WW_7":    37,
	"WW_8":    38,
	"WW_9":    39,
	"WW_10":   40,
}

func init() {
	basicSetting(MAX_COL, MAX_ROW, MAX_PAYLINE, MAX_SYMBOL)
}

func basicSetting(cn int, rn int, pn int, sn int) {
	colNum = cn
	rowNum = rn
	payLineNum = pn
	symbolNum = sn
}
