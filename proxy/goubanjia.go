package proxy

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/parnurzeal/gorequest"
	"GoJobSkills/model"
)

// GBJ get ip from goubanjia.com
func GBJ() (result []*model.IP) {
	log.Println("start retrieving proxies from goubanjia...")

	pollURL := "http://www.goubanjia.com/free/gngn/index"
	for i := 1; i <= 10; i++ {
		resp, _, errs := gorequest.New().Get(pollURL + strconv.Itoa(i) + ".shtml").End()
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

		doc.Find("#list > table > tbody > tr").Each(func(_ int, s *goquery.Selection) {
			sf, _ := s.Find(".ip").Html()
			tee := regexp.MustCompile("<pstyle=\"display:none;\">.?.?</p>").ReplaceAllString(strings.Replace(sf, " ", "", -1), "")
			re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
			ip := model.NewIP()
			ip.Data = re.ReplaceAllString(tee, "")
			ip.Type = s.Find("td:nth-child(3) > a").Text()
			result = append(result, ip)
		})
	}
	log.Println("GBJ done.")
	return
}
