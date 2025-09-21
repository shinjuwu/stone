package slot_MegShark

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
	MAX_PAYLINE = 50 //這裡代表unitbet
	MAX_SYMBOL  = 41
	unitbet     = 50
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
		slot4002ng,
		slot4002fg,
	)
	// ng
	slot4002ng = NewWayGames(
		ngTable, reelDef, payTable, SymbolList, ScatterPosition, unitbet)

	// Bonus
	slot4002fg = NewWayGames(
		fgTable, reelDef, payTable, SymbolList, ScatterPosition, unitbet)
)
