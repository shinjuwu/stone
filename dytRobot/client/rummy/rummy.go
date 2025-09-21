package rummy

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

type RummyClient struct {
	*client.BaseMatchClient
	Seat        int
	Seats       []int
	CurSeat     int
	Cards       []int
	Points      int
	Hints       map[int]int
	BestCards   []int
	WildPoint   int
	DropPoint   int
	Discard     int
	Win         bool
	PickCard    int
	PickCardWin int
	DiscardWin  int

	btDecks   []*widget.Button
	btCards   []*widget.Button
	SelfInfo  *widget.Label
	TableInfo *widget.Label
}

const (
	UNDEFINED     = -1
	DECK_PICKCARD = 0
	DECK_DISCARD  = 1
	DROP_CARD     = 2
)

func NewClient(setting client.ClientConfig) *RummyClient {
	matchClient := client.NewMatchClient(setting)
	t := &RummyClient{
		BaseMatchClient: matchClient,
	}

	t.CheckResponse = t.CheckGameResponse

	/* t.CustomMessage = append(t.CustomMessage, "{\"MatchGameBet\":{\"BetInfo\":1}}")
	t.CustomMessage = append(t.CustomMessage, "{\"PlayOperate\":{\"instruction\":7}}")
	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[11,10,9,25,38],[12,8],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage) */
	return t
}

func (t *RummyClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateMatchSection(c)
	t.CreateGameSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *RummyClient) CreateGameSection(c *fyne.Container) {
	t.btDecks = make([]*widget.Button, 3)
	t.btDecks[DECK_PICKCARD] = widget.NewButton("[D]", func() { t.SendPickCard(0) })
	t.btDecks[DECK_DISCARD] = widget.NewButton("", func() { t.SendPickCard(1) })
	t.btDecks[DROP_CARD] = widget.NewButton("[Drop]\nP.20", func() { t.SendDrop() })
	t.SetButton(t.btDecks, false)

	t.btCards = make([]*widget.Button, 14)
	t.btCards[0] = widget.NewButton("", func() { t.SendDiscard(0) })
	t.btCards[1] = widget.NewButton("", func() { t.SendDiscard(1) })
	t.btCards[2] = widget.NewButton("", func() { t.SendDiscard(2) })
	t.btCards[3] = widget.NewButton("", func() { t.SendDiscard(3) })
	t.btCards[4] = widget.NewButton("", func() { t.SendDiscard(4) })
	t.btCards[5] = widget.NewButton("", func() { t.SendDiscard(5) })
	t.btCards[6] = widget.NewButton("", func() { t.SendDiscard(6) })
	t.btCards[7] = widget.NewButton("", func() { t.SendDiscard(7) })
	t.btCards[8] = widget.NewButton("", func() { t.SendDiscard(8) })
	t.btCards[9] = widget.NewButton("", func() { t.SendDiscard(9) })
	t.btCards[10] = widget.NewButton("", func() { t.SendDiscard(10) })
	t.btCards[11] = widget.NewButton("", func() { t.SendDiscard(11) })
	t.btCards[12] = widget.NewButton("", func() { t.SendDiscard(12) })
	t.btCards[13] = widget.NewButton("", func() { t.SendDiscard(13) })
	t.SetButton(t.btCards, false)

	t.TableInfo = widget.NewLabel("")
	t.SelfInfo = widget.NewLabel("")

	/* section := container.NewHBox(
	t.btCards[0], t.btCards[1], t.btCards[2], t.btCards[3], t.btCards[4],
	t.btCards[5], t.btCards[6], t.btCards[7], t.btCards[8], t.btCards[9],
	t.btCards[10], t.btCards[11], t.btCards[12], t.btCards[13], t.SelfInfo,
	t.btDecks[0], t.btDecks[1], t.btDecks[2], t.TableInfo) */

	selfInfo := container.NewHBox(
		t.SelfInfo,
		t.btCards[0], t.btCards[1], t.btCards[2], t.btCards[3], t.btCards[4],
		t.btCards[5], t.btCards[6], t.btCards[7], t.btCards[8], t.btCards[9],
		t.btCards[10], t.btCards[11], t.btCards[12], t.btCards[13])
	deskInfo := container.NewHBox(
		t.btDecks[0], t.btDecks[1], t.btDecks[2], t.TableInfo)
	section := container.NewVBox(
		selfInfo,
		deskInfo,
	)
	c.Add(section)
}

// Send Commands
func (t *RummyClient) SendPickCard(index int) (bool, error) {
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

func (t *RummyClient) SendDiscard(index int) (bool, error) {
	var data struct {
		PlayOperate struct {
			Instruction int         `json:"instruction"`
			Data        interface{} `json:"data"`
		}
	}
	data.PlayOperate.Instruction = 14
	data.PlayOperate.Data = t.Cards[index]
	t.btCards[index].SetText("")
	t.SetButton(t.btCards, false)
	t.SetButton(t.btDecks, false)
	return t.SendMessage(data)
}

func (t *RummyClient) SendDrop() (bool, error) {
	var data struct {
		PlayOperate struct {
			Instruction int `json:"instruction"`
		}
	}
	data.PlayOperate.Instruction = 15
	t.SetButton(t.btDecks, false)
	return t.SendMessage(data)
}

const (
	ACT_TABLE_STATUS = "ActTableStatus"
	ACT_DEAL_CARD    = "ActDealCard"
	ACT_DROP_POINTS  = "ActDropPoints"
	ACT_PICKCARD     = "ActPickCard"
	ACT_PICKCARD_WIN = "ActPickCardWin"
	ACT_DISCARD      = "ActDiscard"
	ACT_DROP         = "ActDrop"
	ACT_BESTCARD     = "ActBestCard"
	ACT_SETTLE_INFO  = "ActSettleInfo"
)

func (t *RummyClient) CheckGameResponse(response *utils.RespBase) bool {
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
		case "PickCard":
			t.GetCurrentSeat(info)
			if t.CurSeat == t.Seat {
				if t.PickCardWin == UNDEFINED {
					t.SetButton(t.btDecks, true)
				} else {
					t.btDecks[DECK_DISCARD].Enable()
				}
			}
			t.UpdateTableInfo()
		case "Discard":
			t.GetCurrentSeat(info)
			t.SetButton(t.btDecks, false)
			if t.CurSeat == t.Seat {
				t.btDecks[DROP_CARD].Enable()
			}
			// t.UpdateTableInfo()
		}
		return true
	}

	switch response.Ret {
	case client.RET_JOIN_GAME:
		t.CurSeat = UNDEFINED
		t.Cards = []int{}
		t.PickCard = UNDEFINED
		t.PickCardWin = UNDEFINED
		t.DiscardWin = UNDEFINED
		t.Hints = make(map[int]int)
		t.GetJoinInfo(info)
		return true
	case ACT_TABLE_STATUS:
		t.GetTableInfo(info)
		return true
	case ACT_DEAL_CARD:
		t.GetDealCard(info)
		return true
	case ACT_DROP_POINTS:
		t.btDecks[DROP_CARD].SetText("[Drop]\nP.40")
		return true
	case ACT_PICKCARD:
		t.GetPickCard(info)
		return true
	case ACT_PICKCARD_WIN:
		t.GetPickCardWin(info)
		return true
	case ACT_DISCARD:
		t.GetDiscard(info)
		return true
	case ACT_DROP:
		t.GetDrop(info)
		return true
	case ACT_BESTCARD:
		t.GetBestCard(info)
		return true
	case ACT_SETTLE_INFO:
		t.GetSettleInfo(info)
		return true
	}

	return false
}

func (t *RummyClient) GetJoinInfo(data map[string]interface{}) {
	t.Seat = int(data["OwnSeat"].(float64))
	if data["PlayerInfo"] != nil {
		for _, info := range data["PlayerInfo"].([]interface{}) {
			seatInfo := info.(map[string]interface{})
			if int(seatInfo["seatId"].(float64)) == t.Seat {
			}
		}
	}
	if data["WildCard"] != nil {
		wild := int(data["WildCard"].(float64))
		if wild >= 52 {
			t.WildPoint = 0
		} else {
			t.WildPoint = wild % 13
		}
	}
	if data["Discards"] != nil {

	}
	if data["BetCards"] != nil {
		t.GetBestCard(data)
	} else {
		for i := 0; i < 14; i++ {
			t.btCards[i].SetText("")
		}
	}
	if data["DropPoints"] != nil {
		t.btDecks[DROP_CARD].SetText(fmt.Sprintf("[Drop]\nP.%d", int(data["DropPoints"].(float64))))
	}
}

func (t *RummyClient) GetTableInfo(data map[string]interface{}) {
	if data["tablePlayerInfo"] != nil {
		t.Seats = []int{}
		for _, info := range data["tablePlayerInfo"].([]interface{}) {
			seatInfo := info.(map[string]interface{})
			if seatInfo["seatId"] != nil {
				seat := int(seatInfo["seatId"].(float64))
				t.Seats = append(t.Seats, seat)
			}
		}
		t.UpdateSeatInfo()
	}
}

// Get Wild & Discard
func (t *RummyClient) GetDealCard(data map[string]interface{}) {
	if data["Wild"] != nil {
		wild := int(data["Wild"].(float64))
		if wild >= 52 {
			t.WildPoint = 0
		} else {
			t.WildPoint = wild % 13
		}
	}
	if data["Discard"] != nil {
		t.Discard = int(data["Discard"].(float64))
		t.btDecks[DECK_DISCARD].SetText(t.PrintCard(t.Discard))
	}
}

func (t *RummyClient) GetPickCard(data map[string]interface{}) {
	card := int(data["Card"].(float64))
	if card == -1 {
		return
	}
	seat := int(data["SeatId"].(float64))
	if seat == t.Seat {
		t.PickCard = card
	}
}

func (t *RummyClient) GetPickCardWin(data map[string]interface{}) {
	card := int(data["Card"].(float64))
	if card == -1 {
		return
	}
	t.PickCardWin = card

	selfInfo := fmt.Sprintf("Point %d\nSeat %d\nDiscard [Win]", t.Points, t.Seat)

	t.SelfInfo.SetText(selfInfo)
}

func (t *RummyClient) GetDiscard(data map[string]interface{}) {
	t.Discard = int(data["Card"].(float64))
	discardInfo := t.PrintCard(t.Discard)
	if t.Discard == t.PickCardWin {
		discardInfo = ">" + discardInfo + "<"
	}
	t.btDecks[DECK_DISCARD].SetText(discardInfo)

	if data["Win"] != nil {
		win := data["Win"].(bool)
		if win {
			seat := int(data["SeatId"].(float64))
			if seat == t.Seat {
				t.SetButton(t.btCards, false)
				t.btCards[13].SetText("")
				t.DiscardWin = UNDEFINED
			}
			t.TableInfo.SetText(fmt.Sprintf("[Winner] ..... Seat %d", int(data["SeatId"].(float64))))
		}
	}
}

func (t *RummyClient) GetDrop(data map[string]interface{}) {
	if data["SeatId"] != nil {
		dropSeat := int(data["SeatId"].(float64))
		for i, seat := range t.Seats {
			if dropSeat == seat {
				t.Seats = append(t.Seats[0:i], t.Seats[i+1])
				t.UpdateSeatInfo()
				break
			}
		}
	}
}

func (t *RummyClient) GetBestCard(data map[string]interface{}) {
	if info, ok := data["Hints"].(map[string]interface{}); ok {
		t.Hints = map[int]int{}
		for card, value := range info {
			c, _ := strconv.Atoi(card)
			num := int(value.(float64))
			t.Hints[c] = num
		}
	}
	if info, ok := data["Points"].(interface{}); ok {
		t.Points = int(info.(float64))
	} else {
		t.Points = -1
	}
	t.UpdateTableInfo()

	if discard, ok := data["Discard"].(interface{}); ok {
		t.DiscardWin = int(discard.(float64))
	}

	index := 0

	items := []string{"Runs", "ImpureRuns", "Sets", "ImpureSets"}
	groups := []string{" Runs", "iRuns", " Sets", "iSets"}
	t.Cards = []int{}
	for it, item := range items {
		if info, ok := data[item].([]interface{}); ok {
			j := 0
			for _, cards := range info {
				l := len(cards.([]interface{})) - 1
				for k, value := range cards.([]interface{}) {
					card := int(value.(float64))
					t.Cards = append(t.Cards, card)
					cardInfo := ""
					if len(t.Hints) > 0 {
						for hint, num := range t.Hints {
							if hint == card {
								cardInfo = strconv.Itoa(num)
							}
						}
						cardInfo += "\n"
					}

					if card == t.PickCard {
						cardInfo += "*"
					}
					// cardInfo += t.PrintCard(card) + "\n" + groups[it] + "." + strconv.Itoa(j)
					cardInfo += t.PrintCard(card) + "\n"
					if k == l {
						cardInfo += groups[it] + "." + strconv.Itoa(j)
					}
					t.btCards[index].SetText(cardInfo)
					if t.Fsm == "PickCard" && t.DiscardWin == UNDEFINED {
						t.btCards[index].Enable()
					}
					index++
				}
				j++
			}
		}
	}

	if info, ok := data["Others"].([]interface{}); ok {
		others := []int{}
		for _, value := range info {
			card := int(value.(float64))
			others = append(others, card)
		}
		sort.Ints(others)
		for _, card := range others {
			cardInfo := ""
			if len(t.Hints) > 0 {
				for hint, num := range t.Hints {
					if hint == card {
						cardInfo = strconv.Itoa(num)
					}
				}
				cardInfo += "\n"
			}

			if card == t.PickCard {
				cardInfo += "*"
			}
			cardInfo += t.PrintCard(card) + "\n"

			t.btCards[index].SetText(cardInfo)
			if t.Fsm == "PickCard" && t.DiscardWin == UNDEFINED {
				t.btCards[index].Enable()
			}
			index++
		}
		t.Cards = append(t.Cards, others...)
	}

	if index == 13 {
		t.btCards[index].SetText("")
		t.btCards[index].Disable()
		t.PickCard = UNDEFINED
	} else if index == 14 {
		t.btDecks[DROP_CARD].Enable()
	}

	selfInfo := fmt.Sprintf("Point %d\nSeat %d", t.Points, t.Seat)

	if t.DiscardWin != UNDEFINED {
		selfInfo += "\n[Win]"

		if index == 13 {
			t.Cards = append(t.Cards, t.DiscardWin)
			cardInfo := "\n>" + t.PrintCard(t.DiscardWin) + "<\n"
			t.btCards[index].SetText(cardInfo)
			t.btCards[index].Enable()
		}
	}

	t.SelfInfo.SetText(selfInfo)
}

func (t *RummyClient) GetCurrentSeat(data map[string]interface{}) {
	if info, ok := data["SeatId"].(interface{}); ok {
		t.CurSeat = int(info.(float64))
		t.UpdateSeatInfo()
	}
}

func (t *RummyClient) GetSettleInfo(data map[string]interface{}) {
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

func (t *RummyClient) UpdateTableInfo() {
	if t.WildPoint < 0 {
		t.TableInfo.SetText("")
		return
	}
	wildInfo := "Wild (" + CardPoint[t.WildPoint] + ")"
	if t.CurSeat >= 0 {
		wildInfo += "\n" + t.UpdateSeatInfo()
	}

	t.TableInfo.SetText(wildInfo)
}

func (t *RummyClient) UpdateSeatInfo() (seatStr string) {
	if len(t.Seats) <= 1 {
		return
	}

	sort.Ints(t.Seats)
	for _, seat := range t.Seats {
		if len(seatStr) != 0 {
			seatStr += " -> "
		}
		if t.CurSeat == seat {
			seatStr += "[" + strconv.Itoa(seat) + "]"
		} else {
			seatStr += strconv.Itoa(seat)
		}
	}
	return
}

func (t *RummyClient) SetButton(buttons []*widget.Button, enable bool) {
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

var CardSuit = [4]string{"C", "d", "h", "S"}
var CardPoint = [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

func (t *RummyClient) PrintCard(card int) (cardStr string) {
	if card == 53 {
		return "(JK)"
	} else if card == 52 {
		return "(jk)"
	}
	suit := card / 13
	point := card % 13
	cardStr = CardSuit[suit] + "." + CardPoint[point]
	if point == t.WildPoint {
		return "(" + cardStr + ")"
	}
	return
}
