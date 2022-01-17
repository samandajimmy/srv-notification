package notification

import (
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/convert"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
)

func (s *ServiceContext) CreateApplication(payload dto.Application) (*dto.ApplicationResponse, error) {
	// Initialize data to insert
	xid, err := gonanoid.Generate(constant.AlphaNumUpperCharSet, 8)
	if err != nil {
		panic(fmt.Errorf("failed to generate xid. Error = %w", err))
	}

	apl := model.Application{
		XID:          xid,
		Name:         payload.Name,
		Metadata:     []byte("{}"),
		ItemMetadata: model.NewItemMetadata(convert.ModifierDTOToModel(payload.Subject.ModifiedBy)),
	}

	// Persist application
	err = s.repo.InsertApplication(apl)
	if err != nil {
		log.Errorf("unable to insert application: %v", err)

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

	return composeDetailApplicationResponse(&apl)

}

func (s *ServiceContext) GetApplication(payload dto.GetApplication) (*dto.ApplicationResponse, error) {
	// Get application by xid
	res, err := s.repo.FindApplicationByXID(payload.XID)
	if err != nil {
		log.Errorf("error when get data application. err: %v", err)
		return nil, err
	}

	return composeDetailApplicationResponse(res)
}

func composeDetailApplicationResponse(row *model.Application) (*dto.ApplicationResponse, error) {
	return &dto.ApplicationResponse{
		Name:                 row.Name,
		XID:                  row.XID,
		ItemMetadataResponse: convert.ItemMetadataModelToResponse(row.ItemMetadata),
	}, nil
}
