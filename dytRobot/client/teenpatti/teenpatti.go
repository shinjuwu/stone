package teenpatti

import (
	"dytRobot/client"
	"dytRobot/utils"
	"encoding/json"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TeenpattiClient struct {
	*client.BaseMatchClient
	selectedSeat int
}

func NewClient(setting client.ClientConfig) *TeenpattiClient {
	matchClient := client.NewMatchClient(setting)
	t := &TeenpattiClient{
		BaseMatchClient: matchClient,
	}

	t.CheckResponse = t.CheckGfResponse
	// t.CustomMessage = append(t.CustomMessage, "{\"MatchGameBet\":{\"BetInfo\":\"[{\\\"SeatId\\\":3,\\\"Bet\\\":500}]\"}}")
	// t.CustomMessage = append(t.CustomMessage, "{\"MatchGameBet\":{\"BetInfo\":\"[{\\\"SeatId\\\":1,\\\"Bet\\\":100},{\\\"SeatId\\\":2,\\\"Bet\\\":100}]\"}}")
	// t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[4,9],[[0,8],[0,8]],[[0,2],[0,2]],[[0,8],[0,8]]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}

func (t *TeenpattiClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateMatchSection(c)
	t.CreateGfControlSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *TeenpattiClient) CheckGfResponse(response *utils.RespBase) bool {
	return t.CheckMatchResponse(response)
}

func (t *TeenpattiClient) CreateGfControlSection(c *fyne.Container) {

	buttonFold := widget.NewButton("棄牌", func() {
		t.SendPlayOperate(7)
	})
	buttonCall := widget.NewButton("跟注", func() {
		t.SendPlayOperate(8)
	})
	buttonRaise := widget.NewButton("加注", func() {
		t.SendPlayOperate(9)
	})
	buttonCompare := widget.NewButton("比牌", func() {
		t.SendPlayOperateWithInfo(10, 0)
	})
	buttonLookAtCard := widget.NewButton("看牌", func() {
		t.SendPlayOperate(12)
	})
	buttonAccept := widget.NewButton("接受比牌", func() {
		t.SendPlayOperateWithInfo(10, 2)
	})
	buttonReject := widget.NewButton("拒絕比牌", func() {
		t.SendPlayOperateWithInfo(10, 1)
	})

	section3 := container.NewHBox(buttonLookAtCard, buttonFold, buttonCall, buttonRaise, buttonCompare, buttonAccept, buttonReject)
	section := container.NewVBox(section3)
	c.Add(section)
}

func (t *TeenpattiClient) SendMatchGameBet(info string) {
	bet, err := strconv.Atoi(info)
	if err != nil {
		return
	}

	var data struct {
		MatchGameBet struct {
			BetInfo string `json:"BetInfo"`
		}
	}

	type Bet struct {
		SeatId int `json:"SeatId"`
		Bet    int `json:"Bet"`
	}

	var betinfo Bet
	var bets []Bet

	betinfo.Bet = bet
	betinfo.SeatId = t.selectedSeat
	bets = append(bets, betinfo)

	m, _ := json.Marshal(bets)

	data.MatchGameBet.BetInfo = string(m)

	t.SendMessage(data)
}

func (t *TeenpattiClient) SendPlayOperate(action int) {
	var data struct {
		PlayOperate struct {
			Instruction int         `json:"instruction"`
		}
	}
	data.PlayOperate.Instruction = action
	t.SendMessage(data)
}

func (t *TeenpattiClient) SendPlayOperateWithInfo(action int, info interface{}) {
	var data struct {
		PlayOperate struct {
			Instruction int         `json:"instruction"`
			Data        interface{} `json:"data"`
		}
	}
	data.PlayOperate.Instruction = action
	data.PlayOperate.Data = info
	t.SendMessage(data)
}
