package slot_4006PryTreasure

import (
	"slotEditor/utils/random"
)

// Games Game Structure
type BGames struct {
	bonus_num *WeightGames //每層顆數權重
	bonus_pt  *WeightGames //每顆分數權重
	bonus_GU  *WeightGames //上或停
	Win       []float64    //中獎符號 [中獎]
	UpOver    []int        //[0=GG/1=UP]
	reelDef   int          //盤面大小
}

// NewBonusGame - 建立 Bonus Game
func NewBonusGame(num *WeightGames, pt *WeightGames, GU *WeightGames, def uint) *BGames {
	return &BGames{
		bonus_num: num,
		bonus_pt:  pt,
		bonus_GU:  GU,
		reelDef:   int(def),
	}
}
func (b *BGames) SpinBonus(level int) ([]float64, []int) {
	// 中獎顆數
	dice := random.Intn(b.bonus_num.Sum())
	pick, _ := b.bonus_num.Pick(dice)

	//中獎分數
	for i := 0; i < pick; i++ {
		dice := random.Intn(b.bonus_pt.Sum())
		pick, _ := b.bonus_pt.Pick(dice)
		winpt, _ := BgPayTable[pick].Float64()
		b.Win = append(b.Win, winpt)
	}
	var next int
	// 向上或結束
	if level != 5 || (level == 5 && pick == 0) {
		dice_2 := random.Intn(b.bonus_GU.Sum())
		pick_2, _ := b.bonus_GU.Pick(dice_2)
		if pick_2 == 300 {
			next = 0
		} else if pick_2 == 301 {
			next = 1
		}
		b.UpOver = append(b.UpOver, next)
	}
	return b.Win, b.UpOver
}
