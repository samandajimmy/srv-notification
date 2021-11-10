package service

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nfirebase"
)

type Notification struct {
	FirebaseConfig contract.FirebaseConfig
	Firebase       *nfirebase.NucleoFirebase
	responses      *ncore.ResponseMap
}

func (s *Notification) HasInitialized() bool {
	return true
}

func (s *Notification) Init(app *contract.PdsApp) error {
	var err error
	s.FirebaseConfig.Key = app.Config.Firebase.Key
	s.Firebase, err = nfirebase.NewNucleoFirebase(s.FirebaseConfig.Key)
	if err != nil {
		log.Errorf("Error initialise firebase config %v", err)
		return err
	}
	s.responses = app.Responses

	return nil
}

func (s *Notification) SendNotificationByToken(payload dto.NotificationCreate) error {

	// Send notification by token
	_, err := s.Firebase.SendToTarget(payload)
	if err != nil {
		log.Errorf("Error when sending notification by token %v", err)
		return err
	}

	return nil
}
