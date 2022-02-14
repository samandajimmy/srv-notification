package contract

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/model"
)

type ClientConfigure interface {
	FindClientConfigByXID(xid string) (*model.ClientConfigVO, error)
	FindByKey(key string, appId int64) (*model.ClientConfig, error)
	FindClientConfig(params *dto.FindOptions) (*model.ClientConfigSearchResult, error)
	InsertClientConfig(row model.ClientConfig) error
	UpdateClientConfig(row *model.ClientConfig) error
	DeleteClientConfigById(id int64) error
}
