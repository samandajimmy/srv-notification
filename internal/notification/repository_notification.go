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

func (rc *RepositoryContext) FindNotification(params *dto.NotificationFindOptions) (*model.FindNotificationResult, error) {
	// Prepare where
	var args []interface{}
	var whereQuery []string

	if userId, ok := params.Filters["userId"]; ok {
		whereQuery = append(whereQuery, `"userRefId" = ?`)
		args = append(args, userId)
	}

	if applicationId, ok := params.Filters["applicationId"]; ok {
		whereQuery = append(whereQuery, `"applicationId" = ?`)
		args = append(args, applicationId)
	}

	where := ""
	if len(whereQuery) > 0 {
		where = "WHERE " + strings.Join(whereQuery, " AND ")
	}

	// Prepare query
	columns := `"id","createdAt","updatedAt","modifiedBy","metadata","version","applicationId","userRefId","isRead","readAt","options"`
	from := `FROM "Notification"`
	queryList := fmt.Sprintf(`SELECT %s %s %s ORDER BY %s LIMIT %d OFFSET %d`,
		columns,
		from,
		where,
		rc.GetOrderByQuery(params.SortBy, params.SortDirection),
		params.Limit,
		params.Skip)

	// Execute query
	queryList = rc.conn.Rebind(queryList)
	var rows []model.Notification
	err := rc.conn.SelectContext(rc.ctx, &rows, queryList, args...)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Prepare query
	queryCount := fmt.Sprintf(`SELECT COUNT(id) %s %s`, from, where)
	queryCount = rc.conn.Rebind(queryCount)
	var count int64
	err = rc.conn.GetContext(rc.ctx, &count, queryCount, args...)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	return &model.FindNotificationResult{
		Rows:  rows,
		Count: count,
	}, nil
}

func (rc *RepositoryContext) CountNotificationNotRead(options dto.GetCountNotification) (int64, error) {
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
	queryCount := fmt.Sprintf(`SELECT COUNT(id) FROM "Notification" %s AND "readAt" IS NULL AND "isRead" = 'false'`, where)
	queryCount = rc.conn.Rebind(queryCount)
	var count int64
	err := rc.conn.GetContext(rc.ctx, &count, queryCount, args...)
	if err != nil {
		return 0, ncore.TraceError(err)
	}

	return count, nil
}
