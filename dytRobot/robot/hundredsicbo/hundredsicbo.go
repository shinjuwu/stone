package hundredsicbo

import (
	"dytRobot/robot"
	"dytRobot/utils"
	"fmt"
	"math/rand"
	"time"
)

const (
	BET_AREA_ODD   = 21 // 單
	BET_AREA_EVEN  = 22 // 雙
	BET_AREA_BIG   = 23 // 大
	BET_AREA_SMALL = 24 // 小
)

var areaRate = []float64{
	30,                           // 任意豹子
	180, 180, 180, 180, 180, 180, // 指定豹子
	60, 30, 18, 12, 8, 6, 6, // 4 ~ 10 點
	6, 6, 8, 12, 18, 30, 60, // 11 ~ 17 點
	2, 2, // 單雙
	2, 2} // 大小

type HundredsicboRobot struct {
	*robot.BaseBetRobot
}

func NewRobot(setting robot.RobotConfig) *HundredsicboRobot {
	betRobot := robot.NewBetRobot(setting)
	t := &HundredsicboRobot{
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

func (t *HundredsicboRobot) goBet() {
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
		//單雙不能同時下
		if area == BET_AREA_ODD && t.BetArea[BET_AREA_EVEN] > 0 {
			continue
		}
		if area == BET_AREA_EVEN && t.BetArea[BET_AREA_ODD] > 0 {
			continue
		}
		//大小不能同時下
		if area == BET_AREA_BIG && t.BetArea[BET_AREA_SMALL] > 0 {
			continue
		}
		if area == BET_AREA_SMALL && t.BetArea[BET_AREA_BIG] > 0 {
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
