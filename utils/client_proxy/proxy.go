package client_proxy

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
)

// NewReqClient 使用单例模式创建一个请求客户端,向 reddit 发送请求, 由于 reddit 的外网环境,需要使用本机代理
func NewReqClient() *http.Client {
	var ProxyURI, _ = url.Parse("http://127.0.0.1:41091")
	
	var transport = &http.Transport{
		Proxy:           http.ProxyURL(ProxyURI),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	
	// 使用这种办法可以在启动应用程序时以配置环境变量的方式启动程序.
	// The environment values may be either a complete URL or a "host[:port]", in which case the "http" scheme is assumed.
	// The schemes "http", "https", and "socks5" are supported. An error is returned if the value is a different form.
	//var transport = &http.Transport{
	//	Proxy:           http.ProxyFromEnvironment,
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	
	var once sync.Once
	var reqClient *http.Client
	once.Do(func() {
		reqClient = &http.Client{
			Transport: transport,
		}
	})
	return reqClient
}
