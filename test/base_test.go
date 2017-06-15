package test

import (
	"testing"
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"regexp"
	"bytes"
	"GoJobSkills/website/lagou"
)


//func TestCentipede(t *testing.T) {
//	centipede := crawler.NewCentipede()
//	centipede.Run();
//}
//
//func TestRunTime(t *testing.T){
//	_, file, _, _ := runtime.Caller(0);
//	fmt.Println(file)
//}
//
//func TestRegex(t *testing.T) {
//	totalPageRx := "[page=\"*\" class=\"pager_not_current\"]"
//	compile, _ := regexp.Compile(totalPageRx)
//	matched := compile.MatchString("page=\"30\" class=\"pager_not_current\"")
//	fmt.Println(matched)
//
//}

func TestHttpPost(t *testing.T) {
	values := url.Values{}
	values.Add("first", "true")
	values.Add("pn", "1")
	values.Add("kd", "ui")
	resp, err:= http.PostForm("https://www.lagou.com/jobs/positionAjax.json?needAddtionalResult=false", values)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
	all, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
	laGouResponse := &lagou.PositionResponse{}
	err = json.Unmarshal(all, laGouResponse)
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
	fmt.Println(laGouResponse.Content.PositionResult.Result[5].PositionLables)
}

func TestHttpGet(t *testing.T) {
	resp, err := http.Get("https://www.lagou.com/jobs/2728503.html")
	if err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println("error")
		fmt.Println(err)
	}
	body := string(all[:])

	//fmt.Println(body)

	fmt.Println("=====================")

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
	multiHtml := "<p>&nbsp;</p>"
	multiHtmlsCompile := regexp.MustCompile(multiHtml)
	split := multiHtmlsCompile.Split(matched, -1)

	fmt.Println(len(split))
	fmt.Println(split[1])

	fmt.Println("=====================")

	htmlLabelRx := "<.+?>"
	htmlCompile := regexp.MustCompile(htmlLabelRx)
	matched = htmlCompile.ReplaceAllString(matched, "")

	//matched = strings.Replace(matched, head, "", -1)
	//matched = strings.Replace(matched, tail, "", -1)


	fmt.Println(matched)
}

func TestRx(t *testing.T) {
	emptyRx := "\\s+"
	emptyCompile := regexp.MustCompile(emptyRx)
	fmt.Println(emptyCompile.ReplaceAllString("1 2 3    4", ""))
}

func TestPositionId(t *testing.T) {
	ids := lagou.GetPositionIds("ui")
	fmt.Println(ids)
}