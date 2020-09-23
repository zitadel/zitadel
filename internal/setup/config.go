package setup

import (
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type IAMSetUp struct {
	Step1 *Step1
	//TODO: label policy
	// Step2 *Step2
}

func (setup *IAMSetUp) steps(currentDone iam_model.Step) ([]step, error) {
	steps := make([]step, 0)
	missingSteps := make([]iam_model.Step, 0)

	for _, step := range []step{
		setup.Step1,
		//TODO: label policy
		// setup.Step2,
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

//TODO: label policy
// type LabelPolicy struct {
// 	PrimaryColor  string
// 	SecondayColor string
// }

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
