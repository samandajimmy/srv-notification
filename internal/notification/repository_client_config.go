package notification

import (
	"github.com/nbs-go/errx"
	"github.com/nbs-go/nsql"
	"github.com/nbs-go/nsql/op"
	"github.com/nbs-go/nsql/option"
	"github.com/nbs-go/nsql/pq/query"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/statement"
	nsqlDep "repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"strings"
)

// Define filters
var clientConfigFilters = map[string]nsql.FilterParser{
	constant.ApplicationXIDKey: newEqualFilter(statement.ApplicationSchema, "xid"),
	constant.KeyKey:            newEqualFilter(statement.ClientConfigSchema, "key"),
}

func (rc *RepositoryContext) HasInitialized() bool {
	return true
}

func (rc *RepositoryContext) FindByKey(key string, appId int64) (*model.ClientConfig, error) {
	var row model.ClientConfig
	err := rc.RepositoryStatement.ClientConfig.FindByKey.Get(&row, key, appId)
	return &row, err
}

func (rc *RepositoryContext) FindDefaultClientConfigByKey(key string) (*model.ClientConfig, error) {
	var row model.ClientConfigDetailed
	err := rc.RepositoryStatement.ClientConfig.FindDefaultByKey.Get(&row, constant.DefaultConfig, key)
	return row.ClientConfig, err
}

func (rc *RepositoryContext) FindClientConfigByXID(xid string) (*model.ClientConfigDetailed, error) {
	var row model.ClientConfigDetailed
	err := rc.RepositoryStatement.ClientConfig.FindJoinApplicationByXID.GetContext(rc.ctx, &row, xid)
	return &row, err
}

func (rc *RepositoryContext) InsertClientConfig(row model.ClientConfig) error {
	_, err := rc.RepositoryStatement.ClientConfig.Insert.Exec(row)
	return err
}

func (rc *RepositoryContext) FindClientConfig(params *dto.ListPayload) (*model.ClientConfigListResult, error) {
	// Init query builder
	b := query.From(statement.ClientConfigSchema, option.As("cc")).
		Join(statement.ApplicationSchema, query.Equal(query.Column("applicationId"), query.On("id")),
			option.As("a"), option.JoinMethod(op.InnerJoin))

	// Set where
	filters := query.NewFilter(params.Filters, clientConfigFilters)
	b.Where(filters.Conditions())

	// Set order by
	switch strings.ToLower(params.SortBy) {
	case constant.SortByCreated:
		b.OrderBy("createdAt")
	default:
		b.OrderBy("updatedAt", option.SortDirection(op.Descending))
	}

	// Create select query
	selectQuery := b.Select(
		query.Column("*", option.Schema(statement.ClientConfigSchema)),
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
	var rows []model.ClientConfigDetailed
	err := rc.conn.SelectContext(rc.ctx, &rows, selectQuery, args...)
	if err != nil {
		return nil, errx.Trace(err)
	}

	var count int64
	err = rc.conn.GetContext(rc.ctx, &count, countQuery, args...)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// Prepare result
	result := model.ClientConfigListResult{
		Rows:  rows,
		Count: count,
	}
	return &result, err
}

func (rc *RepositoryContext) UpdateClientConfig(row *model.ClientConfig, currentVersion int64) error {
	result, err := rc.RepositoryStatement.ClientConfig.UpdateByID.ExecContext(rc.ctx,
		row.UpdatedAt, row.ModifiedBy, row.Version, row.Value, row.ID, currentVersion)
	if err != nil {
		return err
	}
	return nsqlDep.IsUpdated(result)
}

func (rc *RepositoryContext) DeleteClientConfigById(id int64) error {
	_, err := rc.RepositoryStatement.ClientConfig.DeleteByID.ExecContext(rc.ctx, id)
	return err
}

func (rc *RepositoryContext) IsClientConfigExists(appId int64, key string) (bool, error) {
	var result bool
	err := rc.RepositoryStatement.ClientConfig.IsExistsByKey.GetContext(rc.ctx, &result, appId, key)
	return result, err
}
