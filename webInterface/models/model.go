package models

import "net/http"

type ProxyInfo struct {
	Url            string         `redis:"Url"`
	HeaderInfo     http.Header    `redis:"HeaderInfo"`
	Method         string         `redis:"Method"`
	RequestBody    []byte         `redis:"RequestBody"`
	RequestCookies []*http.Cookie `redis:"RequestCookies"`
}

type Request struct {
	ID   string    `redis:"ID"`
	Info ProxyInfo `redis:"Info"`
}
