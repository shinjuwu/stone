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
	RET_INTO_ROOM  = "IntoRoom"
	RET_OUTTO_ROOM = "OuttoRoom"
	RET_INTO_GAME  = "IntoGame"
	RET_OUTTO_GAME = "OuttoGame"

	RET_PLAYER_ACTION = "PlayerAction"
)

type BaseElecClient struct {
	*BaseClient
}

func NewElecClient(setting ClientConfig) *BaseElecClient {
	baseClient := NewBaseClient(setting)
	t := &BaseElecClient{
		BaseClient: baseClient,
	}
	t.CheckResponse = t.CheckElecResponse

	t.CustomMessage = append(t.CustomMessage, "{\"IntoRoom\":{\"GameID\":"+strconv.Itoa(t.GameId)+"}}")
	t.CustomMessage = append(t.CustomMessage, "{\"IntoGame\":{\"TableId\":"+strconv.Itoa(t.GameId)+"0}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}
func (t *BaseElecClient) IntoRoom(gameId int) (bool, error) {
	var data struct {
		IntoRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.IntoRoom.GameID = gameId

	return t.SendMessage(data)
}
func (t *BaseElecClient) OuttoRoom(gameId int) (bool, error) {
	var data struct {
		OuttoRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.OuttoRoom.GameID = gameId

	return t.SendMessage(data)
}
func (t *BaseElecClient) IntoGame(tableId int) (bool, error) {
	var data struct {
		IntoGame struct {
			TableId   int     `json:"tableId"`
			BringGold float64 `json:"bringGold,omitempty"` // singleWallet
		}
	}
	if t.WalletType == constant.SW_TYPE_SINGLE {
		minBringin := t.EnterInfo[tableId]
		swGoldInt := int(t.SWGold)
		bringIn := math.Min(minBringin*100, float64(swGoldInt))
		data.IntoGame.BringGold = bringIn
	}
	data.IntoGame.TableId = tableId

	return t.SendMessage(data)
}
func (t *BaseElecClient) OuttoGame(tableId int) (bool, error) {
	var data struct {
		OuttoGame struct {
			TableId int `json:"tableId"`
		}
	}
	data.OuttoGame.TableId = tableId

	return t.SendMessage(data)
}

func (t *BaseElecClient) ElecStop() (bool, error) {
	var data struct {
		ElecStop struct {
		}
	}

	return t.SendMessage(data)
}

func (t *BaseElecClient) GetTableStatus(response *utils.RespBase) string {
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

func (t *BaseElecClient) CreateElecSection(c *fyne.Container) {
	//進入遊戲按鈕、房間按鈕
	buttonJoinRoom := widget.NewButton("進入大廳", func() {
		t.IntoRoom(t.GameId)
	})
	buttonQuitRoom := widget.NewButton("離開大廳", func() {
		t.OuttoRoom(t.GameId)
	})

	gameLobby := container.NewHBox(buttonJoinRoom, buttonQuitRoom)
	c.Add(gameLobby)
	c.Add(t.LabalRoom)
	c.Add(t.LabelFsm)

	gameRoom := container.NewHBox()
	for roomType := 0; roomType < 4; roomType++ {
		roomName := constant.RoomType2Name[roomType]
		buttonIntoGame := widget.NewButton(roomName, nil)
		buttonIntoGame.OnTapped = func() {
			name := buttonIntoGame.Text
			t.SetTableStatus("")
			t.TableId = t.GameId*10 + constant.GameName2Type[name]
			t.IntoGame(t.TableId)
		}
		gameRoom.Add(buttonIntoGame)
	}
	buttonOuttoGame := widget.NewButton("離開牌桌", func() {
		t.OuttoGame(t.TableId)
	})
	gameRoom.Add(buttonOuttoGame)

	c.Add(gameRoom)
	c.Add(t.EntryTableStatus)
}

func (t *BaseElecClient) CheckElecResponse(response *utils.RespBase) bool {
	if t.CheckBaseResponse(response) {
		return true
	}

	switch response.Ret {
	case RET_INTO_ROOM:
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
		return true
	case ACT_EXIT_GAME:
		fmt.Println(response)
		return true
	case ACT_JOIN_GAME:
		fmt.Println(response)
		return true
	}

	return false
}
