package fetch

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"

	"spider/engine"
	"spider/types"
)

func FetchSubject(all string) (engine.Result,error) {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(all))
	if err != nil {
		log.Error(err)
	}
	resource:=make([]engine.Resource,0)
	dom.Find("#db-rec-section dd > a").Each(func(i int, selection *goquery.Selection) {
		if url, exists := selection.Attr("href");exists{
			resource=append(resource,engine.Resource{
				Url:       url,
				FetchFunc: FetchSubject,
			} )
		}
	})
	average := dom.Find(".rating_self > strong").Text()
	num := dom.Find(".rating_sum span a span").Text()
	title := dom.Find("#wrapper h1 span").Text()
	fmt.Println(title, average, num)
	book:=types.BookInfo{
		Name:  title,
		Score: average,
		Num:   num,
	}
	marshal, _ := json.Marshal(book)
	item:=make([]engine.Item,0)
	item=append(item, engine.Item{
		Type:   reflect.TypeOf(book),
		Source: string(marshal),
	})
	return engine.Result{
		Resources: nil,
		Items:     item,
	},nil
}