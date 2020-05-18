package mock

import (
	"github.com/golang/mock/gomock"

	"github.com/caos/zitadel/internal/auth_request/repository"
)

func NewMockAuthRequestRepository(ctrl *gomock.Controller) repository.Repository {
	repo := NewMockRepository(ctrl)
	return repo
}
