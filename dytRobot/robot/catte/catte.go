package catte

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/utils"
	"sort"
)

const (
	FOLD    = 14
	COMPARE = 10
)

type PlayerAction struct {
	PlayOperate struct {
		Instruction int         `json:"instruction"`
		Data        interface{} `json:"data"`
	}
}

type CatteRobot struct {
	*robot.BaseMatchRobot
	ownSeatId    int
	currentSeat  int
	ownHandcards []int
	roundMaxCard int
	round        int
}

func NewRobot(setting robot.RobotConfig) *CatteRobot {
	matchRobot := robot.NewMatchRobot(setting)
	t := &CatteRobot{
		BaseMatchRobot: matchRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	return t
}

func (t *CatteRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.BaseMatchRobot.CheckMatchCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case "JoinGame":
		t.CheckJoinGame(response)
	case "ActSettleData":
		t.CheckMatchPlayCount()
	case "ActDealCard":
		t.getHandCards(response)
	case "ActGamePeriod":
		if t.Fsm == "Play" || t.Fsm == "PlayRoundSix" {
			t.getRound(response)
		}
	// case "ActTokenPlayerSeat":
		// t.PlayCard(response)
	case "ActAction":
		t.getMaxCard(response)
	}
	return robot.RESPONSE_EXCUTED_SUCCESS
}
func (t *CatteRobot) getHandCards(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}

	var handCards []int
	for _, cardNum := range data["cards"].([]interface{}) {
		c := int(cardNum.(float64))
		handCards = append(handCards, c)
	}
	sort.Ints(handCards)
	t.ownHandcards = handCards

}

func (t *CatteRobot) CheckJoinGame(response *utils.RespBase) {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.ownSeatId = int(data["OwnSeat"].(float64))

}

func (t *CatteRobot) PlayCard(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.currentSeat = int(info["SeatId"].(float64))
	if t.ownSeatId != t.currentSeat {
		return
	}

	var cardNum int
	var actionType int
	data := PlayerAction{}
	switch {
	case t.round <= 4:
		cardNum, actionType = t.getCard()
		t.tidyHandCards(cardNum)
		data.PlayOperate.Instruction = actionType
		data.PlayOperate.Data = cardNum
	case t.round == 5:
		data.PlayOperate.Instruction = t.ownHandcards[0]
		data.PlayOperate.Data = COMPARE
	case t.round == 6:
		data.PlayOperate.Instruction = t.ownHandcards[1]
		data.PlayOperate.Data = COMPARE
	}

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}

// 找出牌
func (t *CatteRobot) getCard() (card int, actionType int) {
	playCard := -1
	maxCardSuit := cardSuit(t.roundMaxCard)
	maxCardValue := cardValue(t.roundMaxCard)

	for _, card := range t.ownHandcards {
		if cardSuit(card) == maxCardSuit {
			if cardValue(card) > maxCardValue {
				playCard = card
			}
		}
	}

	if playCard != -1 {
		return playCard, COMPARE
	} else {
		return t.ownHandcards[0], FOLD
	}

}
func (t *CatteRobot) getMaxCard(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	actiontype := data["actiontype"].(string)
	if actiontype == "compare" {
		t.roundMaxCard = int(data["card"].(float64))

	}
}

func (t *CatteRobot) getRound(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.round = int(data["Round"].(float64))

}

func cardSuit(card int) int {
	return int(card / 13)
}

func cardValue(card int) int {
	value := int(card % 13)
	if value == 0 {
		return 13 // set a larger value to Ace
	}
	return value
}

func (t *CatteRobot) tidyHandCards(Card int) {
	cardIndex := sort.SearchInts(t.ownHandcards, Card)
	t.ownHandcards = append(t.ownHandcards[:cardIndex], t.ownHandcards[cardIndex+1:]...)
}
