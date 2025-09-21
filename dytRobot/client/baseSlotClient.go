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
	RET_GOIN_ROOM  = "GoInRoom"
	RET_LEAVE_ROOM = "LeaveRoom"
	RET_GOIN_GAME  = "GoInGame"
	RET_LEAVE_GAME = "LeaveGame"

	RET_SLOT_ACTION = "SlotAction"
)

type BaseSlotClient struct {
	*BaseClient
}

func NewSlotClient(setting ClientConfig) *BaseSlotClient {
	baseClient := NewBaseClient(setting)
	t := &BaseSlotClient{
		BaseClient: baseClient,
	}
	t.CheckResponse = t.CheckSlotResponse

	t.CustomMessage = append(t.CustomMessage, "{\"GoInRoom\":{\"GameID\":"+strconv.Itoa(t.GameId)+"}}")
	t.CustomMessage = append(t.CustomMessage, "{\"GoInGame\":{\"TableId\":"+strconv.Itoa(t.GameId)+"0}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}
func (t *BaseSlotClient) GoInRoom(gameId int) (bool, error) {
	var data struct {
		GoInRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.GoInRoom.GameID = gameId

	return t.SendMessage(data)
}
func (t *BaseSlotClient) LeaveRoom(gameId int) (bool, error) {
	var data struct {
		LeaveRoom struct {
			GameID int `json:"GameID"`
		}
	}
	data.LeaveRoom.GameID = gameId

	return t.SendMessage(data)
}
func (t *BaseSlotClient) GoInGame(tableId int) (bool, error) {
	var data struct {
		GoInGame struct {
			TableId int `json:"tableId"`
		}
	}
	data.GoInGame.TableId = tableId

	return t.SendMessage(data)
}
func (t *BaseSlotClient) LeaveGame(tableId int) (bool, error) {
	var data struct {
		LeaveGame struct {
			TableId int `json:"tableId"`
		}
	}
	data.LeaveGame.TableId = tableId

	return t.SendMessage(data)
}

func (t *BaseSlotClient) SlotStop() (bool, error) {
	var data struct {
		SlotStop struct {
		}
	}

	return t.SendMessage(data)
}

func (t *BaseSlotClient) GetTableStatus(response *utils.RespBase) string {
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

func (t *BaseSlotClient) CreateSlotSection(c *fyne.Container) {
	//進入遊戲按鈕、房間按鈕
	buttonJoinRoom := widget.NewButton("進入大廳", func() {
		t.GoInRoom(t.GameId)
	})
	buttonQuitRoom := widget.NewButton("離開大廳", func() {
		t.LeaveRoom(t.GameId)
	})

	gameLobby := container.NewHBox(buttonJoinRoom, buttonQuitRoom)
	c.Add(gameLobby)
	c.Add(t.LabalRoom)
	c.Add(t.LabelFsm)

	gameRoom := container.NewHBox()
	for roomType := 0; roomType < constant.RoomTypeNum[t.GameId]; roomType++ {
		roomName := constant.RoomType2Name[roomType]
		buttonGoInGame := widget.NewButton(roomName, nil)
		buttonGoInGame.OnTapped = func() {
			name := buttonGoInGame.Text
			t.SetTableStatus("")
			t.TableId = t.GameId*10 + constant.GameName2Type[name]
			t.GoInGame(t.TableId)
		}
		gameRoom.Add(buttonGoInGame)
	}
	buttonLeaveGame := widget.NewButton("離開牌桌", func() {
		t.LeaveGame(t.TableId)
	})
	gameRoom.Add(buttonLeaveGame)

	c.Add(gameRoom)
	c.Add(t.EntryTableStatus)
}

func (t *BaseSlotClient) CheckSlotResponse(response *utils.RespBase) bool {
	if t.CheckBaseResponse(response) {
		return true
	}

	switch response.Ret {
	case RET_GOIN_ROOM:
		var roomInfo string
		data, ok := response.Data.([]interface{})
		if !ok {
			return true
		}
		for i, info := range data {
			detail, ok := info.(map[string]interface{})
			if !ok {
				return true
			}
			tableId := int(detail["TableId"].(float64))
			_, tableName := constant.GetGameRoomName(t.GameId, tableId)
			tableStatus := t.CheckTableStatus(int(detail["Status"].(float64)))
			roomInfo += fmt.Sprintf("%s  %s", tableName, tableStatus)
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
