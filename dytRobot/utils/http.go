package utils

import (
	"dytRobot/constant"
	"dytRobot/pkg/encrypt/aescbc"
	md5 "dytRobot/pkg/encrypt/md5hash"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	DCC_AGENTID  = 3
	HTTP_CODE_OK = 200
)

const (
	DCC_ACTION_LOGIN        = 0
	DCC_ACTION_DEPOSIT      = 2
	DCC_CHANNEL             = "channel/channelHandle?"
	DCC_REDIRECTION_KEY     = "token="
	DCC_REDIRECTION_KEY_END = "&lang="
	DCC_ASEKEY              = "ddxbst648uf7hdbc" //16bit
	DCC_MD5KEY              = "f7c637cd39679e99"
)

type DccResp struct {
	Code int         `json:"code,omitempty"`
	Msg  string      `json:"msg,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type DccAgentInfo struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Md5key string `json:"md5_key"`
	Aeskey string `json:"aes_key"`
}

var DccAgentList map[string]DccAgentInfo

func GameDeposit(depositUrl string, userId int64, gold int) (bool, error) {
	v := make(map[string][]string)
	v["method"] = []string{"up"}
	v["gold"] = []string{strconv.Itoa(gold)}
	v["userid"] = []string{strconv.FormatInt(userId, 10)}

	qs := url.Values(v)
	resp, err := http.PostForm(depositUrl, qs)
	if err != nil {
		return false, fmt.Errorf(fmt.Sprintf("Deposit failed: : %v", err))
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	response := &RespBase{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, fmt.Errorf("Deposit解析失敗")
	}

	if response.Code != constant.ERROR_CODE_SUCCESS {
		return false, fmt.Errorf(fmt.Sprintf("Deposit失敗.錯誤代碼=%d,錯誤訊息=%s", response.Code, response.Message))
	}
	return true, nil
}

func DccLogin(loginUrl string, agentName string, loginName string) (bool, string, error) {
	requestUrl := createRequestURL(loginUrl, agentName, DCC_ACTION_LOGIN, loginName, 0)
	if requestUrl == "" {
		return false, "", fmt.Errorf("DCC 登入失敗:產生request url失敗")
	}

	resp, err := http.Get(requestUrl)
	if err != nil {
		return false, "", fmt.Errorf(fmt.Sprintf("DCC 登入失敗:%v", err))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != HTTP_CODE_OK {
		return false, "", fmt.Errorf(fmt.Sprintf("DCC http 登入失敗,status code=%d,錯誤訊息=%s", resp.StatusCode, string(body)))
	}

	response := &DccResp{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, "", fmt.Errorf("DCC 登入 json解析失敗,內容=" + string(body))
	}

	if response.Code != constant.ERROR_CODE_SUCCESS {
		return false, "", fmt.Errorf(fmt.Sprintf("DCC 登入失敗.錯誤代碼=%d,錯誤訊息=%s,內容=%s", response.Code, response.Msg, string(body)))
	}

	redirectionURL := response.Data.(map[string]interface{})["d"].(map[string]interface{})["url"].(string)
	if redirectionURL == "" {
		return false, "", fmt.Errorf("DCC 登入失敗,不存在轉址=" + string(body))
	}

	index := strings.Index(redirectionURL, DCC_REDIRECTION_KEY)
	endindex := strings.Index(redirectionURL, DCC_REDIRECTION_KEY_END)
	if index == -1 {
		return false, "", fmt.Errorf("DCC 登入失敗.無法找到token. URL:" + redirectionURL)
	}
	token := redirectionURL[index+len(DCC_REDIRECTION_KEY) : endindex]

	return true, token, nil
}

func DccDeposit(depositUrl string, agentName string, loginName string, gold int) (bool, error) {
	requestUrl := createRequestURL(depositUrl, agentName, DCC_ACTION_DEPOSIT, loginName, gold)
	if requestUrl == "" {
		return false, fmt.Errorf("DCC 上分失敗:產生request url失敗")
	}

	resp, err := http.Get(requestUrl)
	if err != nil {
		return false, fmt.Errorf(fmt.Sprintf("DCC上分失敗:%v", err))
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != HTTP_CODE_OK {
		return false, fmt.Errorf(fmt.Sprintf("DCC http 上分失敗,status code=%d,錯誤訊息=%s", resp.StatusCode, string(body)))
	}

	response := &DccResp{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, fmt.Errorf("DCC上分json解析失敗,內容=" + string(body))
	}

	if response.Code != constant.ERROR_CODE_SUCCESS {
		return false, fmt.Errorf(fmt.Sprintf("DCC上分失敗.錯誤代碼=%d,錯誤訊息=%s,內容=%s", response.Code, response.Msg, string(body)))
	}
	return true, nil
}

func createRequestURL(actionURL string, agentName string, dccAction int, loginName string, gold int) string {
	agentId := DCC_AGENTID
	aeskey := DCC_ASEKEY
	md5key := DCC_MD5KEY
	if agentName != "" {
		agentId = DccAgentList[agentName].ID
		aeskey = DccAgentList[agentName].Aeskey
		md5key = DccAgentList[agentName].Md5key
	}

	orderid := fmt.Sprintf("%d%s%s", agentId, time.Now().Format("20060102150405"), loginName)
	paramData := fmt.Sprintf("s=%d&account=%s&money=%d&orderid=%s&kind=0", dccAction, loginName, gold, orderid)

	timestamp := time.Now().UnixMilli()
	paramDataAesEncoding, err := aescbc.AesEncrypt([]byte(paramData), []byte(aeskey))
	if err != nil {
		return ""
	}

	b64Encoding := base64.StdEncoding.EncodeToString([]byte(paramDataAesEncoding))

	ss := strconv.Itoa(agentId) + strconv.FormatInt(timestamp, 10) + md5key
	s32 := md5.Hash32bit(ss)

	rawApiQuery := make(url.Values)
	rawApiQuery.Add("agent", strconv.Itoa(agentId))
	rawApiQuery.Add("timestamp", strconv.FormatInt(timestamp, 10))
	rawApiQuery.Add("param", b64Encoding)
	rawApiQuery.Add("key", s32)

	return actionURL + DCC_CHANNEL + rawApiQuery.Encode()
}

func GetDccAgentInfo(url string) (errorMessage string, agentList []string) {
	requestUrl := url + "api/v1/intercom/getagentlist"
	resp, err := http.Get(requestUrl)
	if err != nil {
		errorMessage = fmt.Sprintf("DCC獲取agent失敗:%v", err)
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != HTTP_CODE_OK {
		errorMessage = fmt.Sprintf("DCC http 獲取agent失敗,status code=%d,錯誤訊息=%s", resp.StatusCode, string(body))
		return
	}

	response := &DccResp{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		errorMessage = "DCC獲取agent於json解析失敗,內容=" + string(body)
		return
	}

	if response.Code != constant.ERROR_CODE_SUCCESS {
		errorMessage = fmt.Sprintf("DCC獲取agent失敗.錯誤代碼=%d,錯誤訊息=%s,內容=%s", response.Code, response.Msg, string(body))
		return
	}

	info, ok := response.Data.([]interface{})
	if !ok {
		errorMessage = fmt.Sprintf("DCC獲取agent的資料解析失敗,內容=%s", string(body))
		return
	}

	ClearDccAgentInfo()

	for _, value := range info {
		agentInfo, ok := value.(map[string]interface{})
		if !ok {
			continue
		}
		id := int(agentInfo["id"].(float64))
		if id == 1 || id == 2 { //不使用
			continue
		}
		var info DccAgentInfo
		info.ID = id
		info.Name = agentInfo["name"].(string)
		info.Md5key = agentInfo["md5_key"].(string)
		info.Aeskey = agentInfo["aes_key"].(string)

		DccAgentList[info.Name] = info
		agentList = append(agentList, info.Name)
	}
	return
}

func ClearDccAgentInfo() {
	DccAgentList = make(map[string]DccAgentInfo, 0)
}
