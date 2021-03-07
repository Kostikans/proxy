package repositorySave

import (
	"github.com/Kostikans/proxy/webInterface/models"
)

type HttpRequestSaver interface {
	SaveRequest(info models.ProxyInfo) error
}
