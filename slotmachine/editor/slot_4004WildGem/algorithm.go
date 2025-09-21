package slot_4004WildGem

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
	var ReelsLen []int //每組滾輪長度
	switch slotType {
	case 0: //NG滾輪
		ReelsLen, _ = GetNGReelsLen(rtp)
		pos = random.Intsn(ReelsLen)
		for i, v := range *reelDef {
			pt[i] = make([]int, v+2)
			for j := 0; j < v+2; j++ { //抽出不包含上下兩排:for j := 0; j < rowNum; j++ {
				currentPosition := pos[i] + j //e.g.第一輪0-9共10個獎圖，抽到9號時，應依序列出9、0、1、2、3號
				if currentPosition >= ReelsLen[i] {
					currentPosition -= ReelsLen[i]
				}
				//抽取WW倍數並替換
				symbol := ngReelStrips[i][currentPosition]
				if symbol == WW { //換WW倍數
					symbol = getWWPrize(rtp, slotType, i)
				}
				pt[i][j] = symbol.Int()
				//----
			}
		}
	case 1: //FG滾輪
		ReelsLen, _ = GetFGReelsLen(rtp)
		pos = random.Intsn(ReelsLen)
		for i := 0; i < colNum; i++ {
			pt[i] = make([]int, rowNum+2)
			for j := 0; j < rowNum+2; j++ { //抽出不包含上下兩排:for j := 0; j < rowNum; j++ {
				currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1、2、3號
				if currentPosition >= ReelsLen[i] {
					currentPosition -= ReelsLen[i]
				}
				symbol := fgReelStrips[i][currentPosition]
				if symbol == WW { //換WW倍數
					symbol = getWWPrize(rtp, slotType, i)
				}
				pt[i][j] = symbol.Int()
			}
		}
	case 3: //NG必中FG滾輪
		ReelsLen, _ = GetNGReelsLen(rtp)
		dice := random.Intn(Rand_SFcount.Sum())
		pick, _ := Rand_SFcount.Pick(dice)
		SFpos := GetSFpos(pick)
		pos = GetFGpos(SFpos)
		for i := 0; i < MAX_COL; i++ {
			pt[i] = make([]int, MAX_ROW+2)
			for j := 0; j < MAX_ROW+2; j++ { //抽出不包含上下兩排:for j := 0; j < MAX_ROW; j++ {
				currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1號
				if currentPosition >= ReelsLen[i] {
					currentPosition -= ReelsLen[i]
				}
				symbol := ngReelStrips[i][currentPosition]
				if symbol == WW { //換WW倍數
					symbol = getWWPrize(rtp, 0, i)
				}
				pt[i][j] = symbol.Int()
			}
		}
	}
	return
}
func DebugReels(debugIndex [][]int, rtp string) ([]Reels, [][]int, bool) {
	var ok bool
	ngReelsLen, _ := GetNGReelsLen(rtp) //每組滾輪長度
	fgReelsLen, _ := GetFGReelsLen(rtp) //每組滾輪長度
	debugpos := make([]Reels, len(debugIndex))
	for id, pos := range debugIndex {
		pt := make(Reels, colNum)
		if isDebug(pos) {
			if id == 0 {
				for i := 0; i < colNum; i++ {
					pt[i] = make([]int, rowNum+2)
					for j := 0; j < rowNum+2; j++ {
						currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1號
						if currentPosition >= ngReelsLen[i] {
							currentPosition -= ngReelsLen[i]
						}
						//抽取WW倍數並替換
						symbol := ngReelStrips[i][currentPosition]
						if symbol == WW { //換WW倍數
							symbol = getWWPrize(rtp, 0, i)
						}
						pt[i][j] = symbol.Int()
						//----
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
						if symbol == WW { //換WW倍數
							symbol = getWWPrize(rtp, 1, i)
						}
						pt[i][j] = symbol.Int()
					}
					debugpos[id] = pt
				}
			}
		} else if !isDebug(pos) {
			if id == 0 {
				pt, pos = CreateReels(0, rtp)
				debugpos[id] = pt
				debugIndex[id] = pos
			} else if id != 0 {
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
		if v < 0 || v > 75 {
			def = false
		}
	}
	return
}

// 計算分數--LINEgame							               pt       wildprize          freeround
func CalcNgWinLine(pt Reels, rtp string) ([]Windetail, decimal.Decimal, []float64, int) {
	var windetail []Windetail
	var win = decimal.Zero
	var freeround int
	var isgetww [5][5]int
	//var wildpos [][]int
	for id, payLine := range PayLines {
		firstSymbol := pt[0][payLine[0]+1] //row=1,2,3
		curSymbol := firstSymbol
		count := 0
		for i := 1; i < colNum; i++ {
			curSymbol = pt[i][payLine[i]+1]
			if curSymbol == firstSymbol || isSymbolWild(curSymbol) {
				count++
				if isSymbolWild(curSymbol) {
					isgetww[i][payLine[i]+1] = 1
					//wildpos = append(wildpos, []int{i, payLine[i] + 1})
				}
			} else if firstSymbol != curSymbol && !isSymbolWild(curSymbol) {
				break
			}
		}
		if len(symbolPayout[firstSymbol]) > count {
			if symbolPayout[firstSymbol][count] > 0 {
				winpoint := decimal.NewFromInt(int64(symbolPayout[firstSymbol][count]))
				winfloat, _ := winpoint.Div(decimal.NewFromInt(int64(payLineNum))).Float64()
				win = win.Add(winpoint)
				newWin := Windetail{firstSymbol, id, count + 1, winfloat}
				windetail = append(windetail, newWin)
			}
		}
	}
	win = win.Div(decimal.NewFromInt(int64(payLineNum))) // 除上中獎線數量

	//判斷是否觸發FG
	SFCount := CountSF(pt) // 數Sf個數
	if SFCount >= fg_sym_def {
		freeround = FGtimes[SFCount]
	}

	wildprize, ww_pt := accWildPrize(isgetww, pt)
	win = win.Add(ww_pt) //加上WW_prize
	return windetail, win, wildprize, freeround
}

// 計算分數--LINEgame								                        pt
func CalcFgWinLine(pt Reels, rtp string) ([]Windetail, decimal.Decimal, []float64) {
	var windetail []Windetail
	var win = decimal.Zero
	//var wildpos [][]int
	var isgetww [5][5]int
	for id, payLine := range PayLines {
		firstSymbol := pt[0][payLine[0]+1] //row=1,2,3
		curSymbol := firstSymbol
		count := 0
		for i := 1; i < colNum; i++ {
			curSymbol = pt[i][payLine[i]+1]
			if curSymbol == firstSymbol || isSymbolWild(curSymbol) {
				count++
				if isSymbolWild(curSymbol) {
					isgetww[i][payLine[i]+1] = 1
				}
			} else if firstSymbol != curSymbol && !isSymbolWild(curSymbol) {
				break
			}
		}
		if len(symbolPayout[firstSymbol]) > count {
			if symbolPayout[firstSymbol][count] > 0 {
				winpoint := decimal.NewFromInt(int64(symbolPayout[firstSymbol][count]))
				winfloat, _ := winpoint.Div(decimal.NewFromInt(int64(payLineNum))).Float64()
				win = win.Add(winpoint)
				newWin := Windetail{firstSymbol, id, count + 1, winfloat}
				windetail = append(windetail, newWin)
			}
		}
	}

	win = win.Div(decimal.NewFromInt(int64(payLineNum))) // 除上中獎線數量

	//ww_prize
	wildprize, ww_pt := accWildPrize(isgetww, pt)
	win = win.Add(ww_pt) //加上WW_prize

	return windetail, win, wildprize
}

func isSymbolWild(symbol int) bool {
	return symbolTypeMap[symbol] == SYMBOL_TYPE_WILD
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
	}
	ng.Bet = round.TotalBet
	windetail, winpoint, ww_prize, freeround := CalcNgWinLine(pt, rtp)
	ng.Point, _ = winpoint.Mul(decimal.NewFromFloat(ng.Bet)).Round(4).Float64()
	ng.Windetail = windetail
	ng.WildPrize = ww_prize

	ng.Point_Deci = decimal.NewFromFloat(ng.Point)
	if winpoint.GreaterThan(decimal.Zero) {
		ng.Case = ng.Case.Push(Win)
	}
	if freeround > 0 {
		ng.Case = ng.Case.Push(FreeGame)
		ng.FreeRound = freeround
	}
	round.Result[0] = ng
	round.FreeSpin = freeround
	round.TotalPoint, _ = decimal.NewFromFloat(round.TotalPoint).Add(decimal.NewFromFloat(ng.Point)).Round(4).Float64()
	round.TotalPoint_Deci = decimal.NewFromFloat(round.TotalPoint)
	return round
}

func FGflow(stage int, pt Reels, pos []int, round *Rounds) *Rounds {
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
		FreeRound:    round.Result[0].FreeRound,
		Retrigger:    0,
		Case:         Lose,
		ReelPosition: pos,
		Bet:          round.TotalBet,
		TotalPoint:   round.TotalPoint,
	}

	fg.PayReel = nowPT
	windetail, point, ww_prize := CalcFgWinLine(pt, round.Rtp)
	fg.Point, _ = point.Mul(decimal.NewFromFloat(fg.Bet)).Float64()
	fg.Point_Deci = decimal.NewFromFloat(fg.Point).Round(4)
	fg.Windetail = windetail
	fg.WildPrize = ww_prize

	if point.GreaterThan(decimal.Zero) {

		fg.Case = fg.Case.Push(Win)
	}
	round.TotalPoint, _ = decimal.NewFromFloat(round.TotalPoint).Add(decimal.NewFromFloat(fg.Point)).Float64()
	round.TotalPoint_Deci = decimal.NewFromFloat(round.TotalPoint)
	fg.TotalPoint, _ = decimal.NewFromFloat(fg.TotalPoint).Add(point).Float64()

	SFCount := CountSF(fg.PayReel)
	if SFCount >= fg_sym_def {
		fg.Retrigger = FGtimes[SFCount]
		round.FreeSpin += fg.Retrigger
		fg.FreeRound += fg.Retrigger
		fg.Case = fg.Case.Push(FreeGame)
	}
	round.Result[stage] = fg
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
	fgnum := round.FreeSpin
	//free game------------------
	stage := 0
	//fgwin = decimal.Zero
	for round.FreeSpin > 0 && len(position) > 1 {
		round.FreeSpin--
		stage++
		if stage == len(reel) {
			round.FreeSpin++
			stage--
			break
		}
		round = FGflow(stage, reel[stage], position[stage], round)
	}
	for round.FreeSpin > 0 {
		round.FreeSpin--
		stage++
		fgreel, newpos := CreateReels(1, rtp)
		round = FGflow(stage, fgreel, newpos, round)
	}
	fgwin = round.TotalPoint_Deci.Sub(round.Result[0].TotalPoint_Deci)
	if fgnum > 0 {
		if fgwin.LessThan(decimal.NewFromInt(5)) {
			round = Spin(rtp, bet, debugIndex, DebugSwitch, 3)
		}
	}
	//section of free game ----------

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
	for i := range r {
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
	for _, v := range col {
		if v == sym {
			count++
		}
	}
	return
}

// 計算SF個數
func CountSF(pt Reels) (sfCount int) {
	for _, rowLine := range pt {
		if rowLine[1] == SF.Int() || rowLine[2] == SF.Int() || rowLine[3] == SF.Int() {
			sfCount++
		}
	}
	return
}
func accWildPrize(isgetww [5][5]int, pt Reels) ([]float64, decimal.Decimal) { //wildpos [][]int
	accwildprize := decimal.Zero
	var wildPrize []float64
	// for _, v := range wildpos {
	// 	sym := pt[v[0]][v[1]]
	// 	prize := WW_number[sym]
	// 	accwildprize = accwildprize.Add(prize)
	// 	prizefloat, _ := prize.Float64()
	// 	wildPrize = append(wildPrize, prizefloat)
	// }
	for i, v := range isgetww {
		for j, k := range v {
			if k == 1 {
				sym := pt[i][j]
				prize := WW_number[sym]
				accwildprize = accwildprize.Add(prize)
				prizefloat, _ := prize.Float64()
				wildPrize = append(wildPrize, prizefloat)
			}
		}
	}
	return wildPrize, accwildprize
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

func getWWPrize(rtp string, slottype int, reelnum int) (k Symbol) {
	if slottype == 0 {
		dice := random.Intn(NGwild[rtp][reelnum].Sum())
		pick, _ := NGwild[rtp][reelnum].Pick(dice)
		k = Symbol(pick)
	} else if slottype == 1 {
		dice := random.Intn(FGwild[rtp][reelnum].Sum())
		pick, _ := FGwild[rtp][reelnum].Pick(dice)
		k = Symbol(pick)
	}
	return
}
