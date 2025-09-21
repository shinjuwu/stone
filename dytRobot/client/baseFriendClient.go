package client

import (
	"dytRobot/constant"
	"dytRobot/utils"
	"fmt"
	"math"
	"strconv"

	"encoding/json"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	RET_FCREATE_GAME = "FCreateGame"
	RET_FJOIN_ROOM   = "FJoinRoom"
	RET_FQUIT_ROOM   = "FQuitRoom"
	RET_FJOIN_GAME   = "FJoinGame"
	RET_FQUIT_GAME   = "FQuitGame"
	RET_FMATCH_STOP  = "FMatchStop"

	ACT_FPLAYER_INFO = "ActFriendsPlayerInfo"
)

type BaseFriendsClient struct {
	*BaseClient
	EntryBringIn, EntryInviteCode                                                *widget.Entry
	SelectPlayer, SelectRound, SelectAnte, SelectBetSec, SelectBringInLowerBound *widget.Select
}

func NewFriendsClient(setting ClientConfig) *BaseFriendsClient {
	baseClient := NewBaseClient(setting)
	t := &BaseFriendsClient{
		BaseClient: baseClient,
	}
	t.CheckResponse = t.CheckFriendResponse

	t.CustomMessage = append(t.CustomMessage, "{\"JoinRoom\":{\"GameID\":"+strconv.Itoa(t.GameId)+"}}")
	t.CustomMessage = append(t.CustomMessage, "{\"JoinGame\":{\"TableId\":"+strconv.Itoa(t.GameId)+"0}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}
func (t *BaseFriendsClient) JoinRoom(gameId int) (bool, error) {
	var data struct {
		FJoinRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.FJoinRoom.GameID = gameId

	return t.SendMessage(data)
}
func (t *BaseFriendsClient) QuitRoom(gameId int) (bool, error) {
	var data struct {
		FQuitRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.FQuitRoom.GameID = gameId

	return t.SendMessage(data)
}

func (t *BaseFriendsClient) CreateGame(tableId int) (bool, error) {
	var data struct {
		FCreateGame struct {
			TableId   int     `json:"TableId"`
			Info      string  `json:"Info"`
			BringGold float64 `json:"bringGold,omitempty"` //Texas & singleWallet
		}
	}
	var info struct {
		PlayerNum         int `json:"playernum"`
		Rounds            int `json:"rounds"`
		Ante              int `json:"ante"`
		BetSec            int `json:"betsec"`
		BringInLowerBound int `json:"bringinlowerbound"`
	}
	data.FCreateGame.TableId = tableId
	info.PlayerNum = t.SelectPlayer.SelectedIndex()
	info.Rounds = t.SelectRound.SelectedIndex()
	info.Ante = t.SelectAnte.SelectedIndex()
	info.BetSec = t.SelectBetSec.SelectedIndex()
	info.BringInLowerBound = t.SelectBringInLowerBound.SelectedIndex()
	infoStr, _ := json.Marshal(info)
	data.FCreateGame.Info = string(infoStr)
	data.FCreateGame.BringGold, _ = strconv.ParseFloat(t.EntryBringIn.Text, 64) // TODO
	// if t.WalletType == constant.SW_TYPE_SINGLE {
	// 	minBringin := t.EnterInfo[tableId]
	// 	swGoldInt := int(t.SWGold)
	// 	bringIn := math.Min(minBringin*100, float64(swGoldInt))
	// 	data.FCreateGame.BringGold = bringIn
	// } else if (t.TableId / 10) == 2004 {
	// 	data.FCreateGame.BringGold = constant.TexasBringGold[(t.TableId % 10)]
	// }

	return t.SendMessage(data)
}

func (t *BaseFriendsClient) JoinGame(tableId int) (bool, error) {
	var data struct {
		FJoinGame struct {
			InviteCode string  `json:"InviteCode"`
			TableId    int     `json:"TableId"`
			BringGold  float64 `json:"BringGold"` // singleWallet
		}
	}
	data.FJoinGame.TableId = tableId
	if t.WalletType == constant.SW_TYPE_SINGLE {
		minBringin := t.EnterInfo[tableId]
		swGoldInt := int(t.SWGold)
		bringIn := math.Min(minBringin*100, float64(swGoldInt))
		data.FJoinGame.BringGold = bringIn
	}
	data.FJoinGame.InviteCode = t.EntryInviteCode.Text

	return t.SendMessage(data)
}

func (t *BaseFriendsClient) QuitGame(tableId int) (bool, error) {
	var data struct {
		FQuitGame struct {
			TableId int `json:"tableId"`
		}
	}
	data.FQuitGame.TableId = tableId

	return t.SendMessage(data)
}

func (t *BaseFriendsClient) MatchStop() (bool, error) {
	var data struct {
		FMatchStop struct {
		}
	}

	return t.SendMessage(data)
}

func (t *BaseFriendsClient) GetTableStatus(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	info, ok := data["tablePlayerInfo"].([]interface{})
	if !ok {
		return ""
	}
	var message string
	for _, player := range info {
		detail, ok := player.(map[string]interface{})
		if !ok {
			return ""
		}
		seatID := int(detail["seatId"].(float64))
		name := detail["name"].(string)
		gold := detail["gold"].(float64)

		message += fmt.Sprintf("座位:%d 名字:%7s 金額:%8.4f\n", seatID, name, gold)
	}
	return message
}

func (t *BaseFriendsClient) CreateFriendSection(c *fyne.Container) {
	//進入遊戲按鈕、房間按鈕
	buttonJoinRoom := widget.NewButton("進入大廳", func() {
		t.JoinRoom(t.GameId)
	})
	buttonQuitRoom := widget.NewButton("離開大廳", func() {
		t.QuitRoom(t.GameId)
	})

	gameLobby := container.NewHBox(buttonJoinRoom, buttonQuitRoom)
	c.Add(gameLobby)
	c.Add(t.LabalRoom)
	c.Add(t.LabelFsm)

	gameRoom := container.NewHBox()
	buttonCreateGame := widget.NewButton("創建牌桌", func() {
		t.TableId = t.GameId * 10
		t.CreateGame(t.TableId)
	})
	gameRoom.Add(buttonCreateGame)
	for roomType := 0; roomType < constant.RoomTypeNum[t.GameId]; roomType++ {
		buttonJoinGame := widget.NewButton("加入房間", nil)
		buttonJoinGame.OnTapped = func() {
			t.SetTableStatus("")
			t.TableId = t.GameId * 10
			t.JoinGame(t.TableId)
		}
		gameRoom.Add(buttonJoinGame)
	}
	buttonQuitGame := widget.NewButton("離開牌桌", func() {
		t.QuitGame(t.TableId)
	})
	buttonMatchStop := widget.NewButton("停止對戰", func() {
		t.MatchStop()
	})
	gameRoom.Add(buttonQuitGame)
	gameRoom.Add(buttonMatchStop)

	// 帶入(單一錢包)
	// 加入 邀請馬

	c.Add(gameRoom)
	c.Add(t.CreateGameSection())
	c.Add(t.JoinGameSection())
	c.Add(t.BringInSection())
	c.Add(t.EntryTableStatus)
}

func (t *BaseFriendsClient) CreateGameSection() (c *fyne.Container) {
	c = container.NewHBox()
	// 人數選項
	labelPlayer := widget.NewLabel("人數:")
	t.SelectPlayer = widget.NewSelect([]string{}, func(value string) {
	})
	c.Add(labelPlayer)
	c.Add(t.SelectPlayer)
	// 局數選項
	labelRound := widget.NewLabel("局數:")
	t.SelectRound = widget.NewSelect([]string{}, func(value string) {
	})
	c.Add(labelRound)
	c.Add(t.SelectRound)
	// 底注選項
	labelAnte := widget.NewLabel("底注:")
	t.SelectAnte = widget.NewSelect([]string{}, func(value string) {
	})
	c.Add(labelAnte)
	c.Add(t.SelectAnte)
	// 下注時間選項
	labelBetSec := widget.NewLabel("下注時間:")
	t.SelectBetSec = widget.NewSelect([]string{}, func(value string) {
	})
	c.Add(labelBetSec)
	c.Add(t.SelectBetSec)
	// 准入選項
	labelBringInLowerBound := widget.NewLabel("准入:")
	t.SelectBringInLowerBound = widget.NewSelect([]string{}, func(value string) {
	})
	c.Add(labelBringInLowerBound)
	c.Add(t.SelectBringInLowerBound)

	labelInfo := widget.NewLabel("  請選擇")
	c.Add(labelInfo)
	return
}

func (t *BaseFriendsClient) JoinGameSection() (c *fyne.Container) {
	c = container.NewHBox()
	// 邀請碼
	labelInviteCode := widget.NewLabel("邀請碼:")
	t.EntryInviteCode = widget.NewEntry()
	t.EntryInviteCode.SetPlaceHolder("Invite Code...")
	t.EntryInviteCode.Resize(fyne.NewSize(150, 36))
	c.Add(labelInviteCode)
	c.Add(container.NewWithoutLayout(t.EntryInviteCode))
	return
}

func (t *BaseFriendsClient) BringInSection() (c *fyne.Container) {
	c = container.NewHBox()
	// 帶入金額
	labelBringIn := widget.NewLabel("帶入金額:")
	t.EntryBringIn = widget.NewEntry()
	c.Add(labelBringIn)
	c.Add(t.EntryBringIn)
	labelInfo := widget.NewLabel("  轉帳錢包無須填入")
	c.Add(labelInfo)
	return
}

func (t *BaseFriendsClient) CheckFriendResponse(response *utils.RespBase) bool {
	if t.CheckBaseResponse(response) {
		return true
	}

	switch response.Ret {
	case RET_FJOIN_ROOM:
		data := response.Data

		t.EnterInfo = make(map[int]float64)
		detail, ok := data.(map[string]interface{})
		if !ok {
			return true
		}
		playernums, _ := detail["playernum"].([]interface{})
		rounds, _ := detail["rounds"].([]interface{})
		antes, _ := detail["ante"].([]interface{})
		betsecs, _ := detail["betsec"].([]interface{})
		bringinlowerbounds, _ := detail["bringinlowerbound"].([]interface{})
		t.SelectPlayer.Options = []string{}
		for _, player := range playernums {
			playerNum := int(player.(float64))
			t.SelectPlayer.Options = append(t.SelectPlayer.Options, strconv.Itoa(playerNum))
		}
		t.SelectPlayer.SetSelectedIndex(0)
		t.SelectRound.Options = []string{}
		for _, round := range rounds {
			roundNum := int(round.(float64))
			t.SelectRound.Options = append(t.SelectRound.Options, strconv.Itoa(roundNum))
		}
		t.SelectRound.SetSelectedIndex(0)
		t.SelectAnte.Options = []string{}
		for _, ante := range antes {
			anteNum := int(ante.(float64))
			t.SelectAnte.Options = append(t.SelectAnte.Options, strconv.Itoa(anteNum))
		}
		t.SelectAnte.SetSelectedIndex(0)
		t.SelectBetSec.Options = []string{}
		for _, betsec := range betsecs {
			betsecNum := int(betsec.(float64))
			t.SelectBetSec.Options = append(t.SelectBetSec.Options, strconv.Itoa(betsecNum))
		}
		t.SelectBetSec.SetSelectedIndex(0)
		t.SelectBringInLowerBound.Options = []string{}
		for _, bringinlowerbound := range bringinlowerbounds {
			bringinlowerboundNum := int(bringinlowerbound.(float64))
			t.SelectBringInLowerBound.Options = append(t.SelectBringInLowerBound.Options, strconv.Itoa(bringinlowerboundNum))
		}
		t.SelectBringInLowerBound.SetSelectedIndex(0)
		return true
	case ACT_TABLE_STATUS:
		t.AddTableStatus(t.GetTableStatus(response))
		return false
	case RET_FCREATE_GAME, RET_FJOIN_GAME:
		t.SetSetting(response)
		return true
	case RET_FQUIT_GAME:
		if response.Code == constant.ERROR_CODE_SUCCESS || response.Code == constant.ERROR_CODE_ERROR_NO_NEED_TO_QUITGAME {
			t.EnableSetting(response)
		}
		return false
	}

	return false
}

func (t *BaseFriendsClient) SendFriendsReady(ready bool) (bool, error) {
	var data struct {
		FReady struct {
			Ready bool `json:"Ready"`
		}
	}

	data.FReady.Ready = ready
	return t.SendMessage(data)
}

func (t *BaseFriendsClient) SendStartGame() (bool, error) {
	var data struct {
		FStartGame struct {
		}
	}
	return t.SendMessage(data)
}

func (t *BaseFriendsClient) SendRoundStop() (bool, error) {
	var data struct {
		FRoundStop struct {
		}
	}
	return t.SendMessage(data)
}

func (t *BaseFriendsClient) SetSetting(response *utils.RespBase) {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}
	if data["Setting"] != nil {
		setting := data["Setting"].(map[string]interface{})
		if setting["InviteCode"] != nil {
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
		}
	}
}

func (t *BaseFriendsClient) EnableSetting(response *utils.RespBase) {
	t.EntryInviteCode.Enable()
	t.SelectPlayer.Enable()
	t.SelectRound.Enable()
	t.SelectAnte.Enable()
	t.SelectBetSec.Enable()
	t.SelectBringInLowerBound.Enable()
}
