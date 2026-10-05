package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommands_AddClientIDMetadataDocumentAllowedURL(t *testing.T) {
	type want struct {
		details *domain.ObjectDetails
		err     error
	}
	tests := []struct {
		name       string
		eventstore func(*testing.T) *eventstore.Eventstore
		url        string
		want       want
	}{
		{
			name:       "invalid url, error",
			eventstore: expectEventstore(),
			url:        "https://app.example.com/oauth/../client",
			want: want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-p4Lz8", "Errors.Instance.SecurityPolicy.ClientIDMetadataDocumentAllowedURL.Invalid"),
			},
		},
		{
			name:       "url without path, error",
			eventstore: expectEventstore(),
			url:        "https://app.example.com",
			want: want{
				err: zerrors.ThrowInvalidArgument(nil, "COMMA-p4Lz8", "Errors.Instance.SecurityPolicy.ClientIDMetadataDocumentAllowedURL.Invalid"),
			},
		},
		{
			name: "url already allowed by setting the policy, error",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(securityPolicySetEvent(t,
						instance.ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLs([]string{"https://app.example.com/client"}),
					)),
				),
			),
			url: "https://app.example.com/client",
			want: want{
				err: zerrors.ThrowAlreadyExists(nil, "COMMA-Vq7tY", "Errors.Instance.SecurityPolicy.ClientIDMetadataDocumentAllowedURL.AlreadyExists"),
			},
		},
		{
			name: "url already added, error",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent(context.Background(),
						&instance.NewAggregate("instanceID").Aggregate, "https://app.example.com/client")),
				),
			),
			url: "https://app.example.com/client",
			want: want{
				err: zerrors.ThrowAlreadyExists(nil, "COMMA-Vq7tY", "Errors.Instance.SecurityPolicy.ClientIDMetadataDocumentAllowedURL.AlreadyExists"),
			},
		},
		{
			name: "url added again after removal, ok",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent(context.Background(),
						&instance.NewAggregate("instanceID").Aggregate, "https://app.example.com/client")),
					eventFromEventPusher(instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent(context.Background(),
						&instance.NewAggregate("instanceID").Aggregate, "https://app.example.com/client")),
				),
				expectPush(
					instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent(context.Background(),
						&instance.NewAggregate("instanceID").Aggregate, "https://app.example.com/client"),
				),
			),
			url: "https://app.example.com/client",
			want: want{
				details: &domain.ObjectDetails{ResourceOwner: "instanceID"},
			},
		},
		{
			name: "prefix url added next to other urls, ok",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(securityPolicySetEvent(t,
						instance.ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLs([]string{"https://app.example.com/client"}),
					)),
				),
				expectPush(
					instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent(context.Background(),
						&instance.NewAggregate("instanceID").Aggregate, "https://clients.example.com/connectors/"),
				),
			),
			url: "https://clients.example.com/connectors/",
			want: want{
				details: &domain.ObjectDetails{ResourceOwner: "instanceID"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.eventstore(t),
			}
			got, err := c.AddClientIDMetadataDocumentAllowedURL(authz.WithInstanceID(context.Background(), "instanceID"), tt.url)
			assert.ErrorIs(t, err, tt.want.err)
			assertObjectDetails(t, tt.want.details, got)
		})
	}
}

func TestCommands_RemoveClientIDMetadataDocumentAllowedURL(t *testing.T) {
	tests := []struct {
		name        string
		eventstore  func(*testing.T) *eventstore.Eventstore
		url         string
		wantDetails *domain.ObjectDetails
	}{
		{
			name:       "url not allowed, no error and no event",
			eventstore: expectEventstore(expectFilter()),
			url:        "https://app.example.com/client",
		},
		{
			name: "url removed before, no error and no event",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent(context.Background(),
						&instance.NewAggregate("instanceID").Aggregate, "https://app.example.com/client")),
					eventFromEventPusher(instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent(context.Background(),
						&instance.NewAggregate("instanceID").Aggregate, "https://app.example.com/client")),
				),
			),
			url: "https://app.example.com/client",
		},
		{
			name: "url allowed by setting the policy, removed",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(securityPolicySetEvent(t,
						instance.ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLs([]string{"https://app.example.com/client", "https://clients.example.com/"}),
					)),
				),
				expectPush(
					instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent(context.Background(),
						&instance.NewAggregate("instanceID").Aggregate, "https://app.example.com/client"),
				),
			),
			url:         "https://app.example.com/client",
			wantDetails: &domain.ObjectDetails{ResourceOwner: "instanceID"},
		},
		{
			name: "added url, removed",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent(context.Background(),
						&instance.NewAggregate("instanceID").Aggregate, "https://app.example.com/client")),
				),
				expectPush(
					instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLRemovedEvent(context.Background(),
						&instance.NewAggregate("instanceID").Aggregate, "https://app.example.com/client"),
				),
			),
			url:         "https://app.example.com/client",
			wantDetails: &domain.ObjectDetails{ResourceOwner: "instanceID"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.eventstore(t),
			}
			got, err := c.RemoveClientIDMetadataDocumentAllowedURL(authz.WithInstanceID(context.Background(), "instanceID"), tt.url)
			assert.NoError(t, err)
			if tt.wantDetails == nil {
				assert.Nil(t, got)
				return
			}
			assertObjectDetails(t, tt.wantDetails, got)
		})
	}
}

func TestCommands_SetSecurityPolicy_clientIDMetadataDocument(t *testing.T) {
	t.Run("invalid allowed url, error", func(t *testing.T) {
		c := &Commands{
			eventstore: expectEventstore()(t),
		}
		_, err := c.SetSecurityPolicy(authz.WithInstanceID(context.Background(), "instanceID"), &SecurityPolicy{
			EnableClientIDMetadataDocument:      true,
			ClientIDMetadataDocumentAllowedURLs: []string{"https://app.example.com/client", "https://app.example.com/client#fragment"},
		})
		assert.ErrorIs(t, err, zerrors.ThrowInvalidArgument(nil, "COMMA-p4Lz8", "Errors.Instance.SecurityPolicy.ClientIDMetadataDocumentAllowedURL.Invalid"))
	})
	t.Run("allowed urls and allow any url, set", func(t *testing.T) {
		c := &Commands{
			eventstore: expectEventstore(
				expectFilter(),
				expectPush(securityPolicySetEvent(t,
					instance.ChangeSecurityPolicyEnableClientIDMetadataDocument(true),
					instance.ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLs([]string{"https://clients.example.com/"}),
					instance.ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLConstraints(nil, []string{"https://clients.example.com/"}),
					instance.ChangeSecurityPolicyClientIDMetadataDocumentAllowAnyURL(true),
				)),
			)(t),
		}
		_, err := c.SetSecurityPolicy(authz.WithInstanceID(context.Background(), "instanceID"), &SecurityPolicy{
			EnableClientIDMetadataDocument:      true,
			ClientIDMetadataDocumentAllowedURLs: []string{"https://clients.example.com/"},
			ClientIDMetadataDocumentAllowAnyURL: true,
		})
		assert.NoError(t, err)
	})
	t.Run("set list replaces added urls", func(t *testing.T) {
		c := &Commands{
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent(context.Background(),
						&instance.NewAggregate("instanceID").Aggregate, "https://app.example.com/client")),
				),
				expectPush(securityPolicySetEvent(t,
					instance.ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLs(nil),
					instance.ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLConstraints([]string{"https://app.example.com/client"}, nil),
				)),
			)(t),
		}
		_, err := c.SetSecurityPolicy(authz.WithInstanceID(context.Background(), "instanceID"), &SecurityPolicy{})
		assert.NoError(t, err)
	})
}

func securityPolicySetEvent(t *testing.T, changes ...instance.SecurityPolicyChanges) *instance.SecurityPolicySetEvent {
	t.Helper()
	event, err := instance.NewSecurityPolicySetEvent(context.Background(), &instance.NewAggregate("instanceID").Aggregate, changes)
	if err != nil {
		t.Fatal(err)
	}
	return event
}

// TestCommands_SetLegacySecurityPolicy pins that the admin v1 and settings v2beta APIs only
// change the settings they know, so a call through them does not reset dynamic client
// registration or client ID metadata documents.
func TestCommands_SetLegacySecurityPolicy(t *testing.T) {
	existing := func(t *testing.T) eventstore.Event {
		return eventFromEventPusher(securityPolicySetEvent(t,
			instance.ChangeSecurityPolicyEnableIframeEmbedding(false),
			instance.ChangeSecurityPolicyAllowedOrigins([]string{"https://origin.example.com"}),
			instance.ChangeSecurityPolicyEnableDynamicClientRegistration(true),
			instance.ChangeSecurityPolicyAllowUnauthenticatedDynamicClientRegistration(true),
			instance.ChangeSecurityPolicyEnableClientIDMetadataDocument(true),
			instance.ChangeSecurityPolicyClientIDMetadataDocumentAllowedURLs([]string{"https://app.example.com/client"}),
			instance.ChangeSecurityPolicyClientIDMetadataDocumentAllowAnyURL(true),
		))
	}

	t.Run("only the legacy settings change", func(t *testing.T) {
		c := &Commands{
			eventstore: expectEventstore(
				expectFilter(
					existing(t),
					eventFromEventPusher(instance.NewSecurityPolicyClientIDMetadataDocumentAllowedURLAddedEvent(context.Background(),
						&instance.NewAggregate("instanceID").Aggregate, "https://clients.example.com/")),
				),
				expectPush(securityPolicySetEvent(t,
					instance.ChangeSecurityPolicyEnableIframeEmbedding(true),
					instance.ChangeSecurityPolicyEnableImpersonation(true),
				)),
			)(t),
		}
		_, err := c.SetLegacySecurityPolicy(authz.WithInstanceID(context.Background(), "instanceID"), &SecurityPolicy{
			EnableIframeEmbedding: true,
			AllowedOrigins:        []string{"https://origin.example.com"},
			EnableImpersonation:   true,
		})
		assert.NoError(t, err)
	})

	t.Run("settings the legacy apis do not know are ignored", func(t *testing.T) {
		c := &Commands{
			eventstore: expectEventstore(
				expectFilter(existing(t)),
			)(t),
		}
		_, err := c.SetLegacySecurityPolicy(authz.WithInstanceID(context.Background(), "instanceID"), &SecurityPolicy{
			AllowedOrigins: []string{"https://origin.example.com"},
		})
		assert.ErrorIs(t, err, zerrors.ThrowPreconditionFailed(nil, "POLICY-EWsf3", "Errors.NoChangesFound"))
	})
}
