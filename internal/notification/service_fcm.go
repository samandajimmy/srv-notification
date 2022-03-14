package notification

import (
	"encoding/json"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"google.golang.org/api/option"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"

	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
)

func (s *ServiceContext) SendPushNotificationByTarget(payload *dto.SendPushNotification) error {
	var config json.RawMessage
	err := s.loadClientConfig(payload.ApplicationId, constant.Firebase, &config)
	if err != nil {
		return ncore.TraceError(err)
	}

	client, err := s.newFcmClient(string(config))
	if err != nil {
		s.log.Error("failed to initialize new client", logOption.Error(err))
		return ncore.TraceError(err)
	}

	// Compose message
	message := composeFcmMessage(payload)

	// Send
	result, err := client.Send(s.ctx, &message)
	if err != nil {
		s.log.Error("failed to send fcm message", logOption.Error(err))
		return ncore.TraceError(err)
	}
	s.log.Debugf("success sending fcm message. Result = %s", result)

	return nil
}

func (s *ServiceContext) newFcmClient(credential string) (*messaging.Client, error) {
	opt := option.WithCredentialsJSON([]byte(credential))
	fcm, err := firebase.NewApp(s.ctx, nil, opt)
	if err != nil {
		s.log.Error("failed to init fcm client", logOption.Error(err))
		return nil, ncore.TraceError(err)
	}
	return fcm.Messaging(s.ctx)
}

func composeFcmMessage(payload *dto.SendPushNotification) messaging.Message {
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
