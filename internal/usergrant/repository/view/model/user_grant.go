package model

import (
	"encoding/json"
	"time"

	"github.com/caos/logging"
	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/usergrant/model"
	es_model "github.com/caos/zitadel/internal/usergrant/repository/eventsourcing/model"
)

const (
	UserGrantKeyID            = "id"
	UserGrantKeyUserID        = "user_id"
	UserGrantKeyProjectID     = "project_id"
	UserGrantKeyGrantID       = "grant_id"
	UserGrantKeyResourceOwner = "resource_owner"
	UserGrantKeyState         = "state"
	UserGrantKeyOrgName       = "org_name"
	UserGrantKeyRole          = "role_keys"
	UserGrantKeyUserName      = "user_name"
	UserGrantKeyFirstName     = "first_name"
	UserGrantKeyLastName      = "last_name"
	UserGrantKeyEmail         = "email"
	UserGrantKeyOrgDomain     = "org_primary_domain"
	UserGrantKeyProjectName   = "project_name"
	UserGrantKeyDisplayName   = "display_name"
)

type UserGrantView struct {
	ID                string         `json:"-" gorm:"column:id;primary_key"`
	ResourceOwner     string         `json:"-" gorm:"resource_owner"`
	UserID            string         `json:"userId" gorm:"user_id"`
	ProjectID         string         `json:"projectId" gorm:"column:project_id"`
	GrantID           string         `json:"grantId" gorm:"column:grant_id"`
	UserName          string         `json:"-" gorm:"column:user_name"`
	FirstName         string         `json:"-" gorm:"column:first_name"`
	LastName          string         `json:"-" gorm:"column:last_name"`
	DisplayName       string         `json:"-" gorm:"column:display_name"`
	Email             string         `json:"-" gorm:"column:email"`
	ProjectName       string         `json:"-" gorm:"column:project_name"`
	ProjectOwner      string         `json:"-" gorm:"column:project_owner"`
	OrgName           string         `json:"-" gorm:"column:org_name"`
	OrgPrimaryDomain  string         `json:"-" gorm:"column:org_primary_domain"`
	RoleKeys          pq.StringArray `json:"roleKeys" gorm:"column:role_keys"`
	AvatarKey         string         `json:"-" gorm:"column:avatar_key"`
	UserResourceOwner string         `json:"-" gorm:"column:user_resource_owner"`

	CreationDate time.Time `json:"-" gorm:"column:creation_date"`
	ChangeDate   time.Time `json:"-" gorm:"column:change_date"`
	State        int32     `json:"-" gorm:"column:grant_state"`

	Sequence uint64 `json:"-" gorm:"column:sequence"`
}

func UserGrantToModel(grant *UserGrantView, prefixAvatarURL string) *model.UserGrantView {
	return &model.UserGrantView{
		ID:               grant.ID,
		ResourceOwner:    grant.ResourceOwner,
		UserID:           grant.UserID,
		ProjectID:        grant.ProjectID,
		ChangeDate:       grant.ChangeDate,
		CreationDate:     grant.CreationDate,
		State:            model.UserGrantState(grant.State),
		UserName:         grant.UserName,
		FirstName:        grant.FirstName,
		LastName:         grant.LastName,
		DisplayName:      grant.DisplayName,
		Email:            grant.Email,
		ProjectName:      grant.ProjectName,
		OrgName:          grant.OrgName,
		OrgPrimaryDomain: grant.OrgPrimaryDomain,
		RoleKeys:         grant.RoleKeys,
		AvatarURL:        domain.AvatarURL(prefixAvatarURL, grant.ResourceOwner, grant.AvatarKey),
		Sequence:         grant.Sequence,
		GrantID:          grant.GrantID,
	}
}

func UserGrantsToModel(grants []*UserGrantView, prefixAvatarURL string) []*model.UserGrantView {
	result := make([]*model.UserGrantView, len(grants))
	for i, g := range grants {
		result[i] = UserGrantToModel(g, prefixAvatarURL)
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
