package webInterface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Kostikans/proxy/proxyServer"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Handler struct {
	proxy   *proxyServer.MyProxyServer
	usecase UsecaseInterface
}

func InitHandler(r *mux.Router, proxy *proxyServer.MyProxyServer, usecaseInterface UsecaseInterface) *Handler {
	handler := &Handler{proxy: proxy, usecase: usecaseInterface}
	r.HandleFunc("/requests", handler.HandlerListRequests).Methods("GET")
	r.HandleFunc("/requests/{id}", handler.HandlerGetRequest).Methods("GET")
	r.HandleFunc("/repeat/{id}", handler.HandlerRepeatRequest).Methods("GET")
	r.HandleFunc("/scan/{id}", handler.HandlerScanRequest).Methods("GET")
	return handler
}

func (h *Handler) HandlerListRequests(w http.ResponseWriter, req *http.Request) {
	list, err := h.usecase.GetListRequests()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(list)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) HandlerGetRequest(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	request, err := h.usecase.GetRequest(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	err = json.NewEncoder(w).Encode(request)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) HandlerRepeatRequest(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	request, err := h.usecase.GetRequest(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	req, err = http.NewRequest(request.Info.Method, request.Info.Url, bytes.NewReader(request.Info.RequestBody))
	if err != nil {
		log.Println(err)
	}
	req.Header = request.Info.HeaderInfo
	h.proxy.HandleHTTP(w, req)
}

func (h *Handler) HandlerScanRequest(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	request, err := h.usecase.GetRequest(id)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	url := h.usecase.GetXXSInjectionUrl(request)
	req, err = http.NewRequest(request.Info.Method, url, bytes.NewReader(request.Info.RequestBody))
	if err != nil {
		log.Println(err)
	}
	req.Header = request.Info.HeaderInfo
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	injectionString := `vulnerable'"><img src onerror=alert()>`
	if strings.Contains(string(bodyBytes), injectionString) {
		w.Write([]byte(fmt.Sprintf("Have XSS injection in %s", request.Info.Url)))
		return
	}

	w.Write([]byte(fmt.Sprintf("No XSS injection in %s", request.Info.Url)))
}
