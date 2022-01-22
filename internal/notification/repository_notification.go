package notification

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
)

func (rc *RepositoryContext) InsertNotification(row model.Notification) error {
	_, err := rc.RepositoryStatement.Notification.Insert.ExecContext(rc.ctx, row)
	return err
}
