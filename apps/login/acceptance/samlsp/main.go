package main

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/crewjam/saml/samlsp"
)

var keyPair = func() tls.Certificate {
	cert := []byte(`-----BEGIN CERTIFICATE-----
MIIDITCCAgmgAwIBAgIUKjAUmxsHO44X+/TKBNciPgNl1GEwDQYJKoZIhvcNAQEL
BQAwIDEeMBwGA1UEAwwVbXlzZXJ2aWNlLmV4YW1wbGUuY29tMB4XDTI0MTIxOTEz
Mzc1MVoXDTI1MTIxOTEzMzc1MVowIDEeMBwGA1UEAwwVbXlzZXJ2aWNlLmV4YW1w
bGUuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0QYuJsayILRI
hVT7G1DlitVSXnt1iw3gEXJZfe81Egz06fUbvXF6Yo1LJmwYpqe/rm+hf4FNUb8e
2O+LH2FieA9FkVe4P2gKOzw87A/KxvpV8stgNgl4LlqRCokbc1AzeE/NiLr5TcTD
RXm3DUcYxXxinprtDu2jftFysaOZmNAukvE/iL6qS3X6ggVEDDM7tY9n5FV2eJ4E
p0ImKfypi2aZYROxOK+v5x9ryFRMl4y07lMDvmtcV45uXYmfGNCgG9PNf91Kk/mh
JxEQbxycJwFoSi9XWljR8ahPdO11LXG7Dsj/RVbY8k2LdKNstl6Ae3aCpbe9u2Pj
vxYs1bVJuQIDAQABo1MwUTAdBgNVHQ4EFgQU+mRVN5HYJWgnpopReaLhf2cMcoYw
HwYDVR0jBBgwFoAU+mRVN5HYJWgnpopReaLhf2cMcoYwDwYDVR0TAQH/BAUwAwEB
/zANBgkqhkiG9w0BAQsFAAOCAQEABJpHVuc9tGhD04infRVlofvqXIUizTlOrjZX
vozW9pIhSWEHX8o+sJP8AMZLnrsdq+bm0HE0HvgYrw7Lb8pd4FpR46TkFHjeukoj
izqfgckjIBl2nwPGlynbKA0/U/rTCSxVt7XiAn+lgYUGIpOzNdk06/hRMitrMNB7
t2C97NseVC4b1ZgyFrozsefCfUmD8IJF0+XJ4Wzmsh0jRrI8koCtVmPYnKn6vw1b
cZprg/97CWHYrsavd406wOB60CMtYl83Q16ucOF1dretDFqJC5kY+aFLvuqfag2+
kIaoPV1MnGsxveQyyHdOsEatS5XOv/1OWcmnvePDPxcvb9jCcw==
-----END CERTIFICATE-----
`)
	key := []byte(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDRBi4mxrIgtEiF
VPsbUOWK1VJee3WLDeARcll97zUSDPTp9Ru9cXpijUsmbBimp7+ub6F/gU1Rvx7Y
74sfYWJ4D0WRV7g/aAo7PDzsD8rG+lXyy2A2CXguWpEKiRtzUDN4T82IuvlNxMNF
ebcNRxjFfGKemu0O7aN+0XKxo5mY0C6S8T+IvqpLdfqCBUQMMzu1j2fkVXZ4ngSn
QiYp/KmLZplhE7E4r6/nH2vIVEyXjLTuUwO+a1xXjm5diZ8Y0KAb081/3UqT+aEn
ERBvHJwnAWhKL1daWNHxqE907XUtcbsOyP9FVtjyTYt0o2y2XoB7doKlt727Y+O/
FizVtUm5AgMBAAECggEACak+l5f6Onj+u5vrjc4JyAaXW6ra6loSM9g8Uu3sHukW
plwoA7Pzp0u20CAxrP1Gpqw984/hSCCcb0Q2ItWMWLaC/YZni5W2WFnOyo3pzlPa
hmH4UNMT+ReCSfF/oW8w69QLcNEMjhfEu0i2iWBygIlA4SoRwC2Db6yEX7nLMwUB
6AICid9hfeACNRz/nq5ytdcHdmcB7Ptgb9jLiXr6RZw26g5AsRPHU3LdcyZAOXjP
aUHriHuHQFKAVkoEUxslvCB6ePCTCpB0bSAuzQbeGoY8fmvmNSCvJ1vrH5hiSUYp
Axtl5iNgFl5o9obb0eBYlY9x3pMSz0twdbCwfR7HAQKBgQDtWhmFm0NaJALoY+tq
lIIC0EOMSrcRIlgeXr6+g8womuDOMi5m/Nr5Mqt4mPOdP4HytrQb+a/ZmEm17KHh
mQb1vwH8ffirCBHbPNC1vwSNoxDKv9E6OysWlKiOzxPFSVZr3dKl2EMX6qi17n0l
LBrGXXaNPgYiHSmwBA5CZvvouQKBgQDhclGJfZfuoubQkUuz8yOA2uxalh/iUmQ/
G8ac6/w7dmnL9pXehqCWh06SeC3ZvW7yrf7IIGx4sTJji2FzQ+8Ta6pPELMyBEXr
1VirIFrlNVMlMQEbZcbzdzEhchM1RUpZJtl3b4amvH21UcRB69d9klcDRisKoFRm
k0P9QLHpAQKBgQDh5J9nphZa4u0ViYtTW1XFIbs3+R/0IbCl7tww67TRbF3KQL4i
7EHna88ALumkXf3qJvKRsXgoaqS0jSqgUAjst8ZHLQkOldaQxneIkezedDSWEisp
9YgTrJYjnHefiyXB8VL63jE0wPOiewEF8Mzmv6sFz+L8cq7rQ2Di16qmmQKBgQDH
bvCwVxkrMpJK2O2GH8U9fOzu6bUE6eviY/jb4mp8U7EdjGJhuuieoM2iBoxQ/SID
rmYftYcfcWlo4+juJZ99p5W+YcCTs3IDQPUyVOnzr6uA0Avxp6RKxhsBQj+5tTUj
Dpn77P3JzB7MYqvhwPcdD3LH46+5s8FWCFpx02RPAQKBgARbngtggfifatcsMC7n
lSv/FVLH7LYQAHdoW/EH5Be7FeeP+eQvGXwh1dgl+u0VZO8FvI8RwFganpBRR2Nc
ZSBRIb0fSUlTvIsckSWjpEvUJUomJXyi4PIZAfNvd9/u1uLInQiCDtObwb6hnLTU
FHHEZ+dR4eMaJp6PhNm8hu2O
-----END PRIVATE KEY-----
`)

	kp, err := tls.X509KeyPair(cert, key)
	if err != nil {
		panic(err)
	}
	kp.Leaf, err = x509.ParseCertificate(kp.Certificate[0])
	if err != nil {
		panic(err)
	}
	return kp
}()

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", samlsp.AttributeFromContext(r.Context(), "UserName"))
}

func main() {
	apiURL := os.Getenv("API_URL")
	pat := readPAT(os.Getenv("PAT_FILE"))
	domain := os.Getenv("API_DOMAIN")
	loginURL := os.Getenv("LOGIN_URL")
	idpURL := os.Getenv("IDP_URL")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	idpMetadataURL, err := url.Parse(idpURL)
	if err != nil {
		panic(err)
	}
	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient,
		*idpMetadataURL)
	if err != nil {
		panic(fmt.Errorf("failed to fetch IDP metadata from %s: %w", idpURL, err))
	}
	fmt.Printf("idpMetadata: %+v\n", idpMetadata)
	rootURL, err := url.Parse(host + ":" + port)
	if err != nil {
		panic(err)
	}

	samlSP, err := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata,
	})
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr: ":" + port,
	}
	app := http.HandlerFunc(hello)
	http.Handle("/hello", samlSP.RequireAccount(app))
	http.Handle("/saml/", samlSP)
	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	metadata, err := xml.MarshalIndent(samlSP.ServiceProvider.Metadata(), "", "  ")
	if err != nil {
		panic(err)
	}
	if err := createZitadelResources(apiURL, pat, domain, metadata, loginURL); err != nil {
		panic(err)
	}

	http.Handle("/healthy", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { return }))
	fmt.Println("/healthy returns 200 OK")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
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

func createZitadelResources(apiURL, pat, domain string, metadata []byte, loginURL string) error {
	projectID, err := CreateProject(apiURL, pat, domain)
	if err != nil {
		return err
	}
	return CreateApp(apiURL, pat, domain, projectID, metadata, loginURL)
}

type project struct {
	ID string `json:"id"`
}
type createProject struct {
	Name                   string `json:"name"`
	ProjectRoleAssertion   bool   `json:"projectRoleAssertion"`
	ProjectRoleCheck       bool   `json:"projectRoleCheck"`
	HasProjectCheck        bool   `json:"hasProjectCheck"`
	PrivateLabelingSetting string `json:"privateLabelingSetting"`
}

func CreateProject(apiURL, pat, domain string) (string, error) {
	createProject := &createProject{
		Name:                   "SAML",
		ProjectRoleAssertion:   false,
		ProjectRoleCheck:       false,
		HasProjectCheck:        false,
		PrivateLabelingSetting: "PRIVATE_LABELING_SETTING_UNSPECIFIED",
	}
	resp, err := doRequestWithHeaders(apiURL+"/management/v1/projects", pat, domain, createProject)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	p := new(project)
	if err := json.Unmarshal(data, p); err != nil {
		return "", err
	}
	return p.ID, nil
}

type createApp struct {
	Name         string  `json:"name"`
	MetadataXml  string  `json:"metadataXml"`
	LoginVersion version `json:"loginVersion"`
}
type version struct {
	LoginV2 loginV2 `json:"loginV2"`
}
type loginV2 struct {
	BaseUri string `json:"baseUri"`
}

func CreateApp(apiURL, pat, domain, projectID string, spMetadata []byte, loginURL string) error {
	encoded := make([]byte, base64.URLEncoding.EncodedLen(len(spMetadata)))
	base64.URLEncoding.Encode(encoded, spMetadata)

	createApp := &createApp{
		Name:        "SAML",
		MetadataXml: string(encoded),
		LoginVersion: version{
			LoginV2: loginV2{
				BaseUri: loginURL,
			},
		},
	}
	_, err := doRequestWithHeaders(apiURL+"/management/v1/projects/"+projectID+"/apps/saml", pat, domain, createApp)
	if err != nil {
		return fmt.Errorf("error creating saml app with request %+v: %v", *createApp, err)
	}
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
