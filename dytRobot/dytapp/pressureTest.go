//go:build !auto

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
	"dytRobot/robot/roulette"
	"dytRobot/robot/rummy"
	"dytRobot/robot/sangong"
	"dytRobot/robot/teenpatti"
	"dytRobot/robot/texas"
	"dytRobot/utils"
	"fmt"
	"image/color"
	"strconv"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var robotEnv int = constant.ENV_NONE
var accessENV string
var robotName string = "dayon"
var robotStartNum int = 1
var robotIntervalSec int = 3

var robotSttings map[int]*constant.GameSetting

var messageChannel chan string
var aliveChannel chan interface{}

var robotAlive sync.Map
var labelAlive *widget.Label
var labelAlive2 *widget.Label
var checkDCC *widget.Check
var robotList *widget.Table
var messageApp *widget.Entry
var comboAgent *widget.Select
var checkRooms *widget.Check

var tabsPressTest *container.AppTabs
var agentEnv int = constant.ENV_NONE

const (
	ALIVE_SECONDS = 30
)

func setPressureTestEnv() {
	initRobotSettings()

	messageApp = widget.NewMultiLineEntry()
	//畫面上顯示的行數
	messageApp.SetMinRowsVisible(25)

	startPage := makeStartPage()
	envConfigPage := makeNormalSettingPage()
	roomConfigPage := makeRoomSettingPage()
	runPage := makeExecutablePage()

	tabsPressTest = container.NewAppTabs()
	tabsPressTest.Append(container.NewTabItem("說明頁", startPage))
	tabsPressTest.Append(container.NewTabItem("基本設定", envConfigPage))
	tabsPressTest.Append(container.NewTabItem("遊戲設定", roomConfigPage))
	tabsPressTest.Append(container.NewTabItem("測試運行", runPage))
	tabsPressTest.Append(container.NewTabItem("機器人列表", widget.NewLabel("機器人數量:0 更新時間:")))
	tabsPressTest.SetTabLocation(container.TabLocationLeading)
}

func initRobotSettings() {
	robotSttings = make(map[int]*constant.GameSetting)
	for _, gameID := range constant.GameIDList {
		setting := &constant.GameSetting{}
		robotSttings[gameID] = setting
	}
}

func makeStartPage() fyne.CanvasObject {
	message := "歡迎使用遊戲壓測機器人\n" +
		"煩請遵照以下步驟設定\n" +
		"1.至基本設定頁:\n" +
		" (1)設定測試環境\n" +
		" (2)設置測試代理：點擊'讀取'得到可選取列表；未設置則用預設總代理\n" +
		" (3)設定機器人名稱及其初始編號，並會於下方顯示初始機器人名稱\n" +
		" (4)設定機器人登入間隔\n" +
		" (5)在DEV、QA、ETE、PROD皆可透過後台登入\n" +
		"2.至遊戲設定頁:\n" +
		" (1)設定所要測試之遊戲的房間及機器人數量，可多遊戲同時設定\n" +
		" (2)當選擇多房間時，所設的機器人數量亦為該房間機器人數量\n" +
		" (3)百家樂可多設定是否要把機器人分到該房間的4個小房間\n" +
		" (4)各遊戲可多設定機器人玩幾局就結束，設置0或不填為無限制\n" +
		"3.至測試運行頁\n" +
		" (1)點擊開始測試鈕開始進行壓力測試\n" +
		" (2)按鈕下方顯示訊息\n" +
		" (3)按鈕上方統計目前運行之機器人數量\n" +
		" (4)機器人在線情況可至機器人列表頁面確認\n" +
		"4.機器人動作\n" +
		" (1)機器人於百人場時，每4秒會隨機從其中一種籌碼下於任一個區域\n" +
		"     各區域下注機率會與賠率成反比\n" +
		" (2)機器人於配對場時，會隨機下注或動作\n" +
		" (3)機器人當金額不足時，便會自動進行上分\n" +
		" (4)機器人於水果機下注時，只會有下注幾個特定區域；猜大小會隨機加倍減半下分\n" +
		"5.測試結束\n" +
		" (1)請直接關閉程式\n" +
		"6.歷程檔\n" +
		" (1)操作紀錄與錯誤訊息皆會寫在log資料夾中\n"
	return widget.NewLabel(message)
}

func ClearAgentInfo() {
	comboAgent.Options = []string{}
	comboAgent.Selected = ""
	agentEnv = constant.ENV_NONE
	comboAgent.Refresh()
}

func makeNormalSettingPage() fyne.CanvasObject {
	label1 := widget.NewLabel("選擇測試環境")
	checkDCC = widget.NewCheck("透過後台登入", nil)
	checkDCC.Checked = true
	checkDCC.Disable()
	labelEnv := container.NewHBox(label1, checkDCC)

	label6 := widget.NewLabel("選擇登入代理(不選則使用預設)")
	label7 := widget.NewLabel("")
	comboAgent = widget.NewSelect([]string{}, func(value string) {
	})

	combo := widget.NewSelect([]string{"PROD", "ETE", "QA", "QC", "DEV", "LOCAL"}, func(value string) {
		setRobotEnv(value)
	})

	button1 := widget.NewButton("讀取列表", func() {
		if combo.Selected == "PROD" || combo.Selected == "ETE" || combo.Selected == "QA" || combo.Selected == "QC" || (combo.Selected == "DEV" && checkDCC.Checked) {
			errorMessage, list := utils.GetDccAgentInfo(constant.DccAgentURL[robotEnv])
			if errorMessage != "" {
				SetMessage(utils.LOG_ERROR, errorMessage)
				ClearAgentInfo()
				label7.SetText("讀取失敗")
			} else {
				comboAgent.Options = list
				agentEnv = robotEnv
				label7.SetText("讀取成功")
			}
		} else {
			ClearAgentInfo()
			label7.SetText("")
		}
	})

	button2 := widget.NewButton("清空列表", func() {
		ClearAgentInfo()
		label7.SetText("")
	})

	sectionAgent := container.NewHBox(label6, button1, button2, label7)

	label2 := widget.NewLabel("設置測試機器人基本資料")
	message := widget.NewLabel("")
	ShowRobotName(message)

	label3 := widget.NewLabel("機器人名稱")
	nameEntry := widget.NewEntry()
	nameEntry.SetText(robotName)
	nameEntry.OnChanged = func(content string) {
		robotName = content
		ShowRobotName(message)
	}

	label4 := widget.NewLabel("機器人起始編號")
	numEntry := widget.NewEntry()
	numEntry.SetText(strconv.Itoa(robotStartNum))
	numEntry.OnChanged = func(content string) {
		if content == "" {
			robotStartNum = 0
			numEntry.SetText(strconv.Itoa(robotStartNum))
		}
		num, err := strconv.Atoi(content)
		if err != nil {
			numEntry.SetText(strconv.Itoa(robotStartNum))
			return
		}
		robotStartNum = num
		ShowRobotName(message)
	}

	label5 := widget.NewLabel("機器人間隔秒數")
	intervalEntry := widget.NewEntry()
	intervalEntry.SetText(strconv.Itoa(robotIntervalSec))
	intervalEntry.OnChanged = func(content string) {
		if content == "" {
			robotIntervalSec = 0
			intervalEntry.SetText(strconv.Itoa(robotIntervalSec))
			return
		}
		num, err := strconv.Atoi(content)
		if err != nil {
			intervalEntry.SetText(strconv.Itoa(robotIntervalSec))
			return
		}
		robotIntervalSec = num
	}

	return container.NewVBox(labelEnv, combo, sectionAgent, comboAgent, label2, label3, nameEntry,
		label4, numEntry, label5, intervalEntry, message)
}

func setRobotEnv(value string) {
	switch value {
	case "LOCAL":
		robotEnv = constant.ENV_LOCAL
		checkDCC.Checked = false
		checkDCC.Disable()
		checkDCC.Refresh()
	case "DEV":
		robotEnv = constant.ENV_DEV
		checkDCC.Checked = true
		checkDCC.Enable()
		checkDCC.Refresh()
	case "QC":
		robotEnv = constant.ENV_QC
		checkDCC.Checked = true
		checkDCC.Enable()
		checkDCC.Refresh()
	case "QA":
		robotEnv = constant.ENV_QA
		checkDCC.Checked = true
		checkDCC.Disable()
		checkDCC.Refresh()
	case "ETE":
		robotEnv = constant.ENV_ETE
		checkDCC.Checked = true
		checkDCC.Disable()
		checkDCC.Refresh()
	case "PROD":
		robotEnv = constant.ENV_PROD
		checkDCC.Checked = true
		checkDCC.Disable()
		checkDCC.Refresh()
	}
	accessENV = value
}

func ShowRobotName(message *widget.Label) {
	message.SetText("測試機器人名稱：" + fmt.Sprintf("%s%04d", robotName, robotStartNum))
}

func makeRoomSettingPage() fyne.CanvasObject {
	accordion := widget.NewAccordion()

	for _, gameID := range constant.GameIDList {
		roomNum := constant.RoomTypeNum[gameID]
		item := widget.NewAccordionItem(constant.GameIDNameTable[gameID], nil)
		check01 := widget.NewCheck("開啟機器人", func(value bool) {
			gameID := getGameID(item.Title)
			if value {
				item.Title = constant.GameIDNameTable[gameID] + "(已開啟)"
				robotSttings[gameID].IsEnable = true
				accordion.Refresh()
			} else {
				item.Title = constant.GameIDNameTable[gameID]
				robotSttings[gameID].IsEnable = false
				accordion.Refresh()
			}
		})

		check02 := widget.NewCheck("新手房", func(value bool) {
			gameID := getGameID(item.Title)
			robotSttings[gameID].Room[0] = value
		})

		check03 := widget.NewCheck("普通房", func(value bool) {
			gameID := getGameID(item.Title)
			robotSttings[gameID].Room[1] = value
		})

		check04 := widget.NewCheck("高級房", func(value bool) {
			gameID := getGameID(item.Title)
			robotSttings[gameID].Room[2] = value
		})

		check05 := widget.NewCheck("大師房", func(value bool) {
			gameID := getGameID(item.Title)
			robotSttings[gameID].Room[3] = value
		})

		check06 := widget.NewCheck("初级场", func(value bool) {
			gameID := getGameID(item.Title)
			robotSttings[gameID].Room[4] = value
		})

		check07 := widget.NewCheck("中级场", func(value bool) {
			gameID := getGameID(item.Title)
			robotSttings[gameID].Room[5] = value
		})

		check08 := widget.NewCheck("高级场", func(value bool) {
			gameID := getGameID(item.Title)
			robotSttings[gameID].Room[6] = value
		})

		check09 := widget.NewCheck("至尊场", func(value bool) {
			gameID := getGameID(item.Title)
			robotSttings[gameID].Room[7] = value
		})

		check10 := widget.NewCheck("全部房間", func(value bool) {
			if value {
				check02.SetChecked(true)
				if roomNum > 2 {
					check03.SetChecked(true)
					if roomNum > 2 {
						check04.SetChecked(true)
						if roomNum > 3 {
							check05.SetChecked(true)
							if roomNum > 4 {
								check06.SetChecked(true)
								if roomNum > 5 {
									check07.SetChecked(true)
									if roomNum > 6 {
										check08.SetChecked(true)
										if roomNum > 7 {
											check09.SetChecked(true)
										}
									}
								}
							}
						}
					}
				}
			} else {
				check02.SetChecked(false)
				check03.SetChecked(false)
				check04.SetChecked(false)
				check05.SetChecked(false)
				check06.SetChecked(false)
				check07.SetChecked(false)
				check08.SetChecked(false)
				check09.SetChecked(false)
			}
		})

		var room *fyne.Container
		if roomNum == 4 {
			room = container.NewHBox(widget.NewLabel("選擇房間"), check10, check02, check03, check04, check05)
		} else if roomNum == 8 {
			room = container.NewHBox(widget.NewLabel("選擇房間"), check10, check02, check03, check04, check05, check06, check07, check08, check09)
		} else if roomNum == 3 {
			room = container.NewHBox(widget.NewLabel("選擇房間"), check10, check02, check03, check04)
		} else if roomNum == 1 {
			room = container.NewHBox(widget.NewLabel("選擇房間"), check10, check02)
		}

		robotNumEntry := widget.NewEntry()
		robotNumEntry.SetPlaceHolder("請輸入每個房間的機器人數量")

		robotNumEntry.OnChanged = func(content string) {
			if content == "" {
				return
			}
			gameID := getGameID(item.Title)
			num, err := strconv.Atoi(content)
			if err != nil {
				robotNumEntry.SetText(strconv.Itoa(robotSttings[gameID].RobotNum))
				return
			}
			robotSttings[gameID].RobotNum = num
		}

		check := container.NewHBox(check01)
		if gameID == 1001 {
			checkRooms = widget.NewCheck("是否分發至房間裡的4小房", nil)
			check.Add(checkRooms)
		}

		playCountEntry := widget.NewEntry()
		playCountEntry.SetPlaceHolder("請輸入機器人玩的場次，0或不填為無限制")

		playCountEntry.OnChanged = func(content string) {
			if content == "" {
				return
			}
			gameID := getGameID(item.Title)
			num, err := strconv.Atoi(content)
			if err != nil {
				robotNumEntry.SetText(strconv.Itoa(robotSttings[gameID].PlayCount))
				return
			}
			robotSttings[gameID].PlayCount = num
		}

		detail := container.NewVBox(check, room, robotNumEntry, playCountEntry)
		item.Detail = detail
		accordion.Append(item)
	}
	// 將 accordion 放入垂直滾動條容器中
	scrollContainer := container.NewVScroll(accordion)
	scrollContainer.SetMinSize(fyne.NewSize(300, 800)) // 這裡的寬度和高度只是示例，請根據實際需要調整
	label1 := widget.NewLabel("設置測試廳館")
	return container.NewVBox(label1, scrollContainer)
}

func getGameID(title string) int {
	return constant.NameTableGameID[strings.Replace(title, "(已開啟)", "", -1)]
}

func makeExecutablePage() fyne.CanvasObject {
	label1 := widget.NewLabel("執行壓力測試")
	labelAlive = widget.NewLabel("存活的機器人:0 更新時間:")
	button := widget.NewButton("開始測試", nil)
	button.OnTapped = func() {
		//開啟訊息頻道
		if messageChannel == nil {
			messageChannel = make(chan string, 50)
			go ReceiveMessage()
		}

		ClearMessage()
		SetMessage(utils.LOG_INFO, "開始執行壓力測試:"+accessENV)
		//檢查設定
		if error, ok := checkSetting(); !ok {
			SetMessage(utils.LOG_ERROR, error)
			return
		}
		button.Disable()

		//建立RobotList
		makeRobotListPage()

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
				if gameId == 1001 && checkRooms.Checked {
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
	return container.NewVBox(label1, labelAlive, button, messageApp)
}

func checkSetting() (string, bool) {
	var error string
	if robotEnv == constant.ENV_NONE {
		error += "未設置測試環境"
		return error, len(error) == 0
	}

	if agentEnv != constant.ENV_NONE && agentEnv != robotEnv {
		error += "測試環境與代理環境不一致"
		return error, len(error) == 0
	}

	var isSetGame bool

	for gameId, data := range robotSttings {
		if !data.IsEnable {
			continue
		}
		isSetGame = true

		var isSetRoom bool
		for _, room := range data.Room {
			if room {
				isSetRoom = room
			}
		}
		if !isSetRoom {
			error += fmt.Sprintf("%s已設置開啟，但未設定測試廳館。", constant.GameIDNameTable[gameId])
		}

		if data.RobotNum == 0 {
			error += fmt.Sprintf("%s已設置開啟，但未設定機器人數量。", constant.GameIDNameTable[gameId])
		}
	}

	if !isSetGame {
		error += "未設置任何遊戲"
	}

	return error, len(error) == 0
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

			//設資料到table
			num := startNum + i + 1 - robotStartNum
			robotStatus[num][1] = gameName
			robotStatus[num][2] = roomName
			robotStatus[num][3] = "Login"
			robotStatus[num][4] = "0"
			robotStatus[num][5] = "在線"

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
		AccessDcc:      checkDCC.Checked,
		AgentName:      comboAgent.Selected,
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
	case 1010:
		robot := roulette.NewRobot(setting)
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
func ClearMessage() {
	messageApp.SetText("")
}

func SetMessage(logLevel int, message string) {
	utils.LogInfo(logLevel, message, messageChannel)
}

func ReceiveMessage() {
	for {
		message := <-messageChannel
		content := messageApp.Text
		if len(content) > 5000 {
			messageApp.SetText(message)
			continue
		}
		messageApp.SetText(message + content)
	}
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
				robotStatus[number+1][3] = data.Fsm
				robotStatus[number+1][4] = strconv.Itoa(data.PlayCount + 1)
				robotStatus[number+1][5] = "在線"
			} else {
				robotStatus[number+1][3] = data.Fsm
				robotStatus[number+1][4] = strconv.Itoa(data.PlayCount + 1)
				robotStatus[number+1][5] = "離線"
			}

			return true
		})
		mseage := fmt.Sprintf("機器人數量:%d 更新時間:%s", count, time.Now().Format("2006/01/02 15:04:05"))
		labelAlive.SetText(mseage)
		labelAlive2.SetText(mseage)
		robotList.Refresh()
	}
}

var robotStatus = [][]string{{"機器人名稱", "遊戲名稱", "遊戲房間", "遊戲階段", "局", "在線/離線"}}

func makeRobotListPage() {
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
			if gameId == 1001 && checkRooms.Checked {
				//放到類型房的4小間房
				roomCount = 4
			}

			for roomNum := 1; roomNum <= roomCount; roomNum++ {
				sum += data.RobotNum
			}
		}
	}

	for i := robotStartNum; i < robotStartNum+sum; i++ {
		name := fmt.Sprintf("%s%04d", robotName, i)
		robotStatus = append(robotStatus, []string{name, "", "", "", "", "離線"})
	}

	aliveChannel = make(chan interface{}, sum/4)

	robotList = widget.NewTable(
		func() (int, int) {
			return sum + 1, len(robotStatus[0])
		},
		func() fyne.CanvasObject {
			return canvas.NewText("", color.White)
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*canvas.Text).Text = robotStatus[i.Row][i.Col]
			if i.Row == 0 {
				o.(*canvas.Text).Color = color.NRGBA{B: 255, G: 255, R: 0, A: 255}
			} else if robotStatus[i.Row][5] == "在線" && i.Col == 5 {
				o.(*canvas.Text).Color = color.NRGBA{B: 0, G: 255, R: 0, A: 255}
			} else if robotStatus[i.Row][5] == "離線" && i.Col == 5 {
				o.(*canvas.Text).Color = color.NRGBA{B: 0, G: 0, R: 255, A: 255}
			} else {
				o.(*canvas.Text).Color = color.White
			}
		})
	for i := 0; i < len(robotStatus[0]); i++ {
		robotList.SetColumnWidth(i, 100)
	}

	go ReceiveAlive()
	go CheckRobotAlive()

	labelAlive2 = widget.NewLabel("機器人數量:0 更新時間:")
	page := container.NewBorder(labelAlive2, nil, nil, nil, robotList)
	tabsPressTest.Items[len(tabsPressTest.Items)-1] = container.NewTabItem("機器人列表", page)
	tabsPressTest.Refresh()
}
