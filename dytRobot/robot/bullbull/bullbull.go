package bullbull

import (
	"dytRobot/robot"
	"dytRobot/utils"
	"math/rand"
	"strconv"
)

type BullBullRobot struct {
	*robot.BaseMatchRobot
}

func NewRobot(setting robot.RobotConfig) *BullBullRobot {
	matchRobot := robot.NewMatchRobot(setting)
	t := &BullBullRobot{
		BaseMatchRobot: matchRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	return t
}

func (t *BullBullRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.CheckMatchCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case "ActRobInfo":
		t.SetRobBankerBet(response)
	case "ActBetInfo":
		t.SetMatchGameBet(response)
	case "ActGamePeriod":
		t.CheckGamePeriod(response)
	case "ActSettleInfo":
		t.CheckMatchPlayCount()
	}

	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *BullBullRobot) SetRobBankerBet(response *utils.RespBase) {
	var data struct {
		SetRobBankerBet struct {
			Rob int `json:"Rob"`
		}
	}
	data.SetRobBankerBet.Rob = -1

	if info, ok := response.Data.(map[string]interface{}); ok {
		robInfo := info["RobInfo"].([]interface{})
		if len := len(robInfo); len > 0 {
			index := rand.Intn(len)
			data.SetRobBankerBet.Rob = int(robInfo[index].(float64))
		}
	}

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BullBullRobot) SetMatchGameBet(response *utils.RespBase) {
	var data struct {
		MatchGameBet struct {
			BetInfo string `json:"BetInfo"`
		}
	}
	data.MatchGameBet.BetInfo = "1"

	if info, ok := response.Data.(map[string]interface{}); ok {
		betInfo := info["BetInfo"].([]interface{})
		if len := len(betInfo); len > 0 {
			index := rand.Intn(len)
			data.MatchGameBet.BetInfo = strconv.Itoa(int(betInfo[index].(float64)))
		}
	}

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BullBullRobot) CheckGamePeriod(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	fsm, ok := info["Fsm"].(string)
	if !ok {
		return
	}

	switch fsm {
	case "Result":
		var data struct {
			PlayOperate struct {
				Instruction int `json:"instruction"`
			}
		}
		data.PlayOperate.Instruction = 6
		utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
	}
}
