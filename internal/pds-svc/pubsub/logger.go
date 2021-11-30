package pubsub

import "repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nlogger"

var logger nlogger.Logger

func init() {
	logger = nlogger.Get()
}
