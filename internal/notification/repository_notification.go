package notification

import (
	"fmt"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"strings"
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

func (rc *RepositoryContext) CountNotification(options dto.GetCountNotification) (int64, error) {
	// Prepare where
	var args []interface{}
	var whereQuery []string

	if options.UserRefId > 0 {
		whereQuery = append(whereQuery, `"userRefId" = ?`)
		args = append(args, options.UserRefId)
	}

	if options.Application.ID > 0 {
		whereQuery = append(whereQuery, `"applicationId" = ?`)
		args = append(args, options.Application.ID)
	}

	where := ""
	if len(whereQuery) > 0 {
		where = "WHERE " + strings.Join(whereQuery, " AND ")
	}

	// Prepare query
	queryCount := fmt.Sprintf(`SELECT COUNT(id) FROM "Notification" %s`, where)
	queryCount = rc.conn.Rebind(queryCount)
	var count int64
	err := rc.conn.GetContext(rc.ctx, &count, queryCount, args...)
	if err != nil {
		return 0, ncore.TraceError(err)
	}

	return count, nil
}
