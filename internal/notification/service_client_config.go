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
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/convert"
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

	clientConfigFindByXID, err := s.repo.FindClientConfigByXID(xid)
	if err != nil {
		return nil, err
	}

	return composeDetailClientConfigResponse(clientConfigFindByXID)
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
			ApplicationXid:       row.ApplicationXid,
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

	modelClientConfig := model.ClientConfig{
		ID:           clientConfig.ID,
		XID:          clientConfig.XID,
		Key:          clientConfig.Key,
		Value:        clientConfig.Value,
		Metadata:     clientConfig.Metadata,
		ItemMetadata: clientConfig.ItemMetadata,
	}

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
			modelClientConfig.Key = d.Key
			changesCount += 1
		case "applicationXid":
			// If title is empty, or value is still the same, then skip
			if d.ApplicationXid == "" || d.ApplicationXid == clientConfig.ApplicationXid {
				changelog[k] = false
				continue
			}
			// check application by applicationXid
			application, err := s.repo.FindApplicationByXID(d.ApplicationXid)
			if err != nil {
				log.Error("application not found when Update ClientConfig.", nlogger.Error(err))
				return nil, s.responses.GetError("E_RES_1")
			}
			// Set updated value
			modelClientConfig.ApplicationId = application.ID
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
			modelClientConfig.Value = value
			changesCount += 1
		}
	}

	// If changes count more than 0, then persist update
	if changesCount > 0 {
		// Update metadata
		modifiedBy := convert.ModifierDTOToModel(payload.Subject.ModifiedBy)
		modelClientConfig.UpdatedAt = time.Now()
		modelClientConfig.ModifiedBy = &modifiedBy
		modelClientConfig.Version += 1
		// Update client config
		err = s.repo.UpdateClientConfig(&modelClientConfig)
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

func composeDetailClientConfigResponse(row *model.ClientConfigVO) (*dto.ClientConfigItemResponse, error) {
	return &dto.ClientConfigItemResponse{
		ApplicationXid:       row.ApplicationXid,
		Key:                  row.Key,
		Value:                row.Value,
		XID:                  row.XID,
		ItemMetadataResponse: convert.ItemMetadataModelToResponse(row.ItemMetadata),
	}, nil
}
