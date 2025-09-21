//go:build !auto

package dytapp

import (
	"math/rand"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

var myApp fyne.App
var myWindow fyne.Window

type DefaultApp struct {
	version   string
	processID string
}

func NewApp(version string) *DefaultApp {
	app := new(DefaultApp)
	app.processID = "development"
	app.version = version
	return app
}

func (dytapp *DefaultApp) Run(debug bool) {
	dytapp.setAppEnv()
	setAppMenu()
	setPressureTestEnv()
	setClientTestEnv()

	myWindow.SetContent(tabsPressTest)

	myWindow.ShowAndRun()
}

func (dytapp *DefaultApp) setAppEnv() {
	os.Setenv("FYNE_FONT", "C:/Windows/Fonts/msyh.ttc")
	myApp = app.New()
	myWindow = myApp.NewWindow("大鏞遊戲測試工具 版本:" + dytapp.version)
	myWindow.Resize(fyne.NewSize(1280, 600))
	myApp.Settings().SetTheme(theme.DarkTheme())

	rand.Seed(time.Now().UnixNano())
}

func setAppMenu() {

	menuItems1 := fyne.NewMenuItem("壓測工具", func() {
		myWindow.SetContent(tabsPressTest)
	})
	menuItems2 := fyne.NewMenuItem("client測試工具", func() {
		myWindow.SetContent(clientContainer)
	})

	newMenu := fyne.NewMenu("模式", menuItems1, menuItems2)

	myWindow.SetMainMenu(fyne.NewMainMenu(newMenu))
}
