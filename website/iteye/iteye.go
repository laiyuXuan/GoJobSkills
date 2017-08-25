package iteye

import (
	"goJobSkills/client"
	"github.com/parnurzeal/gorequest"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"goJobSkills/log"
	"goJobSkills/utils"
	"time"
	"math/rand"
	"goJobSkills/proxy"
)

const (
	RATE_LIMIT = 95
	RATE_LIMIT_PERIOD = time.Minute * 30
)

var logger = log.GetLogger()

func CrawNews() {

	conn := client.REDIS.Get()
	defer conn.Close()

	url := "http://www.iteye.com/news/"
	newsId := 16327
	fail := 0
	times := 0

	areaList := []string{"United Kingdom", "United States", "Hong Kong", "Taiwan", "Germany", "Thailand", "Japan", "India", "France", "Canada"}
	//are??aList := make([]string, 10)
	//areappend(areaList, "United Kingdom")



	freeProxy := proxy.GetFreeProxy(areaList[rand.Intn(len(areaList))])
	for newsId > 0 {
		fullUrl := url + strconv.Itoa(newsId)
		logger.Println("fetching -- " + fullUrl)


		resp, _, err := gorequest.New().
			Proxy(freeProxy).
			Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36").
			Set("Host", "www.iteye.com").
			Set("Referer", "http://www.iteye.com/").
			Set("Accept-Encoding", "gzip, deflate").
			Set("Accept-Language", "zh-CN,zh;q=0.8,en;q=0.6").
			Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8").
			Set("Cache-Control", "max-age=0").
			Set("Connection", "keep-alive").
			Set("Upgrade-Insecure-Requests", "1").
			Timeout(time.Second * 10).
			Get(fullUrl).
			End()
		if err != nil{
			logger.Println(err)

			//if strings.Contains(err[0].Error(), "proxyconnect") {
			//	logger.Println("changing proxy due to proxy connect failed ...")
				freeProxy = proxy.GetFreeProxy(areaList[rand.Intn(len(areaList))])
			//}
			//if  strings.Contains(err[0].Error(), "timeout") {
			//	logger.Println("changing proxy due time out...")
			//	freeProxy = proxy.GetFreeProxy(areaList[rand.Intn(len(areaList))])
			//}
			continue
		}

		if resp == nil {
			logger.Println("no response on " + strconv.Itoa(newsId))
			continue
		}
		if resp.StatusCode == 401 {
			logger.Println("you've been forbidden!!!")
			freeProxy = proxy.GetFreeProxy(areaList[rand.Intn(len(areaList))])
			times = 0
			continue
		}
		if resp.StatusCode == 403 {
			logger.Println("response 403")
			logger.Println(resp)
			freeProxy = proxy.GetFreeProxy(areaList[rand.Intn(len(areaList))])
			times = 0
			continue
		}
		if resp.StatusCode == 503 {
			logger.Println("proxy failed")
			freeProxy = proxy.GetFreeProxy(areaList[rand.Intn(len(areaList))])
			times = 0
			continue
		}
		logger.Printf("response: %d", resp.StatusCode)

		doc, errs := goquery.NewDocumentFromReader(resp.Body)
		if errs != nil {
			logger.Println(errs)
			continue
		}


		div := doc.Find("#news_content")
		content := div.Text()

		if len(content) == 0 {
			fail ++
		}

		content = utils.RemoveHtmlTag(content)
		logger.Println("fetched -- " + strconv.Itoa(newsId))

		utils.Save2File("/Users/lyons/doc/iteye/news_content_"  + time.Now().Format("2006-01-02"), content)
		time.Sleep(time.Second * time.Duration(2 + rand.Intn(2)))

		if times == 50 {
			logger.Println("changing proxy")
			freeProxy = proxy.GetFreeProxy(areaList[rand.Intn(len(areaList))])
			times = 0
		}
		times ++
		newsId --
	}
	logger.Printf("ALL DONE! failed %d times", fail)
}

func FindRateLimit() {
	times := 0
	//fiveMinLater := time.Now().Add(time.Minute  * 20)
	conn := client.REDIS.Get()
	defer  conn.Close()


	for true  {
		times ++
		resp, _, err := gorequest.New().Proxy("http://217.15.85.202:8080").Get("http://www.iteye.com/news/32571").End()
		if err != nil {
			logger.Panic(err)
		}

		logger.Println(resp.StatusCode)

		if resp.StatusCode == 401 {
			logger.Panicf(strconv.Itoa(resp.StatusCode), times)
		}
		//if times == 90 {
		//	logger.Println("let's sleep a while")
		//	time.Sleep(fiveMinLater.Sub(time.Now()))
		//	fiveMinLater = time.Now().Add(time.Minute * 20)
		//	times = 0
		//}
	}

}