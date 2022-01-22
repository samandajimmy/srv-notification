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
	payloadSMTP := payload.Options.SMTP
	//
	now := time.Now()
	var metadata map[string]interface{}
	// Initialize data to insert
	notification := model.Notification{
		ID:        id,
		UserRefId: payload.UserId,
		Title:     payload.Title,
		IsRead:    false,
		ReadAt:    sql.NullTime{},
		Metadata:  nil,
		ItemMetadata: model.ItemMetadata{
			CreatedAt:  now,
			UpdatedAt:  now,
			ModifiedBy: &model.Modifier{ID: "", Role: "", FullName: ""},
			Version:    1,
		},
	}

	if payload.Content != "" {
		notification.Content = payload.Content
	}

	if payload.ContentShort != "" {
		notification.ContentShort = payload.ContentShort
	}

	if payload.ContentEncoded != "" {
		notification.ContentEncoded = payload.ContentEncoded
	}

	if payloadFCM != nil {
		if payloadFCM.Data != nil {
			payloadFcmData := make(map[string]interface{}, len(payloadFCM.Data))
			for keyData, valueData := range payloadFCM.Data {
				payloadFcmData[keyData] = valueData
			}
			metadata = payloadFcmData
		}

		additionalButton := dto.AdditionalButton{
			ButtonLabel:   payloadFCM.AdditionalButton.ButtonLabel,
			TransactionId: payloadFCM.AdditionalButton.TransactionId,
			ScreenName:    payloadFCM.AdditionalButton.ScreenName,
		}

		additionalBtn, err := json.Marshal(additionalButton)
		if err != nil {
			log.Errorf("error marshal additionalBtn.", nlogger.Error(err))
			return err
		}
		// assign additionalButton to Metadata
		metadata["additionalButton"] = additionalBtn
		metadata["token"] = payloadFCM.Token
		metadata["imageUrl"] = payloadFCM.ImageUrl

		metaData, err := json.Marshal(metadata)
		if err != nil {
			log.Errorf("error marshal metaData.", nlogger.Error(err))
			return err
		}
		notification.Metadata = metaData
	}

	if payloadSMTP != nil {
		dtoFrom := dto.FromFormat{
			Name:  payloadSMTP.From.Name,
			Email: payloadSMTP.From.Email,
		}
		from, err := json.Marshal(dtoFrom)
		if err != nil {
			log.Errorf("error marshal fromEmail.", nlogger.Error(err))
			return err
		}
		metadata["from"] = from
		metadata["to"] = payloadSMTP.To
		metadata["attachment"] = payloadSMTP.Attachment
		metadata["mimeType"] = payloadSMTP.MimeType
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
