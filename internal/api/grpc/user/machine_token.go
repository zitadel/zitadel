package user

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/pkg/grpc/user"
)

func MachineTokensToPb(tokens []*query.MachineToken) []*user.MachineToken {
	t := make([]*user.MachineToken, len(tokens))
	for i, token := range tokens {
		t[i] = MachineTokenToPb(token)
	}
	return t
}
func MachineTokenToPb(token *query.MachineToken) *user.MachineToken {
	return &user.MachineToken{
		Id:             token.ID,
		Details:        object.ChangeToDetailsPb(token.Sequence, token.ChangeDate, token.ResourceOwner),
		ExpirationDate: timestamppb.New(token.Expiration),
		Scopes:         token.Scopes,
	}
}
