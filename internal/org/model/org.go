package model

import (
	"strings"

	"github.com/golang/protobuf/ptypes/timestamp"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type Org struct {
	es_models.ObjectRoot

	State   OrgState
	Name    string
	Domains []*OrgDomain

	Members                  []*OrgMember
	OrgIamPolicy             *iam_model.OrgIAMPolicy
	LoginPolicy              *iam_model.LoginPolicy
	LabelPolicy              *iam_model.LabelPolicy
	MailTemplate             *iam_model.MailTemplate
	MailTexts                []*iam_model.MailText
	PasswordComplexityPolicy *iam_model.PasswordComplexityPolicy
	PasswordAgePolicy        *iam_model.PasswordAgePolicy
	PasswordLockoutPolicy    *iam_model.PasswordLockoutPolicy

	IDPs []*iam_model.IDPConfig
}
type OrgChanges struct {
	Changes      []*OrgChange
	LastSequence uint64
}

type OrgChange struct {
	ChangeDate        *timestamp.Timestamp `json:"changeDate,omitempty"`
	EventType         string               `json:"eventType,omitempty"`
	Sequence          uint64               `json:"sequence,omitempty"`
	ModifierId        string               `json:"modifierUser,omitempty"`
	ModifierName      string               `json:"-"`
	ModifierLoginName string               `json:"-"`
	Data              interface{}          `json:"data,omitempty"`
}

type OrgState int32

const (
	OrgStateActive OrgState = iota
	OrgStateInactive
)

func NewOrg(id string) *Org {
	return &Org{ObjectRoot: es_models.ObjectRoot{AggregateID: id}, State: OrgStateActive}
}

func (o *Org) IsActive() bool {
	return o.State == OrgStateActive
}

func (o *Org) IsValid() bool {
	return o.Name != ""
}

func (o *Org) GetDomain(domain *OrgDomain) (int, *OrgDomain) {
	for i, d := range o.Domains {
		if d.Domain == domain.Domain {
			return i, d
		}
	}
	return -1, nil
}

func (o *Org) GetIDP(idpID string) (int, *iam_model.IDPConfig) {
	for i, idp := range o.IDPs {
		if idp.IDPConfigID == idpID {
			return i, idp
		}
	}
	return -1, nil
}

func (o *Org) GetPrimaryDomain() *OrgDomain {
	for _, d := range o.Domains {
		if d.Primary {
			return d
		}
	}
	return nil
}

func (o *Org) MemeberByUserID(userID string) (*OrgMember, int) {
	for i, member := range o.Members {
		if member.UserID == userID {
			return member, i
		}
	}
	return nil, -1
}

func (o *Org) nameForDomain(iamDomain string) string {
	return strings.ToLower(strings.ReplaceAll(o.Name, " ", "-") + "." + iamDomain)
}

func (o *Org) AddIAMDomain(iamDomain string) {
	o.Domains = append(o.Domains, &OrgDomain{Domain: o.nameForDomain(iamDomain), Verified: true, Primary: true})
}
