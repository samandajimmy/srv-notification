package notification

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nbs-go/errx"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"reflect"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	svcError "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/error"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"time"
)

func (s *ServiceContext) CreateClientConfig(payload *dto.ClientConfigRequest) (*dto.ClientConfigItemResponse, error) {
	// Get application
	application, err := s.repo.FindApplicationByXID(payload.ApplicationXid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error when get data application.", logOption.Error(err))
			return nil, svcError.ResourceNotFound.Trace(errx.Source(err))
		}
		log.Error("error when get data application.", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Check if config for certain key has been added
	isExists, err := s.repo.IsClientConfigExists(application.ID, payload.Key)
	if err != nil {
		log.Error("error while query IsClientConfigExists", logOption.Error(err))
		return nil, err
	}

	if isExists {
		log.Errorf("client config for Key %s is already exists", payload.Key)
		return nil, svcError.DuplicatedResource.Trace()
	}

	// Initialize data to insert
	xid, err := gonanoid.Generate(constant.AlphaNumUpperCharSet, 8)
	if err != nil {
		panic(fmt.Errorf("failed to generate xid. Error = %w", err))
	}

	// Initialize data to insert
	value, err := json.Marshal(payload.Value)
	if err != nil {
		return nil, err
	}

	clientConfig := model.ClientConfig{
		XID:           xid,
		Key:           payload.Key,
		Value:         value,
		ApplicationId: application.ID,
		BaseField:     model.NewBaseField(model.ToModifier(payload.Subject.ModifiedBy())),
	}

	// Persist client config
	err = s.repo.InsertClientConfig(clientConfig)
	if err != nil {
		log.Error("unable to insert clientConfig.", logOption.Error(err))
		// Handle pq.Error
		errCode, _ := nsql.GetPostgresError(err)
		switch errCode {
		case nsql.UniqueError:
			return nil, svcError.DuplicatedResource.Trace()
		default:
			return nil, err
		}
	}

	return composeDetailClientConfigResponse(&model.ClientConfigDetailed{
		ClientConfig: &clientConfig,
		Application:  application,
	})
}

func (s *ServiceContext) ListClientConfig(options *dto.ListPayload) (*dto.ClientConfigListResponse, error) {
	// Query
	queryResult, err := s.repo.FindClientConfig(options)
	if err != nil {
		log.Error("failed to find data client config.", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	// Compose response
	rowsResp := make([]dto.ClientConfigItemResponse, len(queryResult.Rows))
	for idx, row := range queryResult.Rows {
		var rowItem = dto.ClientConfigItemResponse{
			XID:            row.ClientConfig.XID,
			Key:            row.ClientConfig.Key,
			Value:          row.ClientConfig.Value,
			ApplicationXid: row.Application.XID,
			BaseField:      model.ToBaseFieldDTO(row.ClientConfig.BaseField),
		}
		rowsResp[idx] = rowItem
	}

	return &dto.ClientConfigListResponse{
		ClientConfig: rowsResp,
		Metadata:     dto.ToListMetadata(options, queryResult.Count),
	}, nil
}

func (s *ServiceContext) GetDetailClientConfig(payload *dto.ClientConfigRequest) (*dto.ClientConfigItemResponse, error) {
	// Get client config by xid
	res, err := s.repo.FindClientConfigByXID(payload.XID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, svcError.ResourceNotFound.Trace(errx.Source(err))
		}
		log.Error("failed to client config", logOption.Error(err))
		return nil, errx.Trace(err)
	}

	return composeDetailClientConfigResponse(res)
}

func (s *ServiceContext) UpdateClientConfig(payload *dto.ClientConfigUpdateOptions) (*dto.ClientConfigItemResponse, error) {
	// Get client config by xid
	row, err := s.repo.FindClientConfigByXID(payload.XID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, svcError.ResourceNotFound.Trace(errx.Source(err))
		}
		log.Error("error when get data client config.", logOption.Error(err))
		return nil, err
	}

	// Get model
	m := row.ClientConfig

	// Validate version
	if m.Version != payload.Version {
		return nil, svcError.StaleResource.Trace()
	}

	// Copy values from payload
	d := payload.Data
	changelog := payload.Changelog
	changesCount := 0
	for k, changed := range changelog {
		// If not changed, then continue
		if !changed {
			continue
		}

		switch k {
		case "value":
			var updatedValue map[string]string
			jErr := json.Unmarshal(m.Value, &updatedValue)
			if jErr != nil {
				return nil, errx.Trace(err)
			}

			// comparing
			payloadValue := d.Value
			if payloadValue == nil || reflect.DeepEqual(payloadValue, updatedValue) {
				changelog[k] = false
				continue
			}

			// convert to byte
			value, jErr := json.Marshal(d.Value)
			if jErr != nil {
				return nil, errx.Trace(err)
			}
			// Set updated value
			m.Value = value
			changesCount += 1
		}
	}

	// If changes count more than 0, then persist update
	if changesCount > 0 {
		// Update metadata
		modifiedBy := model.ToModifier(payload.Subject.ModifiedBy())
		m.UpdatedAt = time.Now()
		m.ModifiedBy = modifiedBy
		m.Version += 1

		// Update client config
		err = s.repo.UpdateClientConfig(m, payload.Version)
		if err != nil {
			return nil, errx.Trace(err)
		}
	}

	return composeDetailClientConfigResponse(row)
}

func (s *ServiceContext) DeleteClientConfig(payload *dto.GetClientConfig) error {
	// Get application by xid
	res, err := s.repo.FindClientConfigByXID(payload.XID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error when get data client config", logOption.Error(err))
			return svcError.ResourceNotFound.Trace(errx.Source(err))
		}
		log.Error("error when get data client config", logOption.Error(err))
		return err
	}

	// Delete client config
	err = s.repo.DeleteClientConfigById(res.ClientConfig.ID)
	if err != nil {
		panic(fmt.Errorf("failed to delete client config. Error = %w", err))
	}

	return nil
}

func (s *ServiceContext) loadClientConfig(appId int64, key string, dest interface{}) error {
	// Get client config from database
	clientConfig, err := s.repo.FindByKey(key, appId)
	if err != nil {
		if err != sql.ErrNoRows {
			s.log.Error("failed to get ClientConfig from db", logOption.Error(err),
				logOption.AddMetadata("key", key), logOption.AddMetadata("applicationId", appId))
			return errx.Trace(err)
		}

		// Get default config
		clientConfig, err = s.repo.FindDefaultClientConfigByKey(key)
		if err != nil {
			if err == sql.ErrNoRows {
				return errx.Trace(fmt.Errorf("default configuration not set for key %s", key))
			}
			s.log.Error("failed to get default ClientConfig from db", logOption.Error(err),
				logOption.AddMetadata("key", key), logOption.AddMetadata("applicationId", appId))
			return errx.Trace(err)
		}
	}

	err = json.Unmarshal(clientConfig.Value, dest)
	if err != nil {
		s.log.Error("Error when unmarshaling ClientConfig value", logOption.Error(err),
			logOption.AddMetadata("key", key), logOption.AddMetadata("applicationId", appId))
		return errx.Trace(err)
	}

	return nil
}

func composeDetailClientConfigResponse(row *model.ClientConfigDetailed) (*dto.ClientConfigItemResponse, error) {
	return &dto.ClientConfigItemResponse{
		ApplicationXid: row.Application.XID,
		Key:            row.ClientConfig.Key,
		Value:          row.ClientConfig.Value,
		XID:            row.ClientConfig.XID,
		BaseField:      model.ToBaseFieldDTO(row.ClientConfig.BaseField),
	}, nil
}
