package sangong

import (
	"dytRobot/client"
	"dytRobot/utils"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SangongClient struct {
	*client.BaseMatchClient
	SeatId int
	Banker int
	Bet    []string
	MyInfo *widget.Label
	btRob  []*widget.Button
	btBet  []*widget.Button
}

func NewClient(setting client.ClientConfig) *SangongClient {
	matchClient := client.NewMatchClient(setting)
	t := &SangongClient{
		BaseMatchClient: matchClient,
	}

	t.CheckResponse = t.CheckGameResponse

	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[12,25,38],[-1,-1,-1],[-1,-1,-1],[-1,-1,-1],[-1,-1,-1]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}

func (t *SangongClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateMatchSection(c)
	t.CreateGameSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *SangongClient) CreateGameSection(c *fyne.Container) {
	t.btRob = make([]*widget.Button, 2)
	t.btBet = make([]*widget.Button, 5)
	t.btRob[0] = widget.NewButton("搶庄", func() {
		t.SetRobBanker(true)
	})
	t.btRob[1] = widget.NewButton("不搶", func() {
		t.SetRobBanker(false)
	})
	t.btBet[0] = widget.NewButton("Bet", func() {
		t.SendMatchGameBet(0)
	})
	t.SetButton(t.btRob, false)

	t.btBet[1] = widget.NewButton("Bet", func() {
		t.SendMatchGameBet(1)
	})
	t.btBet[2] = widget.NewButton("Bet", func() {
		t.SendMatchGameBet(2)
	})
	t.btBet[3] = widget.NewButton("Bet", func() {
		t.SendMatchGameBet(3)
	})
	t.btBet[4] = widget.NewButton("Bet", func() {
		t.SendMatchGameBet(4)
	})
	t.SetButton(t.btBet, false)

	t.MyInfo = widget.NewLabel("")
	section := container.NewHBox(t.btRob[0], t.btRob[1], t.btBet[0], t.btBet[1], t.btBet[2], t.btBet[3], t.btBet[4], t.MyInfo)
	c.Add(section)
}

const (
	ACT_BET_INFO = "ActBetInfo"
	ACT_BANKER   = "ActBanker"
)

func (t *SangongClient) CheckGameResponse(response *utils.RespBase) bool {
	if t.CheckMatchResponse(response) {
		return true
	}

	if response.Ret == client.ACT_GAME_PERIOD {
		switch t.Fsm {
		case "Rob":
			t.SetButton(t.btRob, true)
		case "Banker":
			t.SetButton(t.btRob, false)
		case "Result":
			t.SetButton(t.btBet, false)
		}
		return true
	}

	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return true
	}

	switch response.Ret {
	case client.RET_JOIN_GAME:
		t.GetJoinInfo(info)
		t.SetButton(t.btRob, false)
		t.SetButton(t.btBet, false)
		return true
	case ACT_BANKER:
		t.GetBanker(info)
		return true
	case ACT_BET_INFO:
		t.GetBetInfo(info)
		return true
	}

	return false
}

func (t *SangongClient) GetJoinInfo(data map[string]interface{}) {
	t.SeatId = int(data["OwnSeat"].(float64))
	if data["BankerId"] != nil {
		t.Banker = int(data["BankerId"].(float64))
	}
	t.UpdateMyInfo()
}

func (t *SangongClient) GetBanker(data map[string]interface{}) {
	if data["SeatId"] != nil {
		t.Banker = int(data["SeatId"].(float64))
		t.UpdateMyInfo()
	}
}

func (t *SangongClient) SendMatchGameBet(index int) {
	if t.Fsm != "Bet" {
		return
	}
	var data struct {
		MatchGameBet struct {
			BetInfo string `json:"BetInfo"`
		}
	}

	if index < len(t.Bet) {
		data.MatchGameBet.BetInfo = t.Bet[index]
		t.SendMessage(data)
		t.SetButton(t.btBet, false)
	}
}

func (t *SangongClient) SetRobBanker(rob bool) {
	if t.Fsm != "Rob" {
		return
	}
	var data struct {
		SetRobBanker struct {
			Rob bool `json:"rob"`
		}
	}
	data.SetRobBanker.Rob = rob
	t.SendMessage(data)
	t.SetButton(t.btRob, false)
}

func (t *SangongClient) GetBetInfo(data map[string]interface{}) {
	info := data["BetInfo"].([]interface{})
	t.Bet = []string{}
	for i, bet := range info {
		v := strconv.Itoa(int(bet.(float64)))
		t.Bet = append(t.Bet, v)
		t.btBet[i].SetText("x" + v)
		t.btBet[i].Enable()
	}
}

func (t *SangongClient) SetButton(buttons []*widget.Button, enable bool) {
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

func (t *SangongClient) UpdateMyInfo() {
	t.MyInfo.SetText(fmt.Sprintf("Seat %d / Banker %d", t.SeatId, t.Banker))
}
