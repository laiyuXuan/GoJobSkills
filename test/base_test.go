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
	"GoJobSkills/proxy"
	"github.com/PuerkitoBio/goquery"
	"os"
	"strconv"
	"time"
	"github.com/garyburd/redigo/redis"
	"github.com/parnurzeal/gorequest"
	"net"
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
func TestRegex(t *testing.T) {
	totalPageRx := "^(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9])\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[0-9]):\\d{0,5}$"
	compile, _ := regexp.Compile(totalPageRx)
	matched := compile.MatchString("IP:端口")
	fmt.Println(matched)

}

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
	resp, err := http.Get("https://www.lagou.com/jobs/2365703.html")
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
	////multiHtml := "<p>&nbsp;</p>"
	//multiHtmlsCompile := regexp.MustCompile(multiHtml)
	//split := multiHtmlsCompile.Split(matched, -1)
	//
	//fmt.Println(len(split))
	//fmt.Println(split[1])

	fmt.Println("=====================")

	htmlLabelRx := "<.+?>"
	htmlCompile := regexp.MustCompile(htmlLabelRx)
	matched = htmlCompile.ReplaceAllString(matched, "")

	//matched = strings.Replace(matched, head, "", -1)
	//matched = strings.Replace(matched, tail, "", -1)

	htmlSpaceRx := "&nbsp"
	htmlSpaceComplie := regexp.MustCompile(htmlSpaceRx)
	matched = htmlSpaceComplie.ReplaceAllString(matched, "")

	fmt.Println(matched)
}

func TestGoRequest(t *testing.T) {
	pollURL := "http://httpbin.org/get"
	resp, _, errs := gorequest.New().Proxy("http://167.114.211.122:8246").Get(pollURL).Timeout(time.Second * 20).End()
	if errs != nil{
		fmt.Println(errs)
		return
	}
	fmt.Println(resp.StatusCode)
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

func TestProxy(t *testing.T) {
	client, err := redis.Dial("tcp", "127.0.0.1:6379")
	defer client.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	ipRx := "^(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9])\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[1-9]|0)\\.(25[0-5]|2[0-4][0-9]|[0-1]{1}[0-9]{2}|[1-9]{1}[0-9]{1}|[0-9]):\\d{0,5}$"
	compile := regexp.MustCompile(ipRx)

	results := proxy.Data5u()
	for idx, result := range results {
		fmt.Println(result.Data)
		if !compile.MatchString(result.Data){
			continue
		}
		if !proxy.CheckIP(result){
			continue
		}
		fmt.Println(idx)
		client.Do("SADD", "proxy_pool", result.Data)
	}
}

func TestGoQuery(t *testing.T) {
	doc, err := goquery.NewDocument("http://www.jianshu.com/p/d1dc265c9dd3")
	if err != nil {
		fmt.Println(err)
	}
	find := doc.Find(".show-content")
	println(find.Text())

}

func TestJianShuLink(t *testing.T) {
	client, error := redis.Dial("tcp", "127.0.0.1:6379")
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



	//nodes := ul.Nodes
	//for i := 0; i < len(nodes); i++{
	//	fmt.Println(nodes[0].FirstChild.Attr[0].Val)
	//}

}

func TestJianShuContent(t *testing.T) {
	jianshuPrefix := "http://www.jianshu.com"
	client := getRedis()
	defer client.Close()

	if client == nil {
		return
	}
	hrefs, error := redis.Strings(client.Do("SMEMBERS", "go_hrefs"))
	if error != nil {
		fmt.Println(error)
		return
	}


	for i := 0; i < len(hrefs); i ++{
		if len(hrefs[i]) == 0{
			continue
		}
		url := jianshuPrefix + hrefs[i]

		fmt.Println("spiding -->> ", url)
		doc, err := goquery.NewDocument(url)
		if err != nil {
			fmt.Println(err)
		}
		find := doc.Find(".show-content")
		fileName := "jianshu-" + strconv.Itoa(i) + ".txt";
		save2File(fileName, find.Text())

	}
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

func TestFile(t *testing.T) {
	save2File("", "test")
}

func TestRedis(t *testing.T) {
	client, error := redis.Dial("tcp", "127.0.0.1:6379")
	defer client.Close()

	if error != nil {
		fmt.Println(error)
		return
	}

	client.Do("SET", "testGo", "goTest")
	reply, err := redis.String(client.Do("GET", "testGo"))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(reply)
}

func TestTime(t *testing.T) {
	fmt.Println(time.Now().Unix())
}

func getRedis() (client redis.Conn) {
	client, error := redis.Dial("tcp", "127.0.0.1:6379")

	if error != nil {
		fmt.Println(error)
		return nil
	}
	return
}

func TestTimeOut(t *testing.T) {
	agent := gorequest.New()
	resp, str, err := agent.Get("http://127.0.0.1:9085/test/testOT").Send("").Timeout(time.Second * 10).End()
	if err, ok := err[0].(net.Error); ok && err.Timeout(){
		fmt.Println("yes")
		return
	}
	if err != nil {
		fmt.Println("error -->>   ", err)
	}
	fmt.Println(str)
	fmt.Println(resp)
}

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for _ = range ticker.C{
			fmt.Println("tick tock")
		}
	}()

	ch := make(chan int)
	value := <- ch
	fmt.Println("value =", value)
}

func TestProxySources(t *testing.T) {
	//fmt.Println(proxy.IP66()) bad
	//results := proxy.GBJ() good
	//results := proxy.IP181() bad
	//results := proxy.KDL() empty
	//results := proxy.PLP() bad
	//results := proxy.XDL() good
	//results := proxy.Xici() good
	results := proxy.YDL() // bad
	for _, result := range results {
		fmt.Println(result.Data)
	}
}

func TestArray(t *testing.T) {
	args := make([]string, 0)
	args = append(args, "set_test_key")
	args = append(args, "1")
	args = append(args, "2")
	args = append(args, "3")

	client, err := redis.Dial("tcp", "127.0.0.1:6379")
	defer client.Close()

	if err != nil {
		fmt.Println(err)
		return
	}

	client.Do("SADD", args...)
}