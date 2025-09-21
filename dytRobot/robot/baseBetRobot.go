package robot

import (
	"dytRobot/constant"
	"dytRobot/utils"
)

const (
	RET_ENTER_ROOM = "EnterRoom"
	RET_ENTER_GAME = "EnterGame"
	RET_BET        = "Bet"

	FSM_BET = "Bet"

	BET_WAIT_SECOND = 4
)

type BaseBetRobot struct {
	*BaseRobot

	BetLimit []interface{}
	Chips    []interface{}
	BetArea  []int

	//下注時會依區域倍率大而下注機率小
	AreaWeight  []int
	TotalWeight int

	Bet      func()
	CheckBet func(response *utils.RespBase) bool
	BetOp    func()
}

func NewBetRobot(setting RobotConfig) *BaseBetRobot {
	baseRobot := NewBaseRobot(setting)
	t := &BaseBetRobot{
		BaseRobot: baseRobot,
	}
	t.CheckCommand = t.CheckBetCommand
	t.CheckBet = t.CheckBaseBet
	return t
}

func (t *BaseBetRobot) CheckBetCommand(response *utils.RespBase) int {
	result := t.CheckBaseCommand(response)
	if result != RESPONSE_NO_SUTIALBE {
		return result
	}
	switch response.Ret {
	case RET_LOGIN:
		if data, ok := response.Data.(map[string]interface{}); ok {
			t.UserId = int64(data["UserId"].(float64))
			if response.Code == constant.ERROR_CODE_SUCCESS {
				result = RESPONSE_EXCUTED_SUCCESS
				t.EnterRoom()
			}
		}
	case RET_ENTER_ROOM:
		if response.Code == constant.ERROR_CODE_SUCCESS {
			result = RESPONSE_EXCUTED_SUCCESS
			t.EnterGame()
		}
	case RET_ENTER_GAME:
		switch response.Code {
		case constant.ERROR_CODE_SUCCESS: //成功
			if data, ok := response.Data.(map[string]interface{}); ok {
				t.BetLimit = data["BetLimit"].([]interface{})
				t.Chips = data["Chips"].([]interface{})
				t.BetArea = make([]int, len(t.BetLimit))
				result = RESPONSE_EXCUTED_SUCCESS
			}
		}
	case RET_BET:
		switch response.Code {
		case constant.ERROR_CODE_SUCCESS: //成功
			result = RESPONSE_EXCUTED_SUCCESS
		case constant.ERROR_CODE_ERROR_MONEY_NOT_ENOUGH:
			ok, error := t.Deposit()
			if ok {
				t.LogInfo(utils.LOG_INFO, "上分成功")
				result = RESPONSE_EXCUTED_SUCCESS
			} else {
				t.LogInfo(utils.LOG_ERROR, error.Error())
			}
		case constant.ERROR_CODE_ERROR_BET_LIMIT:
			result = RESPONSE_EXCUTED_SUCCESS
		}
	case ACT_GAME_PERIOD:
		t.CheckGamePeriod(response)
		result = RESPONSE_EXCUTED_SUCCESS
	default:
		result = RESPONSE_NO_SUTIALBE
	}
	return result
}

func (t *BaseBetRobot) EnterRoom() {
	var data struct {
		EnterRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.EnterRoom.GameID = t.GameId

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BaseBetRobot) EnterGame() {
	var data struct {
		EnterGame struct {
			TableId int `json:"tableId"`
		}
	}
	data.EnterGame.TableId = t.TableId

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BaseBetRobot) CheckGamePeriod(response *utils.RespBase) {
	switch t.Fsm {
	case FSM_BET:
		if !t.CheckBet(response) {
			return
		}

		//進到Bet階段才檢查
		if t.PlayLimit == 0 || t.PlayCount < t.PlayLimit {
			go t.Bet()
			t.PlayCount++
		} else {
			t.SetDisConnect()
		}
	}
}

func (t *BaseBetRobot) CheckBaseBet(response *utils.RespBase) bool {
	return true
}
