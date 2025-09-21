package chinesepoker

import (
	"dytRobot/robot"
	"dytRobot/utils"
)

type PlayerAction struct {
	PlayOperate struct {
		Instruction int         `json:"instruction"`
		Data        interface{} `json:"data"`
	}
}

type CardInfo struct {
	HandCardType int     `json:"handCardType"`
	CardsArray   [][]int `json:"cardsArray"`
	ArrayType    []int   `json:"arrayType"`
}

type ChinesepokerRobot struct {
	*robot.BaseMatchRobot
	handCardType   []int
	cardsArray     [][][]int
	cardsArrayType [][]int
}

func NewRobot(setting robot.RobotConfig) *ChinesepokerRobot {
	matchRobot := robot.NewMatchRobot(setting)
	t := &ChinesepokerRobot{
		BaseMatchRobot: matchRobot,
	}
	t.CheckCommand = t.GoCheckCommand
	return t
}

func (t *ChinesepokerRobot) GoCheckCommand(response *utils.RespBase) int {
	result := t.BaseMatchRobot.CheckMatchCommand(response)
	if result != robot.RESPONSE_NO_SUTIALBE {
		return result
	}

	switch response.Ret {
	case "ActCardInfo":
		t.getHandCard(response)
		t.PlayCard(response)
	}
	return robot.RESPONSE_EXCUTED_SUCCESS
}

func (t *ChinesepokerRobot) getHandCard(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	gameInfo, ok := (data["gameInfo"].(map[string]interface{}))
	if !ok {
		return
	}

	//擺牌
	pokerInfo, ok := gameInfo["pokerInfo"].([]interface{})
	if !ok {
		return
	}

	var handCardsArray [][][]int
	var transArrayType [][]int
	t.cardsArrayType = make([][]int, 0)
	t.handCardType = make([]int, 0)
	t.cardsArray = make([][][]int, 0)

	for i := 0; i < len(pokerInfo); i++ {
		var eachHandCardsArray [][]int
		eachArrayInfo := pokerInfo[i].(map[string]interface{})
		handCardType := int(eachArrayInfo["handCardType"].(float64))
		t.handCardType = append(t.handCardType, handCardType)

		cardsArray := eachArrayInfo["cardsArray"].([]interface{})
		for _, array := range cardsArray {
			var transArray []int
			for _, card := range array.([]interface{}) {
				transCard := int(card.(float64))
				transArray = append(transArray, transCard)
			}
			eachHandCardsArray = append(eachHandCardsArray, transArray)
		}
		handCardsArray = append(handCardsArray, eachHandCardsArray)

		if eachArrayInfo["arrayType"] != nil {
			var eachArrayType []int
			cardsArrayType := eachArrayInfo["arrayType"].([]interface{})
			for _, arrayType := range cardsArrayType {
				eachArrayType = append(eachArrayType, int(arrayType.(float64)))
			}
			transArrayType = append(transArrayType, eachArrayType)
		} else {
			transArrayType = append(transArrayType, []int{})
		}

	}
	t.cardsArray = handCardsArray
	t.cardsArrayType = transArrayType
}

func (t *ChinesepokerRobot) PlayCard(response *utils.RespBase) {

	info := CardInfo{
		HandCardType: t.handCardType[0],
		CardsArray:   t.cardsArray[0],
	}
	if t.handCardType[0] != 0 {
		info.ArrayType = nil
	} else {
		info.ArrayType = t.cardsArrayType[0]
	}
	data := PlayerAction{}
	data.PlayOperate.Instruction = 10
	data.PlayOperate.Data = info

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
}
