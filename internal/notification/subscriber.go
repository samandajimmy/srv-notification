package notification

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/contract"
	pubsub2 "repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/pubsub"
)

func SetUpSubscriber(sub message.Subscriber, service *contract.Service) {
	// Init subscriber handlers
	sendEmailHandler := pubsub2.NewSendEmailHandler(sub, service)
	sendFcmHandler := pubsub2.NewSendFcmPushHandler(sub, service)

	// Start listening
	go sendEmailHandler.Listen()
	go sendFcmHandler.Listen()
}
