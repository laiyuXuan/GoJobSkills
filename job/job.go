package job

import (
	"github.com/robfig/cron"
	"GoJobSkills/log"
	"GoJobSkills/proxy"
)


var logger = log.GetLogger()


func Register()  {
	cron := cron.New()

	cron.AddFunc("0 0 * * ? *", ProxyPoolMaintainJob)

	cron.Start()
}

func ProxyPoolMaintainJob() {
	logger.Println("ProxyPoolMaintainJob started")

	proxy.FillProxyPool()
}

