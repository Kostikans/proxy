package repositoryImpl

import (
	"encoding/json"
	"github.com/Kostikans/proxy/webInterface/models"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

type ProxyRepo struct {
	RedisClient redis.Conn
}

func NewProxyRepo(client redis.Conn) *ProxyRepo{
	return &ProxyRepo{RedisClient: client}
}

func(pr *ProxyRepo) SaveRequest(info models.ProxyInfo) error {
	key := uuid.NewV4()
	request := models.Request{ID: key.String(),Info: info}
	requestRawBinary,err  := json.Marshal(&request)
	if err != nil {
		return err
	}
	_,err = pr.RedisClient.Do("RPUSH", "1", requestRawBinary)
	if err != nil {
		return err
	}
	return nil
}
