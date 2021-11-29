package repository

import (
	"context"
)

type ApplicationRepository interface {
	AuthorizeClientIDSecret(ctx context.Context, clientID, secret string) error
}
