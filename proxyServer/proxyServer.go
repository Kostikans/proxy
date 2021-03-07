package proxyServer

import (
	"crypto/tls"
	"fmt"
	"github.com/Kostikans/proxy/saver/repositorySave"
	"github.com/Kostikans/proxy/saver/repositorySave/repositoryImpl"
	"github.com/Kostikans/proxy/webInterface/models"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MyProxyServer struct {
	Server *http.Server
	saver   repositorySave.HttpRequestSaver
}

func NewMyProxyServer(saverImpl *repositoryImpl.ProxyRepo) *MyProxyServer {
	return &MyProxyServer{Server: &http.Server{
		Addr: ":80",
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler))},
		saver: saverImpl,
	}
}

func (proxy* MyProxyServer) InitHandler(){
	proxy.Server.Handler = 	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.HandleHTTP(w, r)
	})
}

func (proxy* MyProxyServer) HandleHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("fsd")
	delete(req.Header, "Proxy-Connection")
	req.RequestURI = ""

	info,err := proxy.GetProxyInfo(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = proxy.saver.SaveRequest(info)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := &http.Client{}
	resp,err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.Header().Set("Server", "nginx/1.14.1")
	w.Header().Set("Connection", "close")
	w.Header().Set("Location", GetRedirectUrl(req))
	w.WriteHeader(http.StatusMovedPermanently)
	w.Write([]byte("<html>\n<head><title>301 Moved Permanently</title></head>" +
		"\n<body bgcolor=\"white\">\n<center><h1>301 Moved Permanently</h1></center>\n<hr><center>nginx/1.14.1</center>\n</body>\n</html>\n"))

}


func (proxy* MyProxyServer) GetProxyInfo (r *http.Request) (models.ProxyInfo, error) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return models.ProxyInfo{}, err
	}

	info := models.ProxyInfo{
		Url:                r.URL.String(),
		HeaderInfo:         r.Header,
		Method:             r.Method,
		RequestBody:        requestBody,
	}
	return info, nil
}

func GetRedirectUrl(r *http.Request) string {
	var redirectedUrl string
	redirectedUrl = r.URL.Query().Get("url")
	if redirectedUrl == "" {
		redirectedUrl = r.RequestURI
	}
	redirectedUrl = strings.ReplaceAll(redirectedUrl, "%", "")
	return redirectedUrl
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
