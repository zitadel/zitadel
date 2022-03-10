package client

import (
	"errors"
	"github.com/caos/zitadel/operator"
	"strings"

	"github.com/caos/orbos/pkg/labels"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/secret"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate/certificates"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate/pem"
)

const (
	clientSecretPrefix     = "cockroachdb.client."
	caCertKey              = "ca.crt"
	clientCertKeyPrefix    = "client."
	clientCertKeySuffix    = ".crt"
	clientPrivKeyKeyPrefix = "client."
	clientPrivKeyKeySuffix = ".key"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	componentLabels *labels.Component,
) (
	func(client string) operator.QueryFunc,
	func(client string) operator.DestroyFunc,
	error,
) {

	return func(client string) operator.QueryFunc {
			clientSecret := clientSecretPrefix + client
			nameLabels := labels.MustForName(componentLabels, strings.ReplaceAll(clientSecret, "_", "-"))
			clientCertKey := clientCertKeyPrefix + client + clientCertKeySuffix
			clientPrivKeyKey := clientPrivKeyKeyPrefix + client + clientPrivKeyKeySuffix

			return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				queriers := make([]operator.QueryFunc, 0)

				currentDB, err := core.ParseQueriedForDatabase(queried)
				if err != nil {
					return nil, err
				}

				caCert := currentDB.GetCertificate()
				caKey := currentDB.GetCertificateKey()
				if caKey == nil || caCert == nil || len(caCert) == 0 {
					return nil, errors.New("no ca-certificate found")
				}

				clientPrivKey, clientCert, err := certificates.NewClient(caKey, caCert, client)
				if err != nil {
					return nil, err
				}

				pemClientPrivKey, err := pem.EncodeKey(clientPrivKey)
				if err != nil {
					return nil, err
				}

				pemClientCert, err := pem.EncodeCertificate(clientCert)
				if err != nil {
					return nil, err
				}

				pemCaCert, err := pem.EncodeCertificate(caCert)
				if err != nil {
					return nil, err
				}

				clientSecretData := map[string]string{
					caCertKey:        string(pemCaCert),
					clientPrivKeyKey: string(pemClientPrivKey),
					clientCertKey:    string(pemClientCert),
				}

				queryClientSecret, err := secret.AdaptFuncToEnsure(namespace, labels.AsSelectable(nameLabels), clientSecretData)
				if err != nil {
					return nil, err
				}
				queriers = append(queriers, operator.ResourceQueryToZitadelQuery(queryClientSecret))

				return operator.QueriersToEnsureFunc(monitor, false, queriers, k8sClient, queried)
			}
		}, func(client string) operator.DestroyFunc {
			clientSecret := clientSecretPrefix + client

			destroy, err := secret.AdaptFuncToDestroy(namespace, clientSecret)
			if err != nil {
				return nil
			}
			return operator.ResourceDestroyToZitadelDestroy(destroy)
		},
		nil
}
