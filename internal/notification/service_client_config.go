package notification

import (
	"encoding/json"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/convert"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

func (s *ServiceContext) CreateClientConfig(payload dto.ClientConfig) (*dto.ClientConfigItemResponse, error) {
	// Initialize data to insert
	xid, err := gonanoid.Generate(constant.AlphaNumUpperCharSet, 8)
	if err != nil {
		panic(fmt.Errorf("failed to generate xid. Error = %w", err))
	}

	value, err := json.Marshal(payload.Value)
	if err != nil {
		return nil, err
	}

	clientConfig := model.ClientConfig{
		XID:           xid,
		Key:           payload.Key,
		Value:         value,
		ApplicationId: payload.ApplicationId,
		Metadata:      []byte("{}"),
		ItemMetadata:  model.NewItemMetadata(convert.ModifierDTOToModel(payload.Subject.ModifiedBy)),
	}

	// Persist application
	err = s.repo.InsertClientConfig(clientConfig)
	if err != nil {
		log.Errorf("unable to insert client config: %v", err)

		// Handle pq.Error
		errCode, _ := nsql.GetPostgresError(err)
		switch errCode {
		case nsql.UniqueError:
			// TODO: get response from ncore.Response s.responses.GetError("E_UAL_1").Wrap(err)
			return nil, err
		default:
			return nil, err
		}
	}

	return composeClientConfigResponse(&clientConfig)

}

func (s *ServiceContext) ListClientConfig(params dto.ClientConfigFindOptions) (*dto.ClientConfigListResponse, error) {
	// Query
	queryResult, err := s.repo.Find(&params.FindOptions)
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
			FindOptions: params.FindOptions,
		},
	}, nil
}

func composeClientConfigResponse(row *model.ClientConfig) (*dto.ClientConfigItemResponse, error) {
	return &dto.ClientConfigItemResponse{
		Key:                  row.Key,
		Value:                row.Value,
		XID:                  row.XID,
		ItemMetadataResponse: convert.ItemMetadataModelToResponse(row.ItemMetadata),
	}, nil
}
