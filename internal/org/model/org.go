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

	DomainPolicy *iam_model.DomainPolicy
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
	ModifierAvatarURL string               `json:"-"`
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

func (o *Org) GetPrimaryDomain() *OrgDomain {
	for _, d := range o.Domains {
		if d.Primary {
			return d
		}
	}
	return nil
}

func (o *Org) nameForDomain(iamDomain string) string {
	return strings.ToLower(strings.ReplaceAll(o.Name, " ", "-") + "." + iamDomain)
}

func (o *Org) AddIAMDomain(iamDomain string) {
	o.Domains = append(o.Domains, &OrgDomain{Domain: o.nameForDomain(iamDomain), Verified: true, Primary: true})
}
