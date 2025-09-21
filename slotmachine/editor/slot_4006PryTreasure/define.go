package slot_4006PryTreasure

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
		"98": &ngReelStrips,
		"97": &ngReelStrips,
		"92": &ngReelStrips,
	}

	gameplay map[string]*Games = map[string]*Games{
		"98": NewGames(
			slot4005ng,
			slot4005bg[98],
		),
		"97": NewGames(
			slot4005ng,
			slot4005bg[97],
		),
		"92": NewGames(
			slot4005ng,
			slot4005bg[92],
		),
	}
	// ng
	slot4005ng = NewWayGames(
		ngTable, reelDef, payTable, SymbolList, ScatterPosition, unitbet)

	// Bonus
	slot4005bg = map[int]*BGames{
		98: NewBonusGame(BgBallNumber["98"][5], BgPt["98"][5], BgOverUp[5], 19),
		97: NewBonusGame(BgBallNumber["97"][5], BgPt["97"][5], BgOverUp[5], 19),
		92: NewBonusGame(BgBallNumber["92"][5], BgPt["92"][5], BgOverUp[5], 19),
	}
)
