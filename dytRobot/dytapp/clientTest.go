//go:build !auto

package dytapp

import (
	"dytRobot/client"
	"dytRobot/client/andarbahar"
	"dytRobot/client/baccarat"
	"dytRobot/client/blackjack"
	"dytRobot/client/bullbull"
	"dytRobot/client/catte"
	"dytRobot/client/chinesepoker"
	"dytRobot/client/cockfight"
	"dytRobot/client/colordisc"
	"dytRobot/client/dogracing"
	"dytRobot/client/fantan"
	"dytRobot/client/friendstexas"
	"dytRobot/client/fruit777slot"
	"dytRobot/client/fruitslot"
	"dytRobot/client/goldenflower"
	"dytRobot/client/hundredsicbo"
	"dytRobot/client/jackpot"
	"dytRobot/client/megsharkslot"
	"dytRobot/client/midasslot"
	"dytRobot/client/okey"
	"dytRobot/client/plinko"
	"dytRobot/client/pokdeng"
	"dytRobot/client/prawncrab"
	"dytRobot/client/rcfish"
	"dytRobot/client/rocket"
	"dytRobot/client/roulette"
	"dytRobot/client/rummy"
	"dytRobot/client/sangong"
	"dytRobot/client/teenpatti"
	"dytRobot/client/texas"
	"dytRobot/constant"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var clientContainer *fyne.Container

func setClientTestEnv() {
	var gameList []string
	var allPage []fyne.CanvasObject
	var selectedId widget.ListItemID

	for _, gameID := range constant.GameIDList {
		page := createClinetGamePage(gameID)
		if page != nil {
			gameList = append(gameList, constant.GameIDNameTable[gameID])
			allPage = append(allPage, page)
		}
	}

	list := widget.NewList(
		func() int {
			return len(gameList)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(gameList[id])
		},
	)

	list.Select(0)
	selectPage := container.NewMax(allPage[0])
	scrollPage := container.NewScroll(selectPage)
	clientContainer = container.NewBorder(nil, nil, list, nil, scrollPage)

	list.OnSelected = func(id widget.ListItemID) {
		if client.IsConnect { //讓client測試工具一次只能操作一個遊戲
			list.Select(selectedId)
			return
		}
		selectPage.Objects[0] = container.NewMax(allPage[id])
		selectedId = id
		clientContainer.Refresh()
	}

}

func createClinetGamePage(gameID int) fyne.CanvasObject {
	var page fyne.CanvasObject
	config := client.ClientConfig{
		GameId: gameID,
	}
	switch gameID {
	case 1001:
		object := baccarat.NewClient(config)
		page = object.CreateSection()
	case 1002:
		object := fantan.NewClient(config)
		page = object.CreateSection()
	case 1003:
		object := colordisc.NewClient(config)
		page = object.CreateSection()
	case 1004:
		object := prawncrab.NewClient(config)
		page = object.CreateSection()
	case 1005:
		object := hundredsicbo.NewClient(config)
		page = object.CreateSection()
	case 1006:
		object := cockfight.NewClient(config)
		page = object.CreateSection()
	case 1007:
		object := dogracing.NewClient(config)
		page = object.CreateSection()
	case 1008:
		object := rocket.NewClient(config)
		page = object.CreateSection()
	case 1009:
		object := andarbahar.NewClient(config)
		page = object.CreateSection()
	case 1010:
		object := roulette.NewClient(config)
		page = object.CreateSection()
	case 2001:
		object := blackjack.NewClient(config)
		page = object.CreateSection()
	case 2002:
		object := sangong.NewClient(config)
		page = object.CreateSection()
	case 2003:
		object := bullbull.NewClient(config)
		page = object.CreateSection()
	case 2004:
		object := texas.NewClient(config)
		page = object.CreateSection()
	case 2005:
		object := rummy.NewClient(config)
		page = object.CreateSection()
	case 2006:
		object := goldenflower.NewClient(config)
		page = object.CreateSection()
	case 2007:
		object := pokdeng.NewClient(config)
		page = object.CreateSection()
	case 2008:
		object := catte.NewClient(config)
		page = object.CreateSection()
	case 2009:
		object := chinesepoker.NewClient(config)
		page = object.CreateSection()
	case 2010:
		object := okey.NewClient(config)
		page = object.CreateSection()
	case 2011:
		object := teenpatti.NewClient(config)
		page = object.CreateSection()
	case 3001:
		object := fruitslot.NewClient(config)
		page = object.CreateSection()
	case 3002:
		object := rcfish.NewClient(config)
		page = object.CreateSection()
	case 3003:
		object := plinko.NewEnhancedClient(config)
		page = object.CreateSection()
	case 4001:
		object := fruit777slot.NewClient(config)
		page = object.CreateSection()
	case 4002:
		object := megsharkslot.NewClient(config)
		page = object.CreateSection()
	case 4003:
		object := midasslot.NewClient(config)
		page = object.CreateSection()
	case 5001:
		object := friendstexas.NewClient(config)
		page = object.CreateSection()
	case 9001:
		object := jackpot.NewClient(config)
		page = object.CreateSection()
	}

	return page
}
