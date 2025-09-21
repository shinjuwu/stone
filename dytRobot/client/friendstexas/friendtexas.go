package friendstexas

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
	*client.BaseFriendsClient
	SeatId     int
	CurSeatId  int
	Call       float64
	LockGold   float64
	Action     string
	CardType   int
	Cards      []int
	BestCards  []int
	Board      []int
	Button     []*widget.Button
	TableOwner int
	Ready      bool
	BetLimit   float64
	RoundBet   float64
	TotalBet   float64
	Balance    float64
	CardsInfo  *widget.Label
}

const (
	BTN_CHECK_CALL int = iota
	BTN_BET
	BTN_ALLIN
	BTN_FOLD
	BTN_READY
	BTN_START_GAME
	BTN_ROUND_STOP

	BTN_COUNT
)

func NewClient(setting client.ClientConfig) *TexasClient {
	friendsClient := client.NewFriendsClient(setting)
	t := &TexasClient{
		BaseFriendsClient: friendsClient,
	}

	t.CheckResponse = t.CheckGameResponse

	t.CustomMessage = append(t.CustomMessage, "{\"FGameBet\":{\"BetInfo\":1}}")
	t.CustomMessage = append(t.CustomMessage, "{\"FPlayOperate\":{\"instruction\":7}}")
	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[11,10,9,25,38],[12,8],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1],[-1,-1]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}

func (t *TexasClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateFriendSection(c)
	t.CreateGameSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *TexasClient) CreateGameSection(c *fyne.Container) {
	t.Button = make([]*widget.Button, BTN_COUNT)
	t.Button[BTN_READY] = widget.NewButton("UnReady", func() {
		ok, err := t.SendReady()
		if !ok {
			fmt.Println(err)
		}
	})
	t.Button[BTN_CHECK_CALL] = widget.NewButton("Check/Call", func() {
		t.SendCall()
	})

	entryBet := widget.NewEntry()
	t.Button[BTN_BET] = widget.NewButton("Bet", func() {
		bet, err := strconv.ParseFloat(entryBet.Text, 64)
		if err != nil {
			return
		}
		t.SendFriendsGameBet(bet)
	})

	//下注倍數按鈕
	t.Button[BTN_ALLIN] = widget.NewButton("AllIn", func() {
		t.SetBetButtons(false)
		t.SendFriendsGameBet(t.Balance)
	})

	//攤牌按鈕
	t.Button[BTN_FOLD] = widget.NewButton("Fold", func() {
		t.SendFold()
	})

	t.Button[BTN_START_GAME] = widget.NewButton("StartGame", func() {
		t.SendStartGame()
	})

	t.Button[BTN_ROUND_STOP] = widget.NewButton("RoundStop", func() {
		t.SendRoundStop()
	})

	t.CardsInfo = widget.NewLabel("")
	section := container.NewHBox(t.Button[BTN_READY],
		t.Button[BTN_START_GAME], t.Button[BTN_ROUND_STOP],
		t.Button[BTN_CHECK_CALL], t.Button[BTN_BET], entryBet, t.Button[BTN_ALLIN], t.Button[BTN_FOLD], t.CardsInfo)
	c.Add(section)

	t.SetBetButtons(false)
}

func (t *TexasClient) SendFold() (bool, error) {
	var data struct {
		FPlayOperate struct {
			Instruction int `json:"instruction"`
		}
	}
	data.FPlayOperate.Instruction = 7
	t.SetBetButtons(false)
	return t.SendMessage(data)
}

func (t *TexasClient) SendReady() (bool, error) {
	return t.SendFriendsReady(!t.Ready)
}

/* func (t *TexasClient) SendStartGame() (bool, error) {
	return t.SendStartGame()
}

func (t *TexasClient) SendRoundStop() (bool, error) {
	return t.SendRoundStop()
} */

func (t *TexasClient) SendCall() (bool, error) {
	t.SetBetButtons(false)
	return t.SendFriendsGameBet(t.Call)
}

/* func (t *TexasClient) SendFriendsReady() (bool, error) {
	var data struct {
		FReady struct {
			Ready bool `json:"Ready"`
		}
	}

	data.FReady.Ready = !t.Ready
	return t.SendMessage(data)
} */

func (t *TexasClient) SendFriendsGameBet(bet float64) (bool, error) {
	var data struct {
		FGameBet struct {
			BetInfo string `json:"BetInfo"`
		}
	}

	data.FGameBet.BetInfo = strconv.FormatFloat(bet, 'f', 4, 64)
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
	if t.CheckFriendResponse(response) {
		return true
	}

	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return true
	}

	switch response.Ret {
	case client.ACT_GAME_PERIOD:
		if t.Fsm != "Match" {
			t.Button[BTN_READY].Disable()
		}
		if t.SeatId == t.TableOwner {
			if t.Fsm == "Match" {
				t.Button[BTN_START_GAME].Enable()
			} else if t.Fsm == "NextRound" {
				t.Button[BTN_ROUND_STOP].Enable()
			} else {
				t.Button[BTN_START_GAME].Disable()
				t.Button[BTN_ROUND_STOP].Disable()
			}
		}
	case client.RET_FCREATE_GAME, client.RET_FJOIN_GAME:
		t.Board = []int{}
		t.Cards = []int{}
		t.Action = ""
		t.GetJoinInfo(info)
		return true
	case client.ACT_GOLD:
		t.GetGoldInfo(info)
		return true
	case client.ACT_FPLAYER_INFO:
		t.GetFriendsReady(info)
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
	t.Balance = data["Balance"].(float64)
	t.TableOwner = int(data["TableOwner"].(float64))

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
	if data["Setting"] != nil {
		setting := data["Setting"].(map[string]interface{})
		if setting["BetLimit"] != nil {
			t.BetLimit = setting["BetLimit"].(float64)
		}
		/* if setting["InviteCode"] != nil {
			t.EntryInviteCode.SetText(setting["InviteCode"].(string))
			t.EntryInviteCode.Disable()
		}
		if setting["Player"] != nil {
			t.SelectPlayer.SetSelected(strconv.Itoa(int(setting["Player"].(float64))))
			t.SelectPlayer.Disable()
		}
		if setting["Rounds"] != nil {
			t.SelectRound.SetSelected(strconv.Itoa(int(setting["Rounds"].(float64))))
			t.SelectRound.Disable()
		}
		if setting["Ante"] != nil {
			t.SelectAnte.SetSelected(strconv.Itoa(int(setting["Ante"].(float64))))
			t.SelectAnte.Disable()
		}
		if setting["BetSec"] != nil {
			t.SelectBetSec.SetSelected(strconv.Itoa(int(setting["BetSec"].(float64))))
			t.SelectBetSec.Disable()
		}
		if setting["BringInLowerBound"] != nil {
			t.SelectBringInLowerBound.SetSelected(strconv.Itoa(int(setting["BringInLowerBound"].(float64))))
			t.SelectBringInLowerBound.Disable()
		} */
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

func (t *TexasClient) GetFriendsReady(data map[string]interface{}) {
	if t.SeatId != t.TableOwner {
		t.Button[BTN_READY].Enable()
	}
	if t.SeatId == t.TableOwner {
		t.Button[BTN_START_GAME].Enable()
		t.Button[BTN_ROUND_STOP].Enable()
	} else {
		t.Button[BTN_START_GAME].Disable()
		t.Button[BTN_ROUND_STOP].Disable()
	}
	players := data["Players"].([]interface{})
	for _, player := range players {
		info := player.(map[string]interface{})
		if info["seatId"] != nil {
			if int(info["seatId"].(float64)) == t.SeatId {
				t.Ready = info["ready"].(bool)
				if t.Ready {
					t.Button[BTN_READY].SetText("Ready")
				} else {
					t.Button[BTN_READY].SetText("Unready")
				}
			}
		}
	}
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
	t.SetBetButtons(enable)
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
		t.Balance = data["Balance"].(float64)
		t.RoundBet = bet
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

func (t *TexasClient) SetBetButtons(enable bool) {
	if enable {
		for i, button := range t.Button {
			button.Enable()
			if i == BTN_FOLD {
				break
			}
		}
	} else {
		for i, button := range t.Button {
			button.Disable()
			if i == BTN_FOLD {
				break
			}
		}
	}
}
