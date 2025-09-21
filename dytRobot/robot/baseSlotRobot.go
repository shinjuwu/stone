package robot

import (
	"dytRobot/constant"
	"dytRobot/utils"
)

type BaseSlotRobot struct {
	*BaseRobot
}

func NewSlotRobot(setting RobotConfig) *BaseSlotRobot {
	baseRobot := NewBaseRobot(setting)
	t := &BaseSlotRobot{
		BaseRobot: baseRobot,
	}
	t.CheckCommand = t.CheckSlotCommand
	return t
}

func (t *BaseSlotRobot) CheckSlotCommand(response *utils.RespBase) int {
	result := t.CheckBaseCommand(response)
	if result != RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case "Login":
		if data, ok := response.Data.(map[string]interface{}); ok {
			t.UserId = int64(data["UserId"].(float64))
			if response.Code == constant.ERROR_CODE_SUCCESS {
				result = RESPONSE_EXCUTED_SUCCESS
				t.GoInRoom()
			}
		}
	case "GoInRoom":
		if response.Code == constant.ERROR_CODE_SUCCESS {
			result = RESPONSE_EXCUTED_SUCCESS
			t.GoInGame()
		}
	case "GoInGame":
		switch response.Code {
		case constant.ERROR_CODE_SUCCESS: //成功
			result = RESPONSE_NO_SUTIALBE //讓各遊戲可以處理
		case constant.ERROR_CODE_BRING_MONEY_LOWER_LIMIT: //金額不足
			ok, error := t.Deposit()
			if ok {
				t.LogInfo(utils.LOG_INFO, "上分成功")
				result = RESPONSE_EXCUTED_SUCCESS
				t.GoInGame()
			} else {
				t.LogInfo(utils.LOG_ERROR, error.Error())
			}
		case constant.ERROR_CODE_ERROR_CAME_FINISHED: //遊戲已結束
			result = RESPONSE_EXCUTED_SUCCESS
			t.GoInGame()
		}
	default:
		result = RESPONSE_NO_SUTIALBE
	}
	return result
}

func (t *BaseSlotRobot) GoInRoom() {
	var data struct {
		GoInRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.GoInRoom.GameID = t.GameId

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BaseSlotRobot) GoInGame() {
	var data struct {
		GoInGame struct {
			TableId int `json:"tableId"`
		}
	}
	data.GoInGame.TableId = t.TableId

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BaseSlotRobot) CheckSlotPlayCount() bool {
	t.PlayCount++
	//動作完檢查
	if t.PlayLimit != 0 && t.PlayCount >= t.PlayLimit {
		t.SetDisConnect()
		return false
	}
	return true
}
