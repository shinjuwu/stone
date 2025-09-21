package megsharkslot

import (
	"dytRobot/client"
	"dytRobot/utils"
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	BUTTON_COUNT = 8 // 下注區域數量

	ACTION_SPIN = 1
)

type MegSharkSlotClient struct {
	*client.BaseSlotClient
	BetID       *widget.Label
	chipsSelect *widget.Select
	buttonBet   *widget.Button
	Bet         float64
	debugSwitch *widget.Check
}

func NewClient(setting client.ClientConfig) *MegSharkSlotClient {
	slotClient := client.NewSlotClient(setting)
	t := &MegSharkSlotClient{
		BaseSlotClient: slotClient,
	}

	t.CheckResponse = t.CheckMegSharkSlotResponse

	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}

func (t *MegSharkSlotClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateSlotSection(c)
	t.CreateMegSharkSlotSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *MegSharkSlotClient) CheckMegSharkSlotResponse(response *utils.RespBase) bool {
	if t.CheckSlotResponse(response) {
		return true
	}
	switch response.Ret {
	case client.RET_GOIN_GAME:
		t.AddTableStatus(t.GetJoinInfo(response))
		return true
	case client.ACT_GOLD:
		t.AddTableStatus(t.GetGoldInfo(response))
		return true
	case client.RET_SLOT_ACTION:
		t.AddTableStatus(t.GetActionInfo(response))
		return true
	}
	return false
}

func (t *MegSharkSlotClient) CreateMegSharkSlotSection(c *fyne.Container) {
	t.BetID = widget.NewLabel("下注選項")
	t.chipsSelect = widget.NewSelect([]string{"0", "1", "2", "3", "4"}, nil)
	t.buttonBet = widget.NewButton("SPIN", func() {
		betStr := t.chipsSelect.Selected
		betID, _ := strconv.Atoi(betStr)

		t.SendGameSpin(betID)
	})

	t.debugSwitch = widget.NewCheck("Debug", nil)
	t.debugSwitch.Checked = false
	buttonDebug := widget.NewButton("Debug", func() {
		t.SendDebugInfo()
	})
	buttonDebugNG := widget.NewButton("Debug NG", func() {
		t.SendDebugNGInfo()
	})

	section1 := container.NewHBox(t.BetID, t.chipsSelect, t.buttonBet, t.debugSwitch, buttonDebug, buttonDebugNG)
	section := container.NewVBox(section1)
	c.Add(section)
}

func (t *MegSharkSlotClient) SendGameSpin(betID int) {
	var data struct {
		SlotAction struct {
			Action int         `json:"Action"`
			Data   interface{} `json:"Data,omitempty"`
		}
	}

	var info struct {
		BetPos int `json:"BetPos"`
	}
	info.BetPos = betID

	data.SlotAction.Action = ACTION_SPIN
	data.SlotAction.Data = info

	t.SendMessage(data)
}

func (t *MegSharkSlotClient) GetJoinInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}
	fsm := data["Fsm"].(string)
	labelMessage := "Fsm:" + fsm

	if gid, ok := data["Gid"].(string); ok {
		labelMessage += " Gid:" + gid
	}
	t.LabelFsm.SetText(labelMessage)

	// baseBet := int(data["BaseBet"].(float64))
	// t.labelBaseBet.SetText("底注:" + strconv.Itoa(baseBet))

	var message string

	return message
}

func (t *MegSharkSlotClient) GetGoldInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}

	gold := data["Gold"].(float64)
	return fmt.Sprintf("玩家金額:%.4f\n", gold)
}

func (t *MegSharkSlotClient) GetActionInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}

	/* fsm := data["Fsm"].(string)
	labelMessage := "Fsm:" + fsm */
	gid := data["Gid"].(string)
	labelMessage := " Gid:" + gid
	t.LabelFsm.SetText(labelMessage)

	t.SetTableStatusClear()

	var message string
	return message
}

func (t *MegSharkSlotClient) SendDebugInfo() {
	var data struct {
		DebugInfo struct {
			Data struct {
				DebugIndex  [][]int `json:"DebugIndex"`
				DebugSwitch bool    `json:"DebugSwitch"`
				DebugRTP    string  `json:"DebugRTP"`
			}
		}
	}
	data.DebugInfo.Data.DebugIndex = [][]int{{41, 1, 13, 1, 10}, {13, 1, 15, 1, 18}}
	data.DebugInfo.Data.DebugSwitch = t.debugSwitch.Checked
	data.DebugInfo.Data.DebugRTP = "98"
	t.SendMessage(data)
}

func (t *MegSharkSlotClient) SendDebugNGInfo() {
	var data struct {
		DebugInfo struct {
			Data struct {
				DebugNGIndex [][]int `json:"DebugNGIndex"`
				DebugSwitch  bool    `json:"DebugSwitch"`
				DebugRTP     string  `json:"DebugRTP"`
			}
		}
	}
	data.DebugInfo.Data.DebugNGIndex = [][]int{{5, 5, 5, 5, 5}, {5, 5, 5, 5, 6}}
	data.DebugInfo.Data.DebugSwitch = t.debugSwitch.Checked
	data.DebugInfo.Data.DebugRTP = "98"
	t.SendMessage(data)
}
