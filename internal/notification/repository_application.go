package notification

import "repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"

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
