//go:build integration

package webkey_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/feature/v2"
	webkey "github.com/zitadel/zitadel/pkg/grpc/webkey/v2beta"
)

var (
	CTX context.Context
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()
		CTX = ctx
		return m.Run()
	}())
}

func TestServer_Feature_Disabled(t *testing.T) {
	instance, iamCtx, _ := createInstance(t, false)
	client := instance.Client.WebKeyV2Beta

	t.Run("CreateWebKey", func(t *testing.T) {
		_, err := client.CreateWebKey(iamCtx, &webkey.CreateWebKeyRequest{})
		assertFeatureDisabledError(t, err)
	})
	t.Run("ActivateWebKey", func(t *testing.T) {
		_, err := client.ActivateWebKey(iamCtx, &webkey.ActivateWebKeyRequest{
			Id: "1",
		})
		assertFeatureDisabledError(t, err)
	})
	t.Run("DeleteWebKey", func(t *testing.T) {
		_, err := client.DeleteWebKey(iamCtx, &webkey.DeleteWebKeyRequest{
			Id: "1",
		})
		assertFeatureDisabledError(t, err)
	})
	t.Run("ListWebKeys", func(t *testing.T) {
		_, err := client.ListWebKeys(iamCtx, &webkey.ListWebKeysRequest{})
		assertFeatureDisabledError(t, err)
	})
}

func TestServer_ListWebKeys(t *testing.T) {
	instance, iamCtx, creationDate := createInstance(t, true)
	// After the feature is first enabled, we can expect 2 generated keys with the default config.
	checkWebKeyListState(iamCtx, t, instance, 2, "", &webkey.WebKey_Rsa{
		Rsa: &webkey.RSA{
			Bits:   webkey.RSABits_RSA_BITS_2048,
			Hasher: webkey.RSAHasher_RSA_HASHER_SHA256,
		},
	}, creationDate)
}

func TestServer_CreateWebKey(t *testing.T) {
	instance, iamCtx, creationDate := createInstance(t, true)
	client := instance.Client.WebKeyV2Beta

	_, err := client.CreateWebKey(iamCtx, &webkey.CreateWebKeyRequest{
		Key: &webkey.CreateWebKeyRequest_Rsa{
			Rsa: &webkey.RSA{
				Bits:   webkey.RSABits_RSA_BITS_2048,
				Hasher: webkey.RSAHasher_RSA_HASHER_SHA256,
			},
		},
	})
	require.NoError(t, err)

	checkWebKeyListState(iamCtx, t, instance, 3, "", &webkey.WebKey_Rsa{
		Rsa: &webkey.RSA{
			Bits:   webkey.RSABits_RSA_BITS_2048,
			Hasher: webkey.RSAHasher_RSA_HASHER_SHA256,
		},
	}, creationDate)
}

func TestServer_ActivateWebKey(t *testing.T) {
	instance, iamCtx, creationDate := createInstance(t, true)
	client := instance.Client.WebKeyV2Beta

	resp, err := client.CreateWebKey(iamCtx, &webkey.CreateWebKeyRequest{
		Key: &webkey.CreateWebKeyRequest_Rsa{
			Rsa: &webkey.RSA{
				Bits:   webkey.RSABits_RSA_BITS_2048,
				Hasher: webkey.RSAHasher_RSA_HASHER_SHA256,
			},
		},
	})
	require.NoError(t, err)

	_, err = client.ActivateWebKey(iamCtx, &webkey.ActivateWebKeyRequest{
		Id: resp.GetId(),
	})
	require.NoError(t, err)

	checkWebKeyListState(iamCtx, t, instance, 3, resp.GetId(), &webkey.WebKey_Rsa{
		Rsa: &webkey.RSA{
			Bits:   webkey.RSABits_RSA_BITS_2048,
			Hasher: webkey.RSAHasher_RSA_HASHER_SHA256,
		},
	}, creationDate)
}

func TestServer_DeleteWebKey(t *testing.T) {
	instance, iamCtx, creationDate := createInstance(t, true)
	client := instance.Client.WebKeyV2Beta

	keyIDs := make([]string, 2)
	for i := 0; i < 2; i++ {
		resp, err := client.CreateWebKey(iamCtx, &webkey.CreateWebKeyRequest{
			Key: &webkey.CreateWebKeyRequest_Rsa{
				Rsa: &webkey.RSA{
					Bits:   webkey.RSABits_RSA_BITS_2048,
					Hasher: webkey.RSAHasher_RSA_HASHER_SHA256,
				},
			},
		})
		require.NoError(t, err)
		keyIDs[i] = resp.GetId()
	}
	_, err := client.ActivateWebKey(iamCtx, &webkey.ActivateWebKeyRequest{
		Id: keyIDs[0],
	})
	require.NoError(t, err)

	ok := t.Run("cannot delete active key", func(t *testing.T) {
		_, err := client.DeleteWebKey(iamCtx, &webkey.DeleteWebKeyRequest{
			Id: keyIDs[0],
		})
		require.Error(t, err)
		s := status.Convert(err)
		assert.Equal(t, codes.FailedPrecondition, s.Code())
		assert.Contains(t, s.Message(), "COMMAND-Chai1")
	})
	if !ok {
		return
	}

	start := time.Now()
	ok = t.Run("delete inactive key", func(t *testing.T) {
		resp, err := client.DeleteWebKey(iamCtx, &webkey.DeleteWebKeyRequest{
			Id: keyIDs[1],
		})
		require.NoError(t, err)
		require.WithinRange(t, resp.GetDeletionDate().AsTime(), start, time.Now())
	})
	if !ok {
		return
	}

	ok = t.Run("delete inactive key again", func(t *testing.T) {
		resp, err := client.DeleteWebKey(iamCtx, &webkey.DeleteWebKeyRequest{
			Id: keyIDs[1],
		})
		require.NoError(t, err)
		require.WithinRange(t, resp.GetDeletionDate().AsTime(), start, time.Now())
	})
	if !ok {
		return
	}

	ok = t.Run("delete not existing key", func(t *testing.T) {
		resp, err := client.DeleteWebKey(iamCtx, &webkey.DeleteWebKeyRequest{
			Id: "not-existing",
		})
		require.NoError(t, err)
		require.Nil(t, resp.DeletionDate)
	})
	if !ok {
		return
	}

	// There are 2 keys from feature setup, +2 created, -1 deleted = 3
	checkWebKeyListState(iamCtx, t, instance, 3, keyIDs[0], &webkey.WebKey_Rsa{
		Rsa: &webkey.RSA{
			Bits:   webkey.RSABits_RSA_BITS_2048,
			Hasher: webkey.RSAHasher_RSA_HASHER_SHA256,
		},
	}, creationDate)
}

func createInstance(t *testing.T, enableFeature bool) (*integration.Instance, context.Context, *timestamppb.Timestamp) {
	instance := integration.NewInstance(CTX)
	creationDate := timestamppb.Now()
	iamCTX := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	if enableFeature {
		_, err := instance.Client.FeatureV2.SetInstanceFeatures(iamCTX, &feature.SetInstanceFeaturesRequest{
			WebKey: proto.Bool(true),
		})
		require.NoError(t, err)
	}

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamCTX, time.Minute)
	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		resp, err := instance.Client.WebKeyV2Beta.ListWebKeys(iamCTX, &webkey.ListWebKeysRequest{})
		if enableFeature {
			assert.NoError(collect, err)
			assert.Len(collect, resp.GetWebKeys(), 2)
		} else {
			assert.Error(collect, err)
		}
	}, retryDuration, tick)

	return instance, iamCTX, creationDate
}

func assertFeatureDisabledError(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	s := status.Convert(err)
	assert.Equal(t, codes.FailedPrecondition, s.Code())
	assert.Contains(t, s.Message(), "WEBKEY-Ohx6E")
}

func checkWebKeyListState(ctx context.Context, t *testing.T, instance *integration.Instance, nKeys int, expectActiveKeyID string, config any, creationDate *timestamppb.Timestamp) {
	t.Helper()

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	assert.EventuallyWithT(t, func(collect *assert.CollectT) {
		resp, err := instance.Client.WebKeyV2Beta.ListWebKeys(ctx, &webkey.ListWebKeysRequest{})
		require.NoError(collect, err)
		list := resp.GetWebKeys()
		assert.Len(collect, list, nKeys)

		now := time.Now()
		var gotActiveKeyID string
		for _, key := range list {
			assert.WithinRange(collect, key.GetCreationDate().AsTime(), now.Add(-time.Minute), now.Add(time.Minute))
			assert.WithinRange(collect, key.GetChangeDate().AsTime(), now.Add(-time.Minute), now.Add(time.Minute))
			assert.NotEqual(collect, webkey.State_STATE_UNSPECIFIED, key.GetState())
			assert.NotEqual(collect, webkey.State_STATE_REMOVED, key.GetState())
			assert.Equal(collect, config, key.GetKey())

			if key.GetState() == webkey.State_STATE_ACTIVE {
				gotActiveKeyID = key.GetId()
			}
		}
		assert.NotEmpty(collect, gotActiveKeyID)
		if expectActiveKeyID != "" {
			assert.Equal(collect, expectActiveKeyID, gotActiveKeyID)
		}
	}, retryDuration, tick)
}
