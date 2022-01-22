package pds_svc

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/pubsub"
)

func SetUpSubscriber(sub message.Subscriber, service *contract.Service) {
	// Init subscriber handlers
	sendEmailHandler := pubsub.NewSendEmailHandler(sub, service)
	sendFcmHandler := pubsub.NewSendFcmPushHandler(sub, service)
	sendNotification := pubsub.NewSendNotificationHandler(sub, service)

	// Start listening
	go sendEmailHandler.Listen()
	go sendFcmHandler.Listen()
	go sendNotification.Listen()
}
