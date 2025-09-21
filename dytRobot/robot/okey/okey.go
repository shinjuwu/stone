package okey

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
	"math/rand"
)

type OkeyRobot struct {
	*robot.BaseMatchRobot
	Gid         string
	Seat        int
	PrevSeat    int
	PrevDiscard int
	Gold        float64
	HoleCards   []int
	Deck        int
	Card        int
	DiscardWin  []int
	Others      []int
	Hints       map[int]int

	Win  int
	Lose int
}

const (
	UNDEFINED = -1

	ACTION_PICKCARD    = 13
	ACTION_DISCARD     = 14
	ACTION_DISCARD_WIN = 15
)

const (
	ACT_TABLE_STATUS = "ActTableStatus"
	ACT_HOLE_CARD    = "ActHoleCard"
	ACT_PICKCARD     = "ActPickCard"
	ACT_DISCARD      = "ActDiscard"
	ACT_SETTLE_INFO  = "ActSettleInfo"
)

func NewRobot(setting robot.RobotConfig) *OkeyRobot {
	matchRobot := robot.NewMatchRobot(setting)
	t := &OkeyRobot{
		BaseMatchRobot: matchRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	return t
}

func (t *OkeyRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.CheckMatchCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case "JoinGame":
		t.CheckJoinGame(response)
	case "ActGamePeriod":
		t.CheckGamePeriod(response)
	case ACT_HOLE_CARD:
		t.GetHoleCard(response)
	case "ActPickCard":
		t.PickCard(response)
	case "ActDiscard":
		t.Discard(response)
	case "ActSettleInfo":
		t.Settle(response)
		t.CheckMatchPlayCount()
	}

	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *OkeyRobot) CheckJoinGame(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	if gid, ok := data["Gid"].(string); ok {
		t.Gid = gid
	}
	t.Seat = int(data["OwnSeat"].(float64))
	t.Card = UNDEFINED
	t.Deck = 0
}

func (t *OkeyRobot) CheckGamePeriod(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	fsm, ok := info["Fsm"].(string)
	if !ok {
		return
	}

	if gid, ok := info["Gid"].(string); ok {
		t.Gid = gid
	}

	seat, ok := info["SeatId"].(float64)
	if !ok {
		return
	}

	if int(seat) != t.Seat {
		return
	}
	switch fsm {
	case "Play":
		if len(t.HoleCards) == 14 {
			t.DoPickCard()
		} else if len(t.HoleCards) == 15 {
			t.DoDiscard()
		}
	}
}

func (t *OkeyRobot) DoPickCard() {
	var data struct {
		PlayOperate struct {
			Instruction int         `json:"instruction"`
			Data        interface{} `json:"data"`
		}
	}

	data.PlayOperate.Instruction = ACTION_PICKCARD
	data.PlayOperate.Data = t.Deck

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}
func (t *OkeyRobot) DoDiscard() {
	var data struct {
		PlayOperate struct {
			Instruction int         `json:"instruction"`
			Data        interface{} `json:"data"`
		}
	}
	if len(t.DiscardWin) > 0 {
		data.PlayOperate.Instruction = ACTION_DISCARD_WIN
		data.PlayOperate.Data = t.DiscardWin[0]
	} else {
		data.PlayOperate.Instruction = ACTION_DISCARD
		if t.Card == UNDEFINED && len(t.HoleCards) == 15 {
			t.Card = t.HoleCards[rand.Intn(15)]
		}
		data.PlayOperate.Data = t.Card
	}
	// fmt.Printf("Discard %+v\n", data)
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *OkeyRobot) GetHoleCard(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}

	items := []string{"Runs", "ImRuns", "Sets", "ImSets", "Pairs", "ImPairs", "CouldRuns", "CouldSets", "Others"}
	t.HoleCards = []int{}

	for _, item := range items {
		if info, ok := data[item].([]interface{}); ok {
			if item == "Others" {
				for _, value := range info {
					card := int(value.(float64))
					t.HoleCards = append(t.HoleCards, card)
				}
			} else {
				if len(info) == 0 {
					continue
				}
				for _, cards := range info {
					if cards == nil {
						break
					}
					for _, value := range cards.([]interface{}) {
						card := int(value.(float64))
						t.HoleCards = append(t.HoleCards, card)
					}
				}
			}
		}
	}
	// fmt.Printf("Seat [%d] Card %d %v\n", t.Seat, len(t.HoleCards), t.HoleCards)
}

func (t *OkeyRobot) PickCard(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	if gid, ok := data["Gid"].(string); ok {
		t.Gid = gid
	}
	seat := int(data["SeatId"].(float64))
	// index := int(data["Index"].(float64))
	card := int(data["Card"].(float64))
	// num := int(data["Num"].(float64))

	// discardWinStr := ""
	if seat == t.Seat && card >= 0 {
		t.Card = card // pickcard
		t.HoleCards = append(t.HoleCards, card)

		t.DiscardWin = []int{}
		if cards, ok := data["DiscardWin"].([]interface{}); ok {
			for _, value := range cards {
				card := int(value.(float64))
				t.DiscardWin = append(t.DiscardWin, card)
			}
			if len(t.DiscardWin) > 0 {
				// discardWinStr = fmt.Sprintf("DiscardWin %v", t.DiscardWin)
			}
		}
	}

	/* fmt.Printf("Seat [%d] Deck[%d/%d] Card %d %s\n", seat, index, num, card, discardWinStr) // */

	if seat == t.Seat {
		t.DoDiscard()
	}
}

func (t *OkeyRobot) Discard(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	seat := int(data["SeatId"].(float64))
	discard := int(data["Card"].(float64))
	// str := ""
	if seat == t.Seat {
		for i, card := range t.HoleCards {
			if card == discard {
				t.HoleCards = append(t.HoleCards[:i], t.HoleCards[i+1:]...)
				break
			}
		}
		// str = "<---"
	} else if seat == t.PrevSeat {
		t.PrevDiscard = discard
	}

	// fmt.Printf("Seat [%d] Discard %d %s\n", seat, discard, str)
}

func (t *OkeyRobot) Settle(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	if gameEnd, ok := data["GameEnd"].(bool); ok {
		if gameEnd {
			// fmt.Printf("-----< GameEnd >-----")
			return
		}
	}

	/* if settle, ok := data["SettleInfo"].([]interface{}); ok {
		for _, seatInfo := range settle {
			if data, ok := seatInfo.(map[string]interface{}); ok {
				if seatId, ok := data["SeatId"].(float64); ok {
					if t.SeatId == int(seatId) {
						if win, ok := data["Win"].(float64); ok {
							winStr := "Win"
							if win <= 0 {
								winStr = "Lose"
								t.Lose++
							} else {
								t.Win++
							}
							fmt.Printf("-----< Settle [%s] ... %.2f >-----   %s ( %d / %d )\n", winStr, win, t.Gid, t.Win, t.Lose)
							return
						}
					}
				}
			}
		}
	} */
	// fmt.Printf("Settle\n")
}
