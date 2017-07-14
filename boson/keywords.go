package boson

import (
	"github.com/parnurzeal/gorequest"
	"goJobSkills/log"
	"io/ioutil"
	"encoding/json"
	"strings"
	"path/filepath"
)

var KEYWORDS_URL = "http://api.bosonnlp.com/keywords/analysis"
var APP_KEY = "8SD_XL0x.16636.9U69N1EbgBqo"

var logger = log.GetLogger()

func GetKeyWords(filePath string)  {

	text := "\"text3 年以上 Java 开发经验，熟悉 MVC 编程思想，熟练掌握 Spring 框架；熟悉XML、JSON 等的技术规范；精通数据库设计，熟练使用 SQL 语言，熟悉常用的关系型数据库产品如 Mysql；熟悉服务器端缓存系统，如 Memcached、Redis等；熟练使用 Linux 相关命令，Linux 服务器下相关软件的安装和环境搭建；有手机端应用接口开发经验 具有互联网电商，社交等项目开发经验者优先考虑；有解决大数据量、高并发经验者优先考虑。\""
	response, _, errs := gorequest.New().
		Post(KEYWORDS_URL).
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Token", APP_KEY).
		Send(text).
		End()
	if errs != nil {
		logger.Println(errs)
		return
	}
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Println(err)
		return
	}

	//keywordMaps := make(map[float32]interface{})
	keywords := make([][]interface{}, 0)
	err = json.Unmarshal(all, &keywords)
	if err != nil {
		logger.Println(err)
		return
	}

	keywordMap := make(map[string]float64)
	for _, keyword := range keywords{
		str, ok := keyword[1].(string)
		if !ok {
			logger.Printf("it's not ok when asserting %s", keyword[1])
			continue
		}
		value, ok := keyword[0].(float64)
		if !ok {
			logger.Printf("it's not ok when asserting %f", keyword[0])
			continue
		}
		keywordMap[str] = value
	}
	removeStopwords(keywordMap)
	logger.Println(keywordMap)

}

func removeStopwords(keywordMap map[string]float64)   {
	path, _ := filepath.Abs("../stopwords")

	file, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Println("stopword file open failed, abort.", err)
		return
	}
	stopwords := string(file)

	for k ,_ := range keywordMap{
		if strings.Contains(stopwords, k){
			delete(keywordMap, k)
		}
	}
}

//type keyword struct {
//	weight float32
//}