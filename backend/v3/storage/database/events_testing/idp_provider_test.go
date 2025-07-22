//go:build integration

package events_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/admin"
	"github.com/zitadel/zitadel/pkg/grpc/idp"
	idp_grpc "github.com/zitadel/zitadel/pkg/grpc/idp"
)

func TestServer_TestIDProviderReduces(t *testing.T) {
	instanceID := Instance.ID()

	t.Run("test idp add reduces", func(t *testing.T) {
		name := gofakeit.Name()

		beforeCreate := time.Now()
		addOCID, err := AdminClient.AddOIDCIDP(CTX, &admin.AddOIDCIDPRequest{
			Name:               name,
			StylingType:        idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE,
			ClientId:           "clientID",
			ClientSecret:       "clientSecret",
			Issuer:             "issuer",
			Scopes:             []string{"scope"},
			DisplayNameMapping: idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			UsernameMapping:    idp.OIDCMappingField_OIDC_MAPPING_FIELD_EMAIL,
			AutoRegister:       true,
		})
		require.NoError(t, err)
		afterCreate := time.Now()

		idpRepo := repository.IDProviderRepository(pool)

		retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
		assert.EventuallyWithT(t, func(t *assert.CollectT) {
			idp, err := idpRepo.Get(CTX,
				idpRepo.NameCondition(name),
				instanceID,
				nil,
			)
			require.NoError(t, err)

			// event iam.idp.config.added
			assert.Equal(t, addOCID.IdpId, idp.ID)
			assert.Equal(t, name, idp.Name)
			assert.Equal(t, instanceID, idp.InstanceID)
			assert.Equal(t, domain.OrgStateActive.String(), idp.State)
			assert.Equal(t, int16(idp_grpc.IDPStylingType_STYLING_TYPE_GOOGLE), idp.StylingType)
			assert.WithinRange(t, idp.CreatedAt, beforeCreate, afterCreate)
			assert.WithinRange(t, idp.UpdatedAt, beforeCreate, afterCreate)
		}, retryDuration, tick)
	})
}
