package node

import (
	"crypto/rsa"
	"errors"
	"github.com/caos/zitadel/operator"
	"reflect"

	"github.com/caos/orbos/pkg/labels"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/secret"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate/certificates"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate/pem"
)

const (
	caCertKey      = "ca.crt"
	caPrivKeyKey   = "ca.key"
	nodeCertKey    = "node.crt"
	nodePrivKeyKey = "node.key"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	nameLabels *labels.Name,
	clusterDns string,
	generateIfNotExists bool,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {

	caPrivKey := new(rsa.PrivateKey)
	caCert := make([]byte, 0)
	nodeSecretSelector := labels.MustK8sMap(labels.DeriveNameSelector(nameLabels, false))

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			queriers := make([]operator.QueryFunc, 0)

			currentDB, err := core.ParseQueriedForDatabase(queried)
			if err != nil {
				return nil, err
			}

			allNodeSecrets, err := k8sClient.ListSecrets(namespace, nodeSecretSelector)
			if err != nil {
				return nil, err
			}

			if len(allNodeSecrets.Items) == 0 {

				if !generateIfNotExists {
					return nil, errors.New("node secret not found")
				}

				emptyCert := true
				emptyKey := true
				if currentCaCert := currentDB.GetCertificate(); currentCaCert != nil && len(currentCaCert) != 0 {
					emptyCert = false
					caCert = currentCaCert
				}
				if currentCaCertKey := currentDB.GetCertificateKey(); currentCaCertKey != nil && !reflect.DeepEqual(currentCaCertKey, &rsa.PrivateKey{}) {
					emptyKey = false
					caPrivKey = currentCaCertKey
				}

				if emptyCert || emptyKey {
					caPrivKeyInternal, caCertInternal, err := certificates.NewCA()
					if err != nil {
						return nil, err
					}
					caPrivKey = caPrivKeyInternal
					caCert = caCertInternal

					nodePrivKey, nodeCert, err := certificates.NewNode(caPrivKey, caCert, namespace, clusterDns)
					if err != nil {
						return nil, err
					}

					pemNodePrivKey, err := pem.EncodeKey(nodePrivKey)
					if err != nil {
						return nil, err
					}
					pemCaPrivKey, err := pem.EncodeKey(caPrivKey)
					if err != nil {
						return nil, err
					}

					pemCaCert, err := pem.EncodeCertificate(caCert)
					if err != nil {
						return nil, err
					}

					pemNodeCert, err := pem.EncodeCertificate(nodeCert)
					if err != nil {
						return nil, err
					}

					nodeSecretData := map[string]string{
						caPrivKeyKey:   string(pemCaPrivKey),
						caCertKey:      string(pemCaCert),
						nodePrivKeyKey: string(pemNodePrivKey),
						nodeCertKey:    string(pemNodeCert),
					}
					queryNodeSecret, err := secret.AdaptFuncToEnsure(namespace, labels.AsSelectable(nameLabels), nodeSecretData)
					if err != nil {
						return nil, err
					}
					queriers = append(queriers, operator.ResourceQueryToZitadelQuery(queryNodeSecret))
				}
			} else {
				key, err := pem.DecodeKey(allNodeSecrets.Items[0].Data[caPrivKeyKey])
				if err != nil {
					return nil, err
				}
				caPrivKey = key

				cert, err := pem.DecodeCertificate(allNodeSecrets.Items[0].Data[caCertKey])
				if err != nil {
					return nil, err
				}
				caCert = cert
			}

			currentDB.SetCertificate(caCert)
			currentDB.SetCertificateKey(caPrivKey)

			return operator.QueriersToEnsureFunc(monitor, false, queriers, k8sClient, queried)
		}, func(k8sClient kubernetes.ClientInt) error {
			allNodeSecrets, err := k8sClient.ListSecrets(namespace, nodeSecretSelector)
			if err != nil {
				return err
			}
			for _, deleteSecret := range allNodeSecrets.Items {
				destroyer, err := secret.AdaptFuncToDestroy(namespace, deleteSecret.Name)
				if err != nil {
					return err
				}
				if err := destroyer(k8sClient); err != nil {
					return err
				}
			}
			return nil
		}, nil
}
