package robot

import (
	"dytRobot/constant"
	"dytRobot/utils"
	"math"
)

type BaseElecRobot struct {
	*BaseRobot
}

func NewElecRobot(setting RobotConfig) *BaseElecRobot {
	baseRobot := NewBaseRobot(setting)
	t := &BaseElecRobot{
		BaseRobot: baseRobot,
	}
	t.CheckCommand = t.CheckElecCommand
	return t
}

func (t *BaseElecRobot) CheckElecCommand(response *utils.RespBase) int {
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
				t.IntoRoom()
			}
		}
	case "IntoRoom":
		if response.Code == constant.ERROR_CODE_SUCCESS {
			result = RESPONSE_EXCUTED_SUCCESS
			t.IntoGame()
		}
	case "IntoGame":
		switch response.Code {
		case constant.ERROR_CODE_SUCCESS: //成功
			result = RESPONSE_NO_SUTIALBE //讓各遊戲可以處理
		case constant.ERROR_CODE_BRING_MONEY_LOWER_LIMIT: //金額不足
			ok, error := t.Deposit()
			if ok {
				t.LogInfo(utils.LOG_INFO, "上分成功")
				result = RESPONSE_EXCUTED_SUCCESS
				t.IntoGame()
			} else {
				t.LogInfo(utils.LOG_ERROR, error.Error())
			}
		case constant.ERROR_CODE_ERROR_CAME_FINISHED: //遊戲已結束
			result = RESPONSE_EXCUTED_SUCCESS
			t.IntoGame()
		}
	default:
		result = RESPONSE_NO_SUTIALBE
	}
	return result
}

func (t *BaseElecRobot) IntoRoom() {
	var data struct {
		IntoRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.IntoRoom.GameID = t.GameId

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BaseElecRobot) IntoGame() {
	var data struct {
		IntoGame struct {
			TableId   int     `json:"tableId"`
			BringGold float64 `json:"bringGold,omitempty"` //Texas & SingleWallet
		}
	}
	data.IntoGame.TableId = t.TableId
	if t.WalletType == constant.SW_TYPE_SINGLE {
		minBringin := t.EnterInfo[t.TableId]
		swGoldInt := int(t.SWGold)
		bringIn := math.Min(minBringin*100, float64(swGoldInt))
		data.IntoGame.BringGold = bringIn
	}

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BaseElecRobot) CheckElecPlayCount() bool {
	t.PlayCount++
	//動作完檢查
	if t.PlayLimit != 0 && t.PlayCount >= t.PlayLimit {
		t.SetDisConnect()
		return false
	}
	return true
}
