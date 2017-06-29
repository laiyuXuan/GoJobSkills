package jianshu

import (
	"fmt"
	"strconv"
	"github.com/PuerkitoBio/goquery"
	"github.com/garyburd/redigo/redis"
	"time"
	"GoJobSkills/constant"
)

func TestJianShuLink() {
	client, error := redis.Dial("tcp", constant.REDIS_SERVER)
	defer client.Close()

	if error != nil {
		fmt.Println(error)
		return
	}


	hrefs := make([]string, 0)
	totalArticleNum := 1000
	timeStamp := time.Now().Unix()
	fmt.Println("using timeStamp :", timeStamp)
	for len(hrefs) < totalArticleNum {
		timeStamp = timeStamp - 1000
		url := "http://www.jianshu.com/recommendations/notes?category_id=56&max_id=" + strconv.FormatInt(timeStamp, 10)
		fmt.Println("spidering -->" + url)
		document, err := goquery.NewDocument(url)
		if err != nil {
			fmt.Println(err)
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
	fmt.Println(hrefs)
	fmt.Println(len(hrefs))
	for i := 0; i < len(hrefs); i++ {
		client.Do("SADD", "go_hrefs", hrefs[i])
	}
}
