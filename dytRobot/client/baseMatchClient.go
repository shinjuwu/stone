package client

import (
	"dytRobot/constant"
	"dytRobot/utils"
	"fmt"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	RET_JOIN_ROOM  = "JoinRoom"
	RET_QUIT_ROOM  = "QuitRoom"
	RET_JOIN_GAME  = "JoinGame"
	RET_QUIT_GAME  = "QuitGame"
	RET_MATCH_STOP = "MatchStop"

	ACT_TABLE_STATUS = "ActTableStatus"
	ACT_ROOM_STATUS  = "ActRoomStatus"

	//魚機
	ACT_EXIT_GAME = "ActExitGame"
	ACT_JOIN_GAME = "ActJoinGame"
)

type BaseMatchClient struct {
	*BaseClient
}

func NewMatchClient(setting ClientConfig) *BaseMatchClient {
	baseClient := NewBaseClient(setting)
	t := &BaseMatchClient{
		BaseClient: baseClient,
	}
	t.CheckResponse = t.CheckMatchResponse

	t.CustomMessage = append(t.CustomMessage, "{\"JoinRoom\":{\"GameID\":"+strconv.Itoa(t.GameId)+"}}")
	t.CustomMessage = append(t.CustomMessage, "{\"JoinGame\":{\"TableId\":"+strconv.Itoa(t.GameId)+"0}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}
func (t *BaseMatchClient) JoinRoom(gameId int) (bool, error) {
	var data struct {
		JoinRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.JoinRoom.GameID = gameId

	return t.SendMessage(data)
}
func (t *BaseMatchClient) QuitRoom(gameId int) (bool, error) {
	var data struct {
		QuitRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.QuitRoom.GameID = gameId

	return t.SendMessage(data)
}
func (t *BaseMatchClient) JoinGame(tableId int) (bool, error) {
	var data struct {
		JoinGame struct {
			TableId   int     `json:"tableId"`
			BringGold float64 `json:"bringGold,omitempty"` //Texas & singleWallet
		}
	}
	data.JoinGame.TableId = tableId
	if t.WalletType == constant.SW_TYPE_SINGLE {
		minBringin := t.EnterInfo[tableId]
		swGoldInt := int(t.SWGold)
		bringIn := math.Min(minBringin*100, float64(swGoldInt))
		data.JoinGame.BringGold = bringIn
	} else if (t.TableId / 10) == 2004 {
		data.JoinGame.BringGold = constant.TexasBringGold[(t.TableId % 10)]
	}

	return t.SendMessage(data)
}
func (t *BaseMatchClient) QuitGame(tableId int) (bool, error) {
	var data struct {
		QuitGame struct {
			TableId int `json:"tableId"`
		}
	}
	data.QuitGame.TableId = tableId

	return t.SendMessage(data)
}

func (t *BaseMatchClient) MatchStop() (bool, error) {
	var data struct {
		MatchStop struct {
		}
	}

	return t.SendMessage(data)
}

func (t *BaseMatchClient) GetTableStatus(response *utils.RespBase) string {
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

func (t *BaseMatchClient) CreateMatchSection(c *fyne.Container) {
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
	for roomType := 0; roomType < constant.RoomTypeNum[t.GameId]; roomType++ {
		roomName := constant.RoomType2Name[roomType]
		buttonJoinGame := widget.NewButton(roomName, nil)
		buttonJoinGame.OnTapped = func() {
			name := buttonJoinGame.Text
			t.SetTableStatus("")
			t.TableId = t.GameId*10 + constant.GameName2Type[name]
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

	c.Add(gameRoom)
	c.Add(t.EntryTableStatus)
}

func (t *BaseMatchClient) CheckMatchResponse(response *utils.RespBase) bool {
	if t.CheckBaseResponse(response) {
		return true
	}

	switch response.Ret {
	case RET_JOIN_ROOM:
		var roomInfo string
		data, ok := response.Data.([]interface{})
		if !ok {
			return true
		}
		t.EnterInfo = make(map[int]float64)
		for i, info := range data {
			detail, ok := info.(map[string]interface{})
			if !ok {
				return true
			}
			tableId := int(detail["TableId"].(float64))
			_, tableName := constant.GetGameRoomName(t.GameId, tableId)
			tableStatus := t.CheckTableStatus(int(detail["Status"].(float64)))
			tableAnte := int(detail["Ante"].(float64))
			tableEnter := detail["EnterGold"].(float64)
			t.EnterInfo[tableId] = tableEnter

			roomInfo += fmt.Sprintf("%s  %s  底注:%2d  准入：%4d", tableName, tableStatus, tableAnte, int(tableEnter))
			if i != len(data)-1 {
				roomInfo += "\n"
			}
		}
		t.LabalRoom.SetText(roomInfo)
		return true
	case ACT_TABLE_STATUS:
		t.AddTableStatus(t.GetTableStatus(response))
		return false
	}

	return false
}
