package cockfight

import (
	"dytRobot/robot"
	"dytRobot/utils"
	"fmt"
	"math/rand"
	"time"
)

type CockfightRobot struct {
	*robot.BaseBetRobot
}

func NewRobot(setting robot.RobotConfig) *CockfightRobot {
	betRobot := robot.NewBetRobot(setting)
	t := &CockfightRobot{
		BaseBetRobot: betRobot,
	}
	t.Bet = t.goBet

	return t
}

func (t *CockfightRobot) goBet() {
	ticker := time.NewTicker(robot.BET_WAIT_SECOND * time.Second)
	for range ticker.C {
		if t.Fsm != robot.FSM_BET {
			t.BetArea = make([]int, len(t.BetLimit))
			return
		}

		if len(t.BetLimit) == 0 || len(t.Chips) == 0 {
			return
		}
		
		area := rand.Intn(3)
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
