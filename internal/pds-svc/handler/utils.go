package handler

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/constant"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/dto"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/ncore"
	"code.nbs.dev/pegadaian/pds/microservice/internal/pkg/nucleo/nhttp"
)

func GetSubject(rx *nhttp.Request) (*dto.Subject, error) {
	v := rx.GetContextValue(constant.SubjectKey)
	subject, ok := v.(*dto.Subject)
	if !ok {
		err := ncore.NewError("no Subject found in request context")
		log.Errorf("an error occurred on getting subject in request context. Error => %s",err)
		return nil, err
	}
	return subject, nil
}
