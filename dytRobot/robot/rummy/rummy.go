package rummy

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
	"strconv"
)

type RummyRobot struct {
	*robot.BaseMatchRobot
	Gid        string
	SeatId     int
	Gold       float64
	Cards      []int
	Wild       int
	WildPoint  int
	Deck       int
	Card       int
	DiscardWin int
	Others     []int
	Hints      map[int]int
	Win        int
	Lose       int
}

const (
	UNDEFINED = -1
)

func NewRobot(setting robot.RobotConfig) *RummyRobot {
	matchRobot := robot.NewMatchRobot(setting)
	t := &RummyRobot{
		BaseMatchRobot: matchRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	return t
}

func (t *RummyRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.CheckMatchCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case "JoinGame":
		t.CheckJoinGame(response)
	case "ActGamePeriod":
		t.CheckGamePeriod(response)
	case "ActDealCard":
		t.DealCard(response)
	case "ActPickCard":
		t.PickCard(response)
	/* case "ActDiscard":
		t.Discard(response)
	case "ActDrop":
		t.Drop(response) */
	case "ActPickCardWin":
		t.PickCardWin(response)
	case "ActBestCard":
		t.BestCard(response)
	case "ActSettleInfo":
		t.Settle(response)
		t.CheckMatchPlayCount()
	}

	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *RummyRobot) CheckJoinGame(response *utils.RespBase) {
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
	t.SeatId = int(data["OwnSeat"].(float64))
	t.DiscardWin = UNDEFINED
	t.Deck = 0
}

func (t *RummyRobot) CheckGamePeriod(response *utils.RespBase) {
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

	if int(seat) != t.SeatId {
		return
	}
	switch fsm {
	case "PickCard":
		var data struct {
			PlayOperate struct {
				Instruction int         `json:"instruction"`
				Data        interface{} `json:"data"`
			}
		}
		data.PlayOperate.Instruction = 13
		data.PlayOperate.Data = t.Deck
		utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)

	case "Discard":
		var data struct {
			PlayOperate struct {
				Instruction int         `json:"instruction"`
				Data        interface{} `json:"data"`
			}
		}
		data.PlayOperate.Instruction = 14
		if t.DiscardWin != UNDEFINED {
			data.PlayOperate.Data = t.DiscardWin
		} else {
			data.PlayOperate.Data = t.Card
		}
		utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)

	}
}

func (t *RummyRobot) DealCard(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.Wild = int(data["Wild"].(float64))
	if t.Wild >= 52 {
		t.WildPoint = 0
	} else {
		t.WildPoint = t.Wild % 13
	}
	if cards, ok := data["Cards"].([]interface{}); ok {
		t.Cards = []int{}
		for _, value := range cards {
			card := int(value.(float64))
			t.Cards = append(t.Cards, card)
		}
		// fmt.Printf("Seat [%d], WildPoint %d, Cards %+v\n", t.SeatId, t.WildPoint, t.Cards)
	}
}

func (t *RummyRobot) PickCard(response *utils.RespBase) {
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
	/* var str string
	if seat == t.SeatId {
		str = " vvvvvvvvvvv"
	}
	if t.DiscardWin != UNDEFINED {
		str = " .......... [Win]"
	}
	fmt.Printf("Seat [%d] PickCard[%d] %d%s\n", seat, index, card, str) */
	if seat == t.SeatId {
		t.Card = card // pickcard
		// t.Card = UNDEFINED
	}
}

/* func (t *RummyRobot) Discard(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	seat := int(data["SeatId"].(float64))
	card := int(data["Card"].(float64))
	var str string
	if seat == t.SeatId {
		str = " ^^^^^^^^^^"
	}
	if t.DiscardWin != UNDEFINED {
		str = " .......... [Win]"
	}
	fmt.Printf("Seat [%d] Discard %d%s\n", seat, card, str)
}

func (t *RummyRobot) Drop(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	seat := int(data["SeatId"].(float64))
	fmt.Printf("Seat [%d] Drop...\n", seat)
} */

func (t *RummyRobot) PickCardWin(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	if gid, ok := data["Gid"].(string); ok {
		t.Gid = gid
	}
	t.Card = int(data["Card"].(float64)) // pickcard win
	t.Deck = int(data["Index"].(float64))
	if discard, ok := data["Discard"].(float64); ok {
		t.DiscardWin = int(discard)
	}

	// fmt.Printf("PickCard [%d] %d discard %d .......... [Win]\n", t.Deck, t.Card, t.DiscardWin)
}

func (t *RummyRobot) BestCard(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	/* if point, ok := data["Points"].(float64); ok {
		fmt.Printf("Seat [%d] Point %d, Best:%+v\n", t.SeatId, int(point), data)
	} */
	if info, ok := data["Hints"].(map[string]interface{}); ok {
		// hints := map[string]int{}
		max := 0
		for card, value := range info {
			num := int(value.(float64))
			if max < num {
				max = num
				t.Card, _ = strconv.Atoi(card)
			}

		}
		// fmt.Printf("Seat [%d] Hints %+v, discard %2d\n", t.SeatId, info, t.Card)
		return
	} else if cards, ok := data["Others"].([]interface{}); ok {
		items := []string{"Runs", "ImpureRuns", "Sets", "ImpureSets"}
		others := []int{}
		for _, value := range cards {
			card := int(value.(float64))
			others = append(others, card)
		}
		cardNum := len(others)
		for _, item := range items {
			if datas, ok := data[item].([]interface{}); ok {
				for _, cards := range datas {
					cardNum += len(cards.([]interface{}))
				}
			}
		}

		if cardNum == 14 {
			index := 0
			for {
				/* index := rand.Intn(len(others))
				if others[index] >= 52 {
					continue
				}
				if (others[index] % 13) != t.WildPoint {
					t.Card = others[index] // discard
					break
				} */
				if others[index] >= 52 {
					continue
				}
				if (others[index] % 13) != t.WildPoint {
					t.Card = others[index] // discard
					break
				}
				index++
			}
			// fmt.Printf("Seat [%d] Others %+v, discard[%d] %d\n", t.SeatId, others, index, t.Card)
		}
	}

	/* pickcard, ok := data["PickCard"].(float64)
	if !ok {
		return
	}
	if int(pickcard) == t.Card {
		if discard, ok := data["Discard"].(int); ok {
			t.DiscardWin = discard
		}
	} */
	if discard, ok := data["Discard"].(float64); ok {
		t.DiscardWin = int(discard)
	}
}

func (t *RummyRobot) Settle(response *utils.RespBase) {
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
}
