package notification

import (
	"fmt"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"strings"
	"time"
)

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

func (rc *RepositoryContext) FindApplication(params *dto.ApplicationFindOptions) (*model.ApplicationFindResult, error) {
	// Prepare where
	args := []interface{}{constant.DefaultConfig}
	whereQuery := []string{"xid != ?"}

	if xid, ok := params.Filters["xid"]; ok {
		whereQuery = append(whereQuery, `xid = ?`)
		args = append(args, xid)
	}

	if name, ok := params.Filters["name"]; ok {
		whereQuery = append(whereQuery, `name LIKE ?`)
		args = append(args, fmt.Sprintf(`%%%s%%`, name))
	}

	if filterCreatedFrom, ok := params.Filters["createdFrom"]; ok {
		whereQuery = append(whereQuery, "createdAt >= ?")
		// Convert unix timestamp to date string
		t := time.Unix(filterCreatedFrom.(int64), 0).UTC()
		args = append(args, t.Format("2006-01-02 15:04:05"))
	}

	if filterCreatedUntil, ok := params.Filters["createdUntil"]; ok {
		whereQuery = append(whereQuery, "createdAt <= ?")
		// Convert unix timestamp to date string
		t := time.Unix(filterCreatedUntil.(int64), 0).UTC()
		args = append(args, t.Format("2006-01-02 15:04:05"))
	}

	where := ""
	if len(whereQuery) > 0 {
		where = "WHERE " + strings.Join(whereQuery, " AND ")
	}

	// Prepare query
	columns := `"id","metadata","createdAt","updatedAt","modifiedBy","version","xid","name","webhookUrl"`
	from := `Application`
	queryList := fmt.Sprintf(`SELECT %s FROM "%s" %s ORDER BY %s LIMIT %d OFFSET %d`,
		columns,
		from,
		where,
		rc.GetOrderByQuery(params.SortBy, params.SortDirection),
		params.Limit,
		params.Skip)

	// Execute query
	queryList = rc.conn.Rebind(queryList)

	var rows []model.Application
	err := rc.conn.SelectContext(rc.ctx, &rows, queryList, args...)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Count all
	queryCount := fmt.Sprintf(`SELECT COUNT(id) FROM "%s" %s`, from, where)
	queryCount = rc.conn.Rebind(queryCount)
	var count int64
	err = rc.conn.GetContext(rc.ctx, &count, queryCount, args...)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Prepare result
	result := model.ApplicationFindResult{
		Rows:  rows,
		Count: count,
	}
	return &result, err
}

func (rc *RepositoryContext) UpdateApplication(row *model.Application) error {
	result, err := rc.Application.Update.ExecContext(rc.ctx, row)
	if err != nil {
		return err
	}
	return nsql.IsUpdated(result)
}
