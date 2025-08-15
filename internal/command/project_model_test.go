package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func TestProjectWriteModel_Reduce(t *testing.T) {
	type fields struct {
		WriteModel             eventstore.WriteModel
		Name                   string
		ProjectRoleAssertion   bool
		ProjectRoleCheck       bool
		HasProjectCheck        bool
		PrivateLabelingSetting domain.PrivateLabelingSetting
		State                  domain.ProjectState
	}
	type args struct {
		events []eventstore.Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ProjectWriteModel
	}{
		{
			name: "org removed event resets project write model",
			fields: fields{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "project-id",
					ResourceOwner: "org-id",
				},
				Name:                   "Test Project",
				ProjectRoleAssertion:   true,
				ProjectRoleCheck:       true,
				HasProjectCheck:        true,
				PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
				State:                  domain.ProjectStateActive,
			},
			args: args{
				events: []eventstore.Event{
					org.NewOrgRemovedEvent(
						context.Background(),
						&eventstore.Aggregate{
							ID:            "org-id",
							Type:          org.AggregateType,
							ResourceOwner: "org-id",
						},
						"Test Org",
						[]string{},
						false,
						[]string{},
						nil,
						[]string{},
					),
				},
			},
			want: &ProjectWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "project-id",
					ResourceOwner: "org-id",
				},
				Name:                   "",
				ProjectRoleAssertion:   false,
				ProjectRoleCheck:       false,
				HasProjectCheck:        false,
				PrivateLabelingSetting: domain.PrivateLabelingSetting(0),
				State:                  domain.ProjectState(0),
			},
		},
		{
			name: "project added event sets project data",
			fields: fields{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "project-id",
					ResourceOwner: "org-id",
				},
			},
			args: args{
				events: []eventstore.Event{
					project.NewProjectAddedEvent(
						context.Background(),
						&eventstore.Aggregate{
							ID:            "project-id",
							Type:          project.AggregateType,
							ResourceOwner: "org-id",
						},
						"Test Project",
						true,
						true,
						true,
						domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
					),
				},
			},
			want: &ProjectWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "project-id",
					ResourceOwner: "org-id",
				},
				Name:                   "Test Project",
				ProjectRoleAssertion:   true,
				ProjectRoleCheck:       true,
				HasProjectCheck:        true,
				PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
				State:                  domain.ProjectStateActive,
			},
		},
		{
			name: "project added then org removed resets model",
			fields: fields{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "project-id",
					ResourceOwner: "org-id",
				},
			},
			args: args{
				events: []eventstore.Event{
					project.NewProjectAddedEvent(
						context.Background(),
						&eventstore.Aggregate{
							ID:            "project-id",
							Type:          project.AggregateType,
							ResourceOwner: "org-id",
						},
						"Test Project",
						true,
						true,
						true,
						domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
					),
					org.NewOrgRemovedEvent(
						context.Background(),
						&eventstore.Aggregate{
							ID:            "org-id",
							Type:          org.AggregateType,
							ResourceOwner: "org-id",
						},
						"Test Org",
						[]string{},
						false,
						[]string{},
						nil,
						[]string{},
					),
				},
			},
			want: &ProjectWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "project-id",
					ResourceOwner: "org-id",
				},
				Name:                   "",
				ProjectRoleAssertion:   false,
				ProjectRoleCheck:       false,
				HasProjectCheck:        false,
				PrivateLabelingSetting: domain.PrivateLabelingSetting(0),
				State:                  domain.ProjectState(0),
			},
		},
		{
			name: "project removed event sets state to removed",
			fields: fields{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "project-id",
					ResourceOwner: "org-id",
				},
				Name:                   "Test Project",
				ProjectRoleAssertion:   true,
				ProjectRoleCheck:       true,
				HasProjectCheck:        true,
				PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
				State:                  domain.ProjectStateActive,
			},
			args: args{
				events: []eventstore.Event{
					project.NewProjectRemovedEvent(
						context.Background(),
						&eventstore.Aggregate{
							ID:            "project-id",
							Type:          project.AggregateType,
							ResourceOwner: "org-id",
						},
						"Test Project",
						[]*eventstore.UniqueConstraint{},
					),
				},
			},
			want: &ProjectWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID:   "project-id",
					ResourceOwner: "org-id",
				},
				Name:                   "Test Project",
				ProjectRoleAssertion:   true,
				ProjectRoleCheck:       true,
				HasProjectCheck:        true,
				PrivateLabelingSetting: domain.PrivateLabelingSettingAllowLoginUserResourceOwnerPolicy,
				State:                  domain.ProjectStateRemoved,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wm := &ProjectWriteModel{
				WriteModel:             tt.fields.WriteModel,
				Name:                   tt.fields.Name,
				ProjectRoleAssertion:   tt.fields.ProjectRoleAssertion,
				ProjectRoleCheck:       tt.fields.ProjectRoleCheck,
				HasProjectCheck:        tt.fields.HasProjectCheck,
				PrivateLabelingSetting: tt.fields.PrivateLabelingSetting,
				State:                  tt.fields.State,
			}

			// Set the events for the write model to process
			wm.Events = tt.args.events

			err := wm.Reduce()
			assert.NoError(t, err)

			assert.Equal(t, tt.want.AggregateID, wm.AggregateID)
			assert.Equal(t, tt.want.ResourceOwner, wm.ResourceOwner)
			assert.Equal(t, tt.want.Name, wm.Name)
			assert.Equal(t, tt.want.ProjectRoleAssertion, wm.ProjectRoleAssertion)
			assert.Equal(t, tt.want.ProjectRoleCheck, wm.ProjectRoleCheck)
			assert.Equal(t, tt.want.HasProjectCheck, wm.HasProjectCheck)
			assert.Equal(t, tt.want.PrivateLabelingSetting, wm.PrivateLabelingSetting)
			assert.Equal(t, tt.want.State, wm.State)
		})
	}
}
