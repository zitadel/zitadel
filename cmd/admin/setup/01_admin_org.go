package setup

import (
	"context"

	command "github.com/caos/zitadel/internal/command/v2"
	"github.com/caos/zitadel/internal/command/v2/user"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"golang.org/x/text/language"
)

type AdminOrg struct {
	cmd  *command.Command
	Org  admin.SetUpOrgRequest_Org
	User struct {
		Human admin.SetUpOrgRequest_Human
	}
}

func (mig *AdminOrg) Execute(ctx context.Context) error {
	_, err := mig.cmd.SetUpOrg(ctx, &command.OrgSetup{
		Name:   mig.Org.Name,
		Domain: mig.Org.Domain,
		Human: user.AddHuman{
			Username:  mig.User.Human.UserName,
			FirstName: mig.User.Human.Profile.FirstName,
			LastName:  mig.User.Human.Profile.LastName,
			Email:     mig.User.Human.Email.Email,
			//TODO: email verified
			NickName:      mig.User.Human.Profile.NickName,
			DisplayName:   mig.User.Human.Profile.DisplayName,
			PreferredLang: language.Make(mig.User.Human.Profile.PreferredLanguage),
			Gender:        domain.Gender(mig.User.Human.Profile.Gender),
			Phone:         mig.User.Human.Phone.Phone,
			//TODO: phone verified
			// Password: human.Password,
		},
	})

	return err
}

func (mig *AdminOrg) String() string {
	return "01_admin_org"
}
