package notification

import (
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/logger"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
)

func (s *ServiceContext) SendPushNotificationByTarget(payload dto.SendPushNotification) error {
	// TODO: Load credential by config id from database
	client, err := s.newFcmClient(s.config.Firebase.ServiceAccountCredential)
	if err != nil {
		return ncore.TraceError(err)
	}

	// Compose message
	message := composeFcmMessage(payload)

	// Send
	result, err := client.Send(s.ctx, &message)
	if err != nil {
		s.log.Error("failed to send fcm message", logger.Error(err))
		return ncore.TraceError(err)
	}
	s.log.Debugf("success sending fcm message. Result = %s", result)

	return nil
}

func (s *ServiceContext) newFcmClient(credential string) (*messaging.Client, error) {
	opt := option.WithCredentialsJSON([]byte(credential))
	fcm, err := firebase.NewApp(s.ctx, nil, opt)
	if err != nil {
		s.log.Error("failed to init fcm client", logger.Error(err))
		return nil, ncore.TraceError(err)
	}
	return fcm.Messaging(s.ctx)
}

func composeFcmMessage(payload dto.SendPushNotification) messaging.Message {
	return messaging.Message{
		Notification: &messaging.Notification{
			Title:    payload.Title,
			Body:     payload.Body,
			ImageURL: payload.ImageURL,
		},
		Token: payload.Token,
		Data:  payload.Data,
	}
}
