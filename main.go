package main

import (
	"github.com/comeonjy/library/conf"
	"github.com/comeonjy/library/iredis"

	"spider/engine"
	"spider/scheduler"
)

const (
	REDIS_URLS    = "redis_urls"
	REDIS_DO_LIST = "redis_do_list"
)

func main() {
	e:=&engine.Engine{
		ScanUrl:   "https://book.douban.com/subject/34907964/?icn=index-latestbook-subject",
		WorkerNum: 10,
		Scheduler: &scheduler.QueuedScheduler{},
		Cache:     iredis.New(&conf.RedisConf{Address: "localhost:6379"}),
	}
	e.Run()
}


