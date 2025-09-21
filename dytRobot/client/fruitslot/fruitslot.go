package fruitslot

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

	ACTION_BET         = 1
	ACTION_GUESS_DOWN  = 4
	ACTION_GUESS_SMALL = 5
	ACTION_GUESS_BIG   = 6
	ACTION_GUESS_TIME  = 7

	ACT_GUESS_WARNING = "ActGuessWarning"
	ACT_GUESS_DOWN    = "ActGuessDown"
)

// 單獨獎項
const (
	SLOT_RESULT_APPLEx3      int = iota // 蘋果x3
	SLOT_RESULT_APPLE                   // 蘋果
	SLOT_RESULT_LEMONx3                 // 檸檬x3
	SLOT_RESULT_LEMON                   // 檸檬
	SLOT_RESULT_GRAPEx3                 // 葡萄x3
	SLOT_RESULT_GRAPE                   // 葡萄
	SLOT_RESULT_BELLx3                  // 鈴鐺x3
	SLOT_RESULT_BELL                    // 鈴鐺
	SLOT_RESULT_MANGOSTEENx3            // 山竹x3
	SLOT_RESULT_MANGOSTEEN              // 山竹
	SLOT_RESULT_DIAMONDx3               // 鑽石x3
	SLOT_RESULT_DIAMOND                 // 鑽石
	SLOT_RESULT_DOUBLE7x3               // 雙7x3
	SLOT_RESULT_DOUBLE7                 // 雙7
	SLOT_RESULT_CROWNx50                // 皇冠x50
	SLOT_RESULT_CROWN                   // 皇冠
	SLOT_RESULT_LUCKY                   // 幸運四葉

	SLOT_RESULT_ICON_COUNT
)

const (
	SLOT_SP_RESULT_NONE        int = iota // 隨機0燈
	SLOT_SP_RESULT_ONE                    // 隨機1燈
	SLOT_SP_RESULT_TWO                    // 隨機2燈
	SLOT_SP_RESULT_THREE                  // 隨機3燈
	SLOT_SP_RESULT_FOUR                   // 隨機4燈
	SLOT_SP_RESULT_SMALL_THREE            // 小三元獎
	SLOT_SP_RESULT_BIG_THREE              // 大三元獎
	SLOT_SP_RESULT_BIG_FOUR               // 大四喜獎
	SLOT_SP_RESULT_OCEAN_UP               // 縱橫四海上排
	SLOT_SP_RESULT_OCEAN_DOWN             // 縱橫四海下排
	SLOT_SP_RESULT_OCEAN_LEFT             // 縱橫四海左排
	SLOT_SP_RESULT_OCEAN_RIGHT            // 縱橫四海右排

	SLOT_SP_RESULT_COUNT
)

var (
	SlotIconName = []string{
		"蘋果x3",
		"蘋果",
		"檸檬x3",
		"檸檬",
		"葡萄x3",
		"葡萄",
		"鈴鐺x3",
		"鈴鐺",
		"山竹x3",
		"山竹",
		"鑽石x3",
		"鑽石",
		"雙7x3",
		"雙7",
		"皇冠x3",
		"皇冠",
	}

	SlotResultName = []string{
		"皇冠_1",
		"蘋果_1",
		"蘋果x3_1",
		"葡萄_1",
		"山竹_1",
		"山竹x3_1",
		"幸運四葉_1",
		"蘋果_2",
		"檸檬x3_1",
		"檸檬_1",
		"鈴鐺_1",
		"雙7x3_1",
		"雙7_1",
		"蘋果_3",
		"葡萄x3_1",
		"葡萄_2",
		"鑽石_1",
		"鑽石x3_1",
		"幸運四葉_2",
		"蘋果_4",
		"鈴鐺x3_1",
		"檸檬_2",
		"鈴鐺_2",
		"皇冠x50_1",
	}

	SlotSpResultName = []string{
		"隨機0燈",
		"隨機1燈",
		"隨機2燈",
		"隨機3燈",
		"隨機4燈",
		"小三元獎",
		"大三元獎",
		"大四喜獎",
		"縱橫四海上排",
		"縱橫四海下排",
		"縱橫四海左排",
		"縱橫四海右排",
	}

	WinRate = []int{
		3,   //蘋果x3
		5,   //蘋果
		3,   //檸檬x3
		10,  //檸檬
		3,   //葡萄x3
		15,  //葡萄
		3,   //鈴鐺x3
		20,  //鈴鐺
		3,   //山竹x3
		20,  //山竹
		3,   //鑽石x3
		30,  //鑽石
		3,   //雙7x3
		40,  //雙7
		50,  //皇冠x50
		100, //皇冠
		0,   //幸運四葉
	} // 單一圖標倍率

	WinSlotResult = []int{
		SLOT_RESULT_CROWN,        //皇冠_1
		SLOT_RESULT_APPLE,        //蘋果_1
		SLOT_RESULT_APPLEx3,      //蘋果x3_1
		SLOT_RESULT_GRAPE,        //葡萄_1
		SLOT_RESULT_MANGOSTEEN,   //山竹_1
		SLOT_RESULT_MANGOSTEENx3, //山竹x3_1
		SLOT_RESULT_LUCKY,        //幸運四葉_1
		SLOT_RESULT_APPLE,        //蘋果_2
		SLOT_RESULT_LEMONx3,      //檸檬x3_1
		SLOT_RESULT_LEMON,        //檸檬_1
		SLOT_RESULT_BELL,         //鈴鐺_1
		SLOT_RESULT_DOUBLE7x3,    //雙7x3_1
		SLOT_RESULT_DOUBLE7,      //雙7_1
		SLOT_RESULT_APPLE,        //蘋果_3
		SLOT_RESULT_GRAPEx3,      //葡萄x3_1
		SLOT_RESULT_GRAPE,        //葡萄_2
		SLOT_RESULT_DIAMOND,      //鑽石_1
		SLOT_RESULT_DIAMONDx3,    //鑽石x3_1
		SLOT_RESULT_LUCKY,        //幸運四葉_2
		SLOT_RESULT_APPLE,        //蘋果_4
		SLOT_RESULT_BELLx3,       //鈴鐺x3_1
		SLOT_RESULT_LEMON,        //檸檬_2
		SLOT_RESULT_BELL,         //鈴鐺_2
		SLOT_RESULT_CROWNx50,     //皇冠x50_1
	} // 單一圖標倍率

	WinSmallThreeIcon = []int{
		SLOT_RESULT_BELL,  //鈴鐺"
		SLOT_RESULT_GRAPE, //葡萄"
		SLOT_RESULT_LEMON, //檸檬"
	} // 小三元圖標

	WinBigThreeIcon = []int{
		SLOT_RESULT_DOUBLE7,    //雙7"
		SLOT_RESULT_DIAMOND,    //鑽石"
		SLOT_RESULT_MANGOSTEEN, //山竹"
	} // 大三元圖標

	WinBigFourIcon = []int{
		SLOT_RESULT_APPLE, //蘋果"
		SLOT_RESULT_APPLE, //蘋果"
		SLOT_RESULT_APPLE, //蘋果"
		SLOT_RESULT_APPLE, //蘋果"
	} // 大四喜圖標

	WinOceanUpIcon = []int{
		SLOT_RESULT_BELLx3,     //鈴鐺x3
		SLOT_RESULT_LEMON,      //檸檬
		SLOT_RESULT_BELL,       //鈴鐺
		SLOT_RESULT_CROWNx50,   //皇冠x50
		SLOT_RESULT_CROWN,      //皇冠
		SLOT_RESULT_APPLE,      //蘋果
		SLOT_RESULT_APPLEx3,    //蘋果x3
		SLOT_RESULT_GRAPE,      //葡萄
		SLOT_RESULT_MANGOSTEEN, //山竹
	} //縱橫四海上排

	WinOceanRightIcon = []int{
		SLOT_RESULT_MANGOSTEEN,   //山竹_1
		SLOT_RESULT_MANGOSTEENx3, //山竹x3_1
		SLOT_RESULT_APPLE,        //蘋果_2
		SLOT_RESULT_LEMONx3,      //檸檬x3
	} //縱橫四海右排

	WinOceanDownIcon = []int{
		SLOT_RESULT_LEMONx3,   //檸檬x3
		SLOT_RESULT_LEMON,     //檸檬
		SLOT_RESULT_BELL,      //鈴鐺
		SLOT_RESULT_DOUBLE7x3, //雙7x3
		SLOT_RESULT_DOUBLE7,   //雙7
		SLOT_RESULT_APPLE,     //蘋果
		SLOT_RESULT_GRAPEx3,   //葡萄x3
		SLOT_RESULT_GRAPE,     //葡萄
		SLOT_RESULT_DIAMOND,   //鑽石
	} //縱橫四海下排

	WinOceanLeftIcon = []int{
		SLOT_RESULT_DIAMOND,   //鑽石
		SLOT_RESULT_DIAMONDx3, //鑽石x3
		SLOT_RESULT_APPLE,     //蘋果
		SLOT_RESULT_BELLx3,    //鈴鐺x3
	} //縱橫四海左排
)

type FruitslotClient struct {
	*client.BaseElecClient
	entryApple, entryLemon, entryGrape       *widget.Entry
	entryBell, entryMangosteen, entryDiamond *widget.Entry
	entryDouble7, entryCrown                 *widget.Entry
	labelBaseBet                             *widget.Label
	selectGuess                              *widget.Select
}

func NewClient(setting client.ClientConfig) *FruitslotClient {
	elecClient := client.NewElecClient(setting)
	t := &FruitslotClient{
		BaseElecClient: elecClient,
	}

	t.CheckResponse = t.CheckFruitslotResponse

	t.CustomMessage = append(t.CustomMessage, "{\"PlayerAction\":{\"Action\":1,\"Data\":{\"BetInfo\":[0,0,5,10,5,15,20,0]}}}")
	t.CustomMessage = append(t.CustomMessage, "{\"PlayerAction\":{\"Action\":4}}")
	t.CustomMessage = append(t.CustomMessage, "{\"PlayerAction\":{\"Action\":5,\"Data\":{\"BetInfo\":1}}}")
	t.CustomMessage = append(t.CustomMessage, "{\"PlayerAction\":{\"Action\":6,\"Data\":{\"BetInfo\":1}}}")
	t.CustomMessage = append(t.CustomMessage, "{\"DebugInfo\":{\"Data\":{\"price\":[6,2,0,23]}}}")
	t.EntrySendMessage.SetOptions(t.CustomMessage)
	return t
}

func (t *FruitslotClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateElecSection(c)
	t.CreateFruitSlotSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *FruitslotClient) CheckFruitslotResponse(response *utils.RespBase) bool {
	if t.CheckElecResponse(response) {
		return true
	}
	switch response.Ret {
	case client.RET_INTO_GAME:
		t.AddTableStatus(t.GetJoinInfo(response))
		return true
	case client.ACT_GOLD:
		t.AddTableStatus(t.GetGoldInfo(response))
		return true
	case client.RET_PLAYER_ACTION:
		t.AddTableStatus(t.GetActionInfo(response))
		return true
	case ACT_GUESS_WARNING:
		t.AddTableStatus(t.GetGuessWarning(response))
		return true
	case ACT_GUESS_DOWN:
		t.AddTableStatus(t.GetGuessDownInfo(response))
		return true
	}
	return false
}

func (t *FruitslotClient) CreateFruitSlotSection(c *fyne.Container) {
	//水果機按鈕
	t.labelBaseBet = widget.NewLabel("底注:")

	labelApple := widget.NewLabel("蘋果")
	t.entryApple = widget.NewEntry()
	t.entryApple.SetText("10")

	labelLemon := widget.NewLabel("檸檬")
	t.entryLemon = widget.NewEntry()
	t.entryLemon.SetText("10")

	labelGrape := widget.NewLabel("葡萄")
	t.entryGrape = widget.NewEntry()
	t.entryGrape.SetText("10")

	labelBell := widget.NewLabel("鈴鐺")
	t.entryBell = widget.NewEntry()
	t.entryBell.SetText("10")

	labelMangosteen := widget.NewLabel("山竹")
	t.entryMangosteen = widget.NewEntry()
	t.entryMangosteen.SetText("10")

	labelDiamond := widget.NewLabel("鑽石")
	t.entryDiamond = widget.NewEntry()
	t.entryDiamond.SetText("10")

	labelDouble7 := widget.NewLabel("雙7")
	t.entryDouble7 = widget.NewEntry()
	t.entryDouble7.SetText("10")

	labelCrown := widget.NewLabel("皇冠")
	t.entryCrown = widget.NewEntry()
	t.entryCrown.SetText("10")

	buttonBet := widget.NewButton("押注", func() {
		t.SendPlayerAction(ACTION_BET)
	})

	section1 := container.NewHBox(t.labelBaseBet, labelApple, t.entryApple, labelLemon, t.entryLemon, labelGrape, t.entryGrape, labelBell, t.entryBell,
		labelMangosteen, t.entryMangosteen, labelDiamond, t.entryDiamond, labelDouble7, t.entryDouble7, labelCrown, t.entryCrown, buttonBet)
	c.Add(section1)

	t.selectGuess = widget.NewSelect([]string{}, nil)

	buttonGuessDown := widget.NewButton("猜:下分", func() {
		t.SendPlayerAction(ACTION_GUESS_DOWN)
	})

	buttonGuessSmall := widget.NewButton("猜:小1-7", func() {
		t.SendPlayerAction(ACTION_GUESS_SMALL)
	})

	buttonGuessBig := widget.NewButton("猜:大8-14", func() {
		t.SendPlayerAction(ACTION_GUESS_BIG)
	})

	buttonGuessTime := widget.NewButton("重設時間", func() {
		t.SendPlayerAction(ACTION_GUESS_TIME)
	})

	section2 := container.NewHBox(t.selectGuess, buttonGuessSmall, buttonGuessBig, buttonGuessDown, buttonGuessTime)
	c.Add(section2)

	labelDebugPriceType := widget.NewLabel("自訂中獎")
	comboDebugPriceType := widget.NewSelect(SlotResultName, nil)
	labelDebugSpPriceType := widget.NewLabel("自訂特別獎")
	comboDebugSpPriceType := widget.NewSelect(SlotSpResultName, nil)
	section3 := container.NewHBox(labelDebugPriceType, comboDebugPriceType, labelDebugSpPriceType, comboDebugSpPriceType)
	c.Add(section3)

	labelDebugSpPrice := widget.NewLabel("自訂隨機獎")
	comboDebugSpPrice1 := widget.NewSelect(SlotResultName, nil)
	comboDebugSpPrice2 := widget.NewSelect(SlotResultName, nil)
	comboDebugSpPrice3 := widget.NewSelect(SlotResultName, nil)
	comboDebugSpPrice4 := widget.NewSelect(SlotResultName, nil)

	buttonDebug := widget.NewButton("自訂送出", func() {
		priceType := comboDebugPriceType.SelectedIndex()
		spPriceType := comboDebugSpPriceType.SelectedIndex()
		spPrice1 := comboDebugSpPrice1.SelectedIndex()
		spPrice2 := comboDebugSpPrice2.SelectedIndex()
		spPrice3 := comboDebugSpPrice3.SelectedIndex()
		spPrice4 := comboDebugSpPrice4.SelectedIndex()
		t.SendDebug(priceType, spPriceType, spPrice1, spPrice2, spPrice3, spPrice4)

	})

	section4 := container.NewHBox(labelDebugSpPrice, comboDebugSpPrice1, comboDebugSpPrice2, comboDebugSpPrice3, comboDebugSpPrice4, buttonDebug)
	c.Add(section4)

}

func (t *FruitslotClient) SendPlayerAction(action int) (bool, error) {
	var data struct {
		PlayerAction struct {
			Action int         `json:"action"`
			Data   interface{} `json:"data,omitempty"`
		}
	}

	data.PlayerAction.Action = action
	if action == ACTION_BET {
		var info struct {
			BetInfo []int `json:"BetInfo"`
		}

		apple, _ := strconv.Atoi(t.entryApple.Text)
		lemon, _ := strconv.Atoi(t.entryLemon.Text)
		grape, _ := strconv.Atoi(t.entryGrape.Text)
		bell, _ := strconv.Atoi(t.entryBell.Text)
		mangosteen, _ := strconv.Atoi(t.entryMangosteen.Text)
		diamond, _ := strconv.Atoi(t.entryDiamond.Text)
		double7, _ := strconv.Atoi(t.entryDouble7.Text)
		crown, _ := strconv.Atoi(t.entryCrown.Text)

		info.BetInfo = append(info.BetInfo, apple)
		info.BetInfo = append(info.BetInfo, lemon)
		info.BetInfo = append(info.BetInfo, grape)
		info.BetInfo = append(info.BetInfo, bell)
		info.BetInfo = append(info.BetInfo, mangosteen)
		info.BetInfo = append(info.BetInfo, diamond)
		info.BetInfo = append(info.BetInfo, double7)
		info.BetInfo = append(info.BetInfo, crown)

		data.PlayerAction.Data = info
	}

	if action == ACTION_GUESS_BIG || action == ACTION_GUESS_SMALL {
		var info struct {
			BetInfo int `json:"BetInfo"`
		}

		bet, _ := strconv.Atoi(t.selectGuess.Selected)
		info.BetInfo = bet
		data.PlayerAction.Data = info
	}

	return t.SendMessage(data)
}

func (t *FruitslotClient) SendDebug(priceType int, spPriceType int, spPrice1 int, spPrice2 int, spPrice3 int, spPrice4 int) {
	var data struct {
		DebugInfo struct {
			Data struct {
				Price []int `json:"Price"`
			}
		}
	}
	if priceType != -1 {
		data.DebugInfo.Data.Price = append(data.DebugInfo.Data.Price, priceType)
	}
	if spPriceType != -1 {
		data.DebugInfo.Data.Price = append(data.DebugInfo.Data.Price, spPriceType)
	}
	if spPrice1 != -1 {
		data.DebugInfo.Data.Price = append(data.DebugInfo.Data.Price, spPrice1)
	}
	if spPrice2 != -1 {
		data.DebugInfo.Data.Price = append(data.DebugInfo.Data.Price, spPrice2)
	}
	if spPrice3 != -1 {
		data.DebugInfo.Data.Price = append(data.DebugInfo.Data.Price, spPrice3)
	}
	if spPrice4 != -1 {
		data.DebugInfo.Data.Price = append(data.DebugInfo.Data.Price, spPrice4)
	}

	t.SendMessage(data)
}

func (t *FruitslotClient) GetJoinInfo(response *utils.RespBase) string {
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

	baseBet := int(data["BaseBet"].(float64))
	t.labelBaseBet.SetText("底注:" + strconv.Itoa(baseBet))

	var message string
	gold := data["Gold"].(float64)
	message += fmt.Sprintf("玩家金額:%.4f\n", gold)

	if priceType, ok := data["PriceType"].(float64); ok {
		price := int(priceType)
		message += fmt.Sprintf("獎項:%s(%d)\n", SlotResultName[price], WinRate[WinSlotResult[price]])
	}
	if SpPriceType, ok := data["SpPriceType"].(float64); ok {
		spPrice := int(SpPriceType)
		message += "Lucky:" + SlotSpResultName[spPrice] + "\n"
		message += t.GetSpPriceName(spPrice)
	}
	if SpPrice, ok := data["SpPrice"].([]interface{}); ok {
		if len(SpPrice) > 0 {
			message += "隨機獎項:"
			for _, value := range SpPrice {
				price := int(value.(float64))
				message += fmt.Sprintf("%s(%d) ", SlotResultName[price], WinRate[WinSlotResult[price]])
			}
			message += "\n"
		}
	}
	if num, ok := data["Num"].(float64); ok {
		message += fmt.Sprintf("玩家上一次數字:%d", int(num))
	}
	if isBig, ok := data["IsBig"].(bool); ok {
		big := "為大\n"
		if !isBig {
			big = "為小\n"
		}
		message += big
	}
	if win, ok := data["Win"].(float64); ok {
		message += fmt.Sprintf("累積獎金:%d\n", int(win))
	}

	if guessArea, ok := data["GuessArea"].([]interface{}); ok {
		if len(guessArea) > 0 {
			message += fmt.Sprintf("猜大小區域:%v\n", guessArea)
			var option []string
			for _, value := range guessArea {
				bet := int(value.(float64))
				option = append(option, strconv.Itoa(bet))
			}
			t.selectGuess.Options = option
			t.selectGuess.SetSelectedIndex(1)
			t.selectGuess.Refresh()
		}
	}

	if betArea, ok := data["BetArea"].([]interface{}); ok {
		if len(betArea) > 0 {
			message += fmt.Sprintf("押注區域:%v\n", betArea)
			for i, value := range betArea {
				bet := int(value.(float64))
				switch i {
				case 0:
					t.entryApple.SetText(strconv.Itoa(bet))
				case 1:
					t.entryLemon.SetText(strconv.Itoa(bet))
				case 2:
					t.entryGrape.SetText(strconv.Itoa(bet))
				case 3:
					t.entryBell.SetText(strconv.Itoa(bet))
				case 4:
					t.entryMangosteen.SetText(strconv.Itoa(bet))
				case 5:
					t.entryDiamond.SetText(strconv.Itoa(bet))
				case 6:
					t.entryDouble7.SetText(strconv.Itoa(bet))
				case 7:
					t.entryCrown.SetText(strconv.Itoa(bet))
				}
			}
		}
	}

	return message
}

func (t *FruitslotClient) GetGoldInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}

	gold := data["Gold"].(float64)
	return fmt.Sprintf("玩家金額:%.4f\n", gold)
}

func (t *FruitslotClient) GetActionInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}

	fsm := data["Fsm"].(string)
	labelMessage := "Fsm:" + fsm
	gid := data["Gid"].(string)
	labelMessage += " Gid:" + gid
	t.LabelFsm.SetText(labelMessage)

	action := int(data["Action"].(float64))
	var message string

	switch action {
	case ACTION_BET:
		t.SetTableStatusClear()

		priceType := int(data["PriceType"].(float64))
		message += fmt.Sprintf("獎項:%s(%d)\n", SlotResultName[priceType], WinRate[WinSlotResult[priceType]])
		if SpPriceType, ok := data["SpPriceType"].(float64); ok {
			spPrice := int(SpPriceType)
			message += "Lucky:" + SlotSpResultName[spPrice] + "\n"
			message += t.GetSpPriceName(spPrice)
		}
		if SpPrice, ok := data["SpPrice"].([]interface{}); ok {
			if len(SpPrice) > 0 {
				message += "隨機獎項:"
				for _, value := range SpPrice {
					price := int(value.(float64))
					message += fmt.Sprintf("%s(%d) ", SlotResultName[price], WinRate[WinSlotResult[price]])
				}
				message += "\n"
			}
		}

		if guessArea, ok := data["GuessArea"].([]interface{}); ok {
			if len(guessArea) > 0 {
				var option []string
				message += "猜大小區域:"
				for _, value := range guessArea {
					bet := int(value.(float64))
					option = append(option, strconv.Itoa(bet))
				}
				t.selectGuess.Options = option
				t.selectGuess.SetSelectedIndex(1)
				t.selectGuess.Refresh()
				message += fmt.Sprintf("%v\n", option)
			}
		} else {
			t.selectGuess.Options = []string{}
			t.selectGuess.ClearSelected()
			t.selectGuess.Refresh()
		}

		win := int(data["Win"].(float64))
		message += fmt.Sprintf("累積獎金:%d\n", win)
	case ACTION_GUESS_DOWN:
		win := int(data["Win"].(float64))
		t.selectGuess.Options = []string{}
		t.selectGuess.ClearSelected()
		t.selectGuess.Refresh()
		message += fmt.Sprintf("玩家下分，累積獎金:%d\n", win)
	case ACTION_GUESS_SMALL, ACTION_GUESS_BIG:
		num := int(data["Num"].(float64))
		win := int(data["Win"].(float64))
		isBig := data["IsBig"].(bool)
		big := "為大"
		if !isBig {
			big = "為小"
		}
		guess := "猜大"
		if action == ACTION_GUESS_SMALL {
			guess = "猜小"
		}
		message += fmt.Sprintf("玩家%s，數字:%d%s,累積獎金:%d\n", guess, num, big, win)

		if guessArea, ok := data["GuessArea"].([]interface{}); ok {
			if len(guessArea) > 0 {
				var option []string
				message += "猜大小區域:"
				for _, value := range guessArea {
					bet := int(value.(float64))
					option = append(option, strconv.Itoa(bet))
				}
				t.selectGuess.Options = option
				t.selectGuess.SetSelectedIndex(1)
				t.selectGuess.Refresh()
				message += fmt.Sprintf("%v\n", option)
			}
		} else {
			t.selectGuess.Options = []string{}
			t.selectGuess.ClearSelected()
			t.selectGuess.Refresh()
		}
	}
	return message
}

func (t *FruitslotClient) GetGuessWarning(response *utils.RespBase) string {
	return "收到玩家未進行比大小警告\n"
}

func (t *FruitslotClient) GetGuessDownInfo(response *utils.RespBase) string {
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		return ""
	}

	fsm := data["Fsm"].(string)
	labelMessage := "Fsm:" + fsm
	gid := data["Gid"].(string)
	labelMessage += " Gid:" + gid
	t.LabelFsm.SetText(labelMessage)

	win := int(data["Win"].(float64))
	return fmt.Sprintf("自動下分，累積獎金:%d\n", win)
}

func (t *FruitslotClient) GetSpPriceName(spPrice int) string {
	if spPrice < SLOT_SP_RESULT_SMALL_THREE {
		return ""
	}

	message := "獎項:"
	switch spPrice {
	case SLOT_SP_RESULT_SMALL_THREE:
		for _, value := range WinSmallThreeIcon {
			message += fmt.Sprintf("%s(%d) ", SlotIconName[value], WinRate[value])
		}
	case SLOT_SP_RESULT_BIG_THREE:
		for _, value := range WinBigThreeIcon {
			message += fmt.Sprintf("%s(%d) ", SlotIconName[value], WinRate[value])
		}
	case SLOT_SP_RESULT_BIG_FOUR:
		for _, value := range WinBigFourIcon {
			message += fmt.Sprintf("%s(%d) ", SlotIconName[value], WinRate[value])
		}
	case SLOT_SP_RESULT_OCEAN_UP:
		for _, value := range WinOceanUpIcon {
			message += fmt.Sprintf("%s(%d) ", SlotIconName[value], WinRate[value])
		}
	case SLOT_SP_RESULT_OCEAN_RIGHT:
		for _, value := range WinOceanRightIcon {
			message += fmt.Sprintf("%s(%d) ", SlotIconName[value], WinRate[value])
		}
	case SLOT_SP_RESULT_OCEAN_DOWN:
		for _, value := range WinOceanDownIcon {
			message += fmt.Sprintf("%s(%d) ", SlotIconName[value], WinRate[value])
		}
	case SLOT_SP_RESULT_OCEAN_LEFT:
		for _, value := range WinOceanLeftIcon {
			message += fmt.Sprintf("%s(%d) ", SlotIconName[value], WinRate[value])
		}
	}
	message += "\n"

	return message
}
