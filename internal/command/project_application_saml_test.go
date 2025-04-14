package command

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/id"
	id_mock "github.com/zitadel/zitadel/internal/id/mock"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var testMetadata = []byte(`<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata"
                     validUntil="2022-08-26T14:08:16Z"
                     cacheDuration="PT604800S"
                     entityID="https://test.com/saml/metadata">
    <md:SPSSODescriptor AuthnRequestsSigned="false" WantAssertionsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
                                     Location="https://test.com/saml/acs"
                                     index="1" />
        
    </md:SPSSODescriptor>
</md:EntityDescriptor>
`)
var testMetadataChangedEntityID = []byte(`<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata"
                     validUntil="2022-08-26T14:08:16Z"
                     cacheDuration="PT604800S"
                     entityID="https://test2.com/saml/metadata">
    <md:SPSSODescriptor AuthnRequestsSigned="false" WantAssertionsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
                                     Location="https://test.com/saml/acs"
                                     index="1" />
        
    </md:SPSSODescriptor>
</md:EntityDescriptor>
`)

func TestCommandSide_AddSAMLApplication(t *testing.T) {
	type fields struct {
		eventstore  func(t *testing.T) *eventstore.Eventstore
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
				eventstore: expectEventstore(),
			},
			args: args{
				ctx:           authz.WithInstanceID(context.Background(), "instanceID"),
				samlApp:       &domain.SAMLApp{},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "project not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
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
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "invalid app, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(
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
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
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
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "create saml app, metadata not parsable",
			fields: fields{
				eventstore: expectEventstore(
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
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
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
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "create saml app, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingUnspecified),
						),
					),
					expectPush(
						project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
						),
						project.NewSAMLConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"https://test.com/saml/metadata",
							testMetadata,
							"",
							domain.LoginVersionUnspecified,
							"",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:     "app",
					EntityID:    "https://test.com/saml/metadata",
					Metadata:    testMetadata,
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
					Metadata:    testMetadata,
					MetadataURL: "",
					State:       domain.AppStateActive,
				},
			},
		},
		{
			name: "create saml app, loginversion, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingUnspecified),
						),
					),
					expectPush(
						project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
						),
						project.NewSAMLConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"https://test.com/saml/metadata",
							testMetadata,
							"",
							domain.LoginVersion2,
							"https://test.com/login",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1"),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:      "app",
					EntityID:     "https://test.com/saml/metadata",
					Metadata:     testMetadata,
					MetadataURL:  "",
					LoginVersion: domain.LoginVersion2,
					LoginBaseURI: "https://test.com/login",
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:        "app1",
					AppName:      "app",
					EntityID:     "https://test.com/saml/metadata",
					Metadata:     testMetadata,
					MetadataURL:  "",
					State:        domain.AppStateActive,
					LoginVersion: domain.LoginVersion2,
					LoginBaseURI: "https://test.com/login",
				},
			},
		},
		{
			name: "create saml app metadataURL, ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							project.NewProjectAddedEvent(context.Background(),
								&project.NewAggregate("project1", "org1").Aggregate,
								"project", true, true, true,
								domain.PrivateLabelingSettingUnspecified),
						),
					),
					expectPush(
						project.NewApplicationAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"app",
						),
						project.NewSAMLConfigAddedEvent(context.Background(),
							&project.NewAggregate("project1", "org1").Aggregate,
							"app1",
							"https://test.com/saml/metadata",
							testMetadata,
							"http://localhost:8080/saml/metadata",
							domain.LoginVersionUnspecified,
							"",
						),
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "app1"),
				httpClient:  newTestClient(200, testMetadata),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:     "app",
					EntityID:    "https://test.com/saml/metadata",
					Metadata:    nil,
					MetadataURL: "http://localhost:8080/saml/metadata",
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
					Metadata:    testMetadata,
					MetadataURL: "http://localhost:8080/saml/metadata",
					State:       domain.AppStateActive,
				},
			},
		},
		{
			name: "create saml app metadataURL, http error",
			fields: fields{
				eventstore: expectEventstore(
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
				ctx: authz.WithInstanceID(context.Background(), "instanceID"),
				samlApp: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID: "project1",
					},
					AppName:     "app",
					EntityID:    "https://test.com/saml/metadata",
					Metadata:    nil,
					MetadataURL: "http://localhost:8080/saml/metadata",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:  tt.fields.eventstore(t),
				idGenerator: tt.fields.idGenerator,
				httpClient:  tt.fields.httpClient,
			}
			c.setMilestonesCompletedForTest("instanceID")
			got, err := c.AddSAMLApplication(tt.args.ctx, tt.args.samlApp, tt.args.resourceOwner)
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
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
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
				eventstore: expectEventstore(),
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
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing appid, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
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
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "missing aggregateid, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
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
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "app not existing, not found error",
			fields: fields{
				eventstore: expectEventstore(
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
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error, metadataURL",
			fields: fields{
				eventstore: expectEventstore(
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
								testMetadata,
								"http://localhost:8080/saml/metadata",
								domain.LoginVersionUnspecified,
								"",
							),
						),
					),
				),
				httpClient: newTestClient(http.StatusOK, testMetadata),
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
					MetadataURL: "http://localhost:8080/saml/metadata",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "no changes, precondition error, metadata",
			fields: fields{
				eventstore: expectEventstore(
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
								testMetadata,
								"",
								domain.LoginVersionUnspecified,
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
					Metadata:    testMetadata,
					MetadataURL: "",
				},
				resourceOwner: "org1",
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "change saml app, ok, metadataURL",
			fields: fields{
				eventstore: expectEventstore(
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
								testMetadata,
								"http://localhost:8080/saml/metadata",
								domain.LoginVersionUnspecified,
								"",
							),
						),
					),
					expectPush(
						newSAMLAppChangedEventMetadataURL(context.Background(),
							"app1",
							"project1",
							"org1",
							"https://test.com/saml/metadata",
							"https://test2.com/saml/metadata",
							testMetadataChangedEntityID,
						),
					),
				),
				httpClient: newTestClient(http.StatusOK, testMetadataChangedEntityID),
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
					MetadataURL: "http://localhost:8080/saml/metadata",
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
					Metadata:    testMetadataChangedEntityID,
					MetadataURL: "http://localhost:8080/saml/metadata",
					State:       domain.AppStateActive,
				},
			},
		},
		{
			name: "change saml app, ok, metadata",
			fields: fields{
				eventstore: expectEventstore(
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
								testMetadata,
								"",
								domain.LoginVersionUnspecified,
								"",
							),
						),
					),
					expectPush(
						newSAMLAppChangedEventMetadata(context.Background(),
							"app1",
							"project1",
							"org1",
							"https://test.com/saml/metadata",
							"https://test2.com/saml/metadata",
							testMetadataChangedEntityID,
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
					AppID:       "app1",
					AppName:     "app",
					EntityID:    "https://test2.com/saml/metadata",
					Metadata:    testMetadataChangedEntityID,
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
					Metadata:    testMetadataChangedEntityID,
					MetadataURL: "",
					State:       domain.AppStateActive,
				},
			},
		}, {
			name: "change saml app, ok, loginversion",
			fields: fields{
				eventstore: expectEventstore(
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
								testMetadata,
								"",
								domain.LoginVersionUnspecified,
								"",
							),
						),
					),
					expectPush(
						newSAMLAppChangedEventLoginVersion(context.Background(),
							"app1",
							"project1",
							"org1",
							"https://test.com/saml/metadata",
							"https://test2.com/saml/metadata",
							testMetadataChangedEntityID,
							domain.LoginVersion2,
							"https://test.com/login",
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
					AppID:        "app1",
					AppName:      "app",
					EntityID:     "https://test2.com/saml/metadata",
					Metadata:     testMetadataChangedEntityID,
					MetadataURL:  "",
					LoginVersion: domain.LoginVersion2,
					LoginBaseURI: "https://test.com/login",
				},
				resourceOwner: "org1",
			},
			res: res{
				want: &domain.SAMLApp{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "project1",
						ResourceOwner: "org1",
					},
					AppID:        "app1",
					AppName:      "app",
					EntityID:     "https://test2.com/saml/metadata",
					Metadata:     testMetadataChangedEntityID,
					MetadataURL:  "",
					State:        domain.AppStateActive,
					LoginVersion: domain.LoginVersion2,
					LoginBaseURI: "https://test.com/login",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore(t),
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

func newSAMLAppChangedEventMetadata(ctx context.Context, appID, projectID, resourceOwner, oldEntityID, entityID string, metadata []byte) *project.SAMLConfigChangedEvent {
	changes := []project.SAMLConfigChanges{
		project.ChangeEntityID(entityID),
		project.ChangeMetadata(metadata),
	}
	event, _ := project.NewSAMLConfigChangedEvent(ctx,
		&project.NewAggregate(projectID, resourceOwner).Aggregate,
		appID,
		oldEntityID,
		changes,
	)
	return event
}

func newSAMLAppChangedEventMetadataURL(ctx context.Context, appID, projectID, resourceOwner, oldEntityID, entityID string, metadata []byte) *project.SAMLConfigChangedEvent {
	changes := []project.SAMLConfigChanges{
		project.ChangeEntityID(entityID),
		project.ChangeMetadata(metadata),
	}
	event, _ := project.NewSAMLConfigChangedEvent(ctx,
		&project.NewAggregate(projectID, resourceOwner).Aggregate,
		appID,
		oldEntityID,
		changes,
	)
	return event
}

func newSAMLAppChangedEventLoginVersion(ctx context.Context, appID, projectID, resourceOwner, oldEntityID, entityID string, metadata []byte, loginVersion domain.LoginVersion, loginURI string) *project.SAMLConfigChangedEvent {
	changes := []project.SAMLConfigChanges{
		project.ChangeEntityID(entityID),
		project.ChangeMetadata(metadata),
		project.ChangeSAMLLoginVersion(loginVersion),
		project.ChangeSAMLLoginBaseURI(loginURI),
	}
	event, _ := project.NewSAMLConfigChangedEvent(ctx,
		&project.NewAggregate(projectID, resourceOwner).Aggregate,
		appID,
		oldEntityID,
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
			Body:       io.NopCloser(bytes.NewBuffer(metadata)),
			Header:     make(http.Header), //must be non-nil value
		}
	})
	return &http.Client{
		Transport: fn,
	}
}
