package plinko

import (
	"dytRobot/client"
	"dytRobot/utils"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// 增強版 Plinko 客戶端，支援新功能
type EnhancedPlinkoClient struct {
	*client.BaseElecClient
	
	// 原有控件
	buttonBet         *widget.Button
	buttonResult      *widget.Button
	buttonDebug       *widget.Button
	BetID             *widget.Label
	chipsSelect       *widget.Select
	debugBallIDSelect *widget.Select
	debugHoleIDSelect *widget.Select
	debugSwitch       *widget.Check

	// 新增：模式切換控件
	buttonSwitchMode   *widget.Button
	payoutModeSelect   *widget.Select
	currentModeLabel   *widget.Label
	
	// 新增：統計顯示（移除 statsDisplay 和 holeCountLabel, payoutRatesLabel）
	betCountLabel      *widget.Label
	winRateLabel       *widget.Label
	
	// 數據追蹤
	tableBalls         map[int]BallInfo
	currentPayoutMode  int
	holeCount         int
	payoutRates       []float64
	
	// 統計數據
	totalBets         int
	totalWins         int
	totalWinAmount    float64
	totalBetAmount    float64
}

const (
	ACTION_SWITCH_MODE_ENHANCED = 3

	PAYOUT_MODE_CLASSIC = 0
	PAYOUT_MODE_HIGH    = 1
	PAYOUT_MODE_EXTREME = 2
)

func NewEnhancedClient(setting client.ClientConfig) *EnhancedPlinkoClient {
	elecClient := client.NewElecClient(setting)
	t := &EnhancedPlinkoClient{
		BaseElecClient: elecClient,
		tableBalls:     make(map[int]BallInfo),
		currentPayoutMode: PAYOUT_MODE_CLASSIC,
		holeCount:      11,
	}
	t.CheckResponse = t.CheckEnhancedPlinkoResponse
	
	// 新增自定義訊息選項
	t.CustomMessage = append(t.CustomMessage, 
		"{\"PlayerAction\":{\"Action\":3,\"PayoutMode\":%d}}", // 模式切換
	)
	t.EntrySendMessage.SetOptions(t.CustomMessage)

	return t
}

func (t *EnhancedPlinkoClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateElecSection(c)
	t.CreateOptimizedPlinkoSection(c)
	t.CreateTableBalls(c)
	t.CreateBottomSection(c)
	return c
}

func (t *EnhancedPlinkoClient) CreateTableBalls(c *fyne.Container) {
	t.tableBalls = make(map[int]BallInfo)
}

// 優化版：整合所有功能到一個緊湊的區域
func (t *EnhancedPlinkoClient) CreateOptimizedPlinkoSection(c *fyne.Container) {
	// === 遊戲控制區 ===
	t.BetID = widget.NewLabel("💰 籌碼選擇:")
	t.chipsSelect = widget.NewSelect([]string{
		"0 - 最小注 (10金)", 
		"1 - 小注 (50金)", 
		"2 - 中注 (100金)", 
		"3 - 大注 (500金)", 
		"4 - 最大注 (1000金)",
	}, nil)
	t.chipsSelect.SetSelected("1 - 小注 (50金)")
	
	t.buttonBet = widget.NewButton("🎯 投注", func() {
		betStr := t.chipsSelect.Selected
		betID, _ := strconv.Atoi(betStr[:1])
		t.SendGameBet(betID)
		t.totalBets++
		t.updateStats()
	})
	t.buttonResult = widget.NewButton("🏆 領獎", func() {
		t.SendGetResult()
	})

	// === 模式控制區 ===
	t.currentModeLabel = widget.NewLabel("📊 當前: 經典模式 (11洞)")
	t.payoutModeSelect = widget.NewSelect([]string{
		"0 - 🟢 經典 (11洞, 穩定倍率)",
		"1 - 🟡 高倍 (17洞, 高風險)",
		"2 - 🔴 極限 (17洞, 超高倍)",
	}, nil)
	t.payoutModeSelect.SetSelected("0 - 🟢 經典 (11洞, 穩定倍率)")
	
	t.buttonSwitchMode = widget.NewButton("🔄 切換", func() {
		selectedStr := t.payoutModeSelect.Selected
		if selectedStr != "" {
			modeID, _ := strconv.Atoi(selectedStr[:1])
			t.SendSwitchMode(modeID)
		}
	})

	// === 統計顯示區 ===
	t.betCountLabel = widget.NewLabel("📈 投注: 0次")
	t.winRateLabel = widget.NewLabel("💹 勝率: 0% | RTP: 0%")
	resetButton := widget.NewButton("🔄", func() {
		t.resetStats()
	})
	resetButton.Resize(fyne.NewSize(40, 30))

	// === Debug 區域 (摺疊) ===
	t.debugSwitch = widget.NewCheck("🔧 Debug模式", nil)
	t.debugBallIDSelect = widget.NewSelect([]string{"0", "1", "2"}, nil)
	t.debugHoleIDSelect = widget.NewSelect([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16"}, nil)
	t.debugBallIDSelect.SetSelected("0")
	t.debugHoleIDSelect.SetSelected("5")
	t.buttonDebug = widget.NewButton("🎯", func() {
		t.SendDebugInfo()
	})
	t.buttonDebug.Resize(fyne.NewSize(40, 30))

	// === 緊湊布局 ===
	// 第一行：遊戲控制
	gameSection := container.NewHBox(
		t.BetID, t.chipsSelect, t.buttonBet, t.buttonResult,
	)
	
	// 第二行：模式控制
	modeSection := container.NewHBox(
		widget.NewLabel("🎮 模式:"), t.payoutModeSelect, t.buttonSwitchMode,
	)
	
	// 第三行：統計資訊
	statsSection := container.NewHBox(
		t.betCountLabel, t.winRateLabel, resetButton,
	)
	
	// 第四行：當前狀態
	statusSection := container.NewHBox(
		t.currentModeLabel,
	)
	
	// 第五行：Debug (可選)
	debugSection := container.NewHBox(
		t.debugSwitch, 
		widget.NewLabel("球:"), t.debugBallIDSelect,
		widget.NewLabel("洞:"), t.debugHoleIDSelect, 
		t.buttonDebug,
	)

	// 添加分隔線和整體布局
	separator := widget.NewSeparator()
	mainSection := container.NewVBox(
		widget.NewCard("🎲 Plinko 測試控制台", "", 
			container.NewVBox(
				gameSection,
				modeSection, 
				statusSection,
				statsSection,
				separator,
				debugSection,
			),
		),
	)

	c.Add(mainSection)
	
	// 初始化統計顯示
	t.tableBalls = make(map[int]BallInfo)
	t.updateStats()
}

// 移除：此功能已整合到 CreateOptimizedPlinkoSection

// 移除：此功能已整合到 CreateOptimizedPlinkoSection

func (t *EnhancedPlinkoClient) CheckEnhancedPlinkoResponse(response *utils.RespBase) bool {
	if t.CheckBaseResponse(response) {
		return true
	}

	switch response.Ret {
	case client.RET_PLAYER_ACTION:
		t.AddTableStatus(t.GetEnhancedActionInfo(response))
		return true
	case "IntoGame":
		t.updateGameState(response)
		return true
	}

	return t.CheckElecResponse(response)
}

func (t *EnhancedPlinkoClient) GetEnhancedActionInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	
	action := int(data["Action"].(float64))
	var message string
	
	switch action {
	case ACTION_BET:
		if b, ok := data["BallInfo"].(map[string]interface{}); ok {
			sId := b["SID"].(float64)
			ballID := b["BallID"].(float64)
			bet := b["BetID"].(float64)
			holeID := b["HoleID"].(float64)

			ballInfo := BallInfo{
				SID:    int(sId),
				BallID: int(ballID),
				HoleID: int(holeID),
				BetID:  int(bet),
			}
			t.tableBalls[ballInfo.SID] = ballInfo
			
			// 計算預期賠率
			expectedPayout := t.getExpectedPayout(int(holeID))
			message = fmt.Sprintf("🎯 球落下 | 模式:%d | SID:%d | 球種:%d | 洞口:%d | 下注等級:%d | 預期倍率:%.2fx", 
				t.currentPayoutMode, ballInfo.SID, ballInfo.BallID, ballInfo.HoleID, ballInfo.BetID, expectedPayout)
		}
		
	case ACTION_GET_RESULT:
		sId := int(data["SID"].(float64))
		if win, ok := data["Win"].(float64); ok {
			t.totalWinAmount += win
			if win > 0 {
				t.totalWins++
			}
			message = fmt.Sprintf("💰 球兌獎 | SID:%d | 獲勝金額:%.2f", sId, win)
		}
		delete(t.tableBalls, sId)
		t.updateStats()
		
	case ACTION_SWITCH_MODE_ENHANCED:
		// 處理模式切換回應
		if payoutMode, ok := data["PayoutMode"].(float64); ok {
			t.currentPayoutMode = int(payoutMode)
		}
		if holeCount, ok := data["HoleCount"].(float64); ok {
			t.holeCount = int(holeCount)
		}
		if rates, ok := data["PayoutRates"].([]interface{}); ok {
			t.payoutRates = make([]float64, len(rates))
			for i, rate := range rates {
				t.payoutRates[i] = rate.(float64)
			}
		}
		
		t.updateModeDisplay()
		message = fmt.Sprintf("🔄 模式切換成功 | 模式:%d | 洞口數:%d", t.currentPayoutMode, t.holeCount)
	}
	
	message += "\n"
	return message
}

// 新增：發送模式切換請求
func (t *EnhancedPlinkoClient) SendSwitchMode(payoutMode int) {
	var data struct {
		PlayerAction struct {
			Action int                 `json:"Action"`
			Data   map[string]interface{} `json:"Data"`
		}
	}
	
	data.PlayerAction.Action = ACTION_SWITCH_MODE_ENHANCED
	data.PlayerAction.Data = map[string]interface{}{
		"PayoutMode": payoutMode,
	}
	
	utils.SendMessage(t.Connect, t.LoginName, data, nil)
}

// 更新遊戲狀態
func (t *EnhancedPlinkoClient) updateGameState(response *utils.RespBase) {
	if data, ok := response.Data.(map[string]interface{}); ok {
		if payoutMode, exists := data["PayoutMode"].(float64); exists {
			t.currentPayoutMode = int(payoutMode)
		}
		if holeCount, exists := data["HoleCount"].(float64); exists {
			t.holeCount = int(holeCount)
		}
		if rates, exists := data["PayoutRates"].([]interface{}); exists {
			t.payoutRates = make([]float64, len(rates))
			for i, rate := range rates {
				t.payoutRates[i] = rate.(float64)
			}
		}
	}
	t.updateModeDisplay()
}

// 更新模式顯示 - 優化版
func (t *EnhancedPlinkoClient) updateModeDisplay() {
	modeInfo := []struct{
		name string
		icon string
		desc string
	}{
		{"經典", "🟢", "穩定倍率"},
		{"高倍", "🟡", "高風險"},
		{"極限", "🔴", "超高倍"},
	}
	
	if t.currentPayoutMode < len(modeInfo) {
		mode := modeInfo[t.currentPayoutMode]
		t.currentModeLabel.SetText(fmt.Sprintf("📊 當前: %s %s模式 (%d洞, %s)", 
			mode.icon, mode.name, t.holeCount, mode.desc))
	} else {
		t.currentModeLabel.SetText("📊 當前: 未知模式")
	}
}

// 更新統計 - 優化版
func (t *EnhancedPlinkoClient) updateStats() {
	t.betCountLabel.SetText(fmt.Sprintf("📈 投注: %d次", t.totalBets))
	
	winRate := 0.0
	rtp := 0.0
	if t.totalBets > 0 {
		winRate = float64(t.totalWins) / float64(t.totalBets) * 100
		if t.totalBetAmount > 0 {
			rtp = t.totalWinAmount / t.totalBetAmount * 100
		}
	}
	
	// 添加趨勢指示器
	rtpIcon := "📊"
	if rtp > 100 {
		rtpIcon = "💹" // 盈利
	} else if rtp < 50 {
		rtpIcon = "📉" // 虧損
	}
	
	t.winRateLabel.SetText(fmt.Sprintf("%s 勝率: %.1f%% | RTP: %.1f%%", rtpIcon, winRate, rtp))
}

// 重置統計 - 優化版
func (t *EnhancedPlinkoClient) resetStats() {
	t.totalBets = 0
	t.totalWins = 0
	t.totalWinAmount = 0.0
	t.totalBetAmount = 0.0
	t.tableBalls = make(map[int]BallInfo)
	t.updateStats()
	
	// 簡化提示，不需要額外的文字區域
	t.AddTableStatus("🔄 統計數據已重置")
}

// 移除重複的 resetStats 函數，保留優化版

// 獲取預期賠率
func (t *EnhancedPlinkoClient) getExpectedPayout(holeID int) float64 {
	if holeID < len(t.payoutRates) {
		return t.payoutRates[holeID]
	}
	return 0.0
}

// 實現原有接口方法
func (t *EnhancedPlinkoClient) SendGameBet(betID int) {
	// 實現發送下注的邏輯
	var data struct {
		PlayerAction struct {
			Action int
			Data   struct {
				BetInfo int `json:"BetInfo"`
			}
		}
	}

	data.PlayerAction.Action = ACTION_BET
	data.PlayerAction.Data.BetInfo = betID

	utils.SendMessage(t.Connect, t.LoginName, data, nil)
}

func (t *EnhancedPlinkoClient) SendGetResult() {
	for sId, _ := range t.tableBalls {
		var data struct {
			PlayerAction struct {
				Action int
				Data   struct {
					SID int `json:"SID"`
				}
			}
		}
		data.PlayerAction.Action = ACTION_GET_RESULT
		data.PlayerAction.Data.SID = sId

		utils.SendMessage(t.Connect, t.LoginName, data, nil)
	}
}

func (t *EnhancedPlinkoClient) SendDebugInfo() {
	// 實現 Debug 功能
	ballID, _ := strconv.Atoi(t.debugBallIDSelect.Selected)
	holeID, _ := strconv.Atoi(t.debugHoleIDSelect.Selected)
	
	debugMsg := fmt.Sprintf("{\"DebugBall\":{\"BallID\":%d,\"HoleID\":%d,\"Enable\":%t}}", 
		ballID, holeID, t.debugSwitch.Checked)
	
	utils.SendMessage(t.Connect, t.LoginName, debugMsg, nil)
}