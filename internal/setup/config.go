package setup

import (
	"github.com/caos/zitadel/internal/v2/command"
	"github.com/caos/zitadel/internal/v2/domain"
)

type IAMSetUp struct {
	Step1 *command.Step1
	Step2 *command.Step2
	Step3 *command.Step3
	Step4 *command.Step4
	Step5 *command.Step5
	Step6 *command.Step6
	Step7 *command.Step7
	Step8 *command.Step8
	Step9 *command.Step9
}

func (setup *IAMSetUp) Steps(currentDone domain.Step) ([]command.Step, error) {
	steps := make([]command.Step, 0)

	for _, step := range []command.Step{
		setup.Step1,
		setup.Step2,
		setup.Step3,
		setup.Step4,
		setup.Step5,
		setup.Step6,
		setup.Step7,
		setup.Step8,
		setup.Step9,
	} {
		if step.Step() <= currentDone {
			continue
		}
		steps = append(steps, step)
	}
	return steps, nil
}

//
//type LoginPolicy struct {
//	AllowRegister         bool
//	AllowUsernamePassword bool
//	AllowExternalIdp      bool
//}
//
//type User struct {
//	FirstName string
//	LastName  string
//	UserName  string
//	Email     string
//	Password  string
//}
//
//type Org struct {
//	Name         string
//	Domain       string
//	OrgIamPolicy bool
//	Users        []User
//	Owners       []string
//	Projects     []Project
//}
//
//type Project struct {
//	Name     string
//	Users    []User
//	Members  []string
//	OIDCApps []OIDCApp
//}
//
//type OIDCApp struct {
//	Name                   string
//	RedirectUris           []string
//	ResponseTypes          []string
//	GrantTypes             []string
//	ApplicationType        string
//	AuthMethodType         string
//	PostLogoutRedirectUris []string
//	DevMode                bool
//}
