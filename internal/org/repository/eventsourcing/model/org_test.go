package model

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/org/model"
)

func TestOrgFromEvents(t *testing.T) {
	type args struct {
		event []*es_models.Event
		org   *Org
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "org from events, ok",
			args: args{
				event: []*es_models.Event{
					{AggregateID: "ID", Sequence: 1, Type: OrgAdded},
				},
				org: &Org{Name: "OrgName"},
			},
			result: &Org{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, State: int32(model.OrgStateActive), Name: "OrgName"},
		},
		{
			name: "org from events, nil org",
			args: args{
				event: []*es_models.Event{
					{AggregateID: "ID", Sequence: 1, Type: OrgAdded},
				},
				org: nil,
			},
			result: &Org{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, State: int32(model.OrgStateActive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.org != nil {
				data, _ := json.Marshal(tt.args.org)
				tt.args.event[0].Data = data
			}
			result, _ := OrgFromEvents(tt.args.org, tt.args.event...)
			if result.Name != tt.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.result.Name, result.Name)
			}
		})
	}
}

func TestAppendEvent(t *testing.T) {
	type args struct {
		event *es_models.Event
		org   *Org
	}
	tests := []struct {
		name   string
		args   args
		result *Org
	}{
		{
			name: "append added event",
			args: args{
				event: &es_models.Event{AggregateID: "ID", Sequence: 1, Type: OrgAdded},
				org:   &Org{Name: "OrgName"},
			},
			result: &Org{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, State: int32(model.OrgStateActive), Name: "OrgName"},
		},
		{
			name: "append change event",
			args: args{
				event: &es_models.Event{AggregateID: "ID", Sequence: 1, Type: OrgChanged, Data: []byte(`{"name": "OrgName}`)},
				org:   &Org{Name: "OrgNameChanged"},
			},
			result: &Org{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, State: int32(model.OrgStateActive), Name: "OrgNameChanged"},
		},
		{
			name: "append deactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "ID", Sequence: 1, Type: OrgDeactivated},
			},
			result: &Org{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, State: int32(model.OrgStateInactive)},
		},
		{
			name: "append reactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "ID", Sequence: 1, Type: OrgReactivated},
			},
			result: &Org{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, State: int32(model.OrgStateActive)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.org != nil {
				data, _ := json.Marshal(tt.args.org)
				tt.args.event.Data = data
			}
			result := &Org{}
			result.AppendEvent(tt.args.event)
			if result.State != tt.result.State {
				t.Errorf("got wrong result state: expected: %v, actual: %v ", tt.result.State, result.State)
			}
			if result.Name != tt.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.result.Name, result.Name)
			}
			if result.ObjectRoot.AggregateID != tt.result.ObjectRoot.AggregateID {
				t.Errorf("got wrong result id: expected: %v, actual: %v ", tt.result.ObjectRoot.AggregateID, result.ObjectRoot.AggregateID)
			}
		})
	}
}

func TestChanges(t *testing.T) {
	type args struct {
		existingOrg *Org
		newOrg      *Org
	}
	type res struct {
		changesLen int
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			name: "org name changes",
			args: args{
				existingOrg: &Org{Name: "Name"},
				newOrg:      &Org{Name: "NameChanged"},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "no changes",
			args: args{
				existingOrg: &Org{Name: "Name"},
				newOrg:      &Org{Name: "Name"},
			},
			res: res{
				changesLen: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existingOrg.Changes(tt.args.newOrg)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
			}
		})
	}
}
