package jackpot

import (
	"dytRobot/client"
	"dytRobot/constant"
	"dytRobot/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type jackpotClient struct {
	*client.BaseBetClient
}

func NewClient(setting client.ClientConfig) *jackpotClient {
	betClient := client.NewBetClient(setting)
	t := &jackpotClient{
		BaseBetClient: betClient,
	}

	t.CheckResponse = t.CheckjackpotResponse
	return t
}

func (t *jackpotClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateBetSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *jackpotClient) CheckjackpotResponse(response *utils.RespBase) bool {
	return t.CheckBetResponse(response)
}

func (t *jackpotClient) CreateBetSection(c *fyne.Container) {
	//進入遊戲按鈕、房間按鈕
	buttonJoinRoom := widget.NewButton("進入JP", func() {
		t.InJackpot(t.GameId)
	})
	buttonQuitRoom := widget.NewButton("離開JP", func() {
		t.OutJackpot(t.GameId)
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

func (t *jackpotClient) InJackpot(gameId int) (bool, error) {
	var data struct {
		InJackpot struct {
			TableID int `json:"tableId"`
		}
	}
	data.InJackpot.TableID = 900101

	return t.SendMessage(data)
}
func (t *jackpotClient) OutJackpot(gameId int) (bool, error) {
	var data struct {
		OutJackpot struct {
			TableID int `json:"tableId"`
		}
	}
	data.OutJackpot.TableID = 900101

	return t.SendMessage(data)
}
