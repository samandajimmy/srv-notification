package contract

import (
	"code.nbs.dev/pegadaian/pds/microservice/internal/pds-svc/dto"
)

type AuthService interface {
	ValidateClient(payload dto.ClientCredential) error
}
