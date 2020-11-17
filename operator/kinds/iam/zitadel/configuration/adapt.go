package configuration

import (
	"encoding/json"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/configuration/users"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/database"
	"strconv"
	"strings"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/configmap"
	"github.com/caos/orbos/pkg/kubernetes/resources/secret"
	"github.com/caos/zitadel/operator"
)

type ConsoleEnv struct {
	AuthServiceURL string `json:"authServiceUrl"`
	MgmtServiceURL string `json:"mgmtServiceUrl"`
	Issuer         string `json:"issuer"`
	ClientID       string `json:"clientid"`
}

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	labels map[string]string,
	desired *Configuration,
	cmName string,
	certPath string,
	secretName string,
	secretPath string,
	consoleCMName string,
	secretVarsName string,
	secretPasswordName string,
	necessaryUsers map[string]string,
	getClientID func() string,
	repoURL string,
	repoKey string,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	func(k8sClient kubernetes.ClientInt) map[string]string,
	error,
) {
	internalMonitor := monitor.WithField("component", "configuration")

	googleServiceAccountJSONPath := "google-serviceaccount-key.json"
	zitadelKeysPath := "zitadel-keys.yaml"

	literalsSecret := literalsSecret(desired, googleServiceAccountJSONPath, zitadelKeysPath)
	literalsSecretVars := literalsSecretVars(desired)

	destroyCM, err := configmap.AdaptFuncToDestroy(namespace, cmName)
	if err != nil {
		return nil, nil, nil, err
	}
	destroyS, err := secret.AdaptFuncToDestroy(namespace, secretName)
	if err != nil {
		return nil, nil, nil, err
	}
	destroyCCM, err := configmap.AdaptFuncToDestroy(namespace, consoleCMName)
	if err != nil {
		return nil, nil, nil, err
	}
	destroySV, err := secret.AdaptFuncToDestroy(namespace, secretVarsName)
	if err != nil {
		return nil, nil, nil, err
	}
	destroySP, err := secret.AdaptFuncToDestroy(namespace, secretPasswordName)
	if err != nil {
		return nil, nil, nil, err
	}

	queryUser, destroyUser, err := users.AdaptFunc(internalMonitor, necessaryUsers, repoURL, repoKey)
	if err != nil {
		return nil, nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		destroyUser,
		operator.ResourceDestroyToZitadelDestroy(destroyS),
		operator.ResourceDestroyToZitadelDestroy(destroyCM),
		operator.ResourceDestroyToZitadelDestroy(destroyCCM),
		operator.ResourceDestroyToZitadelDestroy(destroySV),
		operator.ResourceDestroyToZitadelDestroy(destroySP),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {

			queryS, err := secret.AdaptFuncToEnsure(namespace, secretName, labels, literalsSecret)
			if err != nil {
				return nil, err
			}
			querySV, err := secret.AdaptFuncToEnsure(namespace, secretVarsName, labels, literalsSecretVars)
			if err != nil {
				return nil, err
			}
			querySP, err := secret.AdaptFuncToEnsure(namespace, secretPasswordName, labels, necessaryUsers)
			if err != nil {
				return nil, err
			}

			queryCCM, err := configmap.AdaptFuncToEnsure(
				namespace,
				consoleCMName,
				labels,
				literalsConsoleCM(
					getClientID(),
					desired.DNS,
					k8sClient,
					namespace,
					consoleCMName,
				),
			)
			if err != nil {
				return nil, err
			}

			queryCM, err := configmap.AdaptFuncToEnsure(
				namespace,
				cmName,
				labels,
				literalsConfigMap(
					monitor,
					desired,
					necessaryUsers,
					certPath,
					secretPath,
					googleServiceAccountJSONPath,
					zitadelKeysPath,
					k8sClient,
					repoURL,
					repoKey,
				),
			)
			if err != nil {
				return nil, err
			}

			queriers := []operator.QueryFunc{
				queryUser,
				operator.ResourceQueryToZitadelQuery(queryS),
				operator.ResourceQueryToZitadelQuery(queryCCM),
				operator.ResourceQueryToZitadelQuery(querySV),
				operator.ResourceQueryToZitadelQuery(querySP),
				operator.ResourceQueryToZitadelQuery(queryCM),
			}

			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		func(k8sClient kubernetes.ClientInt) map[string]string {
			return map[string]string{
				secretName:         getHash(literalsSecret),
				secretVarsName:     getHash(literalsSecretVars),
				secretPasswordName: getHash(necessaryUsers),
				cmName: getHash(
					literalsConfigMap(
						monitor,
						desired,
						necessaryUsers,
						certPath,
						secretPath,
						googleServiceAccountJSONPath,
						zitadelKeysPath,
						k8sClient,
						repoURL,
						repoKey,
					),
				),
				consoleCMName: getHash(
					literalsConsoleCM(
						getClientID(),
						desired.DNS,
						k8sClient,
						namespace,
						consoleCMName,
					),
				),
			}
		},
		nil
}

func literalsConfigMap(
	monitor mntr.Monitor,
	desired *Configuration,
	users map[string]string,
	certPath, secretPath, googleServiceAccountJSONPath, zitadelKeysPath string,
	k8sClient kubernetes.ClientInt,
	repoURL, repoKey string,
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
		for _, user := range users {
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

	url, port, err := database.GetConnectionInfo(monitor, k8sClient, repoURL, repoKey)
	if err == nil {
		literalsConfigMap["ZITADEL_EVENTSTORE_HOST"] = url
		literalsConfigMap["ZITADEL_EVENTSTORE_PORT"] = port
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

	jsonName := "environment.json"

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
