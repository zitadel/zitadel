package setup

import (
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type IAMSetUp struct {
	Step1  *Step1
	Step2  *Step2
	Step3  *Step3
	Step4  *Step4
	Step5  *Step5
	Step6  *Step6
	Step7  *Step7
	Step8  *Step8
	Step9  *Step9
	Step10 *Step10
}

func (setup *IAMSetUp) steps(currentDone iam_model.Step) ([]step, error) {
	steps := make([]step, 0)
	missingSteps := make([]iam_model.Step, 0)

	for _, step := range []step{
		setup.Step1,
		setup.Step2,
		setup.Step3,
		setup.Step4,
		setup.Step5,
		setup.Step6,
		setup.Step7,
		setup.Step8,
		setup.Step9,
		setup.Step10,
	} {
		if step.step() <= currentDone {
			continue
		}

		if step.isNil() {
			missingSteps = append(missingSteps, step.step())
			continue
		}
		steps = append(steps, step)
	}

	if len(missingSteps) > 0 {
		return nil, errors.ThrowPreconditionFailedf(nil, "SETUP-1nk49", "steps %v not configured", missingSteps)
	}

	return steps, nil
}

type LoginPolicy struct {
	AllowRegister         bool
	AllowUsernamePassword bool
	AllowExternalIdp      bool
}

type User struct {
	FirstName string
	LastName  string
	UserName  string
	Email     string
	Password  string
}

type Org struct {
	Name         string
	Domain       string
	OrgIamPolicy bool
	Users        []User
	Owners       []string
	Projects     []Project
}

type Project struct {
	Name     string
	Users    []User
	Members  []string
	OIDCApps []OIDCApp
}

type OIDCApp struct {
	Name                   string
	RedirectUris           []string
	ResponseTypes          []string
	GrantTypes             []string
	ApplicationType        string
	AuthMethodType         string
	PostLogoutRedirectUris []string
	DevMode                bool
}
