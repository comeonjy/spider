package engine

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type Engine struct {
	ScanUrl   string
	WorkerNum int
	Scheduler Scheduler
	Cache     *redis.Client
}

// 调度器：负责调度worker
type Scheduler interface {
	WorkChan() chan Resource
	WorkReady(chan Resource)
	Worker(chan Result)
	Submit(Resource)
	Run()
}

type Result struct {
	Resources []Resource
	Items     []Item
}

type Item struct {
	Type   reflect.Type
	Source string
}

type Resource struct {
	Url       string
	FetchFunc func(string) (Result,error)
}

func (e *Engine) Run() {
	out := make(chan Result)

	for i := 0; i < e.WorkerNum; i++ {
		e.Scheduler.Worker(out)
	}

	e.Scheduler.Run()

	for  {
		result:=<-out
		for _,v:=range result.Items{
			fmt.Println(v.Source)
		}
		for _,v:= range result.Resources{
			e.Scheduler.Submit(v)
		}
	}

}



func Download(url string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(err)
		return "", err
	}
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Add("Accept-Encoding", "gzip")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Connection", "keep-alive")
	//req.Header.Add("Cookie","bid=IRgP3_9OokQ; douban-fav-remind=1; __gads=ID=5f1295ff15879cdd:T=1588055806:S=ALNI_MY5hYZKq_awUjDsbmlz8rVhcQvbiw; ll=\"118318\"; __utma=30149280.167317502.1588055806.1588055806.1588214405.2; __utmc=30149280; __utmz=30149280.1588214405.2.2.utmcsr=baidu|utmccn=(organic)|utmcmd=organic; gr_user_id=65802814-eeed-40a4-b220-3c96f0205041; gr_session_id_22c937bbd8ebd703f2d8e9445f7dfd03=c74ef1c1-3f7a-435d-84a3-cf3983216ce8; gr_cs1_c74ef1c1-3f7a-435d-84a3-cf3983216ce8=user_id%3A0; ap_v=0,6.0; __utma=81379588.787616364.1588214412.1588214412.1588214412.1; __utmc=81379588; __utmz=81379588.1588214412.1.1.utmcsr=douban.com|utmccn=(referral)|utmcmd=referral|utmcct=/; _pk_ref.100001.3ac3=%5B%22%22%2C%22%22%2C1588214412%2C%22https%3A%2F%2Fwww.douban.com%2F%22%5D; _pk_ses.100001.3ac3=*; gr_session_id_22c937bbd8ebd703f2d8e9445f7dfd03_c74ef1c1-3f7a-435d-84a3-cf3983216ce8=true; _vwo_uuid_v2=D7698B5DFD18AE1EC7AF2563DB6223E80|6a21ba2e6c2009517156bf6d4c6f1fff; __yadk_uid=5PO8vTKDbijstzmB3XfjSghis6jPcYBH; viewed=\"34995610\"; _pk_id.100001.3ac3=5c15645b6f7a8181.1588214412.1.1588214474.1588214412.; __utmb=30149280.7.9.1588214405; __utmb=81379588.5.10.1588214412")
	req.Header.Add("Host", "book.douban.com")
	req.Header.Add("Referer", "https://www.douban.com/")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	parseGzip, err := ParseGzip(all)
	if err != nil {
		return "", err
	}
	return string(parseGzip), err
}





func ParseGzip(data []byte) ([]byte, error) {
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, data)
	r, err := gzip.NewReader(b)
	if err != nil {
		log.Info("[ParseGzip] NewReader error: %v, maybe data is ungzip", err)
		return nil, err
	} else {
		defer r.Close()
		undatas, err := ioutil.ReadAll(r)
		if err != nil {
			log.Warn("[ParseGzip]  ioutil.ReadAll error: %v", err)
			return nil, err
		}
		return undatas, nil
	}
}
