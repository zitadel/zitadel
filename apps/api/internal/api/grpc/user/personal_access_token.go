package user

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/user"
)

func PersonalAccessTokensToPb(tokens []*query.PersonalAccessToken) []*user.PersonalAccessToken {
	t := make([]*user.PersonalAccessToken, len(tokens))
	for i, token := range tokens {
		t[i] = PersonalAccessTokenToPb(token)
	}
	return t
}
func PersonalAccessTokenToPb(token *query.PersonalAccessToken) *user.PersonalAccessToken {
	return &user.PersonalAccessToken{
		Id:             token.ID,
		Details:        object.ToViewDetailsPb(token.Sequence, token.CreationDate, token.ChangeDate, token.ResourceOwner),
		ExpirationDate: timestamppb.New(token.Expiration),
		Scopes:         token.Scopes,
	}
}
