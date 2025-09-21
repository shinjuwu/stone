package colordisc

import (
	"dytRobot/client"
	"dytRobot/utils"

	"fyne.io/fyne/v2"
)

type ColordiscClient struct {
	*client.BaseBetClient
}

func NewClient(setting client.ClientConfig) *ColordiscClient {
	betClient := client.NewBetClient(setting)
	t := &ColordiscClient{
		BaseBetClient: betClient,
	}

	t.CheckResponse = t.CheckColordiscResponse
	return t
}

func (t *ColordiscClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateBetSection(c)
	//t.CreateBullBullSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *ColordiscClient) CheckColordiscResponse(response *utils.RespBase) bool {
	return t.CheckBetResponse(response)
}
