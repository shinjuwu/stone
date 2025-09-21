package robot

const (
	UP_GOLD                 = 1000000
	SEND_ALIVE_SECONDS      = 5
	SEND_HEART_BEAT_SECONDS = 150
	PLAY_ALIVE_TIME         = 60 //單位秒。超過時間都沒收到玩牌階段，就當未在玩
	RECONNECT_SECONDS       = 5
	RECONNECT_COUNT         = 3
)

const (
	RESPONSE_NO_SUTIALBE     int = iota //server回應未有合適的對應
	RESPONSE_EXCUTED_SUCCESS            //server回應有找到合適的對應，並不需繼續找
	RESPONSE_EXCUTED_FAILED
)

type RobotConfig struct {
	Env            int
	AccessDcc      bool
	AgentName      string
	LoginName      string
	GameId         int
	TableId        int
	MessageChannel chan string
	AliveChannel   chan interface{}
	Count          int //下注次數，達到時停止，設0為無限制
}

type RobotAlive struct {
	RobotName  string
	Fsm        string
	UpdateTime int64
	PlayCount  int
}
