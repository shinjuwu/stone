package gate

import (
	"net"
)

type Agent interface {
	WriteMsg(msg interface{})
	WriteMsgByte(data [][]byte)
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
	GetSession() Session
	SetAgentID(id int64)
	GetAgentID() int64
	IsClosed() bool
	Heartbeat()
}

type Session interface {
	GetIP() string
	GetNetwork() string
	GetUserid() string
	GetSessionid() string
	GetServerid() string
	GetSettings() map[string]string
	SetIP(ip string)
	SetNetwork(network string)
	SetUserid(userid string)
	SetSessionid(sessionid string)
	SetServerid(serverid string)
	SetSettings(settings map[string]string)
	Serializable() ([]byte, error)
	Update() (err error)
	Bind(Userid string) string
	UnBind() string
	Push() string
	Set(key string, value string) string
	SetPush(key string, value string) string
	Get(key string) string
	Remove(key string) string
	Send(id string, data interface{}) string
	SendNR(id string, data interface{}) string
	SendBatch(Sessionids string, data interface{}) string
	IsConnect(Userid string) (bool, string)
	IsGuest() bool
	JudgeGuest(judgeGuest func(session Session) bool)
	Close() (err error)
	CloseMultiSession(key string) string
	BroadCast(data interface{}) string
}

/*
*
Session 持久化
*/
type StorageHandler interface {
	Storage(Userid string, session Session) (err error)
	Delete(Userid string) (err error)
	Query(Userid string) (data []byte, err error)
	Heartbeat(Userid string)
	GetRedisSessionID(Userid string) string
}

type SessionLearner interface {
	Connect(a Agent)
	DisConnect(a Agent)
	CloseMultiSession(a Agent)
}

type AgentLearner interface {
	Connect(a Agent)
	DisConnect(a Agent)
	OnHeartbeat(a Agent)
}

/*
*
net代理服务 处理器
*/
type GateHandler interface {
	Bind(args []interface{}) interface{}   //Bind the session with the the Userid.
	UnBind(args []interface{}) interface{} //UnBind the session with the the Userid.
	Set(args []interface{}) interface{}    //Set values (one or many) for the session.
	Remove(args []interface{}) interface{} //Remove value from the session.
	Push(args []interface{}) interface{}
	Send(args []interface{}) interface{} //Send message
	SendBatch(args []interface{}) interface{}
	BroadCast(args []interface{}) interface{}
	FilterBroadCast(args []interface{}) interface{}

	IsConnect(args []interface{}) interface{}
	Close(args []interface{}) interface{}
	Update(args []interface{}) interface{}
	OnDestory()
	CloseMultiSession(args []interface{}) interface{}
	NewAgent(args []interface{})
	CloseAll(args []interface{})
}
