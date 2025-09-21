package gate

import (
	"fmt"
	"strings"

	"collie/log"
	"collie/util"
)

type handler struct {
	gate     Gate
	sessions *util.BeeMap //use sessionID be key
}

func NewGateHandler(gate Gate) *handler {
	handler := &handler{
		gate:     gate,
		sessions: util.NewBeeMap(),
	}

	return handler
}

func (h *handler) OnDestory() {
	for _, v := range h.sessions.Items() {
		v.(Agent).Close()
	}
	h.sessions.DeleteAll()
}

// Update : f(sessionID string) error
func (h *handler) Update(args []interface{}) interface{} {
	sessionID := args[0].(string)
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		return "No agent found"
	}
	return ""
}

// Bind : f(sessionID string, userID string) error
func (h *handler) Bind(agrs []interface{}) interface{} {
	Sessionid := agrs[0].(string)
	Userid := agrs[1].(string)
	agent := h.sessions.Get(Sessionid)
	if agent == nil {
		return "No Sesssion found"
	}
	agent.(Agent).GetSession().SetUserid(Userid)

	if h.gate.GetStorageHandler() != nil && agent.(Agent).GetSession().GetUserid() != "" {
		data, err := h.gate.GetStorageHandler().Query(Userid)
		if err == nil && data != nil {

			imSession, err := h.gate.NewSession(data)
			if err == nil {
				if agent.(Agent).GetSession().GetSettings() == nil {
					agent.(Agent).GetSession().SetSettings(imSession.GetSettings())
				} else {
					settings := imSession.GetSettings()
					if settings != nil {
						for k, v := range settings {
							if _, ok := agent.(Agent).GetSession().GetSettings()[k]; ok {
							} else {
								agent.(Agent).GetSession().GetSettings()[k] = v
							}
						}
					}
					h.gate.GetStorageHandler().Storage(Userid, agent.(Agent).GetSession())
				}
			} else {
			}
		}
	}
	return agent.(Agent).GetSession()
}

func (h *handler) IsConnect(args []interface{}) interface{} {
	Userid := args[0].(string)
	for _, agent := range h.sessions.Items() {
		if agent.(Agent).GetSession().GetUserid() == Userid {
			return !agent.(Agent).IsClosed()
		}
	}
	return false
}

func (h *handler) UnBind(args []interface{}) interface{} {
	sessionID := args[0].(string)
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		return "No Session found"
	}
	agent.(Agent).GetSession().SetUserid("")
	return ""
}

func (h *handler) Push(args []interface{}) interface{} {
	sessionID := args[0].(string)
	settings := args[1].(map[string]string)
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		return "No Session found"
	}
	agent.(Agent).GetSession().SetSettings(settings)
	result := agent.(Agent).GetSession()
	if h.gate.GetStorageHandler() != nil && agent.(Agent).GetSession().GetUserid() != "" {
		err := h.gate.GetStorageHandler().Storage(agent.(Agent).GetSession().GetUserid(), agent.(Agent).GetSession())
		if err != nil {
			log.Debug("gate session storage failure : %s", err.Error())
		}
	}
	return result
}

func (h *handler) Set(args []interface{}) interface{} {
	sessionID := args[0].(string)
	key := args[1].(string)
	value := args[2].(string)
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		return "No Session found"
	}
	agent.(Agent).GetSession().GetSettings()[key] = value

	if h.gate.GetStorageHandler() != nil && agent.(Agent).GetSession().GetUserid() != "" {
		err := h.gate.GetStorageHandler().Storage(agent.(Agent).GetSession().GetUserid(), agent.(Agent).GetSession())
		if err != nil {
			log.Error("gate session storage failure : %s", err.Error())
		}
	}
	return ""
}

func (h *handler) Remove(args []interface{}) interface{} {
	sessionID := args[0].(string)
	key := args[1].(string)
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		return "No Session found"
	}
	delete(agent.(Agent).GetSession().GetSettings(), key)

	if h.gate.GetStorageHandler() != nil && agent.(Agent).GetSession().GetUserid() != "" {
		err := h.gate.GetStorageHandler().Storage(agent.(Agent).GetSession().GetUserid(), agent.(Agent).GetSession())
		if err != nil {
			log.Error("gate session storage failure :%s", err.Error())
		}
	}

	return ""
}

func (h *handler) Send(args []interface{}) interface{} {
	sessionID := args[0].(string)
	data := args[1].(*[]byte)
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		return "No session found"
	}
	agent.(Agent).WriteMsgByte([][]byte{*data})
	return ""
}

func (h *handler) SendBatch(args []interface{}) interface{} {
	sessionIDStr := args[0].(string)
	data := args[1].(*[]byte)
	sessionIDs := strings.Split(sessionIDStr, ",")
	var count int = 0
	for _, sessionID := range sessionIDs {
		agent := h.sessions.Get(sessionID)
		if agent == nil {
			continue
		}
		agent.(Agent).WriteMsgByte([][]byte{*data})
		count++
	}
	return count
}

func (h *handler) BroadCast(args []interface{}) interface{} {
	data := args[0]
	var count int64 = 0
	for _, agent := range h.sessions.Items() {
		agent.(Agent).WriteMsg(data)
		count++
	}
	return count
}

func (h *handler) FilterBroadCast(args []interface{}) interface{} {
	data := args[0]
	filterFunc := args[1].(func(a Agent) bool)
	var count int64 = 0
	for _, agent := range h.sessions.Items() {
		if filterFunc(agent.(Agent)) {
			agent.(Agent).WriteMsg(data)
			count++
		}
	}
	return count
}

func (h *handler) Close(args []interface{}) interface{} {
	sessionID := args[0].(string)
	agent := h.sessions.Get(sessionID)
	if agent == nil {
		return "No session found"
	}
	agent.(Agent).Close()
	return ""
}

func (h *handler) CloseMultiSession(args []interface{}) interface{} {
	key := args[0].(string)
	if key == "" {
		return "Invaild key"
	}

	sessionID := h.gate.GetStorageHandler().GetRedisSessionID(key)
	if sessionID == "" {
		return ""
	}

	agent := h.sessions.Get(sessionID)
	if agent != nil {
		h.gate.GetSessionLearner().CloseMultiSession(agent.(Agent))
		//agent.(Agent).Close()
	}
	return ""
}

func (h *handler) Connect(a Agent) {
	if a.GetSession() != nil {
		h.sessions.Set(a.GetSession().GetSessionid(), a)
	}
	if h.gate.GetSessionLearner() != nil {
		h.gate.GetSessionLearner().Connect(a)
	}
}

func (h *handler) DisConnect(a Agent) {
	if a.GetSession() != nil {
		h.sessions.Delete(a.GetSession().GetSessionid())
	}
	if h.gate.GetSessionLearner() != nil {
		h.gate.GetSessionLearner().DisConnect(a)
	}
}

func (h *handler) OnHeartbeat(a Agent) {
	if h.gate.GetStorageHandler() != nil {
		h.gate.GetStorageHandler().Heartbeat(a.GetSession().GetUserid())
	}
}

func (h *handler) NewAgent(args []interface{}) {
	fmt.Println("init agent.................")
}

func (h *handler) CloseAll([]interface{}) {
	h.OnDestory()
}
