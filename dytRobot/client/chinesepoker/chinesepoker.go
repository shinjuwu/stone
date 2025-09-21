package chinesepoker

import (
	"dytRobot/client"
	"dytRobot/utils"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ChinesepokerClient struct {
	*client.BaseMatchClient

	playButtom    *widget.Button
	allCardArray  [][]*widget.Label
	allArrayType  []*widget.Button
	handCardsType []*widget.Label

	handCardType    []int
	cardsArray      [][][]int
	cardsArrayType  [][]int
	cardsArrayIndex int
}

const (
	WU_LONG       int = 0
	DUI_ZI        int = 1
	LIANG_DUI     int = 2
	SAN_TIAO      int = 3
	SHUN_ZI       int = 4
	TONG_HUA      int = 5
	HU_LU         int = 6
	TIE_ZHI       int = 7
	TONG_HUA_SHUN int = 8
)

func NewClient(setting client.ClientConfig) *ChinesepokerClient {
	matchClient := client.NewMatchClient(setting)
	t := &ChinesepokerClient{
		BaseMatchClient: matchClient,
	}
	t.CheckResponse = t.CheckChinesepokerResponse
	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[0,13,26,2,12,25]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)

	return t
}

func (t *ChinesepokerClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateMatchSection(c)
	t.CreateGameSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *ChinesepokerClient) CheckChinesepokerResponse(response *utils.RespBase) bool {
	if t.CheckMatchResponse(response) {
		return true
	}
	switch response.Ret {
	case "ActCardInfo":
		t.getHandCard(response)
		t.putValue()

	case "ActGamePeriod":
		if t.Fsm == "Result" {
			t.resetValue()
		}
		if t.Fsm == "Array" {
			t.SetButton(t.playButtom, true)
		}
	}
	return false
}

func (t *ChinesepokerClient) CreateGameSection(c *fyne.Container) {
	t.playButtom = widget.NewButton("比牌", func() { t.SendPlayOperate(10, t.cardsArrayIndex) })
	t.SetButton(t.playButtom, false)
	t.allCardArray = make([][]*widget.Label, 5)
	t.handCardsType = make([]*widget.Label, 5)
	for i := range t.allCardArray {
		t.allCardArray[i] = make([]*widget.Label, 3)
	}

	for groupIndex := 0; groupIndex < len(t.allCardArray); groupIndex++ {
		t.handCardsType[groupIndex] = widget.NewLabel("")
	}

	pokertype := widget.NewLabel("Type:")
	front := widget.NewLabel("前墩:")
	middle := widget.NewLabel("中墩:")
	tail := widget.NewLabel("後墩:")
	poker1 := widget.NewLabel("第一組牌:")
	poker2 := widget.NewLabel("第二組牌:")
	poker3 := widget.NewLabel("第三組牌:")
	poker4 := widget.NewLabel("第四組牌:")
	poker5 := widget.NewLabel("第五組牌:")

	for groupIndex := 0; groupIndex < len(t.allCardArray); groupIndex++ {
		for arrayIndex := 0; arrayIndex < 3; arrayIndex++ {
			t.allCardArray[groupIndex][arrayIndex] = widget.NewLabel("")
		}
	}

	t.allArrayType = make([]*widget.Button, 5)
	arrayType := widget.NewLabel("牌型:")
	t.allArrayType[0] = widget.NewButton("", func() { t.cardsArrayIndex = 0 })
	t.allArrayType[1] = widget.NewButton("", func() { t.cardsArrayIndex = 1 })
	t.allArrayType[2] = widget.NewButton("", func() { t.cardsArrayIndex = 2 })
	t.allArrayType[3] = widget.NewButton("", func() { t.cardsArrayIndex = 3 })
	t.allArrayType[4] = widget.NewButton("", func() { t.cardsArrayIndex = 4 })

	//新增手牌牌型

	pokerInfo1 := container.NewHBox(pokertype, t.handCardsType[0], arrayType, t.allArrayType[0], poker1, front, t.allCardArray[0][0], middle, t.allCardArray[0][1], tail, t.allCardArray[0][2])
	pokerInfo2 := container.NewHBox(pokertype, t.handCardsType[1], arrayType, t.allArrayType[1], poker2, front, t.allCardArray[1][0], middle, t.allCardArray[1][1], tail, t.allCardArray[1][2])
	pokerInfo3 := container.NewHBox(pokertype, t.handCardsType[2], arrayType, t.allArrayType[2], poker3, front, t.allCardArray[2][0], middle, t.allCardArray[2][1], tail, t.allCardArray[2][2])
	pokerInfo4 := container.NewHBox(pokertype, t.handCardsType[3], arrayType, t.allArrayType[3], poker4, front, t.allCardArray[3][0], middle, t.allCardArray[3][1], tail, t.allCardArray[3][2])
	pokerInfo5 := container.NewHBox(pokertype, t.handCardsType[4], arrayType, t.allArrayType[4], poker5, front, t.allCardArray[4][0], middle, t.allCardArray[4][1], tail, t.allCardArray[4][2])

	section := container.NewVBox(pokerInfo1, pokerInfo2, pokerInfo3, pokerInfo4, pokerInfo5, t.playButtom)
	c.Add(section)
}

func (t *ChinesepokerClient) resetValue() {
	for i := 0; i < 5; i++ {
		t.allArrayType[i].SetText("")
		t.handCardsType[i].SetText("")
	}

	for i := 0; i < 5; i++ {
		for arrayIndex := range t.allCardArray[i] {
			t.allCardArray[i][arrayIndex].SetText("")
		}
	}

}

func (t *ChinesepokerClient) getHandCard(response *utils.RespBase) {
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
	var handCardType []int
	for i := 0; i < len(pokerInfo); i++ {
		var eachHandCardsArray [][]int
		eachArrayInfo := pokerInfo[i].(map[string]interface{})
		handCardType = append(handCardType, int(eachArrayInfo["handCardType"].(float64)))

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
	t.handCardType = handCardType
	t.cardsArray = handCardsArray
	t.cardsArrayType = transArrayType
}

func (t *ChinesepokerClient) putValue() {
	for i := 0; i < len(t.cardsArrayType); i++ {
		for arrayIndex, array := range t.cardsArray[i] {
			var strArray []string
			for _, card := range array {
				strArray = append(strArray, PrintCard(card))
				combinedstrCardArray := strings.Join(strArray, " ")
				t.allCardArray[i][arrayIndex].SetText(combinedstrCardArray)
			}
		}

		t.handCardsType[i].SetText(strconv.Itoa(t.handCardType[i]))

		if t.handCardType[i] != 0 {
			t.allArrayType[i].SetText("~我是特殊牌型~")
		} else {
			var strArrayType []string
			for _, arrayType := range t.cardsArrayType[i] {
				strType := printArrayType(arrayType)
				strArrayType = append(strArrayType, strType)
			}
			combinedstrArrayType := strings.Join(strArrayType, " ")
			t.allArrayType[i].SetText(combinedstrArrayType)
		}
	}
}

func (t *ChinesepokerClient) SendPlayOperate(action int, index int) {
	type CardInfo struct {
		HandCardType int     `json:"handCardType"`
		CardsArray   [][]int `json:"cardsArray"`
		ArrayType    []int   `json:"arrayType"`
	}

	var data struct {
		PlayOperate struct {
			Instruction int      `json:"instruction"`
			Data        CardInfo `json:"data"`
		}
	}

	info := CardInfo{
		HandCardType: t.handCardType[index],
		CardsArray:   t.cardsArray[index],
		ArrayType:    t.cardsArrayType[index],
	}

	data.PlayOperate.Instruction = action
	data.PlayOperate.Data = info

	t.SendMessage(data)
}

func (t *ChinesepokerClient) SetButton(button *widget.Button, enable bool) {
	if enable {
		button.Enable()
	} else {
		button.Disable()
	}
}

var CardSuit = [4]string{"方塊", "梅花", "紅桃", "黑桃"}
var CardPoint = [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

func PrintCard(card int) (cardStr string) {

	suit := card / 13
	point := card % 13
	cardStr = CardSuit[suit] + CardPoint[point]
	return
}

func printArrayType(typeNum int) string {
	switch typeNum {
	case WU_LONG:
		return "烏龍"
	case DUI_ZI:
		return "對子"
	case LIANG_DUI:
		return "兩對"
	case SAN_TIAO:
		return "三條"
	case SHUN_ZI:
		return "順子"
	case TONG_HUA:
		return "同花"
	case HU_LU:
		return "葫蘆"
	case TIE_ZHI:
		return "鐵支"
	case TONG_HUA_SHUN:
		return "同花順"
	default:
		return "乌龙"
	}
}
