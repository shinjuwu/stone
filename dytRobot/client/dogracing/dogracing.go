package dogracing

import (
	"dytRobot/client"
	"dytRobot/utils"

	"fyne.io/fyne/v2"
)

type DogracingClient struct {
	*client.BaseBetClient
}

func NewClient(setting client.ClientConfig) *DogracingClient {
	betClient := client.NewBetClient(setting)
	t := &DogracingClient{
		BaseBetClient: betClient,
	}
	t.CheckResponse = t.CheckCockfightResponse

	return t
}

func (t *DogracingClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateBetSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *DogracingClient) CheckCockfightResponse(response *utils.RespBase) bool {
	return t.CheckBetResponse(response)
}
