package bullbull

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
	RET_SET_ROB_BANKER_BET = "SetRobBankerBet"
	RET_MATCH_GAME_BET     = "MatchGameBet"
	RET_SHOW_CARD          = "PlayOperate"

	ACT_ROB_INFO        = "ActRobInfo"
	ACT_PLAYER_ROB      = "ActPlayerRob"
	ACT_BANKER          = "ActBanker"
	ACT_BET_INFO        = "ActBetInfo"
	ACT_PLAYER_BET_INFO = "ActPlayerBetInfo"
	ACT_CARD_INFO       = "ActCardInfo"
	ACT_SHOW_CARDS      = "ActShowCards"
	ACT_SETTLE_INFO     = "ActSettleInfo"
)

var CardTypeName = [15]string{
	" 沒牛",
	" 牛一",
	" 牛二",
	" 牛三",
	" 牛四",
	" 牛五",
	" 牛六",
	" 牛七",
	" 牛八",
	" 牛九",
	" 牛牛",
	"四花牛",
	"五花牛",
	"炸彈牛",
	"五小牛"}

type BullBullClient struct {
	*client.BaseMatchClient
	buttonRob1     *widget.Button
	buttonRob2     *widget.Button
	buttonRob3     *widget.Button
	buttonRob4     *widget.Button
	buttonRob5     *widget.Button
	buttonBet1     *widget.Button
	buttonBet2     *widget.Button
	buttonBet3     *widget.Button
	buttonBet4     *widget.Button
	buttonBet5     *widget.Button
	buttonShowCard *widget.Button
}

func NewClient(setting client.ClientConfig) *BullBullClient {
	matchClient := client.NewMatchClient(setting)
	t := &BullBullClient{
		BaseMatchClient: matchClient,
	}

	t.CheckResponse = t.CheckBullBullResponse

	t.CustomMessage = append(t.CustomMessage, "{\"SetRobBankerBet\":{\"Rob\":-1}}")
	t.CustomMessage = append(t.CustomMessage, "{\"MatchGameBet\":{\"BetInfo\":1}}")
	t.CustomMessage = append(t.CustomMessage, "{\"PlayOperate\":{\"instruction\":6}}")
	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[12,25,38,51,50],[-1,-1,-1,-1,-1],[-1,-1,-1,-1,-1],[-1,-1,-1,-1,-1],[-1,-1,-1,-1,-1]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}

func (t *BullBullClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateMatchSection(c)
	t.CreateBullBullSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *BullBullClient) GetJoinInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	seatId := int(data["OwnSeat"].(float64))

	return fmt.Sprintf("玩家座位:%d\n", seatId)
}

func (t *BullBullClient) GetRobInfo(response *utils.RespBase) string {
	var message = "搶庄資訊: "
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return message
	}
	info, ok := data["RobInfo"].([]interface{})
	if !ok {
		return message
	}

	t.buttonRob1.SetText("-1")
	t.buttonRob1.Enable()

	for i, value := range info {
		rob := int(value.(float64))
		message += strconv.Itoa(rob) + " "

		switch i {
		case 1:
			t.buttonRob2.SetText(strconv.Itoa(rob))
			t.buttonRob2.Enable()
		case 2:
			t.buttonRob3.SetText(strconv.Itoa(rob))
			t.buttonRob3.Enable()
		case 3:
			t.buttonRob4.SetText(strconv.Itoa(rob))
			t.buttonRob4.Enable()
		case 4:
			t.buttonRob5.SetText(strconv.Itoa(rob))
			t.buttonRob5.Enable()
		}

	}
	return message + "\n"
}

func (t *BullBullClient) SendRobBankerBet(rob int) (bool, error) {
	var data struct {
		SetRobBankerBet struct {
			Rob int `json:"Rob"`
		}
	}
	data.SetRobBankerBet.Rob = rob
	t.disableRobButton()
	return t.SendMessage(data)
}

func (t *BullBullClient) GetBankerInfo(response *utils.RespBase) string {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}

	t.disableRobButton()
	seatId := int(info["SeatId"].(float64))
	bankerBet := int(info["BankerBet"].(float64))
	return fmt.Sprintf("莊家座位:%d，莊家倍數:%d\n", seatId, bankerBet)
}

func (t *BullBullClient) GetBetInfo(response *utils.RespBase) string {
	var message = "押注資訊: "
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return message
	}
	info, ok := data["BetInfo"].([]interface{})
	if !ok {
		return message
	}

	for i, value := range info {
		bet := int(value.(float64))
		message += strconv.Itoa(bet) + " "

		switch i {
		case 0:
			t.buttonBet1.SetText(strconv.Itoa(bet))
			t.buttonBet1.Enable()
		case 1:
			t.buttonBet2.SetText(strconv.Itoa(bet))
			t.buttonBet2.Enable()
		case 2:
			t.buttonBet3.SetText(strconv.Itoa(bet))
			t.buttonBet3.Enable()
		case 3:
			t.buttonBet4.SetText(strconv.Itoa(bet))
			t.buttonBet4.Enable()
		case 4:
			t.buttonBet5.SetText(strconv.Itoa(bet))
			t.buttonBet5.Enable()
		}
	}
	return message + "\n"
}

func (t *BullBullClient) GetSettleInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	info, ok := data["SettleInfo"].([]interface{})
	if !ok {
		return ""
	}

	t.buttonShowCard.Disable()

	message := "牌局結果:\n"
	for _, detail := range info {
		result := detail.(map[string]interface{})
		seatId := int(result["SeatId"].(float64))
		data := result["Cards"].([]interface{})
		cards := []int{}
		for _, detail := range data {
			card := int(detail.(float64))
			cards = append(cards, card)
		}
		cardType := int(result["Type"].(float64))
		winLose := result["WinLose"].(float64)
		gold := result["Gold"].(float64)
		message += fmt.Sprintf("座位:%d 手牌:%s 牌型：%s 輸贏：%.4f 結算:%.4f\n", seatId, utils.PrintCard(cards), CardTypeName[cardType], winLose, gold)
	}
	return message
}

func (t *BullBullClient) SendMatchGameBet(bet int) (bool, error) {
	var data struct {
		MatchGameBet struct {
			BetInfo string `json:"BetInfo"`
		}
	}
	data.MatchGameBet.BetInfo = strconv.Itoa(bet)

	t.disableBetButton()

	return t.SendMessage(data)
}

func (t *BullBullClient) SendShowCard() (bool, error) {
	var data struct {
		PlayOperate struct {
			Instruction int `json:"instruction"`
		}
	}
	data.PlayOperate.Instruction = 6
	t.buttonShowCard.Disable()
	return t.SendMessage(data)
}

func (t *BullBullClient) CreateBullBullSection(c *fyne.Container) {
	labelRob := widget.NewLabel("搶庄")
	//搶庄倍數按鈕
	t.buttonRob1 = widget.NewButton(" ", func() {
		rob, err := strconv.Atoi(t.buttonRob1.Text)
		if err != nil {
			return
		}
		t.SendRobBankerBet(rob)
	})
	t.buttonRob1.Disable()
	t.buttonRob2 = widget.NewButton(" ", func() {
		rob, err := strconv.Atoi(t.buttonRob2.Text)
		if err != nil {
			return
		}
		t.SendRobBankerBet(rob)
	})
	t.buttonRob2.Disable()
	t.buttonRob3 = widget.NewButton(" ", func() {
		rob, err := strconv.Atoi(t.buttonRob3.Text)
		if err != nil {
			return
		}
		t.SendRobBankerBet(rob)
	})
	t.buttonRob3.Disable()
	t.buttonRob4 = widget.NewButton(" ", func() {
		rob, err := strconv.Atoi(t.buttonRob4.Text)
		if err != nil {
			return
		}
		t.SendRobBankerBet(rob)
	})
	t.buttonRob4.Disable()
	t.buttonRob5 = widget.NewButton(" ", func() {
		rob, err := strconv.Atoi(t.buttonRob5.Text)
		if err != nil {
			return
		}
		t.SendRobBankerBet(rob)
	})
	t.buttonRob5.Disable()

	//下注倍數按鈕
	labelBet := widget.NewLabel("押注")
	t.buttonBet1 = widget.NewButton(" ", func() {
		bet, err := strconv.Atoi(t.buttonBet1.Text)
		if err != nil {
			return
		}
		t.SendMatchGameBet(bet)
	})
	t.buttonBet1.Disable()
	t.buttonBet2 = widget.NewButton(" ", func() {
		bet, err := strconv.Atoi(t.buttonBet2.Text)
		if err != nil {
			return
		}
		t.SendMatchGameBet(bet)
	})
	t.buttonBet2.Disable()
	t.buttonBet3 = widget.NewButton(" ", func() {
		bet, err := strconv.Atoi(t.buttonBet3.Text)
		if err != nil {
			return
		}
		t.SendMatchGameBet(bet)
	})
	t.buttonBet3.Disable()
	t.buttonBet4 = widget.NewButton(" ", func() {
		bet, err := strconv.Atoi(t.buttonBet4.Text)
		if err != nil {
			return
		}
		t.SendMatchGameBet(bet)
	})
	t.buttonBet4.Disable()
	t.buttonBet5 = widget.NewButton(" ", func() {
		bet, err := strconv.Atoi(t.buttonBet5.Text)
		if err != nil {
			return
		}
		t.SendMatchGameBet(bet)
	})
	t.buttonBet5.Disable()

	//攤牌按鈕
	t.buttonShowCard = widget.NewButton("攤牌", func() {
		t.SendShowCard()
	})
	t.buttonShowCard.Disable()

	section := container.NewHBox(labelRob, t.buttonRob1, t.buttonRob2, t.buttonRob3, t.buttonRob4, t.buttonRob5,
		labelBet, t.buttonBet1, t.buttonBet2, t.buttonBet3, t.buttonBet4, t.buttonBet5, t.buttonShowCard)
	c.Add(section)
}

func (t *BullBullClient) GetPlayerRobInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	seatid := int(data["SeatId"].(float64))
	rob := int(data["Rob"].(float64))
	return fmt.Sprintf("座位:%d 搶庄倍數:%d\n", seatid, rob)
}

func (t *BullBullClient) GetPlayerBetInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	seatid := int(data["SeatId"].(float64))
	bet := int(data["Bet"].(float64))

	return fmt.Sprintf("座位:%d 押注倍數:%d\n", seatid, bet)
}

func (t *BullBullClient) GetCardInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	t.disableBetButton()

	info := data["Cards"].([]interface{})
	cards := []int{}
	for _, detail := range info {
		card := int(detail.(float64))
		cards = append(cards, card)
	}

	t.buttonShowCard.Enable()
	return fmt.Sprintf("玩家手牌:%s\n", utils.PrintCard(cards))
}

func (t *BullBullClient) GetShowCardInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	seatid := int(data["SeatId"].(float64))
	return fmt.Sprintf("座位:%d攤牌\n", seatid)
}

func (t *BullBullClient) CheckBullBullResponse(response *utils.RespBase) bool {
	if t.CheckMatchResponse(response) {
		return true
	}

	switch response.Ret {
	case client.RET_JOIN_GAME:
		t.AddTableStatus(t.GetJoinInfo(response))
		return true
	case ACT_BANKER:
		t.AddTableStatus(t.GetBankerInfo(response))
		return true
	case ACT_ROB_INFO:
		t.AddTableStatus(t.GetRobInfo(response))
		return true
	case ACT_BET_INFO:
		t.AddTableStatus(t.GetBetInfo(response))
		return true
	case ACT_SETTLE_INFO:
		t.AddTableStatus(t.GetSettleInfo(response))
		return true
	case ACT_PLAYER_ROB:
		t.AddTableStatus(t.GetPlayerRobInfo(response))
		return true
	case ACT_PLAYER_BET_INFO:
		t.AddTableStatus(t.GetPlayerBetInfo(response))
		return true
	case ACT_CARD_INFO:
		t.AddTableStatus(t.GetCardInfo(response))
		return true
	case ACT_SHOW_CARDS:
		t.AddTableStatus(t.GetShowCardInfo(response))
		return true
	}

	return false
}

func (t *BullBullClient) disableRobButton() {
	t.buttonRob1.Disable()
	t.buttonRob2.Disable()
	t.buttonRob3.Disable()
	t.buttonRob4.Disable()
	t.buttonRob5.Disable()
}

func (t *BullBullClient) disableBetButton() {
	t.buttonBet1.Disable()
	t.buttonBet2.Disable()
	t.buttonBet3.Disable()
	t.buttonBet4.Disable()
	t.buttonBet5.Disable()
}
