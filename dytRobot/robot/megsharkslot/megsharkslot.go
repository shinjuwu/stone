package megsharkslot

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
	"math/rand"
	"time"
)

const (
	ACTION_SPIN = 1

	FSM_BET = "Bet"
)

type MegSharkSlotRobot struct {
	*robot.BaseSlotRobot

	BetsLen int

	Bets []float64
}

func NewRobot(setting robot.RobotConfig) *MegSharkSlotRobot {
	slotRobot := robot.NewSlotRobot(setting)
	t := &MegSharkSlotRobot{
		BaseSlotRobot: slotRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	return t
}

func (t *MegSharkSlotRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.CheckSlotCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case "GoInGame":
		t.CheckGoInGame(response)

	case "SlotAction":
		switch response.Code {
		case constant.ERROR_CODE_ERROR_MONEY_NOT_ENOUGH:
			ok, error := t.Deposit()
			if ok {
				t.LogInfo(utils.LOG_INFO, "上分成功")
				t.Spin()
			} else {
				t.LogInfo(utils.LOG_ERROR, error.Error())
			}
		default:
			t.CheckSlotAction(response)
		}
	}

	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *MegSharkSlotRobot) CheckGoInGame(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}

	t.Fsm = data["Fsm"].(string)

	if bets, ok := data["Bets"].([]interface{}); ok {
		t.BetsLen = len(bets)
	}

	switch t.Fsm {
	case FSM_BET:
		t.Spin()
	}
}

func (t *MegSharkSlotRobot) CheckSlotAction(response *utils.RespBase) {
	t.ContinueSpin()
}

func (t *MegSharkSlotRobot) ContinueSpin() {
	if !t.CheckSlotPlayCount() {
		return
	}
	t.Spin()
}

func (t *MegSharkSlotRobot) Spin() {
	betPos := rand.Intn(t.BetsLen)
	var data struct {
		SlotAction struct {
			Action int         `json:"Action"`
			Data   interface{} `json:"Data,omitempty"`
		}
	}
	var info struct {
		BetPos int `json:"BetPos"`
	}
	info.BetPos = betPos

	data.SlotAction.Action = ACTION_SPIN
	data.SlotAction.Data = info

	go t.SendMessage(data)
}

func (t *MegSharkSlotRobot) SendMessage(data interface{}) {
	time.Sleep(1 * time.Second)
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}
