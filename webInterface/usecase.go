package webInterface

import "github.com/Kostikans/proxy/webInterface/models"

type UsecaseInterface interface {
	GetListRequests() ([]models.Request, error)
	GetRequest(ID string) (models.Request, error)
	GetXXSInjectionUrl(request models.Request) string
}
