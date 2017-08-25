package lagou

import (
	"regexp"
	"strconv"
	"encoding/json"
	"github.com/parnurzeal/gorequest"
	"goJobSkills/log"
	"goJobSkills/client"
	"github.com/garyburd/redigo/redis"
	"time"
	"github.com/google/uuid"
	"bytes"
	"goJobSkills/utils"
	"goJobSkills/proxy"
)

const (
	totalPageRx       		= "[page=\"*\" class=\"pager_not_current\"]"

	MAX_PAGE_INDEX    		= 35
	MAX_POST_DURATION 		= time.Hour * 24 * 7
	INTERVAL          		= time.Second * 20
	TIME_FORMAT       		= "2006-01-02 15:04:05"
	JD_FILE_PATH	  		= "/Users/lyons/doc/lagou/"

	KEY_POSITION_ID_PREFIX 	= "position_id_lagou_"
)

var logger = log.GetLogger()
var KEY_DATE_POSTFIX = "@" + time.Now().Format("2006-01-02")

//type LaGouPageProcessor struct {
//
//}
//
//func (processor *LaGouPageProcessor) Process(p *page.Page) {
//	body := p.GetBodyStr()
//	totalPage := getTotalPage(body)
//	logger.Println("total page num is ", totalPage)
//}
//
//func (processor *LaGouPageProcessor) Finish()  {
//
//}

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

func GetPositionIds(keyword, city, workYear string) {
	conn := client.REDIS.Get()
	defer conn.Close()

	params := "kd=" + keyword
	request := gorequest.New()
	positionUrl := getPositionUrl(workYear, city)
	for pageNum := 1; pageNum <= MAX_PAGE_INDEX; pageNum ++ {
		params = params + "&pn=" + strconv.Itoa(pageNum)
		newUUID, _ := uuid.NewUUID()
		_, body, errs := request.
		Proxy(proxy.GetRandomProxy(conn)).
			Post(positionUrl).
			Set("REQUEST_ID", newUUID.String()).
			Set("Origin", "https://www.lagou.com").
			Set("Referer", "https://www.lagou.com/jobs/list_"+keyword+"?city=%E5%8C%97%E4%BA%AC&cl=false&fromSearch=true&labelWords=&suginput=").
			Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36").
			Set("X-Forwarded-For", proxy.GetRandomIP(conn)).
			Send(params).
			End()
		if errs != nil {
			logger.Println("GetPositionIds error, ", errs)
			continue
		}

		logger.Println("resulet body -->>, " + body)
		positionResponse := &PositionResponse{}
		err := json.Unmarshal([]byte(body), positionResponse)
		if err != nil {
			logger.Println("ioutil.ReadAll error, ", err)
			continue
		}

		positionIds := make([]int, 0)
		positionInfos := positionResponse.Content.PositionResult.Result
		for i := 0; i < len(positionInfos)-1; i++ {
			positionIds = append(positionIds, positionInfos[i].PositionId)
		}
		logger.Printf("the %d loop result: %d", pageNum, positionIds)

		for _, id := range positionIds {
			conn.Do("SADD", KEY_POSITION_ID_PREFIX+keyword+KEY_DATE_POSTFIX, id)
		}

		createTime, err := time.Parse(TIME_FORMAT, positionInfos[len(positionInfos)-1].CreateTime)
		if err != nil {
			logger.Println(err)
		}

		if createTime.Add(MAX_POST_DURATION).Before(time.Now()) {
			logger.Println("reach the post of one week ago, abort...")
			break
		}

		time.Sleep(INTERVAL)
	}
	logger.Println("ALL DONE!")
}
func getPositionUrl(workYear string, city string) string {
	positionUrl := "https://www.lagou.com/jobs/positionAjax.json?px=default&needAddtionalResult=false"
	if workYear != "" {
		positionUrl += "&gj=" + workYear
	}
	if city != "" {
		positionUrl += "&city=" + city
	}
	return positionUrl
}

func GetJobDescription(keyword string) {
	conn := client.REDIS.Get()
	defer conn.Close()

	positionIds, err := redis.Strings(conn.Do("SMEMBERS", KEY_POSITION_ID_PREFIX + keyword + KEY_DATE_POSTFIX))
	if err != nil {
		logger.Panic(err)
	}
	limits := 0
	for _, id := range positionIds {
		logger.Println("position id: " + id)
		newUUID, _ := uuid.NewUUID()
		resp, body, errs := gorequest.
			New().
			Proxy(proxy.GetRandomProxy(conn)).
			Set("REQUEST_ID", newUUID.String()).
			Set("Origin","https://www.lagou.com").
			Set("Referer","https://www.lagou.com/jobs/list_" + keyword + "?city=%E5%8C%97%E4%BA%AC&cl=false&fromSearch=true&labelWords=&suginput=").
			Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36").
			Set("X-Forwarded-For", proxy.GetRandomIP(conn)).
			Get("https://www.lagou.com/jobs/" + id + ".html").
			End()
		if errs != nil {
			logger.Println(errs)
			continue
		}

		body = utils.RemoveBlanks(body)

		requirementRx := getJDRegex()
		compile := regexp.MustCompile(requirementRx)
		matched := compile.FindString(body)

		matched = utils.RemoveHtmlTag(matched)
		matched = utils.RemoveSpace(matched)

		logger.Println(resp.Status + " " + matched)

		if matched == "" {
			limits ++
		}

		err = utils.Save2File(JD_FILE_PATH + "job_description_" + keyword + time.Now().Format("2006-01-02"), matched)
		if err != nil {
			logger.Panic(err)
		}

		time.Sleep(time.Second * 10)
	}
	logger.Printf("ALL DONE! failed %d times", limits)
}

func getJDRegex() string {
	buffer := &bytes.Buffer{}
	head := "<h3class=\"description\">"
	tail := "</div>"
	buffer.WriteString(head)
	buffer.WriteString("(.*?)")
	buffer.WriteString(tail)
	requirementRx := buffer.String()
	return requirementRx
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
