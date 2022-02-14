package notification

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nbs-go/nlogger"
	"reflect"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	dto "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	model "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
	"time"
)

func (s *ServiceContext) CreateClientConfig(payload dto.ClientConfigRequest) (*dto.ClientConfigItemResponse, error) {
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

	// check application exist
	application, err := s.repo.FindApplicationByXID(payload.ApplicationXid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error when get data application.", nlogger.Error(err))
			return nil, s.responses.GetError("E_RES_1")
		}
		log.Error("error when get data application.", nlogger.Error(err))
		return nil, err
	}

	clientConfig := model.ClientConfig{
		XID:           xid,
		Key:           payload.Key,
		Value:         value,
		ApplicationId: application.ID,
		BaseField:     model.NewBaseField(model.ToModifier(payload.Subject.ModifiedBy)),
	}

	// Persist client config
	err = s.repo.InsertClientConfig(clientConfig)
	if err != nil {
		log.Error("unable to insert clientConfig.", nlogger.Error(err))
		// Handle pq.Error
		errCode, _ := nsql.GetPostgresError(err)
		switch errCode {
		case nsql.UniqueError:
			return nil, s.responses.GetError("E_UAL_1").Wrap(err)
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
	// Handle sort request
	rulesSortBy := []string{
		"createdAt",
		"updatedAt",
		"key",
	}

	// Get orderBy
	sortBy, sortDirection := s.GetOrderBy(
		nval.ParseStringFallback(options.SortBy, `createdAt`),
		nval.ParseStringFallback(options.SortDirection, `desc`),
		rulesSortBy,
	)

	// Set sort by and direction
	options.SortBy = sortBy
	options.SortDirection = sortDirection
	// Query
	queryResult, err := s.repo.FindClientConfig(options)
	if err != nil {
		log.Error("failed to find data client config.", nlogger.Error(err))
		return nil, ncore.TraceError(err)
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

func (s *ServiceContext) GetDetailClientConfig(payload dto.ClientConfigRequest) (*dto.ClientConfigItemResponse, error) {
	// Get client config by xid
	res, err := s.repo.FindClientConfigByXID(payload.XID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, s.responses.GetError("E_RES_1")
		}
		log.Error("failed to client config", nlogger.Error(err))
		return nil, ncore.TraceError(err)
	}

	return composeDetailClientConfigResponse(res)
}

func (s *ServiceContext) UpdateClientConfig(payload dto.ClientConfigUpdateOptions) (*dto.ClientConfigItemResponse, error) {
	// Get client config by xid
	row, err := s.repo.FindClientConfigByXID(payload.XID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, s.responses.GetError("E_RES_1")
		}
		log.Error("error when get data client config.", nlogger.Error(err))
		return nil, err
	}

	// Get model
	m := row.ClientConfig

	// Validate version
	if m.Version != payload.Version {
		return nil, s.responses.GetError("E_RES_2").Wrap(err)
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
				return nil, ncore.TraceError(err)
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
				return nil, ncore.TraceError(err)
			}
			// Set updated value
			m.Value = value
			changesCount += 1
		}
	}

	// If changes count more than 0, then persist update
	if changesCount > 0 {
		// Update metadata
		modifiedBy := model.ToModifier(payload.Subject.ModifiedBy)
		m.UpdatedAt = time.Now()
		m.ModifiedBy = modifiedBy
		m.Version += 1

		// Update client config
		err = s.repo.UpdateClientConfig(m, payload.Version)
		if err != nil {
			return nil, ncore.TraceError(err)
		}
	}

	return composeDetailClientConfigResponse(row)
}

func (s *ServiceContext) DeleteClientConfig(payload dto.GetClientConfig) error {
	// Get application by xid
	res, err := s.repo.FindClientConfigByXID(payload.XID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error when get data client config. err: %v", err)
			return s.responses.GetError("E_RES_1")
		}
		log.Error("error when get data client config. err: %v", err)
		return err
	}

	// Delete client config
	err = s.repo.DeleteClientConfigById(res.ClientConfig.ID)
	if err != nil {
		panic(fmt.Errorf("failed to delete client config. Error = %w", err))
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
