package notification

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nbs-go/nlogger"
	dto "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/logger"
	model "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/model"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
	"time"
)

func (s *ServiceContext) CreateNotification(payload *dto.SendNotificationOptionsRequest) (*dto.DetailNotificationResponse, error) {
	// Generate temporary id
	id, err := uuid.NewUUID()
	if err != nil {
		s.log.Error("error when get uuid", nlogger.Error(err))
		return nil, err
	}

	payloadFCM := payload.Options.FCM

	// init options
	options := map[string]interface{}{}
	// Initialize data to insert
	notification := model.Notification{
		ID:            id,
		ApplicationId: payload.Auth.ID,
		UserRefId:     payload.UserId,
		IsRead:        false,
		ReadAt:        sql.NullTime{},
		BaseField:     model.NewBaseField(&model.Modifier{ID: "", Role: "", FullName: ""}),
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

func (s *ServiceContext) GetDetailNotification(payload *dto.GetNotification) (*dto.DetailNotificationResponse, error) {
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

func (s *ServiceContext) DeleteNotification(payload *dto.GetNotification) error {
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

func (s *ServiceContext) CountNotification(payload *dto.GetCountNotification) (*dto.DetailCountNotificationResponse, error) {
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

func (s *ServiceContext) ListNotification(options *dto.ListPayload) (*dto.ListNotificationResponse, error) {
	// Handle sort request
	rulesSortBy := []string{
		"createdAt",
		"updatedAt",
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

	// Get list notification
	result, err := s.repo.FindNotification(options)
	if err != nil {
		s.log.Error("failed to find data notification", nlogger.Error(err))
		return nil, ncore.TraceError(err)
	}

	// Set item response
	items := make([]*dto.DetailNotificationResponse, len(result.Rows))
	for idx, row := range result.Rows {
		item := composeDetailNotification(&row)
		items[idx] = item
	}

	return &dto.ListNotificationResponse{
		Items:    items,
		Metadata: dto.ToListMetadata(options, result.Count),
	}, nil
}

func (s *ServiceContext) UpdateIsRead(payload *dto.UpdateIsReadNotification) (*dto.DetailNotificationResponse, error) {
	// Get detail notification
	notification, err := s.repo.FindNotificationByID(payload.ID)
	if err != nil {
		s.log.Error("error when get notification data", nlogger.Error(err))
		if err == sql.ErrNoRows {
			return nil, s.responses.GetError("E_RES_1")
		}
		return nil, err
	}

	// Set notification read at and is read true.
	notification.ReadAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	notification.IsRead = true
	notification.UpdatedAt = time.Now()
	notification.Version += 1

	err = s.repo.UpdateNotificationByID(notification)
	if err != nil {
		s.log.Error("error when update notification data", nlogger.Error(err))
		return nil, s.responses.GetError("E_RES_2")
	}

	return composeDetailNotification(notification), nil
}

func composeDetailNotification(m *model.Notification) *dto.DetailNotificationResponse {
	var readAt int64
	if m.ReadAt.Valid {
		readAt = m.ReadAt.Time.Unix()
	}

	return &dto.DetailNotificationResponse{
		Id:            m.ID,
		ApplicationId: m.ApplicationId,
		UserRefId:     m.UserRefId,
		IsRead:        m.IsRead,
		ReadAt:        readAt,
		Options:       m.Options,
		BaseField:     model.ToBaseFieldDTO(m.BaseField),
	}
}
