package pokdeng

import (
	"dytRobot/client"
	"dytRobot/utils"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type PokdengClient struct {
	*client.BaseMatchClient
	SeatId    int
	Banker    int
	SeatInfo  [2]BetCardsInfo
	Chips     []string
	TableInfo *widget.Label
	btDraw    []*widget.Button
	btBet     []*widget.Button
}

type BetCardsInfo struct {
	Bet   int
	Cards []int
	Type  int
	Type2 int
	Odds  int
	Win   float64
}

func NewClient(setting client.ClientConfig) *PokdengClient {
	matchClient := client.NewMatchClient(setting)
	t := &PokdengClient{
		BaseMatchClient: matchClient,
	}

	t.CheckResponse = t.CheckGameResponse

	// t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[12,25,38],[-1,-1,-1],[-1,-1,-1],[-1,-1,-1],[-1,-1,-1]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}

func (t *PokdengClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateMatchSection(c)
	t.CreateGameSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *PokdengClient) CreateGameSection(c *fyne.Container) {
	t.btBet = make([]*widget.Button, 4)
	t.btDraw = make([]*widget.Button, 2)
	t.btBet[0] = widget.NewButton("Bet", func() {
		t.SendMatchGameBet(0)
	})
	t.btBet[1] = widget.NewButton("Bet", func() {
		t.SendMatchGameBet(1)
	})
	t.btBet[2] = widget.NewButton("Bet", func() {
		t.SendMatchGameBet(2)
	})
	t.btBet[3] = widget.NewButton("Bet", func() {
		t.SendMatchGameBet(3)
	})
	t.SetButton(t.btBet, false)

	t.btDraw[0] = widget.NewButton("補牌", func() {
		t.SetDraw(true)
	})
	t.btDraw[1] = widget.NewButton("不補", func() {
		t.SetDraw(false)
	})
	t.SetButton(t.btDraw, false)

	t.TableInfo = widget.NewLabel("")
	section := container.NewHBox(t.btBet[0], t.btBet[1], t.btBet[2], t.btBet[3], t.btDraw[0], t.btDraw[1], t.TableInfo)
	c.Add(section)
}

const (
	ACT_BET_INFO       = "ActBetInfo"
	ACT_CARD_INFO      = "ActCardInfo"
	ACT_BET_CARDS_INFO = "ActBetCardsInfo"
	ACT_DRAW_INFO      = "ActDrawInfo"
	ACT_DRAWCARD_INFO  = "ActDrawCardInfo"
	ACT_SETTLE_INFO    = "ActSettleInfo"
)

func (t *PokdengClient) CheckGameResponse(response *utils.RespBase) bool {
	if t.CheckMatchResponse(response) {
		return true
	}

	if response.Ret == client.ACT_GAME_PERIOD {
		switch t.Fsm {
		case "Bet":
			t.SetButton(t.btBet, true)
		case "Deal":
			t.SetButton(t.btBet, false)
		case "Draw":
			t.SetButton(t.btDraw, true)
		case "BankerDraw":
			t.SetButton(t.btDraw, false)
		case "Settle":
			t.SetButton(t.btDraw, false)
		}
		return true
	}

	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return true
	}

	switch response.Ret {
	case client.RET_JOIN_GAME:
		t.ResetData()
		t.GetJoinInfo(info)
		t.TableInfo.SetText("")
		t.SetButton(t.btBet, false)
		return true
	case ACT_BET_CARDS_INFO:
		t.GetBetCardsInfo(info)
	case ACT_CARD_INFO:
		t.GetCardsInfo(info)
		return true
	case ACT_BET_INFO:
		// t.GetBetInfo(info)
	case ACT_DRAW_INFO:
		t.GetDrawInfo(info)
		return true
	case ACT_DRAWCARD_INFO:
		t.GetDrawCardInfo(info)
		return true
	case ACT_SETTLE_INFO:
		t.GetSettleInfo(info)
		return true
	}

	return false
}

const (
	CARDTYPE_POINT_0        int = iota // 點數 0 - (對子/同花)
	CARDTYPE_POINT_1                   // 點數 1 - (對子/同花)
	CARDTYPE_POINT_2                   // 點數 2 - (對子/同花)
	CARDTYPE_POINT_3                   // 點數 3 - (對子/同花)
	CARDTYPE_POINT_4                   // 點數 4 - (對子/同花)
	CARDTYPE_POINT_5                   // 點數 5 - (對子/同花)
	CARDTYPE_POINT_6                   // 點數 6 - (對子/同花)
	CARDTYPE_POINT_7                   // 點數 7 - (對子/同花)
	CARDTYPE_POINT_8                   // 點數 8 - (對子/同花)
	CARDTYPE_POINT_9                   // 點數 9 - (對子/同花)
	CARDTYPE_SANGONG                   //  三公  -
	CARDTYPE_STRAIGHT                  //  順子  - AKQ...432 (32A 不算順)
	CARDTYPE_STRAIGHT_FLUSH            // 同花順 - AKQ...432 (32A 不算順)
	CARDTYPE_TRIPLE                    //  三條  -
	CARDTYPE_POK_8                     //  博八  - 兩張牌 (對子/同花)
	CARDTYPE_POK_9                     //  博九  - 兩張牌 (同花)

	CARDTYPE_COUNT
)

var (
	CardTypeName = [CARDTYPE_COUNT]string{
		"零點",
		"一點",
		"二點",
		"三點",
		"四點",
		"五點",
		"六點",
		"七點",
		"八點",
		"九點",
		"三公",
		"順子",
		"同花順",
		"三條",
		"博八",
		"博九"}

	CardType2Name = [3]string{
		"",     // 一般
		"(對子)", // 對子
		"(同花)", // 同花
	}
)

func (t *PokdengClient) ResetData() {
	for i, _ := range t.SeatInfo {
		t.SeatInfo[i].Bet = 0
		t.SeatInfo[i].Cards = []int{}
		t.SeatInfo[i].Type = 0
		t.SeatInfo[i].Type2 = 0
		t.SeatInfo[i].Win = 0
		t.SeatInfo[i].Odds = 0
	}
}

func (t *PokdengClient) GetJoinInfo(data map[string]interface{}) {
	t.SeatId = int(data["OwnSeat"].(float64))
	// t.Chips = data["Chips"].(float64)
	t.GetChipsInfo(data)
	if data["BankerId"] != nil {
		t.Banker = int(data["BankerId"].(float64))
	}
	t.UpdateTableInfo()
}

func (t *PokdengClient) SendMatchGameBet(index int) {
	if t.Fsm != "Bet" {
		return
	}
	var data struct {
		MatchGameBet struct {
			BetInfo string `json:"BetInfo"`
		}
	}

	if index < len(t.Chips) {
		data.MatchGameBet.BetInfo = t.Chips[index]
		t.SendMessage(data)
		t.SetButton(t.btBet, false)
	}
}

func (t *PokdengClient) GetChipsInfo(data map[string]interface{}) {
	info := data["Chips"].([]interface{})
	t.Chips = []string{}
	for i, chip := range info {
		v := strconv.Itoa(int(chip.(float64)))
		t.Chips = append(t.Chips, v)
		t.btBet[i].SetText(v)
		t.btBet[i].Enable()
	}
}

func (t *PokdengClient) SetDraw(draw bool) {
	if t.Fsm != "Draw" {
		return
	}
	var data struct {
		PlayOperate struct {
			Instruction int `json:"instruction"`
			Data        int `json:"data"`
		}
	}
	data.PlayOperate.Instruction = 16
	if draw {
		data.PlayOperate.Data = 1
	} else {
		data.PlayOperate.Data = 2
	}

	t.SendMessage(data)
	t.SetButton(t.btDraw, false)
}

func (t *PokdengClient) SetButton(buttons []*widget.Button, enable bool) {
	if enable {
		if t.Fsm == "Bet" {
			for _, button := range buttons {
				button.Enable()
			}
		} else if t.Fsm == "Draw" {
			if t.SeatInfo[0].Type >= CARDTYPE_POK_8 {
				t.btDraw[0].Disable()
				t.btDraw[1].Disable()
			} else {
				t.btDraw[0].Enable()
				if t.SeatInfo[0].Type < CARDTYPE_POINT_4 && t.SeatInfo[0].Type2 == 0 {
					t.btDraw[1].Disable()
				} else {
					t.btDraw[1].Enable()
				}
			}
		}
	} else {
		for _, button := range buttons {
			button.Disable()
		}
	}
}

func (t *PokdengClient) GetBetCardsInfo(data map[string]interface{}) {
	info := int(data["SeatId"].(float64))
	if info == t.SeatId {
		t.SeatInfo[0].Bet = int(data["Bet"].(float64))
		t.UpdateTableInfo()
	}
}

func (t *PokdengClient) GetCardsInfo(data map[string]interface{}) {
	var cards []int
	info := data["Cards"].([]interface{})
	for _, card := range info {
		cards = append(cards, int(card.(float64)))
	}
	type1 := int(data["Type"].(float64))
	type2 := int(data["Type2"].(float64))
	seat := int(data["SeatId"].(float64))
	if seat == t.SeatId || seat == t.Banker {
		var index int
		if seat == t.Banker {
			index = 1
		}
		t.SeatInfo[index].Cards = cards
		t.SeatInfo[index].Type = type1
		t.SeatInfo[index].Type2 = type2
		t.UpdateTableInfo()
	}
}

func (t *PokdengClient) GetDrawInfo(data map[string]interface{}) {
	if data["SeatId"] != nil {
		seatId := int(data["SeatId"].(float64))
		if seatId == t.SeatId {
			t.SetButton(t.btDraw, false)
		}
	}
}

func (t *PokdengClient) GetDrawCardInfo(data map[string]interface{}) {
	update := false
	if data["Card"] != nil {
		card := int(data["Card"].(float64))
		t.SeatInfo[0].Cards = append(t.SeatInfo[0].Cards, card)
		update = true
	}
	if data["Type"] != nil {
		t.SeatInfo[0].Type = int(data["Type"].(float64))
		update = true
	}
	if data["Type2"] != nil {
		t.SeatInfo[0].Type2 = int(data["Type2"].(float64))
		update = true
	}
	if update {
		t.UpdateTableInfo()
	}
}

func (t *PokdengClient) GetSettleInfo(data map[string]interface{}) string {
	info, ok := data["SettleInfo"].([]interface{})
	if !ok {
		return ""
	}

	for _, detail := range info {
		cards := []int{}

		seatInfo := detail.(map[string]interface{})
		seatId := int(seatInfo["SeatId"].(float64))
		if seatId != t.Banker && seatId != t.SeatId {
			continue
		}
		index := 0
		if seatId == t.Banker {
			index = 1
		}

		if seatInfo["Cards"] != nil {
			data := seatInfo["Cards"].([]interface{})
			for _, card := range data {
				cards = append(cards, int(card.(float64)))
			}
			t.SeatInfo[index].Cards = cards
		}

		if seatInfo["Type"] != nil {
			t.SeatInfo[index].Type = int(seatInfo["Type"].(float64))
		}
		if seatInfo["Odds"] != nil {
			t.SeatInfo[index].Odds = int(seatInfo["Odds"].(float64))
		}
		if seatInfo["Win"] != nil {
			t.SeatInfo[index].Win = seatInfo["Win"].(float64)
		}
	}
	t.UpdateTableInfo()

	return ""
}

func (t *PokdengClient) UpdateTableInfo() {
	t.TableInfo.SetText(fmt.Sprintf("Banker %d %s / Seat %d %s", t.Banker, t.UpdateInfo(1), t.SeatId, t.UpdateInfo(0)))
}

func (t *PokdengClient) UpdateInfo(index int) (str string) {
	if index == 0 && t.SeatInfo[0].Bet > 0 {
		str = "Bet " + strconv.Itoa(t.SeatInfo[0].Bet) + " "
	}
	if len(t.SeatInfo[index].Cards) > 0 {
		str += utils.PrintCard(t.SeatInfo[index].Cards)
		str += " [" + CardTypeName[t.SeatInfo[index].Type]
		if t.SeatInfo[index].Type2 > 0 {
			str += " " + CardType2Name[t.SeatInfo[index].Type2]
		}
		str += "]"

		if t.SeatInfo[index].Odds > 1 {
			str += " " + strconv.Itoa(t.SeatInfo[index].Odds) + "x"
		}
		if t.SeatInfo[index].Win > 0 {
			str = fmt.Sprintf("%s Win %f", str, t.SeatInfo[index].Win)
		}
	}
	return
}
