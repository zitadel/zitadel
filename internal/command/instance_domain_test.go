package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddInstanceDomain(t *testing.T) {
	type fields struct {
		eventstore     *eventstore.Eventstore
		externalSecure bool
	}
	type args struct {
		ctx    context.Context
		domain string
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
			name: "invalid domain, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				domain: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid domain ', error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				domain: "hodor's-org.localhost",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid domain umlaut, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				domain: "bücher.ch",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid domain other unicode, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				domain: "🦒.ch",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "domain already exists, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"domain.ch",
								false,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				domain: "domain.ch",
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "domain add, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"consoleApplicationID",
								"app",
							),
						),
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							project.NewOIDCConfigAddedEvent(context.Background(),
								&project.NewAggregate("projectID", "org1").Aggregate,
								domain.OIDCVersionV1,
								"consoleApplicationID",
								"client1@project",
								"secret",
								[]string{"https://test.ch"},
								[]domain.OIDCResponseType{domain.OIDCResponseTypeCode},
								[]domain.OIDCGrantType{domain.OIDCGrantTypeAuthorizationCode},
								domain.OIDCApplicationTypeWeb,
								domain.OIDCAuthMethodTypePost,
								[]string{"https://test.ch/logout"},
								true,
								domain.OIDCTokenTypeBearer,
								true,
								true,
								true,
								time.Second*1,
								[]string{"https://sub.test.ch"},
								false,
								"",
							),
						),
					),
					expectPush(
						instance.NewDomainAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"domain.ch",
							false,
						),
						newOIDCAppChangedEventInstanceDomain(
							context.Background(),
							"consoleApplicationID",
							"projectID",
							"org1",
						),
					),
				),
				externalSecure: true,
			},
			args: args{
				ctx:    authz.WithInstance(context.Background(), new(mockInstance)),
				domain: "domain.ch",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:     tt.fields.eventstore,
				externalSecure: tt.fields.externalSecure,
			}
			got, err := r.AddInstanceDomain(tt.args.ctx, tt.args.domain)
			if tt.res.err == nil {
				assert.NoError(t, err)
			} else if !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
				return
			}
			assertObjectDetails(t, tt.res.want, got)
		})
	}
}

func TestCommandSide_SetPrimaryInstanceDomain(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		domain string
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
			name: "invalid domain, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				domain: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "domain not exists, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				domain: "domain.ch",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "set primary domain, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"domain.ch",
								false,
							),
						),
					),
					expectPush(
						instance.NewDomainPrimarySetEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"domain.ch",
						),
					),
				),
			},
			args: args{
				ctx:    authz.WithInstanceID(context.Background(), "INSTANCE"),
				domain: "domain.ch",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.SetPrimaryInstanceDomain(tt.args.ctx, tt.args.domain)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RemoveInstanceDomain(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		domain string
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
			name: "invalid domain, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				domain: "",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "domain not exists, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:    context.Background(),
				domain: "domain.ch",
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove domain, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"domain.ch",
								false,
							),
						),
					),
					expectPush(
						instance.NewDomainRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"domain.ch",
						),
					),
				),
			},
			args: args{
				ctx:    authz.WithInstanceID(context.Background(), "INSTANCE"),
				domain: "domain.ch",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "remove generated domain, precondition failed",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"domain.ch",
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:    context.Background(),
				domain: "domain.ch",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveInstanceDomain(tt.args.ctx, tt.args.domain)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func newOIDCAppChangedEventInstanceDomain(ctx context.Context, appID, projectID, resourceOwner string) *project.OIDCConfigChangedEvent {
	changes := []project.OIDCConfigChanges{
		project.ChangeRedirectURIs([]string{"https://test.ch", "https://domain.ch/ui/console/auth/callback"}),
		project.ChangePostLogoutRedirectURIs([]string{"https://test.ch/logout", "https://domain.ch/ui/console/signedout"}),
	}

	aggregate := project.NewAggregate(projectID, resourceOwner).Aggregate
	aggregate.InstanceID = "INSTANCE"

	event, _ := project.NewOIDCConfigChangedEvent(ctx,
		&aggregate,
		appID,
		changes,
	)
	return event
}
