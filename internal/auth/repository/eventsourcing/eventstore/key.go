package eventstore

import (
	"context"
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/caos/logging"
	"gopkg.in/square/go-jose.v2"

	"github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/v1/spooler"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/key/model"
	key_view "github.com/caos/zitadel/internal/key/repository/view"
)

type KeyRepository struct {
	Commands                 *command.Commands
	Eventstore               *eventstore.Eventstore
	View                     *view.View
	SigningKeyRotationCheck  time.Duration
	SigningKeyGracefulPeriod time.Duration
	KeyAlgorithm             crypto.EncryptionAlgorithm
	KeyChan                  <-chan *model.KeyView
	CaCertChan               <-chan *model.CertificateAndKeyView
	MetadataCertChan         <-chan *model.CertificateAndKeyView
	ResponseCertChan         <-chan *model.CertificateAndKeyView
	Locker                   spooler.Locker
	lockID                   string
	currentKeyID             string
	currentKeyExpiration     time.Time
}

type CertificateAndKey struct {
	Certificate *jose.SigningKey
	Key         *jose.SigningKey
}

const (
	signingKey             = "signing_key"
	samlMetadataSigningKey = "saml_metadata_singing_key"
	samlResponseSigningKey = "saml_response_singing_key"
	samlCaKey              = "saml_ca_key"
)

func (k *KeyRepository) GetSigningKey(ctx context.Context, keyCh chan<- jose.SigningKey, algorithm string, usage model.KeyUsage) {
	renewTimer := time.After(0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case key := <-k.KeyChan:
				refreshed, err := k.refreshSigningKey(ctx, key, keyCh, algorithm, usage)
				logging.Log("KEY-asd5g").OnError(err).Error("could not refresh signing key on key channel push")
				renewTimer = time.After(k.getRenewTimer(refreshed))
			case <-renewTimer:
				key, err := k.latestSigningKey(usage)
				logging.Log("KEY-DAfh4-1").OnError(err).Error("could not check for latest signing key")
				refreshed, err := k.refreshSigningKey(ctx, key, keyCh, algorithm, usage)
				logging.Log("KEY-DAfh4-2").OnError(err).Error("could not refresh signing key when ensuring key")
				renewTimer = time.After(k.getRenewTimer(refreshed))
			}
		}
	}()
}

func (k *KeyRepository) GetCertificateAndKey(ctx context.Context, certAndKeyCh chan<- CertificateAndKey, algorithm string, usage model.KeyUsage) {
	renewTimer := time.After(0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case certAndKey := <-k.CaCertChan:
				if certAndKey.Key.Usage != 0 && certAndKey.Key.Usage == usage {
					refreshed, err := k.refreshCertificate(ctx, certAndKey, certAndKeyCh, algorithm, usage)
					logging.Log("KEY-asd8g").OnError(err).Error("could not refresh certificate on key channel push")
					renewTimer = time.After(k.getRenewTimer(refreshed))
				}
			case certAndKey := <-k.MetadataCertChan:
				if certAndKey.Key.Usage != 0 && certAndKey.Key.Usage == usage {
					refreshed, err := k.refreshCertificate(ctx, certAndKey, certAndKeyCh, algorithm, usage)
					logging.Log("KEY-asd6g").OnError(err).Error("could not refresh certificate on key channel push")
					renewTimer = time.After(k.getRenewTimer(refreshed))
				}
			case certAndKey := <-k.ResponseCertChan:
				if certAndKey.Key.Usage != 0 && certAndKey.Key.Usage == usage {
					refreshed, err := k.refreshCertificate(ctx, certAndKey, certAndKeyCh, algorithm, usage)
					logging.Log("KEY-asd7g").OnError(err).Error("could not refresh certificate on key channel push")
					renewTimer = time.After(k.getRenewTimer(refreshed))
				}
			case <-renewTimer:
				certAndKey, err := k.latestCertificateAndKey(usage)
				logging.Log("KEY-DAfh4-1").OnError(err).Error("could not check for latest certificate")
				refreshed, err := k.refreshCertificate(ctx, certAndKey, certAndKeyCh, algorithm, usage)
				logging.Log("KEY-DAfh4-2").OnError(err).Error("could not refresh certificate when ensuring")
				renewTimer = time.After(k.getRenewTimer(refreshed))
			}
		}
	}()
}

func (k *KeyRepository) GetKeySet(ctx context.Context, usage model.KeyUsage) (*jose.JSONWebKeySet, error) {
	keys, err := k.View.GetActiveKeySet(usage)
	if err != nil {
		return nil, err
	}
	webKeys := make([]jose.JSONWebKey, len(keys))
	for i, key := range keys {
		webKeys[i] = jose.JSONWebKey{KeyID: key.ID, Algorithm: key.Algorithm, Use: key.Usage.String(), Key: key.Key}
	}
	return &jose.JSONWebKeySet{Keys: webKeys}, nil
}

func (k *KeyRepository) getRenewTimer(refreshed bool) time.Duration {
	duration := k.SigningKeyRotationCheck
	if refreshed {
		duration = k.currentKeyExpiration.Sub(time.Now().Add(k.SigningKeyGracefulPeriod + k.SigningKeyRotationCheck*2))
	}
	logging.LogWithFields("EVENT-dK432", "in", duration).Info("next signing key check")
	return duration
}

func (k *KeyRepository) latestCertificateAndKey(usage model.KeyUsage) (shortRefresh *model.CertificateAndKeyView, err error) {
	var ret *model.CertificateAndKeyView
	switch usage {
	case model.KeyUsageSAMLCA, model.KeyUsageSAMLResponseSinging, model.KeyUsageSAMLMetadataSigning:
		certAndKey, errView := k.View.GetActiveCertificateAndKeyForSigning(time.Now().UTC().Add(k.SigningKeyGracefulPeriod), usage)
		if errView != nil && !errors.IsNotFound(errView) {
			logging.Log("EVENT-GEd4h").WithError(errView).Warn("could not get signing key")
			return nil, errView
		}
		ret = certAndKey
	}
	return ret, nil
}

func (k *KeyRepository) latestSigningKey(usage model.KeyUsage) (shortRefresh *model.KeyView, err error) {
	var ret *model.KeyView
	switch usage {
	case model.KeyUsageSigning:
		key, errView := k.View.GetActivePrivateKeyForSigning(time.Now().UTC().Add(k.SigningKeyGracefulPeriod), usage)
		if errView != nil && !errors.IsNotFound(errView) {
			logging.Log("EVENT-GEd4h").WithError(errView).Warn("could not get signing key")
			return nil, errView
		}
		ret = key
	}
	return ret, nil
}

func (k *KeyRepository) ensureIsLatestKey(ctx context.Context) (bool, error) {
	sequence, err := k.View.GetLatestKeySequence()
	if err != nil {
		return false, err
	}
	events, err := k.getKeyEvents(ctx, sequence.CurrentSequence)
	if err != nil {
		logging.Log("EVENT-der5g").WithError(err).Warn("error retrieving new events")
		return false, err
	}
	if len(events) > 0 {
		logging.Log("EVENT-GBD23").Warn("view not up to date, retrying later")
		return false, nil
	}
	return true, nil
}

func (k *KeyRepository) refreshCertificate(
	ctx context.Context,
	certAndKey *model.CertificateAndKeyView,
	certAndKeyCh chan<- CertificateAndKey,
	algorithm string,
	usage model.KeyUsage,
) (refreshed bool, err error) {
	if certAndKey == nil {
		if k.currentKeyExpiration.Before(time.Now().UTC()) {
			logging.Log("EVENT-ADg26").Info("unset current signing key")
			certAndKeyCh <- CertificateAndKey{
				Certificate: &jose.SigningKey{},
				Key:         &jose.SigningKey{},
			}
		}
		if ok, err := k.ensureIsLatestKey(ctx); !ok && err == nil {
			return false, err
		}
		logging.Log("EVENT-sdz53").Info("lock and generate signing key pair")
		err = k.lockAndGenerateCertificateAndKey(ctx, algorithm, usage)
		logging.Log("EVENT-B4d21").OnError(err).Warn("could not create signing key")
		return false, err
	}

	if k.currentKeyID == certAndKey.Key.ID {
		logging.Log("EVENT-Abb3e").Info("no new signing key")
		return false, nil
	}
	if ok, err := k.ensureIsLatestKey(ctx); !ok && err == nil {
		logging.Log("EVENT-HJd92").Info("signing key in view is not latest key")
		return false, err
	}

	switch usage {
	case model.KeyUsageSAMLMetadataSigning, model.KeyUsageSAMLResponseSinging, model.KeyUsageSAMLCA:
		certAndKeyView, err := model.CertificateAndKeyFromCertificateAndKeyView(certAndKey, k.KeyAlgorithm)
		if err != nil {
			logging.Log("EVENT-HJd92").WithError(err).Error("cert cannot be decrypted -> immediate refresh")
			return k.refreshCertificate(ctx, nil, certAndKeyCh, algorithm, usage)
		}
		certAndKeyCh <- CertificateAndKey{
			Key: &jose.SigningKey{
				Algorithm: jose.SignatureAlgorithm(certAndKeyView.Key.Algorithm),
				Key: jose.JSONWebKey{
					KeyID: certAndKeyView.Key.ID,
					Key:   certAndKeyView.Key.Key,
				},
			},
			Certificate: &jose.SigningKey{
				Algorithm: jose.SignatureAlgorithm(certAndKeyView.Certificate.Algorithm),
				Key: jose.JSONWebKey{
					KeyID: certAndKeyView.Certificate.ID,
					Key:   certAndKeyView.Certificate.Certificate,
				},
			},
		}

		logging.LogWithFields("EVENT-dsg54", "keyID", certAndKey.Key.ID).Info("refreshed certificate")
	default:
	}

	return true, nil
}

func (k *KeyRepository) refreshSigningKey(ctx context.Context, key *model.KeyView, keyCh chan<- jose.SigningKey, algorithm string, usage model.KeyUsage) (refreshed bool, err error) {
	if key == nil {
		if k.currentKeyExpiration.Before(time.Now().UTC()) {
			logging.Log("EVENT-ADg26").Info("unset current signing key")
			keyCh <- jose.SigningKey{}
		}
		if ok, err := k.ensureIsLatestKey(ctx); !ok && err == nil {
			return false, err
		}
		logging.Log("EVENT-sdz53").Info("lock and generate signing key pair")
		err = k.lockAndGenerateSigningKeyPair(ctx, algorithm, usage)
		logging.Log("EVENT-B4d21").OnError(err).Warn("could not create signing key")
		return false, err
	}

	if k.currentKeyID == key.ID {
		logging.Log("EVENT-Abb3e").Info("no new signing key")
		return false, nil
	}
	if ok, err := k.ensureIsLatestKey(ctx); !ok && err == nil {
		logging.Log("EVENT-HJd92").Info("signing key in view is not latest key")
		return false, err
	}

	switch usage {
	case model.KeyUsageSigning:
		signingKey, err := model.SigningKeyFromKeyView(key, k.KeyAlgorithm)
		if err != nil {
			logging.Log("EVENT-HJd92").WithError(err).Error("signing key cannot be decrypted -> immediate refresh")
			return k.refreshSigningKey(ctx, nil, keyCh, algorithm, usage)
		}
		k.currentKeyID = signingKey.ID
		k.currentKeyExpiration = key.Expiry
		keyCh <- jose.SigningKey{
			Algorithm: jose.SignatureAlgorithm(signingKey.Algorithm),
			Key: jose.JSONWebKey{
				KeyID: signingKey.ID,
				Key:   signingKey.Key,
			},
		}
		logging.LogWithFields("EVENT-dsg54", "keyID", signingKey.ID).Info("refreshed signing key")
	}

	return true, nil
}

func (k *KeyRepository) lockAndGenerateCertificateAndKey(ctx context.Context, algorithm string, usage model.KeyUsage) error {
	keyType := ""
	switch usage {
	case model.KeyUsageSAMLMetadataSigning:
		keyType = samlMetadataSigningKey
	case model.KeyUsageSAMLResponseSinging:
		keyType = samlResponseSigningKey
	case model.KeyUsageSAMLCA:
		keyType = samlCaKey
	default:
		return fmt.Errorf("unknown certificate usage")
	}

	err := k.Locker.Renew(k.lockerID(), keyType, k.SigningKeyRotationCheck*2)
	if err != nil {
		if errors.IsErrorAlreadyExists(err) {
			return nil
		}
		return err
	}
	switch usage {
	case model.KeyUsageSAMLMetadataSigning, model.KeyUsageSAMLResponseSinging:
		done, err := k.ensureIsLatestKey(ctx)
		if err != nil || !done {
			//TODO
		}

		certAndKeyView, err := k.latestCertificateAndKey(model.KeyUsageSAMLCA)
		if err != nil {
			return err
		}
		certAndKey, err := model.CertificateAndKeyFromCertificateAndKeyView(certAndKeyView, k.KeyAlgorithm)
		if err != nil {
			return err
		}

		certData := certAndKey.Certificate.Certificate.([]byte)
		privateKey := certAndKey.Key.Key.(*rsa.PrivateKey)

		switch usage {
		case model.KeyUsageSAMLMetadataSigning:
			return k.Commands.GenerateSAMLMetadataCertificate(ctx, "metadata", algorithm, privateKey, certData)
		case model.KeyUsageSAMLResponseSinging:
			return k.Commands.GenerateSAMLResponseCertificate(ctx, "response", algorithm, privateKey, certData)
		default:
			return fmt.Errorf("unknown usage")
		}
	case model.KeyUsageSAMLCA:
		return k.Commands.GenerateSAMLCACertificate(ctx)
	default:
		return fmt.Errorf("unknown certificate usage")
	}
}

func (k *KeyRepository) lockAndGenerateSigningKeyPair(ctx context.Context, algorithm string, usage model.KeyUsage) error {
	keyType := ""
	switch usage {
	case model.KeyUsageSigning:
		keyType = signingKey
	default:
		return fmt.Errorf("unknown key usage")
	}

	err := k.Locker.Renew(k.lockerID(), keyType, k.SigningKeyRotationCheck*2)
	if err != nil {
		if errors.IsErrorAlreadyExists(err) {
			return nil
		}
		return err
	}
	switch usage {
	case model.KeyUsageSigning:
		return k.Commands.GenerateOIDCSigningPair(ctx, algorithm)
	default:
		return fmt.Errorf("unknown key usage")
	}
}

func (k *KeyRepository) lockerID() string {
	if k.lockID == "" {
		var err error
		k.lockID, err = os.Hostname()
		if err != nil || k.lockID == "" {
			k.lockID, err = id.SonyFlakeGenerator.Next()
			logging.Log("EVENT-bsdf6").OnError(err).Panic("unable to generate lockID")
		}
	}
	return k.lockID
}

func (k *KeyRepository) getKeyEvents(ctx context.Context, sequence uint64) ([]eventstore.EventReader, error) {
	return k.Eventstore.FilterEvents(ctx, key_view.KeyPairQuery(sequence))
}
