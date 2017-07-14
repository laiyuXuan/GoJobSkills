package jianshu

import (
	"fmt"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"time"
	"goJobSkills/log"
	"github.com/parnurzeal/gorequest"
	"goJobSkills/client"
	"github.com/garyburd/redigo/redis"
	"goJobSkills/utils"
	"goJobSkills/constant"
	"os"
)

var logger = log.GetLogger()
var jianshuPrefix = "http://www.jianshu.com"
var path = "/Users/lyons/doc/jianshu/"

func GetArticleLinks() {

	conn := client.REDIS.Get()
	defer conn.Close()

	hrefs := make([]string, 0)
	totalArticleNum := 10000
	timeStamp := time.Now().Unix()
	logger.Println("using timeStamp :", timeStamp)
	for len(hrefs) < totalArticleNum {
		timeStamp = timeStamp - 1000
		url := "http://www.jianshu.com/recommendations/notes?category_id=56&max_id=" + strconv.FormatInt(timeStamp, 10)
		logger.Println("spidering -->" + url)
		proxy, err := redis.String(conn.Do("SRANDMEMBER", "proxy_pool"))
		if err != nil {
			logger.Panic(err)
		}
		logger.Println("using proxy -->>" + proxy)
		resp, _, errs := gorequest.New().Proxy(proxy).Get(url).Timeout(time.Second * 20).End()
		if errs != nil && !utils.IsTimeOut(errs[0]){
			logger.Println("time out for " + url + "with proxy" + proxy)
		}
		if errs != nil {
			continue
		}
		document, err := goquery.NewDocumentFromResponse(resp)

		if err != nil {
			logger.Println(err)
			continue
		}
		ul := document.Find("#list-container").Find(".note-list")

		list := ul.Find("a.title")

		idx := 0
		for i := 0; i < list.Length(); i++ {
			for j := 0; j < len(list.Get(i).Attr); j++ {
				if list.Get(i).Attr[j].Key == "href" {

					hrefs = append(hrefs, list.Get(i).Attr[j].Val)
					if len(hrefs) >= totalArticleNum{
						goto outsideLoop
					}
					idx ++
				}
			}
		}
	}

outsideLoop:
	fmt.Printf("%d article links retrieved", len(hrefs))
	for _, link := range hrefs{
		conn.Do("SADD", constant.KEY_JIANSHU_ARTICLES_LINKS, link)
	}

	getAndSaveContents(hrefs, conn)
}

func getAndSaveContents(links []string, conn redis.Conn) {



	for idx, link := range links{
		if len(link) == 0{
			continue
		}
		url := jianshuPrefix + link
		proxy, err := redis.String(conn.Do("SRANDMEMBER", "proxy_pool"))
		if err != nil {
			logger.Panic(err)
		}
		logger.Println("using proxy -->>" + proxy)
		fmt.Println("spiding -->> ", url)
		resp, _, errs := gorequest.New().Proxy(proxy).Get(url).Timeout(time.Second * 20).End()
		if errs != nil && !utils.IsTimeOut(errs[0]){
			logger.Println("time out for " + url + "with proxy" + proxy)
		}
		if errs != nil {
			continue
		}
		doc, err := goquery.NewDocumentFromResponse(resp)
		find := doc.Find(".show-content")
		fileName := "jianshu-" + strconv.Itoa(idx) + ".txt";
		save2File(fileName, find.Text())
	}

	logger.Println("---ALL DONE---")
}

func save2File(fileName string, content string)  {
	path := "/Users/lyons/doc/jianshu/"
	file, err := os.Create(path + fileName)
	defer file.Close()

	if err != nil{
		fmt.Println(file, err)
		return
	}
	file.WriteString(content)
}