package plinko

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
	"fmt"
	"math/rand"
	"time"
)

// 增強版 Plinko 機器人，支援新功能
type EnhancedPlinkoRobot struct {
	*robot.BaseElecRobot
	tableBalls map[int]BallInfo
	
	// 新增：倍率模式相關
	currentPayoutMode int
	holeCount        int
	payoutRates      []float64
	
	// 新增：測試配置
	testConfig TestConfig
}

// 測試配置
type TestConfig struct {
	EnableModeSwitch   bool    // 是否啟用模式切換測試
	ModeSwitchInterval int     // 模式切換間隔（秒）
	TestAllModes      bool    // 是否測試所有模式
	BetVariation      bool    // 是否使用隨機下注金額
	MaxBetsPerMode    int     // 每個模式最大下注次數
}

// 新增常量 (使用不同名稱避免重複定義)
const (
	ACTION_SWITCH_MODE = 3

	PAYOUT_MODE_CLASSIC = 0  // 經典模式 (11個洞)
	PAYOUT_MODE_HIGH    = 1  // 高倍模式 (17個洞) 
	PAYOUT_MODE_EXTREME = 2  // 極限模式 (17個洞)
	PAYOUT_MODE_MAX     = 3
)

// 倍率模式切換請求
type SwitchModeRequest struct {
	Action     int `json:"Action"`
	PayoutMode int `json:"PayoutMode"`
}

// 倍率模式回應
type SwitchModeResponse struct {
	Action      int       `json:"Action"`
	PayoutMode  int       `json:"PayoutMode"`
	HoleCount   int       `json:"HoleCount"`
	PayoutRates []float64 `json:"PayoutRates"`
}

func NewEnhancedRobot(setting robot.RobotConfig) *EnhancedPlinkoRobot {
	elecRobot := robot.NewElecRobot(setting)
	t := &EnhancedPlinkoRobot{
		BaseElecRobot: elecRobot,
		tableBalls:    make(map[int]BallInfo),
		
		// 預設配置
		currentPayoutMode: PAYOUT_MODE_CLASSIC,
		holeCount:        11,
		testConfig: TestConfig{
			EnableModeSwitch:   true,
			ModeSwitchInterval: 30, // 30秒切換一次
			TestAllModes:      true,
			BetVariation:      true,
			MaxBetsPerMode:    10,
		},
	}
	t.CheckCommand = t.GoCheckCommand

	// 啟動模式切換定時器
	if t.testConfig.EnableModeSwitch {
		go t.modeSwitchTimer()
	}

	return t
}

// StartRobot 啟動機器人 (兼容原始接口)
func (t *EnhancedPlinkoRobot) StartRobot() {
	t.BaseElecRobot.StartRobot()
}

func (t *EnhancedPlinkoRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.BaseElecRobot.CheckElecCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}
	
	switch response.Ret {
	case "IntoGame":
		t.CheckIntoGame(response)
	case "PlayerAction":
		t.CheckPlayerAction(response)
	}
	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *EnhancedPlinkoRobot) CheckPlayerAction(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	
	action := int(data["Action"].(float64))
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
			
			// 記錄測試數據
			fmt.Printf("球落下 - 模式:%d, SID:%d, 球種:%d, 洞口:%d, 下注:%d\n", 
				t.currentPayoutMode, ballInfo.SID, ballInfo.BallID, ballInfo.HoleID, ballInfo.BetID)
		}
		t.goGetResult()
		
	case ACTION_GET_RESULT:
		sId := int(data["SID"].(float64))
		if win, ok := data["Win"].(float64); ok {
			fmt.Printf("球兌獎 - SID:%d, 獲勝金額:%.2f\n", sId, win)
		}
		delete(t.tableBalls, sId)
		t.goBet()
		
	case ACTION_SWITCH_MODE:
		// 新增：處理模式切換回應
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
		fmt.Printf("模式切換成功 - 模式:%d, 洞口數:%d, 倍率:%v\n", 
			t.currentPayoutMode, t.holeCount, t.payoutRates)
	}
}

func (t *EnhancedPlinkoRobot) CheckIntoGame(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	// 解析遊戲狀態
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

	fmt.Printf("進入遊戲 - 當前模式:%d, 洞口數:%d\n", t.currentPayoutMode, t.holeCount)
	t.goBet()
}

func (t *EnhancedPlinkoRobot) goGetResult() {
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

		utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
	}
}

func (t *EnhancedPlinkoRobot) goBet() {
	var data struct {
		PlayerAction struct {
			Action int
			Data   struct {
				BetInfo int `json:"BetInfo"`
			}
		}
	}

	data.PlayerAction.Action = ACTION_BET
	
	// 新增：隨機下注變化
	if t.testConfig.BetVariation {
		data.PlayerAction.Data.BetInfo = rand.Intn(5) // 0-4 隨機下注等級
	} else {
		data.PlayerAction.Data.BetInfo = 0
	}

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

// 新增：倍率模式切換功能
func (t *EnhancedPlinkoRobot) goSwitchMode(targetMode int) {
	if targetMode < 0 || targetMode >= PAYOUT_MODE_MAX {
		return
	}

	var data struct {
		PlayerAction SwitchModeRequest
	}
	
	data.PlayerAction.Action = ACTION_SWITCH_MODE
	data.PlayerAction.PayoutMode = targetMode

	fmt.Printf("請求切換模式到: %d\n", targetMode)
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

// 新增：自動模式切換定時器
func (t *EnhancedPlinkoRobot) modeSwitchTimer() {
	ticker := time.NewTicker(time.Duration(t.testConfig.ModeSwitchInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if t.testConfig.TestAllModes {
				// 循環測試所有模式
				nextMode := (t.currentPayoutMode + 1) % PAYOUT_MODE_MAX
				t.goSwitchMode(nextMode)
			}
		}
	}
}

// 新增：測試統計功能
func (t *EnhancedPlinkoRobot) GetTestStats() map[string]interface{} {
	return map[string]interface{}{
		"current_mode":   t.currentPayoutMode,
		"hole_count":     t.holeCount,
		"payout_rates":   t.payoutRates,
		"active_balls":   len(t.tableBalls),
		"test_config":    t.testConfig,
	}
}