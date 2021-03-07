package repositoryInterfaceImpl

import (
	"encoding/json"
	"errors"
	"github.com/Kostikans/proxy/webInterface/models"
	"github.com/gomodule/redigo/redis"
)

type InterfaceRepo struct {
	RedisClient redis.Conn
}

func NewInterfaceRep(client redis.Conn) *InterfaceRepo{
	return &InterfaceRepo{RedisClient: client}
}

func (interfaceRepo *InterfaceRepo) GetListRequests() ([]models.Request,error){
	var requests []models.Request

	listInterface,err := interfaceRepo.RedisClient.Do("LRANGE","1","1", "-1")
	v, err := redis.Values(listInterface,err)

	if err != nil {
		return requests,err
	}

	for _,el := range v {
		var request models.Request
		if err := json.Unmarshal(el.([]byte), &request); err != nil {
			return requests, err
		}
		requests = append(requests, request)
	}

	return requests,nil
}

func (interfaceRepo *InterfaceRepo) GetRequest(ID string) (models.Request,error){
	var requests []models.Request
	listInterface,err := interfaceRepo.RedisClient.Do("LRANGE","1","1", "-1")
	v, err := redis.Values(listInterface,err)

	if err != nil {
		return models.Request{},err
	}

	for _,el := range v {
		var request models.Request
		if err := json.Unmarshal(el.([]byte), &request); err != nil {
			return models.Request{}, err
		}
		requests = append(requests, request)
	}
	for _,request := range requests {
		if request.ID == ID {
			return request,nil
		}
	}
	return models.Request{},errors.New("id doesn't exist")

}

