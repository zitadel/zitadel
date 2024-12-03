package integration

import (
	"context"
	"fmt"
	"net/url"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"

	oidc_internal "github.com/zitadel/zitadel/internal/api/oidc"
	"github.com/zitadel/zitadel/pkg/grpc/management"
)

func (i *Instance) CreateSAMLClient(ctx context.Context, projectID, entityID, acsURL, logoutURL string) (*management.AddSAMLAppResponse, error) {
	samlSPMetadata := `<?xml version="1.0"?>
<md:EntityDescriptor xmlns:md="urn:oasis:names:tc:SAML:2.0:metadata"
                     validUntil="2024-12-05T17:23:27Z"
                     cacheDuration="PT604800S"
                     entityID="` + entityID + `">
    <md:SPSSODescriptor AuthnRequestsSigned="true" WantAssertionsSigned="true" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
        <md:SingleLogoutService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect"
                                Location="` + logoutURL + `" />
        <md:NameIDFormat>urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified</md:NameIDFormat>
        <md:AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST"
                                     Location="` + acsURL + `"
                                     index="1" />
    </md:SPSSODescriptor>
</md:EntityDescriptor>`

	resp, err := i.Client.Mgmt.AddSAMLApp(ctx, &management.AddSAMLAppRequest{
		ProjectId: projectID,
		Name:      fmt.Sprintf("app-%s", gofakeit.AppName()),
		Metadata:  &management.AddSAMLAppRequest_MetadataXml{MetadataXml: []byte(samlSPMetadata)},
	})
	if err != nil {
		return nil, err
	}
	return resp, await(func() error {
		_, err := i.Client.Mgmt.GetProjectByID(ctx, &management.GetProjectByIDRequest{
			Id: projectID,
		})
		if err != nil {
			return err
		}
		_, err = i.Client.Mgmt.GetAppByID(ctx, &management.GetAppByIDRequest{
			ProjectId: projectID,
			AppId:     resp.GetAppId(),
		})
		return err
	})
}

func (i *Instance) CreateSAMLAuthRequest(ctx context.Context, entityID, loginClient, acsURL, relayState string) (authRequestID string, err error) {
	binding := saml.HTTPRedirectBinding
	entityDescriptor := new(saml.EntityDescriptor)
	rootURL, err := url.Parse(entityID)
	if err != nil {
		return "", err
	}

	m, _ := samlsp.New(samlsp.Options{
		URL: *rootURL,
		/* TODO
		Key:            keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:    keyPair.Leaf,
		*/
		IDPMetadata: entityDescriptor,
	})

	authReq, err := m.ServiceProvider.MakeAuthenticationRequest(acsURL, binding, m.ResponseBinding)
	if err != nil {
		return "", err
	}

	redirectURL, err := authReq.Redirect(relayState, &m.ServiceProvider)
	if err != nil {
		return "", err
	}

	req, err := GetRequest(redirectURL.String(), map[string]string{oidc_internal.LoginClientHeader: loginClient})
	if err != nil {
		return "", fmt.Errorf("get request: %w", err)
	}

	loc, err := CheckRedirect(req)
	if err != nil {
		return "", fmt.Errorf("check redirect: %w", err)
	}

	//TODO get id from loc
	return loc.String(), nil
}
