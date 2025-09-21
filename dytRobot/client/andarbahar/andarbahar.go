package andarbahar

import (
	"dytRobot/client"
	"dytRobot/utils"

	"fyne.io/fyne/v2"
)

type AndarbaharClient struct {
	*client.BaseBetClient
}

func NewClient(setting client.ClientConfig) *AndarbaharClient {
	betClient := client.NewBetClient(setting)
	t := &AndarbaharClient{
		BaseBetClient: betClient,
	}
	t.CheckResponse = t.CheckCockfightResponse

	return t
}

func (t *AndarbaharClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateBetSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *AndarbaharClient) CheckCockfightResponse(response *utils.RespBase) bool {
	return t.CheckBetResponse(response)
}
