package client

import (
	"dytRobot/constant"
	"dytRobot/utils"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gorilla/websocket"
)

const (
	RET_LOGIN = "Login"

	ACT_BULLETIN    = "ActBulletin"
	ACT_GAME_PERIOD = "ActGamePeriod"
	ACT_GOLD        = "ActGold"
	ACT_SW_GOLD     = "ActSingleWalletGold"
)

var IsConnect bool

type ClientConfig struct {
	GameId int
}

type BaseClient struct {
	Env           int
	AgentEnv      int
	GameId        int
	TableId       int
	Connect       *websocket.Conn
	LoginName     string
	isLogIn       bool
	CustomMessage []string
	Fsm           string
	WalletType    int
	Gold          float64
	SWGold        float64
	EnterInfo     map[int]float64

	CheckResponse func(response *utils.RespBase) bool

	CheckActBulletin *widget.Check       //是否要秀跑燈訊息
	ClientMessage    *widget.Entry       //顯示收到的封包
	EntryClientName  *widget.Entry       //登入帳號
	LabelFsm         *widget.Label       //FSM資訊
	EntryTableStatus *widget.Entry       //牌桌資訊
	LabalRoom        *widget.Label       //大廳資訊
	EntrySendMessage *widget.SelectEntry //可送出的自定義Json
	comboAgent       *widget.Select      // 使用的Agent
}

func NewBaseClient(setting ClientConfig) *BaseClient {
	baseClient := &BaseClient{
		GameId: setting.GameId,
	}

	baseClient.EntryTableStatus = widget.NewMultiLineEntry()
	baseClient.EntryTableStatus.SetText("")
	baseClient.EntryTableStatus.SetMinRowsVisible(3)
	baseClient.EntryTableStatus.Wrapping = fyne.TextWrapBreak
	baseClient.LabelFsm = widget.NewLabel("Fsm:")
	baseClient.LabalRoom = widget.NewLabel("大廳資訊:")
	baseClient.CustomMessage = append(baseClient.CustomMessage, "{\"Login\":{\"Token\":\"\",\"Name\":\"Danny\"}}")
	baseClient.EntrySendMessage = widget.NewSelectEntry(baseClient.CustomMessage)
	return baseClient
}

func (t *BaseClient) ConnectServer() (bool, error) {
	dialer := websocket.Dialer{}
	connect, _, err := dialer.Dial(constant.GameServerURL[t.Env], nil)
	if nil != err {
		t.LogInfo(utils.LOG_ERROR, err.Error())
		return false, err
	}
	t.Connect = connect
	t.isLogIn = true
	go t.ListenMessage()
	return true, nil
}

func (t *BaseClient) ListenMessage() {
	for {
		if !t.isLogIn {
			break
		}

		if t.Connect == nil {
			return
		}

		messageType, bmsg, err := t.Connect.ReadMessage()
		if err != nil {
			t.LogInfo(utils.LOG_ERROR, err.Error())
			break
		}
		messageData, err := utils.Unpackage(bmsg)
		if err != nil {
			t.LogInfo(utils.LOG_ERROR, "Unpackage failed. data:"+string(bmsg))
		}
		switch messageType {
		case websocket.TextMessage:
			response := &utils.RespBase{}
			err = json.Unmarshal(messageData, &response)
			if err != nil {
				t.LogInfo(utils.LOG_ERROR, "json unmarshal failed. data:"+string(messageData))
				continue
			}
			t.LogInfo(utils.LOG_INFO, "Receive Message:"+string(messageData))
			t.CheckResponse(response)
		default:
			t.LogInfo(utils.LOG_ERROR, "Error type message. data:"+string(messageData))
		}
	}
}

func (t *BaseClient) LogInServer() (bool, error) {
	var data struct {
		Login struct {
			Token string
			Name  string
		}
	}

	t.LoginName = t.EntryClientName.Text
	if t.comboAgent.Selected != "" {
		ok, token, error := utils.DccLogin(constant.DccURL[t.Env], t.comboAgent.Selected, t.LoginName)
		if !ok {
			return false, error
		}
		data.Login.Token = token
	} else {
		data.Login.Name = t.LoginName
	}
	// data.Login.Token
	t.SendMessage(data)
	return true, nil
}

func (t *BaseClient) DisconnectServer() bool {
	t.Connect.Close()
	t.isLogIn = false
	return true
}

func (t *BaseClient) GetGamePeriod(response *utils.RespBase) string {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}

	t.Fsm = info["Fsm"].(string)
	gid := info["Gid"].(string)
	sec := int(info["RemainSec"].(float64))

	return fmt.Sprintf("Fsm:%s Gid:%s Sec:%d", t.Fsm, gid, sec)
}

func (t *BaseClient) LogInfo(logLevel int, message string) {
	if !t.CheckActBulletin.Checked {
		if strings.Contains(message, ACT_BULLETIN) {
			return
		}
	}
	content := time.Now().Format("2006/01/02 15:04:05 ") + t.LoginName + " " + message + "\n"
	if len(t.ClientMessage.Text) > 5000 {
		t.ClientMessage.SetText(content)
	} else {
		t.ClientMessage.SetText(content + t.ClientMessage.Text)
	}
	utils.LogInfo(logLevel, message, nil)
}

func (t *BaseClient) CreateTopSection() *fyne.Container {
	label1 := widget.NewLabel("選擇測試環境")
	comboClientEnv := widget.NewSelect([]string{"LOCAL", "DEV_TEST", "DEV", "QC", "QA"}, func(value string) {
		t.setClientEnv(value)
	})
	buttonConnect := widget.NewButton("Connect", nil)
	t.comboAgent = widget.NewSelect([]string{}, func(value string) {
	})

	label6 := widget.NewLabel("選擇登入代理(不選則使用預設)")
	label7 := widget.NewLabel("")

	button1 := widget.NewButton("讀取列表", func() {
		if comboClientEnv.Selected == "PROD" || comboClientEnv.Selected == "ETE" ||
			comboClientEnv.Selected == "QA" || comboClientEnv.Selected == "QC" ||
			comboClientEnv.Selected == "DEV" || comboClientEnv.Selected == "DEV_TEST" {
			errorMessage, list := utils.GetDccAgentInfo(constant.DccAgentURL[t.Env])
			if errorMessage != "" {
				t.ClearAgentInfo()
				label7.SetText("讀取失敗")
			} else {
				t.comboAgent.Options = list
				t.AgentEnv = t.Env
				label7.SetText("讀取成功")
			}
		} else {
			t.ClearAgentInfo()
			label7.SetText("")
		}
	})

	button2 := widget.NewButton("清空列表", func() {
		t.ClearAgentInfo()
		label7.SetText("")
	})

	sectionAgent := container.NewHBox(label6, button1, button2, label7)

	label2 := widget.NewLabel("輸入帳號名稱:")
	t.EntryClientName = widget.NewEntry()
	buttonLogin := widget.NewButton("Login", nil)

	buttonConnect.OnTapped = func() {
		if buttonConnect.Text == "Connect" {
			if t.Env == constant.ENV_NONE {
				t.LogInfo(utils.LOG_ERROR, "未選擇環境")
				return
			}
			if ok, err := t.ConnectServer(); !ok {
				t.LogInfo(utils.LOG_ERROR, "Connect failed:"+err.Error())
				return
			}

			IsConnect = true
			t.LogInfo(utils.LOG_INFO, "Connect Success")
			buttonConnect.Text = "Disconnect"
			buttonConnect.Refresh()
		} else {
			t.DisconnectServer()
			buttonLogin.Enable()
			IsConnect = false
			t.LogInfo(utils.LOG_INFO, "連線中斷")
			buttonConnect.Text = "Connect"
			buttonConnect.Refresh()
		}
	}

	buttonLogin.OnTapped = func() {
		if buttonLogin.Text == "Login" {
			if t.EntryClientName.Text == "" {
				t.LogInfo(utils.LOG_ERROR, "未輸入帳號")
				return
			}
			if t.Connect == nil {
				t.LogInfo(utils.LOG_ERROR, "尚未連線")
				return
			}
			if ok, err := t.LogInServer(); !ok {
				t.LogInfo(utils.LOG_ERROR, "Login failed:"+err.Error())
				return
			}
			buttonLogin.Disable()
		}
	}

	boxConnect := container.NewHBox(label1, comboClientEnv, buttonConnect)
	boxLogin := container.NewBorder(nil, nil, label2, buttonLogin, t.EntryClientName)

	return container.NewVBox(boxConnect, sectionAgent, t.comboAgent, boxLogin)
}

func (t *BaseClient) ClearAgentInfo() {
	t.comboAgent.Options = []string{}
	t.comboAgent.Selected = ""
	t.AgentEnv = constant.ENV_NONE
	t.comboAgent.Refresh()
}

func (t *BaseClient) CreateBottomSection(c *fyne.Container) {

	label1 := widget.NewLabel("自定義Json訊息")
	buttonSendMessage := widget.NewButton("送出", func() {
		message := t.EntrySendMessage.Text
		t.ShowMessage(message)
		if ok, err := utils.SendCustomizeMessage(t.Connect, t.LoginName, message, nil); !ok {
			t.ShowMessage("不是json的格式:" + message + " Error message:" + err.Error())
		}
	})
	c.Add(label1)
	c.Add(t.EntrySendMessage)
	c.Add(buttonSendMessage)

	t.ClientMessage = widget.NewMultiLineEntry()
	t.ClientMessage.SetMinRowsVisible(12)
	t.ClientMessage.Wrapping = fyne.TextWrapBreak
	c.Add(t.ClientMessage)

	buttonClientMessageClear := widget.NewButton("清除訊息", func() {
		t.ClientMessage.SetText("")
	})

	t.CheckActBulletin = widget.NewCheck("顯示跑馬燈訊息", nil)
	t.CheckActBulletin.Checked = false

	box := container.NewHBox(buttonClientMessageClear, t.CheckActBulletin)
	c.Add(box)

}

func (t *BaseClient) setClientEnv(value string) {
	switch value {
	case "LOCAL":
		t.Env = constant.ENV_LOCAL
	case "DEV_TEST":
		t.Env = constant.ENV_DEV_TEST
	case "DEV":
		t.Env = constant.ENV_DEV
	case "QC":
		t.Env = constant.ENV_QC
	case "QA":
		t.Env = constant.ENV_QA
	}
}

func (t *BaseClient) SendMessage(message interface{}) (bool, error) {
	t.ShowMessage(message)
	return utils.SendMessage(t.Connect, t.LoginName, message, nil)
}

func (t *BaseClient) ShowMessage(message interface{}) {
	now := time.Now().Format("2006/01/02 15:04:05")
	content := fmt.Sprintf("%s %s SendMessage=%+v\n", now, t.LoginName, message)
	if len(t.ClientMessage.Text) > 5000 {
		t.ClientMessage.SetText(content)
	} else {
		t.ClientMessage.SetText(content + t.ClientMessage.Text)
	}
}

func (t *BaseClient) CheckBaseResponse(response *utils.RespBase) bool {
	if response.Code != constant.ERROR_CODE_SUCCESS {
		return true
	}
	switch response.Ret {
	case ACT_GAME_PERIOD:
		t.LabelFsm.SetText(t.GetGamePeriod(response))
	case ACT_GOLD:
		detail, ok := response.Data.(map[string]interface{})
		if !ok {
			return true
		}
		t.Gold = detail["Gold"].(float64)
		return true
	case RET_LOGIN:
		detail, ok := response.Data.(map[string]interface{})
		if !ok {
			return true
		}
		t.Gold = detail["Gold"].(float64)
		t.SWGold = detail["single_wallet_gold"].(float64)
		t.WalletType = int(detail["wallet_type"].(float64))
	case ACT_SW_GOLD:
		detail, ok := response.Data.(map[string]interface{})
		if !ok {
			return true
		}
		t.SWGold = detail["SWGold"].(float64)
	}

	return false
}

func (t *BaseClient) AddTableStatus(message string) {
	t.EntryTableStatus.SetText(t.EntryTableStatus.Text + message)
	t.EntryTableStatus.CursorRow = 100
}
func (t *BaseClient) SetTableStatus(message string) {
	t.EntryTableStatus.SetText(message)
}
func (t *BaseClient) SetTableStatusClear() {
	t.EntryTableStatus.SetText("")
}

func (t *BaseClient) CheckTableStatus(status int) string {
	tableStatus := "不明"
	switch status {
	case 0:
		tableStatus = "關閉"
	case 1:
		tableStatus = "開啟"
	case 2:
		tableStatus = "維護"
	}
	return tableStatus
}
