package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

var myApp fyne.App
var myWindow fyne.Window

var version string = "3.0.5"

func main() {

	setAppEnv()
	setAllGamesBtn()

	myWindow.SetContent(gamesContainer)

	myWindow.ShowAndRun()
}

func setAppEnv() {
	os.Setenv("FYNE_FONT", "C:/Windows/Fonts/msyh.ttc")
	myApp = app.New()
	myWindow = myApp.NewWindow("老虎機編輯器 版本:" + version)
	myWindow.Resize(fyne.NewSize(1280, 600))
	myApp.Settings().SetTheme(theme.DarkTheme())
}

// func run() {
// 	payTable := slotgame.createPayTable()
// 	slotgame.printSymbolTable(payTable)
// 	win, slotType, winLines := slotgame.calcWin(payTable)
// 	fmt.Printf("win : %d, slotType : %d, winLines : %v\n", win, slotType, winLines)
// }
