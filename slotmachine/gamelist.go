package main

import (
	"slotEditor/editor"
	WildGem "slotEditor/editor/slot_4004WildGem"
	Jumphigh "slotEditor/editor/slot_4005Jumphigh"
	PryTreasure "slotEditor/editor/slot_4006PryTreasure"
	MegShark "slotEditor/editor/slot_MegShark"
	Midas "slotEditor/editor/slot_Midas"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var gameIDList []int = []int{4002, 4003, 4004, 4005, 4006}

var GameIDNameTable map[int]string = map[int]string{
	4002: "巨齒鯊",
	4003: "邁達斯之手",
	4004: "狂野寶石",
	4005: "跳高高",
	4006: "金字塔寶藏",
}

var NameTableGameID map[string]int = map[string]int{
	"巨齒鯊":   4002,
	"邁達斯之手": 4003,
	"狂野寶石":  4004,
	"跳高高":   4005,
	"金字塔寶藏": 4006,
}

var gamesContainer *fyne.Container

func setAllGamesBtn() {
	gamesContainer = container.NewVBox()
	count := 0
	c := container.NewHBox()
	for i := 0; i < len(gameIDList); i++ {
		count++
		gameID := gameIDList[i]
		c.Add(widget.NewButton(GameIDNameTable[gameID], func() {
			window := fyne.CurrentApp().NewWindow(GameIDNameTable[gameID])
			window.SetContent(createGamePage(gameID, &window))
			window.Show()
		}))
		if count > 5 {
			count -= 5
			gamesContainer.Add(c)
			c = container.NewHBox()
		} else if i == len(gameIDList)-1 {
			gamesContainer.Add(c)
		}
	}

}

func createGamePage(gameID int, curWindow *fyne.Window) fyne.CanvasObject {
	var page fyne.CanvasObject
	config := editor.BaseConfig{
		GameId: gameID,
		Window: curWindow,
		Rtp:    "97", //default:"97",可選"98","97","92"
	}

	switch gameID {
	case 4002:
		object := MegShark.NewGame(config)
		page = object.CreateAllSection()
	case 4003:
		object := Midas.NewGame(config)
		page = object.CreateAllSection()
	case 4004:
		object := WildGem.NewGame(config)
		page = object.CreateAllSection()
	case 4005:
		object := Jumphigh.NewGame(config)
		page = object.CreateAllSection()
	case 4006:
		object := PryTreasure.NewGame(config)
		page = object.CreateAllSection()
	}
	return page
}
