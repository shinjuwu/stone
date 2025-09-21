package slot_Midas

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
					symbol = getWWmulti(rtp, slotType)
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
				currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1號
				if currentPosition >= ReelsLen[i] {
					currentPosition -= ReelsLen[i]
				}
				symbol := fgReelStrips[i][currentPosition]
				if symbol == WW { //換WW倍數
					symbol = getWWmulti(rtp, slotType)
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
			pt[i] = make([]int, rowNum+2)
			for j := 0; j < rowNum+2; j++ { //抽出不包含上下兩排:for j := 0; j < MAX_ROW; j++ {
				currentPosition := pos[i] + j //e.g.第一輪共10個獎圖，抽到9號時，應依序列出9、10、1號
				if currentPosition >= ReelsLen[i] {
					currentPosition -= ReelsLen[i]
				}
				symbol := ngReelStrips[i][currentPosition]
				if symbol == WW { //換WW倍數
					symbol = getWWmulti(rtp, slotType)
				}
				pt[i][j] = symbol.Int()
			}
		}
	}
	return
}
func DebugReels(debugIndex [][]int, rtp string) ([]Reels, [][]int, int, bool) {
	var ok bool
	var fgnum int
	ngReelsLen, _ := GetNGReelsLen(rtp) //每組滾輪長度
	fgReelsLen, _ := GetFGReelsLen(rtp) //每組滾輪長度
	debugpos := make([]Reels, len(debugIndex))
	for id, pos := range debugIndex {
		pt := make(Reels, colNum)
		if isDebug(pos[:5]) {
			ok = true
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
							symbol = getWWmulti(rtp, 0)
						}
						pt[i][j] = symbol.Int()
						//----
					}
					debugpos[id] = pt
				}
				scatterNum, _ := CountSF(pt)
				if len(pos) > 5 {
					if debug_isFG(pos[5], scatterNum) {
						fgnum = pos[5]
					}
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
							symbol = getWWmulti(rtp, 1)
						}
						pt[i][j] = symbol.Int()
					}
					debugpos[id] = pt
				}
			}
		} else if !isDebug(pos) {
			ok = false
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
	return debugpos, debugIndex, fgnum, ok
}

func isDebug(pos []int) (def bool) {
	def = true
	for _, v := range pos {
		if v < 0 || v > 70 {
			def = false
		}
	}
	return
}
func debug_isFG(fgtimes int, scatterNum int) bool {
	ok := false
	if scatterNum > 2 && scatterNum < 6 {
		if fgtimes >= scatterNum*3 && fgtimes <= scatterNum*9 {
			ok = true
		}
	}
	return ok
}

// 計算分數--LINEgame
func CalcNgWinLine(pt Reels, rtp string) ([]Windetail, decimal.Decimal, int, SfInfo, int) {
	var windetail []Windetail
	var win = decimal.Zero
	var isWild, getWild bool = false, false
	var WildMulti, freeround int
	var sfinfo SfInfo
	WildMulti = 1
	for id, payLine := range PayLines {
		isWild = false
		firstSymbol := pt[0][payLine[0]+1] //row=1,2,3
		curSymbol := firstSymbol
		count := 0
		for i := 1; i < colNum; i++ {
			curSymbol = pt[i][payLine[i]+1]
			if curSymbol == firstSymbol || isSymbolWild(curSymbol) {
				count++
				if isSymbolWild(curSymbol) {
					isWild = true
				}
			} else if firstSymbol != curSymbol && !isSymbolWild(curSymbol) {
				break
			}
		}
		if len(symbolPayout[firstSymbol]) > count {
			if symbolPayout[firstSymbol][count] > 0 {
				winpoint := decimal.NewFromInt(int64(symbolPayout[firstSymbol][count]))
				winfloat, _ := winpoint.Div(decimal.NewFromInt(int64(payLineNum))).Round(2).Float64()
				win = win.Add(winpoint)
				newWin := Windetail{firstSymbol, id, count + 1, winfloat}
				windetail = append(windetail, newWin)
				if isWild {
					getWild = true
				}
			}
		}
	}
	win = win.Div(decimal.NewFromInt(int64(payLineNum))) // 除上中獎線數量

	if getWild {
		WildMulti = countWildMulti(pt)
		win = win.Mul(decimal.NewFromInt(int64(WildMulti)))
	}

	//判斷是否觸發FG
	ScatterCount, sfreel := CountSF(pt) // 數Sf個數
	sfinfo.ScatterCount = ScatterCount
	if sfinfo.ScatterCount >= fg_sym_def {
		freeround, sfinfo.ScSmallReel = FgTimes(sfinfo.ScatterCount, rtp, sfreel)
	}

	return windetail, win, WildMulti, sfinfo, freeround
}

// 計算分數--LINEgame								    pt
func CalcFgWinLine(pt Reels, rtp string, wildMulti int) ([]Windetail, decimal.Decimal) {
	var windetail []Windetail
	var win = decimal.Zero

	for id, payLine := range PayLines {
		firstSymbol := pt[0][payLine[0]+1] //row=1,2,3
		curSymbol := firstSymbol
		count := 0
		for i := 1; i < colNum; i++ {
			curSymbol = pt[i][payLine[i]+1]
			if curSymbol == firstSymbol || isSymbolWild(curSymbol) {
				count++
			} else if firstSymbol != curSymbol && !isSymbolWild(curSymbol) {
				break
			}
		}
		if len(symbolPayout[firstSymbol]) > count {
			if symbolPayout[firstSymbol][count] > 0 {
				winpoint := decimal.NewFromInt(int64(symbolPayout[firstSymbol][count]))
				winfloat, _ := winpoint.Div(decimal.NewFromInt(int64(payLineNum))).Round(2).Float64()
				win = win.Add(winpoint)
				newWin := Windetail{firstSymbol, id, count + 1, winfloat}
				windetail = append(windetail, newWin)
			}
		}
	}

	win = win.Div(decimal.NewFromInt(int64(payLineNum))) // 除上中獎線數量
	win = win.Mul(decimal.NewFromInt(int64(wildMulti)))

	return windetail, win
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
	windetail, winpoint, wwmulti, sfinfo, freeround := CalcNgWinLine(pt, rtp)

	ng.Point_Deci = winpoint
	ng.Point, _ = ng.Point_Deci.Round(2).Float64()
	ng.Windetail = windetail
	ng.WWMulti = wwmulti
	ng.TotalPoint_Deci = ng.TotalPoint_Deci.Add(ng.Point_Deci)
	ng.TotalPoint, _ = ng.TotalPoint_Deci.Round(2).Float64()
	if winpoint.GreaterThan(decimal.Zero) {
		ng.Case = ng.Case.Push(Win)
	}
	if freeround > 0 {
		ng.ScatterInfo = sfinfo
		ng.Case = ng.Case.Push(FreeGame)
		ng.FreeRound = freeround
	}
	round.Result[0] = ng
	round.FreeSpin = freeround
	round.TotalPoint_Deci = round.TotalPoint_Deci.Add(ng.Point_Deci)
	round.TotalPoint, _ = round.TotalPoint_Deci.Round(2).Float64()
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
		Case:         Lose,
		ReelPosition: pos,
		Bet:          round.TotalBet,
	}
	if stage == 1 {
		fg.FgWildReel.Before = nil
	} else {
		fg.FgWildReel.Before = round.Result[stage-1].FgWildReel.After
	}
	pt, fg.FgWildReel.After = StickyWild(pt, fg.FgWildReel.Before, fg.FreeRound-stage)
	fg.WWMulti = countWildMulti(pt)
	fg.PayReel = nowPT
	windetail, point := CalcFgWinLine(pt, round.Rtp, fg.WWMulti)
	fg.Point_Deci = point
	fg.Point, _ = fg.Point_Deci.Round(2).Round(2).Float64()
	fg.Windetail = windetail

	if stage != 1 {
		fg.TotalPoint_Deci = round.Result[stage-1].TotalPoint_Deci.Add(fg.Point_Deci)
		fg.TotalPoint, _ = fg.TotalPoint_Deci.Round(2).Float64()
	} else if stage == 1 {
		fg.TotalPoint_Deci = point
		fg.TotalPoint, _ = fg.TotalPoint_Deci.Round(2).Float64()
	}
	if point.GreaterThan(decimal.Zero) {
		fg.Case = fg.Case.Push(Win)
	}
	round.TotalPoint_Deci = round.TotalPoint_Deci.Add(fg.Point_Deci)
	round.TotalPoint, _ = round.TotalPoint_Deci.Round(2).Float64()
	round.Result[stage] = fg
	return round
}

// Spin
//
//	@return round 回合資料
func Spin(rtp string, bet int, debugIndex [][]int, DebugSwitch bool, retrigger int) *Rounds {
	var reel []Reels
	var position [][]int
	var fgNum int
	var slottype int
	fgwin := decimal.Zero
	if DebugSwitch {
		newreel, newposition, debugfgtime, ok := DebugReels(debugIndex, rtp)
		position = newposition
		reel = newreel
		fgNum = debugfgtime
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
	if fgNum > 0 {
		round.FreeSpin = fgNum
		round.Result[0].FreeRound = fgNum
		round.Result[0].ScatterInfo.ScSmallReel = debug_SCReel(fgNum, round.Result[0].ScatterInfo.ScSmallReel)
	}

	//free game------------------
	fgNum = round.FreeSpin
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
		round = FGflow(stage, reel[stage], position[stage], round)
	}
	for round.FreeSpin > 0 {
		round.FreeSpin--
		stage++
		fgreel, newpos := CreateReels(1, rtp)
		round = FGflow(stage, fgreel, newpos, round)
	}
	fgwin = round.TotalPoint_Deci.Sub(round.Result[0].TotalPoint_Deci)

	if fgNum > 0 {
		if fgwin.LessThan(decimal.NewFromInt(10)) {
			round = Spin(rtp, bet, debugIndex, DebugSwitch, 3)
		}
	}
	//part of free game----------
	//	if len(round.Result) > 1 {
	//	j_round, _ := json.Marshal(round)
	//		fmt.Printf("\nRound: %s", string(j_round))
	//}
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
func CountSF(pt Reels) (sfCount int, sfreel [5]int) {
	for i, rowLine := range pt {
		if rowLine[1] == SF.Int() || rowLine[2] == SF.Int() || rowLine[3] == SF.Int() {
			sfCount++
			sfreel[i] = 1
		} else {
			sfreel[i] = 0
		}
	}
	return
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

func getWWmulti(rtp string, slottype int) (k Symbol) {
	if slottype == 0 {
		dice := random.Intn(WWMulti[rtp].Sum())
		pick, _ := WWMulti[rtp].Pick(dice)
		switch pick {
		case 2:
			k = WW_2
		case 3:
			k = WW_3
		default:
			k = WW
		}
	} else if slottype == 1 {
		dice := random.Intn(WWMulti_fg[rtp].Sum())
		pick, _ := WWMulti_fg[rtp].Pick(dice)
		switch pick {
		case 2:
			k = WW_2
		case 3:
			k = WW_3
		default:
			k = WW
		}
	}
	return
}

func countWildMulti(pt Reels) (WildMulti int) {
	WildMulti = 1
	for i := 0; i < colNum; i++ {
		for j := 0; j < rowNum; j++ {
			switch pt[i][j+1] {
			case 31:
				WildMulti += 1
			case 32:
				WildMulti += 2
			case 33:
				WildMulti += 3
			default:
			}
		}
	}
	return
}

func FgTimes(countSF int, rtp string, sfreel [5]int) (fgtimes int, ScSmallReel []int) {
	fgtimes = 0
	for i := 0; i < 3; i++ { //每個SF會有3個小滾輪
		for _, v := range sfreel {
			if v == 1 {
				dice := random.Intn(fgTimes[rtp].Sum())
				pick, _ := fgTimes[rtp].Pick(dice)
				fgtimes += pick
				ScSmallReel = append(ScSmallReel, pick)
			} else {
				ScSmallReel = append(ScSmallReel, 0)
			}
		}
	}
	return
}

func StickyWild(nowpt Reels, previousWild Reels, fg_remain int) (subpt Reels, afterWild Reels) {

	afterWild = AfterWW(nowpt)
	subpt = nowpt
	for i, v := range previousWild {
		for j, k := range v {
			if k != 0 {
				if fg_remain >= 0 {
					afterWild[i][j] = k
				}
				switch k {
				case 1:
					subpt[i][j] = Symbol(31).Int()
				case 2:
					subpt[i][j] = Symbol(32).Int()
				case 3:
					subpt[i][j] = Symbol(33).Int()
				default:
				}
			}
		}
	}
	//fmt.Printf("\nori: %+v\nnew:%+v\nbef:%+v\naft:%+v", original_pt, new_pt, previousWild, afterWild)
	return
}

func AfterWW(pt Reels) (WWReel Reels) {
	for i := 0; i < colNum; i++ {
		WWReel = append(WWReel, []int{0, 0, 0, 0, 0})
		for j := 0; j < rowNum+2; j++ {
			switch pt[i][j] {
			case 31:
				WWReel[i][j] = 1
			case 32:
				WWReel[i][j] = 2
			case 33:
				WWReel[i][j] = 3
			default:
				WWReel[i][j] = 0
			}
		}
	}
	return
}

func debug_SCReel(debug_fg int, scatterreel []int) []int {
	var sc []int
	for i := 0; i < 5; i++ {
		if scatterreel[i] > 0 {
			sc = append(sc, 1)
		} else {
			sc = append(sc, 0)
		}
	}
	debugSC := make([]int, 15)
	for m := 0; m < 3; m++ {
		for j := 0; j < 15; j++ {
			if debug_fg <= 0 {
				break
			}
			if sc[j%5] == 1 {
				debugSC[j]++
				debug_fg--
			}
		}
	}
	return debugSC
}
