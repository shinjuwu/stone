package main

import (
	"dytRobot/dytapp"
)

var version string = "2.10"

func main() {

	clientApp := dytapp.NewApp(version)
	clientApp.Run(true)
}
