package slot_MegShark

import (
	"encoding/json"
	"fmt"
	"slotEditor/constant"
	"slotEditor/editor"
	server "slotserver/slot_MegShark"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/shopspring/decimal"
)

type WayData struct {
	ColNum         int      `json:"ColNum"`
	RowNum         int      `json:"RowNum"`
	PayLineNum     int      `json:"PayLineNum"`
	SymbolNum      int      `json:"SymbolNum"`
	SymbolNameList []string `json:"SymbolName"`
	SymbolPayout   [][]int  `json:"SymbolPayout"`
	SymbolWeight   [][]int  `json:"SymbolWeight"`
	PayLines       [][]int  `json:"PayLines"`
}

type WayEditor struct {
	*editor.BaseEditor
	DebugIndex     [][]int
	curRound       *Rounds
	curServerRound *server.Rounds
}

func NewGame(config editor.BaseConfig) *WayEditor {
	baseEditor := editor.NewBaseEditor(config)
	t := &WayEditor{
		BaseEditor: baseEditor,
	}
	t.Name = "slot4002"

	t.updateDebugTable(colNum+1, rowNum)
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

func (t *WayEditor) CreateAllSection() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("總覽", t.CreateOverallSection()),
		container.NewTabItem("測試", t.CreateDebugSection()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	return container.NewMax(tabs)
}

func (t *WayEditor) wayInitState() {
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

func (t *WayEditor) waySave() {
	data := WayData{
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

func (t *WayEditor) wayLoad() {
	win := *t.MainWindow
	t.EditorLoad(win)
}

func (t *WayEditor) wayUpdate(m []byte) {
	newWayData := WayData{}
	if err := json.Unmarshal(m, &newWayData); err != nil {
		fmt.Println(err)
		return
	}
	t.SetBasicInfo(newWayData.ColNum, newWayData.RowNum, newWayData.PayLineNum, newWayData.SymbolNum)
	symbolNameList = newWayData.SymbolNameList
	symbolPayout = newWayData.SymbolPayout
	PayLines = newWayData.PayLines
	t.Refresh()
}

func (t *WayEditor) wayServerRunOnce() {
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
	win *= bet
	t.FreeRound = len(curRound.Result) - 1
	if t.SlotType == SLOT_NORMAL {
		if t.FreeRound > 0 {
			t.StatBonusWin += win
		} else {
			t.StatNormalWin += win
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
	if curRound.Result[0].Point > 0 {
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

func (t *WayEditor) wayRunOnce() {
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
		//t.FreeTimes++
	}

	curRound := Spin(t.Rtp, int(bet), t.DebugIndex, t.DebugSwitch, 0)
	win, _ = decimal.NewFromFloat(curRound.TotalPoint).Round(2).Float64()
	t.FreeRound = len(curRound.Result) - 1
	if t.SlotType == SLOT_NORMAL {
		if t.FreeRound > 0 {
			t.StatBonusWin += win
		} else {
			t.StatNormalWin += win
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
	t.CurNG = curRound.Result[0].TotalPoint
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

func (t *WayEditor) wayGetColNum() int {
	return colNum
}

func (t *WayEditor) wayGetRowNum() int {
	return rowNum
}

func (t *WayEditor) wayGetPaylineNum() int {
	return payLineNum
}

func (t *WayEditor) wayGetSymbolNum() int {
	return symbolNum
}

func (t *WayEditor) waySetBasicInfo(cn int, rn int, pn int, sn int) {
	t.updateDebugTable(cn, rn)
	basicSetting(cn, rn, pn, sn)
}

func (t *WayEditor) wayGetSymbolName(index int) string {
	return symbolNameList[index]
}

func (t *WayEditor) wayGetSymbolNameList() []string {
	return symbolNameList
}

func (t *WayEditor) waySetSymbolName(index int, name string) {
	symbolNameList[index] = name
}

func (t *WayEditor) waySetSymbolNameList(nameList []string) {
	symbolNameList = nameList
}

func (t *WayEditor) wayGetSymbolTypeNameList() []string {
	return symbolTypeStr
}

func (t *WayEditor) wayGetSymbolType(index int) int {
	return symbolTypeMap[index]
}

func (t *WayEditor) wayDisplay() {
	str := storeSymbolTable(t.curRound)
	str += fmt.Sprintf("當前中獎金額:%.4f\n", t.CurWin)
	t.UpdateOverallInfo(str)
}

func (t *WayEditor) wayServerDisplay() {
	str := storeServerSymbolTable(t.curServerRound)
	str += fmt.Sprintf("當前中獎金額:%.4f\n", t.CurWin)
	t.UpdateOverallInfo(str)
}

func (t *WayEditor) wayGetSymbolPayout(symbol int, col int) int {
	return symbolPayout[symbol][col]
}

func (t *WayEditor) waySetSymbolPayout(symbol int, col int, num int) {
	symbolPayout[symbol][col] = num
}

func (t *WayEditor) waySetSymbolPayLine(payLinesNum int, col int, num int) {
	PayLines[payLinesNum][col] = num
}

func (t *WayEditor) wayDebugConfirm() {
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
		preReel := round.Result[i].PreReel

		if len(preReel) == 0 {
			return str
		}

		for j := 0; j < rowNum+2; j++ {
			for i := 0; i < colNum; i++ {
				symbol := preReel[i][j]
				tmpstr := symbolNameList[symbol] + " "
				str = str + tmpstr
			}
			str = str + "\n"
		}
		if i == 0 && round.Result[0].Point > 0 {
			str += "\n中獎內容:"
			for _, v := range round.Result[0].Windetail {
				str += fmt.Sprintf("\nSymbol: %+v 得到 %+v Way %+v 連線", v.Symbol, v.WinWay, v.Match)
			}
			str += "\n"
		}

		if i > 0 {
			wwReel := round.Result[i].WwReel
			str += "	||\n	V\n"
			str += fmt.Sprintf("抽中%+v當作WW\n", round.Result[i].WwSym)
			str += "	||\n	V\n"
			for j := 0; j < rowNum+2; j++ {
				for i := 0; i < colNum; i++ {
					Sym_2 := wwReel[i][j]
					tmpstr_2 := symbolNameList[Sym_2] + " "
					str += tmpstr_2
				}
				str = str + "\n"
			}
			if round.Result[i].Point > 0 {
				str += "中獎內容:"
				for _, v := range round.Result[i].Windetail {
					str += fmt.Sprintf("\nSymbol: %+v 得到 %+v Way %+v 連線", v.Symbol, v.WinWay, v.Match)
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
		preReel := round.Result[i].PreReel

		if len(preReel) == 0 {
			return str
		}

		for j := 0; j < MAX_ROW+2; j++ {
			for i := 0; i < MAX_COL; i++ {
				symbol := preReel[i][j]
				tmpstr := symbolNameList[symbol] + " "
				str = str + tmpstr
			}
			str = str + "\n"
		}
		if i == 0 && round.Result[0].Point > 0 {
			str += "\n中獎內容:"
			for _, v := range round.Result[0].Windetail {
				str += fmt.Sprintf("\nSymbol: %+v 得到 %+v Way %+v 連線", v.Symbol, v.WinWay, v.Match)
			}
			str += "\n"
		}

		if i > 0 {
			wwReel := round.Result[i].WwReel
			str += "	||\n	V\n"
			str += fmt.Sprintf("抽中%+v當作WW\n", round.Result[i].WwSym)
			str += "	||\n	V\n"
			for j := 0; j < MAX_ROW+2; j++ {
				for i := 0; i < MAX_COL; i++ {
					Sym_2 := wwReel[i][j]
					tmpstr_2 := symbolNameList[Sym_2] + " "
					str += tmpstr_2
				}
				str = str + "\n"
			}
			if round.Result[i].Point > 0 {
				str += "中獎內容:"
				for _, v := range round.Result[i].Windetail {
					str += fmt.Sprintf("\nSymbol: %+v 得到 %+v Way %+v 連線", v.Symbol, v.WinWay, v.Match)
				}

			}
			str += fmt.Sprintf("\n此局得分:%+v\n累積得分:%+v\n", round.Result[i].Point, round.Result[i].TotalPoint)
		}
		str += fmt.Sprintf("此局index:%+v\n", round.Result[i].ReelPosition)
	}
	str += "\n"
	return str
}

func (t *WayEditor) updateDebugTable(cn int, rn int) {
	//fmt.Printf("\nhit~~")
	t.DebugIndex = make([][]int, 11)
}
