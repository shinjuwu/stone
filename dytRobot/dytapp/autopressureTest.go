//go:build auto

package dytapp

import (
	"dytRobot/constant"
	"dytRobot/robot"
	"dytRobot/robot/andarbahar"
	"dytRobot/robot/baccarat"
	"dytRobot/robot/blackjack"
	"dytRobot/robot/bullbull"
	"dytRobot/robot/catte"
	"dytRobot/robot/chinesepoker"
	"dytRobot/robot/cockfight"
	"dytRobot/robot/colordisc"
	"dytRobot/robot/dogracing"
	"dytRobot/robot/fantan"
	"dytRobot/robot/fruit777slot"
	"dytRobot/robot/fruitslot"
	"dytRobot/robot/goldenflower"
	"dytRobot/robot/hundredsicbo"
	"dytRobot/robot/megsharkslot"
	"dytRobot/robot/midasslot"
	"dytRobot/robot/okey"
	"dytRobot/robot/plinko"
	"dytRobot/robot/pokdeng"
	"dytRobot/robot/prawncrab"
	"dytRobot/robot/rocket"
	"dytRobot/robot/rummy"
	"dytRobot/robot/sangong"
	"dytRobot/robot/teenpatti"
	"dytRobot/robot/texas"
	"dytRobot/utils"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

var robotEnv int = constant.ENV_NONE
var robotName string = "dayon"
var robotStartNum int = 1
var robotIntervalSec int = 3

var robotSttings map[int]*constant.GameSetting

var messageChannel chan string
var aliveChannel chan interface{}
var robotAlive sync.Map

const (
	ALIVE_SECONDS = 30
)

func setPressureTestEnv() {
	initRobotSettings()
	makeRobotRun()
}

func initRobotSettings() {
	robotSttings = make(map[int]*constant.GameSetting)
	robotEnv = constant.ENV_LOCAL
	for _, gameID := range constant.GameIDList {
		setting := &constant.GameSetting{}
		roomNum := constant.RoomTypeNum[gameID]
		setting.IsEnable = true
		setting.RobotNum = 2
		for i := 0; i < roomNum; i++ {
			setting.Room[i] = true
		}
		robotSttings[gameID] = setting
	}
}

func makeRobotRun() {
	//開啟訊息頻道
	if messageChannel == nil {
		messageChannel = make(chan string, 50)
		go ReceiveMessage()
	}

	//建立RobotList
	makeRobotList()

	startNum := robotStartNum
	//開始放機器人進遊戲
	for gameId, data := range robotSttings {
		if !data.IsEnable {
			continue
		}
		robotNum := data.RobotNum
		playCount := data.PlayCount

		for roomType, roomSet := range data.Room {
			if !roomSet {
				continue
			}

			roomCount := 1
			if gameId == 1001 {
				//放到類型房的4小間房
				roomCount = 4
			}

			for roomNum := 1; roomNum <= roomCount; roomNum++ {
				tableID := constant.GetTableID(gameId, roomType, roomNum)
				go startRobot(gameId, tableID, startNum, robotNum, playCount)
				startNum += robotNum
			}
		}
	}

}

func ReceiveMessage() {
	for {
		message := <-messageChannel
		fmt.Println(message)
	}
}

func makeRobotList() {
	var sum int
	for gameId, data := range robotSttings {
		if !data.IsEnable {
			continue
		}
		for _, roomSet := range data.Room {
			if !roomSet {
				continue
			}

			roomCount := 1
			if gameId == 1001 {
				//放到類型房的4小間房
				roomCount = 4
			}

			for roomNum := 1; roomNum <= roomCount; roomNum++ {
				sum += data.RobotNum
			}
		}
	}

	aliveChannel = make(chan interface{}, sum/4)

	go ReceiveAlive()
	go CheckRobotAlive()

}

func ReceiveAlive() {
	for {
		robotData := <-aliveChannel
		if data, ok := robotData.(robot.RobotAlive); ok {
			robotAlive.Store(data.RobotName, data)
		}
	}
}

func CheckRobotAlive() {
	ticker := time.NewTicker(10 * time.Second)
	for range ticker.C {
		var count int
		robotAlive.Range(func(k, v interface{}) bool {
			data, ok := v.(robot.RobotAlive)
			if !ok {
				return true
			}
			number, _ := strconv.Atoi(strings.Replace(data.RobotName, robotName, "", -1))
			number -= robotStartNum

			if time.Now().Unix()-data.UpdateTime < ALIVE_SECONDS {
				count++

			}
			return true
		})
		fmt.Printf("機器人數量:%d 更新時間:%s\n", count, time.Now().Format("2006/01/02 15:04:05"))
	}
}

func startRobot(gameId int, tableID int, startNum int, robotNum int, playCount int) {
	for i := 0; i < robotNum; i++ {
		name := fmt.Sprintf("%s%04d", robotName, startNum+i)
		if ok, err := runRobot(gameId, tableID, name, playCount); !ok {
			SetMessage(utils.LOG_ERROR, err)
		} else {
			gameName, roomName := constant.GetGameRoomName(gameId, tableID)
			message := fmt.Sprintf("%s加入「%s」的「%s」(%d)進行測試", name, gameName, roomName, tableID)
			if playCount > 0 {
				message += fmt.Sprintf("，共玩%d局", playCount)
			}
			SetMessage(utils.LOG_INFO, message)

			//設資料到alive的map
			var initData = robot.RobotAlive{
				RobotName:  name,
				Fsm:        "",
				UpdateTime: time.Now().Unix(),
			}
			robotAlive.Store(name, initData)
		}
		time.Sleep(time.Duration(robotIntervalSec) * time.Second)
	}
}

func runRobot(gameId int, tableId int, robotName string, playCount int) (bool, string) {
	setting := robot.RobotConfig{
		Env:            robotEnv,
		AccessDcc:      false,
		AgentName:      "",
		LoginName:      robotName,
		GameId:         gameId,
		TableId:        tableId,
		MessageChannel: messageChannel,
		AliveChannel:   aliveChannel,
		Count:          playCount,
	}
	switch gameId {
	case 1001:
		robot := baccarat.NewRobot(setting)
		go robot.StartRobot()
	case 1002:
		robot := fantan.NewRobot(setting)
		go robot.StartRobot()
	case 1003:
		robot := colordisc.NewRobot(setting)
		go robot.StartRobot()
	case 1004:
		robot := prawncrab.NewRobot(setting)
		go robot.StartRobot()
	case 1005:
		robot := hundredsicbo.NewRobot(setting)
		go robot.StartRobot()
	case 1006:
		robot := cockfight.NewRobot(setting)
		go robot.StartRobot()
	case 1007:
		robot := dogracing.NewRobot(setting)
		go robot.StartRobot()
	case 1008:
		robot := rocket.NewRobot(setting)
		go robot.StartRobot()
	case 1009:
		robot := andarbahar.NewRobot(setting)
		go robot.StartRobot()
	case 2001:
		robot := blackjack.NewRobot(setting)
		go robot.StartRobot()
	case 2002:
		robot := sangong.NewRobot(setting)
		go robot.StartRobot()
	case 2003:
		robot := bullbull.NewRobot(setting)
		go robot.StartRobot()
	case 2004:
		robot := texas.NewRobot(setting)
		go robot.StartRobot()
	case 2005:
		robot := rummy.NewRobot(setting)
		go robot.StartRobot()
	case 2006:
		robot := goldenflower.NewRobot(setting)
		go robot.StartRobot()
	case 2007:
		robot := pokdeng.NewRobot(setting)
		go robot.StartRobot()
	case 2008:
		robot := catte.NewRobot(setting)
		go robot.StartRobot()
	case 2009:
		robot := chinesepoker.NewRobot(setting)
		go robot.StartRobot()
	case 2010:
		robot := okey.NewRobot(setting)
		go robot.StartRobot()
	case 2011:
		robot := teenpatti.NewRobot(setting)
		go robot.StartRobot()
	case 3001:
		robot := fruitslot.NewRobot(setting)
		go robot.StartRobot()
	case 3003:
		robot := plinko.NewEnhancedRobot(setting)
		go robot.StartRobot()
	case 4001:
		robot := fruit777slot.NewRobot(setting)
		go robot.StartRobot()
	case 4002:
		robot := megsharkslot.NewRobot(setting)
		go robot.StartRobot()
	case 4003:
		robot := midasslot.NewRobot(setting)
		go robot.StartRobot()
	default:
		err := fmt.Sprintf("%s啟動失敗GameId=%d,TableId=%d", robotName, gameId, tableId)
		return false, err
	}
	return true, ""
}

func SetMessage(logLevel int, message string) {
	utils.LogInfo(logLevel, message, messageChannel)
}
