package client

import (
	"errors"
	"github.com/caos/zitadel/pkg/databases/db"

	"github.com/caos/zitadel/operator/database/kinds/databases/managed/current"

	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/secret"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/user/certificates"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/user/pem"
)

const (
	clientSecretPrefix = "cockroachdb.client."
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	componentLabels *labels.Component,
) (
	func(client, secretName, userCrtFilename, userKeyFilename string) operator.QueryFunc,
	func(secretName string) operator.DestroyFunc,
	error,
) {

	return func(client, secretName, userCrtFilename, userKeyFilename string) operator.QueryFunc {
			nameLabels := labels.MustForName(componentLabels, secretName)

			return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				queriers := make([]operator.QueryFunc, 0)

				currentDB, err := db.ParseQueriedForDatabase(queried)
				if err != nil {
					return nil, err
				}

				managedDB := currentDB.(*current.Current)

				caCert := managedDB.GetCertificate()
				caKey := managedDB.GetCertificateKey()
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
					db.CACert:       string(pemCaCert),
					userKeyFilename: string(pemClientPrivKey),
					userCrtFilename: string(pemClientCert),
				}

				queryClientSecret, err := secret.AdaptFuncToEnsure(namespace, labels.AsSelectable(nameLabels), clientSecretData)
				if err != nil {
					return nil, err
				}
				queriers = append(queriers, operator.ResourceQueryToZitadelQuery(queryClientSecret))

				return operator.QueriersToEnsureFunc(monitor, false, queriers, k8sClient, queried)
			}
		}, func(secretName string) operator.DestroyFunc {

			destroy, err := secret.AdaptFuncToDestroy(namespace, secretName)
			if err != nil {
				return nil
			}
			return operator.ResourceDestroyToZitadelDestroy(destroy)
		},
		nil
}
