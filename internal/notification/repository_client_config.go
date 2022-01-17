package notification

import (
	"fmt"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"strings"
)

func (rc *RepositoryContext) HasInitialized() bool {
	return true
}

func (rc *RepositoryContext) FindByKey(key string, appId int) (*model.ClientConfig, error) {
	var row model.ClientConfig
	err := rc.ClientConfig.FindByKey.Get(&row, key, appId)
	return &row, err
}

func (rc *RepositoryContext) InsertClientConfig(row model.ClientConfig) error {
	_, err := rc.ClientConfig.Insert.Exec(row)
	return err
}

func (rc *RepositoryContext) Find(params *dto.FindOptions) (*model.ClientConfigSearchResult, error) {
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
	orderBy := rc.getOrderByQuery(params)
	q := fmt.Sprintf(`SELECT %s FROM "%s" %s ORDER BY %s LIMIT %d OFFSET %d`,
		columns, from, where, orderBy, params.Limit, params.Skip,
	)

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
