package slot_Next3

const (
	MAX_COL     = 5
	MAX_ROW     = 3
	MAX_PAYLINE = 20
	MAX_SYMBOL  = 50
	Unitbet     = 1
)

// slotType
const (
	SLOT_NORMAL = iota
	SLOT_BONUS_FREE
	SLOT_BONUS_MONEY
)

const (
	SYMBOL_TYPE_NORMAL = iota
	SYMBOL_TYPE_WILD
	SYMBOL_TYPE_BONUS
)

type symbol int

type ReelStrip [MAX_COL][MAX_ROW]Symbol

type PayLinesType [MAX_PAYLINE][MAX_COL]int

type symbolPayoutType [MAX_SYMBOL][MAX_COL]int

type symbolWeightType [MAX_COL][MAX_SYMBOL]int

var (
	ReelDef = &ReelStripsDef{3, 3, 3, 3, 3}

	NgTable = ReelStripList{
		"98": &ngReelStrips,
		"97": &ngReelStrips,
		"92": &ngReelStrips,
	}
	FgTable_50 = ReelStripList{
		"98": &fgReelStrips,
		"97": &fgReelStrips,
		"92": &fgReelStrips,
	}

	Gameplay *Games = NewGames(
		slot4004ng,
		slot4004fg,
	)
	// ng
	slot4004ng = NewLineGames(
		NgTable, ReelDef, payTable, SymbolList, ScatterPosition, Unitbet)

	// Bonus
	slot4004fg = NewLineGames(
		FgTable_50, ReelDef, payTable, SymbolList, ScatterPosition, Unitbet)
)
