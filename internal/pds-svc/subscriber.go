package pds_svc

import (
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/pubsub"
)

func setUpSubscriber(sub *gochannel.GoChannel, services contract.ServiceMap) {
	// Init subscriber handlers
	sendEmailHandler := pubsub.NewSendEmailHandler(sub, services.Email)

	// Start listening
	go sendEmailHandler.Listen()
}
