package setup

import (
	"context"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/domain"
)

type AdminOrg struct {
	cmds *command.Commands
	Org  struct {
		Name   string
		Domain string
	}
	User struct {
		Username  string
		FirstName string
		LastName  string
		UserName  string
		Email     string
		Password  string
	}
}

func (mig *AdminOrg) Execute() error {
	_, err := mig.cmds.SetUpOrg(context.TODO(),
		&domain.Org{
			Name:    mig.Org.Name,
			Domains: []*domain.OrgDomain{{Domain: mig.Org.Domain}},
		},
		&domain.Human{
			Username: mig.User.Username,
			Profile: &domain.Profile{
				FirstName: mig.User.FirstName,
				LastName:  mig.User.LastName,
			},
			Email: &domain.Email{
				EmailAddress:    mig.User.Email,
				IsEmailVerified: true,
			},
			Password: domain.NewPassword(mig.User.Password),
		}, nil, nil, nil, false)

	return err
}
