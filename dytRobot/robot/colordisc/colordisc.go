package colordisc

import (
	"dytRobot/robot"
	"dytRobot/utils"
	"fmt"
	"math/rand"
	"time"
)

type ColordiscRobot struct {
	*robot.BaseBetRobot
}

var areaRate = []float64{1.0, 1.0, 15.0, 3.0, 3.0, 15.0}

func NewRobot(setting robot.RobotConfig) *ColordiscRobot {
	betRobot := robot.NewBetRobot(setting)
	t := &ColordiscRobot{
		BaseBetRobot: betRobot,
	}
	t.Bet = t.goBet

	for _, value := range areaRate {
		weight := (int)(1 / value * 100)
		t.AreaWeight = append(t.AreaWeight, weight)
		t.TotalWeight += weight
	}

	return t
}

func (t *ColordiscRobot) goBet() {
	ticker := time.NewTicker(robot.BET_WAIT_SECOND * time.Second)
	for range ticker.C {
		if t.Fsm != robot.FSM_BET {
			t.BetArea = make([]int, len(t.BetLimit))
			return
		}

		if len(t.BetLimit) == 0 || len(t.Chips) == 0 {
			return
		}
		var area int
		weight := rand.Intn(t.TotalWeight)
		for i, value := range t.AreaWeight {
			if weight < value {
				area = i
				break
			}
			weight -= value
		}
		bet := int(t.Chips[rand.Intn(len(t.Chips))].(float64))
		limit := int(t.BetLimit[area].([]interface{})[1].(float64))
		if t.BetArea[area]+bet > limit {
			continue
		}
		t.BetArea[area] += bet

		betInfo := fmt.Sprintf("[{\"AreaID\":%d,\"Bet\":%d}]", area, bet)
		var data struct {
			Bet struct {
				BetInfo string
			}
		}
		data.Bet.BetInfo = betInfo

		if ok, _ := utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel); !ok {
			return
		}
	}
}
