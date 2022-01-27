package notification

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nbs-go/nlogger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/logger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/convert"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
	"time"
)

func (s *ServiceContext) CreateNotification(payload dto.SendNotificationOptionsRequest) (*dto.DetailNotificationResponse, error) {
	// Generate temporary id
	id, err := uuid.NewUUID()
	if err != nil {
		s.log.Error("error when get uuid", nlogger.Error(err))
		return nil, err
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
		s.log.Errorf("error marshalling options.", nlogger.Error(err))
		return nil, err
	}

	if opt != nil {
		notification.Options = opt
	}

	// Persist Notification
	err = s.repo.InsertNotification(notification)
	if err != nil {
		s.log.Error("unable to insert notification. err: %v", logger.Error(err))
		return nil, err
	}

	return composeDetailNotification(&notification), nil
}

func (s *ServiceContext) GetDetailNotification(payload dto.GetNotification) (*dto.DetailNotificationResponse, error) {
	// Get detail notification
	notification, err := s.repo.FindNotificationByID(payload.ID)
	if err != nil {
		s.log.Error("error when get notification data", nlogger.Error(err))
		if err == sql.ErrNoRows {
			return nil, s.responses.GetError("E_RES_1")
		}
		return nil, err
	}

	return composeDetailNotification(notification), nil
}

func (s *ServiceContext) DeleteNotification(payload dto.GetNotification) error {
	// Get notification by xid
	notification, err := s.repo.FindNotificationByID(payload.ID)
	if err != nil {
		log.Error("error when get data notification. err: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return s.responses.GetError("E_RES_1")
		}
		return err
	}

	// Delete application
	err = s.repo.DeleteNotificationByID(nval.ParseStringFallback(notification.ID, ""))
	if err != nil {
		panic(fmt.Errorf("failed to delete notification. Error = %w", err))
	}

	return nil
}

func (s *ServiceContext) GetCountNotification(payload dto.GetCountNotification) (*dto.DetailCountNotificationResponse, error) {
	// Get count notification
	count, err := s.repo.CountNotification(payload)
	if err != nil {
		s.log.Error("error when get data count notification. err: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, s.responses.GetError("E_RES_1")
		}
		return nil, err
	}

	return &dto.DetailCountNotificationResponse{
		Count: count,
	}, err
}

func composeDetailNotification(m *model.Notification) *dto.DetailNotificationResponse {
	var readAt int64
	if m.ReadAt.Valid {
		readAt = m.ReadAt.Time.Unix()
	}

	return &dto.DetailNotificationResponse{
		Id:                   m.ID,
		ApplicationId:        m.ApplicationId,
		UserRefId:            m.UserRefId,
		IsRead:               m.IsRead,
		ReadAt:               readAt,
		Options:              m.Options,
		ItemMetadataResponse: convert.ItemMetadataModelToResponse(m.ItemMetadata),
	}
}
