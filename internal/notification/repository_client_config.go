package notification

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
)

func (rc *RepositoryContext) HasInitialized() bool {
	return true
}

func (rc *RepositoryContext) FindByKey(key string, appId int) (*model.ClientConfig, error) {
	var row model.ClientConfig
	err := rc.ClientConfig.FindByKey.Get(&row, key, appId)
	return &row, err
}
