package slot_Next3

import (
	"encoding/json"
	"fmt"
	"slotEditor/constant"
	"slotEditor/editor"
	server "slotserver/slot_Next3"
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

type Next3Editor struct {
	*editor.BaseEditor
	DebugIndex     [][]int
	curRound       *Rounds
	curServerRound *server.Rounds
}

func NewGame(config editor.BaseConfig) *Next3Editor {
	baseEditor := editor.NewBaseEditor(config)
	t := &Next3Editor{
		BaseEditor: baseEditor,
	}
	t.Name = "slot4004"

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

func (t *Next3Editor) CreateAllSection() fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("總覽", t.CreateOverallSection()),
		container.NewTabItem("測試", t.CreateDebugSection()),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	return container.NewMax(tabs)
}

func (t *Next3Editor) wayInitState() {
	t.Gold = constant.DEFAULT_GOLD
	t.Round = 0
	t.FreeRound = 0
	t.FreeTimes = 0
	t.SlotType = SLOT_NORMAL
	t.StatBet = 0
	t.StatWin = 0
	t.StatNormalWin = 0
	t.StatBonusWin = 0
	t.Stage = 0
	t.LabelGold.SetText(fmt.Sprintf("%.2f", t.Gold))
	t.LabelCurRound.SetText(fmt.Sprintf("%d", t.Round))
	t.LabelFreeRound.SetText(fmt.Sprintf("%d", t.FreeRound))
	t.EntryDisplay.SetText("")
	t.LabelStatBet.SetText("0")
	t.LabelStatWin.SetText("0")
	t.LabelRTP.SetText("0")
	t.LabelStatNormalWin.SetText("0")
	t.LabelStatBonusWin.SetText("0")

}

func (t *Next3Editor) waySave() {
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

func (t *Next3Editor) wayLoad() {
	win := *t.MainWindow
	t.EditorLoad(win)
}

func (t *Next3Editor) wayUpdate(m []byte) {
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

func (t *Next3Editor) wayServerRunOnce() {
	win := 0.0
	bet := 1.0

	if t.SlotType == SLOT_NORMAL {
		t.StatBet += bet
		t.Gold -= bet
		t.Round++
	} else if t.SlotType == SLOT_BONUS_FREE {
		t.FreeRound--
	}

	curRound := server.Spin(t.Rtp, int(bet), t.DebugIndex, t.DebugSwitch)
	if len(curRound.Result) > 1 {
		t.FreeTimes++
	}
	win, _ = decimal.NewFromFloat(curRound.TotalPoint).Round(2).Float64()
	win_ng, _ := decimal.NewFromFloat(curRound.Result[0].Point).Round(2).Float64()
	win *= bet
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
	t.curServerRound = curRound
}

func (t *Next3Editor) wayRunOnce() {
	win := 0.0
	bet := 1.0
	if t.SlotType == SLOT_NORMAL {
		t.StatBet += bet
		t.Gold -= bet
		t.Round++
	} else if t.SlotType == SLOT_BONUS_FREE {
		t.FreeRound--
	}

	curRound := Spin(t.Rtp, int(bet), t.DebugIndex, t.DebugSwitch)
	if len(curRound.Result) > 1 {
		t.FreeTimes++
	}
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

	if t.FreeRound > 0 {
		t.SlotType = SLOT_BONUS_FREE
	} else {
		t.SlotType = SLOT_NORMAL
		t.FreeRound = 0
	}
	t.curRound = curRound
}

func (t *Next3Editor) wayGetColNum() int {
	return colNum
}

func (t *Next3Editor) wayGetRowNum() int {
	return rowNum
}

func (t *Next3Editor) wayGetPaylineNum() int {
	return payLineNum
}

func (t *Next3Editor) wayGetSymbolNum() int {
	return symbolNum
}

func (t *Next3Editor) waySetBasicInfo(cn int, rn int, pn int, sn int) {
	t.updateDebugTable(cn, rn)
	basicSetting(cn, rn, pn, sn)
}

func (t *Next3Editor) wayGetSymbolName(index int) string {
	return symbolNameList[index]
}

func (t *Next3Editor) wayGetSymbolNameList() []string {
	return symbolNameList
}

func (t *Next3Editor) waySetSymbolName(index int, name string) {
	symbolNameList[index] = name
}

func (t *Next3Editor) waySetSymbolNameList(nameList []string) {
	symbolNameList = nameList
}

func (t *Next3Editor) wayGetSymbolTypeNameList() []string {
	return symbolTypeStr
}

func (t *Next3Editor) wayGetSymbolType(index int) int {
	return symbolTypeMap[index]
}

func (t *Next3Editor) wayDisplay() {
	str := storeSymbolTable(t.curRound)
	str += fmt.Sprintf("當前中獎金額:%.4f\n", t.CurWin)
	t.UpdateOverallInfo(str)
}

func (t *Next3Editor) wayServerDisplay() {
	str := storeServerSymbolTable(t.curServerRound)
	t.UpdateOverallInfo(str)
}

func (t *Next3Editor) wayGetSymbolPayout(symbol int, col int) int {
	return symbolPayout[symbol][col]
}

func (t *Next3Editor) waySetSymbolPayout(symbol int, col int, num int) {
	symbolPayout[symbol][col] = num
}

func (t *Next3Editor) waySetSymbolPayLine(payLinesNum int, col int, num int) {
	PayLines[payLinesNum][col] = num
}

func (t *Next3Editor) wayDebugConfirm() {
	t.DebugCheckBox.SetChecked(false)
	t.DebugSwitch = false
	for j := range t.DebugIndex {
		for i := 0; i < colNum; i++ {
			text := t.DebugSymbolSelect[j][i].Text
			t.DebugIndex[j][i], _ = strconv.Atoi(text)
		}
	}
}

func (t *Next3Editor) waySymbolConfirm() {
	for i := 0; i < symbolNum; i++ {
		// type
		text := t.SymbolTypeSelect[i].Selected
		symbolType, ok := stringToTypeMap[text]
		if !ok {
			return
		}
		symbolTypeMap[i] = symbolType

		// name
		name := t.SymbolNameEntry[i].Text
		lastName := symbolNameList[i]
		delete(stringToSymbol, lastName)
		stringToSymbol[name] = i
		symbolNameList[i] = name
	}
}

func (t *Next3Editor) wayPayLineConfirm() {
	for i := 0; i < payLineNum; i++ {
		for j := 0; j < colNum; j++ {
			PayLines[i][j] = t.PayLineSelect[i][j].SelectedIndex()
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
			if len(round.Result[i].WildPrize) > 0 {
				str += fmt.Sprintf(" wildprize%+v", round.Result[i].WildPrize)
			}
			for _, v := range round.Result[0].WinLine {
				str += fmt.Sprintf("\nSymbol: %+v 得到第 %+v 條線的 %+v 連線 得分：%+v ", v.Symbol, v.WinLine, v.Match, v.Multi)
			}
			str += "\n"

		}
		if i > 0 {
			if round.Result[i].Point > 0 {
				str += "\n中獎內容:"
				if len(round.Result[i].WildPrize) > 0 {
					str += fmt.Sprintf(" wildprize%+v", round.Result[i].WildPrize)
				}
				for _, v := range round.Result[i].WinLine {
					str += fmt.Sprintf("\nSymbol:"+"%+v"+" 得到第 %+v 條線的 %+v 連線 得分：%+v ", v.Symbol, v.WinLine, v.Match, v.Multi)
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
			if len(round.Result[i].WildPrize) > 0 {
				str += fmt.Sprintf(" wildprize%+v", round.Result[i].WildPrize)
			}
			for _, v := range round.Result[0].WinLine {
				str += fmt.Sprintf("\nSymbol:"+"%+v"+" 得到第 %+v 條線的 %+v 連線 得分：%+v ", v.Symbol, v.WinLine, v.Match, v.Multi)
			}
			str += "\n"

		}
		if i > 0 {

			if round.Result[i].Point > 0 {
				str += "\n中獎內容:"
				if len(round.Result[i].WildPrize) > 0 {
					str += fmt.Sprintf(" wildprize%+v", round.Result[i].WildPrize)
				}
				for _, v := range round.Result[i].WinLine {
					str += fmt.Sprintf("\nSymbol:"+"%+v"+" 得到第 %+v 條線的 %+v 連線 得分：%+v ", v.Symbol, v.WinLine, v.Match, v.Multi)
				}
			}
			str += fmt.Sprintf("\n此局得分:%+v\n累積得分:%+v\n", round.Result[i].Point, round.Result[i].TotalPoint)
		}
		str += fmt.Sprintf("此局index:%+v\n", round.Result[i].ReelPosition)
	}
	str += "\n"
	return str
}

func (t *Next3Editor) updateDebugTable(cn int, rn int) {
	t.DebugIndex = make([][]int, 11)
	for i := 0; i < 11; i++ {
		t.DebugIndex[i] = make([]int, cn)
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
