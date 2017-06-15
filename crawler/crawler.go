package crawler

import (
	"github.com/hu17889/go_spider/core/spider"
	"github.com/hu17889/go_spider/core/pipeline"
	"GoJobSkills/website/lagou"
	"github.com/hu17889/go_spider/core/common/request"
)

type Centipede struct {

}
func (crawler *Centipede) Run() {
	lagouRequest := request.NewRequest(
		"https://www.lagou.com/jobs/positionAjax.json?needAddtionalResult=false",
		"json",
		"lagou",
		"POST",
		"first=true&pn=1&kd=ui",
		nil,
		nil,
		nil,
		nil)
	spider.NewSpider(new(lagou.LaGouPageProcessor), "TaskName").
		AddRequest(lagouRequest).
		AddPipeline(pipeline.NewPipelineConsole()).                    // Print result on screen
		SetThreadnum(1).                                               // Crawl request by three Coroutines
		Run()
}

func NewCentipede() *Centipede {
	return &Centipede{};
}


