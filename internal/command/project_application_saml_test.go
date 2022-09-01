package command

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/project"
)

func TestCommandSide_AddSAMLApplication(t *testing.T) {
	port := "8081"
	path := "/metadataadd"

	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
		httpClient  *http.Client
	}
	type args struct {
		ctx           context.Context
		samlApp       *domain.SAMLApp
		resourceOwner string
	}
	type res struct {
		want *domain.SAMLApp
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "no aggregate id, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:           context.Background(),
				samlApp:       &domain.SAMLApp{},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:   "app1",
					AppName: "app",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid app, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingUnspecified),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:   "app1",
					AppName: "",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "create saml app basic, metadata not parsable",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingUnspecified),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t),
			},
			args: args{
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:     "app",
					EntityID:    "https://test.com/saml/metadata",
					Metadata:    []byte("test metadata"),
					MetadataURL: "",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "create saml app basic, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingUnspecified),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								project.NewApplicationAddedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"app1",
									"app",
								),
							),
							eventFromEventPusher(
								project.NewSAMLConfigAddedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"app1",
									"https://test.com/saml/metadata",
									[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
									"",
								),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewAddApplicationUniqueConstraint("app", "project1")),
						uniqueConstraintsFromEventConstraint(project.NewAddSAMLConfigEntityIDUniqueConstraint("https://test.com/saml/metadata")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1"),
			},
			args: args{
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:     "app",
					EntityID:    "https://test.com/saml/metadata",
					Metadata:    []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
					MetadataURL: "",
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:       "app1",
					AppName:     "app",
					EntityID:    "https://test.com/saml/metadata",
					Metadata:    []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
					MetadataURL: "",
					State:       domain.AppStateActive,
				},
			},
		},
		{
			name: "create saml app metadataURL, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingUnspecified),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								project.NewApplicationAddedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"app1",
									"app",
								),
							),
							eventFromEventPusher(
								project.NewSAMLConfigAddedEvent(context.Background(),
									&project.NewAggregate("project1", "org1").Aggregate,
									"app1",
									"https://test.com/saml/metadata",
									[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
									"http://localhost:"+port+path,
								),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewAddApplicationUniqueConstraint("app", "project1")),
						uniqueConstraintsFromEventConstraint(project.NewAddSAMLConfigEntityIDUniqueConstraint("https://test.com/saml/metadata")),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1"),
				httpClient:  newTestClient(200, []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>")),
			},
			args: args{
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:     "app",
					EntityID:    "https://test.com/saml/metadata",
					Metadata:    nil,
					MetadataURL: "http://localhost:" + port + path,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:       "app1",
					AppName:     "app",
					EntityID:    "https://test.com/saml/metadata",
					Metadata:    []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
					MetadataURL: "http://localhost:" + port + path,
					State:       domain.AppStateActive,
				},
			},
		},
		{
			name: "create saml app metadataURL, http error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingUnspecified),
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t),
				httpClient:  newTestClient(http.StatusNotFound, nil),
			},
			args: args{
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:     "app",
					EntityID:    "https://test.com/saml/metadata",
					Metadata:    nil,
					MetadataURL: "http://localhost:" + port + path,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
				httpClient:  tt.fields.httpClient,
			}

			got, err := r.AddSAMLApplication(tt.args.ctx, tt.args.samlApp, tt.args.resourceOwner)
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

func TestCommandSide_ChangeSAMLApplication(t *testing.T) {
	port := "8082"
	path := "/metadatachange"

	type fields struct {
		eventstore *eventstore.Eventstore
		httpClient *http.Client
	}
	type args struct {
		ctx           context.Context
		samlApp       *domain.SAMLApp
		resourceOwner string
	}
	type res struct {
		want *domain.SAMLApp
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "invalid app, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID: "app1",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
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
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:    "",
					Metadata: []byte("just not empty"),
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing aggregateid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "",
					},
					AppID:    "appid",
					Metadata: []byte("just not empty"),
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsErrorInvalidArgument,
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
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppID:    "app1",
					Metadata: []byte("just not empty"),
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error, metadataURL",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
						eventFromEventPusher(
							project.NewSAMLConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"https://test.com/saml/metadata",
								[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
								"http://localhost:"+port+path,
							),
						),
					),
				),
				httpClient: newTestClient(http.StatusOK, []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>")),
			},
			args: args{
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppName:     "app",
					AppID:       "app1",
					EntityID:    "https://test.com/saml/metadata",
					Metadata:    nil,
					MetadataURL: "http://localhost:" + port + path,
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			name: "no changes, precondition error, metadata",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
						eventFromEventPusher(
							project.NewSAMLConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"https://test.com/saml/metadata",
								[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
								"",
							),
						),
					),
				),
				httpClient: nil,
			},
			args: args{
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppName:     "app",
					AppID:       "app1",
					EntityID:    "https://test.com/saml/metadata",
					Metadata:    []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
					MetadataURL: "",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: errors.IsPreconditionFailed,
			},
		},
		{
			name: "change saml app, ok, metadataURL",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
						eventFromEventPusher(
							project.NewSAMLConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"https://test.com/saml/metadata",
								[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
								"http://localhost:"+port+path,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newSAMLAppChangedEventMetadataURL(context.Background(),
									"app1",
									"project1",
									"org1",
									"https://test.com/saml/metadata"),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test.com/saml/metadata")),
						uniqueConstraintsFromEventConstraint(project.NewAddSAMLConfigEntityIDUniqueConstraint("https://test2.com/saml/metadata")),
					),
				),
				httpClient: newTestClient(http.StatusOK, []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-09-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test2.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>")),
			},
			args: args{
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:       "app1",
					AppName:     "app",
					EntityID:    "https://test2.com/saml/metadata",
					Metadata:    nil,
					MetadataURL: "http://localhost:" + port + path,
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:       "app1",
					AppName:     "app",
					EntityID:    "https://test2.com/saml/metadata",
					Metadata:    []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-09-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test2.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
					MetadataURL: "http://localhost:" + port + path,
					State:       domain.AppStateActive,
				},
			},
		},
		{
			name: "change saml app, ok, metadata",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							project.NewApplicationAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"app",
							),
						),
						eventFromEventPusher(
							project.NewSAMLConfigAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"app1",
								"https://test.com/saml/metadata",
								[]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-08-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
								"",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newSAMLAppChangedEventMetadata(context.Background(),
									"app1",
									"project1",
									"org1",
									"https://test.com/saml/metadata"),
							),
						},
						uniqueConstraintsFromEventConstraint(project.NewRemoveSAMLConfigEntityIDUniqueConstraint("https://test.com/saml/metadata")),
						uniqueConstraintsFromEventConstraint(project.NewAddSAMLConfigEntityIDUniqueConstraint("https://test2.com/saml/metadata")),
					),
				),
				httpClient: nil,
			},
			args: args{
				ctx: context.Background(),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:       "app1",
					AppName:     "app",
					EntityID:    "https://test2.com/saml/metadata",
					Metadata:    []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-09-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test2.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
					MetadataURL: "",
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:       "app1",
					AppName:     "app",
					EntityID:    "https://test2.com/saml/metadata",
					Metadata:    []byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-09-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test2.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>"),
					MetadataURL: "",
					State:       domain.AppStateActive,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
				httpClient: tt.fields.httpClient,
			}
			got, err := r.ChangeSAMLApplication(tt.args.ctx, tt.args.samlApp, tt.args.resourceOwner)
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

func newSAMLAppChangedEventMetadata(ctx context.Context, appID, projectID, resourceOwner, entityID string) *project.SAMLConfigChangedEvent {
	changes := []project.SAMLConfigChanges{
		project.ChangeEntityID("https://test2.com/saml/metadata"),
		project.ChangeMetadata([]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-09-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test2.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>")),
	}
	event, _ := project.NewSAMLConfigChangedEvent(ctx,
		&project.NewAggregate(projectID, resourceOwner).Aggregate,
		appID,
		entityID,
		changes,
	)
	return event
}

func newSAMLAppChangedEventMetadataURL(ctx context.Context, appID, projectID, resourceOwner, entityID string) *project.SAMLConfigChangedEvent {
	changes := []project.SAMLConfigChanges{
		project.ChangeEntityID("https://test2.com/saml/metadata"),
		project.ChangeMetadata([]byte("<?xml version=\"1.0\"?>\n<md:EntityDescriptor xmlns:md=\"urn:oasis:names:tc:SAML:2.0:metadata\"\n                     validUntil=\"2022-09-26T14:08:16Z\"\n                     cacheDuration=\"PT604800S\"\n                     entityID=\"https://test2.com/saml/metadata\">\n    <md:SPSSODescriptor AuthnRequestsSigned=\"false\" WantAssertionsSigned=\"false\" protocolSupportEnumeration=\"urn:oasis:names:tc:SAML:2.0:protocol\">\n        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>\n        <md:AssertionConsumerService Binding=\"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST\"\n                                     Location=\"https://test.com/saml/acs\"\n                                     index=\"1\" />\n        \n    </md:SPSSODescriptor>\n</md:EntityDescriptor>")),
	}
	event, _ := project.NewSAMLConfigChangedEvent(ctx,
		&project.NewAggregate(projectID, resourceOwner).Aggregate,
		appID,
		entityID,
		changes,
	)
	return event
}

type roundTripperFunc func(*http.Request) *http.Response

// RoundTrip implements the http.RoundTripper interface.
func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req), nil
}

// NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func newTestClient(httpStatus int, metadata []byte) *http.Client {
	fn := roundTripperFunc(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: httpStatus,
			Body:       ioutil.NopCloser(bytes.NewBuffer(metadata)),
			Header:     make(http.Header), //must be non-nil value
		}
	})
	return &http.Client{
		Transport: fn,
	}
}
