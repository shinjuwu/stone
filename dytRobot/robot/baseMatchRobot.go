package robot

import (
	"dytRobot/constant"
	"dytRobot/utils"
	"math"
)

const (
	RET_JOIN_ROOM = "JoinRoom"
	RET_JOIN_GAME = "JoinGame"
	RET_QUIT_ROOM = "QuitRoom"
	RET_QUIT_GAME = "QuitGame"
)

type BaseMatchRobot struct {
	*BaseRobot
}

func NewMatchRobot(setting RobotConfig) *BaseMatchRobot {
	baseRobot := NewBaseRobot(setting)
	t := &BaseMatchRobot{
		BaseRobot: baseRobot,
	}
	t.CheckCommand = t.CheckMatchCommand
	return t
}

func (t *BaseMatchRobot) CheckMatchCommand(response *utils.RespBase) int {
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
				t.JoinRoom()
			}
		}
	case RET_JOIN_ROOM:
		if response.Code == constant.ERROR_CODE_SUCCESS {
			data, ok := response.Data.([]interface{})
			if !ok {
				return RESPONSE_EXCUTED_FAILED
			}
			t.EnterInfo = make(map[int]float64)
			for _, info := range data {
				detail, ok := info.(map[string]interface{})
				if !ok {
					return RESPONSE_EXCUTED_FAILED
				}
				tableId := int(detail["TableId"].(float64))
				tableEnter := detail["EnterGold"].(float64)
				t.EnterInfo[tableId] = tableEnter
			}
			result = RESPONSE_EXCUTED_SUCCESS
			t.JoinGame()
		}
	case RET_JOIN_GAME:
		switch response.Code {
		case constant.ERROR_CODE_SUCCESS: //成功
			result = RESPONSE_NO_SUTIALBE //讓各遊戲可以處理
		case constant.ERROR_CODE_BRING_MONEY_LOWER_LIMIT: //金額不足
			ok, error := t.Deposit()
			if ok {
				t.LogInfo(utils.LOG_INFO, "上分成功")
				result = RESPONSE_EXCUTED_SUCCESS
				t.JoinGame()
			} else {
				t.LogInfo(utils.LOG_ERROR, error.Error())
			}
		case constant.ERROR_CODE_ERROR_CAME_FINISHED: //遊戲已結束
			result = RESPONSE_EXCUTED_SUCCESS
			t.JoinGame()
		}
	default:
		result = RESPONSE_NO_SUTIALBE
	}
	return result
}

func (t *BaseMatchRobot) JoinRoom() {
	var data struct {
		JoinRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.JoinRoom.GameID = t.GameId

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BaseMatchRobot) JoinGame() {
	var data struct {
		JoinGame struct {
			TableId   int     `json:"tableId"`
			BringGold float64 `json:"bringGold,omitempty"` //Texas & SingleWallet
		}
	}
	data.JoinGame.TableId = t.TableId
	if t.WalletType == constant.SW_TYPE_SINGLE {
		minBringin := t.EnterInfo[t.TableId]
		swGoldInt := int(t.SWGold)
		bringIn := math.Min(minBringin*100, float64(swGoldInt))
		data.JoinGame.BringGold = bringIn
	} else if (t.TableId / 10) == 2004 {
		data.JoinGame.BringGold = constant.TexasBringGold[(t.TableId % 10)]
	}

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BaseMatchRobot) QuitGame() {
	var data struct {
		QuitGame struct {
			TableID int `json:"tableId"`
		}
	}

	data.QuitGame.TableID = t.TableId
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BaseMatchRobot) CheckMatchPlayCount() {
	//玩完一場才檢查
	t.PlayCount++
	if t.PlayLimit == 0 || t.PlayCount < t.PlayLimit {
		t.JoinGame()
	} else {
		t.SetDisConnect()
	}
}
