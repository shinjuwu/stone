package slot_4005Jumphigh

const (
	MAX_COL     = 5
	MAX_ROW     = 3
	MAX_PAYLINE = 10
	MAX_SYMBOL  = 41
	Unitbet     = 10
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

type ReelStrip [MAX_COL][MAX_ROW]Symbol

type PayLinesType [MAX_PAYLINE][MAX_COL]int

type symbolPayoutType [MAX_SYMBOL][MAX_COL]int

type symbolWeightType [MAX_COL][MAX_SYMBOL]int

var (
	ReelDef = &ReelStripsDef{3, 3, 3, 3, 3}

	NgTable = ReelStripList{
		"98": &ngReelStrips_98,
		"97": &ngReelStrips_97,
		"92": &ngReelStrips_92,
	}
	FgTable = ReelStripList{
		"98": &fgReelStrips,
		"97": &fgReelStrips,
		"92": &fgReelStrips,
	}

	Gameplay *Games = NewGames(
		slot4003ng,
		slot4003fg,
	)
	// ng
	slot4003ng = NewWayGames(
		NgTable, ReelDef, payTable, SymbolList, ScatterPosition, Unitbet)

	// Bonus
	slot4003fg = NewWayGames(
		FgTable, ReelDef, payTable, SymbolList, ScatterPosition, Unitbet)
)
