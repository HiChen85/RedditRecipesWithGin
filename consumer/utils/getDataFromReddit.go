package utils

import (
	"encoding/xml"
	"github.com/HiChen85/RedditRecipesWithGin/consumer/utils/client_proxy"
	"io/ioutil"
	"log"
	"net/http"
)

type Feed struct {
	Entries []Entry `xml:"entry"`
}

// 因为.rss 网站返回的是xml 数据,所以需要使用 xml 标签而非 json
type Entry struct {
	// recipe 地址
	Link struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	// 此处的 URL 是食谱的配图地址
	Thumbnail struct {
		URL string `xml:"url,attr"`
	} `xml:"thumbnail"`
	Title string `xml:"title"`
}

func GetDataFromReddit(url string) ([]Entry, error) {
	// 创建请求客户端
	reqClient := client_proxy.NewReqClient()
	
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// 添加请求头
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.81 Safari/537.36")
	rsp, err := reqClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rsp.Body.Close()
	dataBytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var feed Feed
	err = xml.Unmarshal(dataBytes, &feed)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("获取数据中...")
	//log.Println(feed.Entries)
	log.Println("获取完成...")
	return feed.Entries, nil
}
