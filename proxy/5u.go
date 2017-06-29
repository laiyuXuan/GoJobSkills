package proxy

import (
	"log"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"GoJobSkills/model"
	"github.com/parnurzeal/gorequest"

)

// Data5u get ip from data5u.com
func Data5u() (results []*model.IP) {
	pollURL := "http://www.data5u.com/free/index.shtml"
	resp, _, errs := gorequest.New().Get(pollURL).End()
	if errs != nil {
		log.Println(errs)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Println(err.Error())
		return
	}
	doc.Find("body > div.wlist > li:nth-child(2) > ul").Each(func(i int, s *goquery.Selection) {
		node := strconv.Itoa(i + 1)
		ss := s.Find("ul:nth-child(" + node + ") > span:nth-child(1) > li").Text()
		sss := s.Find("ul:nth-child(" + node + ") > span:nth-child(2) > li").Text()
		ssss := s.Find("ul:nth-child(" + node + ") > span:nth-child(4) > li").Text()
		ip := model.NewIP()
		ip.Data = ss + ":" + sss
		ip.Type = ssss
		results = append(results, ip)
	})
	log.Println("Data5u done.")
	for _, result := range results {
		log.Println(result.Data)
	}

	return
}