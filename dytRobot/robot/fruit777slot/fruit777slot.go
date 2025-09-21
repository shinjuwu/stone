package fruit777slot

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
	"fmt"
	"math/rand"
	"time"
)

const (
	ACTION_BET  = 1
	ACTION_TAKE = 2

	FSM_BET = "Bet"
)

type Fruit777SlotRobot struct {
	*robot.BaseSlotRobot

	Line int
	Bet  float64

	Lines int
	Bets  []float64
}

func NewRobot(setting robot.RobotConfig) *Fruit777SlotRobot {
	slotRobot := robot.NewSlotRobot(setting)
	t := &Fruit777SlotRobot{
		BaseSlotRobot: slotRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	return t
}

func (t *Fruit777SlotRobot) GoCheckCommand(response *utils.RespBase) int {
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
				t.ContinueSpin()
			} else {
				t.LogInfo(utils.LOG_ERROR, error.Error())
			}
		default:
			t.CheckSlotAction(response)
		}
	}

	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *Fruit777SlotRobot) CheckGoInGame(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}

	t.Fsm = data["Fsm"].(string)
	t.Lines = int(data["Lines"].(float64))

	if bets, ok := data["Bets"].([]interface{}); ok {
		for _, bet := range bets {
			t.Bets = append(t.Bets, bet.(float64))
		}
	}

	switch t.Fsm {
	case FSM_BET:
		// t.Line = rand.Intn(t.Lines) + 1
		t.Line = 8
		t.Bet = t.Bets[rand.Intn(len(t.Bets))]
		t.LogInfo(utils.LOG_INFO, fmt.Sprintf("Line: %d, Bet: %.1f", t.Line, t.Bet))
		t.Spin()
	}
}

func (t *Fruit777SlotRobot) CheckSlotAction(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}

	action := int(data["Action"].(float64))
	switch action {
	case ACTION_BET:
		win := data["Win"].(float64)
		if win > 0 {
			t.Take()
		} else {
			t.ContinueSpin()
		}
	case ACTION_TAKE:
		t.ContinueSpin()
	}
}

func (t *Fruit777SlotRobot) Spin() {
	var data struct {
		SlotAction struct {
			Action int `json:"Action"`
			Data   struct {
				Line int     `json:"Line"`
				Bet  float64 `json:"Bet"`
			}
		}
	}

	data.SlotAction.Action = ACTION_BET
	data.SlotAction.Data.Line = t.Line
	data.SlotAction.Data.Bet = t.Bet

	go t.SendMessage(data)
}

func (t *Fruit777SlotRobot) ContinueSpin() {
	if !t.CheckSlotPlayCount() {
		return
	}
	t.Spin()
}
func (t *Fruit777SlotRobot) Take() {
	var data struct {
		SlotAction struct {
			Action int `json:"Action"`
		}
	}

	data.SlotAction.Action = ACTION_TAKE

	go t.SendMessage(data)
}

func (t *Fruit777SlotRobot) SendMessage(data interface{}) {
	time.Sleep(1 * time.Second)
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}
