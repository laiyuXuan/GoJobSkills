package job

import (
	timer "github.com/robfig/cron"
	"goJobSkills/log"
	"goJobSkills/proxy"
)


var logger = log.GetLogger()


func Register()  {
	cron := timer.New()

	cron.AddFunc("0 0 0 * ? *", proxyPoolMaintainJob)

	cron.Start()
}

func proxyPoolMaintainJob() {
	logger.Println("ProxyPoolMaintainJob started")

	proxy.CheckAvailablity()
	proxy.FillProxyPool()
}
