package notification

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/convert"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
)

func (s *ServiceContext) CreateClientConfig(payload dto.ClientConfig) (*dto.ClientConfigItemResponse, error) {
	// Initialize data to insert
	value, err := json.Marshal(payload.Value)
	if err != nil {
		return nil, err
	}

	// TODO Check ApplicationId is exists

	clientConfig := model.ClientConfig{
		XID:           payload.XID,
		Key:           payload.Key,
		Value:         value,
		ApplicationId: payload.ApplicationId,
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
		"name",
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
	queryResult, err := s.repo.Find(&options.FindOptions)
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

func (s *ServiceContext) GetClientConfig(payload dto.ClientConfig) (*dto.ClientConfigItemResponse, error) {
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
