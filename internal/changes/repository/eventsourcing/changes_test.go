package eventsourcing

import (
	"encoding/json"
	"testing"

	"github.com/caos/zitadel/internal/changes/model"
	chg_model "github.com/caos/zitadel/internal/changes/model"
	chg_type "github.com/caos/zitadel/internal/changes/types"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_model "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/golang/mock/gomock"
)

func TestChangesUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es            *ChangesEventstore
		aggregateType es_model.AggregateType
		id            string
		secId         string
		lastSequence  uint64
		limit         uint64
	}
	type res struct {
		changes *chg_model.Changes
		user    *chg_type.User
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "changes from events, ok",
			args: args{
				es:            GetMockChangesUserOK(ctrl),
				aggregateType: chg_model.User,
				id:            "1",
				secId:         "",
				lastSequence:  0,
				limit:         0,
			},
			res: res{
				changes: &model.Changes{Changes: []*model.Change{&chg_model.Change{EventType: "", Sequence: 1, Modifier: ""}}, LastSequence: 1},
				user:    &chg_type.User{FirstName: "Hans", LastName: "Muster", EMailAddress: "a@b.ch", Phone: "+41 12 345 67 89", Language: "D", UserName: "HansMuster"},
			},
		},
		{
			name: "changes from events, no events",
			args: args{
				es:            GetMockChangesUserNoEvents(ctrl),
				aggregateType: chg_model.User,
				id:            "2",
				secId:         "",
				lastSequence:  0,
				limit:         0,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.Changes(nil, tt.args.aggregateType, tt.args.id, tt.args.secId, tt.args.lastSequence, tt.args.limit)

			user := &chg_type.User{}
			if result != nil && len(result.Changes) > 0 {
				b, err := json.Marshal(result.Changes[0].Data)
				json.Unmarshal(b, user)
				if err != nil {
				}
			}
			if !tt.res.wantErr && result.LastSequence != tt.res.changes.LastSequence && user.UserName != tt.res.user.UserName {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.changes.LastSequence, result.LastSequence)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangesProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es            *ChangesEventstore
		aggregateType es_model.AggregateType
		id            string
		secId         string
		lastSequence  uint64
		limit         uint64
	}
	type res struct {
		changes *chg_model.Changes
		project *chg_type.Project
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "changes from events, ok",
			args: args{
				es:            GetMockChangesProjectOK(ctrl),
				aggregateType: chg_model.Project,
				id:            "1",
				secId:         "",
				lastSequence:  0,
				limit:         0,
			},
			res: res{
				changes: &model.Changes{Changes: []*model.Change{&chg_model.Change{EventType: "", Sequence: 1, Modifier: ""}}, LastSequence: 1},
				project: &chg_type.Project{Name: "MusterProject"},
			},
		},
		{
			name: "changes from events, no events",
			args: args{
				es:            GetMockChangesProjectNoEvents(ctrl),
				aggregateType: chg_model.Project,
				id:            "2",
				secId:         "",
				lastSequence:  0,
				limit:         0,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.Changes(nil, tt.args.aggregateType, tt.args.id, tt.args.secId, tt.args.lastSequence, tt.args.limit)

			project := &chg_type.Project{}
			if result != nil && len(result.Changes) > 0 {
				b, err := json.Marshal(result.Changes[0].Data)
				json.Unmarshal(b, project)
				if err != nil {
				}
			}
			if !tt.res.wantErr && result.LastSequence != tt.res.changes.LastSequence && project.Name != tt.res.project.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.changes.LastSequence, result.LastSequence)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangesApplication(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es            *ChangesEventstore
		aggregateType es_model.AggregateType
		id            string
		secId         string
		lastSequence  uint64
		limit         uint64
	}
	type res struct {
		changes *chg_model.Changes
		app     *chg_type.App
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "changes from events, ok",
			args: args{
				es:            GetMockChangesApplicationOK(ctrl),
				aggregateType: chg_model.Application,
				id:            "1",
				secId:         "",
				lastSequence:  0,
				limit:         0,
			},
			res: res{
				changes: &model.Changes{Changes: []*model.Change{&chg_model.Change{EventType: "", Sequence: 1, Modifier: ""}}, LastSequence: 1},
				app:     &chg_type.App{Name: "MusterApp", AppId: "AppId", AppType: 3, ClientId: "MyClient"},
			},
		},
		{
			name: "changes from events, no events",
			args: args{
				es:            GetMockChangesApplicationNoEvents(ctrl),
				aggregateType: chg_model.User,
				id:            "2",
				secId:         "",
				lastSequence:  0,
				limit:         0,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.Changes(nil, tt.args.aggregateType, tt.args.id, tt.args.secId, tt.args.lastSequence, tt.args.limit)

			app := &chg_type.App{}
			if result != nil && len(result.Changes) > 0 {
				b, err := json.Marshal(result.Changes[0].Data)
				json.Unmarshal(b, app)
				if err != nil {
				}
			}
			if !tt.res.wantErr && result.LastSequence != tt.res.changes.LastSequence && app.Name != tt.res.app.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.changes.LastSequence, result.LastSequence)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func TestChangesOrg(t *testing.T) {
	ctrl := gomock.NewController(t)
	type args struct {
		es            *ChangesEventstore
		aggregateType es_model.AggregateType
		id            string
		secId         string
		lastSequence  uint64
		limit         uint64
	}
	type res struct {
		changes *chg_model.Changes
		org     *chg_type.Org
		wantErr bool
		errFunc func(err error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "changes from events, ok",
			args: args{
				es:            GetMockChangesOrgOK(ctrl),
				aggregateType: chg_model.Org,
				id:            "1",
				secId:         "",
				lastSequence:  0,
				limit:         0,
			},
			res: res{
				changes: &model.Changes{Changes: []*model.Change{&chg_model.Change{EventType: "", Sequence: 1, Modifier: ""}}, LastSequence: 1},
				org:     &chg_type.Org{Name: "MusterOrg", Domain: "myDomain", UserId: "myUserId"},
			},
		},
		{
			name: "changes from events, no events",
			args: args{
				es:            GetMockChangesOrgNoEvents(ctrl),
				aggregateType: chg_model.User,
				id:            "2",
				secId:         "",
				lastSequence:  0,
				limit:         0,
			},
			res: res{
				wantErr: true,
				errFunc: caos_errs.IsNotFound,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.args.es.Changes(nil, tt.args.aggregateType, tt.args.id, tt.args.secId, tt.args.lastSequence, tt.args.limit)

			org := &chg_type.Org{}
			if result != nil && len(result.Changes) > 0 {
				b, err := json.Marshal(result.Changes[0].Data)
				json.Unmarshal(b, org)
				if err != nil {
				}
			}
			if !tt.res.wantErr && result.LastSequence != tt.res.changes.LastSequence && org.Name != tt.res.org.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.res.changes.LastSequence, result.LastSequence)
			}
			if tt.res.wantErr && !tt.res.errFunc(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}
