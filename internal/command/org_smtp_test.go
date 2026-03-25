package command

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_TestOrgSMTPConfigById(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
		alg        crypto.EncryptionAlgorithm
	}
	type args struct {
		orgID string
		id    string
		email string
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "id empty, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				orgID: "ORG1",
				email: "test@example.com",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-99oki1", "Errors.IDMissing"))
				},
			},
		},
		{
			name: "email empty, invalid argument error",
			fields: fields{
				eventstore: expectEventstore(),
			},
			args: args{
				orgID: "ORG1",
				id:    "ID",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowInvalidArgument(nil, "ORG-99yth1", "Errors.SMTPConfig.TestEmailNotFound"))
				},
			},
		},
		{
			name: "smtp config not found",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args: args{
				orgID: "ORG1",
				id:    "ID",
				email: "test@example.com",
			},
			res: res{
				err: func(err error) bool {
					return errors.Is(err, zerrors.ThrowNotFound(nil, "ORG-99klw1", "Errors.SMTPConfig.NotFound"))
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:     tt.fields.eventstore(t),
				smtpEncryption: tt.fields.alg,
			}
			err := r.TestOrgSMTPConfigById(context.Background(), tt.args.orgID, tt.args.id, tt.args.email)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func newOrgSMTPConfigChangedEvent(
	ctx context.Context,
	orgID, id, description string,
	tls bool,
	fromAddress, fromName, replyTo, host, user string,
	plainAuth *instance.PlainAuth,
	xoauth2Auth *instance.XOAuth2Auth,
) *org.OrgSMTPConfigChangedEvent {
	changes := []org.OrgSMTPConfigChanges{
		org.ChangeOrgSMTPConfigDescription(description),
		org.ChangeOrgSMTPConfigTLS(tls),
		org.ChangeOrgSMTPConfigFromAddress(fromAddress),
		org.ChangeOrgSMTPConfigFromName(fromName),
		org.ChangeOrgSMTPConfigReplyToAddress(replyTo),
		org.ChangeOrgSMTPConfigSMTPHost(host),
		org.ChangeOrgSMTPConfigSMTPUser(user),
	}
	if plainAuth != nil {
		changes = append(changes,
			org.ChangeOrgSMTPConfigSMTPPassword(plainAuth.Password),
		)
	}
	if xoauth2Auth != nil {
		changes = append(changes,
			org.ChangeOrgSMTPConfigXOAuth2TokenEndpoint(xoauth2Auth.TokenEndpoint),
			org.ChangeOrgSMTPConfigXOAuth2Scopes(xoauth2Auth.Scopes),
		)
		if xoauth2Auth.ClientCredentials != nil {
			changes = append(changes,
				org.ChangeOrgSMTPConfigXOAuth2ClientCredentialsClientId(xoauth2Auth.ClientCredentials.ClientId),
				org.ChangeOrgSMTPConfigXOAuth2ClientCredentialsClientSecret(xoauth2Auth.ClientCredentials.ClientSecret),
			)
		}
	}
	event, _ := org.NewOrgSMTPConfigChangeEvent(ctx,
		&org.NewAggregate(orgID).Aggregate,
		id,
		changes,
	)
	return event
}

func newOrgSMTPConfigHTTPChangedEvent(
	ctx context.Context,
	orgID, id, description, endpoint string,
	signingKey *crypto.CryptoValue,
) *org.OrgSMTPConfigHTTPChangedEvent {
	changes := []org.OrgSMTPConfigHTTPChanges{
		org.ChangeOrgSMTPConfigHTTPDescription(description),
		org.ChangeOrgSMTPConfigHTTPEndpoint(endpoint),
		org.ChangeOrgSMTPConfigHTTPSigningKey(signingKey),
	}
	event, _ := org.NewOrgSMTPConfigHTTPChangeEvent(ctx,
		&org.NewAggregate(orgID).Aggregate,
		id,
		changes,
	)
	return event
}
