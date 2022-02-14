package notification

import (
	"fmt"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/statement"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	nsqlDep "repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"strings"
)

var notificationFilters = map[string]nsql.FilterParser{
	constant.UserRefIdKey:     newEqualFilter(statement.NotificationSchema, "userRefId"),
	constant.ApplicationIdKey: newEqualFilter(statement.NotificationSchema, "applicationId"),
}

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
	result, err := rc.RepositoryStatement.Notification.Update.ExecContext(rc.ctx, row.UpdatedAt, row.ModifiedBy,
		row.Version, row.IsRead, row.ReadAt, row.ID)
	if err != nil {
		return err
	}
	return nsqlDep.IsUpdated(result)
}

func (rc *RepositoryContext) DeleteNotificationByID(id string) error {
	_, err := rc.Notification.DeleteByID.ExecContext(rc.ctx, id)
	return err
}

func (rc *RepositoryContext) CountNotification(options *dto.GetCountNotification) (int64, error) {
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

func (rc *RepositoryContext) FindNotification(params *dto.ListPayload) (*model.FindNotificationResult, error) {
	// Init query builder
	b := query.From(statement.NotificationSchema)

	// Set where
	filters := query.NewFilter(params.Filters, notificationFilters)
	b.Where(filters.Conditions())

	// Set order by
	switch params.SortBy {
	case constant.SortByCreated:
		b.OrderBy("createdAt")
	default:
		params.SortBy = constant.SortByLastCreated
		b.OrderBy("createdAt", option.SortDirection(op.Descending))
	}

	// Create select query
	selectQuery := b.Select(
		query.Column("*", option.Schema(statement.NotificationSchema))).
		Limit(params.Limit).Skip(params.Skip).Build()
	selectQuery = rc.conn.Rebind(selectQuery)

	// Create count query
	b.ResetOrderBy().ResetSkip().ResetLimit()
	countQuery := b.Select(
		query.Count("*", option.Schema(statement.NotificationSchema), option.As("count"))).
		Build()
	countQuery = rc.conn.Rebind(countQuery)

	// Get args
	args := filters.Args()

	// Execute query
	var rows []model.Notification
	err := rc.conn.SelectContext(rc.ctx, &rows, selectQuery, args...)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	var count int64
	err = rc.conn.GetContext(rc.ctx, &count, countQuery, args...)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Prepare result
	result := model.FindNotificationResult{
		Rows:  rows,
		Count: count,
	}
	return &result, err
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
