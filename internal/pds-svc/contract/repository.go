package contract

import "repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"

type ClientConfigure interface {
	FindByKey(key string, appId int) (*model.ClientConfig, error)
}
