// Package worker @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/10/20 11:14 下午
package worker_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"spider/internal/scheduler/worker"
)

func TestWork(t *testing.T) {
	//urlStr:="http://metalsucks.net"
	urlStr:="https://juejin.cn/user/3861140565404264/followers"
	info, urls, err := worker.Work(urlStr)
	if err != nil {
		t.Error(err)
	}
	log.Println(info,urls)
}

func TestDemo(t *testing.T)  {
	html := `<body>

				<div>DIV1</div>
				<div>DIV2</div>
				<span>SPAN</span>

			</body>
			`

	dom,err:=goquery.NewDocumentFromReader(strings.NewReader(html))
	if err!=nil{
		log.Fatalln(err)
	}

	dom.Find("span").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
}