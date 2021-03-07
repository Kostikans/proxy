package usecaseInterfaceImpl

import (
	"github.com/Kostikans/proxy/webInterface"
	"github.com/Kostikans/proxy/webInterface/models"
)

type InterfaceUsecase struct {
	Repository  webInterface.RepositoryInterface
}

func NewInterfaceUsecase(Repository  webInterface.RepositoryInterface) *InterfaceUsecase{
	return &InterfaceUsecase{Repository: Repository}
}


func (interfaceUsecase *InterfaceUsecase) GetListRequests() ([]models.Request,error){
	return interfaceUsecase.Repository.GetListRequests()
}

func (interfaceUsecase *InterfaceUsecase) GetRequest(ID string) (models.Request,error){
	return interfaceUsecase.Repository.GetRequest(ID)
}
