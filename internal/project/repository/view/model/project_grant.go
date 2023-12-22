package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/project/model"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	ProjectGrantKeyProjectID     = "project_id"
	ProjectGrantKeyGrantID       = "grant_id"
	ProjectGrantKeyOrgID         = "org_id"
	ProjectGrantKeyResourceOwner = "resource_owner"
	ProjectGrantKeyName          = "project_name"
	ProjectGrantKeyRoleKeys      = "granted_role_keys"
)

type ProjectGrantView struct {
	GrantID           string                     `json:"-" gorm:"column:grant_id;primary_key"`
	ProjectID         string                     `json:"-" gorm:"column:project_id"`
	OrgID             string                     `json:"-" gorm:"column:org_id"`
	Name              string                     `json:"name" gorm:"column:project_name"`
	CreationDate      time.Time                  `json:"-" gorm:"column:creation_date"`
	ChangeDate        time.Time                  `json:"-" gorm:"column:change_date"`
	State             int32                      `json:"-" gorm:"column:project_state"`
	ResourceOwner     string                     `json:"-" gorm:"column:resource_owner"`
	ResourceOwnerName string                     `json:"-" gorm:"column:resource_owner_name"`
	OrgName           string                     `json:"-" gorm:"column:org_name"`
	Sequence          uint64                     `json:"-" gorm:"column:sequence"`
	GrantedRoleKeys   database.TextArray[string] `json:"-" gorm:"column:granted_role_keys"`
}

type ProjectGrant struct {
	GrantID      string   `json:"grantId"`
	GrantedOrgID string   `json:"grantedOrgId"`
	RoleKeys     []string `json:"roleKeys"`
	InstanceID   string   `json:"instanceID"`
}

func (p *ProjectGrantView) AppendEvent(event *models.Event) (err error) {
	p.ChangeDate = event.CreationDate
	p.Sequence = event.Seq
	switch event.Type() {
	case project.GrantAddedType:
		p.State = int32(model.ProjectStateActive)
		p.CreationDate = event.CreationDate
		p.setRootData(event)
		err = p.setProjectGrantData(event)
	case project.GrantChangedType, project.GrantCascadeChangedType:
		err = p.setProjectGrantData(event)
	case project.GrantDeactivatedType:
		p.State = int32(model.ProjectStateInactive)
	case project.GrantReactivatedType:
		p.State = int32(model.ProjectStateActive)
	}
	return err
}

func (p *ProjectGrantView) setRootData(event *models.Event) {
	p.ProjectID = event.AggregateID
	p.ResourceOwner = event.ResourceOwner
}

func (p *ProjectGrantView) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-dlo92").WithError(err).Error("could not unmarshal event data")
		return err
	}
	return nil
}

func (p *ProjectGrantView) setProjectGrantData(event *models.Event) error {
	grant := new(ProjectGrant)
	err := grant.SetData(event)
	if err != nil {
		return err
	}
	if grant.GrantedOrgID != "" {
		p.OrgID = grant.GrantedOrgID
	}
	p.GrantID = grant.GrantID
	p.GrantedRoleKeys = grant.RoleKeys
	return nil
}

func (p *ProjectGrant) SetData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, p); err != nil {
		logging.Log("EVEN-dlo92").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-s9ols", "Could not unmarshal data")
	}
	return nil
}
