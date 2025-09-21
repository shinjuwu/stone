package goldenflower

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
	"fmt"
)

type PlayerAction struct {
	PlayOperate struct {
		Instruction int `json:"instruction"`
	}
}

type GoldenflowerRobot struct {
	*robot.BaseMatchRobot
	SeatId          int
	RaiseArray      [7]float64
	TokenPlayerSeat int
}

func NewRobot(setting robot.RobotConfig) *GoldenflowerRobot {
	matchRobot := robot.NewMatchRobot(setting)
	t := &GoldenflowerRobot{
		BaseMatchRobot: matchRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	t.init()
	return t
}

func (t *GoldenflowerRobot) init() {
	t.SeatId = -1
}

func (t *GoldenflowerRobot) GoCheckCommand(response *utils.RespBase) int {
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
		t.CheckGamePeriod(response)
	case "ActTokenPlayerSeat":
		t.PlayCard(response)
	case "ActRaiseArray":
		t.CheckRaiseArray(response)
	}

	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *GoldenflowerRobot) CheckJoinGame(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.SeatId = int(data["own"].(float64))
	fmt.Println(t.SeatId)
}

func (t *GoldenflowerRobot) CheckGamePeriod(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	fsm, ok := info["Fsm"].(string)
	if !ok {
		return
	}

	switch fsm {
	case "Play":
		//do nothing
	default:
		//do nothing
	}
}

func (t *GoldenflowerRobot) PlayCard(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.TokenPlayerSeat = int(info["SeatId"].(float64))
	if t.SeatId != t.TokenPlayerSeat {
		return
	}
	fmt.Println(info)
	data := PlayerAction{}
	data.PlayOperate.Instruction = 8

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}
func (t *GoldenflowerRobot) CheckRaiseArray(response *utils.RespBase) {
	// data, ok := response.Data.(map[string]interface{})
	// if !ok {
	// 	return
	// }
	// raiseArray := data["raiseArray"].([7]float64)
	// t.RaiseArray = raiseArray
}
