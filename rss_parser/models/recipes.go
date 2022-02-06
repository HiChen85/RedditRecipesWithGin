package models

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
