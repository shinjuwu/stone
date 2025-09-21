package blackjack

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
	"fmt"
	"math/rand"
)

type PlayerAction struct {
	PlayOperate struct {
		Instruction int `json:"instruction"`
	}
}

type BlackJackRobot struct {
	*robot.BaseMatchRobot
	SeatId          int
	Chips           []interface{}
	Point           int
	PointSpilt      int
	CanSpilt        bool
	SpiltCount      int
	TokenPlayerSeat int
}

func NewRobot(setting robot.RobotConfig) *BlackJackRobot {
	matchRobot := robot.NewMatchRobot(setting)
	t := &BlackJackRobot{
		BaseMatchRobot: matchRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	t.init()
	return t
}

func (t *BlackJackRobot) init() {
	t.SeatId = -1
	t.Point = 0
	t.PointSpilt = 0
	t.CanSpilt = false
	t.SpiltCount = 0
}

func (t *BlackJackRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.BaseMatchRobot.CheckMatchCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case "JoinGame":
		t.CheckJoinGame(response)
	case "ActMatchGameResult":
		t.init()
		t.CheckMatchPlayCount()
	case "ActGamePeriod":
		t.CheckGamePeriod(response)
	case "MatchGameBet":
		t.SendBetFinish(response)
	case "ActDealCard":
		t.CheckCardAndPoint(response)
	case "ActTokenPlayerSeat":
		t.PlayCard(response)
	case "ActSplitCard":
		t.CheckSpiltCard(response)
	case "ActHitCard":
		t.CheckHit(response)
	}

	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *BlackJackRobot) CheckJoinGame(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.Chips = data["Chips"].([]interface{})
	playerInfo := data["PlayerInfo"].([]interface{})
	for _, value := range playerInfo {
		info, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		name := t.LoginName
		if !t.AccessDcc {
			name = "Test_" + t.LoginName
		}

		if info["name"].(string) == name {
			t.SeatId = int(info["seatId"].(float64))
			break
		}
	}
}

func (t *BlackJackRobot) CheckGamePeriod(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	fsm, ok := info["Fsm"].(string)
	if !ok {
		return
	}

	switch fsm {
	case "Bet":
		if len(t.Chips) == 0 || t.SeatId == -1 {
			return
		}
		bet := int(t.Chips[rand.Intn(len(t.Chips))].(float64))
		betInfo := fmt.Sprintf("[{\"SeatId\":%d,\"Bet\":%d}]", t.SeatId, bet)
		var data struct {
			MatchGameBet struct {
				BetInfo string
			}
		}
		data.MatchGameBet.BetInfo = betInfo
		utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
	case "Insurance":
		if t.SeatId == -1 {
			return
		}
		var data struct {
			SetInsurance struct {
				seatId int
				isbuy  bool
			}
		}
		data.SetInsurance.seatId = t.SeatId
		if rand.Intn(2) == 0 {
			data.SetInsurance.isbuy = true
		}
		utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
	case "Split":
		if t.TokenPlayerSeat != t.SeatId {
			return
		}
		point := t.Point
		if t.CanSpilt {
			if t.SpiltCount == 1 {
				point = t.PointSpilt
			} else {
				t.SpiltCount++
			}
		}

		data := PlayerAction{}
		if point == 21 {
			return
		} else if point > 17 {
			//Stop
			data.PlayOperate.Instruction = 4
		} else if point < 9 {
			//Hit
			data.PlayOperate.Instruction = 1
		} else {
			if rand.Intn(2) == 0 {
				//Double
				data.PlayOperate.Instruction = 3
			} else {
				//Hit
				data.PlayOperate.Instruction = 1
			}
		}
		utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
	}
}

func (t *BlackJackRobot) SendBetFinish(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}
	data := PlayerAction{}
	data.PlayOperate.Instruction = 5
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BlackJackRobot) CheckCardAndPoint(response *utils.RespBase) {
	if t.SeatId == -1 {
		return
	}
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	points := data["point"].([]interface{})
	t.Point = int(points[t.SeatId].(float64))

	t.CanSpilt = false
	cards, ok := data["playerCards"].([]interface{})
	if !ok {
		return
	}
	deailCard, ok := cards[t.SeatId].([]interface{})
	if !ok {
		return
	}
	card1 := int(deailCard[0].(float64))
	card2 := int(deailCard[1].(float64))
	if card1%13 == card2%13 {
		t.CanSpilt = true
	}
}

func (t *BlackJackRobot) PlayCard(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.TokenPlayerSeat = int(info["SeatId"].(float64))
	if t.SeatId != t.TokenPlayerSeat {
		return
	}

	data := PlayerAction{}
	if t.Point == 21 {
		return
	} else if t.CanSpilt {
		data.PlayOperate.Instruction = 2
	} else if t.Point > 17 {
		//Stop
		data.PlayOperate.Instruction = 4
	} else if t.Point <= 9 {
		//Hit
		data.PlayOperate.Instruction = 1
	} else {
		if rand.Intn(2) == 0 {
			//Double
			data.PlayOperate.Instruction = 3
		} else {
			//Hit
			data.PlayOperate.Instruction = 1
		}
	}

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BlackJackRobot) CheckHit(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	if t.SeatId != int(info["seatId"].(float64)) {
		return
	}
	if info["isFiveDragon"].(bool) {
		return
	}
	point := int(info["point"].(float64))
	if point >= 21 {
		return
	}

	data := PlayerAction{}
	if point > 17 {
		//Stop
		data.PlayOperate.Instruction = 4
	} else {
		//Hit
		data.PlayOperate.Instruction = 1
	}
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *BlackJackRobot) CheckSpiltCard(response *utils.RespBase) {
	if t.TokenPlayerSeat != t.SeatId {
		return
	}
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	points := data["point"].([]interface{})
	t.Point = int(points[0].(float64))
	t.PointSpilt = int(points[1].(float64))
}
