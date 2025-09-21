package prawncrab

import (
	"dytRobot/client"
	"dytRobot/utils"

	"fyne.io/fyne/v2"
)

type PrawncrabClient struct {
	*client.BaseBetClient
}

func NewClient(setting client.ClientConfig) *PrawncrabClient {
	matchClient := client.NewBetClient(setting)
	t := &PrawncrabClient{
		BaseBetClient: matchClient,
	}

	t.CheckResponse = t.CheckPrawncrabResponse

	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[6,6,6]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}

func (t *PrawncrabClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateBetSection(c)
	//t.CreateBullBullSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *PrawncrabClient) CheckPrawncrabResponse(response *utils.RespBase) bool {
	return t.CheckBetResponse(response)
}
