package notification

import (
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
)

// Define filters
var applicationFilters = map[string]nsql.FilterParser{
	constant.NameKey:         query.LikeFilter("name", op.LikeSubString, option.Schema(statement.ApplicationSchema)),
	constant.CreatedFromKey:  newGreaterThanEqualFilter(statement.ApplicationSchema, "createdAt"),
	constant.CreatedUntilKey: newLessThanEqualFilter(statement.ApplicationSchema, "createdAt"),
}

func (rc *RepositoryContext) InsertApplication(row model.Application) error {
	_, err := rc.RepositoryStatement.Application.Insert.ExecContext(rc.ctx, row)
	return err
}

func (rc *RepositoryContext) FindApplicationByXID(xid string) (*model.Application, error) {
	var application model.Application
	err := rc.RepositoryStatement.Application.FindByXID.GetContext(rc.ctx, &application, xid)
	return &application, err
}

func (rc *RepositoryContext) DeleteApplicationById(id int64) error {
	_, err := rc.RepositoryStatement.Application.DeleteByID.ExecContext(rc.ctx, id)
	return err
}

func (rc *RepositoryContext) FindApplication(params *dto.ListPayload) (*model.ApplicationListResult, error) {
	// Init query builder
	b := query.From(statement.ApplicationSchema)

	// Set where
	filters := query.NewFilter(params.Filters, applicationFilters)
	b.Where(filters.Conditions())

	// Set order by
	switch params.SortBy {
	case constant.SortByName:
		b.OrderBy("name")
	case constant.SortByNameDesc:
		b.OrderBy("name", option.SortDirection(op.Descending))
	case constant.SortByCreated:
		b.OrderBy("createdAt")
	default:
		params.SortBy = constant.SortByLastUpdated
		b.OrderBy("updatedAt", option.SortDirection(op.Descending))
	}

	// Create select query
	selectQuery := b.Select(
		query.Column("*", option.Schema(statement.ApplicationSchema))).
		Limit(params.Limit).Skip(params.Skip).Build()
	selectQuery = rc.conn.Rebind(selectQuery)

	// Create count query
	b.ResetOrderBy().ResetSkip().ResetLimit()
	countQuery := b.Select(
		query.Count("*", option.Schema(statement.ClientConfigSchema), option.As("count"))).
		Build()
	countQuery = rc.conn.Rebind(countQuery)

	// Get args
	args := filters.Args()

	// Execute query
	var rows []model.Application
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
	result := model.ApplicationListResult{
		Rows:  rows,
		Count: count,
	}
	return &result, err
}

func (rc *RepositoryContext) UpdateApplication(row *model.Application, currentVersion int64) error {
	// "updatedAt", "modifiedBy", "version", "name", "apiKey", "webhookUrl"
	result, err := rc.Application.Update.ExecContext(rc.ctx, row.UpdatedAt, row.ModifiedBy, row.Version,
		row.Name, row.ApiKey, row.WebhookURL, row.ID, currentVersion)
	if err != nil {
		return err
	}
	return nsqlDep.IsUpdated(result)
}
