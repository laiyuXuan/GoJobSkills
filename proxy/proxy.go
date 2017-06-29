package proxy

import (
	"GoJobSkills/log"
	"math/rand"
	"github.com/parnurzeal/gorequest"
	"GoJobSkills/model"
	"fmt"
	"time"
	"testing"
	"github.com/garyburd/redigo/redis"
	"regexp"
	"golang.org/x/net/proxy"
)

var logger = log.GetLogger()
var ipRx = "^(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9])\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[0-9]):\\d{0,5}$"

func GetRandomProxy() (ip string) {
	result := Data5u()
	ip = result[rand.Intn(10)].Data
	logger.Println("a random proxy is generated: ", ip)
	return
}

func CheckIP(ip *model.IP) bool {
	pollURL := "http://httpbin.org/get"
	resp, _, errs := gorequest.New().Proxy(ip.Data).Get(pollURL).Timeout(time.Second * 20).End()
	fmt.Println(resp)
	if errs != nil {
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func FillProxyPool(t *testing.T) {
	client, err := redis.Dial("tcp", "127.0.0.1:6379")
	defer client.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	compile := regexp.MustCompile(ipRx)

	results := Data5u()
	for idx, result := range results {
		fmt.Println(result.Data)
		if !compile.MatchString(result.Data){
			continue
		}
		if !CheckIP(result){
			continue
		}
		fmt.Println(idx)
		client.Do("SADD", "proxy_pool", result.Data)
	}
}

func getAllProxies()  (results []*model.IP){
	results = append(results, Data5u()...)


	return
}
