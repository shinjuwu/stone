package texas

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
	"strconv"
)

const (
	ACT_ACTION_SEAT   = "ActActionSeat"
	ACT_ACTION_INFO   = "ActActionInfo"
	ACT_BESTCARD_INFO = "ActBestCardInfo"
	ACT_SETTLE_INFO   = "ActSettleInfo"
)

type TexasRobot struct {
	*robot.BaseMatchRobot
	// Gid       string
	SeatId    int
	RoundBet  float64
	EnterGold float64
	LockGold  float64
	CardType  int
	Cards     []int
	//Fsm        string
	NextAction int

	isSendJoin bool
}

const (
	CARDTYPE_HIGH_CARD       int = iota // 0, 散牌.
	CARDTYPE_ONE_PAIR                   // 1, 一對.
	CARDTYPE_TWO_PAIRS                  // 2, 兩對.
	CARDTYPE_THREE_OF_A_KIND            // 3, 三條.
	CARDTYPE_STRAIGHT                   // 4, 順子.
	CARDTYPE_FLUSH                      // 5, 同花.
	CARDTYPE_FULL_HOUSE                 // 6, 葫蘆.
	CARDTYPE_FOUR_OF_A_KIND             // 7, 四條／鐵支／金剛.
	CARDTYPE_STRAIGHT_FLUSH             // 8, 同花順. 花色相同的順子.
	CARDTYPE_ROYAL_FLUSH                // 9, 同花大順／皇家同花順.

	CARDTYPE_COUNT
)

const (
	ACTION_INIT int = iota
	ACTION_FOLD
	ACTION_CHECK
	ACTION_CALL
	ACTION_RAISE
	ACTION_ALL_IN
)

func NewRobot(setting robot.RobotConfig) *TexasRobot {
	matchRobot := robot.NewMatchRobot(setting)
	t := &TexasRobot{
		BaseMatchRobot: matchRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	return t
}

func (t *TexasRobot) init() {
	t.SeatId = -1
	t.RoundBet = 0
	t.NextAction = ACTION_CALL
}

func (t *TexasRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.CheckMatchCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case robot.RET_JOIN_GAME:
		t.init()
		t.CheckJoinGame(response)
	/* case robot.ACT_GAME_PERIOD:
	t.GetGamePeriod(response) */
	case robot.ACT_GOLD:
		t.SetLockGold(response)

	case "ActSelfCard":
		t.isSendJoin = false
	//	t.GetSelfCard(response)
	case ACT_BESTCARD_INFO:
		t.GetCardInfo(response)
	case ACT_ACTION_SEAT:
		t.SetAction(response)
	case ACT_ACTION_INFO:
		t.CheckActionInfo(response)

	case ACT_SETTLE_INFO:
		/* t.Gid = ""
		t.Fsm = "" */
		if !t.isSendJoin { //避免棄牌離開後，收到ACT_SETTLE_INFO和RET_QUIT_GAME而送兩次JoinGame
			t.CheckMatchPlayCount()
			t.isSendJoin = true
		}

	case robot.RET_QUIT_GAME:
		if !t.isSendJoin { //避免棄牌離開後，收到ACT_SETTLE_INFO和RET_QUIT_GAME而送兩次JoinGame
			t.CheckMatchPlayCount()
			t.isSendJoin = true
		}
	}

	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *TexasRobot) CheckJoinGame(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.SeatId = int(data["OwnSeat"].(float64))
	playerInfo := data["PlayerInfo"].([]interface{})
	for _, value := range playerInfo {
		info, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		if int(info["seatId"].(float64)) == t.SeatId {
			gold := info["gold"].(float64)
			t.LockGold = gold
			t.EnterGold = gold
		}
	}
}

func (t *TexasRobot) GetSelfCard(response *utils.RespBase) {
	if info, ok := response.Data.(map[string]interface{}); ok {
		// t.Cards = info["Cards"].([]int)
		data := info["Cards"].([]interface{})
		for _, card := range data {
			t.Cards = append(t.Cards, int(card.(float64)))
		}
	}
}

func (t *TexasRobot) GetCardInfo(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	cardType, ok := info["Type"].(float64)
	if !ok {
		return
	}
	t.CardType = int(cardType)
	pair := -1

	if t.CardType == CARDTYPE_ONE_PAIR /* || t.CardType == CARDTYPE_TWO_PAIRS */ {
		data := info["Cards"].([]interface{})
		for _, card := range data {
			temp := int(card.(float64)) / 13
			if t.CardType == CARDTYPE_ONE_PAIR {
				pair = temp
				break
			} else if t.CardType == CARDTYPE_TWO_PAIRS {
				if temp > pair {
					pair = temp
				}
			}
		}
		if pair == 0 {
			pair = 13
		}
	}

	if t.CardType >= CARDTYPE_THREE_OF_A_KIND {
		t.NextAction = ACTION_ALL_IN

	} else if t.CardType == CARDTYPE_HIGH_CARD {
		if t.Fsm == "PreFlop" || t.Fsm == "FlopRound" || t.Fsm == "TurnRound" {
			t.NextAction = ACTION_CALL
		} else if t.Fsm == "RiverRound" {
			t.NextAction = ACTION_FOLD
		}

	} else if t.CardType == CARDTYPE_ONE_PAIR {
		if pair < 3 && (t.Fsm == "TurnRound" || t.Fsm == "RiverRound") {
			t.NextAction = ACTION_FOLD
		} else {
			t.NextAction = ACTION_CALL
		}

	} else if t.CardType == CARDTYPE_TWO_PAIRS {
		t.NextAction = ACTION_CALL
	}
}

func (t *TexasRobot) SetAction(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	seat, ok := info["SeatId"].(float64)
	if !ok {
		return
	}

	if int(seat) != t.SeatId {
		return
	}
	call, ok := info["Call"].(float64)
	if !ok {
		return
	}
	if t.NextAction == ACTION_FOLD {
		t.ActionFold()
		return
	} else if t.NextAction == ACTION_CALL {
		if t.CardType == CARDTYPE_ONE_PAIR {
			if call > (t.EnterGold / 4) {
				t.ActionFold()
				return
			}
		}
		if call > t.LockGold {
			call = t.LockGold
		}
	} else if t.NextAction == ACTION_ALL_IN {
		call = t.LockGold
	}

	var data struct {
		MatchGameBet struct {
			BetInfo string `json:"BetInfo"`
		}
	}
	data.MatchGameBet.BetInfo = strconv.FormatFloat(call, 'f', 4, 64)
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *TexasRobot) ActionFold() {
	var data struct {
		PlayOperate struct {
			Instruction int `json:"instruction"`
		}
	}
	data.PlayOperate.Instruction = 7
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

/* func (t *TexasRobot) GetGamePeriod(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.Fsm, ok = info["Fsm"].(string)
	if !ok {
		return
	}
	if t.Fsm == "Match" {
		t.Gid, ok = info["Gid"].(string)
	} else if t.Fsm == "Settle" {
		// fmt.Printf("%s %s ______\n", t.Gid, t.Fsm)
	}
} */

func (t *TexasRobot) CheckActionInfo(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	seat, ok := info["SeatId"].(float64)
	if !ok {
		return
	}

	if int(seat) != t.SeatId {
		return
	}
	action, ok := info["Action"].(string)
	if !ok {
		return
	}
	if action == "Fold" { // 自己棄牌, 提早離桌.
		t.QuitGame()
	}
}

func (t *TexasRobot) CheckBet() {
	var data struct {
		PlayOperate struct {
			Instruction int `json:"instruction"`
		}
	}
	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

func (t *TexasRobot) SetLockGold(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.LockGold = data["Gold"].(float64)
}
