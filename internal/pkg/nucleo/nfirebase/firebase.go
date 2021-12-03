package nfirebase

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"google.golang.org/api/option"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
)

type NucleoFirebase struct {
	app *firebase.App
}

func NewNucleoFirebase(serviceAccountCredential string) (*NucleoFirebase, error) {
	ctx := context.Background()
	opt := option.WithCredentialsJSON([]byte(serviceAccountCredential))
	firebaseSvc, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, err
	}

	return &NucleoFirebase{app: firebaseSvc}, nil
}

func (nf *NucleoFirebase) SendToTarget(payload dto.NotificationCreate) (string, error) {
	ctx := context.Background()
	client, err := nf.app.Messaging(ctx)
	if err != nil {
		return "", fmt.Errorf("NucleoFirebase: %w", err)
	}

	// Create message
	message := messaging.Message{
		Notification: &messaging.Notification{
			Title:    payload.Title,
			Body:     payload.Body,
			ImageURL: payload.ImageURL,
		},
		Token: payload.Token,
		Data:  payload.Data,
	}

	response, err := client.Send(ctx, &message)
	return response, err
}
