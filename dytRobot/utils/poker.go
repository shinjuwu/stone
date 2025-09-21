package utils

import (
	"fmt"
	"strconv"
)

var cardSpPoint = map[int]string{
	1:  "A",
	11: "J",
	12: "Q",
	13: "K",
}

var color []string = []string{"梅花", "紅磚", "紅桃", "黑桃"}

func PrintCard(data interface{}) string {
	if cards, ok := data.([]int); ok {
		return transCardString(cards)
	}
	gc := ""
	if groupCards, ok := data.([][]int); ok {
		for _, cards := range groupCards {
			gc += transCardString(cards) + "|"
		}
		return gc
	}
	return ""
}
func transCardString(cards []int) string {
	cs := ""
	for _, v := range cards {
		n := CardNumber(v)
		var ns string
		if n < 2 || n > 10 {
			ns = cardSpPoint[n]
		} else {
			ns = strconv.Itoa(n)
		}
		c := fmt.Sprintf("%s%s, ", color[GetColor(v)], ns)
		cs += c
	}
	return cs
}

func CardNumber(card int) int {
	value := card%13 + 1
	return value
}

func GetColor(card int) int {
	return card / 13
}
