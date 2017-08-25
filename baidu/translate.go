package baidu

import (
	"time"
	"strconv"
	"crypto/md5"
	"goJobSkills/log"
	"github.com/parnurzeal/gorequest"
	"fmt"
	"encoding/hex"
	"strings"
	"encoding/json"
)

const (
	URL = "https://fanyi-api.baidu.com/api/trans/vip/translate"
	APP_ID = "20170727000068718"
	SECRET_KEY = "RIczf0dXKio2duhhk7nv"
)

var logger = log.GetLogger()

func En2Zh(word string)  string {
	fullUrl := URL 	+
		"?q=" + word +
		"&from=en&to=zh&appid=" + APP_ID
	return commonSendRequest(fullUrl, word)
}

func Zh2En(word string) string {
	fullUrl := URL 	+
		"?q=" + word +
		"&from=zh&to=en&appid=" + APP_ID
	return commonSendRequest(fullUrl, word)
}

func commonSendRequest(fullUrl, word string) string {
	salt := strconv.Itoa(time.Now().Nanosecond())
	hashier := md5.New()
	hashier.Write([]byte(APP_ID + word + salt + SECRET_KEY))
	sign := hex.EncodeToString(hashier.Sum(nil))
	fullUrl += 	"&salt=" + salt +
		"&sign=" + sign
	_, body, errs := gorequest.New().Get(fullUrl).End()
	if errs != nil {
		logger.Panic(errs)
	}
	translateResponse := &translateResponse{}
	json.Unmarshal([]byte(body), &translateResponse)
	result := translateResponse.TranslateResult[0].Dest
	logger.Printf("%s translated into %s", word, result)
	return result
	//dests := strings.Split(translateResponse.TranslateResult[0].Dest, "---")
	//srcs := strings.Split(translateResponse.TranslateResult[0].Src, "---")
	//for _, src := range srcs {
	//	src = convertUnicode2Zh(src)
	//}
	//logger.Printf("srcs len: %d", len(srcs))
	//logger.Printf("dest len: %d", len(dests))
	//logger.Println(dests)
	//logger.Println(srcs)
	//resultMap := make(map[string]string)
	//for idx, src := range srcs {
	//	resultMap[dests[idx]] = src
	//}
	//return resultMap
}

func convertUnicode2Zh(word string) string {
	zhs := strings.Split(word, "\\u")
	var zh string
	for _, v := range zhs {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt("82f9", 16, 32)
		if err != nil {
			logger.Panic(err)
		}
		zh += fmt.Sprintf("%c", temp)
	}
	return zh
}

type translateResult struct {
	Src		string	`json:"src"`
	Dest	string 	`json:"dst"`
}

type translateResponse struct {
	TranslateResult []translateResult `json:"trans_result"`
}