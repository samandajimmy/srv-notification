package notification

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nbs-go/nlogger"
	"reflect"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/convert"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
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
		Metadata:      []byte("{}"),
		ItemMetadata:  model.NewItemMetadata(convert.ModifierDTOToModel(payload.Subject.ModifiedBy)),
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

	return composeClientConfigResponse(&clientConfig)
}

func (s *ServiceContext) ListClientConfig(options dto.ClientConfigFindOptions) (*dto.ClientConfigListResponse, error) {
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
	queryResult, err := s.repo.FindClientConfig(&options.FindOptions)
	if err != nil {
		log.Error("failed to find data client config.", nlogger.Error(err))
		return nil, ncore.TraceError(err)
	}

	// Compose response
	rowsResp := make([]dto.ClientConfigItemResponse, len(queryResult.Rows))
	for idx, row := range queryResult.Rows {
		var rowItem = dto.ClientConfigItemResponse{
			XID:                  row.XID,
			Key:                  row.Key,
			Value:                row.Value,
			ApplicationId:        row.ApplicationId,
			ItemMetadataResponse: convert.ItemMetadataModelToResponse(row.ItemMetadata),
		}
		rowsResp[idx] = rowItem
	}

	return &dto.ClientConfigListResponse{
		ClientConfig: rowsResp,
		Metadata: dto.ListMetadata{
			Count:       queryResult.Count,
			FindOptions: options.FindOptions,
		},
	}, nil
}

func (s *ServiceContext) GetClientConfig(payload dto.ClientConfigRequest) (*dto.ClientConfigItemResponse, error) {
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
	clientConfig, err := s.repo.FindClientConfigByXID(payload.XID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("error when get data client config.", nlogger.Error(err))
			return nil, s.responses.GetError("E_RES_1")
		}
		log.Error("error when get data client config.", nlogger.Error(err))
		return nil, err
	}

	// Validate version
	if clientConfig.Version != payload.Version {
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
		case "key":
			// If title is empty, or value is still the same, then skip
			if d.Key == "" || d.Key == clientConfig.Key {
				changelog[k] = false
				continue
			}
			// Set updated value
			clientConfig.Key = d.Key
			changesCount += 1
		case "applicationId":
			// If title is empty, or value is still the same, then skip
			if d.ApplicationId == 0 || d.ApplicationId == clientConfig.ApplicationId {
				changelog[k] = false
				continue
			}
			// Set updated value
			clientConfig.ApplicationId = d.ApplicationId
			changesCount += 1
		case "value":
			var clientConfigValue map[string]string
			err := json.Unmarshal(clientConfig.Value, &clientConfigValue)
			if err != nil {
				return nil, ncore.TraceError(err)
			}
			// comparing
			payloadValue := d.Value
			if payloadValue == nil || reflect.DeepEqual(payloadValue, clientConfigValue) {
				changelog[k] = false
				continue
			}
			// convert to byte
			value, err := json.Marshal(d.Value)
			if err != nil {
				return nil, ncore.TraceError(err)
			}
			// Set updated value
			clientConfig.Value = value
			changesCount += 1
		}
	}

	// If changes count more than 0, then persist update
	if changesCount > 0 {
		// Update metadata
		modifiedBy := convert.ModifierDTOToModel(payload.Subject.ModifiedBy)
		clientConfig.UpdatedAt = time.Now()
		clientConfig.ModifiedBy = &modifiedBy
		clientConfig.Version += 1
		// Update client config
		err = s.repo.UpdateClientConfig(clientConfig)
		if err != nil {
			panic(fmt.Errorf("failed to delete client config. Error = %w", err))
		}
	}

	return composeDetailClientConfigResponse(clientConfig)
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
	err = s.repo.DeleteClientConfigById(res.ID)
	if err != nil {
		panic(fmt.Errorf("failed to delete client config. Error = %w", err))
	}

	return nil
}

func composeClientConfigResponse(row *model.ClientConfig) (*dto.ClientConfigItemResponse, error) {
	return &dto.ClientConfigItemResponse{
		ApplicationId:        row.ApplicationId,
		Key:                  row.Key,
		Value:                row.Value,
		XID:                  row.XID,
		ItemMetadataResponse: convert.ItemMetadataModelToResponse(row.ItemMetadata),
	}, nil
}

func composeDetailClientConfigResponse(row *model.ClientConfig) (*dto.ClientConfigItemResponse, error) {
	return &dto.ClientConfigItemResponse{
		ApplicationId:        row.ApplicationId,
		Key:                  row.Key,
		Value:                row.Value,
		XID:                  row.XID,
		ItemMetadataResponse: convert.ItemMetadataModelToResponse(row.ItemMetadata),
	}, nil
}
