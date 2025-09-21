package slot_4006PryTreasure

import (
	"slotEditor/utils/random"

	"github.com/shopspring/decimal"
)

// 抽盤面------------------
//
//	normal game
//	@param slotType 遊戲類型 : 0 - NG ; 3 - NG必中BG
//	@return pt 盤面
func CreateReels(slotType int, rtp string) (pt Reels, pos []int) {
	pt = make(Reels, colNum)
	var ReelsLen []int //每組滾輪長度//NG滾輪
	switch slotType {
	case 0:
		ReelsLen, _ = GetNGReelsLen(rtp)
		pos = random.Intsn(ReelsLen)
		for i := 0; i < colNum; i++ {
			pt[i] = make([]float64, rowNum+2)
			for j := 0; j < rowNum+2; j++ {
				currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1、2、3號
				if currentPosition >= ReelsLen[i] {
					currentPosition -= ReelsLen[i]
				}
				symbol := ngReelStrips[i][currentPosition]
				pt[i][j] = float64(symbol.Int())
			}
		}
	case 3: //NG必中FG滾輪
		ReelsLen, _ = GetNGReelsLen(rtp)
		pos = GetFGpos()
		for i := 0; i < colNum; i++ {
			pt[i] = make([]float64, rowNum+2)
			for j := 0; j < rowNum+2; j++ { //抽出不包含上下兩排:for j := 0; j < MAX_ROW; j++ {
				currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1號
				if currentPosition >= ReelsLen[i] {
					currentPosition -= ReelsLen[i]
				}
				symbol := ngReelStrips[i][currentPosition]
				pt[i][j] = float64(symbol.Int())
			}
		}
	}
	return
}

// output:
//
//	@return	debug盤面
//	@return	reelpos
//	@return	Bg debug 層數
//	@return isdebug
func DebugReels(debugIndex [][]int, rtp string) ([]Reels, [][]int, int, bool) {
	var ok bool
	var bglevelnum int = 0
	ngReelsLen, _ := GetNGReelsLen(rtp) //每組滾輪長度
	debugpos := make([]Reels, len(debugIndex))
	for id, pos := range debugIndex {
		pt := make(Reels, colNum)
		if isDebug(pos[:5]) {
			if id == 0 {
				for i := 0; i < colNum; i++ {
					pt[i] = make([]float64, rowNum+2)
					for j := 0; j < rowNum+2; j++ {
						currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1號
						if currentPosition >= ngReelsLen[i] {
							currentPosition -= ngReelsLen[i]
						}
						symbol := ngReelStrips[i][currentPosition]
						pt[i][j] = float64(symbol.Int())
					}
					debugpos[id] = pt
				}
			}
		} else if !isDebug(pos[:5]) {
			if id == 0 {
				pt, pos = CreateReels(0, rtp)
				debugpos[id] = pt
				debugIndex[id] = pos
			}
		}
		if pos[5] > 0 && pos[5] < 6 {
			bglevelnum = pos[5]
		}
	}

	return debugpos, debugIndex, bglevelnum, ok
}
func isDebug(pos []int) (def bool) {
	def = true
	for _, v := range pos {
		if v < 0 || v > 80 {
			def = false
		}
	}
	return
}

// 計算分數--waygame
//
//	@return win		得分倍率
//	@return wildReel	Wild獎項倍率
//	@windetail	得分明細
func CalcWin(pt Reels, rtp string, pos []int, wildreel []int) (decimal.Decimal, []int, []Windetail) {
	firstReel := RemoveDuplicates(pt[0][1:4])
	var match []int
	var multi int
	var windetail []Windetail
	var win decimal.Decimal
	max_match := 0
	for _, v := range firstReel {
		if v != SB.Int() {
			match, multi = CalcSymbolsMatchFromLeft(v, pt, wildreel)
			if len(match) > max_match {
				max_match = len(match)
			}
			m_count := len(match) - 1
			if m_count > 1 {
				multiscore := decimal.NewFromInt(int64(payTable.CalcPayTable(v, m_count))).Mul(decimal.NewFromInt(int64(multi)))
				multifloat, _ := multiscore.Div(decimal.NewFromInt(unitbet)).Float64()
				newWin := Windetail{v, multi, len(match), multifloat}
				windetail = append(windetail, newWin)
				win = win.Add(multiscore)
			}
		}
	}
	for i, v := range wildreel {
		if i > max_match-1 || max_match < 3 {
			if v == 1 {
				wildreel[i] = 0
			}
		}
	}
	win = win.Div(decimal.NewFromInt(int64(unitbet))) // 除上unitbet
	return win, wildreel, windetail
}

// ------------------------------------------------
func NGflow(pt Reels, rtp string, pos []int, round *Rounds) *Rounds {
	ng := &Records{
		Id:           0,
		SlotType:     0,
		Case:         Lose,
		ReelPosition: pos,
	}
	ng.PayReel = pt
	ng.Bet = round.TotalBet
	wwr := Wildr(pt)
	point, wwReel, windetail := CalcWin(pt, rtp, pos, wwr)
	ng.NgWildExpand = wwReel

	ng.Point_Deci = point
	ng.Point, _ = point.Round(2).Float64()

	ng.TotalPoint_Deci = ng.TotalPoint_Deci.Add(ng.Point_Deci)
	ng.TotalPoint, _ = ng.TotalPoint_Deci.Round(2).Float64()

	ng.Windetail = windetail
	if point.GreaterThan(decimal.Zero) {
		ng.Case = ng.Case.Push(Win)
	}
	round.TotalPoint_Deci = round.TotalPoint_Deci.Add(ng.Point_Deci)
	round.TotalPoint, _ = round.TotalPoint_Deci.Round(2).Float64()
	round.Result[0] = ng
	return round
}

func BGflow(rtp string, round *Rounds, debug_bg int) *Rounds {
	size := []uint{7, 5, 4, 3, 1}
	var Level []*BGames
	bgid := 0
	for i := 0; i < 5; i++ {
		Level = append(Level, NewBonusGame(BgBallNumber[rtp][i+1], BgPt[rtp][i+1], BgOverUp[i+1], size[i]))
		Prize, UpOver := Level[i].SpinBonus(i + 1)
		if debug_bg > 1 {
			if i+1 < debug_bg {
				UpOver = []int{1}
			} else {
				UpOver = []int{0}
			}
		}
		for j := 0; j < len(Prize)+1; j++ {
			bgid++
			bg := &Records{
				Id:         bgid,
				Bet:        round.TotalBet,
				Stages:     int64(bgid),
				BgLevel:    i + 1,
				SlotType:   1,
				Point_Deci: decimal.Zero,
			}
			bg.Case = bg.Case.Push(Bonus)
			if j < len(Prize) {
				bg.PayReel = append(bg.PayReel, []float64{Prize[j]})
			} else if i == 4 && j == len(Prize) {
				break
			} else if j == len(Prize) {
				bg.BgExtra = UpOver
			}

			pt := decimal.Zero
			if i == 0 && j == 0 {
				bg.TotalPoint_Deci = decimal.Zero
			} else {
				bg.TotalPoint_Deci = round.Result[bgid-1].TotalPoint_Deci
			}
			for _, v := range bg.PayReel {
				for _, k := range v {
					pt = pt.Add(decimal.NewFromFloat(k).Round(4))
				}
			}

			bg.Point_Deci = pt
			bg.Point, _ = bg.Point_Deci.Round(4).Float64()

			if i+1 < 5 {
				if i == 0 && j == 0 {
					bg.TotalPoint_Deci = bg.Point_Deci
					bg.TotalPoint = bg.Point
				} else {
					bg.TotalPoint_Deci = bg.TotalPoint_Deci.Add(bg.Point_Deci)
					bg.TotalPoint, _ = bg.TotalPoint_Deci.Round(4).Float64()
				}

			} else if bg.BgLevel == 5 {
				bg.TotalPoint_Deci = round.Result[bgid-1].TotalPoint_Deci.Mul(bg.Point_Deci)
				bg.TotalPoint, _ = bg.Point_Deci.Round(4).Float64()

				round.TotalPoint_Deci = round.TotalPoint_Deci.Add(bg.TotalPoint_Deci)
				round.TotalPoint, _ = round.TotalPoint_Deci.Float64()
			}

			round.Result[bgid] = bg

		}
		if len(UpOver) > 0 {
			if UpOver[0] == 0 {
				break
			}
		}
	}
	round.FreeSpin = len(round.Result) - 1
	round.TotalPoint_Deci = round.Result[len(round.Result)-1].TotalPoint_Deci.Add(round.Result[0].TotalPoint_Deci)
	round.TotalPoint, _ = round.TotalPoint_Deci.Round(4).Float64()
	return round
}

// Spin
//
//	@param	retrigger	- 0:server請輸入0 ; 1:必中NG
//	@return round 回合資料
func Spin(rtp string, bet int, debugIndex [][]int, DebugSwitch bool, retrigger int) *Rounds {
	var reel []Reels
	var position [][]int
	var debug_bg int
	var BGhit bool = false
	var slottype int
	Bgwin := decimal.Zero
	if DebugSwitch {
		newreel, newposition, bglevelnum, ok := DebugReels(debugIndex, rtp)
		debug_bg = bglevelnum
		position = newposition
		reel = newreel
		if !ok {
			ngreel, newpos := CreateReels(0, rtp)
			position = append(position, newpos)
			reel = append(reel, ngreel)
		}
	} else {
		slottype = 0
		if retrigger != 0 {
			slottype = 3
		}
		ngreel, newpos := CreateReels(slottype, rtp)
		position = append(position, newpos)
		reel = append(reel, ngreel)
	}
	//built new Round
	round := NewRounds()
	round.Rtp = rtp
	round.TotalBet = float64(bet)
	round = NGflow(reel[0], rtp, position[0], round)

	//判斷是否觸發FG
	bonusCount := CountBouns(reel[0]) // 數SB個數
	if bonusCount >= bg_sym_def {
		round.Result[0].Case = round.Result[0].Case.Push(Bonus)
		round = BGflow(rtp, round, debug_bg)
		BGhit = true
	}
	//
	Bgwin = round.TotalPoint_Deci.Sub(round.Result[0].TotalPoint_Deci)
	if BGhit {
		if Bgwin.LessThan(decimal.NewFromInt(5)) {
			round = Spin(rtp, bet, debugIndex, DebugSwitch, 3)
		}
	}
	//
	//part of free game----------
	// if len(round.Result) > 1 {
	//j_round, _ := json.Marshal(round)
	//fmt.Printf("\nRound: %s", string(j_round))
	// }
	return round
}

// 去掉第一輪重覆的獎圖
func RemoveDuplicates(firstReel []float64) (result []int) {
	result = make([]int, 0, len(firstReel))
	temp := map[int]struct{}{}
	for _, item := range firstReel {
		itemint := int(item)
		if _, ok := temp[itemint]; !ok {
			temp[itemint] = struct{}{}
			result = append(result, itemint)
		}
	}
	return result
}

// CalcSymbolsMatchFromLeft - 計算轉輪表左邊連線數量
//
//	@param targetSymbol 多個目標獎圖
//	@return []int 返回每輪共有多少個目標獎圖，陣列個數為中獎數量
//	@return int 返回 Way
func CalcSymbolsMatchFromLeft(v int, r Reels, wwr []int) ([]int, int) {
	match := []int{}
	multi := 1
	count := 0
	for i := 0; i < len(r); i++ {
		if CountSym(WW.Int(), r[i]) > 0 {
			count = 3
		} else {
			count = CountSym(v, r[i])
		}
		if count <= 0 {
			return match, multi
		}
		multi *= count
		match = append(match, count)
	}
	return match, multi
}

func CountSym(sym int, col []float64) (count int) {
	count = 0
	for i := 1; i < 4; i++ {
		if int(col[i]) == sym {
			count++
		}
	}
	return
}

// 計算SB個數
func CountBouns(pt Reels) (bounsCount int) {
	for _, rowLine := range pt {
		for i := 1; i < 4; i++ {
			if int(rowLine[i]) == SB.Int() {
				bounsCount++
			}
		}
	}
	return
}

func GetNGReelsLen(rtp string) ([]int, error) {
	game_math := gameplay[rtp]
	rtps := RTPs(rtp)
	return game_math.ng_game.GetReelsLen(rtps), nil
}

func ExpandWWReel(r Reels, ww [5]int) Reels {
	NewR := CopyReel(r)
	for i, v := range ww {
		if v == 1 {
			for j := range NewR[i] {
				NewR[i][j] = 31
			}
		}
	}
	return NewR
}

func CopyReel(reel Reels) Reels {
	var copyr Reels
	for i := range reel {
		copyr = append(copyr, []float64{})
		for _, k := range reel[i] {
			copyr[i] = append(copyr[i], k)
		}
	}
	return copyr
}

func Wildr(r Reels) []int {
	var wildr []int
	for m := 0; m < 5; m++ {
		wildr = append(wildr, 0)
	}
	for i, v := range r {
		for j, k := range v {
			if k == 31 && (j > 0 && j < 4) {
				wildr[i] = 1
			}
		}
	}
	return wildr
}

func IsExpand(wwr [5]int) (exp bool) {
	exp = false
	count := 0
	for _, v := range wwr {
		count += v
	}
	if count > 0 {
		exp = true
	}
	return
}
