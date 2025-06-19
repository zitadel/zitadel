//go:build integration

package instance_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/zitadel/zitadel/internal/integration"
	project "github.com/zitadel/zitadel/pkg/grpc/project/v2beta"
)

var (
	IAMOwnerCtx          context.Context
	instance             *integration.Instance
	instancePermissionV2 *integration.Instance
	baseURI              = "http://example.com"
	Project              *project.CreateProjectResponse
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		instance = integration.NewInstance(ctx)
		IAMOwnerCtx = instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		Project = instance.CreateProject(IAMOwnerCtx, &testing.T{}, instance.DefaultOrg.GetId(), gofakeit.Name(), false, false)

		return m.Run()
	}())
}

func samlMetadataGen(entityID string) []byte {
	str := fmt.Sprintf(`<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata"
                     validUntil="2022-08-26T14:08:16Z"
                     cacheDuration="PT604800S"
                     entityID="%s">
    <md:SPSSODescriptor AuthnRequestsSigned="false" WantAssertionsSigned="false" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
                                     Location="https://test.com/saml/acs"
                                     index="1" />
        
    </md:SPSSODescriptor>
</md:EntityDescriptor>
`,
		entityID)

	return []byte(str)
}
