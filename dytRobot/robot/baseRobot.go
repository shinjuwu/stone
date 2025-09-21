package robot

import (
	"dytRobot/constant"
	"dytRobot/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	RET_LOGIN = "Login"

	ACT_BULLETIN    = "ActBulletin"
	ACT_GAME_PERIOD = "ActGamePeriod"
	ACT_GOLD        = "ActGold"
)

type BaseRobot struct {
	Env            int
	AccessDcc      bool
	AgentName      string
	LoginName      string
	UserId         int64
	GameId         int
	TableId        int
	Connect        *websocket.Conn
	MessageChannel chan string
	AliveChannel   chan interface{}
	LastAliveTime  int64
	LastHeartTime  int64
	LastPlayTime   int64
	Fsm            string
	PlayLimit      int //下注次數，達到時停止，設0為無限制
	PlayCount      int //實際下注次數
	isDisconnect   bool
	EnterInfo      map[int]float64
	WalletType     int
	SWGold         float64

	CheckCommand func(response *utils.RespBase) int
}

func NewBaseRobot(setting RobotConfig) *BaseRobot {
	baseRobot := &BaseRobot{
		Env:            setting.Env,
		AccessDcc:      setting.AccessDcc,
		AgentName:      setting.AgentName,
		LoginName:      setting.LoginName,
		GameId:         setting.GameId,
		TableId:        setting.TableId,
		MessageChannel: setting.MessageChannel,
		AliveChannel:   setting.AliveChannel,
		LastPlayTime:   time.Now().Unix(),
		Fsm:            "Init",
		PlayLimit:      setting.Count,
	}
	return baseRobot
}

func (t *BaseRobot) StartRobot() {
	dialer := websocket.Dialer{}
	connect, _, err := dialer.Dial(constant.GameServerURL[t.Env], nil)
	if nil != err {
		t.LogInfo(utils.LOG_ERROR, err.Error())
		return
	}
	t.Connect = connect
	defer t.Connect.Close()

	ok, err := t.Login()
	if !ok {
		t.LogInfo(utils.LOG_ERROR, err.Error())
		return
	}

	for {
		if t.DisConnect() {
			return
		}

		t.SendAlive()
		t.SendHeartbeat()

		messageType, bmsg, err := t.Connect.ReadMessage()
		if nil != err {
			t.LogInfo(utils.LOG_ERROR, err.Error())
			if t.ReConnect() {
				continue
			} else {
				return
			}
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

			level := utils.LOG_DEBUG
			if response.Code != constant.ERROR_CODE_SUCCESS {
				level = utils.LOG_ERROR
			}
			t.LogInfo(level, "Receive Message:"+string(messageData))

			if t.CheckCommand(response) == RESPONSE_EXCUTED_FAILED {
				return
			}

		default:
			t.LogInfo(utils.LOG_ERROR, "Error type message. data:"+string(messageData))
		}

		if !t.CheckAlive() {
			t.LogInfo(utils.LOG_ERROR, "一段時間未收到任何遊戲階段訊息")
			t.LastPlayTime = time.Now().Unix()
			if t.ReConnect() {
				continue
			} else {
				return
			}
		}
	}
}

func (t *BaseRobot) ReConnect() bool {
	for i := 0; i < RECONNECT_COUNT; i++ {
		t.LogInfo(utils.LOG_INFO, fmt.Sprintf("重連第%d次，等候%d秒執行", i+1, RECONNECT_SECONDS))
		time.Sleep(RECONNECT_SECONDS * time.Second)

		t.Connect.Close()
		dialer := websocket.Dialer{}
		connect, _, err := dialer.Dial(constant.GameServerURL[t.Env], nil)

		if nil != err {
			t.LogInfo(utils.LOG_ERROR, err.Error())
			continue
		}
		t.Connect = connect

		ok, err := t.Login()
		if !ok {
			t.LogInfo(utils.LOG_ERROR, err.Error())
			continue
		}
		t.LogInfo(utils.LOG_INFO, "重新連線成功")
		return true
	}

	t.LogInfo(utils.LOG_INFO, "重連共"+strconv.Itoa(RECONNECT_COUNT)+"次失敗")
	return false
}

func (t *BaseRobot) CheckBaseCommand(response *utils.RespBase) int {
	switch response.Ret {
	case ACT_GAME_PERIOD:
		t.SetGamePeriod(response)
	case ACT_BULLETIN:
		return RESPONSE_EXCUTED_SUCCESS
	case RET_LOGIN:
		detail, ok := response.Data.(map[string]interface{})
		if !ok {
			return RESPONSE_EXCUTED_FAILED
		}
		t.SWGold = detail["single_wallet_gold"].(float64)
		t.WalletType = int(detail["wallet_type"].(float64))
	}
	t.LastPlayTime = time.Now().Unix()
	return RESPONSE_NO_SUTIALBE
}

func (t *BaseRobot) SetGamePeriod(response *utils.RespBase) {
	info, ok := response.Data.(map[string]interface{})
	if !ok {
		return
	}

	t.Fsm, _ = info["Fsm"].(string)
}

func (t *BaseRobot) LogInfo(logLevel int, message string) {
	content := t.LoginName + " " + message
	utils.LogInfo(logLevel, content, t.MessageChannel)
}

func (t *BaseRobot) SendAlive() {
	now := time.Now().Unix()
	if (now - t.LastAliveTime) > SEND_ALIVE_SECONDS {
		t.LastAliveTime = now
		var data = RobotAlive{
			RobotName:  t.LoginName,
			Fsm:        t.Fsm,
			UpdateTime: t.LastAliveTime,
			PlayCount:  t.PlayCount,
		}
		t.AliveChannel <- data
	}
}

func (t *BaseRobot) CheckAlive() bool {
	now := time.Now().Unix()
	return (now - t.LastPlayTime) < PLAY_ALIVE_TIME
}

func (t *BaseRobot) Login() (bool, error) {
	var data struct {
		Login struct {
			Token string
			Name  string
		}
	}

	if t.AccessDcc {
		ok, token, error := utils.DccLogin(constant.DccURL[t.Env], t.AgentName, t.LoginName)
		if !ok {
			return false, error
		}
		data.Login.Token = token
	} else {
		data.Login.Name = t.LoginName
	}

	utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
	return true, nil
}

func (t *BaseRobot) Deposit() (bool, error) {
	if !t.AccessDcc {
		return utils.GameDeposit(constant.GameDepositURL[t.Env], t.UserId, UP_GOLD)
	}

	ok, error := utils.DccDeposit(constant.DccURL[t.Env], t.AgentName, t.LoginName, UP_GOLD)
	if !ok {
		return false, error
	}
	return true, nil
}

func (t *BaseRobot) SendHeartbeat() {
	now := time.Now().Unix()
	if (now - t.LastHeartTime) > SEND_HEART_BEAT_SECONDS {
		t.LastHeartTime = now
		var data struct {
			Heartbeat struct{}
		}
		utils.SendMessage(t.Connect, t.LoginName, data, t.MessageChannel)
	}
}

func (t *BaseRobot) SetDisConnect() {
	t.isDisconnect = true
}

func (t *BaseRobot) DisConnect() bool {
	return t.isDisconnect
}
