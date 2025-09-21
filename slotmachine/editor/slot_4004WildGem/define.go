package slot_4004WildGem

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

const (
	MAX_COL     = 5
	MAX_ROW     = 3
	MAX_PAYLINE = 20
	MAX_SYMBOL  = 50
	unitbet     = 1
)

var (
	reelDef = &ReelStripsDef{3, 3, 3, 3, 3}

	ngTable = ReelStripList{
		"98": &ngReelStrips,
		"97": &ngReelStrips,
		"92": &ngReelStrips,
	}
	fgTable = ReelStripList{
		"98": &fgReelStrips,
		"97": &fgReelStrips,
		"92": &fgReelStrips,
	}

	gameplay *Games = NewGames(
		slot4004ng,
		slot4004fg,
	)
	// ng
	slot4004ng = NewLineGames(
		ngTable, reelDef, payTable, SymbolList, ScatterPosition, unitbet)

	// Bonus
	slot4004fg = NewLineGames(
		fgTable, reelDef, payTable, SymbolList, ScatterPosition, unitbet)
)
