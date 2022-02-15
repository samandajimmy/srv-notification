package notification

import (
	"database/sql"
	"errors"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	dto "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	model "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"time"
)

func (s *ServiceContext) AuthApplication(username string, password string) (*dto.AuthApplicationResponse, error) {
	applicationXid := username
	apiKey := password

	if applicationXid == constant.DefaultConfig {
		log.Warn("cannot use default configuration as app")
		return nil, nhttp.ForbiddenError
	}

	application, err := s.repo.FindApplicationByXID(applicationXid)
	if err != nil {
		log.Error("application not found", nlogger.Error(err))
		return nil, nhttp.UnauthorizedError
	}

	if application.ApiKey != apiKey {
		log.Error("Incorrect apiKey", nlogger.Error(err))
		return nil, nhttp.UnauthorizedError
	}

	webhookURL := ""
	if application.WebhookURL.Valid {
		webhookURL = application.WebhookURL.String
	}

	result := &dto.AuthApplicationResponse{
		ID:         application.ID,
		XID:        application.XID,
		Name:       application.Name,
		ApiKey:     application.ApiKey,
		WebhookURL: webhookURL,
	}

	return result, err
}

func (s *ServiceContext) CreateApplication(payload *dto.Application) (*dto.ApplicationResponse, error) {
	// Initialize data to insert
	xid, err := gonanoid.Generate(constant.AlphaNumUpperCharSet, 8)
	if err != nil {
		panic(fmt.Errorf("failed to generate xid. Error = %w", err))
	}

	// Initialize data to insert
	apl := model.Application{
		XID:    xid,
		ApiKey: newApplicationApiKey(),
		Name:   payload.Name,
		WebhookURL: sql.NullString{
			String: payload.WebhookURL,
			Valid:  payload.WebhookURL != "",
		},
		BaseField: model.NewBaseField(model.ToModifier(payload.Subject.ModifiedBy())),
	}

	// Persist application
	err = s.repo.InsertApplication(apl)
	if err != nil {
		log.Errorf("unable to insert application. err: %v", err)
		// Handle pq.Error
		errCode, _ := nsql.GetPostgresError(err)

		switch errCode {
		case nsql.UniqueError:
			return nil, s.responses.GetError("E_UAL_1").Wrap(err)
		default:
			return nil, err
		}
	}

	return composeDetailApplicationResponse(&apl)

}

func (s *ServiceContext) GetDetailApplication(payload *dto.GetApplication) (*dto.ApplicationResponse, error) {
	if payload.XID == constant.DefaultConfig {
		log.Warn("did not allowed retrieve default config as app")
		return nil, nhttp.ForbiddenError
	}

	// Get application by xid
	res, err := s.repo.FindApplicationByXID(payload.XID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error when get data application. err: %v", err)
			return nil, s.responses.GetError("E_RES_1")
		}
		log.Error("error when get data application. err: %v", err)
		return nil, err
	}

	return composeDetailApplicationResponse(res)
}

func (s *ServiceContext) DeleteApplication(payload *dto.GetApplication) error {
	if payload.XID == constant.DefaultConfig {
		log.Warn("cannot delete default config app")
		return nhttp.ForbiddenError
	}

	// Get application by xid
	res, err := s.repo.FindApplicationByXID(payload.XID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error when get data application. err: %v", err)
			return s.responses.GetError("E_RES_1")
		}
		log.Error("error when get data application. err: %v", err)
		return err
	}

	// Delete application
	err = s.repo.DeleteApplicationById(res.ID)
	if err != nil {
		panic(fmt.Errorf("failed to delete application. Error = %w", err))
	}

	return nil
}

func (s *ServiceContext) ListApplication(options *dto.ListPayload) (*dto.ListApplicationResponse, error) {
	// Get list application
	result, err := s.repo.FindApplication(options)
	if err != nil {
		log.Error("failed to find data application", nlogger.Error(err))
		return nil, ncore.TraceError(err)
	}

	rows := make([]*dto.ApplicationItem, len(result.Rows))
	for idx, row := range result.Rows {
		rows[idx], _ = composeListApplicationResponse(&row)
	}

	return &dto.ListApplicationResponse{
		Items:    rows,
		Metadata: dto.ToListMetadata(options, result.Count),
	}, err
}

func (s *ServiceContext) UpdateApplication(payload *dto.ApplicationUpdateOptions) (*dto.ApplicationResponse, error) {
	if payload.XID == constant.DefaultConfig {
		log.Warn("cannot update default config app")
		return nil, nhttp.ForbiddenError
	}

	// Get application by xid
	app, err := s.repo.FindApplicationByXID(payload.XID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error when get data application. err: %v", err)
			return nil, s.responses.GetError("E_RES_1")
		}
		log.Error("error when get data application. err: %v", err)
		return nil, err
	}

	// Validate version
	if app.Version != payload.Version {
		log.Errorf("invalid version to change. expected: %v actual: %v",
			app.Version,
			payload.Version,
		)
		return nil, s.responses.GetError("E_RES_2").Wrap(err)
	}

	// Copy values from payload to job
	d := payload.Data
	changelog := payload.Changelog
	changesCount := 0

	for k, changed := range changelog {
		// If not changed, then continue
		if !changed {
			continue
		}
		switch k {
		case "name":
			// If title is empty, or value is still the same, then skip
			if d.Name == "" || d.Name == app.Name {
				changelog[k] = false
				continue
			}

			// Set updated value
			app.Name = d.Name
			changesCount += 1
		case "apiKey":
			// Generate new api key
			app.ApiKey = newApplicationApiKey()
			changesCount += 1
		case "webhookUrl":
			// If title is empty, or value is still the same, then skip
			if d.WebhookURL == "" || d.WebhookURL == app.WebhookURL.String {
				changelog[k] = false
				continue
			}

			// Set updated value
			app.WebhookURL.Valid = true
			app.WebhookURL.String = d.WebhookURL
			changesCount += 1
		}

	}

	// If changes count more than 0, then persist update
	if changesCount > 0 {
		// Update metadata
		modifiedBy := model.ToModifier(payload.Subject.ModifiedBy())
		app.UpdatedAt = time.Now()
		app.ModifiedBy = modifiedBy
		app.Version += 1

		// Persist
		err = s.repo.UpdateApplication(app, payload.Version)
		if err != nil {
			if errors.Is(err, nsql.RowNotUpdatedError) {
				err = s.responses.GetError("E_RES_3").Wrap(err)
			} else {
				log.Errorf("failed to persist application update. err: %v", err)
			}
			return nil, err
		}
	}

	return composeDetailApplicationResponse(app)
}

func composeListApplicationResponse(row *model.Application) (*dto.ApplicationItem, error) {
	return &dto.ApplicationItem{
		Name:      row.Name,
		XID:       row.XID,
		BaseField: model.ToBaseFieldDTO(row.BaseField),
	}, nil
}

func composeDetailApplicationResponse(row *model.Application) (*dto.ApplicationResponse, error) {
	webhookUrl := ""
	if row.WebhookURL.Valid {
		webhookUrl = row.WebhookURL.String
	}

	return &dto.ApplicationResponse{
		Name:       row.Name,
		XID:        row.XID,
		ApiKey:     row.ApiKey,
		WebhookURL: webhookUrl,
		BaseField:  model.ToBaseFieldDTO(row.BaseField),
	}, nil
}

func newApplicationApiKey() string {
	apiKey, err := gonanoid.Generate(constant.AlphaNumUpperCharSet, 32)
	if err != nil {
		panic(fmt.Errorf("failed to generate apiKey. Error = %w", err))
	}
	return apiKey
}
