package roulette

import (
	"dytRobot/client"
	"dytRobot/utils"

	"fyne.io/fyne/v2"
)

type RouletteClient struct {
	*client.BaseBetClient
}

func NewClient(setting client.ClientConfig) *RouletteClient {
	betClient := client.NewBetClient(setting)
	t := &RouletteClient{
		BaseBetClient: betClient,
	}
	t.CheckResponse = t.CheckCockfightResponse

	return t
}

func (t *RouletteClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateBetSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *RouletteClient) CheckCockfightResponse(response *utils.RespBase) bool {
	return t.CheckBetResponse(response)
}
