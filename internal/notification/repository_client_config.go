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

func (rc *RepositoryContext) FindByKey(key string, appId int64) (*model.ClientConfig, error) {
	var row model.ClientConfig
	err := rc.RepositoryStatement.ClientConfig.FindByKey.Get(&row, key, appId)
	return &row, err
}

func (rc *RepositoryContext) FindClientConfigByXID(xid string) (*model.ClientConfigVO, error) {
	var row model.ClientConfigVO
	err := rc.RepositoryStatement.ClientConfig.FindJoinApplicationByXID.Get(&row, xid)
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

	if applicationXid, ok := params.Filters["applicationXid"]; ok {
		whereQuery = append(whereQuery, `"Application"."xid" = ?`)
		args = append(args, applicationXid)
	}

	where := ""
	if len(whereQuery) > 0 {
		where = "WHERE " + strings.Join(whereQuery, " AND ")
	}

	// Prepare query
	columns := `"ClientConfig"."createdAt", "ClientConfig"."updatedAt", "ClientConfig"."metadata", "ClientConfig"."modifiedBy", "ClientConfig"."version", "ClientConfig"."key", "ClientConfig"."value", "ClientConfig"."xid", "Application"."xid" AS "applicationXid"`
	from := `ClientConfig`
	// Join Table
	joinApplication := `LEFT JOIN "Application" ON "Application"."id" = "ClientConfig"."applicationId"`
	// Order By
	orderBy := rc.GetOrderByQuery(params.SortBy, params.SortDirection)
	// query string
	q := fmt.Sprintf(
		`SELECT %s FROM "%s" %s %s ORDER BY %s LIMIT %d OFFSET %d`,
		columns, from, joinApplication, where, orderBy, params.Limit, params.Skip)
	//log2.Fatalf("q: %s", q)
	// count query string
	countQuery := fmt.Sprintf(`SELECT COUNT("ClientConfig"."id") FROM "%s" %s %s`, from, joinApplication, where)
	// Execute query
	q = rc.conn.Rebind(q)
	var rows []model.ClientConfigVO
	err := rc.conn.SelectContext(rc.ctx, &rows, q, args...)
	if err != nil {
		log.Error("Error when execute query.", nlogger.Error(err))
		return nil, ncore.TraceError(err)
	}
	// Count all
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
