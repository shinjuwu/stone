package catte

import (
	"dytRobot/client"
	"dytRobot/utils"
	"sort"
	"strconv"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type CatteClient struct {
	*client.BaseMatchClient
	ownHandcards     []int
	ownSeat          int
	CardIndex        int
	currentSeat      int
	Round            int
	round5BankerCard *widget.Label
	trust            *widget.Button
	playButtom       []*widget.Button
	playCard         []*widget.Button
	roundMaxCard     *widget.Label
	round            *widget.Label
}

var once sync.Once

func NewClient(setting client.ClientConfig) *CatteClient {
	matchClient := client.NewMatchClient(setting)
	t := &CatteClient{
		BaseMatchClient: matchClient,
	}
	t.CheckResponse = t.CheckCatteResponse
	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[0,13,26,2,12,25]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)

	return t
}

func (t *CatteClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateMatchSection(c)
	t.CreateGameSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *CatteClient) CheckCatteResponse(response *utils.RespBase) bool {
	if t.CheckMatchResponse(response) {
		return true
	}
	switch response.Ret {
	case "ActDealCard":
		t.getHandCards(response)
		t.putCardsNum()
	case client.ACT_GAME_PERIOD:
		t.SetButton(t.playButtom, false)
		if t.Fsm == "Play" || t.Fsm == "PlayRoundSix" {
			t.getRound(response)
		}
	case "ActTokenPlayerSeat":
		t.getCurrentSeat(response)
		if t.Fsm == "Play" && t.currentSeat == t.ownSeat {
			t.SetButton(t.playButtom, true)
		}
	case "ActAction":
		t.getRound5BankerCard(response)
		t.getMaxCard(response)
		if t.currentSeat == t.ownSeat {
			t.deletePlayCard(response)
		}
	case "ActSettleData":
		if t.Fsm == "Result" {
			t.resetValue()
		}
	}

	return false
}

func (t *CatteClient) CreateGameSection(c *fyne.Container) {
	t.playButtom = make([]*widget.Button, 2)
	t.playButtom[0] = widget.NewButton("棄牌", func() { t.SendPlayOperate(14, t.CardIndex) })
	t.playButtom[1] = widget.NewButton("比牌", func() { t.SendPlayOperate(10, t.CardIndex) })
	t.SetButton(t.playButtom, false)

	t.playCard = make([]*widget.Button, 6)
	t.playCard[0] = widget.NewButton("", func() { t.choseCard(0) })
	t.playCard[1] = widget.NewButton("", func() { t.choseCard(1) })
	t.playCard[2] = widget.NewButton("", func() { t.choseCard(2) })
	t.playCard[3] = widget.NewButton("", func() { t.choseCard(3) })
	t.playCard[4] = widget.NewButton("", func() { t.choseCard(4) })
	t.playCard[5] = widget.NewButton("", func() { t.choseCard(5) })

	t.trust = widget.NewButton("託管", func() { t.SendPlayOperate(17, -1) })
	maxCard := widget.NewLabel("目前最大牌:")
	t.roundMaxCard = widget.NewLabel("")
	roundLabel := widget.NewLabel("Round:")
	t.round = widget.NewLabel("")
	round5Label := widget.NewLabel("Round5FirstCard:")
	t.round5BankerCard = widget.NewLabel("")

	selfInfo := container.NewHBox(t.playCard[0], t.playCard[1], t.playCard[2], t.playCard[3], t.playCard[4], t.playCard[5])
	deskInfo := container.NewHBox(t.playButtom[0], t.playButtom[1])
	doTrust := container.NewHBox(t.trust)
	tableInfo := container.NewHBox(roundLabel, t.round, maxCard, t.roundMaxCard, round5Label, t.round5BankerCard)
	section := container.NewVBox(selfInfo, deskInfo, doTrust, tableInfo)
	c.Add(section)
}

func (t *CatteClient) getRound5BankerCard(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	if t.Round != 5 {
		return
	}
	once.Do(func() {
		round5Card := PrintCard(int(data["card"].(float64)))
		t.round5BankerCard.SetText(round5Card)
	})

}

func (t *CatteClient) resetValue() {
	for i := 0; i < 6; i++ {
		t.playCard[i].SetText("")
	}
	t.SetButton(t.playCard, true)
	t.round.SetText("")
	t.roundMaxCard.SetText("")
	t.round5BankerCard.SetText("")
}
func (t *CatteClient) getRound(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	t.Round = int(data["Round"].(float64))
	round := strconv.FormatFloat(data["Round"].(float64), 'f', -1, 64)
	t.round.SetText(round)
}

func (t *CatteClient) getMaxCard(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	actiontype := data["actiontype"].(string)
	if actiontype == "compare" {
		roundMaxCard := int(data["card"].(float64))
		cardstr := PrintCard(roundMaxCard)
		t.roundMaxCard.SetText(cardstr)
	}

}

func (t *CatteClient) deletePlayCard(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}

	deleteCard := int(data["card"].(float64))
	cardIndex := sort.SearchInts(t.ownHandcards, deleteCard)
	button := []*widget.Button{t.playCard[cardIndex]}
	t.SetButton(button, false)

}

func (t *CatteClient) getCurrentSeat(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}

	seatId := int(data["SeatId"].(float64))
	t.currentSeat = seatId

}

func (t *CatteClient) getHandCards(response *utils.RespBase) {
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

	ownSeat := int(data["seatid"].(float64))
	t.ownSeat = ownSeat
	t.ownHandcards = handCards
}

func (t *CatteClient) putCardsNum() {
	var cardstr []string
	for _, val := range t.ownHandcards {
		cardstr = append(cardstr, PrintCard(val))
	}

	for i := 0; i < len(cardstr); i++ {
		t.playCard[i].SetText(cardstr[i])
	}

}

func (t *CatteClient) choseCard(cardIndex int) {
	t.CardIndex = cardIndex
}

func (t *CatteClient) SendPlayOperate(action int, index int) {
	var data struct {
		PlayOperate struct {
			Instruction int         `json:"instruction"`
			Data        interface{} `json:"data"`
		}
	}
	data.PlayOperate.Instruction = action
	if index != -1 {
		data.PlayOperate.Data = t.ownHandcards[index]
	}
	t.SendMessage(data)
}

func (t *CatteClient) SetButton(buttons []*widget.Button, enable bool) {
	if enable {
		for _, button := range buttons {
			button.Enable()
		}
	} else {
		for _, button := range buttons {
			button.Disable()
		}
	}
}

var CardSuit = [4]string{"黑桃", "梅花", "方塊", "紅桃"}
var CardPoint = [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

func PrintCard(card int) (cardStr string) {

	suit := card / 13
	point := card % 13
	cardStr = CardSuit[suit] + CardPoint[point]
	return
}
