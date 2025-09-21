package sangong

import (
	"dytRobot/robot"
	"dytRobot/utils"
	"math/rand"
	"strconv"
)

type SangongRobot struct {
	*robot.BaseMatchRobot
}

func NewRobot(setting robot.RobotConfig) *SangongRobot {
	matchRobot := robot.NewMatchRobot(setting)
	t := &SangongRobot{
		BaseMatchRobot: matchRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	return t
}

func (t *SangongRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.CheckMatchCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case "ActBetInfo":
		t.SetMatchGameBet(response)
	case "ActGamePeriod":
		t.CheckGamePeriod(response)
	case "ActSettleInfo":
		t.CheckMatchPlayCount()
	}

	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *SangongRobot) SetMatchGameBet(response *utils.RespBase) {
	var data struct {
		MatchGameBet struct {
			BetInfo string `json:"BetInfo"`
		}
	}
	data.MatchGameBet.BetInfo = "7"

	if info, ok := response.Data.(map[string]interface{}); ok {
		betInfo := info["BetInfo"].([]interface{})
		if len := len(betInfo); len > 0 {
			index := rand.Intn(len)
			data.MatchGameBet.BetInfo = strconv.Itoa(int(betInfo[index].(float64)))
		}
	}

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *SangongRobot) CheckGamePeriod(response *utils.RespBase) {
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
	case "Rob":
		var data struct {
			SetRobBanker struct {
				Rob bool `json:"Rob"`
			}
		}
		data.SetRobBanker.Rob = true

		if rand.Intn(2) == 0 {
			data.SetRobBanker.Rob = false
		}

		utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
	}
}
