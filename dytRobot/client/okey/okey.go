package okey

import (
	"dytRobot/client"
	"dytRobot/utils"
	"fmt"
	"sort"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type OkeyClient struct {
	*client.BaseMatchClient
	Seat       int
	PrevSeat   int
	CurSeat    int
	HoleCards  []int
	Discard    int
	DeckNum    int
	Win        bool
	PickCard   int
	DiscardWin []int

	btDecks    []*widget.Button
	btCards    []*widget.Button
	btWinCards []*widget.Button
	SelfInfo   *widget.Label
	TableInfo  *widget.Label
}

const (
	UNDEFINED     = -1
	DECK_PICKCARD = 0
	DECK_DISCARD  = 1

	PLAYER_HOLECARD_NUM   = 14
	DEALER_START_HOLECARD = PLAYER_HOLECARD_NUM + 1
)

const (
	SUIT_BLUE   int = iota // 藍 B
	SUIT_RED               // 紅 R
	SUIT_BLACK             // 黑 K
	SUIT_YELLOW            // 黃 Y
	SUIT_OKEY              // 鬼牌
	SUIT_COUNT
)

func NewClient(setting client.ClientConfig) *OkeyClient {
	matchClient := client.NewMatchClient(setting)
	t := &OkeyClient{
		BaseMatchClient: matchClient,
	}

	t.CheckResponse = t.CheckGameResponse

	/* t.CustomMessage = append(t.CustomMessage, "{\"MatchGameBet\":{\"BetInfo\":1}}")
	t.CustomMessage = append(t.CustomMessage, "{\"PlayOperate\":{\"instruction\":7}}")
	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[11,10,9,25,38],[12,8],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage) */
	return t
}

func (t *OkeyClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateMatchSection(c)
	t.CreateGameSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *OkeyClient) CreateGameSection(c *fyne.Container) {
	t.btDecks = make([]*widget.Button, 2)
	t.btDecks[DECK_PICKCARD] = widget.NewButton("48", func() { t.SendPickCard(DECK_PICKCARD) })
	t.btDecks[DECK_DISCARD] = widget.NewButton("", func() { t.SendPickCard(DECK_DISCARD) })
	t.SetButton(t.btDecks, false)

	t.btCards = make([]*widget.Button, 15)
	t.btCards[0] = widget.NewButton("\n", func() { t.SendDiscard(0) })
	t.btCards[1] = widget.NewButton("\n", func() { t.SendDiscard(1) })
	t.btCards[2] = widget.NewButton("\n", func() { t.SendDiscard(2) })
	t.btCards[3] = widget.NewButton("\n", func() { t.SendDiscard(3) })
	t.btCards[4] = widget.NewButton("\n", func() { t.SendDiscard(4) })
	t.btCards[5] = widget.NewButton("\n", func() { t.SendDiscard(5) })
	t.btCards[6] = widget.NewButton("\n", func() { t.SendDiscard(6) })
	t.btCards[7] = widget.NewButton("\n", func() { t.SendDiscard(7) })
	t.btCards[8] = widget.NewButton("\n", func() { t.SendDiscard(8) })
	t.btCards[9] = widget.NewButton("\n", func() { t.SendDiscard(9) })
	t.btCards[10] = widget.NewButton("\n", func() { t.SendDiscard(10) })
	t.btCards[11] = widget.NewButton("\n", func() { t.SendDiscard(11) })
	t.btCards[12] = widget.NewButton("\n", func() { t.SendDiscard(12) })
	t.btCards[13] = widget.NewButton("\n", func() { t.SendDiscard(13) })
	t.btCards[14] = widget.NewButton("\n", func() { t.SendDiscard(14) })
	t.SetButton(t.btCards, false)

	t.btWinCards = make([]*widget.Button, 15)
	t.btWinCards[0] = widget.NewButton("", func() { t.SendWin(0) })
	t.btWinCards[1] = widget.NewButton("", func() { t.SendWin(1) })
	t.btWinCards[2] = widget.NewButton("", func() { t.SendWin(2) })
	t.btWinCards[3] = widget.NewButton("", func() { t.SendWin(3) })
	t.btWinCards[4] = widget.NewButton("", func() { t.SendWin(4) })
	t.btWinCards[5] = widget.NewButton("", func() { t.SendWin(5) })
	t.btWinCards[6] = widget.NewButton("", func() { t.SendWin(6) })
	t.btWinCards[7] = widget.NewButton("", func() { t.SendWin(7) })
	t.btWinCards[8] = widget.NewButton("", func() { t.SendWin(8) })
	t.btWinCards[9] = widget.NewButton("", func() { t.SendWin(9) })
	t.btWinCards[10] = widget.NewButton("", func() { t.SendWin(10) })
	t.btWinCards[11] = widget.NewButton("", func() { t.SendWin(11) })
	t.btWinCards[12] = widget.NewButton("", func() { t.SendWin(12) })
	t.btWinCards[13] = widget.NewButton("", func() { t.SendWin(13) })
	t.btWinCards[14] = widget.NewButton("", func() { t.SendWin(14) })
	t.SetButton(t.btWinCards, false)

	t.TableInfo = widget.NewLabel("")
	t.SelfInfo = widget.NewLabel("")

	cardInfo := container.NewHBox(
		container.NewVBox(widget.NewLabel("手牌\n棄牌"), widget.NewLabel("棄牌贏")),
		container.NewVBox(t.btCards[0], t.btWinCards[0]),
		container.NewVBox(t.btCards[1], t.btWinCards[1]),
		container.NewVBox(t.btCards[2], t.btWinCards[2]),
		container.NewVBox(t.btCards[3], t.btWinCards[3]),
		container.NewVBox(t.btCards[4], t.btWinCards[4]),
		container.NewVBox(t.btCards[5], t.btWinCards[5]),
		container.NewVBox(t.btCards[6], t.btWinCards[6]),
		container.NewVBox(t.btCards[7], t.btWinCards[7]),
		container.NewVBox(t.btCards[8], t.btWinCards[8]),
		container.NewVBox(t.btCards[9], t.btWinCards[9]),
		container.NewVBox(t.btCards[10], t.btWinCards[10]),
		container.NewVBox(t.btCards[11], t.btWinCards[11]),
		container.NewVBox(t.btCards[12], t.btWinCards[12]),
		container.NewVBox(t.btCards[13], t.btWinCards[13]),
		container.NewVBox(t.btCards[14], t.btWinCards[14]),
	)

	deskInfo := container.NewHBox(
		widget.NewLabel("棄牌堆"), t.btDecks[DECK_DISCARD],
		widget.NewLabel("牌堆"), t.btDecks[DECK_PICKCARD], t.TableInfo)

	section := container.NewVBox(
		cardInfo,
		// winInfo,
		deskInfo,
	)
	c.Add(section)
}

// Send Commands
func (t *OkeyClient) SendPickCard(index int) (bool, error) {
	var data struct {
		PlayOperate struct {
			Instruction int         `json:"instruction"`
			Data        interface{} `json:"data"`
		}
	}
	data.PlayOperate.Instruction = 13
	data.PlayOperate.Data = index
	if index == DECK_DISCARD {
		t.btDecks[DECK_DISCARD].SetText("")
	}
	t.SetButton(t.btDecks, false)
	return t.SendMessage(data)
}

func (t *OkeyClient) SendDiscard(index int) (bool, error) {
	var data struct {
		PlayOperate struct {
			Instruction int         `json:"instruction"`
			Data        interface{} `json:"data"`
		}
	}
	data.PlayOperate.Instruction = 14
	data.PlayOperate.Data = t.HoleCards[index]
	t.btCards[index].SetText("\n")
	t.SetButton(t.btCards, false)
	t.SetButton(t.btDecks, false)
	return t.SendMessage(data)
}

func (t *OkeyClient) SendWin(index int) (bool, error) {
	var data struct {
		PlayOperate struct {
			Instruction int         `json:"instruction"`
			Data        interface{} `json:"data"`
		}
	}
	data.PlayOperate.Instruction = 15
	data.PlayOperate.Data = t.HoleCards[index]
	t.SetButton(t.btCards, false)
	t.SetButton(t.btWinCards, false)
	t.SetButton(t.btDecks, false)
	return t.SendMessage(data)
}

const (
	ACT_TABLE_STATUS = "ActTableStatus"
	ACT_HOLE_CARD    = "ActHoleCard"
	ACT_PICKCARD     = "ActPickCard"
	ACT_DISCARD      = "ActDiscard"
	ACT_SETTLE_INFO  = "ActSettleInfo"
)

func (t *OkeyClient) CheckGameResponse(response *utils.RespBase) bool {
	if t.CheckMatchResponse(response) {
		return true
	}

	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return true
	}

	// fmt.Printf("%+v\n", response)

	if response.Ret == client.ACT_GAME_PERIOD {
		switch t.Fsm {
		case "Play":
			t.GetCurrentSeat(info)
			if t.CurSeat == t.Seat {
				t.SetButton(t.btDecks, true)
			}
			t.UpdateTableInfo()
		case "Discard":
			t.GetCurrentSeat(info)
			t.SetButton(t.btDecks, false)
			/* if t.CurSeat == t.PrevSeat {
				t.btDecks[DECK_DISCARD].Enable()
			} */
			// t.UpdateTableInfo()
		case "Settle":
			t.TableInfo.SetText("Settle")
		}
		return true
	}

	switch response.Ret {
	case client.RET_JOIN_GAME:
		t.CurSeat = UNDEFINED
		t.HoleCards = []int{}
		t.PickCard = UNDEFINED
		t.GetJoinInfo(info)
		return true
	case ACT_TABLE_STATUS:
		t.GetTableInfo(info)
		return true
	case ACT_HOLE_CARD:
		t.GetHoleCardInfo(info)
		return true
	case ACT_PICKCARD:
		t.GetPickCardInfo(info)
		return true
	case ACT_DISCARD:
		t.GetDiscardInfo(info)
		return true
	case ACT_SETTLE_INFO:
		t.GetSettleInfo(info)
		return true
	}

	return false
}

func (t *OkeyClient) GetJoinInfo(data map[string]interface{}) {
	t.Seat = int(data["OwnSeat"].(float64))
	t.PrevSeat = (t.Seat + 3) % 4
	if data["PlayerInfo"] != nil {
		for _, info := range data["PlayerInfo"].([]interface{}) {
			seatInfo := info.(map[string]interface{})
			if int(seatInfo["seatId"].(float64)) == t.Seat {
			}
		}
	}
}

func (t *OkeyClient) GetTableInfo(data map[string]interface{}) {
	if data["tablePlayerInfo"] != nil {
		t.UpdateTableInfo()
	}
}

// Get Wild & Discard
func (t *OkeyClient) GetHoleCardInfo(data map[string]interface{}) {
	items := []string{"Runs", "ImRuns", "Sets", "ImSets", "CouldRuns", "CouldSets", "Pairs", "ImPairs"}
	t.HoleCards = []int{}

	for _, item := range items {
		if info, ok := data[item].([]interface{}); ok {
			for _, cards := range info {
				t.GetCardData(cards)
			}
		}
	}

	if cards, ok := data["Others"].(interface{}); ok {
		t.GetCardData(cards)
	}

	t.UpdateHoleCards(false, false) // Deal Card
}

func (t *OkeyClient) GetCardData(cards interface{}) {
	for _, value := range cards.([]interface{}) {
		card := int(value.(float64))
		t.HoleCards = append(t.HoleCards, card)
	}
}

func (t *OkeyClient) GetPickCardInfo(data map[string]interface{}) {
	num := int(data["Num"].(float64))
	if num >= 0 {
		t.btDecks[DECK_PICKCARD].SetText(fmt.Sprintf("%d", num))
	}

	card := int(data["Card"].(float64))
	if card == -1 {
		return
	}
	seat := int(data["SeatId"].(float64))
	if seat == t.Seat {
		t.HoleCards = append(t.HoleCards, card)
		t.UpdateHoleCards(false, false) // PickCard
		/* t.btCards[14].SetText(t.PrintCard(card))
		t.SetButton(t.btCards, true) */

		if cards, ok := data["DiscardWin"].([]interface{}); ok {
			for _, value := range cards {
				card := int(value.(float64))
				t.DiscardWin = append(t.DiscardWin, card)
			}
			if len(t.DiscardWin) > 0 {
				t.TableInfo.SetText(fmt.Sprintf("Discard Win %v", t.DiscardWin))
				for i, card := range t.HoleCards {
					if ItemExist(t.DiscardWin, card) {
						t.btWinCards[i].SetText("Win")
						t.btWinCards[i].Enable()
					}
				}
			}
		} else {
			t.DiscardWin = []int{}
		}

		t.UpdateWinCards()
	}
}

// Discard
func (t *OkeyClient) GetDiscardInfo(data map[string]interface{}) {
	seat := -1
	if data["SeatId"] != nil {
		seat = int(data["SeatId"].(float64))
	}

	t.Discard = int(data["Card"].(float64))
	if seat == t.Seat {
		for i, card := range t.HoleCards {
			if card == t.Discard {
				t.HoleCards = append(t.HoleCards[:i], t.HoleCards[i+1:]...)
				if win, ok := data["Win"].(bool); ok && win {
					t.UpdateHoleCards(true, true)
				} else {
					t.UpdateHoleCards(true, false) // Discard
				}
				return
			}
		}
	} else if seat == t.PrevSeat {
		t.btDecks[DECK_DISCARD].SetText(t.PrintCard(t.Discard))
		t.SetButton(t.btDecks, true)
	}
}

func (t *OkeyClient) GetCurrentSeat(data map[string]interface{}) {
	if info, ok := data["SeatId"].(interface{}); ok {
		t.CurSeat = int(info.(float64))
		t.UpdateTableInfo()
	}
}

func (t *OkeyClient) GetSettleInfo(data map[string]interface{}) {
	_, ok := data["SettleInfo"].([]interface{})
	if !ok {
		return
	}
	if data["GameEnd"] != nil {
		if data["GameEnd"].(bool) {
			t.TableInfo.SetText("[Game End]")
		}
	} else {
		t.TableInfo.SetText("[Winner] ...")
	}
	return
}

func (t *OkeyClient) UpdateTableInfo() {
	if t.CurSeat < 0 {
		return
	}
	t.TableInfo.SetText(fmt.Sprintf("%d/%d", t.CurSeat, t.Seat))
}

func (t *OkeyClient) UpdateHoleCards(reSort, win bool) {
	if len(t.HoleCards) <= 0 {
		return
	}
	winText := ""
	if win {
		winText = "Win"
	}
	if reSort {
		sort.Ints(t.HoleCards)
	}
	for i, card := range t.HoleCards {
		t.btCards[i].SetText(t.PrintCard(card))
	}

	if len(t.HoleCards) == PLAYER_HOLECARD_NUM {
		t.SetButton(t.btCards, false)
		t.btCards[14].SetText("\n")
		for i, _ := range t.HoleCards {
			t.btWinCards[i].SetText(winText)
		}
	} else if len(t.HoleCards) == DEALER_START_HOLECARD {
		t.SetButton(t.btCards, true)
		t.SetButton(t.btDecks, false)
		t.UpdateWinCards()
	}

	return
}

func (t *OkeyClient) UpdateWinCards() {
	t.SetButton(t.btWinCards, false)
	if len(t.DiscardWin) == 0 {
		return
	}
	for i, card := range t.HoleCards {
		if ItemExist(t.DiscardWin, card) {
			t.btWinCards[i].Enable()
		}
	}
}

func (t *OkeyClient) SetButton(buttons []*widget.Button, enable bool) {
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

var (
	CardSuit = [SUIT_COUNT - 1]string{"藍", "紅", "黑", "黃"}
)

func (t *OkeyClient) PrintCard(card int) (cardStr string) {
	if card == 52 {
		return "(OK)\n52"
	}
	suit := card / 13
	point := card % 13
	cardStr = CardSuit[suit] + " " + strconv.Itoa(point+1) + "\n" + strconv.Itoa(card)

	return
}

func ItemExist(data []int, item int) bool {
	for _, d := range data {
		if d == item {
			return true
		}
	}
	return false
}
