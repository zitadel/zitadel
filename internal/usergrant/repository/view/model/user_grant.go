package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/usergrant/model"
	es_model "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
	"github.com/lib/pq"
)

const (
	UserGrantKeyID            = "id"
	UserGrantKeyUserID        = "user_id"
	UserGrantKeyProjectID     = "project_id"
	UserGrantKeyResourceOwner = "resource_owner"
	UserGrantKeyState         = "state"
	UserGrantKeyOrgName       = "org_name"
	UserGrantKeyRole          = "role_keys"
)

type UserGrantView struct {
	ID            string         `json:"-" gorm:"column:id;primary_key"`
	ResourceOwner string         `json:"-" gorm:"resource_owner"`
	UserID        string         `json:"userId" gorm:"user_id"`
	ProjectID     string         `json:"projectId" gorm:"column:project_id"`
	UserName      string         `json:"-" gorm:"column:user_name"`
	FirstName     string         `json:"-" gorm:"column:first_name"`
	LastName      string         `json:"-" gorm:"column:last_name"`
	Email         string         `json:"-" gorm:"column:email"`
	ProjectName   string         `json:"-" gorm:"column:project_name"`
	OrgName       string         `json:"-" gorm:"column:org_name"`
	RoleKeys      pq.StringArray `json:"roleKeys" gorm:"column:role_keys"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:grant_state"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func UserGrantFromModel(grant *model.UserGrantView) *UserGrantView {
	return &UserGrantView{
		ID:            grant.ID,
		ResourceOwner: grant.ResourceOwner,
		UserID:        grant.UserID,
		ProjectID:     grant.ProjectID,
		ChangeDate:    grant.ChangeDate,
		CreationDate:  grant.CreationDate,
		State:         int32(grant.State),
		UserName:      grant.UserName,
		FirstName:     grant.FirstName,
		LastName:      grant.LastName,
		Email:         grant.Email,
		ProjectName:   grant.ProjectName,
		OrgName:       grant.OrgName,
		RoleKeys:      grant.RoleKeys,
		Sequence:      grant.Sequence,
	}
}

func UserGrantToModel(grant *UserGrantView) *model.UserGrantView {
	return &model.UserGrantView{
		ID:            grant.ID,
		ResourceOwner: grant.ResourceOwner,
		UserID:        grant.UserID,
		ProjectID:     grant.ProjectID,
		ChangeDate:    grant.ChangeDate,
		CreationDate:  grant.CreationDate,
		State:         model.UserGrantState(grant.State),
		UserName:      grant.UserName,
		FirstName:     grant.FirstName,
		LastName:      grant.LastName,
		Email:         grant.Email,
		ProjectName:   grant.ProjectName,
		OrgName:       grant.OrgName,
		RoleKeys:      grant.RoleKeys,
		Sequence:      grant.Sequence,
	}
}

func UserGrantsToModel(grants []*UserGrantView) []*model.UserGrantView {
	result := make([]*model.UserGrantView, len(grants))
	for i, g := range grants {
		result[i] = UserGrantToModel(g)
	}
	return result
}

func (g *UserGrantView) AppendEvent(event *models.Event) (err error) {
	g.ChangeDate = event.CreationDate
	g.Sequence = event.Sequence
	switch event.Type {
	case es_model.UserGrantAdded:
		g.State = int32(model.UserGrantStateActive)
		g.CreationDate = event.CreationDate
		g.setRootData(event)
		err = g.setData(event)
	case es_model.UserGrantChanged, es_model.UserGrantCascadeChanged:
		err = g.setData(event)
	case es_model.UserGrantDeactivated:
		g.State = int32(model.UserGrantStateInactive)
	case es_model.UserGrantReactivated:
		g.State = int32(model.UserGrantStateActive)
	}
	return err
}

func (u *UserGrantView) setRootData(event *models.Event) {
	u.ID = event.AggregateID
	u.ResourceOwner = event.ResourceOwner
}

func (u *UserGrantView) setData(event *models.Event) error {
	if err := json.Unmarshal(event.Data, u); err != nil {
		logging.Log("EVEN-l9sw4").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(nil, "MODEL-7xhke", "could not unmarshal data")
	}
	return nil
}
