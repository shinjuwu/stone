package baccarat

import (
	"dytRobot/client"
	"dytRobot/constant"
	"dytRobot/utils"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type BaccaratClient struct {
	*client.BaseBetClient
}

func NewClient(setting client.ClientConfig) *BaccaratClient {
	betClient := client.NewBetClient(setting)
	t := &BaccaratClient{
		BaseBetClient: betClient,
	}

	t.CheckResponse = t.CheckBaccaratResponse
	return t
}

func (t *BaccaratClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateBaccaratSection(c)
	//t.CreateBullBullSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *BaccaratClient) CreateBaccaratSection(c *fyne.Container) {
	//進入遊戲按鈕、房間按鈕
	buttonEnterRoom := widget.NewButton("進入大廳", func() {
		t.EnterRoom(t.GameId)
	})
	buttonExitRoom := widget.NewButton("離開大廳", func() {
		t.ExitRoom(t.GameId)
	})

	gameLobby := container.NewHBox(buttonEnterRoom, buttonExitRoom)
	c.Add(gameLobby)
	c.Add(t.LabalRoom)
	c.Add(t.LabelFsm)

	gameRoom := container.NewVBox()
	room := container.NewHBox()
	for i := 1; i < 5; i++ {
		for roomType := 0; roomType < 4; roomType++ {
			roomName := constant.RoomType2Name[roomType] + strconv.Itoa(i)
			buttonEnterGame := widget.NewButton(roomName, nil)
			buttonEnterGame.OnTapped = func() {
				name := buttonEnterGame.Text[:9]
				num := buttonEnterGame.Text[9:]
				count, _ := strconv.Atoi(num)
				t.SetTableStatus("")
				t.TableId = t.GameId*100 + constant.GameName2Type[name]*10 + count
				t.EnterGame(t.TableId)
			}
			room.Add(buttonEnterGame)
		}
		if i == 2 {
			gameRoom.Add(room)
			room = container.NewHBox()
		} else if i == 4 {
			buttonExitGame := widget.NewButton("離開房間", func() {
				t.ExitGame(t.TableId)
			})
			room.Add(buttonExitGame)
			gameRoom.Add(room)
		}
	}
	c.Add(gameRoom)
	c.Add(t.EntryTableStatus)
}

func (t *BaccaratClient) CheckBaccaratResponse(response *utils.RespBase) bool {
	return t.CheckBetResponse(response)
}
