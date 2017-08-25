package boson

import (
	"github.com/parnurzeal/gorequest"
	"goJobSkills/log"
	"io/ioutil"
	"encoding/json"
	"strings"
	"path/filepath"
	"os"
	"fmt"
	"sort"
	"time"
)

const (
	KEYWORDS_URL          = "http://api.bosonnlp.com/keywords/analysis"
	CHARACTERIZING_URL 	= "http://api.bosonnlp.com/tag/analysis?space_mode=0&oov_level=4&t2s=0&&special_char_conv=0"

	APP_KEY	        	= "8SD_XL0x.16636.9U69N1EbgBqo"
	STOPWORDS_PATH 		= "/Users/Lyons/doc/stopwords/stopwords"

	THEARDHOLD 			= 0.1
)


var logger = log.GetLogger()


func GetKeywords(filePath string) WordWeights  {

	//file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	//if err != nil {
	//	logger.Panic(err)
	//	return
	//}
	//text, err := ioutil.ReadAll(file)
	//if err != nil {
	//	logger.Panic(err)
	//	return
	//}

	words := CharacterizeWords(filePath)
	stopwords := getStopWord();
	noStopWord := make([]string, 0)
	for _, word:= range words{
		if !strings.Contains(stopwords, word){
			noStopWord = append(noStopWord, word)
		}
	}

	response, body, errs := gorequest.New().
		Post(KEYWORDS_URL + "?top_k=100&segmented=true").
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Token", APP_KEY).
		//Send(string(text)).
		Send(strings.ToLower(strings.Join(noStopWord, " "))).
		Timeout(time.Second * 10).
		End()
	if errs != nil {
		logger.Panic(errs)
	}
	all, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Panic(err)
	}

	logger.Println(body)

	//keywordMaps := make(map[float32]interface{})
	keywords := make([][][]interface{}, 0)
	err = json.Unmarshal(all, &keywords)
	if err != nil {
		logger.Panic(err)
	}
	weights := make(WordWeights, 0)
	for _, keyword := range keywords[0] {
		str := keyword[0].(float64)
		str1 := keyword[1].(string)
		weight := &WordWeight{str1, str}
		weights = append(weights, *weight)
	}
	sort.Sort(weights)
	logger.Println(weights)
	return weights
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

func getStopWord() string {
	file, err := os.OpenFile(STOPWORDS_PATH, os.O_RDONLY, 0666)
	if err != nil {
		logger.Panic(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Panic(err)
	}
	return string(content)
}

func CharacterizeWords(filePath string) []string  {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		logger.Panic(err)
	}
	text, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Panic(err)
	}

	resp, body, errs := gorequest.New().
		Post(CHARACTERIZING_URL).
		Set("Content-Type", "application/json").
		Set("Accept", "application/json").
		Set("X-Token", APP_KEY).
	//Send(string(text)).
		Send(strings.Replace(string(text), "%", "", -1)).
		End()
	if errs != nil {
		logger.Panic(errs)
	}

	logger.Println(body)
	words := make([]CharacterizedWord, 0)

	all, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(all, &words)

	fmt.Println(len(words))
	fmt.Println(len(words[0].Tags))
	fmt.Println(len(words[0].Words))

	filteredWord := make([]string, 0)
	for _, word := range words{
		for idx, tag := range word.Tags {
			if strings.Contains(tag, "n"){
				filteredWord = append(filteredWord, word.Words[idx])
			}
		}
	}
	logger.Println(filteredWord)
	return filteredWord
}

type CharacterizedWord struct {
	Tags  []string 	`json:"tag"`
	Words []string 	`json:"word"`
}

type WordWeight struct {
	Word string
	Weight float64
}

type WordWeights []WordWeight

func (slice WordWeights) Len() int {
	return len(slice)
}

func (slice WordWeights) Less(i, j int) bool {
	return slice[i].Weight < slice[j].Weight;
}

func (slice WordWeights) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

