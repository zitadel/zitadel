package user

import (
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/user"
)

func GenderToDomain(gender user.Gender) domain.Gender {
	switch gender {
	case user.Gender_GENDER_DIVERSE:
		return domain.GenderDiverse
	case user.Gender_GENDER_MALE:
		return domain.GenderMale
	case user.Gender_GENDER_FEMALE:
		return domain.GenderFemale
	default:
		return -1
	}
}
