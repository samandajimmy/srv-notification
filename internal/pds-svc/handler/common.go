package handler

import (
	"errors"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/ncore"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
	"time"
)

func NewCommon(startTime time.Time, appVersion, appBuildHash string) *Common {
	h := Common{
		startTime: startTime,
		version:   appVersion,
		buildHash: appBuildHash,
	}
	return &h
}

type Common struct {
	startTime time.Time
	version   string
	buildHash string
}

func (c *Common) GetAPIStatus(_ *nhttp.Request) (*nhttp.Response, error) {
	res := nhttp.Success().
		SetData(map[string]string{
			"version":    c.version,
			"build_hash": c.buildHash,
			"uptime":     time.Since(c.startTime).String(),
		})
	return res, nil
}

func (c *Common) ValidateClient(r *nhttp.Request) (*nhttp.Response, error) {

	// Get subject from headers
	subjectID := r.Header.Get(constant.SubjectIDHeader)
	subjectRefID, ok := nval.ParseInt64(subjectID)
	if !ok {
		return nil, errors.New("x-subject-id is required")
	}

	//Get subject role
	subjectRole := r.Header.Get(constant.SubjectRoleHeader)
	role := constant.AdminModifierRole
	if subjectRole != constant.AdminModifierRole {
		role = constant.UserModifierRole
	}

	subject := dto.Subject{
		SubjectID:    subjectID,
		SubjectRefID: subjectRefID,
		SubjectType:  constant.UserSubjectType,
		SubjectRole:  role,
		ModifiedBy: dto.Modifier{
			ID:       subjectID,
			Role:     role,
			FullName: r.Header.Get(constant.SubjectNameHeader),
		},
		Metadata: nil,
	}

	r.SetContextValue(constant.SubjectKey, &subject)

	return nhttp.Continue(), nil
}

func GetSubject(rx *nhttp.Request) (*dto.Subject, error) {
	v := rx.GetContextValue(constant.SubjectKey)
	subject, ok := v.(*dto.Subject)
	if !ok {
		return nil, ncore.NewError("no subject found in request context")
	}
	return subject, nil
}
