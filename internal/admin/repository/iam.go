package repository

import (
	"context"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type IamRepository interface {
	SearchIamMembers(ctx context.Context, request *iam_model.IamMemberSearchRequest) (*iam_model.IamMemberSearchResponse, error)
	AddIamMember(ctx context.Context, member *iam_model.IamMember) (*iam_model.IamMember, error)
	ChangeIamMember(ctx context.Context, member *iam_model.IamMember) (*iam_model.IamMember, error)
	RemoveIamMember(ctx context.Context, userID string) error

	GetIamMemberRoles() []string
}
