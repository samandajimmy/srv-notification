package notification

import "repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"

func (rc *RepositoryContext) InsertApplication(row model.Application) error {
	_, err := rc.RepositoryStatement.Application.Insert.ExecContext(rc.ctx, row)
	return err
}
