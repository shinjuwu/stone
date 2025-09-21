package baccarat

import (
	"dytRobot/robot"
	"dytRobot/utils"
	"fmt"
	"math/rand"
	"time"
)

const (
	CAN_BET_BIG_AND_SMALL_ROUND = 30 // 30 局後不可壓大小

	BET_AREA_BIGGER = 7 // 大
	BET_AREA_SMALL  = 8 // 小
)

type BaccaratRobot struct {
	*robot.BaseBetRobot

	Round int //局數
}

var areaRate = []float64{1.95, 2, 1.95, 1.95, 9, 12, 12, 1.54, 2.5}

func NewRobot(setting robot.RobotConfig) *BaccaratRobot {
	betRobot := robot.NewBetRobot(setting)
	t := &BaccaratRobot{
		BaseBetRobot: betRobot,
	}
	t.Bet = t.goBet
	t.CheckBet = t.CheckBaccaratBet

	for _, value := range areaRate {
		weight := (int)(1 / value * 100)
		t.AreaWeight = append(t.AreaWeight, weight)
		t.TotalWeight += weight
	}

	return t
}

func (t *BaccaratRobot) CheckBaccaratBet(response *utils.RespBase) bool {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return false
	}

	round, ok := info["Round"].(float64)
	if !ok {
		return false
	}
	t.Round = int(round)
	return true
}

func (t *BaccaratRobot) goBet() {
	ticker := time.NewTicker(robot.BET_WAIT_SECOND * time.Second)
	for range ticker.C {
		if t.Fsm != "Bet" {
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
		//超過30局不可壓大小
		if t.Round > CAN_BET_BIG_AND_SMALL_ROUND && area >= BET_AREA_BIGGER {
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
