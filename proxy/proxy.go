package proxy

import (
	"GoJobSkills/log"
	"math/rand"
	"github.com/parnurzeal/gorequest"
	"GoJobSkills/model"
	"fmt"
	"time"
	"github.com/garyburd/redigo/redis"
	"regexp"
	"strconv"
	"GoJobSkills/constant"
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
	logger.Println("proxy" + ip.Data + "checked, result" + strconv.Itoa(resp.StatusCode))
	if errs != nil {
		logger.Println(errs)
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func FillProxyPool() {
	client, err := redis.Dial("tcp", constant.REDIS_SERVER)
	defer client.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	compile := regexp.MustCompile(ipRx)

	results := getAllProxies()
	logger.Printf("%d proxies obtained", len(results))

	ips := make([]string, 0)
	for _, result := range results {
		if !compile.MatchString(result.Data){
			continue
		}
		if !CheckIP(result){
			continue
		}
		ips = append(ips, result.Data)
	}
	logger.Printf("%d proxies are valid, saving to redis", len(ips))

	for _, ip := range ips{
		client.Do("SADD", "proxy_pool", ip)
	}
	size, err := redis.Int(client.Do("SCARD", "proxy_pool"))
	if err != nil {
		logger.Println(err)
	}
	logger.Printf("current size of proxy pool is %d", size)
}

func getAllProxies()  (results []*model.IP){
	results = append(results, Data5u()...)
	results = append(results, GBJ()...)
	results = append(results, Xici()...)
	results = append(results, XDL()...)

	return
}
