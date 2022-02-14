package handler

import (
	"github.com/hetiansu5/urlquery"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

func getListPayload(rx *nhttp.Request) (*dto.ListPayload, error) {
	// Parse query
	var payload dto.ListPayload
	err := urlquery.Unmarshal([]byte(rx.URL.RawQuery), &payload)
	if err != nil {
		return nil, ncore.TraceError(err)
	}

	// Normalize Limit
	if payload.Limit <= 0 {
		payload.Limit = 10
	}

	// Normalize skip
	if payload.Skip < 0 {
		payload.Skip = 0
	}

	return &payload, nil
}
