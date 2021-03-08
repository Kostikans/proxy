package main

import (
	"github.com/Kostikans/proxy/proxyServer"
	"github.com/Kostikans/proxy/saver/repositorySave/repositoryImpl"
	"github.com/Kostikans/proxy/webInterface"
	"github.com/Kostikans/proxy/webInterface/repositoryInterface/repositoryInterfaceImpl"
	"github.com/Kostikans/proxy/webInterface/usecaseInterfaceImpl"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
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

func main() {
	pool := newPool(":6379")
	conn := pool.Get()

	repo := repositoryImpl.NewProxyRepo(conn)
	myProxy := proxyServer.NewMyProxyServer(repo, ":8080")
	myProxy.InitHandler()

	r := mux.NewRouter()
	interfaceRepo := repositoryInterfaceImpl.NewInterfaceRep(conn)
	interfaceUsecase := usecaseInterfaceImpl.NewInterfaceUsecase(interfaceRepo)
	webInterface.InitHandler(r, myProxy, interfaceUsecase)
	go http.ListenAndServe(":8000", r)

	myProxy.Server.ListenAndServe()
	//myProxyHttps := proxyServer.NewMyProxyServer(repo,":8081")
	//myProxyHttps.InitHandler()
	//myProxyHttps.Server.ListenAndServeTLS("ca.crt", "cert.key")

}
