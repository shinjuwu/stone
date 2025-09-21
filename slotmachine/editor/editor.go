package editor

import (
	"fmt"
	"io"
	"log"
	"slotEditor/constant"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
)

type BaseConfig struct {
	GameId int
	Window *fyne.Window
	Rtp    string
}

type BaseInfo struct {
	ColNum    int
	RowNum    int
	SymbolNum int
}

type BaseEditor struct {
	GameId     int
	Name       string
	Info       BaseInfo
	MainWindow *fyne.Window

	Gold         float64
	Round        int
	FreeRound    int //FG總次數
	RtpCurr      float64
	RtpMillion   float64
	SlotType     int
	Rtp          string
	FreeTimes    int //FGspin計數
	Stage        int //FG正數
	Unitbet      string
	MaxScore     float64
	NgGetWin     int
	FgGetWin     int
	TotalFgRound int

	CurWin        float64
	StatBet       float64
	StatWin       float64
	StatNormalWin float64
	StatBonusWin  float64
	CurNG         float64
	CurBG         float64
	CurPos        []int
	CurBGdetail   [][]float64
	Coin          float64
	// widget
	SymbolPayoutEntry [][]*widget.Entry  // symbol 賠倍設置
	SymbolWeightEntry [][]*widget.Entry  // symbol 權重設置
	SymbolPercentage  [][]*widget.Label  // symbol 出現機率設置F
	SymbolNameEntry   []*widget.Entry    // symbol 名稱設置
	SymbolTypeSelect  []*widget.Select   // symbol 型態設置
	BtnInit           *widget.Button     // 重置按鈕
	LabelGold         *widget.Label      // 當前金標示
	LabelCurRound     *widget.Label      // 當前輪次標示
	LabelFreeRound    *widget.Label      // 當前免費倫次標示
	EntryDisplay      *widget.Entry      // 輸出畫面
	PayLineSelect     [][]*widget.Select // 中獎線選項

	LabelStatBet          *widget.Label // 當前下注標示
	LabelStatWin          *widget.Label // 當前贏分標示
	LabelRTP              *widget.Label // 當前RTP
	LabelStatNormalWin    *widget.Label // 當前普通贏分
	LabelStatBonusWin     *widget.Label // 當前特色贏分
	LabelNormalRTP        *widget.Label // 當前普通RTP
	LabelBonusRTP         *widget.Label // 當前特色贏分
	LabelBonusTriggerRate *widget.Label // 當前特色觸發率
	LabelTestMax          *widget.Label //測試最大倍率
	LabelTestNgWin        *widget.Label //測試Ng連線機率
	LabelTestFgWin        *widget.Label //測試Fg連線機率

	DebugSwitch       bool              // Debug切換
	DebugSymbolSelect [][]*widget.Entry // Debug symbol 選項
	DebugCheckBox     *widget.Check     // Debug check box

	// method
	DebugConfirm          func()                   // 確認Debug內容
	InitState             func()                   // 狀態初始化
	RunOnce               func()                   // 拉霸一次
	Display               func()                   // 顯示畫面
	ServerDisplay         func()                   // server 模式顯示畫面
	LoadFile              func()                   // 讀取檔案
	SaveFile              func()                   // 儲存檔案
	UpdateInfo            func([]byte)             // 更新設置到遊戲內部
	CalcSymbolTotalWeight func()                   // 計算加總權重
	RecreateWindow        func() fyne.CanvasObject // 重新刷新window
	SetBasicInfo          func(int, int, int, int) // 設定基礎內容 col row symbol payline
	ServerRunOnce         func()                   // server 拉霸一次

	// get set
	GetColNum             func() int
	GetRowNum             func() int
	GetPayLineNum         func() int
	GetSymbolNum          func() int
	GetSymbolName         func(int) string
	SetSymbolName         func(int, string)
	GetSymbolNameList     func() []string
	SetSymbolNameList     func([]string)
	GetSymbolPayout       func(int, int) int
	SetSymbolPayout       func(int, int, int)
	SetSymbolPayLine      func(int, int, int)
	GetSymbolTypeNameList func() []string
	GetSymbolType         func(int) int
}

func NewBaseEditor(config BaseConfig) *BaseEditor {
	t := &BaseEditor{
		GameId:     config.GameId,
		MainWindow: config.Window,
		Rtp:        config.Rtp,
	}
	return t
}

func (t *BaseEditor) Refresh() {
	curWindow := *t.MainWindow
	curWindow.SetContent(t.RecreateWindow())
}

func (t *BaseEditor) CreateOverallSection() (c *fyne.Container) {
	c = container.NewVBox(widget.NewLabel("總覽內容"))
	labelGoldName := widget.NewLabel("當前金")
	labelGold := widget.NewLabel("0")
	labelCurRoundName := widget.NewLabel("當前輪次")
	labelCurRound := widget.NewLabel("0")
	labelFreeRoundName := widget.NewLabel("當前免費輪次")
	labelFreeRound := widget.NewLabel("0")
	initBnt := widget.NewButton("重置", func() {
		t.InitState()
	})
	debugCheckBox := widget.NewCheck("Debug", func(b bool) {
		t.DebugSwitch = b
	})

	t.DebugCheckBox = debugCheckBox
	t.LabelGold = labelGold
	t.LabelCurRound = labelCurRound
	t.LabelFreeRound = labelFreeRound
	t.BtnInit = initBnt
	tc := container.NewHBox(labelGoldName, labelGold, labelCurRoundName, labelCurRound, labelFreeRoundName, labelFreeRound, initBnt, debugCheckBox)
	c.Add(tc)

	entryDisplay := widget.NewMultiLineEntry()
	entryDisplay.SetMinRowsVisible(10) //DisplayWindowHeight
	t.EntryDisplay = entryDisplay
	c.Add(entryDisplay)
	btnRunOnce := widget.NewButton("單次", func() {
		t.RunOnce()
		t.Display()
	})
	c.Add(btnRunOnce)
	btnRunMulti := widget.NewButton("百萬測試", func() {
		t.RunMulti()
	})
	c.Add(btnRunMulti)
	btnServerRunOnce := widget.NewButton("當前Server單次", func() {
		t.ServerRunOnce()
		t.ServerDisplay()
	})
	c.Add(btnServerRunOnce)
	btnServerRunMulti := widget.NewButton("當前Server百萬測試", func() {
		t.ServerRunMulti()
	})
	c.Add(btnServerRunMulti)
	labelStatWinName := widget.NewLabel("當前玩家贏分")
	labelStatWin := widget.NewLabel("")
	labelStatBetName := widget.NewLabel("當前玩家下注")
	labelStatBet := widget.NewLabel("")
	labelRTPName := widget.NewLabel("當前RTP")
	labelRTP := widget.NewLabel("0")
	t.LabelStatWin = labelStatWin
	t.LabelStatBet = labelStatBet
	t.LabelRTP = labelRTP
	c.Add(container.NewHBox(labelStatWinName, labelStatWin, labelStatBetName, labelStatBet, labelRTPName, labelRTP))

	labelStatNormalWinName := widget.NewLabel("當前NG玩家贏分")
	labelStatNormalWin := widget.NewLabel("")
	labelNormalRTPName := widget.NewLabel("當前NG RTP")
	labelNormalRTP := widget.NewLabel("0")
	labelStatBonusWinName := widget.NewLabel("當前FG玩家贏分")
	labelStatBonusWin := widget.NewLabel("")
	labelBonusRTPName := widget.NewLabel("當前FG RTP")
	labelBonusRTP := widget.NewLabel("0")
	labelBonusTriggerName := widget.NewLabel("當前FG觸發率")
	labelBonusTriggerRate := widget.NewLabel("0")
	t.LabelStatNormalWin = labelStatNormalWin
	t.LabelStatBonusWin = labelStatBonusWin
	t.LabelNormalRTP = labelNormalRTP
	t.LabelBonusRTP = labelBonusRTP
	t.LabelBonusTriggerRate = labelBonusTriggerRate
	c.Add(container.NewHBox(labelStatNormalWinName, labelStatNormalWin, labelNormalRTPName, labelNormalRTP, labelStatBonusWinName, labelStatBonusWin, labelBonusRTPName, labelBonusRTP, labelBonusTriggerName, labelBonusTriggerRate))

	labelTestNgWinName := widget.NewLabel("NG連線機率")
	labelTestNgWin := widget.NewLabel("0")
	labelTestFgWinName := widget.NewLabel("FG連線機率")
	labelTestFgWin := widget.NewLabel("0")
	labelTestMaxName := widget.NewLabel("測試最大倍率")
	labelTestMax := widget.NewLabel("0")
	t.LabelTestNgWin = labelTestNgWin
	t.LabelTestFgWin = labelTestFgWin
	t.LabelTestMax = labelTestMax
	c.Add(container.NewHBox(labelTestNgWinName, labelTestNgWin, labelTestFgWinName, labelTestFgWin, labelTestMaxName, labelTestMax))

	LoadBtn := widget.NewButton("讀取", func() {
		t.LoadFile()
	})
	SaveBtn := widget.NewButton("存檔", func() {
		t.SaveFile()
	})
	ExcelBtn := widget.NewButton(strconv.Itoa(constant.EXCEL_TIMES)+"次Excel輸出", func() {
		t.ExcelExport()
	})
	c.Add(container.NewHBox(LoadBtn, SaveBtn, ExcelBtn))

	t.InitState()
	t.DebugCheckBox.SetChecked(false)
	return
}

func (t *BaseEditor) UpdateOverallInfo(str string) {
	t.LabelGold.SetText(fmt.Sprintf("%.2f", t.Gold))
	t.LabelCurRound.SetText(fmt.Sprintf("%d", t.Round))
	t.LabelFreeRound.SetText(fmt.Sprintf("%d", t.FreeRound))
	t.EntryDisplay.SetText(str)

	rtp := float64(t.StatWin) / float64(t.StatBet)
	nRtp := float64(t.StatNormalWin) / float64(t.StatBet)
	bRtp := float64(t.StatBonusWin) / float64(t.StatBet) //float64(t.FreeTimes) // float64(t.Extra.FGtotalRound) //
	bHitRate := float64(t.FreeTimes) / float64(t.Round)
	NgWinRate := float64(t.NgGetWin) / float64(t.Round)
	FgWinRate := float64(t.FgGetWin) / float64(t.TotalFgRound)
	t.LabelStatBet.SetText(fmt.Sprintf("%.4f", t.StatBet))
	t.LabelStatWin.SetText(fmt.Sprintf("%.4f", t.StatWin))
	t.LabelRTP.SetText(fmt.Sprintf("%.4f", rtp))
	t.LabelStatNormalWin.SetText(fmt.Sprintf("%.4f", t.StatNormalWin))
	t.LabelNormalRTP.SetText(fmt.Sprintf("%.6f", nRtp))
	t.LabelStatBonusWin.SetText(fmt.Sprintf("%.4f", t.StatBonusWin))
	t.LabelBonusRTP.SetText(fmt.Sprintf("%.6f", bRtp))
	t.LabelBonusTriggerRate.SetText(fmt.Sprintf("%.6f", bHitRate))
	t.LabelTestMax.SetText(fmt.Sprintf("%.4f", t.MaxScore))
	t.LabelTestNgWin.SetText(fmt.Sprintf("%.4f", NgWinRate))
	t.LabelTestFgWin.SetText(fmt.Sprintf("%.4f", FgWinRate))
}

func (t *BaseEditor) RunMulti() {
	count := 0
	curTime := time.Now()
	for t.Round < constant.MULTI_TIMES {
		count++
		t.RunOnce()
		if count%1000 == 0 {
			t.Display()
		}
	}
	t.Display()
	usedTime := time.Since(curTime)
	fmt.Printf("usedTime:%s\n", usedTime)

}

func (t *BaseEditor) ServerRunMulti() {
	count := 0
	curTime := time.Now()
	for t.Round < constant.MULTI_TIMES {
		count++
		t.ServerRunOnce()
		if count%1000 == 0 {
			t.ServerDisplay()
		}
	}
	t.ServerDisplay()
	usedTime := time.Since(curTime)
	fmt.Printf("usedTime:%s\n", usedTime)
}

func (t *BaseEditor) ExcelExport() {
	t.InitState()
	count := 0
	curTime := time.Now()
	data := [][]interface{}{
		{"ID", "下注", "獲利", "NG win", "BG win", "累積獲利", "rtp(%)", "期望", "測試時間", "餘額"},
	}
	//preStatWin := float64(0)
	curWin := float64(0)
	for t.Round < constant.EXCEL_TIMES {
		count++
		t.RunOnce()
		//t.ServerRunOnce()
		//curWin = t.StatWin - preStatWin
		t.Coin = float64(constant.EXCEL_TIMES-count) + t.StatWin
		curWin = t.CurWin
		//t.CurNG = curWin - t.CurBG
		data = append(data, []interface{}{count, t.StatBet, curWin, t.CurNG, t.CurBG, t.StatWin, t.StatWin / t.StatBet * 100, t.StatBet * 0.98, curTime, t.Coin})
		//preStatWin = t.StatWin
	}
	f := excelize.NewFile()
	for idx, row := range data {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			fmt.Println(err)
			return
		}
		f.SetSheetRow("Sheet1", cell, &row)
	}
	// Save spreadsheet by the given path.
	if err := f.SaveAs(t.Name + ".xlsx"); err != nil {
		fmt.Println(err)
	}
	usedTime := time.Since(curTime)
	fmt.Printf("usedTime:%s\n", usedTime)
	t.InitState()
}

func (t *BaseEditor) UpdatePayoutSection() {
	for i := 0; i < len(t.SymbolPayoutEntry); i++ {
		for j := 0; j < len(t.SymbolPayoutEntry[i]); j++ {
			t.SymbolPayoutEntry[i][j].SetText(fmt.Sprintf("%d", t.GetSymbolPayout(i, j)))
		}
	}
}

func (t *BaseEditor) CreateDebugSection() *fyne.Container {
	c := container.NewVBox()
	t.DebugSymbolSelect = make([][]*widget.Entry, 11)

	colNum := t.GetColNum()
	for j := 0; j < 11; j++ {
		debugContainer := container.NewHBox()
		for i := 0; i < colNum+1; i++ { //多一個位置放WWSYM
			index := widget.NewEntry()
			index.SetText("-1")
			t.DebugSymbolSelect[j] = append(t.DebugSymbolSelect[j], index)
			debugContainer.Add(index)
		}
		c.Add(debugContainer)
	}
	confirmBtn := widget.NewButton("確認", func() {
		t.DebugConfirm()
	})

	return container.NewBorder(container.NewVBox(widget.NewLabel("Debug:\n1.前5位為要測試的ReelIndex,輸入範圍0-80\n	範圍之外預設自動抽\n2.最後一位為FG次數	可輸入9-27/12-36/15-45\n	輸入範圍之外的數字會歸為自動")), confirmBtn, nil, nil, c)
}

func (t *BaseEditor) EditorSave(m []byte, win fyne.Window) {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if writer == nil {
			log.Println("Cancelled")
			return
		}

		fileSaved(writer, win, m)
	}, win)
}

func (t *BaseEditor) EditorLoad(win fyne.Window) {
	m := []byte{}
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, win)
			return
		}
		if reader == nil {
			log.Println("Cancelled")
			return
		}

		m = fileLoaded(reader, win)
		t.UpdateInfo(m)
	}, win)
	fd.SetFilter(storage.NewExtensionFileFilter([]string{".json"}))
	fd.Show()
}

func fileSaved(f fyne.URIWriteCloser, w fyne.Window, m []byte) {
	defer f.Close()
	_, err := f.Write(m)
	if err != nil {
		dialog.ShowError(err, w)
	}
	err = f.Close()
	if err != nil {
		dialog.ShowError(err, w)
	}
}

func fileLoaded(f fyne.URIReadCloser, w fyne.Window) []byte {
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		dialog.ShowError(err, w)
	}
	err = f.Close()
	if err != nil {
		dialog.ShowError(err, w)
	}
	return data
}
