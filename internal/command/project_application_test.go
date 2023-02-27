package command

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/static"
	"github.com/zitadel/zitadel/internal/static/mock"
)

func TestCommandSide_ChangeApplication(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		projectID     string
		app           *domain.ChangeApp
		resourceOwner string
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
			name: "invalid app missing projectid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "",
				app: &domain.ChangeApp{
					AppID:   "app1",
					AppName: "app",
					ExternalURL: "external-url",
					IsVisibleToEndUser: false,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid app missing appid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				app: &domain.ChangeApp{
					AppName: "app",
					ExternalURL: "external-url",
					IsVisibleToEndUser: false,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid app missing name, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				app: &domain.ChangeApp{
					AppID:   "app1",
					AppName: "",
					ExternalURL: "external-url",
					IsVisibleToEndUser: false,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "invalid app external url, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				app: &domain.ChangeApp{
					AppID:   "app1",
					AppName: "app",
					ExternalURL: "external-url",
					IsVisibleToEndUser: false,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "app not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				app: &domain.ChangeApp{
					AppID:   "app1",
					AppName: "app",
					ExternalURL: "https://zitadel.com",
					IsVisibleToEndUser: true,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "app name not changed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
							"https://zitadel.com",
							false,
						)),
					),
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				app: &domain.ChangeApp{
					AppID:   "app1",
					AppName: "app",
					ExternalURL: "https://zitadel.com",
					IsVisibleToEndUser: false,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},

		{
			name: "app changed, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
							"https://zitadel.com",
							true,
						)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewApplicationChangedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
								"app changed",
								"https://zitadel.com",
								true,
							)),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveApplicationUniqueConstraint("app", "project1")),
						uniqueConstraintsFromEventConstraint(project.NewAddApplicationUniqueConstraint("app changed", "project1")),
					),
				),
			},
			args: args{
				ctx:       context.Background(),
				projectID: "project1",
				app: &domain.ChangeApp{
					AppID:   "app1",
					AppName: "app changed",
					ExternalURL: "https://zitadel.com",
					IsVisibleToEndUser: true,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeApplication(tt.args.ctx, tt.args.projectID, tt.args.app, tt.args.resourceOwner)
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

func TestCommandSide_DeactivateApplication(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		projectID     string
		appID         string
		resourceOwner string
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
			name: "missing projectid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing appid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "app not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "app already inactive, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
							"",
							false,
						)),
						eventFromEventPusher(project.NewApplicationDeactivatedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
						)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app deactivate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
							"",
							false,
						)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewApplicationDeactivatedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
							)),
						},
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.DeactivateApplication(tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.resourceOwner)
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

func TestCommandSide_ReactivateApplication(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		projectID     string
		appID         string
		resourceOwner string
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
			name: "missing projectid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing appid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "app not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "app already active, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
							"",
							false,
						)),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "app reactivate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
							"",
							false,
						)),
						eventFromEventPusher(project.NewApplicationDeactivatedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
						)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewApplicationReactivatedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
							)),
						},
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ReactivateApplication(tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.resourceOwner)
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

func TestCommandSide_RemoveApplication(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		projectID     string
		appID         string
		resourceOwner string
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
			name: "missing projectid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing appid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "app not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "app remove, entityID, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
							"",
							false,
						)),
					),
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
							"",
							false,
						)),
						eventFromEventPusher(project.NewSAMLConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"https://test.com/saml/metadata",
							[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
							"",
						)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewApplicationRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
								"https://test.com/saml/metadata",
							)),
						}, /**/
						uniqueConstraintsFromEventConstraint(project.NewRemoveApplicationUniqueConstraint("app", "project1")),
						uniqueConstraintsFromEventConstraint(project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test.com/saml/metadata")),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "app remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
							"",
							false,
						)),
					),
					// app is not saml, or no saml config available
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(project.NewApplicationRemovedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
								"",
							)),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveApplicationUniqueConstraint("app", "project1")),
					),
				),
			},
			args: args{
				ctx:           context.Background(),
				projectID:     "project1",
				appID:         "app1",
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveApplication(tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.resourceOwner)
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

func TestCommandSide_AddApplicationIcon(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
		storage    static.Storage
	}
	type args struct {
		ctx    context.Context
		projectID string
		appID string
		dark bool // Using dark theme?
		upload *AssetUpload
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
			name: "projectID empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:   context.Background(),
				projectID: "",
				upload: &AssetUpload{
					ResourceOwner: "org1",
					ObjectName:    "logo",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "appID empty, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:   context.Background(),
				projectID: "project1",
				appID: "",
				upload: &AssetUpload{
					ResourceOwner: "org1",
					ObjectName:    "logo",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "application not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				projectID: "project1",
				appID: "app",
				upload: &AssetUpload{
					ResourceOwner: "org1",
					ObjectName:    "logo",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "upload failed, error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
							"https://zitadel.com",
							false,
						)),
					),
				),
				storage: mock.NewStorage(t).ExpectPutObjectError(),
			},
			args: args{
				ctx:   context.Background(),
				projectID: "project1",
				appID: "app1",
				upload: &AssetUpload{
					ResourceOwner: "org1",
					ObjectName:    "logo",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				err: caos_errs.IsInternal,
			},
		},
		{
			name: "light icon added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
							"https://zitadel.com",
							false,
						)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								project.NewApplicationLightIconAddedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"logo",
								),
							),
						},
					),
				),
				storage: mock.NewStorage(t).ExpectPutObject(),
			},
			args: args{
				ctx:   context.Background(),
				projectID: "project1",
				appID: "app1",
				dark: false,
				upload: &AssetUpload{
					ResourceOwner: "org1",
					ObjectName:    "logo",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			name: "dark icon added, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
							"https://zitadel.com",
							false,
						)),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								project.NewApplicationDarkIconAddedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"logo",
								),
							),
						},
					),
				),
				storage: mock.NewStorage(t).ExpectPutObject(),
			},
			args: args{
				ctx:   context.Background(),
				projectID: "project1",
				appID: "app1",
				dark: true,
				upload: &AssetUpload{
					ResourceOwner: "org1",
					ObjectName:    "logo",
					ContentType:   "image",
					ObjectType:    static.ObjectTypeStyling,
					File:          bytes.NewReader([]byte("test")),
					Size:          4,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
				static:     tt.fields.storage,
			}
			got, err := r.AddApplicationIcon(tt.args.ctx, tt.args.projectID, tt.args.appID, tt.args.dark, tt.args.upload)
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