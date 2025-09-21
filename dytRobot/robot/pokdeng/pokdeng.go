package pokdeng

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
	"fmt"
	"math/rand"
	"strconv"
)

type PokdengRobot struct {
	*robot.BaseMatchRobot
	Chips  []interface{}
	SeatId int
	Gold   float64
	Cards  []int
	Type   int
	Type2  int
}

const (
	CARDTYPE_POINT_0        int = iota // 點數 0 - (對子/同花)
	CARDTYPE_POINT_1                   // 點數 1 - (對子/同花)
	CARDTYPE_POINT_2                   // 點數 2 - (對子/同花)
	CARDTYPE_POINT_3                   // 點數 3 - (對子/同花)
	CARDTYPE_POINT_4                   // 點數 4 - (對子/同花)
	CARDTYPE_POINT_5                   // 點數 5 - (對子/同花)
	CARDTYPE_POINT_6                   // 點數 6 - (對子/同花)
	CARDTYPE_POINT_7                   // 點數 7 - (對子/同花)
	CARDTYPE_POINT_8                   // 點數 8 - (對子/同花)
	CARDTYPE_POINT_9                   // 點數 9 - (對子/同花)
	CARDTYPE_SANGONG                   //  三公  -
	CARDTYPE_STRAIGHT                  //  順子  - AKQ...432 (32A 不算順)
	CARDTYPE_STRAIGHT_FLUSH            // 同花順 - AKQ...432 (32A 不算順)
	CARDTYPE_TRIPLE                    //  三條  -
	CARDTYPE_POK_8                     //  博八  - 兩張牌 (對子/同花)
	CARDTYPE_POK_9                     //  博九  - 兩張牌 (同花)

	CARDTYPE_COUNT
)

const (
	ACTION_INIT  int = iota //
	ACTION_DRAW             // Action 補牌
	ACTION_STAND            // Action 不補牌

	ACTION_COUNT
)

var (
	RobotDraw = map[int]int{ // 補牌機率
		4: 90,
		5: 65,
		6: 35,
		7: 10,
	}
)

func NewRobot(setting robot.RobotConfig) *PokdengRobot {
	matchRobot := robot.NewMatchRobot(setting)
	t := &PokdengRobot{
		BaseMatchRobot: matchRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	return t
}

func (t *PokdengRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.CheckMatchCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case "JoinGame":
		t.CheckJoinGame(response)
	case "ActCardsInfo":
		t.GetCardsInfo(response)
	case "ActGamePeriod":
		t.CheckGamePeriod(response)
	case "ActSettleInfo":
		t.CheckMatchPlayCount()
	}

	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *PokdengRobot) CheckJoinGame(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.SeatId = int(data["OwnSeat"].(float64))
	t.Chips = data["Chips"].([]interface{})
	playerInfo := data["PlayerInfo"].([]interface{})
	for _, value := range playerInfo {
		info, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		if int(info["seatId"].(float64)) == t.SeatId {
			gold := info["gold"].(float64)
			t.Gold = gold
			// fmt.Printf("Seat %d, Gold %f\n", t.SeatId, t.Gold)
			break
		}
	}

}

func (t *PokdengRobot) CheckGamePeriod(response *utils.RespBase) {
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
		var data struct {
			MatchGameBet struct {
				BetInfo string `json:"BetInfo"`
			}
		}
		data.MatchGameBet.BetInfo = "1"

		len := len(t.Chips)
		index := rand.Intn(len)
		for index >= 0 {
			fmt.Printf("Chip[%d]  %f\n", index, t.Chips[index].(float64))
			if t.Chips[index].(float64) <= t.Gold {
				data.MatchGameBet.BetInfo = strconv.Itoa(int(t.Chips[index].(float64)))
				break
			}
			index--
		}
		// fmt.Printf("Bet %+v\n", data)
		utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)

	case "Draw":
		if t.Type >= CARDTYPE_POK_8 {
			return
		}
		var data struct {
			PlayOperate struct {
				Instruction int         `json:"instruction"`
				Data        interface{} `json:"data"`
			}
		}
		data.PlayOperate.Instruction = 16
		data.PlayOperate.Data = ACTION_STAND
		if t.Type < CARDTYPE_POINT_4 && t.Type2 == 0 {
			data.PlayOperate.Data = ACTION_DRAW
		} else if t.Type >= CARDTYPE_POINT_4 &&
			t.Type < CARDTYPE_POINT_8 {
			rate := rand.Intn(100)
			if rate < RobotDraw[t.Type] {
				data.PlayOperate.Data = ACTION_DRAW
			}
		}
		// fmt.Printf("Draw %+v\n", data)
		utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)

	}
}

func (t *PokdengRobot) GetCardsInfo(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.Type = int(data["Type"].(float64))
	t.Type2 = int(data["Type2"].(float64))
}
