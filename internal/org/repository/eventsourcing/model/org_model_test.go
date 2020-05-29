package model

import (
	"encoding/json"
	"testing"

	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/org/model"
)

func TestOrgFromEvents(t *testing.T) {
	type args struct {
		event []*es_models.Event
		org   *model2.Org
	}
	tests := []struct {
		name   string
		args   args
		result *model2.Org
	}{
		{
			name: "org from events, ok",
			args: args{
				event: []*es_models.Event{
					{AggregateID: "ID", Sequence: 1, Type: OrgAdded},
				},
				org: &model2.Org{Name: "OrgName"},
			},
			result: &model2.Org{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, State: int32(model.ORGSTATE_ACTIVE), Name: "OrgName"},
		},
		{
			name: "org from events, nil org",
			args: args{
				event: []*es_models.Event{
					{AggregateID: "ID", Sequence: 1, Type: OrgAdded},
				},
				org: nil,
			},
			result: &model2.Org{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, State: int32(model.ORGSTATE_ACTIVE)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.org != nil {
				data, _ := json.Marshal(tt.args.org)
				tt.args.event[0].Data = data
			}
			result, _ := model2.OrgFromEvents(tt.args.org, tt.args.event...)
			if result.Name != tt.result.Name {
				t.Errorf("got wrong result name: expected: %v, actual: %v ", tt.result.Name, result.Name)
			}
		})
	}
}

func TestAppendEvent(t *testing.T) {
	type args struct {
		event *es_models.Event
		org   *model2.Org
	}
	tests := []struct {
		name   string
		args   args
		result *model2.Org
	}{
		{
			name: "append added event",
			args: args{
				event: &es_models.Event{AggregateID: "ID", Sequence: 1, Type: OrgAdded},
				org:   &model2.Org{Name: "OrgName"},
			},
			result: &model2.Org{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, State: int32(model.ORGSTATE_ACTIVE), Name: "OrgName"},
		},
		{
			name: "append change event",
			args: args{
				event: &es_models.Event{AggregateID: "ID", Sequence: 1, Type: OrgChanged, Data: []byte(`{"domain": "OrgDomain"}`)},
				org:   &model2.Org{Name: "OrgName", Domain: "asdf"},
			},
			result: &model2.Org{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, State: int32(model.ORGSTATE_ACTIVE), Name: "OrgName", Domain: "OrgDomain"},
		},
		{
			name: "append deactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "ID", Sequence: 1, Type: OrgDeactivated},
			},
			result: &model2.Org{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, State: int32(model.ORGSTATE_INACTIVE)},
		},
		{
			name: "append reactivate event",
			args: args{
				event: &es_models.Event{AggregateID: "ID", Sequence: 1, Type: OrgReactivated},
			},
			result: &model2.Org{ObjectRoot: es_models.ObjectRoot{AggregateID: "ID"}, State: int32(model.ORGSTATE_ACTIVE)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.org != nil {
				data, _ := json.Marshal(tt.args.org)
				tt.args.event.Data = data
			}
			result := &model2.Org{}
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
		existing *model2.Org
		new      *model2.Org
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
				existing: &model2.Org{Name: "Name"},
				new:      &model2.Org{Name: "NameChanged"},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "org domain changes",
			args: args{
				existing: &model2.Org{Name: "Name", Domain: "old domain"},
				new:      &model2.Org{Name: "Name", Domain: "new domain"},
			},
			res: res{
				changesLen: 1,
			},
		},
		{
			name: "no changes",
			args: args{
				existing: &model2.Org{Name: "Name"},
				new:      &model2.Org{Name: "Name"},
			},
			res: res{
				changesLen: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes := tt.args.existing.Changes(tt.args.new)
			if len(changes) != tt.res.changesLen {
				t.Errorf("got wrong changes len: expected: %v, actual: %v ", tt.res.changesLen, len(changes))
			}
		})
	}
}
