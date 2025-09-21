package fruitslot

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
	"fmt"
	"math/rand"
	"time"
)

const (
	BUTTON_COUNT = 8  // 下注區域數量
	BET_MAX      = 50 // 下注最大量

	ACTION_BET         = 1
	ACTION_GUESS_DOWN  = 4
	ACTION_GUESS_SMALL = 5
	ACTION_GUESS_BIG   = 6

	FSM_BET   = "Bet"
	FSM_GUESS = "Guess"
)

const (
	ROBOT_TYPE_ALL         int = iota // 全押
	ROBOT_TYPE_ALB                    // 只壓蘋果、檸檬、鈴鐺
	ROBOT_TYPE_BM                     // 只押鈴鐺、山竹
	ROBOT_TYPE_DSC                    // 只押鑽石、雙7、皇冠
	ROBOT_TYPE_SMALL_THREE            // 只押小三元(鈴鐺、葡萄、檸檬)
	ROBOT_TYPE_BIG_THREE              // 只押大三元(雙７、鑽石、山竹)
	ROBOT_TYPE_BIG_FOUR               // 只押大四喜(蘋果)

	ROBOT_TYPE_COUNT
)

var RobotType map[int]string = map[int]string{
	ROBOT_TYPE_ALL:         "全押",
	ROBOT_TYPE_ALB:         "只壓蘋果、檸檬、鈴鐺",
	ROBOT_TYPE_BM:          "只押鈴鐺、山竹",
	ROBOT_TYPE_DSC:         "只押鑽石、雙7、皇冠",
	ROBOT_TYPE_SMALL_THREE: "只押小三元(鈴鐺、葡萄、檸檬)",
	ROBOT_TYPE_BIG_THREE:   "只押大三元(雙７、鑽石、山竹)",
	ROBOT_TYPE_BIG_FOUR:    "只押大四喜(蘋果)",
}

// 押注區域列表
const (
	SLOT_ICON_APPLE      int = iota // 蘋果
	SLOT_ICON_LEMON                 // 檸檬
	SLOT_ICON_GRAPE                 // 葡萄
	SLOT_ICON_BELL                  // 鈴鐺
	SLOT_ICON_MANGOSTEEN            // 山竹
	SLOT_ICON_DIAMOND               // 鑽石
	SLOT_ICON_DOUBLE7               // 雙7
	SLOT_ICON_CROWN                 // 皇冠

	SLOT_ICON_NONE //無
)

type FruitSlotRobot struct {
	*robot.BaseElecRobot
	robotType int
	betSum    int
}

func NewRobot(setting robot.RobotConfig) *FruitSlotRobot {
	elecRobot := robot.NewElecRobot(setting)
	t := &FruitSlotRobot{
		BaseElecRobot: elecRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	t.robotType = rand.Intn(ROBOT_TYPE_COUNT)
	t.LogInfo(utils.LOG_INFO, fmt.Sprintf("RobotType=%d, Desc=%s", t.robotType, RobotType[t.robotType]))
	return t
}

func (t *FruitSlotRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.CheckElecCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case "IntoGame":
		t.CheckIntoGame(response)

	case "PlayerAction":
		switch response.Code {
		case constant.ERROR_CODE_ERROR_MONEY_NOT_ENOUGH:
			ok, error := t.Deposit()
			if ok {
				t.LogInfo(utils.LOG_INFO, "上分成功")
				if t.Fsm == FSM_GUESS {
					t.ForceGuessDown()
				} else {
					t.SetBet()
				}
			} else {
				t.LogInfo(utils.LOG_ERROR, error.Error())
			}
		default:
			t.CheckPlayerAction(response)
		}
	}

	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *FruitSlotRobot) CheckIntoGame(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}

	t.Fsm = data["Fsm"].(string)

	switch t.Fsm {
	case FSM_BET:
		t.SetBet()
	case FSM_GUESS:
		win, _ := data["Win"].(float64)
		area, ok := data["GuessArea"].([]interface{})
		if !ok {
			t.ForceGuessDown()
			return
		}
		var guessArea []int
		for _, value := range area {
			bet := int(value.(float64))
			guessArea = append(guessArea, bet)
		}
		t.SetGuess(int(win), guessArea)
	}
}

func (t *FruitSlotRobot) CheckPlayerAction(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}

	t.Fsm = data["Fsm"].(string)
	win := int(data["Win"].(float64))
	if win == 0 { //沒中獎、沒猜中、已下分
		if !t.CheckElecPlayCount() {
			return
		}
		t.SetBet()
	} else {
		area, ok := data["GuessArea"].([]interface{})
		if !ok {
			t.ForceGuessDown()
			return
		}
		var guessArea []int
		for _, value := range area {
			bet := int(value.(float64))
			guessArea = append(guessArea, bet)
		}
		t.SetGuess(win, guessArea)
	}
}

func (t *FruitSlotRobot) SetBet() {
	var data struct {
		PlayerAction struct {
			Action int `json:"Action"`
			Data   struct {
				BetInfo []int `json:"BetInfo"`
			}
		}
	}
	data.PlayerAction.Action = ACTION_BET

	//全圖標隨機
	// for {
	// 	var betInfo []int
	// 	for i := 0; i < BUTTON_COUNT; i++ {
	// 		var bet int
	// 		if rand.Intn(2) == 1 { //代表該圖標要下注
	// 			bet = rand.Intn(BET_MAX) + 1
	// 		}
	// 		betInfo = append(betInfo, bet)
	// 	}

	// 	var isBet bool
	// 	for _, value := range betInfo {
	// 		if value > 0 {
	// 			isBet = true
	// 			break
	// 		}
	// 	}

	// 	if isBet {
	// 		data.PlayerAction.Data.BetInfo = betInfo
	// 		break
	// 	}
	// }

	betInfo := t.fillBet()
	data.PlayerAction.Data.BetInfo = betInfo

	go t.sendMessage(data)
}

func (t *FruitSlotRobot) SetGuess(win int, guessArea []int) {
	var data struct {
		PlayerAction struct {
			Action int `json:"Action"`
			Data   struct {
				BetInfo int `json:"BetInfo"`
			}
		}
	}
	data.PlayerAction.Action = t.randGuess(win)
	if data.PlayerAction.Action != ACTION_GUESS_DOWN {
		data.PlayerAction.Data.BetInfo = t.randGuessBet(guessArea)
	}

	go t.sendMessage(data)
}

func (t *FruitSlotRobot) ForceGuessDown() {
	var data struct {
		PlayerAction struct {
			Action int `json:"Action"`
		}
	}
	data.PlayerAction.Action = ACTION_GUESS_DOWN
	go t.sendMessage(data)
}

func (t *FruitSlotRobot) randGuess(win int) int {
	winTimes := 5
	if rand.Intn(2) == 0 {
		winTimes = 10
	}

	//如果超過倍數就一定下分
	if win >= t.betSum*winTimes {
		t.LogInfo(utils.LOG_DEBUG, fmt.Sprintf("強迫下分win=%d,betSum=%d,winTimes=%d", win, t.betSum, winTimes))
		return ACTION_GUESS_DOWN
	}

	if rand.Intn(5) == 0 {
		t.LogInfo(utils.LOG_DEBUG, "機率下分")
		return ACTION_GUESS_DOWN
	}

	if rand.Intn(2) == 0 {
		t.LogInfo(utils.LOG_DEBUG, "機率猜小")
		return ACTION_GUESS_SMALL
	}
	t.LogInfo(utils.LOG_DEBUG, "機率猜大")
	return ACTION_GUESS_BIG
}

func (t *FruitSlotRobot) randGuessBet(guessArea []int) int {
	if rand.Intn(10) == 0 {
		t.LogInfo(utils.LOG_DEBUG, "機率Half")
		count := len(guessArea) - 2
		if count > 3 {
			count = 3
		}
		p := rand.Intn(count) + 2
		return guessArea[p]
	}

	if rand.Intn(100) == 0 {
		t.LogInfo(utils.LOG_DEBUG, "機率Double")
		return guessArea[0]
	}

	t.LogInfo(utils.LOG_DEBUG, "機率不變")
	return guessArea[1]
}

func (t *FruitSlotRobot) sendMessage(data interface{}) {
	time.Sleep(7 * time.Second)
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *FruitSlotRobot) fillBet() []int {
	var betInfo []int
	t.betSum = 0
	firstBet := rand.Intn(BET_MAX) + 1
	for i := 0; i < BUTTON_COUNT; i++ {
		bet := rand.Intn(11) + firstBet - 5
		if bet < 1 {
			bet = 1
		} else if bet > 50 {
			bet = 50
		}
		switch t.robotType {
		// case ROBOT_TYPE_ALL:
		// 	betInfo = append(betInfo, bet)
		// 	t.betSum += bet
		case ROBOT_TYPE_ALB:
			if i != SLOT_ICON_APPLE && i != SLOT_ICON_LEMON && i != SLOT_ICON_BELL {
				bet = 0
			}
		case ROBOT_TYPE_BM:
			if i != SLOT_ICON_BELL && i != SLOT_ICON_MANGOSTEEN {
				bet = 0
			}
		case ROBOT_TYPE_DSC:
			if i != SLOT_ICON_DIAMOND && i != SLOT_ICON_DOUBLE7 && i != SLOT_ICON_CROWN {
				bet = 0
			}
		case ROBOT_TYPE_SMALL_THREE:
			if i != SLOT_ICON_LEMON && i != SLOT_ICON_GRAPE && i != SLOT_ICON_BELL {
				bet = 0
			}
		case ROBOT_TYPE_BIG_THREE:
			if i != SLOT_ICON_MANGOSTEEN && i != SLOT_ICON_DIAMOND && i != SLOT_ICON_DOUBLE7 {
				bet = 0
			}
		case ROBOT_TYPE_BIG_FOUR:
			if i != SLOT_ICON_APPLE {
				bet = 0
			}
		}
		betInfo = append(betInfo, bet)
		t.betSum += bet
	}
	return betInfo
}
