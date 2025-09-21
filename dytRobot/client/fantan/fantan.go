package fantan

import (
	"dytRobot/client"
	"dytRobot/utils"

	"fyne.io/fyne/v2"
)

type FantanClient struct {
	*client.BaseBetClient
}

func NewClient(setting client.ClientConfig) *FantanClient {
	betClient := client.NewBetClient(setting)
	t := &FantanClient{
		BaseBetClient: betClient,
	}

	t.CheckResponse = t.CheckFantanResponse

	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[20,22,44]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}

func (t *FantanClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateBetSection(c)
	//t.CreateBullBullSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *FantanClient) CheckFantanResponse(response *utils.RespBase) bool {
	return t.CheckBetResponse(response)
}
