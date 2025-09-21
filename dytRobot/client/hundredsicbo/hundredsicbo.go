package hundredsicbo

import (
	"dytRobot/client"
	"dytRobot/utils"

	"fyne.io/fyne/v2"
)

type HundredsicboClient struct {
	*client.BaseBetClient
}

func NewClient(setting client.ClientConfig) *HundredsicboClient {
	betClient := client.NewBetClient(setting)
	t := &HundredsicboClient{
		BaseBetClient: betClient,
	}

	t.CheckResponse = t.CheckHundredsicResponse

	t.CustomMessage = append(t.CustomMessage, "{\"DebugDealCard\":{\"card\":[[6,6,6]]}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)

	return t
}

func (t *HundredsicboClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateBetSection(c)
	//t.CreateBullBullSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *HundredsicboClient) CheckHundredsicResponse(response *utils.RespBase) bool {
	return t.CheckBetResponse(response)
}
