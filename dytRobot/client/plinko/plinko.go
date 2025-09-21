package plinko

import (
	"dytRobot/client"
	"dytRobot/utils"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type BallInfo struct {
	BallID int `json:"BallID"`
	HoleID int `json:"HoleID"`
	BetID  int `json:"BetID"`
	SID    int `json:"SID"`
}

type PlinkoClient struct {
	*client.BaseElecClient
	buttonBet         *widget.Button
	buttonResult      *widget.Button
	buttonDebug       *widget.Button
	BetID             *widget.Label
	chipsSelect       *widget.Select
	debugBallIDSelect *widget.Select
	debugHoleIDSelect *widget.Select
	debugSwitch       *widget.Check

	tableBalls map[int]BallInfo
}

const (
	ACTION_BET        = 1
	ACTION_GET_RESULT = 2
)

func NewClient(setting client.ClientConfig) *PlinkoClient {
	elecClient := client.NewElecClient(setting)
	t := &PlinkoClient{
		BaseElecClient: elecClient,
	}
	t.CheckResponse = t.CheckPlinkoResponse
	// t.CustomMessage = append(t.CustomMessage, "{\"DebugPayout\":%d")
	t.EntrySendMessage.SetOptions(t.CustomMessage)

	return t
}

func (t *PlinkoClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateElecSection(c)
	t.CreatePlinkoSection(c)
	t.CreateTableBalls(c)
	t.CreateBottomSection(c)
	return c
}

func (t *PlinkoClient) CreateTableBalls(c *fyne.Container) {
	t.tableBalls = make(map[int]BallInfo)
}

func (t *PlinkoClient) CheckPlinkoResponse(response *utils.RespBase) bool {
	if t.CheckBaseResponse(response) {
		return true
	}

	switch response.Ret {
	// case client.ACT_GAME_PERIOD:
	// 	t.SetButton(t.buttonBet, false)
	// 	switch t.Fsm {
	// 	case "Bet":
	// 		t.SetButton(t.buttonBet, true)
	// 	}
	// 	return true
	case client.RET_PLAYER_ACTION:
		t.AddTableStatus(t.GetActionInfo(response))
		return true
	}

	return t.CheckElecResponse(response)
}

func (t *PlinkoClient) SetButton(button *widget.Button, enable bool) {
	if enable {
		button.Enable()
	} else {
		button.Disable()
	}
}

func (t *PlinkoClient) GetActionInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	action := int(data["Action"].(float64))
	var message string
	switch action {
	case ACTION_BET:
		if b, ok := data["BallInfo"].(map[string]interface{}); ok {
			sId := b["SID"].(float64)
			ballID := b["BallID"].(float64)
			bet := b["BetID"].(float64)
			holeID := b["HoleID"].(float64)

			ballInfo := BallInfo{
				SID:    int(sId),
				BallID: int(ballID),
				HoleID: int(holeID),
				BetID:  int(bet),
			}
			t.tableBalls[ballInfo.SID] = ballInfo
			message = fmt.Sprintf("球落下 SID : %d", ballInfo.SID)
		}
	case ACTION_GET_RESULT:
		sId := int(data["SID"].(float64))
		message = fmt.Sprintf("球兌獎 SID : %d", sId)
		delete(t.tableBalls, sId)
	}
	message += "\n"
	return message
}

func (t *PlinkoClient) CreatePlinkoSection(c *fyne.Container) {
	t.BetID = widget.NewLabel("下注選項")
	t.chipsSelect = widget.NewSelect([]string{"0", "1", "2", "3", "4"}, nil)
	t.buttonBet = widget.NewButton("下注", func() {
		betStr := t.chipsSelect.Selected
		betID, _ := strconv.Atoi(betStr)

		t.SendGameBet(betID)
	})
	t.buttonResult = widget.NewButton("取得得獎", func() {
		t.SendGetResult()
	})
	t.debugBallIDSelect = widget.NewSelect([]string{"0", "1", "2"}, nil)
	t.debugHoleIDSelect = widget.NewSelect([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17"}, nil)
	t.debugSwitch = widget.NewCheck("Debug", nil)
	t.debugSwitch.Checked = false
	t.buttonDebug = widget.NewButton("Debug", func() {
		t.SendDebugInfo()
	})

	section1 := container.NewHBox(t.BetID, t.chipsSelect)
	section2 := container.NewHBox(t.buttonBet, t.buttonResult)
	section3 := container.NewHBox(t.debugBallIDSelect, t.debugHoleIDSelect, t.debugSwitch, t.buttonDebug)
	section := container.NewVBox(section1, section2, section3)
	c.Add(section)
}

func (t *PlinkoClient) SendDebugInfo() {
	ballStr := t.debugBallIDSelect.Selected
	ballID, _ := strconv.Atoi(ballStr)
	holeStr := t.debugHoleIDSelect.Selected
	holeID, _ := strconv.Atoi(holeStr)
	var data struct {
		DebugInfo struct {
			Data struct {
				BallID      int  `json:"BallID"`
				HoleID      int  `json:"HoleID"`
				DebugSwitch bool `json:"DebugSwitch"`
			}
		}
	}
	data.DebugInfo.Data.BallID = ballID
	data.DebugInfo.Data.HoleID = holeID
	data.DebugInfo.Data.DebugSwitch = t.debugSwitch.Checked
	t.SendMessage(data)
}

func (t *PlinkoClient) SendGameBet(betID int) {
	var data struct {
		PlayerAction struct {
			Action int
			Data   struct {
				BetInfo int `json:"BetInfo"`
			}
		}
	}

	data.PlayerAction.Action = ACTION_BET
	data.PlayerAction.Data.BetInfo = betID

	t.SendMessage(data)
}

func (t *PlinkoClient) SendGetResult() {
	for sId := range t.tableBalls {
		var data struct {
			PlayerAction struct {
				Action int
				Data   struct {
					SID int `json:"SID"`
				}
			}
		}

		data.PlayerAction.Action = ACTION_GET_RESULT
		data.PlayerAction.Data.SID = sId

		t.SendMessage(data)
		break
	}
}
