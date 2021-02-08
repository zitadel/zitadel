package configuration

import (
	"encoding/json"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"
	"strconv"
	"strings"
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

			literalsConfigMap["ZITADEL_ISSUER"] = issuer
			literalsConfigMap["ZITADEL_ACCOUNTS"] = accounts
			literalsConfigMap["ZITADEL_OAUTH"] = oauth
			literalsConfigMap["ZITADEL_AUTHORIZE"] = authorize
			literalsConfigMap["ZITADEL_CONSOLE"] = console
			literalsConfigMap["ZITADEL_ACCOUNTS_DOMAIN"] = accountsDomain
			literalsConfigMap["ZITADEL_COOKIE_DOMAIN"] = accountsDomain
			literalsConfigMap["ZITADEL_DEFAULT_DOMAIN"] = defaultDomain
		}
	}

	db, err := database.GetDatabaseInQueried(queried)
	if err == nil {
		literalsConfigMap["ZITADEL_EVENTSTORE_HOST"] = db.Host
		literalsConfigMap["ZITADEL_EVENTSTORE_PORT"] = db.Port
	}

	return literalsConfigMap
}

func literalsSecret(desired *Configuration, googleServiceAccountJSONPath, zitadelKeysPath string) map[string]string {
	literalsSecret := map[string]string{}
	if desired != nil {
		if desired.Tracing != nil && desired.Tracing.ServiceAccountJSON != nil {
			literalsSecret[googleServiceAccountJSONPath] = desired.Tracing.ServiceAccountJSON.Value
		}
		if desired.Secrets != nil && desired.Secrets.Keys != nil {
			literalsSecret[zitadelKeysPath] = desired.Secrets.Keys.Value
		}
	}
	return literalsSecret
}

func literalsSecretVars(desired *Configuration) map[string]string {
	literalsSecretVars := map[string]string{}
	if desired != nil {
		if desired.Notifications != nil {
			if desired.Notifications.Email.AppKey != nil {
				literalsSecretVars["ZITADEL_EMAILAPPKEY"] = desired.Notifications.Email.AppKey.Value
			}
			if desired.Notifications.GoogleChatURL != nil {
				literalsSecretVars["ZITADEL_GOOGLE_CHAT_URL"] = desired.Notifications.GoogleChatURL.Value
			}
			if desired.Notifications.Twilio.AuthToken != nil {
				literalsSecretVars["ZITADEL_TWILIO_AUTH_TOKEN"] = desired.Notifications.Twilio.AuthToken.Value
			}
			if desired.Notifications.Twilio.SID != nil {
				literalsSecretVars["ZITADEL_TWILIO_SID"] = desired.Notifications.Twilio.SID.Value
			}
		}
	}
	return literalsSecretVars
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

	data, err := json.Marshal(consoleEnv)
	if err != nil {
		return map[string]string{}
	}

	literalsConsoleCM[jsonName] = string(data)
	return literalsConsoleCM
}
