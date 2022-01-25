package notification

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

func (rc *RepositoryContext) InsertNotification(row model.Notification) error {
	_, err := rc.RepositoryStatement.Notification.Insert.ExecContext(rc.ctx, row)
	return err
}

func (rc *RepositoryContext) FindNotificationByID(id string) (*model.Notification, error) {
	var notification model.Notification
	err := rc.Notification.FindByID.GetContext(rc.ctx, &notification, id)
	return &notification, err
}

func (rc *RepositoryContext) UpdateNotificationByID(row *model.Notification) error {
	result, err := rc.RepositoryStatement.Notification.Update.ExecContext(rc.ctx, row)
	if err != nil {
		return err
	}
	return nsql.IsUpdated(result)
}

func (rc *RepositoryContext) DeleteNotificationByID(id string) error {
	_, err := rc.Notification.DeleteByID.ExecContext(rc.ctx, id)
	return err
}
