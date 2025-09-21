package slot_4004WildGem

import (
	"strconv"

	"github.com/shopspring/decimal"
)

type Games struct {
	ng_game *LineGames
	fg_game *LineGames
}
type LineGames struct {
	reelStripsTable ReelStripList
	reelStripsDef   *ReelStripsDef
	reelLen         ReelStripLengthTable
	payTable        *PayTable
	symbolList      []Symbol
	scatter         []Symbol
	unitbet         float64
}

func NewGames(ng, fg *LineGames) *Games {
	return &Games{
		ng_game: ng,
		fg_game: fg,
	}
}

// Windetail
//
//	@Symbol 中獎獎圖
//	@WinWay 連線Way數
//	@Match 連線個數
//	@Point 得分
type Windetail struct {
	Symbol  int
	WinLine int
	Match   int
	Multi   float64
}

// NewLineGames - 建立 WayGames
//
//	@param reelStripsTable			轉輪表 < RTP, 轉輪表 >
//	@param reelStripsDef			轉輪個數，陣列大小為幾輪，陣列內容為每輪的數量
//	@param payTable					賠付表
//	@param symbolList				獎圖列表，可使用預設 slots.SymbolList，也可以自定義
//	@param scatter Scatter			特殊獎圖，可使用預設 slots.Scatter，也可以自定義
//	@param unitbet float64	單位投注
//	@return *WayGames WayGames 物件
func NewLineGames(
	reelStripsTable ReelStripList,
	reelStripsDef *ReelStripsDef,
	payTable *PayTable,
	symbolList []Symbol,
	scatter []Symbol,
	unitbet float64,
) *LineGames {
	// get reel length
	reelLenTable := make(ReelStripLengthTable)
	for rtp, reelStrips := range reelStripsTable {
		reelLenTable[rtp] = reelStrips.Lengths()
	}
	return &LineGames{
		reelStripsTable: reelStripsTable,
		reelStripsDef:   reelStripsDef,
		reelLen:         reelLenTable,
		payTable:        payTable,
		symbolList:      symbolList,
		scatter:         scatter,
		unitbet:         unitbet,
	}
}

// Rounds - 遊戲回合單
//
//	一個 Round 包含一個下注到最後贏分
//	@rtp	- 回合現行rtp
//	@SlotType	-回合遊戲類別: 0 表示NG ; 1 表示FG
//	@TotalBet	- 總下注
//	@TotalPoint	- 總贏分
type Rounds struct {
	Rtp             string          `json:"rtp"`
	Result          Results         `json:"result"`
	FreeSpin        int             `json:"freespin"`
	TotalBet        float64         `json:"totalbet"`
	TotalPoint      float64         `json:"totalpoint"`
	TotalPoint_Deci decimal.Decimal `json:"totalpoint_dec`
}

func NewRounds() *Rounds {
	return &Rounds{
		Rtp:             "",
		Result:          NewResults(),
		FreeSpin:        0,
		TotalBet:        float64(0),
		TotalPoint:      float64(0),
		TotalPoint_Deci: decimal.Zero,
	}
}

type Results map[int]*Records

// NewResults - 建立回合結果
func NewResults() Results {
	return Results{}
}

// Records - 遊戲紀錄
//
//	@Id		- record id
//	@SlotType	- 回合遊戲類別: 0 表示NG ; 1 表示FG
//	@Case		- 0=lose 1=win 16=lose&&getfg 17=win&&getfg
//	@FreeRound	- 總轉次
//	@Retrigger	- FG中再觸發新增轉次
//	@Stages 	- 該回合狀態階段或局數，例如：FreeGame的1(第一局)、2(第二局)
//	@Pickem		- 選擇項目
//	@payReel	- 回合盤面
//	@WinLine	- 贏分明細
//	@WildPrize	- Wild中獎獎項
//	@ReelPosition	- 該局中獎的Reel Index -便於debug mode
//	@Point		- 該局得分
//	@Bet		- 下注
//	@TotalPoint	- 總贏分
type Records struct {
	Id              int             `json:"id"`
	SlotType        int             `json:"slotType"`
	Case            State           `json:"case"`
	FreeRound       int             `json:"freeRound"`
	Stages          int64           `json:"stages"`
	Retrigger       int             `json:"retrigger"`
	PayReel         Reels           `json:"PayReels"`
	Windetail       []Windetail     `json:"winline"`
	WildPrize       []float64       `json:"wildprize"`
	ReelPosition    []int           `json:"reelposition"` //for debug
	Point           float64         `json:"point"`
	Point_Deci      decimal.Decimal `json:"point_deci"`
	Bet             float64         `json:"bet"`
	TotalPoint      float64         `json:"TotalPoint"`
	TotalPoint_Deci decimal.Decimal `json:"TotalPoint_deci"`
}
type State uint16

const (
	Lose     = State(0x0000)
	Win      = State(0x0001)
	FreeGame = State(0x0010)
	Bonus    = State(0x0020)
)

// Push -
// Push the state.
func (s State) Push(state State) State {
	return s | state
}

// Pop -
// Pop the state.
func (s State) Pop(state State) State {
	return s &^ state
}

// IsWin -
// Check if the state is a win.
func (s State) IsWin() bool {
	return (s & Win) == Win
}

// IsLose -
// Check if the state is a lose.
func (s State) IsLose() bool {
	return (s & Lose) == Lose
}

// IsFreeGame -
// Check if the state is a free game.
func (s State) IsFreeGame() bool {
	return (s & FreeGame) == FreeGame
}

// IsBonus -
// Check if the state is a bonus.
func (s State) IsBonus() bool {
	return (s & Bonus) == Bonus
}

// Lengths - 轉輪表長度
func (r ReelStrips) Lengths() []int {
	lengths := make([]int, len(r))
	for i, reel := range r {
		lengths[i] = len(reel)
	}
	return lengths
}

var (
	SymbolList      = []Symbol{DEFAULT, H1, H2, H3, H4, H5, H6, H7, H8, H9, H10, L1, L2, L3, L4, L5, L6, L7, L8, L9, L10, SF, SB, FS3, FS4, FS5, FS6, FS7, FS8, FS9, FS10, WW, WW_0d3, WW_0d5, WW_0d8, WW_1, WW_1d5, WW_2, WW_3, WW_4, WW_5, WW_6, WW_7, WW_8, WW_9, WW_MINI, WW_MINOR, WW_MAJOR, WW_SUPER, WW_GRAND}
	ScatterPosition = []Symbol{SF, SB}
)

const (
	// default
	DEFAULT = Symbol(0)
	// SymbolHighPay1 - High Pay 1
	H1 = Symbol(1)
	// SymbolHighPay2 - High Pay 2
	H2 = Symbol(2)
	// SymbolHighPay3 - High Pay 3
	H3 = Symbol(3)
	// SymbolHighPay4 - High Pay 4
	H4 = Symbol(4)
	// SymbolHighPay5 - High Pay 5
	H5 = Symbol(5)
	// SymbolHighPay6 - High Pay 6
	H6 = Symbol(6)
	//  SymbolHighPay7 - High Pay 7
	H7 = Symbol(7)
	//  SymbolHighPay8 - High Pay 8
	H8 = Symbol(8)
	//  SymbolHighPay9 - High Pay 9
	H9 = Symbol(9)
	//  SymbolHighPay10 - High Pay 10
	H10 = Symbol(10)
	//  SymbolLowPay1 - Low Pay 1
	L1 = Symbol(11)
	//  SymbolLowPay2 - Low Pay 2
	L2 = Symbol(12)
	//  SymbolLowPay3 - Low Pay 3
	L3 = Symbol(13)
	//  SymbolLowPay4 - Low Pay 4
	L4 = Symbol(14)
	//  SymbolLowPay5 - Low Pay 5
	L5 = Symbol(15)
	//  SymbolLowPay6 - Low Pay 6
	L6 = Symbol(16)
	//  SymbolLowPay7 - Low Pay 7
	L7 = Symbol(17)
	//  SymbolLowPay8 - Low Pay 8
	L8 = Symbol(18)
	//  SymbolLowPay9 - Low Pay 9
	L9 = Symbol(19)
	//  SymbolLowPay10 - Low Pay 10
	L10 = Symbol(20)
	// SymbolFreeSpin - Scatter Free Spin
	SF = Symbol(21)
	// SymbolBonus - Scatter Bonus
	SB = Symbol(22)
	// Symbol Feature - Feature Symbol 3
	FS3 = Symbol(23)
	// Symbol Feature - Feature Symbol 4
	FS4 = Symbol(24)
	// Symbol Feature - Feature Symbol 5
	FS5 = Symbol(25)
	// Symbol Feature - Feature Symbol 6
	FS6 = Symbol(26)
	// Symbol Feature - Feature Symbol 7
	FS7 = Symbol(27)
	// Symbol Feature - Feature Symbol 8
	FS8 = Symbol(28)
	// Symbol Feature - Feature Symbol 9
	FS9 = Symbol(29)
	// Symbol Feature - Feature Symbol 10
	FS10 = Symbol(30)
	// SymbolWild -WW
	WW = Symbol(31)
	// SymbolWild -WW_2
	WW_0d3 = Symbol(32)
	// SymbolWild -WW_3
	WW_0d5 = Symbol(33)
	// SymbolWild -WW_4
	WW_0d8 = Symbol(34)
	// SymbolWild -WW_5
	WW_1 = Symbol(35)
	// SymbolWild -WW_6
	WW_1d5 = Symbol(36)
	// SymbolWild -WW_7
	WW_2 = Symbol(37)
	// SymbolWild -WW_8
	WW_3 = Symbol(38)
	// SymbolWild -WW_9
	WW_4 = Symbol(39)
	// SymbolWild -WW_10
	WW_5 = Symbol(40)
	// SymbolWild -WW_2
	WW_6 = Symbol(41)
	// SymbolWild -WW_3
	WW_7 = Symbol(42)
	// SymbolWild -WW_4
	WW_8 = Symbol(43)
	// SymbolWild -WW_5
	WW_9 = Symbol(44)
	// SymbolWild -WW_6
	WW_MINI = Symbol(45)
	// SymbolWild -WW_7
	WW_MINOR = Symbol(46)
	// SymbolWild -WW_8
	WW_MAJOR = Symbol(47)
	// SymbolWild -WW_9
	WW_SUPER = Symbol(48)
	// SymbolWild -WW_10
	WW_GRAND     = Symbol(49)
	SYMBOL_COUNT = 50
)

func (w *LineGames) GetReelsLen(rtp RTPs) []int {
	return w.reelLen[rtp]
}

// Int - 轉換為整數
func (s Symbol) Int() int {
	return int(s)
}

// String - 轉換為字串
func (s Symbol) String() string {
	return strconv.Itoa(int(s))
}
