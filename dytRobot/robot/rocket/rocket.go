package rocket

import (
	"dytRobot/robot"
	"dytRobot/utils"
	"fmt"
	"math/rand"
	"time"
)

type RocketRobot struct {
	*robot.BaseBetRobot
}

func NewRobot(setting robot.RobotConfig) *RocketRobot {
	betRobot := robot.NewBetRobot(setting)
	t := &RocketRobot{
		BaseBetRobot: betRobot,
	}
	t.Bet = t.goBet
	t.CheckCommand = t.GoCheckCommand

	return t
}

func (t *RocketRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.BaseBetRobot.CheckBetCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}
	switch response.Ret {
	case "Forecast":
		t.goFlee()
	}
	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *RocketRobot) goFlee() {
	if len(t.BetArea) == 0 {
		return
	}
	if t.BetArea[0] == 0 {
		return
	}
	fleeTime := rand.Intn(20) + 3
	var data struct {
		BetOp struct {
			Instruction int
		}
	}
	data.BetOp.Instruction = 1
	time.Sleep(time.Duration(fleeTime) * time.Second)
	if t.Fsm != "Open" {
		return
	}
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *RocketRobot) goBet() {
	waitSec := rand.Intn(robot.BET_WAIT_SECOND)
	time.Sleep(time.Duration(waitSec) * time.Second)
	if t.Fsm != robot.FSM_BET {
		t.BetArea = make([]int, len(t.BetLimit))
		return
	}

	if len(t.BetLimit) == 0 || len(t.Chips) == 0 {
		return
	}
	// robotName := "tangRobot"
	// number, _ := strconv.Atoi(strings.Replace(t.LoginName, robotName, "", -1))
	// autoFleePayout := float64((number%20 + 1) * 5)

	autoFleePayout := float64(rand.Intn(10000) / 100.0)
	area := 0
	bet := int(t.Chips[rand.Intn(len(t.Chips))].(float64))
	// bet := int(t.Chips[(len(t.Chips))-1].(float64))
	limit := int(t.BetLimit[area].([]interface{})[1].(float64))
	if bet > limit {
		return
	}
	// t.BetArea[area] += bet

	betInfo := fmt.Sprintf("{\"Bet\":%d,\"AutoFleePayout\":%.2f}", bet, autoFleePayout)
	var data struct {
		Bet struct {
			BetInfo string
		}
	}
	// fmt.Printf("Name : %s, BetInfo : %s\n", t.LoginName, betInfo)
	data.Bet.BetInfo = betInfo
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}
