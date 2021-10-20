// Package worker @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/10/18 8:33 下午
package worker

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func Work(urlStr string) (data string, urls []string, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}

	resp, err := client.Get(urlStr)
	if err != nil {
		return "", nil, err
	}

	if resp.StatusCode != 200 {
		return "", nil, errors.New(fmt.Sprintf("status code error: %d %s", resp.StatusCode, resp.Status))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("meta[itemprop=url]").Each(func(i int, s *goquery.Selection) {
		if url,ok:=s.Attr("content");ok{
			urls=append(urls,url,url+"/followers")
		}
	})
	info:=doc.Find("h1[class=username]").Text()

	return info,urls,nil
}
