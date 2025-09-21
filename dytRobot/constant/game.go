package constant

const (
	ENV_NONE = iota
	ENV_LOCAL
	ENV_DEV_TEST
	ENV_DEV
	ENV_QC
	ENV_QA
	ENV_ETE
	ENV_PROD
)

const (
	// single wallet type
	SW_TYPE_NORMAL = 0
	SW_TYPE_SINGLE = 1
)

var GameServerURL map[int]string = map[int]string{
	ENV_LOCAL:    "ws://localhost:10001",
	ENV_DEV_TEST: "ws://172.30.0.167:20001",
	ENV_DEV:      "ws://172.30.0.167:10001",
	ENV_QC:       "ws://172.30.0.154:10001",
	ENV_QA:       "ws://172.30.0.164:10001",
	ENV_ETE:      "wss://wss.ete-demo.com:10001",
	ENV_PROD:     "wss://wss.86-gaming.com:10001",
}

var GameDepositURL map[int]string = map[int]string{
	ENV_LOCAL:    "http://localhost:9642/deposit",
	ENV_DEV_TEST: "http://172.30.0.167:9742/deposit",
	ENV_DEV:      "http://172.30.0.167:9642/deposit",
	ENV_QC:       "http://172.30.0.154:9642/deposit",
	ENV_QA:       "http://172.30.0.164:9642/deposit",
}

var DccURL map[int]string = map[int]string{
	ENV_DEV_TEST: "http://172.30.0.168:9986/",
	ENV_DEV:      "http://172.30.0.168:9986/",
	ENV_QC:       "http://172.30.0.150:9986/",
	ENV_QA:       "http://172.30.0.165:9986/",
	ENV_ETE:      "http://34.142.251.42:9986/",
	ENV_PROD:     "https://api.86-gaming.com/",
}

var DccAgentURL map[int]string = map[int]string{
	ENV_DEV_TEST: "http://172.30.0.168:9986/",
	ENV_DEV:      "http://172.30.0.168:9986/",
	ENV_QC:       "http://172.30.0.150:9986/",
	ENV_QA:       "http://172.30.0.160:9986/",
	ENV_ETE:      "http://34.142.193.117:9986/",
	ENV_PROD:     "https://backend.86-gaming.com/",
}

const (
	//成功
	ERROR_CODE_SUCCESS int = 0
	//餘額不足
	ERROR_CODE_ERROR_MONEY_NOT_ENOUGH int = 200008
	//超過下注
	ERROR_CODE_ERROR_BET_LIMIT int = 200011
	//攜入分數低於下限
	ERROR_CODE_BRING_MONEY_LOWER_LIMIT int = 200012
	//遊戲已結束
	ERROR_CODE_ERROR_CAME_FINISHED int = 200015
	//已离开游戏，不需再次离开游戏。
	ERROR_CODE_ERROR_NO_NEED_TO_QUITGAME int = 200029
)

var GameIDList []int = []int{
	1001, 1002, 1003, 1004, 1005, 1006, 1007, 1008, 1009, 1010,
	2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009, 2010, 2011,
	3001, 3002, 3003,
	4001, 4002, 4003,
	5001,
	9001}

var GameIDNameTable map[int]string = map[int]string{
	1001: "百家樂",
	1002: "番攤",
	1003: "色碟",
	1004: "魚蝦蟹",
	1005: "百人骰寶",
	1006: "鬥雞",
	1007: "賽狗",
	1008: "火箭",
	1009: "安達巴哈",
	1010: "輪盤",
	2001: "21點",
	2002: "三公",
	2003: "搶庄牛牛",
	2004: "德州撲克",
	2005: "拉密",
	2006: "炸金花",
	2007: "博丁",
	2008: "越南Catte",
	2009: "十三水",
	2010: "土耳其麻將",
	2011: "印度炸金花",
	3001: "水果機",
	3002: "三國捕魚",
	3003: "彈珠檯",
	4001: "水果777",
	4002: "巨齒鯊",
	4003: "邁達斯之手",
	5001: "好友德撲",
	9001: "jackpot",
}

var NameTableGameID map[string]int = map[string]int{
	"百家樂":     1001,
	"番攤":      1002,
	"色碟":      1003,
	"魚蝦蟹":     1004,
	"百人骰寶":    1005,
	"鬥雞":      1006,
	"賽狗":      1007,
	"火箭":      1008,
	"安達巴哈":    1009,
	"輪盤":      1010,
	"21點":     2001,
	"三公":      2002,
	"搶庄牛牛":    2003,
	"德州撲克":    2004,
	"拉密":      2005,
	"炸金花":     2006,
	"博丁":      2007,
	"越南Catte": 2008,
	"十三水":     2009,
	"土耳其麻將":   2010,
	"印度炸金花":   2011,
	"水果機":     3001,
	"三國捕魚":    3002,
	"彈珠檯":     3003,
	"水果777":   4001,
	"巨齒鯊":     4002,
	"邁達斯之手":   4003,
	"好友德撲":    5001,
	"jackpot": 9001,
}

var GameID2GameType map[int]int = map[int]int{
	1001: 1,
	1002: 1,
	1003: 1,
	1004: 1,
	1005: 1,
	1006: 1,
	1007: 1,
	1008: 1,
	1009: 1,
	1010: 1,
	2001: 2,
	2002: 2,
	2003: 2,
	2004: 2,
	2005: 2,
	2006: 2,
	2007: 2,
	2008: 2,
	2009: 2,
	2010: 2,
	2011: 2,
	3001: 3,
	3002: 3,
	3003: 3,
	4001: 4,
	4002: 4,
	4003: 4,
	5001: 6,
	9001: 1,
}

var RoomTypeNum map[int]int = map[int]int{
	1001: 4,
	1002: 4,
	1003: 4,
	1004: 4,
	1005: 4,
	1006: 4,
	1007: 4,
	1008: 4,
	1009: 4,
	1010: 4,
	2001: 4,
	2002: 4,
	2003: 4,
	2004: 8,
	2005: 4,
	2006: 4,
	2007: 4,
	2008: 4,
	2009: 8,
	2010: 4,
	2011: 4,
	3001: 4,
	3002: 3,
	3003: 4,
	4001: 1,
	4002: 1,
	4003: 1,
	5001: 1,
	9001: 1,
}

var RoomType2Name map[int]string = map[int]string{
	0: "新手房",
	1: "普通房",
	2: "高級房",
	3: "大師房",
	4: "初级场",
	5: "中级场",
	6: "高级场",
	7: "至尊场",
}

var GameName2Type map[string]int = map[string]int{
	"新手房": 0,
	"普通房": 1,
	"高級房": 2,
	"大師房": 3,
	"初级场": 4,
	"中级场": 5,
	"高级场": 6,
	"至尊场": 7,
}

var TexasBringGold map[int]float64 = map[int]float64{
	0: 200,
	1: 800,
	2: 1600,
	3: 4000,
	4: 200,
	5: 800,
	6: 1600,
	7: 4000,
}

type GameSetting struct {
	IsEnable  bool
	Room      [8]bool
	RobotNum  int
	PlayCount int
}

func GetTableID(gameID int, roomType int, playRoom int) int {
	var tableId int
	if gameType, ok := GameID2GameType[gameID]; ok {
		if gameType == 1 {
			tableId = gameID*100 + roomType*10 + playRoom
		} else {
			tableId = gameID*10 + roomType
		}
	}
	return tableId
}

func GetGameRoomName(gameID int, tableId int) (string, string) {
	var t int = 1
	if GameID2GameType[gameID] == 1 {
		t = 10
	}
	return GameIDNameTable[gameID], RoomType2Name[tableId/t%gameID]
}
