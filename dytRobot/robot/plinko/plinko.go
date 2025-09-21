package plinko

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
)

type BallInfo struct {
	BallID int `json:"BallID"`
	HoleID int `json:"HoleID"`
	BetID  int `json:"BetID"`
	SID    int `json:"SID"`
}
type PlinkoRobot struct {
	*robot.BaseElecRobot
	tableBalls map[int]BallInfo
}

const (
	ACTION_BET        = 1
	ACTION_GET_RESULT = 2

	FSM_BET = "Bet"
)

func NewRobot(setting robot.RobotConfig) *PlinkoRobot {
	elecRobot := robot.NewElecRobot(setting)
	t := &PlinkoRobot{
		BaseElecRobot: elecRobot,
		tableBalls:    make(map[int]BallInfo),
	}
	t.CheckCommand = t.GoCheckCommand

	return t
}

func (t *PlinkoRobot) GoCheckCommand(response *utils.RespBase) int {
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

func (t *PlinkoRobot) CheckPlayerAction(response *utils.RespBase) {
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
		}
		t.goGetResult()
	case ACTION_GET_RESULT:
		sId := int(data["SID"].(float64))
		delete(t.tableBalls, sId)
		t.goBet()
	}
}

func (t *PlinkoRobot) CheckIntoGame(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	// data, ok := response.Data.(map[string]interface{})
	// if !ok {
	// 	return
	// }

	// t.Fsm = data["Fsm"].(string)

	t.goBet()
	// switch t.Fsm {
	// case FSM_BET:
	// 	t.goBet()
	// }
}

func (t *PlinkoRobot) goGetResult() {
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

func (t *PlinkoRobot) goBet() {
	var data struct {
		PlayerAction struct {
			Action int
			Data   struct {
				BetInfo int `json:"BetInfo"`
			}
		}
	}

	data.PlayerAction.Action = ACTION_BET
	data.PlayerAction.Data.BetInfo = 0

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}
