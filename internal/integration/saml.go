package integration

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
	"github.com/zitadel/logging"
	"github.com/zitadel/saml/pkg/provider"

	http_util "github.com/zitadel/zitadel/internal/api/http"
	oidc_internal "github.com/zitadel/zitadel/internal/api/oidc"
	app_pb "github.com/zitadel/zitadel/pkg/grpc/app"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	saml_pb "github.com/zitadel/zitadel/pkg/grpc/saml/v2"
	session_pb "github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

const spCertificate = `-----BEGIN CERTIFICATE-----
MIIDITCCAgmgAwIBAgIUUo5urYkuUHAe7LQ9sZSL+xXAqBwwDQYJKoZIhvcNAQEL
BQAwIDEeMBwGA1UEAwwVbXlzZXJ2aWNlLmV4YW1wbGUuY29tMB4XDTI0MTIwNDEz
MTE1MFoXDTI1MDEwMzEzMTE1MFowIDEeMBwGA1UEAwwVbXlzZXJ2aWNlLmV4YW1w
bGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoACwbGIh8udK
Um1r+yQoPtfswEX6Cb6Y1KwR6WZDYgzHdMyUC5Sy8Bg1H2puUZZukDLuyu6Pqvum
8kfnzhjUR6nNCoUlidwE+yz020w5oOBofRKgJK/FVUuWD3k6kjdP9CrBFLG0PQQ3
N2e4wilP4czCxizKero2a0e7Eq8OjHAPf8gjM+GWFZgVAbV8uf2Mjt1O2Vfbx5PZ
sLuBZtl5jokx3NiC7my/yj81MbGEDPcQo0emeVBz3J3nVG6Yr4kdCKkvv2dhJ26C
5cL7NIIUY4IQomJNwYC2NaYgSpQOxJHL/HsOPusO4Ia2WtUTXEZUFkxn1u0YuoSx
CkGehF/1OwIDAQABo1MwUTAdBgNVHQ4EFgQUr6S0wA2l3MdfnvfveWDueQtaoJMw
HwYDVR0jBBgwFoAUr6S0wA2l3MdfnvfveWDueQtaoJMwDwYDVR0TAQH/BAUwAwEB
/zANBgkqhkiG9w0BAQsFAAOCAQEAH3Q9obyWJaMKFuGJDkIp1RFot79RWTVcAcwA
XTJNfCseLONRIs4MkRxOn6GQBwV2IEqs1+hFG80dcd/c6yYyJ8bziKEyNMtPWrl6
fdVD+1WnWcD1ZYrS8hgdz0FxXxl/+GjA8Pu6icmnhKgUDTYWns6Rj/gtQtZS8ZoA
JY+T/1mGze2+Xx6pjuArZ7+hnH6EWwo+ckcmXAKyhnkhX7xIo1UFvNY2VWaGl2wU
K2yyJA4Lu/NNmqPnpAcRDsnGP6r4frMhjnPq/ifC3B+6FT3p8dubV9PA0y86bAy5
0yIgNje4DyWLy/DM9EpdPfJmvUAL6hOtyb8Aa9hR+a8stu7h6g==
-----END CERTIFICATE-----`
const spKey = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCgALBsYiHy50pS
bWv7JCg+1+zARfoJvpjUrBHpZkNiDMd0zJQLlLLwGDUfam5Rlm6QMu7K7o+q+6by
R+fOGNRHqc0KhSWJ3AT7LPTbTDmg4Gh9EqAkr8VVS5YPeTqSN0/0KsEUsbQ9BDc3
Z7jCKU/hzMLGLMp6ujZrR7sSrw6McA9/yCMz4ZYVmBUBtXy5/YyO3U7ZV9vHk9mw
u4Fm2XmOiTHc2ILubL/KPzUxsYQM9xCjR6Z5UHPcnedUbpiviR0IqS+/Z2EnboLl
wvs0ghRjghCiYk3BgLY1piBKlA7Ekcv8ew4+6w7ghrZa1RNcRlQWTGfW7Ri6hLEK
QZ6EX/U7AgMBAAECggEAD1aRkwpDO+BdORKhP9WDACc93F647fc1+mk2XFv/yKX1
9uXnqUaLcsW3TfgrdCnKFouzZYPCBP+TzPUErTanHumRrNj/tLwBRDzWijE/8wKg
MaE39dxdu+P/kiMqcLrZsMvqb3vrjc/aJTcNuJsyO7Cf2VSQ4nv4XIdnUQ60A9VR
OmUp//VULZxImnPx/R304/p5VfOhyXfzBeoxUPogBurjtzkyXVG0EG2enJMMiTix
900fTDez0TQ8V6O59vM04fhtPXvH51OkMTW/HU1QQvlnAJuX06I7k4CaBpF3xPII
QpEbFILq5y6yAQJWELRGWzeoxK6kn6bNfI8S0+oKqQKBgQDg2UM7ruMASpY7B4zj
XNztGDOx9BCdYyHH1O05r+ILmltBC7jFImwIYrHbaX+dg52l0PPImZuBz852IqrC
VAEF30yBn2gWyVzIdo7W3mw9Jgqc4LrhStaJxOuXVoT2/PAuDBF8TJMNH9oLNqiD
aPAI0cVn9BRV7AziEsrMlDLLiQKBgQC2K4Z/caAvwx/AescsN6lp+/m7MeLUpZzQ
myZt44bnR5LouUo3vCYl+Bk8wu6PTd41LUYW/SW26HDDFTKgkBb1zVHfk5QRApaB
VPwZnhcUvNapPOnDp75Qoq238wpfayQlKF1xCawS3N5AWkDaEdfzuH7umFJxVss2
1tfDsn01owKBgAYWG3nMHBzv5+0lIS0uYFSSqSOSBbkc69cq7lj3Z9kEjp/OH2xG
qEH52fKkgm3TGDta0p6Fee4jn+UWvySPfY+ZIcsIc5raTIaonuk2EBv/oZ3pf2WF
zxTfnbj1AJhm9GFqtjZ1JC3gxNg03I7iEk1K0FsmAj7pKtgbxh2PjWhxAoGBAKBx
BSwJbwOh3r0vZWvUOilV+0SbUyPmGI7Blr8BvTbFGuZNCsi7tP2L3O5e4Kzl7+b1
0N0+Z5EIdwfaC5TOUup5wroeyDGTDesqZj5JthpVltnHBDuF6WArZsS0EVaojlUL
kACWfC7AyB31X1iwjnng7CpHjZS01JWf8rgw44XxAoGAQ5YYd4WmGYZoJJak7zhb
xnYG7hU7nS7pBPGob1FvjYMw1x/htuJCjxLh08dlzJGM6SFlDn7HVM9ou99w5n+d
xtqmbthw2E9VjSk3zSYb4uFc6mv0C/kRPTDUFH+9CpQTBBx/O016hmcatxlBS6JL
VAV6oE8sEJYHtR6YdZiMWWo=
-----END PRIVATE KEY-----`

func CreateSAMLSP(root string, idpMetadata *saml.EntityDescriptor, binding string) (*samlsp.Middleware, error) {
	rootURL, err := url.Parse(root)
	if err != nil {
		return nil, err
	}
	keyPair, err := tls.X509KeyPair([]byte(spCertificate), []byte(spKey))
	if err != nil {
		return nil, err
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		return nil, err
	}

	sp, err := samlsp.New(samlsp.Options{
		URL:                 *rootURL,
		Key:                 keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate:         keyPair.Leaf,
		IDPMetadata:         idpMetadata,
		UseArtifactResponse: false,
	})
	if err != nil {
		return nil, err
	}
	sp.Binding = binding
	sp.ResponseBinding = binding
	return sp, nil
}

func (i *Instance) CreateSAMLClientLoginVersion(ctx context.Context, projectID string, m *samlsp.Middleware, loginVersion *app_pb.LoginVersion) (*management.AddSAMLAppResponse, error) {
	spMetadata, err := xml.MarshalIndent(m.ServiceProvider.Metadata(), "", "  ")
	if err != nil {
		return nil, err
	}

	if m.ResponseBinding == saml.HTTPRedirectBinding {
		metadata := strings.Replace(string(spMetadata), saml.HTTPPostBinding, saml.HTTPRedirectBinding, 2)
		spMetadata = []byte(metadata)
	}

	resp, err := i.Client.Mgmt.AddSAMLApp(ctx, &management.AddSAMLAppRequest{
		ProjectId:    projectID,
		Name:         ApplicationName(),
		Metadata:     &management.AddSAMLAppRequest_MetadataXml{MetadataXml: spMetadata},
		LoginVersion: loginVersion,
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

func (i *Instance) CreateSAMLClient(ctx context.Context, projectID string, m *samlsp.Middleware) (*management.AddSAMLAppResponse, error) {
	return i.CreateSAMLClientLoginVersion(ctx, projectID, m, nil)
}

func (i *Instance) CreateSAMLAuthRequestWithoutLoginClientHeader(m *samlsp.Middleware, loginBaseURI string, acs saml.Endpoint, relayState, responseBinding string) (now time.Time, authRequestID string, err error) {
	return i.createSAMLAuthRequest(m, "", loginBaseURI, acs, relayState, responseBinding)
}

func (i *Instance) CreateSAMLAuthRequest(m *samlsp.Middleware, loginClient string, acs saml.Endpoint, relayState, responseBinding string) (now time.Time, authRequestID string, err error) {
	return i.createSAMLAuthRequest(m, loginClient, "", acs, relayState, responseBinding)
}

func (i *Instance) createSAMLAuthRequest(m *samlsp.Middleware, loginClient, loginBaseURI string, acs saml.Endpoint, relayState, responseBinding string) (now time.Time, authRequestID string, err error) {
	authReq, err := m.ServiceProvider.MakeAuthenticationRequest(acs.Location, acs.Binding, responseBinding)
	if err != nil {
		return now, "", err
	}

	redirectURL, err := authReq.Redirect(relayState, &m.ServiceProvider)
	if err != nil {
		return now, "", err
	}

	var headers map[string]string
	if loginClient != "" {
		headers = map[string]string{oidc_internal.LoginClientHeader: loginClient}
	}
	req, err := GetRequest(redirectURL.String(), headers)
	if err != nil {
		return now, "", fmt.Errorf("get request: %w", err)
	}

	now = time.Now()
	loc, err := CheckRedirect(req)
	if err != nil {
		return now, "", fmt.Errorf("check redirect: %w", err)
	}

	if loginBaseURI == "" {
		loginBaseURI = i.Issuer() + i.Config.LoginURLV2
	}
	if !strings.HasPrefix(loc.String(), loginBaseURI) {
		return now, "", fmt.Errorf("login location has not prefix %s, but is %s", loginBaseURI, loc.String())
	}
	return now, strings.TrimPrefix(loc.String(), loginBaseURI), nil
}

func (i *Instance) FailSAMLAuthRequest(ctx context.Context, id string, reason saml_pb.ErrorReason) *saml_pb.CreateResponseResponse {
	resp, err := i.Client.SAMLv2.CreateResponse(ctx, &saml_pb.CreateResponseRequest{
		SamlRequestId: id,
		ResponseKind:  &saml_pb.CreateResponseRequest_Error{Error: &saml_pb.AuthorizationError{Error: reason}},
	})
	logging.OnError(err).Panic("create human user")
	return resp
}

func (i *Instance) SuccessfulSAMLAuthRequest(ctx context.Context, userId, id string) *saml_pb.CreateResponseResponse {
	respSession, err := i.Client.SessionV2.CreateSession(ctx, &session_pb.CreateSessionRequest{
		Checks: &session_pb.Checks{
			User: &session_pb.CheckUser{
				Search: &session_pb.CheckUser_UserId{
					UserId: userId,
				},
			},
		},
	})
	logging.OnError(err).Panic("create session")

	resp, err := i.Client.SAMLv2.CreateResponse(ctx, &saml_pb.CreateResponseRequest{
		SamlRequestId: id,
		ResponseKind: &saml_pb.CreateResponseRequest_Session{
			Session: &saml_pb.Session{
				SessionId:    respSession.GetSessionId(),
				SessionToken: respSession.GetSessionToken(),
			},
		},
	})
	logging.OnError(err).Panic("create human user")
	return resp
}

func (i *Instance) GetSAMLIDPMetadata() (*saml.EntityDescriptor, error) {
	issuer := i.Issuer() + "/saml/v2"
	idpEntityID := issuer + "/metadata"

	req, err := http.NewRequestWithContext(provider.ContextWithIssuer(context.Background(), issuer), http.MethodGet, idpEntityID, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	entityDescriptor := new(saml.EntityDescriptor)
	if err := xml.Unmarshal(data, entityDescriptor); err != nil {
		return nil, err
	}

	return entityDescriptor, nil
}

func (i *Instance) Issuer() string {
	return http_util.BuildHTTP(i.Domain, i.Config.Port, i.Config.Secure)
}
