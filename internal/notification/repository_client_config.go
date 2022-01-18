package notification

import (
	"fmt"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"strings"
)

func (rc *RepositoryContext) HasInitialized() bool {
	return true
}

func (rc *RepositoryContext) FindByKey(key string, appId int) (*model.ClientConfig, error) {
	var row model.ClientConfig
	err := rc.RepositoryStatement.ClientConfig.FindByKey.Get(&row, key, appId)
	return &row, err
}

func (rc *RepositoryContext) FindClientConfigByXID(xid string) (*model.ClientConfig, error) {
	var row model.ClientConfig
	err := rc.RepositoryStatement.ClientConfig.FindByXID.Get(&row, xid)
	return &row, err
}

func (rc *RepositoryContext) InsertClientConfig(row model.ClientConfig) error {
	_, err := rc.RepositoryStatement.ClientConfig.Insert.Exec(row)
	return err
}

func (rc *RepositoryContext) FindClientConfig(params *dto.FindOptions) (*model.ClientConfigSearchResult, error) {
	// Prepare where
	var args []interface{}
	var whereQuery []string

	where := ""
	if len(whereQuery) > 0 {
		where = "WHERE " + strings.Join(whereQuery, " AND ")
	}

	// Prepare query
	columns := `"createdAt", "updatedAt", "metadata", "modifiedBy", "version", "key", "value", "applicationId", "xid"`
	from := `ClientConfig`
	orderBy := rc.GetOrderByQuery(params.SortBy, params.SortDirection)
	q := fmt.Sprintf(`SELECT %s FROM "%s" %s ORDER BY %s LIMIT %d OFFSET %d`,
		columns, from, where, orderBy, params.Limit, params.Skip)

	// Execute query
	q = rc.conn.Rebind(q)
	var rows []model.ClientConfig
	err := rc.conn.SelectContext(rc.ctx, &rows, q, args...)
	if err != nil {
		log.Error("Error when execute query.", nlogger.Error(err))
		return nil, ncore.TraceError(err)
	}
	// Count all
	countQuery := fmt.Sprintf(`SELECT COUNT(id) FROM "%s" %s`, from, where)
	countQuery = rc.conn.Rebind(countQuery)
	var count int64
	err = rc.conn.GetContext(rc.ctx, &count, countQuery, args...)
	if err != nil {
		log.Error("Error when execute query count.", nlogger.Error(err))
		return nil, ncore.TraceError(err)
	}

	// Prepare result
	result := model.ClientConfigSearchResult{
		Rows:  rows,
		Count: count,
	}
	return &result, err
}

func (rc *RepositoryContext) UpdateClientConfig(row *model.ClientConfig) error {
	result, err := rc.RepositoryStatement.ClientConfig.UpdateByID.Exec(row)
	if err != nil {
		return err
	}
	return nsql.IsUpdated(result)
}

func (rc *RepositoryContext) DeleteClientConfigById(id int64) error {
	_, err := rc.RepositoryStatement.ClientConfig.DeleteByID.ExecContext(rc.ctx, id)
	return err
}
