package cockfight

import (
	"dytRobot/client"
	"dytRobot/utils"
	"encoding/json"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	BET_AREA_RED    uint = iota //0 紅方獲勝
	BET_AREA_BLUE               //1 藍方獲勝
	BET_AREA_TIE                //2 和局
	BET_AREA_BIGTIE             //3 大和局

	BET_AREA_COUNT
)

type CockfightClient struct {
	*client.BaseBetClient
	entryBigtie, entryTie, entryRedwin, entryBluewin *widget.Entry
	labelBaseBet                                     *widget.Label
}

func NewClient(setting client.ClientConfig) *CockfightClient {
	betClient := client.NewBetClient(setting)
	t := &CockfightClient{
		BaseBetClient: betClient,
	}
	t.CheckResponse = t.CheckCockfightResponse

	return t
}

func (t *CockfightClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateBetSection(c)
	t.CreateCFControlSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *CockfightClient) CheckCockfightResponse(response *utils.RespBase) bool {
	return t.CheckBetResponse(response)
}

func (t *CockfightClient) CreateCFControlSection(c *fyne.Container) {
	t.labelBaseBet = widget.NewLabel("底注:")

	labelRedwin := widget.NewLabel("紅贏")
	t.entryRedwin = widget.NewEntry()
	t.entryRedwin.SetText("0")

	labelBluewin := widget.NewLabel("藍贏")
	t.entryBluewin = widget.NewEntry()
	t.entryBluewin.SetText("0")

	labelTie := widget.NewLabel("和局")
	t.entryTie = widget.NewEntry()
	t.entryTie.SetText("0")

	labelBigtie := widget.NewLabel("大和局")
	t.entryBigtie = widget.NewEntry()
	t.entryBigtie.SetText("0")
	buttonBet := widget.NewButton("押注", func() {
		t.SendPlayerAction()
	})

	section1 := container.NewHBox(t.labelBaseBet, labelRedwin, t.entryRedwin, labelBluewin, t.entryBluewin, labelTie, t.entryTie, labelBigtie, t.entryBigtie,
		buttonBet)
	c.Add(section1)

}

func (t *CockfightClient) SendPlayerAction() (bool, error) {

	var BetInfo struct {
		AreaID int `json:"AreaID"`
		Bet    int `json:"Bet"`
	}

	var data struct {
		Bet struct {
			BetInfo string `json:"BetInfo"`
		}
	}

	Bigtie, _ := strconv.Atoi(t.entryBigtie.Text)
	Tie, _ := strconv.Atoi(t.entryTie.Text)
	Bluewin, _ := strconv.Atoi(t.entryBluewin.Text)
	Redwin, _ := strconv.Atoi(t.entryRedwin.Text)
	betInfo := []interface{}{}
	switch {
	case Redwin != 0:
		BetInfo.AreaID = int(BET_AREA_RED)
		BetInfo.Bet = Redwin
		betInfo = append(betInfo, BetInfo)

	case Bluewin != 0:
		BetInfo.AreaID = int(BET_AREA_BLUE)
		BetInfo.Bet = Bluewin
		betInfo = append(betInfo, BetInfo)

	case Tie != 0:
		BetInfo.AreaID = int(BET_AREA_TIE)
		BetInfo.Bet = Tie
		betInfo = append(betInfo, BetInfo)

	case Bigtie != 0:
		BetInfo.AreaID = int(BET_AREA_BIGTIE)
		BetInfo.Bet = Bigtie
		betInfo = append(betInfo, BetInfo)
	}
	b, _ := json.Marshal(betInfo)
	data.Bet.BetInfo = string(b)
	return t.SendMessage(data)
}
