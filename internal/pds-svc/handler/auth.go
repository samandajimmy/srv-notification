package handler

import (
	"errors"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/constant"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/contract"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pds-svc/dto"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nhttp"
	"repo.pegadaian.co.id/ms-pds/srv-notification/internal/pkg/nucleo/nval"
)

func NewAuth(authService contract.AuthService) *Auth {
	a := Auth{
		authService: authService,
	}
	return &a
}

type Auth struct {
	router      *nhttp.Router
	authService contract.AuthService
}

func (h *Auth) ValidateClient(r *nhttp.Request) (*nhttp.Response, error) {
	// Extract Basic Auth
	clientID, clientSecret, err := nhttp.ExtractBasicAuth(r.Request)
	if err != nil {
		return nil, err
	}

	// Authentication app
	err = h.authService.ValidateClient(dto.ClientCredential{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		log.Errorf("an error occurred on authentication app. Error => %s", err)
		return nil, err
	}

	// Get subject from headers
	subjectID := r.Header.Get(constant.SubjectIDHeader)
	subjectRefID, ok := nval.ParseInt64(subjectID)
	if !ok {
		err = errors.New("x-subject-id is required")
		log.Errorf("an error occurred on getting subject from headers. Error => %s", err)
		return nil, err
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
