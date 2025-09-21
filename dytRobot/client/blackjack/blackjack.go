package blackjack

import (
	"dytRobot/client"
	"dytRobot/utils"
	"encoding/json"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const (
	ACT_TOKEN_PLAYER_INFO = "ActTokenPlayerSeat"
)

type BlackjackClient struct {
	*client.BaseMatchClient
	seatID int

	buttonChip1 *widget.Button
	buttonChip2 *widget.Button
	buttonChip3 *widget.Button
	buttonChip4 *widget.Button
	buttonChip5 *widget.Button
	buttonChip6 *widget.Button

	buttonSeat0 *widget.Button
	buttonSeat1 *widget.Button
	buttonSeat2 *widget.Button
	buttonSeat3 *widget.Button
	buttonSeat4 *widget.Button

	buttonBetFinish *widget.Button

	selectedSeat int
}

func NewClient(setting client.ClientConfig) *BlackjackClient {
	matchClient := client.NewMatchClient(setting)
	t := &BlackjackClient{
		BaseMatchClient: matchClient,
	}

	t.CheckResponse = t.CheckBlackjackResponse
	t.CustomMessage = append(t.CustomMessage, "{\"MatchGameBet\":{\"BetInfo\":\"[{\\\"SeatId\\\":3,\\\"Bet\\\":500}]\"}}")
	t.CustomMessage = append(t.CustomMessage, "{\"MatchGameBet\":{\"BetInfo\":\"[{\\\"SeatId\\\":1,\\\"Bet\\\":100},{\\\"SeatId\\\":2,\\\"Bet\\\":100}]\"}}")
	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[4,9],[[0,8],[0,8]],[[0,2],[0,2]],[[0,8],[0,8]]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)

	return t
}

func (t *BlackjackClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateMatchSection(c)
	t.CreateBlackJackSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *BlackjackClient) CheckBlackjackResponse(response *utils.RespBase) bool {
	if t.CheckMatchResponse(response) {
		return true
	}

	switch response.Ret {
	case client.RET_JOIN_GAME:
		t.GetTableInfo(response)
	case client.ACT_TABLE_STATUS:
		t.GetSeatInfo(response)
	case "ActDealCard":
		t.setAllChipDisable()
		t.buttonBetFinish.Disable()
	case ACT_TOKEN_PLAYER_INFO:
		t.GetTokenInfo(response)
		return true
	}

	return false
}

func (t *BlackjackClient) CreateBlackJackSection(c *fyne.Container) {
	t.buttonSeat0 = widget.NewButton("座位0", func() {
		t.selectSeat(0)
	})
	t.buttonSeat1 = widget.NewButton("座位1", func() {
		t.selectSeat(1)
	})
	t.buttonSeat2 = widget.NewButton("座位2", func() {
		t.selectSeat(2)
	})
	t.buttonSeat3 = widget.NewButton("座位3", func() {
		t.selectSeat(3)
	})
	t.buttonSeat4 = widget.NewButton("座位4", func() {
		t.selectSeat(4)
	})

	section1 := container.NewHBox(t.buttonSeat0, t.buttonSeat1, t.buttonSeat2, t.buttonSeat3, t.buttonSeat4)

	t.buttonChip1 = widget.NewButton(" ", func() {
		t.SendMatchGameBet(t.buttonChip1.Text)
	})
	t.buttonChip2 = widget.NewButton(" ", func() {
		t.SendMatchGameBet(t.buttonChip2.Text)
	})
	t.buttonChip3 = widget.NewButton(" ", func() {
		t.SendMatchGameBet(t.buttonChip3.Text)
	})
	t.buttonChip4 = widget.NewButton(" ", func() {
		t.SendMatchGameBet(t.buttonChip4.Text)
	})
	t.buttonChip5 = widget.NewButton(" ", func() {
		t.SendMatchGameBet(t.buttonChip5.Text)
	})
	t.buttonChip6 = widget.NewButton(" ", func() {
		t.SendMatchGameBet(t.buttonChip6.Text)
	})
	t.buttonBetFinish = widget.NewButton("下注結束", func() {
		t.SendPlayOperate(5)
		t.setAllSeatDisable()
		t.setAllChipDisable()
		t.buttonBetFinish.Disable()
	})
	section2 := container.NewHBox(t.buttonChip1, t.buttonChip2, t.buttonChip3, t.buttonChip4,
		t.buttonChip5, t.buttonChip6, t.buttonBetFinish)

	buttonHit := widget.NewButton("要牌", func() {
		t.SendPlayOperate(1)
	})
	buttonSpilt := widget.NewButton("分牌", func() {
		t.SendPlayOperate(2)
	})
	buttonDouble := widget.NewButton("雙倍 ", func() {
		t.SendPlayOperate(3)
	})
	buttonStop := widget.NewButton("停牌", func() {
		t.SendPlayOperate(4)
	})
	buttonBuy := widget.NewButton("買保險", func() {
		t.SendInstruction(true)
	})
	buttonNoBuy := widget.NewButton("不買保險", func() {
		t.SendInstruction(false)
	})

	section3 := container.NewHBox(buttonHit, buttonSpilt, buttonDouble, buttonStop, buttonBuy, buttonNoBuy)
	section := container.NewVBox(section1, section2, section3)
	c.Add(section)
}

func (t *BlackjackClient) SendPlayOperate(action int) {
	var data struct {
		PlayOperate struct {
			Instruction int `json:"instruction"`
		}
	}
	data.PlayOperate.Instruction = action
	t.SendMessage(data)
}

func (t *BlackjackClient) SendInstruction(isbuy bool) {
	var data struct {
		SetInsurance struct {
			SeatId int  `json:"seatId"`
			Isbuy  bool `json:"isbuy"`
		}
	}
	data.SetInsurance.SeatId = t.selectedSeat
	data.SetInsurance.Isbuy = isbuy
	t.SendMessage(data)
}

func (t *BlackjackClient) GetTokenInfo(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	seatid, ok := data["SeatId"].(float64)
	if !ok {
		return
	}

	t.setSeatDisable(t.selectedSeat)
	t.selectedSeat = int(seatid)
	t.selectSeat(t.selectedSeat)
	t.setSeatEnable(t.selectedSeat)

}

func (t *BlackjackClient) GetTableInfo(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}

	chips, ok := data["Chips"].([]interface{})
	if !ok {
		return
	}

	for i, value := range chips {
		chip := fmt.Sprintf("%d", int(value.(float64)))
		switch i {
		case 0:
			t.buttonChip1.SetText(chip)
		case 1:
			t.buttonChip2.SetText(chip)
		case 2:
			t.buttonChip3.SetText(chip)
		case 3:
			t.buttonChip4.SetText(chip)
		case 4:
			t.buttonChip5.SetText(chip)
		case 5:
			t.buttonChip6.SetText(chip)
		}
	}
	t.setAllChipEnable()
	t.buttonBetFinish.Enable()

	playerInfo, ok := data["PlayerInfo"].([]interface{})
	if !ok {
		return
	}

	for _, value := range playerInfo {
		info, ok := value.(map[string]interface{})
		if !ok {
			continue
		}

		seatid, ok := info["seatId"].(float64)
		if !ok {
			continue
		}

		t.seatID = int(seatid)
	}

	t.setAllSeatEnable()
}

func (t *BlackjackClient) SendMatchGameBet(info string) {
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

func (t *BlackjackClient) selectSeat(seat int) {
	t.selectedSeat = seat

	t.buttonSeat0.SetIcon(nil)
	t.buttonSeat1.SetIcon(nil)
	t.buttonSeat2.SetIcon(nil)
	t.buttonSeat3.SetIcon(nil)
	t.buttonSeat4.SetIcon(nil)

	switch seat {
	case 0:
		t.buttonSeat0.SetIcon(theme.MediaRecordIcon())
	case 1:
		t.buttonSeat1.SetIcon(theme.MediaRecordIcon())
	case 2:
		t.buttonSeat2.SetIcon(theme.MediaRecordIcon())
	case 3:
		t.buttonSeat3.SetIcon(theme.MediaRecordIcon())
	case 4:
		t.buttonSeat4.SetIcon(theme.MediaRecordIcon())
	}
}
func (t *BlackjackClient) GetSeatInfo(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}

	players, ok := data["tablePlayerInfo"].([]interface{})
	if !ok {
		return
	}

	for _, value := range players {
		player, ok := value.(map[string]interface{})
		if !ok {
			return
		}

		seatid, ok := player["seatId"].(float64)
		if !ok {
			return
		}

		if t.seatID == int(seatid) {
			continue
		}
		t.setSeatDisable(int(seatid))
	}
}

func (t *BlackjackClient) setSeatDisable(seat int) {
	switch seat {
	case 0:
		t.buttonSeat0.Disable()
	case 1:
		t.buttonSeat1.Disable()
	case 2:
		t.buttonSeat2.Disable()
	case 3:
		t.buttonSeat3.Disable()
	case 4:
		t.buttonSeat4.Disable()
	}
}

func (t *BlackjackClient) setSeatEnable(seat int) {
	switch seat {
	case 0:
		t.buttonSeat0.Enable()
	case 1:
		t.buttonSeat1.Enable()
	case 2:
		t.buttonSeat2.Enable()
	case 3:
		t.buttonSeat3.Enable()
	case 4:
		t.buttonSeat4.Enable()
	}
}

func (t *BlackjackClient) setAllSeatEnable() {
	t.buttonSeat0.Enable()
	t.buttonSeat1.Enable()
	t.buttonSeat2.Enable()
	t.buttonSeat3.Enable()
	t.buttonSeat4.Enable()
}

func (t *BlackjackClient) setAllSeatDisable() {
	t.buttonSeat0.Disable()
	t.buttonSeat1.Disable()
	t.buttonSeat2.Disable()
	t.buttonSeat3.Disable()
	t.buttonSeat4.Disable()
}

func (t *BlackjackClient) setAllChipEnable() {
	t.buttonChip1.Enable()
	t.buttonChip2.Enable()
	t.buttonChip3.Enable()
	t.buttonChip4.Enable()
	t.buttonChip5.Enable()
	t.buttonChip6.Enable()
}

func (t *BlackjackClient) setAllChipDisable() {
	t.buttonSeat0.Disable()
	t.buttonSeat1.Disable()
	t.buttonSeat2.Disable()
	t.buttonSeat3.Disable()
	t.buttonSeat4.Disable()
}
