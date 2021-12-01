package pds_svc

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/pubsub"
)

func setUpSubscriber(sub message.Subscriber, services contract.ServiceMap) {
	// Init subscriber handlers
	sendEmailHandler := pubsub.NewSendEmailHandler(sub, services.Email)
	sendFcmHandler := pubsub.NewSendFcmPushHandler(sub, services.Notification)

	// Start listening
	go sendEmailHandler.Listen()
	go sendFcmHandler.Listen()
}
