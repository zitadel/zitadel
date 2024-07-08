package command

import (
	"context"
	"testing"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	id_mock "github.com/zitadel/zitadel/internal/id_generator/mock"
	"github.com/zitadel/zitadel/internal/repository/restrictions"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestSetRestrictions(t *testing.T) {
	type fields func(*testing.T) (*eventstore.Eventstore, id_generator.Generator)
	type args struct {
		setRestrictions *SetRestrictions
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "set new restrictions",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								restrictions.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&restrictions.NewAggregate("restrictions1", "INSTANCE", "INSTANCE").Aggregate,
										restrictions.SetEventType,
									),
									restrictions.ChangeDisallowPublicOrgRegistration(true),
								),
							),
						),
					),
					id_mock.NewIDGeneratorExpectIDs(t, "restrictions1")
			},
			args: args{
				setRestrictions: &SetRestrictions{
					DisallowPublicOrgRegistration: gu.Ptr(true),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "change restrictions",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								restrictions.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&restrictions.NewAggregate("restrictions1", "INSTANCE", "INSTANCE").Aggregate,
										restrictions.SetEventType,
									),
									restrictions.ChangeDisallowPublicOrgRegistration(true),
								),
							),
						),
						expectPush(
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								restrictions.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&restrictions.NewAggregate("restrictions1", "INSTANCE", "INSTANCE").Aggregate,
										restrictions.SetEventType,
									),
									restrictions.ChangeDisallowPublicOrgRegistration(false),
								),
							),
						),
					),
					nil
			},
			args: args{
				setRestrictions: &SetRestrictions{
					DisallowPublicOrgRegistration: gu.Ptr(false),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "set restrictions idempotency",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(
						t,
						expectFilter(
							eventFromEventPusher(
								restrictions.NewSetEvent(
									eventstore.NewBaseEventForPush(
										context.Background(),
										&restrictions.NewAggregate("restrictions1", "INSTANCE", "INSTANCE").Aggregate,
										restrictions.SetEventType,
									),
									restrictions.ChangeDisallowPublicOrgRegistration(true),
								),
							),
						),
					),
					nil
			},
			args: args{
				setRestrictions: &SetRestrictions{
					DisallowPublicOrgRegistration: gu.Ptr(true),
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "no restrictions defined",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							restrictions.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&restrictions.NewAggregate("restrictions1", "INSTANCE", "INSTANCE").Aggregate,
									restrictions.SetEventType,
								),
								restrictions.ChangeDisallowPublicOrgRegistration(true),
							),
						),
					),
				), nil
			},
			args: args{
				setRestrictions: &SetRestrictions{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "unsupported language restricted",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							restrictions.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&restrictions.NewAggregate("restrictions1", "INSTANCE", "INSTANCE").Aggregate,
									restrictions.SetEventType,
								),
								restrictions.ChangeAllowedLanguages(SupportedLanguages),
							),
						),
					),
				), nil
			},
			args: args{
				setRestrictions: &SetRestrictions{
					AllowedLanguages: []language.Tag{AllowedLanguage, UnsupportedLanguage},
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "default language not allowed",
			fields: func(*testing.T) (*eventstore.Eventstore, id_generator.Generator) {
				return eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							restrictions.NewSetEvent(
								eventstore.NewBaseEventForPush(
									context.Background(),
									&restrictions.NewAggregate("restrictions1", "INSTANCE", "INSTANCE").Aggregate,
									restrictions.SetEventType,
								),
								restrictions.ChangeAllowedLanguages(OnlyAllowedLanguages),
							),
						),
					),
				), nil
			},
			args: args{
				setRestrictions: &SetRestrictions{
					AllowedLanguages: []language.Tag{DisallowedLanguage},
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_eventstore, _id_generator := tt.fields(t)
			r := &Commands{eventstore: _eventstore}
			id_generator.SetGenerator(_id_generator)
			got, err := r.SetInstanceRestrictions(authz.WithInstance(context.Background(), &mockInstance{}), tt.args.setRestrictions)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}
