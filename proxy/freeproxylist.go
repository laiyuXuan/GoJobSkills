package proxy

import (
	"goJobSkills/model"
	"github.com/PuerkitoBio/goquery"
	"goJobSkills/client"
	"strings"
)

func FreeProxyList() (results []*model.IP) {
	conn := client.REDIS.Get()
	defer  conn.Close()

	document, err := goquery.NewDocument("https://free-proxy-list.net/")
	if err != nil {
		logger.Panic(err)
	}

	proxies := make([]string, 0)
	table := document.Find("#proxylisttable")
	table.Find(" tbody > tr").Each(func(i int, row *goquery.Selection) {
		isHttps := row.Find("td").Eq(6).Text()
		if isHttps == "yes" {
			proxies = append(proxies, "https://" + row.Find("td").Eq(0).Text()  + ":" + row.Find("td").Eq(1).Text())
		} else
		{
			proxies = append(proxies, "http://" + row.Find("td").Eq(0).Text()  + ":" + row.Find("td").Eq(1).Text())
		}
	})
	for _, one := range proxies {
	 	conn.Do("SADD", "proxy_pool", one)
	}
	return
}

func GetFreeProxy(area string) string {
	document, err := goquery.NewDocument("https://free-proxy-list.net/")
	if err != nil {
		logger.Panic(err)
	}

	var lastCheckProxy string
	table := document.Find("#proxylisttable")
	table.Find(" tbody > tr").Each(func(i int, row *goquery.Selection) {
		if strings.Contains(row.Find("td").Eq(3).Text(),  area) {
			isHttps := row.Find("td").Eq(6).Text()
			var prefix string
			if isHttps == "yes" {
				prefix = "https://"
			} else {
				prefix = "http://"
			}
			lastCheckProxy = prefix + row.Find("td").Eq(0).Text() + ":" + row.Find("td").Eq(1).Text()
		}
	})
	logger.Println("using proxy -->> " + lastCheckProxy)
	return lastCheckProxy
}