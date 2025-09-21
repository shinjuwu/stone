package teenpatti

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
	"math/rand"
)

type PlayerAction struct {
	PlayOperate struct {
		Instruction int         `json:"instruction"`
		Data        interface{} `json:"data"`
	}
}

type TeenpattiRobot struct {
	*robot.BaseMatchRobot
	SeatId          int
	TokenPlayerSeat int
	round           int
}

const (
	compare = 10
	raise   = 9
	call    = 8
)

var action = []int{compare, raise, call}

func NewRobot(setting robot.RobotConfig) *TeenpattiRobot {
	matchRobot := robot.NewMatchRobot(setting)
	t := &TeenpattiRobot{
		BaseMatchRobot: matchRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	t.init()
	return t
}

func (t *TeenpattiRobot) init() {
	t.SeatId = -1
	t.round = 0
}

func (t *TeenpattiRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.BaseMatchRobot.CheckMatchCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case "JoinGame":
		t.CheckJoinGame(response)
	case "ActSettleData":
		t.init()
		t.CheckMatchPlayCount()
	case "ActGamePeriod":
		if t.Fsm == "Play" {
			t.getRound(response)
		}
	case "ActCompareSeat":
		t.AcceptCompare(response)

	case "ActTokenPlayerSeat":
		t.PlayCard(response)
	}

	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *TeenpattiRobot) CheckJoinGame(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.SeatId = int(data["own"].(float64))

}
func (t *TeenpattiRobot) getRound(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.round = int(data["Round"].(float64))

}

func (t *TeenpattiRobot) PlayCard(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.TokenPlayerSeat = int(info["SeatId"].(float64))
	if t.SeatId != t.TokenPlayerSeat {
		return
	}

	var instruction int
	if t.round < 3 {
		randIndex := rand.Intn(2) + 1
		instruction = action[randIndex]
	} else {
		randIndex := rand.Intn(3)
		instruction = action[randIndex]
	}

	data := PlayerAction{}
	data.PlayOperate.Instruction = instruction
	if instruction == compare {
		data.PlayOperate.Data = 0
	}

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *TeenpattiRobot) AcceptCompare(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	rivalSeat := int(info["rivalSeat"].(float64))
	if t.SeatId != rivalSeat {
		return
	}

	data := PlayerAction{}
	data.PlayOperate.Instruction = compare
	data.PlayOperate.Data = 2

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)

}
