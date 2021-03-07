package models

import "net/http"

type ProxyInfo struct {
	Url                string `redis:"Url"`
	HeaderInfo         http.Header `redis:"HeaderInfo"`
	Method             string `redis:"Method"`
	RequestBody        []byte `redis:"RequestBody"`
}

type Request struct {
	ID string   `redis:"ID"`
	Info ProxyInfo `redis:"Info"`
}
