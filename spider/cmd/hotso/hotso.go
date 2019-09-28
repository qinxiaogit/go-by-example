package main

import (
	"fmt"
	"sync"
	"github.com/gocolly/colly"
	)

type Spider struct {
	Type int
}

var wg *sync.WaitGroup

var userAgent = "Chrome: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"

//OnWeiBo

func (s *Spider)OnWeiBo()[]map[string]interface{}{
	url := "https://s.weibo.com/top/summary"
	var allData []map[string]interface{}

	c := colly.NewCollector(colly.MaxDepth(1),colly.UserAgent(userAgent))
	c.OnError(func(response *colly.Response, e error) {
		fmt.Println("Request  URL:",response.Request.URL,"failed with response:",response, "\nError:", err)
	})

	c.OnHTML("#pl_top_realtimehot>table>tbody", func(e *colly.HTMLElement) {
		e.ForEach("tbody>tr", func(i int, element *colly.HTMLElement) {
			top := element.ChildText("td .td-01 .ranktop")
			title := element.ChildText("td .td-02 >a")
			reading := element.ChildText("td .td-02 > span")
			state := element.ChildText("td .td-03>i")

			var url = ""
			if state == "荐"{
				url = element.ChildAttr("td .td-02 >a","href_to")
			}else{
				url = element.ChildAttr("td .td-02>a","href")
			}
			allData = append(allData, map[string]interface{}{"top":top,"title": title, "reading": reading, "url": "https://s.weibo.com" + url, "state": state})
		})
	})
	c.Visit(url)
	return allData
}