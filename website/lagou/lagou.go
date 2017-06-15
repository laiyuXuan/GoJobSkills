package lagou

import (
	"github.com/hu17889/go_spider/core/common/page"
	"log"
	"os"
	"runtime"
	"bytes"
	"regexp"
	"strconv"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

const (
	totalPageRx = "[page=\"*\" class=\"pager_not_current\"]"
	positionUrl = "https://www.lagou.com/jobs/positionAjax.json?px=default&needAddtionalResult=false"
)

var logger = getLogger()



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
 get log
 */
func getLogger() *log.Logger{
	_, file, _, _ := runtime.Caller(0);
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(file)
	buffer.WriteString("] ->>>")
	return log.New(os.Stdout, buffer.String(), log.LstdFlags)
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

func GetPositionIds(keyword string) []int {
	values := &url.Values{}
	values.Add("kd", keyword)
	positionIds := make([]int, 900)
	for pageNum := 1; pageNum <= 5; pageNum ++ {
		values.Add("pn", strconv.Itoa(pageNum))
		resp, err := http.PostForm(positionUrl, *values)
		if err != nil {
			logger.Println("GetPositionIds error, ", err)
			return nil
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Println("ioutil.ReadAll error, ", err)
			return nil
		}

		positionResponse := &PositionResponse{}
		err = json.Unmarshal(body, positionResponse)

		if err != nil {
			logger.Println("ioutil.ReadAll error, ", err)
			return nil
		}

		subPositionIds := make([]int, 15)
		positionInfos := positionResponse.Content.PositionResult.Result
		for i := 0; i < len(positionInfos)-1; i++ {
			subPositionIds[i] = positionInfos[i].PositionId
		}
		logger.Println("the %d loop result: %s", pageNum, subPositionIds)
		positionIds = append(positionIds, subPositionIds...)
	}

	return positionIds
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
