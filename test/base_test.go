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
	"goJobSkills/website/lagou"
	"goJobSkills/proxy"
	"github.com/PuerkitoBio/goquery"
	"os"
	"strconv"
	"time"
	"github.com/garyburd/redigo/redis"
	"github.com/parnurzeal/gorequest"
	"net"
	"goJobSkills/client"
	"goJobSkills/website/jianshu"
	"goJobSkills/boson"
	"path/filepath"
	"github.com/google/uuid"
	"github.com/yanyiwu/gojieba"
	"strings"
	"github.com/op/go-logging"
	"goJobSkills/log"
	"goJobSkills/word2vec"
	"math"
	"goJobSkills/baidu"
	"sort"
	"goJobSkills/utils"
	"goJobSkills/website/iteye"
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

var logger = log.GetLogger()
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
	resp, err := http.Get("https://www.lagou.com/jobs/3380622.html")
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

//func TestPositionId(t *testing.T) {
//	client.Init()
//	lagou.GetPositionIds("ui")
//}

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

	crr, err := time.Parse("2006-01-02 15:04:05", "2017-07-25 14:30:38")
	fmt.Println(err, crr)
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

func TestJianshuJob(t *testing.T) {
	client.Init()
	jianshu.GetArticleLinks()

	//dial, err := redis.Dial("tcp", "127.0.0.1:6379")
	//conn := client.REDIS.Get()
	//defer conn.Close()
	//defer dial.Close()
	//proxy, err := redis.String(dial.Do("SRANDMEMBER", "proxy_pool"))
	//proxy, err := redis.Int(dial.Do("GET", "test"))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(proxy)
}

func TestFillProxyPool(t *testing.T) {
	client.Init()
	proxy.FillProxyPool()
}

func TestCheckUnavailable(t *testing.T) {
	client.Init()
	proxy.CheckAvailablity()
}


func TestFilePath(t *testing.T) {
	abs, _ := filepath.Abs("../stopwords")
	fmt.Println(abs)
}

func TestGetLaGouTotalPage(t *testing.T) {
	client.Init()
	lagou.GetPositionIds("php", "广州", "应届毕业生")
}

func TestJD(t *testing.T) {
	client.Init()
	lagou.GetJobDescription("php")
}

func TestUUID(t *testing.T) {
	newUUID, _ := uuid.NewUUID()
	fmt.Println(newUUID)
}

func TestName(t *testing.T) {
	client.Init()
	conn := client.REDIS.Get()
	defer conn.Close()

	conn.Do("SADD", "position_id_lagou", 123)
}

func TestOpenFile(t *testing.T) {
	_, err := os.OpenFile("/Users/Lyons/doc/lagou/job_description", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
}

func TestJieBa(t *testing.T) {
	gojieba.STOP_WORDS_PATH = "/Users/Lyons/doc/stopwords/stopwords";
	jieba := gojieba.NewJieba()
	cut := jieba.Cut("任职要求:1.本科及以上学历，计算机相关专业优先。2.3年以上开发经验，具有互联网行业开发经验者优先。3.精通Java开发，熟悉SpringMVC、Spring、Mybatis或Hibernate框架，熟悉HTML，JavaScript，jQuery，熟悉MySQL，了解Redis，熟悉Git或Svn，熟悉IDEA，Maven，Linux优先。4.熟悉软件开发流程，良好的代码编写风格和文档编写能力。5.热爱软件编程，喜欢钻研技术，善于学习新的技术和理念并应用在工作中。", true)
	fmt.Println(strings.Join(cut, " "))
}

func TestBoSon(t *testing.T) {
	boson.GetKeywords("/Users/Lyons/doc/lagou/job_description_python2017-07-25")
	//boson.CharacterizeWords("/Users/Lyons/doc/lagou/job_description")
}

func TestGlog(t *testing.T) {
	//glog.Info("info")
	//glog.Warning("warning")
	//glog.Error("error")
	//glog.Fatal("fatal")
	logger1 := logging.MustGetLogger("testGlog")
	logger1.Info("info")
	//logger.Debug("debug")
	//logger.Warning("warning")
	//logger.Error("err")
	//logger.Panic("panic")
	//logger.Fatal("fatal")
	logger.Println("test")

}

func TestAll(t *testing.T) {
	client.Init()
	lagou.GetPositionIds("python", "全国", "")
	lagou.GetJobDescription("python")
}

func TestGetBaiduToken(t *testing.T) {
	word2vec.GetBaiduToken();
}

func TestBaiduGetVec(t *testing.T) {
	fmt.Println(word2vec.GetVec("张飞"))
}

func TestMath(t *testing.T) {
	fmt.Println(math.Pow(1.0 - 3.0, 2) + math.Pow(3.0 - 5.0, 2))
	fmt.Println(math.Sqrt(4.0))
}

func TestCalcVecDiff(t *testing.T) {
	fmt.Println(word2vec.CalcVecDiff(word2vec.GetVec("java"), word2vec.GetVec("uncle")))
}

func TestWord2Vec(t *testing.T) {
	word := "java"
	vectorMap := word2vec.LoadModel("/Users/lyons/Downloads/glove.6B/glove.6B.300d.txt")
	jVec := vectorMap[word]
	logger.Println(len(vectorMap))
	logger.Printf( "%s vector: %f",word,  jVec)
	mysqlVec := vectorMap["mysql"]
	hairVec := vectorMap["hair"]
	logger.Printf( "mysql vector: %f",  mysqlVec)
	logger.Printf( "hair vector: %f",  hairVec)
	logger.Println(word2vec.CalcVecDiff(jVec, mysqlVec))
	logger.Println(word2vec.CalcVecDiff(jVec, hairVec))

	//distances := make(word2vec.Distances, 0)
	//for k, v := range vectors {
	//	distance := &word2vec.WordDistance{"java", k, word2vec.CalcVecDistance(jVec, v)}
	//	distances = append(distances, *distance)
	// 	//fmt.Printf("distance of %s: %f", k, word2vec.CalcVecDistance(jVec, v))
	//}
	//sort.Sort(distances)
	//for _, wd := range distances {
	// 	logger.Printf("distance of %s: %f", wd.WordB, wd.Distance)
	//}

}

func TestBaiduTranslate(t *testing.T) {
	fmt.Println(baidu.Zh2En("数据库,架构,开源,并发,测试,线程,缓存,性能,数据,质量,金融,解决方案,保证系统,积极主动,一年,两年,部署,福利,很强,资格,重构,功能,协作,一定,流程,知识,软件产品,产品开发,平台,能够,深入,过程,全日制,意识,行业,详细,快速,习惯,语句,机制,员工,经验者,结构,任务,程序"))
	//fmt.Println(baidu.convertUnicode2Zh("\u82f9\u679c"))
}

func TestFull(t *testing.T) {
	wordWeights := boson.GetKeywords("/Users/Lyons/doc/lagou/job_description_python2017-07-25")
	tooLight := make([]string, 0)
	for _, ww := range wordWeights {
	 	if ww.Weight < boson.THEARDHOLD && utils.IsChineseChar(ww.Word){
			tooLight = append(tooLight, ww.Word)
		}
	}
	tooLightMap := make(map[string]string)
	for _, tt := range tooLight {
	 	tooLightMap[baidu.Zh2En(tt)] = tt
	}
	tooLightEn := make([]string, 0)
	for k, _ := range tooLightMap {
		tooLightEn = append(tooLightEn, k)
	}
	logger.Println("tooLightEn -> ", tooLightEn)
	getDistance("python", tooLightEn, tooLightMap)
}

func getDistance(fromWord string, toWords []string, toWordMap map[string]string) {
	vectorMap := word2vec.LoadModel("/Users/lyons/word2vec/pre-trained-models/GoogleNews-vectors-negative300.txt")

	fromVec := vectorMap[fromWord]
	logger.Printf("from vec is %s\n", fromVec)
	distances := make(word2vec.Distances, 0)
	//for k, v := range vectorMap {
	//	if v == nil || len(v) == 0{
	//		logger.Printf("no vector for %s", k)
	//		continue
	//	}
	//	logger.Printf("comparing %s, whose vector length is %d\n", k, len(v))
	//	distance := &word2vec.WordDistance{fromWord, k, word2vec.CalcVecDiff(fromVec, )}
	//	distances = append(distances, *distance)
	//	//fmt.Printf("distance of %s: %f", k, word2vec.CalcVecDistance(fromVec, v))
	//}
	for _, word := range toWords {
		replace := strings.Replace(strings.TrimSpace(word), " ", "_", -1)
		toVec := vectorMap[strings.ToLower(replace)]
		if toVec == nil || len(toVec) == 0 {
			logger.Printf("%s not found\n", replace)
			continue
		}
		dist := word2vec.CalcVecDiff(fromVec, toVec)
		distance := &word2vec.WordDistance{fromWord, word, dist}
		distances = append(distances, *distance)
		fmt.Printf("distance of %s: %f\n", word, dist)
	}
	sort.Sort(distances)
	for _, wd := range distances {
		logger.Printf("distance of %s: %f", toWordMap[wd.WordB], wd.Distance)
	}
}

func TestLoad(t *testing.T) {
	word2vec.LoadModel("/Users/lyons/Downloads/glove.6B/glove.6B.300d.txt")
}

func TestIteye(t *testing.T) {
	client.Init()
	iteye.CrawNews()
}

func TestIP(t *testing.T) {
	client.Init()
	conn := client.REDIS.Get()
	defer  conn.Close()
	//
	_, body, err := gorequest.New().Proxy("http://203.198.193.3:808").Get("http://1212.ip138.com/ic.asp").End()

	if err != nil {
		contains := strings.Contains(err[0].Error(), "proxyconnect")
		fmt.Println(contains)
		fmt.Println(err[0].Error())
	}

	//doc, _ := goquery.NewDocumentFromResponse(resp)

	//fmt.Println(doc.Find(".col-md-12").Text())

	fmt.Println(body)

	//p := func(_ *http.Request) (*url.URL, error) {
	//	return url.Parse("https://183.151.40.21:808")//根据定义Proxy func(*Request) (*url.URL, error)这里要返回url.URL
	//}

	//proxyUrl, err := url.Parse("http://45.76.154.255:8080")
	//myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}



	//resp, err := myClient.Get("http://1212.ip138.com/ic.asp")

	//transport := &http.Transport{Proxy: p}
	//client1 := &http.Client{Transport: transport}
	//resp, err := client1.Get("http://1212.ip138.com/ic.asp") //请求并获取到对象,使用代理
	//if err != nil {
	//	fmt.Println(err)
	//	panic(err)
	//}
	//all, _ := ioutil.ReadAll(resp.Body)
	//
	//fmt.Println(string(all))
}

func TestStr(t *testing.T) {
	str := "://"
	fmt.Println(str)
}

func TestFreeProxyList(t *testing.T) {
	client.Init()
	fmt.Println(proxy.GetFreeProxy("Hong"))
}

func TestFindRateLimit(t *testing.T) {
	client.Init()
	iteye.FindRateLimit()
}