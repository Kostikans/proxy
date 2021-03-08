package proxyServer

import (
	"fmt"
	"github.com/Kostikans/proxy/saver/repositorySave"
	"github.com/Kostikans/proxy/saver/repositorySave/repositoryImpl"
	"github.com/Kostikans/proxy/webInterface/models"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net"
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
		if r.Method == http.MethodConnect {
			proxy.HandleTunneling(w, r)
		} else {
			proxy.HandleHTTP(w, r)
		}
	})
	proxy.Server.Handler = r
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func (proxy *MyProxyServer) HandleTunneling(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func (proxy *MyProxyServer) HandleHTTP(w http.ResponseWriter, req *http.Request) {
	redirectRequest := proxy.CopyRequest(req)

	fmt.Println(redirectRequest)
	resp, err := proxy.Client.Do(redirectRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(http.StatusMovedPermanently)
	w.Header().Set("Server", "nginx/1.14.1")
	w.Header().Set("Connection", "close")
	w.Header().Set("Location", req.URL.String())

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

	w.Write([]byte("<html>\n<head><title>301 Moved Permanently</title></head>" +
		"\n<body bgcolor=\"white\">\n<center><h1>301 Moved Permanently</h1></center>\n<hr><center>nginx/1.14.1</center>\n</body>\n</html>\n"))

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
