package main

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/logger"
	"github.com/crewjam/saml/samlidp"
	xrv "github.com/mattermost/xml-roundtrip-validator"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/web"
	"golang.org/x/crypto/bcrypt"
)

var key = func() crypto.PrivateKey {
	b, _ := pem.Decode([]byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA0OhbMuizgtbFOfwbK7aURuXhZx6VRuAs3nNibiuifwCGz6u9
yy7bOR0P+zqN0YkjxaokqFgra7rXKCdeABmoLqCC0U+cGmLNwPOOA0PaD5q5xKhQ
4Me3rt/R9C4Ca6k3/OnkxnKwnogcsmdgs2l8liT3qVHP04Oc7Uymq2v09bGb6nPu
fOrkXS9F6mSClxHG/q59AGOWsXK1xzIRV1eu8W2SNdyeFVU1JHiQe444xLoPul5t
InWasKayFsPlJfWNc8EoU8COjNhfo/GovFTHVjh9oUR/gwEFVwifIHihRE0Hazn2
EQSLaOr2LM0TsRsQroFjmwSGgI+X2bfbMTqWOQIDAQABAoIBAFWZwDTeESBdrLcT
zHZe++cJLxE4AObn2LrWANEv5AeySYsyzjRBYObIN9IzrgTb8uJ900N/zVr5VkxH
xUa5PKbOcowd2NMfBTw5EEnaNbILLm+coHdanrNzVu59I9TFpAFoPavrNt/e2hNo
NMGPSdOkFi81LLl4xoadz/WR6O/7N2famM+0u7C2uBe+TrVwHyuqboYoidJDhO8M
w4WlY9QgAUhkPyzZqrl+VfF1aDTGVf4LJgaVevfFCas8Ws6DQX5q4QdIoV6/0vXi
B1M+aTnWjHuiIzjBMWhcYW2+I5zfwNWRXaxdlrYXRukGSdnyO+DH/FhHePJgmlkj
NInADDkCgYEA6MEQFOFSCc/ELXYWgStsrtIlJUcsLdLBsy1ocyQa2lkVUw58TouW
RciE6TjW9rp31pfQUnO2l6zOUC6LT9Jvlb9PSsyW+rvjtKB5PjJI6W0hjX41wEO6
fshFELMJd9W+Ezao2AsP2hZJ8McCF8no9e00+G4xTAyxHsNI2AFTCQcCgYEA5cWZ
JwNb4t7YeEajPt9xuYNUOQpjvQn1aGOV7KcwTx5ELP/Hzi723BxHs7GSdrLkkDmi
Gpb+mfL4wxCt0fK0i8GFQsRn5eusyq9hLqP/bmjpHoXe/1uajFbE1fZQR+2LX05N
3ATlKaH2hdfCJedFa4wf43+cl6Yhp6ZA0Yet1r8CgYEAwiu1j8W9G+RRA5/8/DtO
yrUTOfsbFws4fpLGDTA0mq0whf6Soy/96C90+d9qLaC3srUpnG9eB0CpSOjbXXbv
kdxseLkexwOR3bD2FHX8r4dUM2bzznZyEaxfOaQypN8SV5ME3l60Fbr8ajqLO288
wlTmGM5Mn+YCqOg/T7wjGmcCgYBpzNfdl/VafOROVbBbhgXWtzsz3K3aYNiIjbp+
MunStIwN8GUvcn6nEbqOaoiXcX4/TtpuxfJMLw4OvAJdtxUdeSmEee2heCijV6g3
ErrOOy6EqH3rNWHvlxChuP50cFQJuYOueO6QggyCyruSOnDDuc0BM0SGq6+5g5s7
H++S/wKBgQDIkqBtFr9UEf8d6JpkxS0RXDlhSMjkXmkQeKGFzdoJcYVFIwq8jTNB
nJrVIGs3GcBkqGic+i7rTO1YPkquv4dUuiIn+vKZVoO6b54f+oPBXd4S0BnuEqFE
rdKNuCZhiaE2XD9L/O9KP1fh5bfEcKwazQ23EvpJHBMm8BGC+/YZNw==
-----END RSA PRIVATE KEY-----`))
	k, _ := x509.ParsePKCS1PrivateKey(b.Bytes)
	return k
}()

var cert = func() *x509.Certificate {
	b, _ := pem.Decode([]byte(`-----BEGIN CERTIFICATE-----
MIIDBzCCAe+gAwIBAgIJAPr/Mrlc8EGhMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNV
BAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xNTEyMjgxOTE5NDVaFw0yNTEyMjUxOTE5
NDVaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBANDoWzLos4LWxTn8Gyu2lEbl4WcelUbgLN5zYm4ron8A
hs+rvcsu2zkdD/s6jdGJI8WqJKhYK2u61ygnXgAZqC6ggtFPnBpizcDzjgND2g+a
ucSoUODHt67f0fQuAmupN/zp5MZysJ6IHLJnYLNpfJYk96lRz9ODnO1Mpqtr9PWx
m+pz7nzq5F0vRepkgpcRxv6ufQBjlrFytccyEVdXrvFtkjXcnhVVNSR4kHuOOMS6
D7pebSJ1mrCmshbD5SX1jXPBKFPAjozYX6PxqLxUx1Y4faFEf4MBBVcInyB4oURN
B2s59hEEi2jq9izNE7EbEK6BY5sEhoCPl9m32zE6ljkCAwEAAaNQME4wHQYDVR0O
BBYEFB9ZklC1Ork2zl56zg08ei7ss/+iMB8GA1UdIwQYMBaAFB9ZklC1Ork2zl56
zg08ei7ss/+iMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAVoTSQ5
pAirw8OR9FZ1bRSuTDhY9uxzl/OL7lUmsv2cMNeCB3BRZqm3mFt+cwN8GsH6f3uv
NONIhgFpTGN5LEcXQz89zJEzB+qaHqmbFpHQl/sx2B8ezNgT/882H2IH00dXESEf
y/+1gHg2pxjGnhRBN6el/gSaDiySIMKbilDrffuvxiCfbpPN0NRRiPJhd2ay9KuL
/RxQRl1gl9cHaWiouWWba1bSBb2ZPhv2rPMUsFo98ntkGCObDX6Y1SpkqmoTbrsb
GFsTG2DLxnvr4GdN1BSr0Uu/KV3adj47WkXVPeMYQti/bQmxQB8tRFhrw80qakTL
UzreO96WzlBBMtY=
-----END CERTIFICATE-----`))
	c, _ := x509.ParseCertificate(b.Bytes)
	return c
}()

// Example from https://github.com/crewjam/saml/blob/main/example/idp/idp.go
func main() {
	apiURL := os.Getenv("API_URL")
	pat := readPAT(os.Getenv("PAT_FILE"))
	domain := os.Getenv("API_DOMAIN")
	schema := os.Getenv("SCHEMA")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	baseURL, err := url.Parse(schema + "://" + host + ":" + port)
	if err != nil {

		panic(err)
	}

	idpServer, err := samlidp.New(samlidp.Options{
		URL:         *baseURL,
		Logger:      logger.DefaultLogger,
		Key:         key,
		Certificate: cert,
		Store:       &samlidp.MemoryStore{},
	})
	if err != nil {

		panic(err)
	}

	metadata, err := xml.MarshalIndent(idpServer.IDP.Metadata(), "", "  ")
	if err != nil {
		panic(err)
	}
	idpID, err := createZitadelResources(apiURL, pat, domain, metadata)
	if err != nil {
		panic(err)
	}

	lis := bind.Socket(":" + baseURL.Port())
	goji.Handle("/*", idpServer)

	go func() {
		goji.ServeListener(lis)
	}()

	addService(idpServer, apiURL+"/idps/"+idpID+"/saml/metadata")
	addUsers(idpServer)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	if err := lis.Close(); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
}

func readPAT(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	pat, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return strings.Trim(string(pat), "\n")
}

func addService(idpServer *samlidp.Server, spURLStr string) {
	metadataResp, err := http.Get(spURLStr)
	if err != nil {
		panic(err)
	}
	defer metadataResp.Body.Close()

	idpServer.HandlePutService(
		web.C{URLParams: map[string]string{"id": spURLStr}},
		httptest.NewRecorder(),
		httptest.NewRequest(http.MethodPost, spURLStr, metadataResp.Body),
	)
}

func getSPMetadata(r io.Reader) (spMetadata *saml.EntityDescriptor, err error) {
	var data []byte
	if data, err = io.ReadAll(r); err != nil {
		return nil, err
	}

	spMetadata = &saml.EntityDescriptor{}
	if err := xrv.Validate(bytes.NewBuffer(data)); err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(data, &spMetadata); err != nil {
		if err.Error() == "expected element type <EntityDescriptor> but have <EntitiesDescriptor>" {
			entities := &saml.EntitiesDescriptor{}
			if err := xml.Unmarshal(data, &entities); err != nil {
				return nil, err
			}

			for _, e := range entities.EntityDescriptors {
				if len(e.SPSSODescriptors) > 0 {
					return &e, nil
				}
			}

			// there were no SPSSODescriptors in the response
			return nil, errors.New("metadata contained no service provider metadata")
		}

		return nil, err
	}

	return spMetadata, nil
}

func addUsers(idpServer *samlidp.Server) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("hunter2"), bcrypt.DefaultCost)
	err := idpServer.Store.Put("/users/alice", samlidp.User{Name: "alice",
		HashedPassword: hashedPassword,
		Groups:         []string{"Administrators", "Users"},
		Email:          "alice@example.com",
		CommonName:     "Alice Smith",
		Surname:        "Smith",
		GivenName:      "Alice",
	})
	if err != nil {
		panic(err)
	}

	err = idpServer.Store.Put("/users/bob", samlidp.User{
		Name:           "bob",
		HashedPassword: hashedPassword,
		Groups:         []string{"Users"},
		Email:          "bob@example.com",
		CommonName:     "Bob Smith",
		Surname:        "Smith",
		GivenName:      "Bob",
	})
	if err != nil {
		panic(err)
	}
}

func createZitadelResources(apiURL, pat, domain string, metadata []byte) (string, error) {
	idpID, err := CreateIDP(apiURL, pat, domain, metadata)
	if err != nil {
		return "", err
	}
	return idpID, ActivateIDP(apiURL, pat, domain, idpID)
}

type createIDP struct {
	Name              string          `json:"name"`
	MetadataXml       string          `json:"metadataXml"`
	Binding           string          `json:"binding"`
	WithSignedRequest bool            `json:"withSignedRequest"`
	ProviderOptions   providerOptions `json:"providerOptions"`
	NameIdFormat      string          `json:"nameIdFormat"`
}
type providerOptions struct {
	IsLinkingAllowed  bool   `json:"isLinkingAllowed"`
	IsCreationAllowed bool   `json:"isCreationAllowed"`
	IsAutoCreation    bool   `json:"isAutoCreation"`
	IsAutoUpdate      bool   `json:"isAutoUpdate"`
	AutoLinking       string `json:"autoLinking"`
}

type idp struct {
	ID string `json:"id"`
}

func CreateIDP(apiURL, pat, domain string, idpMetadata []byte) (string, error) {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(idpMetadata)))
	base64.URLEncoding.Encode(encoded, idpMetadata)

	createIDP := &createIDP{
		Name:              "CREWJAM",
		MetadataXml:       string(encoded),
		Binding:           "SAML_BINDING_REDIRECT",
		WithSignedRequest: false,
		ProviderOptions: providerOptions{
			IsLinkingAllowed:  true,
			IsCreationAllowed: true,
			IsAutoCreation:    true,
			IsAutoUpdate:      true,
			AutoLinking:       "AUTO_LINKING_OPTION_USERNAME",
		},
		NameIdFormat: "SAML_NAME_ID_FORMAT_PERSISTENT",
	}

	resp, err := doRequestWithHeaders(apiURL+"/admin/v1/idps/saml", pat, domain, createIDP)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	idp := new(idp)
	if err := json.Unmarshal(data, idp); err != nil {
		return "", err
	}
	return idp.ID, nil
}

type activateIDP struct {
	IdpId string `json:"idpId"`
}

func ActivateIDP(apiURL, pat, domain string, idpID string) error {
	activateIDP := &activateIDP{
		IdpId: idpID,
	}
	_, err := doRequestWithHeaders(apiURL+"/admin/v1/policies/login/idps", pat, domain, activateIDP)
	return err
}

func doRequestWithHeaders(apiURL, pat, domain string, body any) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, io.NopCloser(bytes.NewReader(data)))
	if err != nil {
		return nil, err
	}
	values := http.Header{}
	values.Add("Authorization", "Bearer "+pat)
	values.Add("x-forwarded-host", domain)
	values.Add("Content-Type", "application/json")
	req.Header = values

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
