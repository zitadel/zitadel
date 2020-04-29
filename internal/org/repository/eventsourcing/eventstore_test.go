package eventsourcing

import (
	"context"
	"reflect"
	"testing"

	"github.com/caos/zitadel/internal/api/auth"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/mock"
	"github.com/caos/zitadel/internal/eventstore/models"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/golang/mock/gomock"
)

func mockEventstore(t *testing.T) eventstore.Eventstore {
	ctrl := gomock.NewController(t)
	e := mock.NewMockEventstore(ctrl)

	return e
}

func TestOrgEventstore_OrgByID(t *testing.T) {
	type fields struct {
		Eventstore eventstore.Eventstore
	}
	type res struct {
		org   *org_model.Org
		isErr func(error) bool
	}
	type args struct {
		ctx context.Context
		org *org_model.Org
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name:   "no input org",
			fields: fields{Eventstore: mockEventstore(t)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: nil,
			},
			res: res{
				org:   nil,
				isErr: errors.IsErrorInvalidArgument,
			},
		},
		{
			name:   "no aggregate id in input org",
			fields: fields{Eventstore: mockEventstore(t)},
			args: args{
				ctx: auth.NewMockContext("user", "org"),
				org: &org_model.Org{ObjectRoot: models.ObjectRoot{Sequence: 4}},
			},
			res: res{
				org:   nil,
				isErr: errors.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			es := &OrgEventstore{
				Eventstore: tt.fields.Eventstore,
			}
			got, err := es.OrgByID(tt.args.ctx, tt.args.org)
			if tt.res.isErr == nil && err != nil {
				t.Errorf("no error expected got: %v", err)
			}
			if tt.res.isErr != nil && !tt.res.isErr(err) {
				t.Errorf("wrong error got %T: %v", err, err)
			}
			if !reflect.DeepEqual(got, tt.res.org) {
				t.Errorf("OrgEventstore.OrgByID() = %v, want %v", got, tt.res.org)
			}
		})
	}
}
