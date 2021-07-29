package configuration

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/caos/orbos/mntr"

	"github.com/caos/orbos/pkg/secret/read"

	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"
)

const (
	jsonName = "environment.json"
)

func literalsConfigMap(
	desired *Configuration,
	users map[string]string,
	certPath, secretPath, googleServiceAccountJSONPath, zitadelKeysPath string,
	queried map[string]interface{},
) map[string]string {

	tls := ""
	if desired.Notifications.Email.TLS {
		tls = "TRUE"
	} else {
		tls = "FALSE"
	}

	literalsConfigMap := map[string]string{
		"GOOGLE_APPLICATION_CREDENTIALS": secretPath + "/" + googleServiceAccountJSONPath,
		"ZITADEL_KEY_PATH":               secretPath + "/" + zitadelKeysPath,
		"ZITADEL_LOG_LEVEL":              "info",
		"DEBUG_MODE":                     strconv.FormatBool(desired.DebugMode),
		"ZITADEL_MIGRATE_ES_V1":          strconv.FormatBool(desired.MigrateEventStoreV1),
		"SMTP_TLS":                       tls,
		"CAOS_OIDC_DEV":                  "true",
		"CR_SSL_MODE":                    "require",
		"CR_ROOT_CERT":                   certPath + "/ca.crt",
	}

	if users != nil {
		for user := range users {
			literalsConfigMap["CR_"+strings.ToUpper(user)+"_CERT"] = certPath + "/client." + user + ".crt"
			literalsConfigMap["CR_"+strings.ToUpper(user)+"_KEY"] = certPath + "/client." + user + ".key"
		}
	}

	if desired != nil {
		if desired.Tracing != nil {
			literalsConfigMap["ZITADEL_TRACING_PROJECT_ID"] = desired.Tracing.ProjectID
			literalsConfigMap["ZITADEL_TRACING_FRACTION"] = desired.Tracing.Fraction
			literalsConfigMap["ZITADEL_TRACING_TYPE"] = desired.Tracing.Type
		}
		if desired.Secrets != nil {
			literalsConfigMap["ZITADEL_USER_VERIFICATION_KEY"] = desired.Secrets.UserVerificationID
			literalsConfigMap["ZITADEL_OTP_VERIFICATION_KEY"] = desired.Secrets.OTPVerificationID
			literalsConfigMap["ZITADEL_OIDC_KEYS_ID"] = desired.Secrets.OIDCKeysID
			literalsConfigMap["ZITADEL_COOKIE_KEY"] = desired.Secrets.CookieID
			literalsConfigMap["ZITADEL_CSRF_KEY"] = desired.Secrets.CSRFID
			literalsConfigMap["ZITADEL_DOMAIN_VERIFICATION_KEY"] = desired.Secrets.DomainVerificationID
			literalsConfigMap["ZITADEL_IDP_CONFIG_VERIFICATION_KEY"] = desired.Secrets.IDPConfigVerificationID
		}
		if desired.Notifications != nil {
			literalsConfigMap["TWILIO_SENDER_NAME"] = desired.Notifications.Twilio.SenderName
			literalsConfigMap["SMTP_HOST"] = desired.Notifications.Email.SMTPHost
			literalsConfigMap["SMTP_USER"] = desired.Notifications.Email.SMTPUser
			literalsConfigMap["EMAIL_SENDER_ADDRESS"] = desired.Notifications.Email.SenderAddress
			literalsConfigMap["EMAIL_SENDER_NAME"] = desired.Notifications.Email.SenderName
		}
		if desired.Cache != nil {
			literalsConfigMap["ZITADEL_CACHE_MAXAGE"] = desired.Cache.MaxAge
			literalsConfigMap["ZITADEL_CACHE_SHARED_MAXAGE"] = desired.Cache.SharedMaxAge
			literalsConfigMap["ZITADEL_SHORT_CACHE_MAXAGE"] = desired.Cache.ShortMaxAge
			literalsConfigMap["ZITADEL_SHORT_CACHE_SHARED_MAXAGE"] = desired.Cache.ShortSharedMaxAge
		}

		if desired.LogLevel != "" {
			literalsConfigMap["ZITADEL_LOG_LEVEL"] = desired.LogLevel
		}

		if desired.DNS != nil {
			defaultDomain := desired.DNS.Domain
			accountsDomain := desired.DNS.Subdomains.Accounts + "." + defaultDomain
			accounts := "https://" + accountsDomain
			issuer := "https://" + desired.DNS.Subdomains.Issuer + "." + defaultDomain
			oauth := "https://" + desired.DNS.Subdomains.API + "." + defaultDomain + "/oauth/v2"
			authorize := "https://" + desired.DNS.Subdomains.Accounts + "." + defaultDomain + "/oauth/v2"
			console := "https://" + desired.DNS.Subdomains.Console + "." + defaultDomain
			apiDomain := "https://" + desired.DNS.Subdomains.API + "." + defaultDomain

			literalsConfigMap["ZITADEL_ISSUER"] = issuer
			literalsConfigMap["ZITADEL_ACCOUNTS"] = accounts
			literalsConfigMap["ZITADEL_OAUTH"] = oauth
			literalsConfigMap["ZITADEL_AUTHORIZE"] = authorize
			literalsConfigMap["ZITADEL_CONSOLE"] = console
			literalsConfigMap["ZITADEL_ACCOUNTS_DOMAIN"] = accountsDomain
			literalsConfigMap["ZITADEL_COOKIE_DOMAIN"] = accountsDomain
			literalsConfigMap["ZITADEL_DEFAULT_DOMAIN"] = defaultDomain
			literalsConfigMap["ZITADEL_API_DOMAIN"] = apiDomain
		}
		if desired.AssetStorage != nil {
			literalsConfigMap["ZITADEL_ASSET_STORAGE_TYPE"] = desired.AssetStorage.Type
			literalsConfigMap["ZITADEL_ASSET_STORAGE_ENDPOINT"] = desired.AssetStorage.Endpoint
			literalsConfigMap["ZITADEL_ASSET_STORAGE_SSL"] = strconv.FormatBool(desired.AssetStorage.SSL)
			literalsConfigMap["ZITADEL_ASSET_STORAGE_LOCATION"] = desired.AssetStorage.Location
			literalsConfigMap["ZITADEL_ASSET_STORAGE_BUCKET_PREFIX"] = desired.AssetStorage.BucketPrefix
			literalsConfigMap["ZITADEL_ASSET_STORAGE_MULTI_DELETE"] = strconv.FormatBool(desired.AssetStorage.MultiDelete)
		}
	}

	sentryEnv, _, doIngest := mntr.Environment()
	literalsConfigMap["SENTRY_ENVIRONMENT"] = sentryEnv
	literalsConfigMap["SENTRY_USAGE"] = strconv.FormatBool(doIngest)

	db, err := database.GetDatabaseInQueried(queried)
	if err == nil {
		literalsConfigMap["ZITADEL_EVENTSTORE_HOST"] = db.Host
		literalsConfigMap["ZITADEL_EVENTSTORE_PORT"] = db.Port
	}

	return literalsConfigMap
}

func literalsSecret(k8sClient kubernetes.ClientInt, desired *Configuration, googleServiceAccountJSONPath, zitadelKeysPath string) (map[string]string, error) {
	literalsSecret := map[string]string{}
	if desired != nil {
		if desired.Tracing != nil && (desired.Tracing.ServiceAccountJSON != nil || desired.Tracing.ExistingServiceAccountJSON != nil) {
			value, err := read.GetSecretValue(k8sClient, desired.Tracing.ServiceAccountJSON, desired.Tracing.ExistingServiceAccountJSON)
			if err != nil {
				return nil, err
			}
			literalsSecret[googleServiceAccountJSONPath] = value
		}
		if desired.Secrets != nil && (desired.Secrets.Keys != nil || desired.Secrets.ExistingKeys != nil) {
			value, err := read.GetSecretValue(k8sClient, desired.Secrets.Keys, desired.Secrets.ExistingKeys)
			if err != nil {
				return nil, err
			}
			literalsSecret[zitadelKeysPath] = value
		}
	}
	return literalsSecret, nil
}

func literalsSecretVars(k8sClient kubernetes.ClientInt, desired *Configuration) (map[string]string, error) {
	literalsSecretVars := map[string]string{}
	if desired != nil {
		if desired.Notifications != nil {
			if desired.Notifications.Email.AppKey != nil || desired.Notifications.Email.ExistingAppKey != nil {
				value, err := read.GetSecretValue(k8sClient, desired.Notifications.Email.AppKey, desired.Notifications.Email.ExistingAppKey)
				if err != nil {
					return nil, err
				}
				literalsSecretVars["ZITADEL_EMAILAPPKEY"] = value
			}
			if desired.Notifications.GoogleChatURL != nil || desired.Notifications.ExistingGoogleChatURL != nil {
				value, err := read.GetSecretValue(k8sClient, desired.Notifications.GoogleChatURL, desired.Notifications.ExistingGoogleChatURL)
				if err != nil {
					return nil, err
				}
				literalsSecretVars["ZITADEL_GOOGLE_CHAT_URL"] = value
			}
			if desired.Notifications.Twilio.AuthToken != nil || desired.Notifications.Twilio.ExistingAuthToken != nil {
				value, err := read.GetSecretValue(k8sClient, desired.Notifications.Twilio.AuthToken, desired.Notifications.Twilio.ExistingAuthToken)
				if err != nil {
					return nil, err
				}
				literalsSecretVars["ZITADEL_TWILIO_AUTH_TOKEN"] = value
			}
			if desired.Notifications.Twilio.SID != nil || desired.Notifications.Twilio.ExistingSID != nil {
				value, err := read.GetSecretValue(k8sClient, desired.Notifications.Twilio.SID, desired.Notifications.Twilio.ExistingSID)
				if err != nil {
					return nil, err
				}
				literalsSecretVars["ZITADEL_TWILIO_SID"] = value
			}
		}
		if desired.AssetStorage != nil {
			as := desired.AssetStorage
			if as.AccessKeyID != nil || as.ExistingAccessKeyID != nil {
				value, err := read.GetSecretValue(k8sClient, as.AccessKeyID, as.ExistingAccessKeyID)
				if err != nil {
					return nil, err
				}
				literalsSecretVars["ZITADEL_ASSET_STORAGE_ACCESS_KEY_ID"] = value
			}
			if as.SecretAccessKey != nil || as.ExistingSecretAccessKey != nil {
				value, err := read.GetSecretValue(k8sClient, as.SecretAccessKey, as.ExistingSecretAccessKey)
				if err != nil {
					return nil, err
				}
				literalsSecretVars["ZITADEL_ASSET_STORAGE_SECRET_ACCESS_KEY"] = value
			}
		}

		_, dsns, doIngest := mntr.Environment()
		zitadelDsn := ""
		if doIngest {
			zitadelDsn = dsns["zitadel"]
		}
		literalsSecretVars["SENTRY_DSN"] = zitadelDsn
	}
	return literalsSecretVars, nil
}

func literalsConsoleCM(
	clientID string,
	dns *DNS,
	k8sClient kubernetes.ClientInt,
	namespace string,
	cmName string,
) map[string]string {
	literalsConsoleCM := map[string]string{}
	consoleEnv := ConsoleEnv{
		ClientID: clientID,
	}

	cm, err := k8sClient.GetConfigMap(namespace, cmName)
	//only try to use the old CM when there is a CM found
	if cm != nil && err == nil {
		jsonData, ok := cm.Data[jsonName]
		if ok {
			literalsData := map[string]string{}
			err := json.Unmarshal([]byte(jsonData), &literalsData)
			if err == nil {
				oldClientID, ok := literalsData["clientid"]
				//only use the old ClientID if no new clientID is provided
				if ok && consoleEnv.ClientID == "" {
					consoleEnv.ClientID = oldClientID
				}
			}
		}
	}

	consoleEnv.Issuer = "https://" + dns.Subdomains.Issuer + "." + dns.Domain
	consoleEnv.AuthServiceURL = "https://" + dns.Subdomains.API + "." + dns.Domain
	consoleEnv.MgmtServiceURL = "https://" + dns.Subdomains.API + "." + dns.Domain
	consoleEnv.SubServiceURL = "https://" + dns.Subdomains.Subscription + "." + dns.Domain
	consoleEnv.AssetServiceURL = "https://" + dns.Subdomains.API + "." + dns.Domain

	data, err := json.Marshal(consoleEnv)
	if err != nil {
		return map[string]string{}
	}

	literalsConsoleCM[jsonName] = string(data)
	return literalsConsoleCM
}
