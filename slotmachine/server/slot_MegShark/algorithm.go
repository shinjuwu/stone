package slot_MegShark

import (
	"slotserver/utils/random"

	"github.com/shopspring/decimal"
)

// 抽盤面------------------
//
//	@Param slotType 遊戲類型 : 0 - NG ; 1 - FG；3 - NG且必中FG
//	@return pt 盤面
//	@return pos	reelindex
func CreateReels(slotType int, rtp string) (pt Reels, pos []int) {
	pt = make(Reels, MAX_COL)
	var ReelsLen []int //每組滾輪長度
	switch slotType {
	case 0: //NG滾輪
		ReelsLen, _ = GetNGReelsLen(rtp)
		pos = random.Intsn(ReelsLen)
		for i := 0; i < MAX_COL; i++ {
			pt[i] = make([]int, MAX_ROW+2)
			for j := 0; j < MAX_ROW+2; j++ {
				currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1號
				if currentPosition >= ReelsLen[i] {
					currentPosition -= ReelsLen[i]
				}
				symbol := ngReelStrips[i][currentPosition]
				pt[i][j] = symbol.Int()
			}
		}
	case 1: //FG滾輪
		ReelsLen, _ = GetFGReelsLen(rtp)
		pos = random.Intsn(ReelsLen)
		for i := 0; i < MAX_COL; i++ {
			pt[i] = make([]int, MAX_ROW+2)
			for j := 0; j < MAX_ROW+2; j++ {
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
		//dice := random.Intn(Rand_SFcount.Sum())
		//pick, _ := Rand_SFcount.Pick(dice)
		//SFpos := GetSFpos(pick)
		fgpos := []int{1, 0, 1, 0, 1}
		pos = GetFGpos(fgpos)
		for i := 0; i < MAX_COL; i++ {
			pt[i] = make([]int, MAX_ROW+2)
			for j := 0; j < MAX_ROW+2; j++ { //抽出不包含上下兩排:for j := 0; j < MAX_ROW; j++ {
				currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1號
				if currentPosition >= ReelsLen[i] {
					currentPosition -= ReelsLen[i]
				}
				symbol := ngReelStrips[i][currentPosition]
				pt[i][j] = symbol.Int()
			}
		}
	}
	return
}

// output:
//
//	@return	debugpos	debug盤面
//	@return	wwsym	Debug指定被替換獎圖
//	@return	reelpos	Debug指定滾輪位置
//	@return isdebug	是否符合Debug條件
func DebugReels(debugIndex [][]int, rtp string) ([]Reels, []int, [][]int, bool) {
	var ok bool
	var debugWWsym []int
	ngReelsLen, _ := GetNGReelsLen(rtp) //每組滾輪長度
	fgReelsLen, _ := GetFGReelsLen(rtp) //每組滾輪長度
	debugpos := make([]Reels, len(debugIndex))
	for id, pos := range debugIndex {
		pt := make(Reels, MAX_COL)
		if isDebug(pos[:5]) {
			if id == 0 {
				for i := 0; i < MAX_COL; i++ {
					pt[i] = make([]int, MAX_ROW+2)
					for j := 0; j < MAX_ROW+2; j++ {
						currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1、2、3號
						if currentPosition >= ngReelsLen[i] {
							currentPosition -= ngReelsLen[i]
						}
						symbol := ngReelStrips[i][currentPosition]
						pt[i][j] = symbol.Int()
					}
					debugpos[id] = pt
				}
			} else if id != 0 {
				for i := 0; i < MAX_COL; i++ {
					pt[i] = make([]int, MAX_ROW+2)
					for j := 0; j < MAX_ROW+2; j++ {
						currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1、2、3號
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
				if MatchWW(pos[MAX_COL]) {
					debugWWsym = append(debugWWsym, pos[MAX_COL])
				} else {
					dice := random.Intn(FG_randomSym[rtp].Sum())
					pick, _ := FG_randomSym[rtp].Pick(dice)
					debugWWsym = append(debugWWsym, pick)
				}
			}
		} else if !isDebug(pos[:5]) {
			if id == 0 {
				pt, pos = CreateReels(0, rtp)
				debugpos[id] = pt
				debugIndex[id] = pos
			} else if id != 0 {
				if len(pos) < 6 {
					pos = append(pos, 0)
				}
				if MatchWW(pos[MAX_COL]) {
					debugWWsym = append(debugWWsym, pos[MAX_COL])
				} else {
					dice := random.Intn(FG_randomSym[rtp].Sum())
					pick, _ := FG_randomSym[rtp].Pick(dice)
					debugWWsym = append(debugWWsym, pick)
				}
				pt, pos = CreateReels(1, rtp)
				debugpos[id] = pt
				debugIndex[id] = pos

			}
		}
	}
	return debugpos, debugWWsym, debugIndex, ok
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
func MatchWW(debugwwsym int) bool {
	if (debugwwsym > 0 && debugwwsym < 6) || (debugwwsym > 10 && debugwwsym < 17) {
		return true
	} else {
		return false
	}
}

// 計算分數--waygame
//
//	@return win		得分倍率
//	@return wwReel	如進入FG，輸出替換獎圖後的Reel
//	@windetail	得分細目
//	@return freeRound	FG剩餘Spin數
func CalcWin(pt Reels, rtp string, pos []int, slotType int, WWsym int) (win decimal.Decimal, wwReel Reels, windetail []Windetail, FreeRound int) {
	if slotType == 1 {
		newpt := RandomWW(rtp, pt, WWsym) //抽出該獎圖並換掉
		wwReel = newpt                    //newpt
	}
	firstReel := RemoveDuplicates(pt[0][1:4])
	var match []int
	var multi int
	for _, v := range firstReel {
		if v != SF.Int() {
			if slotType == 1 {
				match, multi = CalcSymbolsMatchFromLeft(v, wwReel, int(WW))
			} else if slotType == 0 {
				match, multi = CalcSymbolsMatchFromLeft(v, pt, int(WW))
			}
			m_count := len(match) - 1
			if m_count > 1 {
				multiscore := decimal.NewFromInt(int64(payTable.CalcPayTable(v, m_count))).Mul(decimal.NewFromInt(int64(multi)))
				multifloat, _ := multiscore.Div(decimal.NewFromInt(Unitbet)).Round(2).Float64()
				newWin := Windetail{v, multi, len(match), multifloat}
				windetail = append(windetail, newWin)
				win = win.Add(multiscore)
			}
		}
	}
	win = win.Div(decimal.NewFromInt(int64(Unitbet))) // 除上unitbet
	if slotType == 1 {
		//判斷是否觸發FG
		bonusCount := CountBouns(wwReel) // 數SB個數
		if bonusCount >= fg_sym_def {
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
		PreReel:      pt,
		ReelPosition: pos,
	}
	ng.Bet = round.TotalBet
	point, _, windetail, _ := CalcWin(pt, rtp, pos, 0, 0)

	ng.Point_Deci = point
	ng.Point, _ = ng.Point_Deci.Round(2).Float64()

	ng.TotalPoint_Deci = ng.TotalPoint_Deci.Add(ng.Point_Deci)
	ng.TotalPoint, _ = ng.TotalPoint_Deci.Round(2).Float64()
	ng.Windetail = windetail
	if point.GreaterThan(decimal.Zero) {
		ng.Case = ng.Case.Push(Win)
	}

	round.TotalPoint_Deci = round.TotalPoint_Deci.Add(ng.TotalPoint_Deci)
	round.TotalPoint, _ = round.TotalPoint_Deci.Round(2).Float64()
	round.Result[0] = ng
	return round
}

func FGflow(stage int, pt Reels, pos []int, round *Rounds, WWsym []int) *Rounds {
	var point decimal.Decimal
	var WWReel Reels
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
		PreReel:      nowPT,
		ReelPosition: pos,
		TotalPoint:   float64(0),
	}
	if stage == 1 {
		fg.TotalPoint_Deci = decimal.Zero
	}
	fg.Bet = round.TotalBet
	rtp := round.Rtp
	var wwsym int
	if stage <= len(WWsym) {
		wwsym = WWsym[stage-1]
	} else {
		dice := random.Intn(FG_randomSym[rtp].Sum())
		pick, _ := FG_randomSym[rtp].Pick(dice)
		wwsym = pick
	}
	point, WWReel, windetail, newFreeRound = CalcWin(pt, rtp, pos, 1, wwsym)

	fg.Point_Deci = point
	fg.Point, _ = fg.Point_Deci.Round(2).Float64()

	fg.WwReel = WWReel
	fg.WwSym = wwsym
	fg.Windetail = windetail
	if newFreeRound > 0 {
		fg.FreeRound += newFreeRound
		round.FreeSpin += newFreeRound
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
//	@param	retriggerTime 記錄進入迴圈次數並限制10次。【機率程式內自我呼叫用，server請填0】
//	@return round 回合資料
func Spin(rtp string, bet int, debugIndex [][]int, DebugSwitch bool, retriggerTime int) *Rounds {
	var reel []Reels
	var WWsym []int
	var position [][]int
	var fgNum int
	var slottype int
	fgwin := decimal.Zero

	if DebugSwitch {
		newreel, newWWsym, newposition, ok := DebugReels(debugIndex, rtp)
		position = newposition
		reel = newreel
		WWsym = newWWsym
		if !ok {
			ngreel, newpos := CreateReels(0, rtp)
			position = append(position, newpos)
			reel = append(reel, ngreel)
		}
	} else {
		slottype = 0
		if retriggerTime != 0 {
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
	if bonusCount >= fg_sym_def {
		round.FreeSpin += fg_round_time
		round.Result[0].Case = round.Result[0].Case.Push(FreeGame)
		round.Result[0].FreeRound += fg_round_time
		fgNum = fg_round_time
	}

	//free game------------------
	fgwin = decimal.Zero
	stage := 0
	for round.FreeSpin > 0 && len(position) > 1 {
		round.FreeSpin--
		stage++
		if stage == len(reel) {
			round.FreeSpin++
			stage--
			break
		}
		round = FGflow(stage, reel[stage], position[stage], round, WWsym)
		round.Result[stage].PreReel = reel[stage]
	}
	for round.FreeSpin > 0 {
		round.FreeSpin--
		stage++
		fgreel, newpos := CreateReels(1, rtp)
		round = FGflow(stage, fgreel, newpos, round, WWsym)
		round.Result[stage].PreReel = fgreel

	}
	fgwin = round.TotalPoint_Deci.Sub(round.Result[0].TotalPoint_Deci)

	if fgNum > 0 {
		if fgwin.LessThan(decimal.NewFromInt(10)) {
			switch {
			case retriggerTime < 10:
				retriggerTime++

				round = Spin(rtp, bet, debugIndex, DebugSwitch, retriggerTime)
			case retriggerTime == 10:
				break
			}
		}
	}

	//part of free game----------
	//// if len(round.Result) > 1 {
	//j_round, _ := json.Marshal(round)
	//fmt.Printf("\nRound: %s", string(j_round))
	// }

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

// 計算SB個數
func CountBouns(pt Reels) (bounsCount int) {
	for _, rowLine := range pt {
		for i, symbol := range rowLine {
			if i != 0 && i != 4 && symbol == int(SF) {
				bounsCount++
			}
		}
	}
	return
}

func GetNGReelsLen(rtp string) ([]int, error) {
	game_math := Gameplay
	rtps := RTPs(rtp)
	return game_math.ng_game.GetReelsLen(rtps), nil
}
func GetFGReelsLen(rtp string) ([]int, error) {
	game_math := Gameplay
	rtps := RTPs(rtp)
	return game_math.fg_game.GetReelsLen(rtps), nil
}

// Feature - FG中每Spin會抽出一個獎圖作為百搭
//
//	@param r 原盤面
//	@return Reels 新盤面
//	@return int 抽出獎圖
func RandomWW(rtp string, r Reels, wwsym int) Reels {
	newR := make(Reels, len(r))
	for i := 0; i < len(r); i++ { //第一輪不變百搭，從第二輪開始替換
		if i == 0 {
			newR[0] = r[0]
		}
		if i > 0 {
			for _, k := range r[i] {
				if k == wwsym {
					newR[i] = append(newR[i], WW.Int())
				} else {
					newR[i] = append(newR[i], k)
				}
			}
		}
	}
	return newR
}
