package main

import (
	"flag"
	"github.com/Kostikans/proxy/proxyServer"
	"github.com/Kostikans/proxy/saver/repositorySave/repositoryImpl"
	"github.com/Kostikans/proxy/webInterface"
	"github.com/Kostikans/proxy/webInterface/repositoryInterface/repositoryInterfaceImpl"
	"github.com/Kostikans/proxy/webInterface/usecaseInterfaceImpl"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func newPool(server string) *redis.Pool {
	return &redis.Pool{

		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func main(){
	var pemPath string
	flag.StringVar(&pemPath, "pem", "server.pem", "path to pem file")

	var keyPath string
	flag.StringVar(&keyPath, "key", "server.key", "path to key file")

	var proto string
	flag.StringVar(&proto, "proto", "https", "Proxy protocol (http or https)")

	flag.Parse()
	if proto != "http" && proto != "https" {
		log.Fatal("Protocol must be either http or https")
	}
	pool := newPool(":6379")
	conn := pool.Get()

	repo := repositoryImpl.NewProxyRepo(conn)
	myProxy := proxyServer.NewMyProxyServer(repo)
	myProxy.InitHandler()

	r := mux.NewRouter()
	interfaceRepo := repositoryInterfaceImpl.NewInterfaceRep(conn)
	interfaceUsecase := usecaseInterfaceImpl.NewInterfaceUsecase(interfaceRepo)
	webInterface.InitHandler(r,myProxy,interfaceUsecase)
	go http.ListenAndServe(":8000",r)
	if proto == "http" {
		log.Fatal(myProxy.Server.ListenAndServe())
	} else {
		log.Fatal(myProxy.Server.ListenAndServeTLS(pemPath, keyPath))
	}

}