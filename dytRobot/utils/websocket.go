package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type RespBase struct {
	Code    int         `json:"code,omitempty"`
	Message string      `json:"message,omitempty"`
	Ret     string      `json:"Ret,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var webData sync.Mutex

func SendMessage(connect *websocket.Conn, loginName string, data interface{}, channel chan string) (bool, error) {
	time.Sleep(1 * time.Second)
	webData.Lock()
	defer webData.Unlock()

	byteMsg := Package(data)
	err := connect.WriteMessage(websocket.TextMessage, byteMsg)
	if nil != err {
		LogInfo(LOG_INFO, fmt.Sprintf("%s SendMessage=%+v", loginName, data), channel)
		LogInfo(LOG_ERROR, loginName+" "+err.Error(), channel)
		return false, err
	}
	LogInfo(LOG_DEBUG, fmt.Sprintf("%s SendMessage=%+v", loginName, data), channel)
	return true, nil
}

func SendCustomizeMessage(connect *websocket.Conn, loginName string, data string, channel chan string) (bool, error) {
	time.Sleep(1 * time.Second)

	var m interface{}
	err := json.Unmarshal([]byte(data), &m)
	if err != nil {
		LogInfo(LOG_ERROR, loginName+" 不是json的格式:"+data+" Error message:"+err.Error(), channel)
		return false, err
	}

	webData.Lock()
	defer webData.Unlock()
	b := []byte(data)
	byteMsg := BytePackage(b)
	err = connect.WriteMessage(websocket.TextMessage, byteMsg)
	if nil != err {
		LogInfo(LOG_INFO, fmt.Sprintf("%s SendMessage=%+v", loginName, data), channel)
		LogInfo(LOG_ERROR, loginName+" "+err.Error(), channel)
		return false, err
	} else {
		LogInfo(LOG_INFO, fmt.Sprintf("%s SendMessage=%+v", loginName, data), channel)
		return true, nil
	}
}

func Unpackage(data []byte) ([]byte, error) {
	m := data[1:]
	dm, err := base64.StdEncoding.DecodeString(string(m))
	if err != nil {
		return []byte{}, err
	}
	//fmt.Println(string(dm))
	return dm, nil
}

func Package(data interface{}) []byte {
	m, _ := json.Marshal(data)
	bm := base64.StdEncoding.EncodeToString(m)
	//fmt.Println(bm)
	abm := []byte("a")
	abm = append(abm, []byte(bm)...)
	//fmt.Println(string(abm))
	return abm
}

func BytePackage(data []byte) []byte {
	bm := base64.StdEncoding.EncodeToString(data)
	//fmt.Println(bm)
	abm := []byte("a")
	abm = append(abm, []byte(bm)...)
	//fmt.Println(string(abm))
	return abm
}
