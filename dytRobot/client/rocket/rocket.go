package rocket

import (
	"dytRobot/client"
	"dytRobot/utils"
	"encoding/json"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type RocketClient struct {
	*client.BaseBetClient
	buttonBet   *widget.Button
	buttonFlee  *widget.Button
	entryDebug  *widget.Entry
	buttonDebug *widget.Button
	chipsLabel  *widget.Label
	chipsSelect *widget.Select
	fleeLabel   *widget.Label
	fleeSelect  *widget.Select
}

func NewClient(setting client.ClientConfig) *RocketClient {
	betClient := client.NewBetClient(setting)
	t := &RocketClient{
		BaseBetClient: betClient,
	}
	t.CheckResponse = t.CheckRocketResponse
	// t.CustomMessage = append(t.CustomMessage, "{\"DebugPayout\":%d")
	t.EntrySendMessage.SetOptions(t.CustomMessage)

	return t
}

func (t *RocketClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateBetSection(c)
	t.CreateRocketSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *RocketClient) CheckRocketResponse(response *utils.RespBase) bool {
	if t.CheckBaseResponse(response) {
		return true
	}

	switch response.Ret {
	case client.ACT_GAME_PERIOD:
		t.SetButton(t.buttonBet, false)
		t.SetButton(t.buttonFlee, false)
		switch t.Fsm {
		case "Bet":
			t.SetButton(t.buttonBet, true)
		case "Open":
			t.SetButton(t.buttonFlee, true)
		}
	}

	return t.CheckBetResponse(response)
}

func (t *RocketClient) SetButton(button *widget.Button, enable bool) {
	if enable {
		button.Enable()
	} else {
		button.Disable()
	}
}

func (t *RocketClient) CreateRocketSection(c *fyne.Container) {
	t.chipsLabel = widget.NewLabel("下注額度")
	t.chipsSelect = widget.NewSelect([]string{"5", "10", "15", "20", "25"}, nil)
	t.fleeLabel = widget.NewLabel("自動逃脫")
	t.fleeSelect = widget.NewSelect([]string{"5", "10", "15", "20", "25", "100"}, nil)
	t.buttonBet = widget.NewButton("下注", func() {
		betStr := t.chipsSelect.Selected
		bet, _ := strconv.ParseFloat(betStr, 64)
		fleeStr := t.fleeSelect.Selected
		flee, _ := strconv.ParseFloat(fleeStr, 64)

		t.SendGameBet(bet, flee)
	})
	t.buttonFlee = widget.NewButton("逃離", func() {
		t.SendGameFlee()
	})
	t.entryDebug = widget.NewEntry()
	t.entryDebug.SetText("0")
	t.buttonDebug = widget.NewButton("Debug", func() {
		t.SendDebugInfo()
	})

	section1 := container.NewHBox(t.chipsLabel, t.chipsSelect, t.fleeLabel, t.fleeSelect)
	section2 := container.NewHBox(t.buttonBet, t.buttonFlee)
	section3 := container.NewHBox(t.entryDebug, t.buttonDebug)
	section := container.NewVBox(section1, section2, section3)
	c.Add(section)
}

func (t *RocketClient) SendGameFlee() {
	var data struct {
		BetOp struct {
			Instruction int `json:"instruction"`
		}
	}
	data.BetOp.Instruction = 1
	t.SendMessage(data)
}

func (t *RocketClient) SendDebugInfo() {
	payout, _ := strconv.Atoi(t.entryDebug.Text)
	fPayout := float64(payout)
	var data struct {
		DebugInfo struct {
			Data struct {
				Payout float64 `json:"Payout"`
			}
		}
	}
	data.DebugInfo.Data.Payout = fPayout
	t.SendMessage(data)
}

func (t *RocketClient) SendGameBet(bet float64, flee float64) {
	var data struct {
		Bet struct {
			BetInfo string `json:"BetInfo"`
		}
	}

	type BetInfo struct {
		Bet            float64 `json:"Bet"`
		AutoFleePayout float64 `json:"AutoFleePayout"`
	}

	var betInfo BetInfo

	betInfo.AutoFleePayout = flee
	betInfo.Bet = bet

	m, _ := json.Marshal(betInfo)

	data.Bet.BetInfo = string(m)

	t.SendMessage(data)
}
