package fruit777slot

import (
	"dytRobot/client"
	"dytRobot/utils"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	BUTTON_COUNT = 8 // 下注區域數量

	ACTION_BET  = 1
	ACTION_TAKE = 2
)

var (
	SymbolName = []string{
		"藍 7",
		"紅 7",
		"3Bar",
		"2Bar",
		"1Bar",
		"鈴鐺",
		"山竹",
		"榴槤",
		"芒果",
		"紅毛丹",
		"混 7",
		"混 Bar",
		"混 水果",
	}

	Lines = []string{
		"1",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
	}

	Bets = []string{
		"0.1",
		"0.2",
		"0.5",
		"0.8",
		"1",
		"2",
		"5",
		"8",
		"10",
		"12",
	}

	S7Name = []string{
		"藍 7",
		"紅 7",
		"混 7",
	}
)

type Fruit777SlotClient struct {
	*client.BaseSlotClient
	Reels         []*widget.Label
	LabelTotalBet *widget.Label
	Line          int
	Bet           float64
}

func NewClient(setting client.ClientConfig) *Fruit777SlotClient {
	slotClient := client.NewSlotClient(setting)
	t := &Fruit777SlotClient{
		BaseSlotClient: slotClient,
	}

	t.CheckResponse = t.CheckFruit777SlotResponse

	t.CustomMessage = append(t.CustomMessage, "{\"SlotAction\":{\"Action\":1,\"Data\":{\"Line\":8,\"Bet\":1}}}")
	t.CustomMessage = append(t.CustomMessage, "{\"SlotAction\":{\"Action\":2}}")
	t.CustomMessage = append(t.CustomMessage, "{\"DebugInfo\":{\"Data\":{\"price\":[6,2,0,23]}}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}

func (t *Fruit777SlotClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateSlotSection(c)
	t.CreateFruit777SlotSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *Fruit777SlotClient) CheckFruit777SlotResponse(response *utils.RespBase) bool {
	if t.CheckSlotResponse(response) {
		return true
	}
	switch response.Ret {
	case client.RET_GOIN_GAME:
		t.AddTableStatus(t.GetJoinInfo(response))
		return true
	case client.ACT_GOLD:
		// t.AddTableStatus(t.GetGoldInfo(response))
		return true
	case client.RET_SLOT_ACTION:
		t.AddTableStatus(t.GetActionInfo(response))
		return true
	}
	return false
}

func (t *Fruit777SlotClient) CreateFruit777SlotSection(c *fyne.Container) {
	//水果機按鈕
	t.Reels = make([]*widget.Label, 3)
	t.Reels[0] = widget.NewLabel("[?]\n[?]\n[?]")
	t.Reels[1] = widget.NewLabel("[?]\n[?]\n[?]")
	t.Reels[2] = widget.NewLabel("[?]\n[?]\n[?]")

	labelLines := widget.NewLabel("Line")
	comboLines := widget.NewSelect(Lines, func(value string) { t.SetLine(value) })
	labelBets := widget.NewLabel("Bet")
	comboBets := widget.NewSelect(Bets, func(value string) { t.SetBet(value) })
	t.LabelTotalBet = widget.NewLabel("TotalBet: 0")

	buttonBet := widget.NewButton("Spin", func() {
		t.SendSlotAction(ACTION_BET)
	})

	section1 := container.NewHBox(t.Reels[0], t.Reels[1], t.Reels[2])
	c.Add(section1)

	section2 := container.NewHBox(labelLines, comboLines, labelBets, comboBets, t.LabelTotalBet, buttonBet)
	c.Add(section2)
}

func (t *Fruit777SlotClient) SendSlotAction(action int) (bool, error) {
	var data struct {
		SlotAction struct {
			Action int         `json:"action"`
			Data   interface{} `json:"data,omitempty"`
		}
	}

	data.SlotAction.Action = action
	if action == ACTION_BET {
		var info struct {
			Line int     `json:"Line"`
			Bet  float64 `json:"Bet"`
		}
		info.Line = t.Line
		info.Bet = t.Bet

		data.SlotAction.Data = info
	}

	return t.SendMessage(data)
}

func (t *Fruit777SlotClient) SendDebug(priceType int, spPriceType int, spPrice1 int, spPrice2 int, spPrice3 int, spPrice4 int) {
	var data struct {
		DebugInfo struct {
			Data struct {
				Price []int `json:"Price"`
			}
		}
	}
	if priceType != -1 {
		data.DebugInfo.Data.Price = append(data.DebugInfo.Data.Price, priceType)
	}
	if spPriceType != -1 {
		data.DebugInfo.Data.Price = append(data.DebugInfo.Data.Price, spPriceType)
	}
	if spPrice1 != -1 {
		data.DebugInfo.Data.Price = append(data.DebugInfo.Data.Price, spPrice1)
	}
	if spPrice2 != -1 {
		data.DebugInfo.Data.Price = append(data.DebugInfo.Data.Price, spPrice2)
	}
	if spPrice3 != -1 {
		data.DebugInfo.Data.Price = append(data.DebugInfo.Data.Price, spPrice3)
	}
	if spPrice4 != -1 {
		data.DebugInfo.Data.Price = append(data.DebugInfo.Data.Price, spPrice4)
	}

	t.SendMessage(data)
}

func (t *Fruit777SlotClient) GetJoinInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	fsm := data["Fsm"].(string)
	labelMessage := "Fsm:" + fsm

	if gid, ok := data["Gid"].(string); ok {
		labelMessage += " Gid:" + gid
	}
	t.LabelFsm.SetText(labelMessage)

	// baseBet := int(data["BaseBet"].(float64))
	// t.labelBaseBet.SetText("底注:" + strconv.Itoa(baseBet))

	var message string
	gold := data["Gold"].(float64)
	message += fmt.Sprintf("玩家金額:%.4f\n", gold)

	return message
}

func (t *Fruit777SlotClient) GetGoldInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}

	gold := data["Gold"].(float64)
	return fmt.Sprintf("玩家金額:%.4f\n", gold)
}

func (t *Fruit777SlotClient) GetActionInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}

	/* fsm := data["Fsm"].(string)
	labelMessage := "Fsm:" + fsm */
	gid := data["Gid"].(string)
	labelMessage := " Gid:" + gid
	t.LabelFsm.SetText(labelMessage)

	action := int(data["Action"].(float64))
	gold := data["Gold"].(float64)
	message := fmt.Sprintf("[%d] Gold: %.1f", action, gold)

	switch action {
	case ACTION_BET:
		t.SetTableStatusClear()

		if totems, ok := data["Totems"].([]interface{}); ok {
			reelStr := make([]string, 3)
			for i, totem := range totems {
				v := int(totem.(float64))
				l := i % 3
				symb := "[ " + SymbolName[v] + " ]"
				if (i / 3) == 0 {
					reelStr[l] = symb
				} else {
					reelStr[l] += "\n" + symb
				}
			}
			for i := 0; i < len(t.Reels); i++ {
				t.Reels[i].SetText(reelStr[i])
			}
		}

		win := int(data["Win"].(float64))
		message += fmt.Sprintf(",  Line: %d,  Bet: %.1f,  Win:%d", t.Line, t.Bet, win)
		if win > 0 {
			if jp, ok := data["JP"].(interface{}); ok {
				message += fmt.Sprintf(",  JP: %s", SymbolName[int(jp.(float64))])
			} else {
				if lines, ok := data["Lines"].(interface{}); ok {
					message += fmt.Sprintf(",  Lines: %+v", lines)
				}
				if s7, ok := data["Special7"].(interface{}); ok {
					message += fmt.Sprintf(",  S7: %s", S7Name[int(s7.(float64))])
				}
				if c7, ok := data["Count7"].(interface{}); ok {
					message += fmt.Sprintf("x%d", int(c7.(float64)))
				}
			}
		}
		t.SendSlotAction(ACTION_TAKE)
	}

	message += "\n"

	return message
}

func (t *Fruit777SlotClient) SetLine(value string) {
	if l, err := strconv.Atoi(value); err == nil {
		t.Line = l
		t.UpdateTotalBet()
	}
}

func (t *Fruit777SlotClient) SetBet(value string) {
	if b, err := strconv.ParseFloat(value, 64); err == nil {
		t.Bet = b
		t.UpdateTotalBet()
	}
}

func (t *Fruit777SlotClient) UpdateTotalBet() {
	totalBet := float64(t.Line) * t.Bet
	message := fmt.Sprintf("TotalBet: %.1f", totalBet)
	t.LabelTotalBet.SetText(message)

}
