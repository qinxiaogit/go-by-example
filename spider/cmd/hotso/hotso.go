package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/qinxiaogit/go-by-example/spider/common"
	"github.com/qinxiaogit/go-by-example/spider/config"
	"github.com/qinxiaogit/go-by-example/spider/internal"
	"github.com/qinxiaogit/go-by-example/spider/internal/cloud"
	"github.com/qinxiaogit/go-by-example/spider/internal/metadata/hotso"
)

type Spider struct {
	Type int
}

var wg *sync.WaitGroup

var userAgent = "Chrome: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.113 Safari/537.36"

//OnWeiBo

func (s *Spider) OnWeiBo() []map[string]interface{} {
	url := "https://s.weibo.com/top/summary"
	var allData []map[string]interface{}

	c := colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(userAgent))
	c.OnError(func(response *colly.Response, e error) {
		fmt.Println("Request  URL:", response.Request.URL, "failed with response:", response, "\nError:", e)
	})

	c.OnHTML("#pl_top_realtimehot>table>tbody", func(e *colly.HTMLElement) {
		e.ForEach("tbody>tr", func(i int, element *colly.HTMLElement) {
			top := element.ChildText("td.td-01.ranktop")
			title := element.ChildText("td.td-02 > a")
			reading := element.ChildText("td.td-02 > span")
			state := element.ChildText("td.td-03 > i")

			var url = ""
			if state == "荐" {
				url = element.ChildAttr("td.td-02 >a", "href_to")
			} else {
				url = element.ChildAttr("td.td-02 >a", "href")
			}
			allData = append(allData, map[string]interface{}{"top": top, "title": title, "reading": reading, "url": "https://s.weibo.com" + url, "state": state})
		})
	})
	c.Visit(url)
	return allData
}

//OnBaiDu 。。。
func (s *Spider) OnBaiDu() []map[string]interface{} {
	url := "http://top.baidu.com/buzz?b=1&c=513&fr=topbuzz_b341_c513"
	var allData []map[string]interface{}

	c := colly.NewCollector(colly.MaxDepth(1), colly.UserAgent(userAgent))
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})
	c.OnHTML("#main > div.mainBody > div > table > tbody", func(e *colly.HTMLElement) {
		e.ForEach("tbody > tr", func(i int, ex *colly.HTMLElement) {
			top := ex.ChildText("td.first > span")
			if top != "" {
				title := ex.ChildText("td.keyword > a.list-title")
				reading := ex.ChildText("td.last > span")
				url := ex.ChildAttr("td.keyword > a.list-title", "href")
				state := "" //ex.ChildText("td.td-03 > i")
				allData = append(allData, map[string]interface{}{"top": top, "title": common.GBK2UTF8(title), "reading": reading, "url": url, "state": state})
			}
		})
	})
	c.Visit(url)
	return allData
}

//OnZhiHu 实时热点
func (s *Spider) OnZhiHu() []map[string]interface{} {

	//ZhiHuOnline ...
	type ZhiHuOnline struct {
		Cookie    string `json:"cookie"`
		UserAgent string `json:"user_agent"`
	}

	var allData []map[string]interface{}
	var success = true

	var zhihu ZhiHuOnline
	if webdavCli, err := cloud.Dial(config.GetConfig().WebDav.Host, config.GetConfig().WebDav.User, config.GetConfig().WebDav.Password); err != nil {
		fmt.Println("zhihu webdav dial error")
		success = false
	} else {
		remoteDir := strings.Replace(config.GetConfig().WebDav.RemoteDir, "\\", "/", -1)
		if remoteDir[len(remoteDir)-1:] != "/" {
			remoteDir = remoteDir + "/"
		}
		if body, err := webdavCli.Download(remoteDir + "zhihu.json"); err != nil {
			fmt.Println("zhihu webdav download error")
			success = false
		} else {
			json.Unmarshal(body, &zhihu)
		}
	}
	if success != true {
		return allData
	}

	c := colly.NewCollector(colly.UserAgent(zhihu.UserAgent), colly.MaxDepth(1))
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", zhihu.Cookie)
	})
	c.OnHTML("#TopstoryContent > div > div > div.HotList-list", func(e *colly.HTMLElement) {
		e.ForEach("div.HotList-list > section.HotItem", func(i int, ex *colly.HTMLElement) {
			top := ex.ChildText("div.HotItem-index > div.HotItem-rank")
			title := ex.ChildText("div.HotItem-content > a > h2.HotItem-title")
			hotread := ex.ChildText("div.HotItem-content > div.HotItem-metrics")
			var reading = 0
			var err error
			ss := strings.Fields(hotread)
			if len(ss) >= 2 {
				if index := strings.Index(hotread, "万"); index == -1 {
					if reading, err = strconv.Atoi(ss[0]); err != nil {
						fmt.Println("zhihu hotnum error")
					}
				} else {
					if reading, err = strconv.Atoi(ss[0]); err != nil {
						fmt.Println("zhihu hotnum error")
					} else {
						reading = reading * 10000
					}
				}
			}
			url := ex.ChildAttr("div.HotItem-content > a ", "href")
			state := ex.ChildText("div.HotItem-index > div.HotItem-label")

			allData = append(allData, map[string]interface{}{"top": top, "title": title, "reading": reading, "url": url, "state": state})
		})
	})
	c.Visit("http://www.zhihu.com/hot")

	return allData
}

func ProduceData(s *Spider) {
	defer wg.Done()
	reflectValue := reflect.ValueOf(s)
	fmt.Println("%s", hotso.HotSoType[s.Type])
	methodValue := reflectValue.MethodByName("On" + hotso.HotSoType[s.Type])
	methodFunc := methodValue.Call(nil)
	fmt.Printf("---- %v ---- ", methodFunc)
	originData := methodFunc[0].Interface().([]map[string]interface{})
	now := time.Now().Unix()

	if len(originData) > 0 {
		switch s.Type {
		case hotso.WEIBO:
			internal.NewMongoDB().OnWeiBoInsert(&hotso.HotData{Type: s.Type, Name: hotso.HotSoType[s.Type], InTime: now, Data: originData})
		case hotso.BAIDU:
			internal.NewMongoDB().OnBaiDuInsert(&hotso.HotData{Type: s.Type, Name: hotso.HotSoType[s.Type], InTime: now, Data: originData})
		case hotso.ZHIHU:
			internal.NewMongoDB().OnZhiHuInsert(&hotso.HotData{Type: s.Type, Name: hotso.HotSoType[s.Type], InTime: now, Data: originData})
		}
	} else {
		fmt.Println("originData nil")
	}
}

func main() {
	wg = &sync.WaitGroup{}
	wg.Add(1)
	if len(os.Args) > 1 {
		for _, v := range os.Args[1:] {
			if n, err := strconv.Atoi(v); err != nil {
				fmt.Println("strconv Atoi error")
			} else {
				wg.Add(1)
				s := &Spider{Type: n}
				go ProduceData(s)
			}
		}
	} else {
		// wg.Add(len(hotso.HotSoType))
		// for k, _ := range hotso.HotSoType {
		s := &Spider{Type: 2}
		ProduceData(s)
		// }
	}
	wg.Wait()
}
