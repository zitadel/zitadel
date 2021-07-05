package zitadel

import (
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
)

func getSecretsMap(desiredKind *DesiredV0) (
	map[string]*secret.Secret,
	map[string]*secret.Existing,
) {

	var (
		secrets  = map[string]*secret.Secret{}
		existing = map[string]*secret.Existing{}
	)

	if desiredKind.Spec == nil {
		desiredKind.Spec = &Spec{}
	}

	if desiredKind.Spec.Configuration == nil {
		desiredKind.Spec.Configuration = &configuration.Configuration{}
	}

	conf := desiredKind.Spec.Configuration

	if conf.Tracing == nil {
		conf.Tracing = &configuration.Tracing{}
	}
	if conf.Tracing.ServiceAccountJSON == nil {
		conf.Tracing.ServiceAccountJSON = &secret.Secret{}
	}
	if conf.Tracing.ExistingServiceAccountJSON == nil {
		conf.Tracing.ExistingServiceAccountJSON = &secret.Existing{}
	}
	sakey := "tracingserviceaccountjson"
	secrets[sakey] = conf.Tracing.ServiceAccountJSON
	existing[sakey] = conf.Tracing.ExistingServiceAccountJSON

	if conf.Secrets == nil {
		conf.Secrets = &configuration.Secrets{}
	}

	if conf.Secrets.Keys == nil {
		conf.Secrets.Keys = &secret.Secret{}
	}
	if conf.Secrets.ExistingKeys == nil {
		conf.Secrets.ExistingKeys = &secret.Existing{}
	}
	keysKey := "keys"
	secrets[keysKey] = conf.Secrets.Keys
	existing[keysKey] = conf.Secrets.ExistingKeys

	if conf.Notifications == nil {
		conf.Notifications = &configuration.Notifications{}
	}

	if conf.Notifications.GoogleChatURL == nil {
		conf.Notifications.GoogleChatURL = &secret.Secret{}
	}
	if conf.Notifications.ExistingGoogleChatURL == nil {
		conf.Notifications.ExistingGoogleChatURL = &secret.Existing{}
	}
	gchatkey := "googlechaturl"
	secrets[gchatkey] = conf.Notifications.GoogleChatURL
	existing[gchatkey] = conf.Notifications.ExistingGoogleChatURL

	if conf.Notifications.Twilio == nil {
		conf.Notifications.Twilio = &configuration.Twilio{}
	}
	if conf.Notifications.Twilio.SID == nil {
		conf.Notifications.Twilio.SID = &secret.Secret{}
	}
	if conf.Notifications.Twilio.ExistingSID == nil {
		conf.Notifications.Twilio.ExistingSID = &secret.Existing{}
	}
	twilKey := "twiliosid"
	secrets[twilKey] = conf.Notifications.Twilio.SID
	existing[twilKey] = conf.Notifications.Twilio.ExistingSID

	if conf.Notifications.Twilio.AuthToken == nil {
		conf.Notifications.Twilio.AuthToken = &secret.Secret{}
	}
	if conf.Notifications.Twilio.ExistingAuthToken == nil {
		conf.Notifications.Twilio.ExistingAuthToken = &secret.Existing{}
	}
	twilOAuthKey := "twilioauthtoken"
	secrets[twilOAuthKey] = conf.Notifications.Twilio.AuthToken
	existing[twilOAuthKey] = conf.Notifications.Twilio.ExistingAuthToken

	if conf.Notifications.Email == nil {
		conf.Notifications.Email = &configuration.Email{}
	}
	if conf.Notifications.Email.AppKey == nil {
		conf.Notifications.Email.AppKey = &secret.Secret{}
	}
	if conf.Notifications.Email.ExistingAppKey == nil {
		conf.Notifications.Email.ExistingAppKey = &secret.Existing{}
	}
	mailKey := "emailappkey"
	secrets[mailKey] = conf.Notifications.Email.AppKey
	existing[mailKey] = conf.Notifications.Email.ExistingAppKey

	if conf.AssetStorage == nil {
		conf.AssetStorage = &configuration.AssetStorage{}
	}
	if conf.AssetStorage.AccessKeyID == nil {
		conf.AssetStorage.AccessKeyID = &secret.Secret{}
	}
	if conf.AssetStorage.ExistingAccessKeyID == nil {
		conf.AssetStorage.ExistingAccessKeyID = &secret.Existing{}
	}
	if conf.AssetStorage.SecretAccessKey == nil {
		conf.AssetStorage.SecretAccessKey = &secret.Secret{}
	}
	if conf.AssetStorage.ExistingSecretAccessKey == nil {
		conf.AssetStorage.ExistingSecretAccessKey = &secret.Existing{}
	}
	accessKey := "accesskeyid"
	secrets[accessKey] = conf.AssetStorage.AccessKeyID
	existing[accessKey] = conf.AssetStorage.ExistingAccessKeyID

	secretKey := "secretaccesskey"
	secrets[secretKey] = conf.AssetStorage.SecretAccessKey
	existing[secretKey] = conf.AssetStorage.ExistingSecretAccessKey

	if conf.Sentry == nil {
		conf.Sentry = &configuration.Sentry{}
	}
	if conf.Sentry.SentryDSN == nil {
		conf.Sentry.SentryDSN = &secret.Secret{}
	}
	if conf.Sentry.ExistingSentryDSN == nil {
		conf.Sentry.ExistingSentryDSN = &secret.Existing{}
	}

	SentryDSN := "sentrydsn"
	secrets[SentryDSN] = conf.Sentry.SentryDSN
	existing[SentryDSN] = conf.Sentry.ExistingSentryDSN

	return secrets, existing
}
