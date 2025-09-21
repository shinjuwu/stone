package texas

import (
	"dytRobot/client"
	"dytRobot/utils"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type TexasClient struct {
	*client.BaseMatchClient
	SeatId    int
	CurSeatId int
	Call      float64
	LockGold  float64
	Action    string
	CardType  int
	Cards     []int
	BestCards []int
	Board     []int
	Button    []*widget.Button
	CardsInfo *widget.Label
}

func NewClient(setting client.ClientConfig) *TexasClient {
	matchClient := client.NewMatchClient(setting)
	t := &TexasClient{
		BaseMatchClient: matchClient,
	}

	t.CheckResponse = t.CheckGameResponse

	t.CustomMessage = append(t.CustomMessage, "{\"MatchGameBet\":{\"BetInfo\":1}}")
	t.CustomMessage = append(t.CustomMessage, "{\"PlayOperate\":{\"instruction\":7}}")
	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[11,10,9,25,38],[12,8],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}

func (t *TexasClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateMatchSection(c)
	t.CreateGameSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *TexasClient) CreateGameSection(c *fyne.Container) {
	t.Button = make([]*widget.Button, 4)
	t.Button[0] = widget.NewButton("Check/Call", func() {
		t.SendCall()
	})

	entryBet := widget.NewEntry()
	t.Button[1] = widget.NewButton("Bet", func() {
		bet, err := strconv.ParseFloat(entryBet.Text, 64)
		if err != nil {
			return
		}
		t.SendMatchGameBet(bet)
	})

	//下注倍數按鈕
	t.Button[2] = widget.NewButton("AllIn", func() {
		t.SetButton(false)
		t.SendMatchGameBet(t.LockGold)
	})

	//攤牌按鈕
	t.Button[3] = widget.NewButton("Fold", func() {
		t.SendFold()
	})

	t.CardsInfo = widget.NewLabel("")
	section := container.NewHBox(t.Button[0], t.Button[1], entryBet, t.Button[2], t.Button[3], t.CardsInfo)
	c.Add(section)

	t.SetButton(false)
}

func (t *TexasClient) SendFold() (bool, error) {
	var data struct {
		PlayOperate struct {
			Instruction int `json:"instruction"`
		}
	}
	data.PlayOperate.Instruction = 7
	t.SetButton(false)
	return t.SendMessage(data)
}

func (t *TexasClient) SendCall() (bool, error) {
	t.SetButton(false)
	return t.SendMatchGameBet(t.Call)

}

func (t *TexasClient) SendMatchGameBet(bet float64) (bool, error) {
	var data struct {
		MatchGameBet struct {
			BetInfo string `json:"BetInfo"`
		}
	}

	data.MatchGameBet.BetInfo = strconv.FormatFloat(bet, 'f', 4, 64)
	return t.SendMessage(data)
}

const (
	ACT_SELF_CARD     = "ActSelfCard"
	ACT_BOARD         = "ActBoard"
	ACT_ACTION_SEAT   = "ActActionSeat"
	ACT_ACTION_INFO   = "ActActionInfo"
	ACT_BESTCARD_INFO = "ActBestCardInfo"
	ACT_SETTLE_INFO   = "ActSettleInfo"
)

func (t *TexasClient) CheckGameResponse(response *utils.RespBase) bool {
	if t.CheckMatchResponse(response) {
		return true
	}

	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return true
	}

	switch response.Ret {
	case client.RET_JOIN_GAME:
		t.Board = []int{}
		t.Cards = []int{}
		t.Action = ""
		t.GetJoinInfo(info)
		return true
	case client.ACT_GOLD:
		t.GetGoldInfo(info)
		return true
	case ACT_SELF_CARD:
		t.GetSelfCard(info)
		return true
	case ACT_BESTCARD_INFO:
		t.GetBestCardInfo(info)
		return true
	case ACT_BOARD:
		t.GetBoardCards(info)
		return true
	case ACT_ACTION_SEAT:
		t.AddTableStatus(t.GetActionSeat(info))
		return true
	case ACT_ACTION_INFO:
		t.AddTableStatus(t.GetActionInfo(info))
		return true
	case ACT_SETTLE_INFO:
		t.AddTableStatus(t.GetSettleInfo(info))
	}

	return false
}

var CardTypeName = [11]string{
	"散牌",
	"一對",
	"兩對",
	"三條",
	"順子",
	"同花",
	"葫蘆",
	"鐵支",
	"同花順",
	"同花大順",
	"比牌前贏",
}

func (t *TexasClient) GetJoinInfo(data map[string]interface{}) {
	t.SeatId = int(data["OwnSeat"].(float64))
	if data["Board"] != nil {
		t.Board = []int{}
		info := data["Board"].([]interface{})
		for _, card := range info {
			t.Board = append(t.Board, int(card.(float64)))
		}
	}
	if data["PlayerInfo"] != nil {
		for _, info := range data["PlayerInfo"].([]interface{}) {
			seatInfo := info.(map[string]interface{})
			if int(seatInfo["seatId"].(float64)) == t.SeatId {
				if seatInfo["gold"] != nil {
					t.LockGold = seatInfo["gold"].(float64)
					break
				}
			}
		}
	}
	if data["BetCards"] != nil {
		for _, info := range data["BetCards"].([]interface{}) {
			seatInfo := info.(map[string]interface{})
			if int(seatInfo["SeatId"].(float64)) == t.SeatId {
				if seatInfo["Cards"] != nil {
					data := seatInfo["Cards"].([]interface{})
					for _, card := range data {
						t.Cards = append(t.Cards, int(card.(float64)))
					}
				}
				if seatInfo["Type"] != nil {
					t.CardType = int(seatInfo["Type"].(float64))
				}
				if t.CardType > 0 {
					if seatInfo["BestCards"] != nil {
						data := seatInfo["BestCards"].([]interface{})
						t.BestCards = []int{}
						for _, card := range data {
							t.BestCards = append(t.BestCards, int(card.(float64)))
						}
					}
				}
				break
			}
		}
	}

	t.UpdateMyInfo()
}

func (t *TexasClient) GetGoldInfo(data map[string]interface{}) {
	if t.Fsm == "Settle" {
		return
	}
	t.LockGold = data["Gold"].(float64)
	t.UpdateMyInfo()
}

func (t *TexasClient) GetSelfCard(data map[string]interface{}) {
	info := data["Cards"].([]interface{})
	for _, card := range info {
		t.Cards = append(t.Cards, int(card.(float64)))
	}

	t.CardType = int(data["Type"].(float64))
	t.UpdateMyInfo()
}

func (t *TexasClient) GetBestCardInfo(data map[string]interface{}) {
	t.CardType = int(data["Type"].(float64))
	if t.CardType != 0 {
		info := data["Cards"].([]interface{})
		t.BestCards = []int{}
		for _, card := range info {
			t.BestCards = append(t.BestCards, int(card.(float64)))
		}
	}
	t.UpdateMyInfo()
}

func (t *TexasClient) GetActionSeat(data map[string]interface{}) string {
	var enable bool
	t.CurSeatId = int(data["SeatId"].(float64))
	call := data["Call"].(float64)
	str := ">"
	if t.CurSeatId == t.SeatId {
		t.Call = call
		str = "輪到自己"
		enable = true
	}
	t.SetButton(enable)
	t.UpdateMyInfo()
	return fmt.Sprintf("%s Seat[%d] %.4f\n", str, t.CurSeatId, call)
}

func (t *TexasClient) GetActionInfo(data map[string]interface{}) string {
	seatId := int(data["SeatId"].(float64))
	action := data["Action"].(string)
	bet := data["RoundBet"].(float64)
	gold := data["Gold"].(float64)
	pot := data["Pot"].(float64)
	if seatId == t.SeatId {
		t.LockGold = gold
		t.Action = action
		t.UpdateMyInfo()
	}
	return fmt.Sprintf("Seat[%d] [%s] %.4f Gold %.4f Pot %.4f\n", seatId, action, bet, gold, pot)
}

func (t *TexasClient) GetSettleInfo(data map[string]interface{}) string {
	info, ok := data["SettleInfo"].([]interface{})
	if !ok {
		return ""
	}

	message := "牌局結果:\n"
	for _, detail := range info {
		cards := []int{}
		bestCards := []int{}

		seatInfo := detail.(map[string]interface{})
		seatId := int(seatInfo["SeatId"].(float64))
		if seatInfo["Cards"] != nil {
			data := seatInfo["Cards"].([]interface{})
			for _, card := range data {
				cards = append(cards, int(card.(float64)))
			}
		}
		if seatInfo["BestCards"] != nil {
			data := seatInfo["BestCards"].([]interface{})
			for _, card := range data {
				bestCards = append(bestCards, int(card.(float64)))
			}
		}
		cardType := 0
		if seatInfo["Type"] != nil {
			cardType = int(seatInfo["Type"].(float64))
		}
		win := seatInfo["Win"].(float64)
		bet := seatInfo["TotalBet"].(float64)
		gold := seatInfo["Gold"].(float64)
		if seatId == t.SeatId {
			t.LockGold = gold
			t.CardType = cardType
			t.BestCards = bestCards
			t.UpdateMyInfo()
		}
		message += fmt.Sprintf("Seat: %d Hold: %s Type: %s Best: %s WinLose: %.4f Gold: %.4f\n", seatId, utils.PrintCard(cards), CardTypeName[cardType], utils.PrintCard(bestCards), (win - bet), gold)
	}
	return message
}

func (t *TexasClient) GetBoardCards(data map[string]interface{}) {
	info := data["Cards"].([]interface{})
	for _, card := range info {
		t.Board = append(t.Board, int(card.(float64)))
	}
}

func (t *TexasClient) UpdateMyInfo() {
	cardInfo := fmt.Sprintf("[%d/%d] G %.2f", t.SeatId, t.CurSeatId, t.LockGold)

	if len(t.Cards) > 0 {
		cardInfo += fmt.Sprintf(", 牌 %s [%s]", utils.PrintCard(t.Cards), CardTypeName[t.CardType])
		if len(t.BestCards) > 0 {
			cardInfo += fmt.Sprintf(" %s", utils.PrintCard(t.BestCards))
		}
	}
	if len(t.Board) > 0 {
		cardInfo += fmt.Sprintf(", 公 %s", utils.PrintCard(t.Board))
	}
	if len(t.Action) > 0 {
		cardInfo += fmt.Sprintf(", [%s]", t.Action)
	}
	t.CardsInfo.SetText(cardInfo)
}

func (t *TexasClient) SetButton(enable bool) {
	if enable {
		for _, button := range t.Button {
			button.Enable()
		}
	} else {
		for _, button := range t.Button {
			button.Disable()
		}
	}
}
