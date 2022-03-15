package handler

import (
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/notification/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
	"time"
)

const (
	AnonymousUserId       = "ANON"
	AnonymousUserRefId    = 0
	AnonymousUserFullName = "Anonymous User"
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

func (c *Common) ParseSubject(r *nhttp.Request) (*nhttp.Response, error) {
	// Get subject from headers
	id := r.Header.Get(constant.SubjectIDHeader)
	if id == "" {
		id = AnonymousUserId
	}

	// Get subject reference id
	refId, ok := nval.ParseInt64(id)
	if !ok {
		refId = AnonymousUserRefId
	}

	// Get subject role and determine subject type
	role := r.Header.Get(constant.SubjectRoleHeader)
	var subjectType constant.SubjectType
	switch role {
	case constant.AdminModifierRole, constant.UserModifierRole:
		subjectType = constant.UserSubjectType
	case constant.SystemModifierRole:
		subjectType = constant.SystemSubjectType
	default:
		// Fallback to anonymous user
		subjectType = constant.UserSubjectType
		role = constant.UserModifierRole
	}

	// Get name
	fullName := r.Header.Get(constant.SubjectNameHeader)
	if fullName == "" {
		fullName = AnonymousUserFullName
	}

	subject := dto.Subject{
		Id:          id,
		RefId:       refId,
		Role:        role,
		FullName:    fullName,
		SubjectType: subjectType,
		Metadata:    map[string]string{},
		SessionID:   0,
	}

	r.SetContextValue(constant.SubjectKey, &subject)

	return nhttp.Continue(), nil
}

func GetSubject(rx *nhttp.Request) *dto.Subject {
	v := rx.GetContextValue(constant.SubjectKey)
	subject, ok := v.(*dto.Subject)
	if !ok {
		// Return anonymous subject
		return &dto.Subject{
			Id:          AnonymousUserId,
			RefId:       AnonymousUserRefId,
			Role:        constant.UserModifierRole,
			FullName:    AnonymousUserFullName,
			SubjectType: constant.UserSubjectType,
			SessionID:   0,
			Metadata:    map[string]string{},
		}
	}
	return subject
}
