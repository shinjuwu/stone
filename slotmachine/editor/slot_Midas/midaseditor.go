package slot_Midas

import (
	"encoding/json"
	"fmt"
	"slotEditor/constant"
	"slotEditor/editor"
	server "slotserver/slot_Midas"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/shopspring/decimal"
)

type LineData struct {
	ColNum         int      `json:"ColNum"`
	RowNum         int      `json:"RowNum"`
	PayLineNum     int      `json:"PayLineNum"`
	SymbolNum      int      `json:"SymbolNum"`
	SymbolNameList []string `json:"SymbolName"`
	SymbolPayout   [][]int  `json:"SymbolPayout"`
	SymbolWeight   [][]int  `json:"SymbolWeight"`
	PayLines       [][]int  `json:"PayLines"`
}

type MidasEditor struct {
	*editor.BaseEditor
	DebugIndex     [][]int
	curRound       *Rounds
	curServerRound *server.Rounds
}

func NewGame(config editor.BaseConfig) *MidasEditor {
	baseEditor := editor.NewBaseEditor(config)
	t := &MidasEditor{
		BaseEditor: baseEditor,
	}
	t.Name = "slot4003"

	t.updateDebugTable(colNum, rowNum)
	t.Unitbet = strconv.Itoa(unitbet)
	//method
	t.DebugConfirm = t.wayDebugConfirm
	t.InitState = t.wayInitState
	t.RunOnce = t.wayRunOnce
	t.ServerRunOnce = t.wayServerRunOnce
	t.Display = t.wayDisplay
	t.ServerDisplay = t.wayServerDisplay
	t.LoadFile = t.wayLoad
	t.SaveFile = t.waySave
	t.UpdateInfo = t.wayUpdate
	t.RecreateWindow = t.CreateAllSection
	t.SetBasicInfo = t.waySetBasicInfo

	// set get
	t.GetColNum = t.wayGetColNum
	t.GetRowNum = t.wayGetRowNum
	t.GetPayLineNum = t.wayGetPaylineNum
	t.GetSymbolNum = t.wayGetSymbolNum
	t.GetSymbolName = t.wayGetSymbolName
	t.SetSymbolName = t.waySetSymbolName
	t.GetSymbolNameList = t.wayGetSymbolNameList
	t.SetSymbolNameList = t.waySetSymbolNameList
	t.GetSymbolPayout = t.wayGetSymbolPayout
	t.SetSymbolPayout = t.waySetSymbolPayout
	t.SetSymbolPayLine = t.waySetSymbolPayLine
	t.GetSymbolTypeNameList = t.wayGetSymbolTypeNameList
	t.GetSymbolType = t.wayGetSymbolType

	return t
}

func (t *MidasEditor) CreateAllSection() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("總覽", t.CreateOverallSection()),
		container.NewTabItem("測試", t.CreateDebugSection()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	return container.NewMax(tabs)
}

func (t *MidasEditor) wayInitState() {
	t.Gold = constant.DEFAULT_GOLD
	t.Round = 0
	t.FreeRound = 0
	t.FreeTimes = 0
	t.SlotType = SLOT_NORMAL
	t.StatBet = 0
	t.StatWin = 0
	t.StatNormalWin = 0
	t.StatBonusWin = 0
	t.LabelGold.SetText(fmt.Sprintf("%.2f", t.Gold))
	t.LabelCurRound.SetText(fmt.Sprintf("%d", t.Round))
	t.LabelFreeRound.SetText(fmt.Sprintf("%d", t.FreeRound))
	t.EntryDisplay.SetText("")
	t.LabelStatBet.SetText("0")
	t.LabelStatWin.SetText("0")
	t.LabelRTP.SetText("0")
	t.LabelStatNormalWin.SetText("0")
	t.LabelStatBonusWin.SetText("0")
	t.LabelTestFgWin.SetText("0")
	t.LabelTestNgWin.SetText("0")
	t.NgGetWin = 0
	t.FgGetWin = 0
}

func (t *MidasEditor) waySave() {
	data := LineData{
		ColNum:         colNum,
		RowNum:         rowNum,
		PayLineNum:     payLineNum,
		SymbolNum:      symbolNum,
		SymbolNameList: symbolNameList,
		SymbolPayout:   symbolPayout,
		PayLines:       PayLines,
	}
	m, _ := json.Marshal(data)
	win := *t.MainWindow
	t.EditorSave(m, win)
}

func (t *MidasEditor) wayLoad() {
	win := *t.MainWindow
	t.EditorLoad(win)
}

func (t *MidasEditor) wayUpdate(m []byte) {
	newLineData := LineData{}
	if err := json.Unmarshal(m, &newLineData); err != nil {
		fmt.Println(err)
		return
	}
	t.SetBasicInfo(newLineData.ColNum, newLineData.RowNum, newLineData.PayLineNum, newLineData.SymbolNum)
	symbolNameList = newLineData.SymbolNameList
	symbolPayout = newLineData.SymbolPayout
	PayLines = newLineData.PayLines
	t.Refresh()
}

func (t *MidasEditor) wayServerRunOnce() {
	win := 0.0
	bet := 1.0
	t.CurBG = 0
	t.CurNG = 0
	t.Round++
	t.StatBet += bet
	if t.SlotType == server.SLOT_NORMAL {
		t.Gold -= bet
	} else if t.SlotType == server.SLOT_BONUS_FREE {
		t.FreeRound--
		t.FreeTimes++
	}

	curRound := server.Spin(t.Rtp, int(bet), t.DebugIndex, t.DebugSwitch, 0)
	win, _ = decimal.NewFromFloat(curRound.TotalPoint).Round(2).Float64()
	win_ng, _ := decimal.NewFromFloat(curRound.Result[0].Point).Float64()
	t.FreeRound = len(curRound.Result) - 1
	if t.SlotType == SLOT_NORMAL {
		if t.FreeRound > 0 {
			t.StatBonusWin += win - win_ng
		} else {
			t.StatNormalWin += win_ng
		}
	} else if t.SlotType == SLOT_BONUS_FREE {
		t.StatBonusWin += win
	}

	t.Gold += win
	t.StatWin += win
	t.CurWin = win

	if t.FreeRound == 0 {
		t.SlotType = SLOT_NORMAL
	} else if t.FreeRound > 0 {
		t.SlotType = SLOT_BONUS_FREE
	} else {
		t.SlotType = SLOT_NORMAL
		t.FreeRound = 0
	}
	//
	if curRound.TotalPoint > t.MaxScore {
		t.MaxScore = curRound.TotalPoint
	}
	if win_ng > 0 {
		t.NgGetWin++
	}
	t.CurNG = curRound.Result[0].TotalPoint
	if len(curRound.Result) > 1 {
		t.CurBG = curRound.TotalPoint - curRound.Result[0].TotalPoint
		for i := 1; i < len(curRound.Result); i++ {
			t.TotalFgRound++
			if curRound.Result[i].Point > 0 {
				t.FgGetWin++
			}
		}
	}
	//
	t.curServerRound = curRound
}

func (t *MidasEditor) wayRunOnce() {
	win := 0.0
	bet := 1.0
	t.CurBG = 0
	t.CurNG = 0
	t.Round++
	t.StatBet += bet
	if t.SlotType == SLOT_NORMAL {
		t.Gold -= bet
	} else if t.SlotType == SLOT_BONUS_FREE {
		t.FreeRound--
	}

	curRound := Spin(t.Rtp, int(bet), t.DebugIndex, t.DebugSwitch, 0)
	win, _ = decimal.NewFromFloat(curRound.TotalPoint).Round(2).Float64()
	win_ng, _ := decimal.NewFromFloat(curRound.Result[0].Point).Float64()
	t.FreeRound = len(curRound.Result) - 1
	if t.SlotType == SLOT_NORMAL {
		if t.FreeRound > 0 {
			t.StatBonusWin += win - win_ng
		} else {
			t.StatNormalWin += win_ng
		}
	} else if t.SlotType == SLOT_BONUS_FREE {
		t.StatBonusWin += win
	}

	t.Gold += win
	t.StatWin += win
	t.CurWin = win

	if t.FreeRound > 0 {
		t.SlotType = SLOT_BONUS_FREE
	} else {
		t.SlotType = SLOT_NORMAL
		t.FreeRound = 0
	}
	//
	if curRound.TotalPoint > t.MaxScore {
		t.MaxScore = curRound.TotalPoint
	}
	if curRound.Result[0].Point > 0 {
		t.NgGetWin++
	}
	t.CurNG = win_ng
	if len(curRound.Result) > 1 {
		t.FreeTimes++
		t.CurBG = curRound.TotalPoint - curRound.Result[0].TotalPoint
		for i := 1; i < len(curRound.Result); i++ {
			t.TotalFgRound++
			if curRound.Result[i].Point > 0 {
				t.FgGetWin++
			}
		}
	}
	//
	t.curRound = curRound
}

func (t *MidasEditor) wayGetColNum() int {
	return colNum
}

func (t *MidasEditor) wayGetRowNum() int {
	return rowNum
}

func (t *MidasEditor) wayGetPaylineNum() int {
	return payLineNum
}

func (t *MidasEditor) wayGetSymbolNum() int {
	return symbolNum
}

func (t *MidasEditor) waySetBasicInfo(cn int, rn int, pn int, sn int) {
	t.updateDebugTable(cn, rn)
	basicSetting(cn, rn, pn, sn)
}

func (t *MidasEditor) wayGetSymbolName(index int) string {
	return symbolNameList[index]
}

func (t *MidasEditor) wayGetSymbolNameList() []string {
	return symbolNameList
}

func (t *MidasEditor) waySetSymbolName(index int, name string) {
	symbolNameList[index] = name
}

func (t *MidasEditor) waySetSymbolNameList(nameList []string) {
	symbolNameList = nameList
}

func (t *MidasEditor) wayGetSymbolTypeNameList() []string {
	return symbolTypeStr
}

func (t *MidasEditor) wayGetSymbolType(index int) int {
	return symbolTypeMap[index]
}

func (t *MidasEditor) wayDisplay() {
	str := storeSymbolTable(t.curRound)
	str += fmt.Sprintf("當前中獎金額:%.4f\n", t.CurWin)
	t.UpdateOverallInfo(str)
}

func (t *MidasEditor) wayServerDisplay() {
	str := storeServerSymbolTable(t.curServerRound)
	str += fmt.Sprintf("當前中獎金額:%.4f\n", t.CurWin)
	t.UpdateOverallInfo(str)
}

func (t *MidasEditor) wayGetSymbolPayout(symbol int, col int) int {
	return symbolPayout[symbol][col]
}

func (t *MidasEditor) waySetSymbolPayout(symbol int, col int, num int) {
	symbolPayout[symbol][col] = num
}

func (t *MidasEditor) waySetSymbolPayLine(payLinesNum int, col int, num int) {
	PayLines[payLinesNum][col] = num
}

func (t *MidasEditor) wayDebugConfirm() {
	t.DebugCheckBox.SetChecked(false)
	t.DebugSwitch = false
	for j := range t.DebugIndex {
		t.DebugIndex[j] = []int{}
		for i := 0; i < colNum+1; i++ {
			text := t.DebugSymbolSelect[j][i].Text
			new, _ := strconv.Atoi(text)
			t.DebugIndex[j] = append(t.DebugIndex[j], new)
		}
	}
}

func storeSymbolTable(round *Rounds) string {
	str := ""
	// 把獎圖轉置輸出
	for i := 0; i < len(round.Result); i++ {
		remain := round.Result[i].FreeRound - int(round.Result[i].Stages)
		if len(round.Result) > 1 {
			if i == 0 {
				str += fmt.Sprint("NG\n")
			} else {
				str += fmt.Sprintf("\nFG第%d局:%d/%d\n", i, remain, round.Result[i].FreeRound)
			}
		} else {
			str += fmt.Sprint("NG\n")
		}
		payReel := round.Result[i].PayReel
		if len(payReel) == 0 {
			return str
		}
		ReelSym := PrintSym(payReel, *reelDef)
		for _, v := range ReelSym {
			for _, k := range v {
				str += k + " 	"
			}
			str += "\n"
		}
		if i == 0 && round.Result[0].Point > 0 {
			str += "\n中獎內容:"
			str += fmt.Sprintf(" 當前倍數%d", round.Result[i].WWMulti)
			for _, v := range round.Result[0].Windetail {
				str += fmt.Sprintf("\nSymbol: %+v 得到 第 %+v line %+v 連線 賠率 %+v", v.Symbol, v.WinLine, v.Match, v.Multi)
			}
			str += "\n"

		}
		if i == 0 && round.Result[0].ScatterInfo.ScatterCount >= 3 {
			str += fmt.Sprintf("\nsfreel:\n%+v\n%+v\n%+v\nget %d freespin\n", round.Result[0].ScatterInfo.ScSmallReel[:5], round.Result[0].ScatterInfo.ScSmallReel[5:10], round.Result[0].ScatterInfo.ScSmallReel[10:], round.Result[0].FreeRound)
		}
		if i > 0 {
			var wwReel [5][5]int
			for col, v := range round.Result[i].FgWildReel.Before {
				for row, k := range v {
					wwReel[col][row] = k
				}
			}
			str += "	||\n	V\n"
			str += fmt.Sprintf("當前倍數%d\n\n停在盤面上的WW:\n", round.Result[i].WWMulti)
			for j := 0; j < rowNum+2; j++ {
				for i := 0; i < colNum; i++ {
					Sym_2 := wwReel[i][j]
					tmpstr_2 := strconv.Itoa(Sym_2) + " "
					str += tmpstr_2
				}
				str = str + "\n"
			}
			if round.Result[i].Point > 0 {
				str += "\n中獎內容:"
				for _, v := range round.Result[i].Windetail {
					str += fmt.Sprintf("\nSymbol:"+"%+v"+" 得到第 %+v 線的 %+v 連線 賠率 %+v", v.Symbol, v.WinLine, v.Match, v.Multi)
				}

			}
			str += fmt.Sprintf("\n此局得分:%+v\n累積得分:%+v\n", round.Result[i].Point, round.Result[i].TotalPoint)
		}
		str += fmt.Sprintf("此局index:%+v\n", round.Result[i].ReelPosition)
	}
	str += "\n"
	return str
}

func storeServerSymbolTable(round *server.Rounds) string {
	str := ""
	// 把獎圖轉置輸出
	for i := 0; i < len(round.Result); i++ {
		remain := round.Result[i].FreeRound - int(round.Result[i].Stages)
		if len(round.Result) > 1 {
			if i == 0 {
				str += fmt.Sprint("NG\n")
			} else {
				str += fmt.Sprintf("\nFG第%d局:%d/%d\n", i, remain, round.Result[i].FreeRound)
			}
		} else {
			str += fmt.Sprint("NG\n")
		}
		payReel := round.Result[i].PayReel

		if len(payReel) == 0 {
			return str
		}

		for j := 0; j < MAX_ROW+2; j++ {
			for i := 0; i < MAX_COL; i++ {
				symbol := payReel[i][j]
				tmpstr := symbolNameList[symbol] + " "
				str = str + tmpstr
			}
			str = str + "\n"
		}
		if i == 0 && round.Result[0].Point > 0 {
			str += "\n中獎內容:"
			str += fmt.Sprintf(" 當前倍數%d", round.Result[i].WWMulti)
			for _, v := range round.Result[0].Windetail {
				str += fmt.Sprintf("\nSymbol: %+v 得到 %+v Line %+v 連線  賠率 %+v", v.Symbol, v.WinLine, v.Match, v.Multi)
			}
			str += "\n"

		}
		if i == 0 && round.Result[0].ScatterInfo.ScatterCount >= 3 {
			str += fmt.Sprintf("\nsfreel:\n%+v\n%+v\n%+v\nget %d freespin\n", round.Result[0].ScatterInfo.ScSmallReel[:5], round.Result[0].ScatterInfo.ScSmallReel[5:10], round.Result[0].ScatterInfo.ScSmallReel[10:], round.Result[0].FreeRound)
		}
		if i > 0 {
			var wwReel [5][5]int
			for col, v := range round.Result[i].FgWildReel.Before {
				for row, k := range v {
					wwReel[col][row] = k
				}
			}
			str += "	||\n	V\n"
			str += fmt.Sprintf("當前倍數%d\n\n停在盤面上的WW:\n", round.Result[i].WWMulti)
			for j := 0; j < rowNum+2; j++ {
				for i := 0; i < colNum; i++ {
					Sym_2 := wwReel[i][j]
					tmpstr_2 := strconv.Itoa(Sym_2) + " "
					str += tmpstr_2
				}
				str = str + "\n"
			}
			if round.Result[i].Point > 0 {
				str += "\n中獎內容:"
				for _, v := range round.Result[i].Windetail {
					str += fmt.Sprintf("\nSymbol:"+"%+v"+" 得到第 %+v 線的 %+v 連線 賠率 %+v", v.Symbol, v.WinLine, v.Match, v.Multi)
				}

			}
			str += fmt.Sprintf("\n此局得分:%+v\n累積得分:%+v\n", round.Result[i].Point, round.Result[i].TotalPoint)
		}
		str += fmt.Sprintf("此局index:%+v\n", round.Result[i].ReelPosition)
	}
	str += "\n"
	return str
}

func (t *MidasEditor) updateDebugTable(cn int, rn int) {
	t.DebugIndex = make([][]int, 11)
	for i := 0; i < 11; i++ {
		t.DebugIndex[i] = make([]int, cn+1)
	}
}

func PrintSym(pt Reels, reelDef ReelStripsDef) [][]string {
	CountRow := 0
	for _, v := range reelDef {
		if v > CountRow {
			CountRow = v
		}
	}
	CountRow += 2
	var SymReel [][]string
	for i := 0; i < CountRow; i++ {
		rowsym := []string{}
		for j := 0; j < len(reelDef); j++ {
			var sym string
			if i <= len(pt[j])-1 {
				sym = symbolNameList[pt[j][i]]
			} else {
				sym = "   "
			}
			rowsym = append(rowsym, sym)
		}
		SymReel = append(SymReel, rowsym)
	}
	return SymReel
}
