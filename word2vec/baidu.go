package word2vec

import (
	"github.com/parnurzeal/gorequest"
	"goJobSkills/log"
	"encoding/json"
	"math"
	"io/ioutil"
	"github.com/sajari/word2vec"
)

const (
	BAIDU_WORD2VEC_URL = "https://aip.baidubce.com/rpc/2.0/nlp/v2/word_emb_vec"
	BAIDU_AUTH_URL = "https://aip.baidubce.com/oauth/2.0/token?"
	BAIDU_API_KEY = "tfALNQtSwQZaakOAgNm1LrDF"
	BAIDU_SECRET_KEY = "Qoro99lcypRoeuApWVTvpyzrW2kgMlN9"
	BAIDU_ACCESS_TOKEN = "24.03746a0144ef95b4beba97997de43ad4.2592000.1503713610.282335-9937346"
	BAIDU_REFRESH_TOKEN = "25.6b87e0ad8ac3a01bbb56b04d0f0269ab.315360000.1816481610.282335-9937346"
)

var logger = log.GetLogger()



func GetVec(word string) []float64 {
	getVecUrl := BAIDU_WORD2VEC_URL +
	"?access_token=" + BAIDU_ACCESS_TOKEN
	param := new(BaiduVecParam)
	param.Word = word
	paramStr, errs := json.Marshal(param)
	if errs != nil {
		logger.Panic(errs)
	}
	params := string(paramStr)
	logger.Println(params)
	resp, _, err := gorequest.New().
		Post(getVecUrl).
		Set("Content-Type", "application/json").
		Send(params).
		End()
	if err != nil {
		logger.Panic(err)
	}
	logger.Println(resp)
	baiduVecResponse := &BaiduVecResponse{}
	body, errs := ioutil.ReadAll(resp.Body)
	if errs != nil {
		logger.Println(errs)
	}
	errs = json.Unmarshal(body, &baiduVecResponse)
	if errs != nil {
		logger.Panic(errs)
	}
	return baiduVecResponse.Vec
}

func GetBaiduToken() string {
	getAccessTokenUrl := BAIDU_AUTH_URL +
		"grant_type=client_credentials" +
		"&client_id=" + BAIDU_API_KEY   +
		"&client_secret=" + BAIDU_SECRET_KEY;
	resp, body, err := gorequest.New().Get(getAccessTokenUrl).End()
	if err != nil {
		logger.Panic(err)
	}
	logger.Println(resp)
	logger.Println(body)
	return ""
}

func CalcVecDiff(vecA, vecB []float64) float64 {
	sum := 0.0
	if len(vecA) != len(vecB) {
	 logger.Panic("dimensions should be equaled")
	}
	dem := len(vecA)
	for idx := 0; idx < dem ; idx ++ {
		sum += math.Pow(vecA[idx] - vecB[idx], 2)
	}
	return math.Sqrt(sum)
}

func CalcVecDistance(vecA, vecB word2vec.Vector) float64  {
	sum := 0.0
	if len(vecA) != len(vecB) {
		logger.Panic("dimensions should be equaled")
	}
	dem := len(vecA)
	for idx := 0; idx < dem ; idx ++ {
		sum += math.Pow(float64(vecA[idx] - vecB[idx]), 2)
	}
	return math.Sqrt(sum)
}

type BaiduVecParam struct {
	Word string `json:"word"`
}

type BaiduVecResponse struct {
	Vec []float64 `json:"vec"`
}

