package job

import (
	timer "github.com/robfig/cron"
	"goJobSkills/log"
	"goJobSkills/proxy"
	"github.com/golang/glog"
)


var logger = log.GetLogger()


func Register()  {
	cron := timer.New()

	cron.AddFunc("0 0 0 * ? *", proxyPoolMaintainJob)

	cron.AddFunc("0 0 11 ? * 7", lagouJDcrawlerJob)

	cron.Start()
}

func proxyPoolMaintainJob() {
	logger.Println("ProxyPoolMaintainJob starts")

	proxy.CheckAvailablity()
	proxy.FillProxyPool()
}

func lagouJDcrawlerJob()  {
	logger.Println("lagouJDcrawlerJob starts")


}
