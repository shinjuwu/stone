package slot_4005Jumphigh

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
	MAX_PAYLINE = 10 //這裡代表unitbet
	MAX_SYMBOL  = 41
	unitbet     = 10
)

var (
	reelDef = &ReelStripsDef{3, 3, 3, 3, 3}

	ngTable = ReelStripList{
		"98": &ngReelStrips_98,
		"97": &ngReelStrips_97,
		"92": &ngReelStrips_92,
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
