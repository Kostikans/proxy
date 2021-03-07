package usecaseInterfaceImpl

import (
	"github.com/Kostikans/proxy/webInterface"
	"github.com/Kostikans/proxy/webInterface/models"
	"net/url"
)

type InterfaceUsecase struct {
	Repository webInterface.RepositoryInterface
}

func NewInterfaceUsecase(Repository webInterface.RepositoryInterface) *InterfaceUsecase {
	return &InterfaceUsecase{Repository: Repository}
}

func (interfaceUsecase *InterfaceUsecase) GetListRequests() ([]models.Request, error) {
	return interfaceUsecase.Repository.GetListRequests()
}

func (interfaceUsecase *InterfaceUsecase) GetRequest(ID string) (models.Request, error) {
	return interfaceUsecase.Repository.GetRequest(ID)
}

func (interfaceUsecase *InterfaceUsecase) GetXXSInjectionUrl(req models.Request) string {
	u, _ := url.Parse(req.Info.Url)

	injectionString := `vulnerable'"><img src onerror=alert()>`
	values, _ := url.ParseQuery(u.RawQuery)
	for key, val := range values {
		values.Set(key, injectionString+val[0])
	}
	u.RawQuery = values.Encode()
	return u.String()
}
