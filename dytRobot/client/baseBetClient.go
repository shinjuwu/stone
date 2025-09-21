package client

import (
	"dytRobot/constant"
	"dytRobot/utils"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	RET_ENTER_ROOM = "EnterRoom"
	RET_EXIT_ROOM  = "ExitRoom"
	RET_ENTER_GAME = "EnterGame"
	RET_EXIT_GAME  = "ExitGame"

	//ACT_TABLE_STATUS = "ActTableStatus"
)

type BaseBetClient struct {
	*BaseClient
}

func NewBetClient(setting ClientConfig) *BaseBetClient {
	baseClient := NewBaseClient(setting)
	t := &BaseBetClient{
		BaseClient: baseClient,
	}
	t.CheckResponse = t.CheckBetResponse

	t.CustomMessage = append(t.CustomMessage, "{\"EnterRoom\":{\"GameID\":"+strconv.Itoa(t.GameId)+"}}")
	t.CustomMessage = append(t.CustomMessage, "{\"EnterGame\":{\"TableId\":"+strconv.Itoa(t.GameId)+"0}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)

	return t
}
func (t *BaseBetClient) EnterRoom(gameId int) (bool, error) {
	var data struct {
		EnterRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.EnterRoom.GameID = gameId

	return t.SendMessage(data)
}
func (t *BaseBetClient) ExitRoom(gameId int) (bool, error) {
	var data struct {
		ExitRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.ExitRoom.GameID = gameId

	return t.SendMessage(data)
}
func (t *BaseBetClient) EnterGame(tableId int) (bool, error) {
	var data struct {
		EnterGame struct {
			TableId int `json:"tableId"`
		}
	}
	data.EnterGame.TableId = tableId

	return t.SendMessage(data)
}
func (t *BaseBetClient) ExitGame(tableId int) (bool, error) {
	var data struct {
		ExitGame struct {
			TableId int `json:"tableId"`
		}
	}
	data.ExitGame.TableId = tableId

	return t.SendMessage(data)
}

func (t *BaseBetClient) CreateBetSection(c *fyne.Container) {
	//進入遊戲按鈕、房間按鈕
	buttonJoinRoom := widget.NewButton("進入大廳", func() {
		t.EnterRoom(t.GameId)
	})
	buttonQuitRoom := widget.NewButton("離開大廳", func() {
		t.ExitRoom(t.GameId)
	})

	gameLobby := container.NewHBox(buttonJoinRoom, buttonQuitRoom)
	c.Add(gameLobby)
	c.Add(t.LabalRoom)
	c.Add(t.LabelFsm)

	gameRoom := container.NewHBox()
	roomNum := constant.RoomTypeNum[t.GameId]
	for roomType := 0; roomType < roomNum; roomType++ {
		roomName := constant.RoomType2Name[roomType]
		buttonEnterGame := widget.NewButton(roomName, nil)
		buttonEnterGame.OnTapped = func() {
			name := buttonEnterGame.Text
			t.SetTableStatus("")
			t.TableId = t.GameId*100 + constant.GameName2Type[name]*10 + 1
			t.EnterGame(t.TableId)
		}
		gameRoom.Add(buttonEnterGame)
	}
	buttonExitGame := widget.NewButton("離開房間", func() {
		t.ExitGame(t.TableId)
	})
	gameRoom.Add(buttonExitGame)
	c.Add(gameRoom)
	c.Add(t.EntryTableStatus)
}

func (t *BaseBetClient) CheckBetResponse(response *utils.RespBase) bool {
	if t.CheckBaseResponse(response) {
		return true
	}

	switch response.Ret {
	case RET_ENTER_ROOM:
		var roomInfo string
		data, ok := response.Data.([]interface{})
		if !ok {
			return true
		}
		for _, info := range data {
			detail, ok := info.(map[string]interface{})
			if !ok {
				return true
			}
			roomType := int(detail["TableId"].(float64)) % t.GameId
			tableName := constant.RoomType2Name[roomType]
			tableStatus := t.CheckTableStatus(int(detail["Status"].(float64)))
			roomInfo += fmt.Sprintf("%s %s；", tableName, tableStatus)
		}
		t.LabalRoom.SetText(roomInfo)
		return true
	}

	return false
}
