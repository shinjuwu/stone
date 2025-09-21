package goldenflower

import (
	"dytRobot/client"
	"dytRobot/utils"
	"encoding/json"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type GoldenflowerClient struct {
	*client.BaseMatchClient
	// seatID int

	// buttonChip1 *widget.Button
	// buttonChip2 *widget.Button
	// buttonChip3 *widget.Button
	// buttonChip4 *widget.Button
	// buttonChip5 *widget.Button
	// buttonChip6 *widget.Button

	// buttonSeat0 *widget.Button
	// buttonSeat1 *widget.Button
	// buttonSeat2 *widget.Button
	// buttonSeat3 *widget.Button
	// buttonSeat4 *widget.Button

	//buttonBetFinish *widget.Button

	selectedSeat int
}

func NewClient(setting client.ClientConfig) *GoldenflowerClient {
	matchClient := client.NewMatchClient(setting)
	t := &GoldenflowerClient{
		BaseMatchClient: matchClient,
	}

	t.CheckResponse = t.CheckGfResponse
	t.CustomMessage = append(t.CustomMessage, "{\"MatchGameBet\":{\"BetInfo\":\"[{\\\"SeatId\\\":3,\\\"Bet\\\":500}]\"}}")
	t.CustomMessage = append(t.CustomMessage, "{\"MatchGameBet\":{\"BetInfo\":\"[{\\\"SeatId\\\":1,\\\"Bet\\\":100},{\\\"SeatId\\\":2,\\\"Bet\\\":100}]\"}}")
	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[4,9],[[0,8],[0,8]],[[0,2],[0,2]],[[0,8],[0,8]]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}

func (t *GoldenflowerClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateMatchSection(c)
	t.CreateGfControlSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *GoldenflowerClient) CheckGfResponse(response *utils.RespBase) bool {
	return t.CheckMatchResponse(response)
}

func (t *GoldenflowerClient) CreateGfControlSection(c *fyne.Container) {
	// t.buttonSeat0 = widget.NewButton("座位0", func() {
	// 	t.selectSeat(0)
	// })
	// t.buttonSeat1 = widget.NewButton("座位1", func() {
	// 	t.selectSeat(1)
	// })
	// t.buttonSeat2 = widget.NewButton("座位2", func() {
	// 	t.selectSeat(2)
	// })
	// t.buttonSeat3 = widget.NewButton("座位3", func() {
	// 	t.selectSeat(3)
	// })
	// t.buttonSeat4 = widget.NewButton("座位4", func() {
	// 	t.selectSeat(4)
	// })

	// section1 := container.NewHBox(t.buttonSeat0, t.buttonSeat1, t.buttonSeat2, t.buttonSeat3, t.buttonSeat4)

	// t.buttonChip1 = widget.NewButton(" ", func() {
	// 	t.SendMatchGameBet(t.buttonChip1.Text)
	// })
	// t.buttonChip2 = widget.NewButton(" ", func() {
	// 	t.SendMatchGameBet(t.buttonChip2.Text)
	// })
	// t.buttonChip3 = widget.NewButton(" ", func() {
	// 	t.SendMatchGameBet(t.buttonChip3.Text)
	// })
	// t.buttonChip4 = widget.NewButton(" ", func() {
	// 	t.SendMatchGameBet(t.buttonChip4.Text)
	// })
	// t.buttonChip5 = widget.NewButton(" ", func() {
	// 	t.SendMatchGameBet(t.buttonChip5.Text)
	// })
	// t.buttonChip6 = widget.NewButton(" ", func() {
	// 	t.SendMatchGameBet(t.buttonChip6.Text)
	// })
	// t.buttonBetFinish = widget.NewButton("下注結束", func() {
	// 	t.SendPlayOperate(5)
	// 	t.setAllSeatDisable()
	// 	t.setAllChipDisable()
	// 	t.buttonBetFinish.Disable()
	// })
	// section2 := container.NewHBox(t.buttonChip1, t.buttonChip2, t.buttonChip3, t.buttonChip4,
	// 	t.buttonChip5, t.buttonChip6, t.buttonBetFinish)

	buttonFold := widget.NewButton("棄牌", func() {
		t.SendPlayOperate(7)
	})
	buttonCall := widget.NewButton("跟注", func() {
		t.SendPlayOperate(8)
	})
	buttonRaise := widget.NewButton("加注", func() {
		t.SendPlayOperateWithInfo(9, 100)
	})
	buttonStop := widget.NewButton("比牌", func() {
		t.SendPlayOperateWithInfo(10, 0)
	})
	buttonAutoplay := widget.NewButton("防超時棄牌", func() {
		t.SendPlayOperate(11)
	})

	section3 := container.NewHBox(buttonFold, buttonCall, buttonRaise, buttonStop, buttonAutoplay)
	section := container.NewVBox(section3)
	c.Add(section)
}

// func (t *GoldenflowerClient) selectSeat(seat int) {
// 	t.selectedSeat = seat

// 	t.buttonSeat0.SetIcon(nil)
// 	t.buttonSeat1.SetIcon(nil)
// 	t.buttonSeat2.SetIcon(nil)
// 	t.buttonSeat3.SetIcon(nil)
// 	t.buttonSeat4.SetIcon(nil)

// 	switch seat {
// 	case 0:
// 		t.buttonSeat0.SetIcon(theme.MediaRecordIcon())
// 	case 1:
// 		t.buttonSeat1.SetIcon(theme.MediaRecordIcon())
// 	case 2:
// 		t.buttonSeat2.SetIcon(theme.MediaRecordIcon())
// 	case 3:
// 		t.buttonSeat3.SetIcon(theme.MediaRecordIcon())
// 	case 4:
// 		t.buttonSeat4.SetIcon(theme.MediaRecordIcon())
// 	}
// }

func (t *GoldenflowerClient) SendMatchGameBet(info string) {
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

func (t *GoldenflowerClient) SendPlayOperate(action int) {
	var data struct {
		PlayOperate struct {
			Instruction int         `json:"instruction"`
			Data        interface{} `json:"data"`
		}
	}
	data.PlayOperate.Instruction = action
	t.SendMessage(data)
}

func (t *GoldenflowerClient) SendPlayOperateWithInfo(action int, info interface{}) {
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

// func (t *GoldenflowerClient) setAllSeatDisable() {
// 	t.buttonSeat0.Disable()
// 	t.buttonSeat1.Disable()
// 	t.buttonSeat2.Disable()
// 	t.buttonSeat3.Disable()
// 	t.buttonSeat4.Disable()
// }

//	func (t *GoldenflowerClient) setAllChipDisable() {
//		t.buttonSeat0.Disable()
//		t.buttonSeat1.Disable()
//		t.buttonSeat2.Disable()
//		t.buttonSeat3.Disable()
//		t.buttonSeat4.Disable()
//	}
// func (t *GoldenflowerClient) SendInstruction(isbuy bool) {
// 	var data struct {
// 		SetInsurance struct {
// 			SeatId int  `json:"seatId"`
// 			Isbuy  bool `json:"isbuy"`
// 		}
// 	}
// 	data.SetInsurance.SeatId = t.selectedSeat
// 	data.SetInsurance.Isbuy = isbuy
// 	t.SendMessage(data)
// }
