package slot_4005Jumphigh

import (
	"slotEditor/utils/random"

	"github.com/shopspring/decimal"
)

// 抽盤面------------------
//
//	@param slotType 遊戲類型 : 0 - NG ; 1 - FG
//	@return pt 盤面
func CreateReels(slotType int, rtp string) (pt Reels, pos []int) {
	pt = make(Reels, colNum)
	Rtp := RTPs(rtp)
	var ReelsLen []int //每組滾輪長度
	switch slotType {
	case 0: //NG滾輪
		ReelsLen, _ = GetNGReelsLen(rtp)
		pos = random.Intsn(ReelsLen)
		for i := 0; i < colNum; i++ {
			pt[i] = make([]int, rowNum+2)
			for j := 0; j < rowNum+2; j++ {
				currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1、2、3號
				if currentPosition >= ReelsLen[i] {
					currentPosition -= ReelsLen[i]
				}
				ngReels := *ngTable[Rtp]
				symbol := ngReels[i][currentPosition]
				pt[i][j] = symbol.Int()
			}
		}
	case 1: //FG滾輪
		ReelsLen, _ = GetFGReelsLen(rtp)
		pos = random.Intsn(ReelsLen)
		for i := 0; i < colNum; i++ {
			pt[i] = make([]int, rowNum+2)
			for j := 0; j < rowNum+2; j++ {
				currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1、2、3號
				if currentPosition >= ReelsLen[i] {
					currentPosition -= ReelsLen[i]
				}
				symbol := fgReelStrips[i][currentPosition]
				pt[i][j] = symbol.Int()
			}
		}
	case 3: //NG必中FG滾輪
		ReelsLen, _ = GetNGReelsLen(rtp)
		pos = GetFGpos()
		for i := 0; i < colNum; i++ {
			pt[i] = make([]int, rowNum+2)
			for j := 0; j < rowNum+2; j++ { //抽出不包含上下兩排:for j := 0; j < MAX_ROW; j++ {
				currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1號
				if currentPosition >= ReelsLen[i] {
					currentPosition -= ReelsLen[i]
				}
				ngReels := *ngTable[Rtp]
				symbol := ngReels[i][currentPosition]
				pt[i][j] = symbol.Int()
			}
		}
	}
	return
}

// output:
//
//	@return	debug盤面
//	@return	wwsym
//	@return	reelpos
//	@return isdebug
func DebugReels(debugIndex [][]int, rtp string) ([]Reels, [][]int, bool) {
	var ok bool
	ngReelsLen, _ := GetNGReelsLen(rtp) //每組滾輪長度
	fgReelsLen, _ := GetFGReelsLen(rtp) //每組滾輪長度
	debugpos := make([]Reels, len(debugIndex))
	Rtp := RTPs(rtp)
	for id, pos := range debugIndex {
		pt := make(Reels, colNum)
		if isDebug(pos[:5]) {
			if id == 0 {
				for i := 0; i < colNum; i++ {
					pt[i] = make([]int, rowNum+2)
					for j := 0; j < rowNum+2; j++ {
						currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1號
						if currentPosition >= ngReelsLen[i] {
							currentPosition -= ngReelsLen[i]
						}
						ngReels := *ngTable[Rtp]
						symbol := ngReels[i][currentPosition]
						pt[i][j] = symbol.Int()
					}
					debugpos[id] = pt
				}
			} else if id != 0 {
				for i := 0; i < colNum; i++ {
					pt[i] = make([]int, rowNum+2)
					for j := 0; j < rowNum+2; j++ {
						currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1號
						if currentPosition >= fgReelsLen[i] {
							currentPosition -= fgReelsLen[i]
						}
						symbol := fgReelStrips[i][currentPosition]
						pt[i][j] = symbol.Int()
					}
					debugpos[id] = pt
				}
				if len(pos) < 6 {
					pos = append(pos, 0)
				}
			}
		} else if !isDebug(pos[:5]) {
			if id == 0 {
				pt, pos = CreateReels(0, rtp)
				debugpos[id] = pt
				debugIndex[id] = pos
			} else if id != 0 {
				if len(pos) < 6 {
					debugIndex[id] = append(pos, 0)
				}
				pt, pos = CreateReels(1, rtp)
				debugpos[id] = pt
				debugIndex[id] = pos
			}
		}
	}
	return debugpos, debugIndex, ok
}
func isDebug(pos []int) (def bool) {
	def = true
	for _, v := range pos {
		if v < 0 || v > 250 {
			def = false
		}
	}
	return
}

// 計算分數--waygame
//
//	@return win		得分倍率
//	@return winLines	得分方式 - 依序獎圖WAY數
//	@return freeRound	免費遊戲次數
//	@return PickSym		被抽中的獎圖
//
// winWay Example:
//
//	此時盤面H2:5連線4way，則會輸出 symbol: H2 得到 4 way 5連線
func CalcWin(pt Reels, rtp string, pos []int, slotType int) (win decimal.Decimal, windetail []Windetail, FreeRound int) {
	firstReel := RemoveDuplicates(pt[0][1:4])
	var match []int
	var multi int
	for _, v := range firstReel {
		if v != SF.Int() {
			match, multi = CalcSymbolsMatchFromLeft(v, pt, int(WW))
			m_count := len(match) - 1
			if m_count > 1 {
				multiscore := decimal.NewFromInt(int64(payTable.CalcPayTable(v, m_count))).Mul(decimal.NewFromInt(int64(multi)))
				multifloat, _ := multiscore.Div(decimal.NewFromInt(unitbet)).Round(2).Float64()
				newWin := Windetail{v, multi, len(match), multifloat}
				windetail = append(windetail, newWin)
				win = win.Add(multiscore)
			}
		}
	}
	win = win.Div(decimal.NewFromInt(int64(unitbet))) // 除上unitbet

	if slotType == 1 {
		//判斷是否觸發FG
		isFG, _ := IsFGWin(pt) // 數SF個數
		if isFG {
			FreeRound += fg_round_time
		}
	}
	return
}

// ------------------------------------------------
func NGflow(pt Reels, rtp string, pos []int, round *Rounds) *Rounds {
	ng := &Records{
		Id:           0,
		SlotType:     0,
		FreeRound:    0,
		Case:         Lose,
		PayReel:      pt,
		ReelPosition: pos,
		Fgtype:       0,
	}
	ng.Bet = round.TotalBet
	point, windetail, _ := CalcWin(pt, rtp, pos, 0)
	round.Result[0] = ng

	ng.Point_Deci = point.Mul(decimal.NewFromFloat(ng.Bet))
	ng.Point, _ = ng.Point_Deci.Round(2).Float64()

	ng.TotalPoint_Deci = ng.TotalPoint_Deci.Add(ng.Point_Deci)
	ng.TotalPoint, _ = ng.TotalPoint_Deci.Round(2).Float64()

	ng.Windetail = windetail
	if point.GreaterThan(decimal.Zero) {
		ng.Case = ng.Case.Push(Win)
	}

	round.TotalPoint_Deci = round.TotalPoint_Deci.Add(ng.TotalPoint_Deci)
	round.TotalPoint, _ = round.TotalPoint_Deci.Round(2).Float64()

	return round
}

func FGflow(stage int, pt Reels, pos []int, round *Rounds, fgmulti int) *Rounds {
	var point decimal.Decimal
	var windetail []Windetail
	var newFreeRound int
	var nowPT Reels
	for i, v := range pt {
		nowPT = append(nowPT, []int{})
		for _, k := range v {
			nowPT[i] = append(nowPT[i], k)
		}
	}
	fg := &Records{
		Id:           stage,
		Stages:       int64(stage),
		SlotType:     1,
		FreeRound:    round.Result[stage-1].FreeRound,
		Case:         Lose,
		PayReel:      nowPT,
		Fg_multi:     fgmulti,
		ReelPosition: pos,
		TotalPoint:   round.TotalPoint,
	}
	//Rtp := round.Rtp
	fg.Bet = round.TotalBet

	point, windetail, newFreeRound = CalcWin(pt, round.Rtp, pos, 1)
	//fgmulti = 1 //test
	fg.Point_Deci = point.Mul(decimal.NewFromInt(int64(fgmulti))).Mul(decimal.NewFromFloat(fg.Bet))
	fg.Point, _ = fg.Point_Deci.Round(2).Float64()

	fg.Windetail = windetail
	if newFreeRound > 0 {
		fg.FreeRound += newFreeRound
		round.FreeSpin += newFreeRound
		//round.Result[0].FreeRound += newFreeRound
		fg.Case = fg.Case.Push(FreeGame)
	}

	if point.GreaterThan(decimal.Zero) {
		fg.Case = fg.Case.Push(Win)
	}
	if stage == 1 {
		fg.TotalPoint_Deci = fg.Point_Deci
	} else if stage > 1 {
		fg.TotalPoint_Deci = round.Result[stage-1].TotalPoint_Deci.Add(fg.Point_Deci)
	}

	fg.TotalPoint, _ = fg.TotalPoint_Deci.Round(2).Float64()

	round.Result[stage] = fg

	round.TotalPoint_Deci = round.TotalPoint_Deci.Add(fg.Point_Deci)
	round.TotalPoint, _ = round.TotalPoint_Deci.Round(2).Float64()

	return round
}

// Spin
//
//	@return round 回合資料
func Spin(rtp string, bet int, debugIndex [][]int, DebugSwitch bool, retrigger int) *Rounds {
	var reel []Reels
	var position [][]int
	var slottype int
	var fgwin decimal.Decimal = decimal.Zero
	if DebugSwitch {
		newreel, newposition, ok := DebugReels(debugIndex, rtp)
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
	fgnum := 0
	//判斷是否觸發FG
	isFG, SFcount := IsFGWin(reel[0])  // 數SF個數
	fgmulti_now := FG_multi[SFcount]   //當前FG倍數
	fgmulti_incr := FG_multi[SFcount]  //每局累加的倍數
	fgmulti_start := FG_multi[SFcount] //起始倍數
	if isFG {
		round.FreeSpin += fg_round_time
		round.Result[0].Case = round.Result[0].Case.Push(FreeGame)
		round.Result[0].Fgtype = fgmulti_start
		round.Result[0].FreeRound += fg_round_time
		fgnum = round.FreeSpin
	}

	//free game------------------
	stage := 0
	for round.FreeSpin > 0 && len(position) > 1 {
		round.FreeSpin--
		stage++
		if stage == len(reel) {
			round.FreeSpin++
			stage--
			break
		}
		round = FGflow(stage, reel[stage], position[stage], round, fgmulti_now)
		round.Result[stage].PayReel = reel[stage]
		round.Result[stage].Fgtype = fgmulti_start
		fgmulti_now += fgmulti_incr
	}
	for round.FreeSpin > 0 {
		round.FreeSpin--
		stage++
		fgreel, newpos := CreateReels(1, rtp)
		round = FGflow(stage, fgreel, newpos, round, fgmulti_now)
		round.Result[stage].PayReel = fgreel
		round.Result[stage].Fgtype = fgmulti_start
		fgmulti_now += fgmulti_incr
	}
	fgwin = round.TotalPoint_Deci.Sub(round.Result[0].TotalPoint_Deci)
	if fgnum > 0 {
		if fgwin.LessThan(decimal.NewFromInt(5)) && fgwin.GreaterThan(decimal.Zero) {
			round = Spin(rtp, bet, debugIndex, DebugSwitch, 3)
			//fmt.Printf("\nhi")
		}
	}
	//j_round, _ := json.Marshal(round)
	//fmt.Printf("\nRound: %s", string(j_round))

	return round
}

// 去掉第一輪重覆的獎圖
func RemoveDuplicates(firstReel []int) (result []int) {
	result = make([]int, 0, len(firstReel))
	temp := map[int]struct{}{}
	for _, item := range firstReel {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// CalcSymbolsMatchFromLeft - 計算轉輪表左邊連線數量
//
//	@param targetSymbol 多個目標獎圖
//	@return []int 返回每輪共有多少個目標獎圖，陣列個數為中獎數量
//	@return int 返回 Way
func CalcSymbolsMatchFromLeft(v int, r Reels, targetSymbol int) ([]int, int) {
	match := []int{}
	multi := 1
	count := 0
	for i := 0; i < len(r); i++ {
		count = CountSym(v, r[i]) + CountSym(targetSymbol, r[i])
		if count <= 0 {
			return match, multi
		}
		multi *= count
		match = append(match, count)

	}
	return match, multi
}

func CountSym(sym int, col []int) (count int) {
	count = 0
	for i := 1; i < 4; i++ {
		if col[i] == sym {
			count++
		}
	}
	return
}
func IsFGWin(reel Reels) (bool, int) {
	countSF, issf := CountScatter(reel, SF.Int())
	if countSF >= fg_sym_def && issf {
		return true, countSF
	} else {
		return false, countSF //0
	}
}

// 計算SF個數
func CountScatter(reels Reels, tar int) (SFCount int, isfg bool) {
	count := 0
	issf := true
	for i := range reels {
		reelcount := 0
		for j, v := range reels[i] {
			if j != 0 && j != 4 {
				if v == tar {
					count++
					reelcount++
				}
			}
		}

		if reelcount > 0 {
			issf = issf && true
		} else if reelcount == 0 {
			issf = issf && false
		}
	}
	return count, issf
}

func GetNGReelsLen(rtp string) ([]int, error) {
	game_math := gameplay
	rtps := RTPs(rtp)
	return game_math.ng_game.GetReelsLen(rtps), nil
}
func GetFGReelsLen(rtp string) ([]int, error) {
	game_math := gameplay
	rtps := RTPs(rtp)
	return game_math.fg_game.GetReelsLen(rtps), nil
}
