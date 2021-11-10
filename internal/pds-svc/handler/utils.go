package handler

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
)

func GetSubject(rx *nhttp.Request) (*dto.Subject, error) {
	v := rx.GetContextValue(constant.SubjectKey)
	subject, ok := v.(*dto.Subject)
	if !ok {
		err := ncore.NewError("no Subject found in request context")
		log.Errorf("an error occurred on getting subject in request context. Error => %s", err)
		return nil, err
	}
	return subject, nil
}
