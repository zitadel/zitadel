package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_model "github.com/caos/zitadel/internal/user/model"
)

func (r *CommandSide) AddHuman(ctx context.Context, user *usr_model.Human) (*usr_model.Human, error) {

	return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-8K0df", "Errors.User.TypeUndefined")
}
