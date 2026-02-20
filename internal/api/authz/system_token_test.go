package authz

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
)

var exampleRsaPrivateKey *rsa.PrivateKey
var exampleRsaPublicKeyBs []byte

func init() {
	exampleRsaPrivateKey, _ = rsa.GenerateKey(rand.Reader, 2048)
	publicBs, _ := x509.MarshalPKIXPublicKey(&exampleRsaPrivateKey.PublicKey)
	writer := &bytes.Buffer{}
	_ = pem.Encode(writer, &pem.Block{Type: "PUBLIC KEY", Bytes: publicBs})
	exampleRsaPublicKeyBs = writer.Bytes()
}

func Test_SystemAPIUser_readKey(t *testing.T) {
	t.Run("rsa key", func(tt *testing.T) {
		// given
		publicKey := exampleRsaPrivateKey.PublicKey
		user := SystemAPIUser{
			KeyData: exampleRsaPublicKeyBs,
		}

		// when
		key, err := user.readKey()

		// then
		assert.NoError(tt, err)
		assert.Nil(tt, key.NotBefore)
		assert.Nil(tt, key.NotAfter)
		assert.True(tt, publicKey.Equal(key.Data))
	})

	t.Run("x.509 cert", func(tt *testing.T) {
		// given
		publicKey := exampleRsaPrivateKey.PublicKey
		now := time.Now().Round(time.Second)
		cert := createExampleX509Cert(now, now.Add(1*time.Hour))
		user := SystemAPIUser{
			KeyData: encodeCert(cert),
		}

		// when
		key, err := user.readKey()

		// then
		assert.NoError(tt, err)
		assert.Equal(tt, cert.NotBefore.UTC(), key.NotBefore.UTC())
		assert.Equal(tt, cert.NotAfter.UTC(), key.NotAfter.UTC())
		assert.True(tt, publicKey.Equal(key.Data))
	})
}

func Test_systemJWTStorage_GetKeyByIDAndClientID_Ok(t *testing.T) {
	type TestCase struct {
		name    string
		storage *systemJWTStorage
		userID  string
		keyID   string
		key     *SystemAPIPublicKey
	}

	testCases := []TestCase{
		func() TestCase {
			key := &SystemAPIPublicKey{Data: &exampleRsaPrivateKey.PublicKey}
			return TestCase{
				name: "get from cache, no notBefore or notAfter",
				storage: &systemJWTStorage{
					cachedKeys: map[string]*SystemAPIPublicKey{
						"user-1": key,
					},
				},
				userID: "user-1",
				key:    key,
			}
		}(),
		func() TestCase {
			key := &SystemAPIPublicKey{
				Data:      &exampleRsaPrivateKey.PublicKey,
				NotBefore: gu.Ptr(time.Now().UTC()),
				NotAfter:  gu.Ptr(time.Now().UTC().Add(time.Second * 2)),
			}
			return TestCase{
				name: "get from cache, with notBefore and notAfter",
				storage: &systemJWTStorage{
					cachedKeys: map[string]*SystemAPIPublicKey{
						"user-2": key,
					},
				},
				userID: "user-2",
				key:    key,
			}
		}(),
		func() TestCase {
			now := time.Now().Add(-1 * time.Second).Round(time.Second).UTC()
			until := now.Add(1 * time.Hour)
			cert := createExampleX509Cert(now, until)
			return TestCase{
				name: "no cache, with notBefore and notAfter",
				storage: &systemJWTStorage{
					cachedKeys: make(map[string]*SystemAPIPublicKey),
					keys: map[string]*SystemAPIUser{
						"user": {
							KeyData: encodeCert(cert),
						},
					},
				},
				userID: "user",
				key: &SystemAPIPublicKey{
					Data:      &exampleRsaPrivateKey.PublicKey,
					NotBefore: &now,
					NotAfter:  &until,
				},
			}
		}(),
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			jwk, err := tc.storage.GetKeyByIDAndClientID(context.Background(), tc.keyID, tc.userID)
			assert.NoError(tt, err)
			assert.Equal(tt, tc.key, jwk.Key)
		})
	}
}

func Test_systemJWTStorage_GetKeyByIDAndClientID_Nok(t *testing.T) {
	type TestCase struct {
		name    string
		storage *systemJWTStorage
		userID  string
		keyID   string
		err     string
	}

	testCases := []TestCase{
		{
			name: "user not found",
			storage: &systemJWTStorage{
				cachedKeys: make(map[string]*SystemAPIPublicKey),
			},
			userID: "does not exist",
			err:    "AUTHZ-asfd3",
		},
		{
			name: "get from cache, not before",
			storage: &systemJWTStorage{
				cachedKeys: map[string]*SystemAPIPublicKey{
					"user": {
						NotBefore: gu.Ptr(time.Now().UTC().Add(time.Second * 2)),
					},
				},
			},
			userID: "user",
			err:    "AUTHZ-NiJstf",
		},
		{
			name: "get from cache, not after",
			storage: &systemJWTStorage{
				cachedKeys: map[string]*SystemAPIPublicKey{
					"user": {
						NotAfter: gu.Ptr(time.Now().UTC()),
					},
				},
			},
			userID: "user",
			err:    "AUTHZ-CGmV4b",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			jwk, err := tc.storage.GetKeyByIDAndClientID(context.Background(), tc.keyID, tc.userID)
			assert.Nil(tt, jwk)
			assert.ErrorContains(tt, err, tc.err)
		})
	}
}

func createExampleX509Cert(notBefore, notAfter time.Time) *x509.Certificate {
	return &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{gofakeit.Company()},
		},
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		SubjectKeyId: []byte{1, 2, 3, 4, 5},
		PublicKey:    &exampleRsaPrivateKey.PublicKey,
	}
}

func encodeCert(pub *x509.Certificate) []byte {
	ca := createExampleX509Cert(time.Now(), time.Now().Add(240*time.Hour))
	certBytes, err := x509.CreateCertificate(rand.Reader, pub, ca, &exampleRsaPrivateKey.PublicKey, exampleRsaPrivateKey)
	if err != nil {
		panic(err)
	}
	buff := &bytes.Buffer{}
	_ = pem.Encode(buff, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	return buff.Bytes()
}
