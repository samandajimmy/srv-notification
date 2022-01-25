package notification

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nsql"
	"time"
)

func (s *ServiceContext) CreateNotification(payload dto.SendNotificationOptionsRequest) error {
	// Generate temporary id
	id, err := uuid.NewUUID()
	if err != nil {
		log.Error("error when get uuid", nlogger.Error(err))
		return err
	}

	payloadFCM := payload.Options.FCM

	now := time.Now()

	// init options
	options := map[string]interface{}{}
	// Initialize data to insert
	notification := model.Notification{
		ID:            id,
		ApplicationId: payload.Auth.ID,
		UserRefId:     payload.UserId,
		IsRead:        false,
		ReadAt:        sql.NullTime{},
		Metadata:      []byte("{}"),
		ItemMetadata: model.ItemMetadata{
			CreatedAt:  now,
			UpdatedAt:  now,
			ModifiedBy: &model.Modifier{ID: "", Role: "", FullName: ""}, // TODO Get Subject
			Version:    1,
		},
	}

	if payloadFCM != nil {
		options["fcm"] = payloadFCM
	}

	opt, err := json.Marshal(options)
	if err != nil {
		log.Errorf("error marshalling options.", nlogger.Error(err))
		return err
	}

	if opt != nil {
		notification.Options = opt
	}

	// Persist Notification
	err = s.repo.InsertNotification(notification)
	if err != nil {
		log.Errorf("unable to insert notification. err: %v", err)
		// Handle pq.Error
		errCode, _ := nsql.GetPostgresError(err)

		switch errCode {
		case nsql.UniqueError:
			return s.responses.GetError("E_UAL_1").Wrap(err)
		default:
			return err
		}
	}

	return nil
}
