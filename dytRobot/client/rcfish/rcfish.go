package rcfish

import (
	"dytRobot/client"
	"dytRobot/utils"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	ACTION_SHOOT = 1
	ACTION_HIT   = 2
)

type RcfishClient struct {
	*client.BaseElecClient
	entryBulletKind, entryBulletId, entrySeatId, entryHit, entrySkillID, entryTarget, entryDebugFishType, entryReward *widget.Entry
}

func NewClient(setting client.ClientConfig) *RcfishClient {
	elecClient := client.NewElecClient(setting)
	t := &RcfishClient{
		BaseElecClient: elecClient,
	}

	t.CheckResponse = t.CheckRcfishResponse
	return t
}

func (t *RcfishClient) CreateSection() fyne.CanvasObject {
	c := t.CreateTopSection()
	t.CreateElecSection(c)
	t.CreateFishControlSection(c)
	t.CreateBottomSection(c)
	return c
}

func (t *RcfishClient) CheckRcfishResponse(response *utils.RespBase) bool {
	return t.CheckElecResponse(response)
}

func (t *RcfishClient) CreateFishControlSection(c *fyne.Container) {
	labelBulletKind := widget.NewLabel("子彈種類")
	t.entryBulletKind = widget.NewEntry()
	t.entryBulletKind.SetText("1")

	labelBulletId := widget.NewLabel("子彈Id")
	t.entryBulletId = widget.NewEntry()
	t.entryBulletId.SetText("1")

	labelSeatId := widget.NewLabel("座位Id")
	t.entrySeatId = widget.NewEntry()
	t.entrySeatId.SetText("0")

	labelHitId := widget.NewLabel("擊中魚Index")
	t.entryHit = widget.NewEntry()
	t.entryHit.SetText("1")

	labelSkillId := widget.NewLabel("技能ID")
	t.entrySkillID = widget.NewEntry()
	t.entrySkillID.SetText("s1")

	labelSkillTarget := widget.NewLabel("技能目標")
	t.entryTarget = widget.NewEntry()
	t.entryTarget.SetText("1")

	labelDebugFishType := widget.NewLabel("指定出魚")
	t.entryDebugFishType = widget.NewEntry()
	t.entryDebugFishType.SetText("1")

	labelReward := widget.NewLabel("指定獎勵")
	t.entryReward = widget.NewEntry()
	t.entryReward.SetText("1")

	buttonShoot := widget.NewButton("射擊", func() {
		t.Shoot()
	})

	buttonHit := widget.NewButton("擊中", func() {
		t.Hit()
	})

	buttonSkillCast := widget.NewButton("使用技能", func() {
		t.SkillCast()
	})

	buttonSkillShoot := widget.NewButton("發射技能", func() {
		t.SkillShoot()
	})

	buttonDebugInfo := widget.NewButton("發送", func() {
		t.DebugInfo()
	})

	section1 := container.NewHBox(labelSeatId, t.entrySeatId, labelBulletKind, t.entryBulletKind, labelBulletId, t.entryBulletId, buttonShoot, labelHitId, t.entryHit, buttonHit)
	c.Add(section1)
	section2 := container.NewHBox(labelSkillId, t.entrySkillID, labelSkillTarget, t.entryTarget, buttonSkillShoot, buttonSkillCast)
	c.Add(section2)
	section3 := container.NewHBox(labelDebugFishType, t.entryDebugFishType, labelReward, t.entryReward, buttonDebugInfo)
	c.Add(section3)
}

func (t *RcfishClient) Shoot() (bool, error) {
	var data struct {
		Shoot struct {
			BulletKind int     `json:"bulletkind"`
			Angle      float64 `json:"angle"`
			LockFish   int     `json:"lockFish"`
		}
	}
	bulletKind, _ := strconv.Atoi(t.entryBulletKind.Text)
	data.Shoot.BulletKind = bulletKind
	return t.SendMessage(data)
}

func (t *RcfishClient) Hit() (bool, error) {
	var data struct {
		Hit struct {
			BulletId   int64 `json:"bulletId"`
			BulletKind int   `json:"bulletkind"`
			Fish       int   `json:"fish"`
			SeatId     int   `json:"seatId"`
		}
	}
	fish, _ := strconv.Atoi(t.entryHit.Text)
	data.Hit.Fish = fish
	bid, _ := strconv.Atoi(t.entryBulletId.Text)
	data.Hit.BulletId = int64(bid)
	seatId, _ := strconv.Atoi(t.entrySeatId.Text)
	data.Hit.SeatId = seatId
	return t.SendMessage(data)
}

func (t *RcfishClient) SkillShoot() (bool, error) {
	var data struct {
		SkillShoot struct {
			SkillIndex int    `json:"skillIndex"`
			SkillId    string `json:"skillId"`
		}
	}
	data.SkillShoot.SkillId = t.entrySkillID.Text
	data.SkillShoot.SkillIndex = 1
	return t.SendMessage(data)
}

func (t *RcfishClient) SkillCast() (bool, error) {
	var data struct {
		SkillCast struct {
			SkillIndex int64   `json:"skillIndex"`
			SkillId    string  `json:"skillId"`
			TargetFish []int64 `json:"targetfish"`
		}
	}
	data.SkillCast.SkillId = t.entrySkillID.Text
	data.SkillCast.SkillIndex = 1
	tar := strings.Split(t.entryTarget.Text, ",")
	var target []int64
	for _, fish := range tar {
		i, _ := strconv.Atoi(fish)
		target = append(target, int64(i))
	}
	data.SkillCast.TargetFish = target
	return t.SendMessage(data)
}

func (t *RcfishClient) DebugInfo() (bool, error) {
	var data struct {
		DebugInfo struct {
			Data interface{} `json:"data"`
		}
	}

	fishtype, _ := strconv.Atoi(t.entryDebugFishType.Text)
	da := map[string]interface{}{
		"fishtype": fishtype,
		"reward":   t.entryReward.Text,
	}
	data.DebugInfo.Data = da
	return t.SendMessage(data)
}
