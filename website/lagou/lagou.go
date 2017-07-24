package lagou

import (
	"github.com/hu17889/go_spider/core/common/page"
	"regexp"
	"strconv"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"goJobSkills/log"
	"goJobSkills/client"
	"github.com/garyburd/redigo/redis"
	"time"
	"strings"
	"github.com/google/uuid"
	"bytes"
	"os"
)

const (
	totalPageRx = "[page=\"*\" class=\"pager_not_current\"]"
	positionUrl = "https://www.lagou.com/jobs/positionAjax.json?px=default&needAddtionalResult=false"
)

var logger = log.GetLogger()



type LaGouPageProcessor struct {

}

func (processor *LaGouPageProcessor) Process(p *page.Page) {
	body := p.GetBodyStr()
	totalPage := getTotalPage(body)
	logger.Println("total page num is ", totalPage)
}

func (processor *LaGouPageProcessor) Finish()  {

}

/**
get total page
 */
func getTotalPage(body string) int  {
	compile := regexp.MustCompile((totalPageRx))
	allMatched := compile.FindAllString(body, -1)
	if len(allMatched) == 0  {
		return 0
	}
	lastMatched := allMatched[(len(allMatched))]
	atoi, e := strconv.Atoi(regexp.MustCompile("[0-9]+").FindString(lastMatched))
	if e != nil{
		logger.Println(e)
		return 0;
	}
	return atoi;
}

func GetPositionIds(keyword string) (positionIds []int) {

	conn := client.REDIS.Get()
	defer conn.Close()

	params := "kd=" + keyword
	request := gorequest.New()

	for pageNum := 1; pageNum <= 100; pageNum ++ {
			params = params + "&pn=" + strconv.Itoa(pageNum)
		newUUID, _ := uuid.NewUUID()

		_, body, errs := request.
			Proxy(getRandomProxy(conn)).
			Post(positionUrl).
			Set("REQUEST_ID", newUUID.String()).
			Set("Origin","https://www.lagou.com").
			Set("Referer","https://www.lagou.com/jobs/list_Java?city=%E5%8C%97%E4%BA%AC&cl=false&fromSearch=true&labelWords=&suginput=").
			Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36").
			Set("X-Forwarded-For", getRandomIP(conn)).
			Send(params).
			End()
		if errs != nil {
			logger.Println("GetPositionIds error, ", errs)
			return
		}

		logger.Println("resulet body -->>, " + body)
		positionResponse := &PositionResponse{}
		err := json.Unmarshal([]byte(body), positionResponse)

		if err != nil {
			logger.Println("ioutil.ReadAll error, ", err)
			return
		}

		subPositionIds := make([]int, 0)
		positionInfos := positionResponse.Content.PositionResult.Result
		for i := 0; i < len(positionInfos)-1; i++ {
			subPositionIds = append(subPositionIds, positionInfos[i].PositionId)
		}
		logger.Printf("the %d loop result: %d", pageNum, subPositionIds)

		for _, id := range subPositionIds {
			conn.Do("SADD", "position_id_lagou", id)
		}

		time.Sleep(time.Second * 20)
	}

	return
}

func GetJobDescription() {
	conn := client.REDIS.Get()
	defer conn.Close()

	positionIds, err := redis.Strings(conn.Do("SMEMBERS", "position_id_lagou"))
	if err != nil {
		logger.Panic(err)
		return
	}
	//positionIds := make([]string, 0)
	//positionIds = append(positionIds, "3159109")
	for _, id := range positionIds {
		logger.Println("position id: " + id)
		newUUID, _ := uuid.NewUUID()
		resp, body, errs := gorequest.
			New().
			Proxy(getRandomProxy(conn)).
			Set("REQUEST_ID", newUUID.String()).
			Set("Origin","https://www.lagou.com").
			Set("Referer","https://www.lagou.com/jobs/list_Java?city=%E5%8C%97%E4%BA%AC&cl=false&fromSearch=true&labelWords=&suginput=").
			Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36").
			Set("X-Forwarded-For", getRandomIP(conn)).
			Get("https://www.lagou.com/jobs/" + id + ".html").
			End()
		if errs != nil {
			logger.Panic(errs)
			continue
		}

		emptyRx := "\\s+"
		emptyCompile := regexp.MustCompile(emptyRx)
		body = emptyCompile.ReplaceAllString(body, "")

		buffer := &bytes.Buffer{}

		head := "<h3class=\"description\">"
		tail := "</div>"
		buffer.WriteString(head)
		buffer.WriteString("(.*?)")
		buffer.WriteString(tail)

		requirementRx := buffer.String()
		compile := regexp.MustCompile(requirementRx)
		matched := compile.FindString(body)

		//multiHtml := "<p><br></p>"
		////multiHtml := "<p>&nbsp;</p>"
		//multiHtmlsCompile := regexp.MustCompile(multiHtml)
		//split := multiHtmlsCompile.Split(matched, -1)
		//
		//fmt.Println(len(split))
		//fmt.Println(split[1])

		htmlLabelRx := "<.+?>"
		htmlCompile := regexp.MustCompile(htmlLabelRx)
		matched = htmlCompile.ReplaceAllString(matched, "")

		//matched = strings.Replace(matched, head, "", -1)
		//matched = strings.Replace(matched, tail, "", -1)

		htmlSpaceRx := "&nbsp"
		htmlSpaceComplie := regexp.MustCompile(htmlSpaceRx)
		matched = htmlSpaceComplie.ReplaceAllString(matched, "")

		matched += "/n"
		logger.Println(strconv.Itoa(resp.StatusCode) + matched)


		save2File("job_description", matched)
		time.Sleep(time.Second * 20)
	}
}

func getRandomProxy(conn redis.Conn) (string) {
	proxy, err := redis.String(conn.Do("SRANDMEMBER", "proxy_pool"))
	if err != nil {
		logger.Panic(err)
	}
	logger.Println("using proxy -->>" + proxy)
	return proxy
}

func getRandomIP(conn redis.Conn) string {
	proxy := getRandomProxy(conn)
	split := strings.Split(proxy, ":")
	return split[0]
}


	type PositionResponse struct {
		Code int

		Msg string

		Content Content
	}

	type Content struct {
		PageNo int

		PageSize int

		PositionResult PositionResult
	}

	type PositionResult struct {

		Result []PositionInfo

		ResultSize int

		TotalCount int
	}

	type PositionInfo struct {
		AdWord int
		AppShow int
		Approve int
		BusinessZones []string
		City string
		CompanyFullName string
		CompanyId int
		CompanyLabelList []string
		CompanyLogo string
		CompanyShortName string
		CompanySize string
		CreateTime string
		Deliver int
		District string
		Education string
		Explain string
		FinanceStage string
		FirstType string
		FormatCreateTime string
		GradeDescription string
		ImState string
		IndustryField string
		IndustryLables []string
		JobNature string
		LastLogin int64
		PcShow int
		Plus string
		PositionAdvantage string
		PositionId int
		PositionLables []string
		PositionName string
		PromotionScoreExplain string
		PublisherId int
		Salary string
		Score int
		SecondType string
		WorkYear string
	}

func save2File(fileName string, content string)  {
	path := "/Users/lyons/doc/lagou/"
	file, err := os.OpenFile(path + fileName, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0666)
	defer file.Close()

	if err != nil{
		logger.Println(file, err)
		return
	}
	_, err = file.WriteString(content)
	if err != nil {
		logger.Println("append to file failed", err)
		return
	}

}