package slot_Midas

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
	MAX_SYMBOL  = 41
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
		slot4003ng,
		slot4003fg,
	)
	// ng
	slot4003ng = NewLineGames(
		ngTable, reelDef, payTable, SymbolList, ScatterPosition, unitbet)

	// Bonus
	slot4003fg = NewLineGames(
		fgTable, reelDef, payTable, SymbolList, ScatterPosition, unitbet)
)
