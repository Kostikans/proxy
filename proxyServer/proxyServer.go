package proxyServer

import (
	"github.com/Kostikans/proxy/saver/repositorySave"
	"github.com/Kostikans/proxy/saver/repositorySave/repositoryImpl"
	"github.com/Kostikans/proxy/webInterface/models"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type MyProxyServer struct {
	Server *http.Server
	saver  repositorySave.HttpRequestSaver
	Client http.Client
}

func NewMyProxyServer(saverImpl *repositoryImpl.ProxyRepo, addres string) *MyProxyServer {
	return &MyProxyServer{Server: &http.Server{
		Addr: addres},
		saver:  saverImpl,
		Client: http.Client{Timeout: 5 * time.Second},
	}
}

func (proxy *MyProxyServer) InitHandler() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.HandleHTTP(w, r)
	})
	proxy.Server.Handler = r
}

func (proxy *MyProxyServer) HandleHTTP(w http.ResponseWriter, req *http.Request) {

	info, err := proxy.GetProxyInfo(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = proxy.saver.SaveRequest(info)
	if err != nil {
		log.Fatal(err)
		return
	}

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (proxy *MyProxyServer) CopyRequest(r *http.Request) *http.Request {
	redirectedRequest, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	delete(r.Header, "Proxy-Connection")
	redirectedRequest.Header = r.Header
	for _, cookie := range r.Cookies() {
		redirectedRequest.AddCookie(cookie)
	}
	return redirectedRequest
}

func (proxy *MyProxyServer) GetProxyInfo(r *http.Request) (models.ProxyInfo, error) {
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return models.ProxyInfo{}, err
	}

	info := models.ProxyInfo{
		Url:            r.URL.String(),
		HeaderInfo:     r.Header,
		Method:         r.Method,
		RequestBody:    requestBody,
		RequestCookies: r.Cookies(),
	}
	return info, nil
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
