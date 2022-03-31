package saml

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/saml/checker"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
	"net/http"
	"time"
)

type LogoutRequestForm struct {
	LogoutRequest string
	Encoding      string
	RelayState    string
}

func (p *IdentityProvider) logoutHandleFunc(w http.ResponseWriter, r *http.Request) {
	checker := checker.Checker{}
	var logoutRequestForm *LogoutRequestForm
	var logoutRequest *samlp.LogoutRequestType
	var err error
	var sp *ServiceProvider
	response := &LogoutResponse{
		LogoutTemplate: p.logoutTemplate,
		ErrorFunc: func(err error) {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		},
		Issuer: p.EntityID,
	}

	// parse from to get logout request
	checker.WithLogicStep(
		func() error {
			logoutRequestForm, err = getLogoutRequestFromRequest(r)
			if err != nil {
				return err
			}
			response.RelayState = logoutRequestForm.RelayState
			return nil
		},
		"SAML-91pokk",
		func() {
			http.Error(w, fmt.Errorf("failed to parse form: %w", err).Error(), http.StatusInternalServerError)
		},
	)

	//decode logout request to internal struct
	checker.WithLogicStep(
		func() error {
			logoutRequest, err = decodeLogoutRequest(logoutRequestForm.Encoding, logoutRequestForm.LogoutRequest)
			if err != nil {
				return err
			}
			response.RelayState = logoutRequestForm.RelayState
			response.RequestID = logoutRequest.Id
			return nil
		},
		"SAML-892umn",
		func() {
			response.sendBackLogoutResponse(w, response.makeUnsupportedlLogoutResponse(fmt.Errorf("failed to decode request: %w", err).Error()))
		},
	)

	//verify required data in request
	checker.WithLogicStep(
		func() error {
			now := time.Now().UTC()
			if logoutRequest.NotOnOrAfter != "" {
				//TODO
				t, err := time.Parse("", logoutRequest.NotOnOrAfter)
				if err != nil {
					return fmt.Errorf("failed to parse NotOnOrAfter: %w", err)
				}
				if t.After(now) {
					return fmt.Errorf("on or after time given by NotOnOrAfter")
				}
			}
			if logoutRequest.NameID == nil || logoutRequest.NameID.Text == "" {
				return fmt.Errorf("no nameID provided")
			}
			return nil
		},
		"SAML-892u3n",
		func() {
			response.sendBackLogoutResponse(w, response.makeDeniedLogoutResponse(fmt.Errorf("failed to validate request: %w", err).Error()))
		},
	)

	// get persisted service provider from issuer out of the request
	checker.WithLogicStep(
		func() error {
			sp, err = p.GetServiceProvider(r.Context(), logoutRequest.Issuer.Text)
			return err
		},
		" SAML-317s2s",
		func() {
			response.sendBackLogoutResponse(w, response.makeDeniedLogoutResponse(fmt.Errorf("failed to find registered serviceprovider: %w", err).Error()))
		},
	)

	// get logoutURL from provided service provider metadata
	checker.WithValueStep(
		func() {
			if sp.metadata.SPSSODescriptor.SingleLogoutService != nil {
				for _, url := range sp.metadata.SPSSODescriptor.SingleLogoutService {
					response.LogoutURL = url.Location
					break
				}
			}
		},
	)

	//check and log errors if necessary
	if checker.CheckFailed() {
		return
	}

	response.sendBackLogoutResponse(
		w,
		response.makeSuccessfulLogoutResponse(),
	)
	logging.Log("SAML-892u3n").Info(fmt.Sprintf("logout request for user %s", logoutRequest.NameID.Text))
}

func decodeLogoutRequest(encoding string, message string) (*samlp.LogoutRequestType, error) {
	reqBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil, err
	}

	req := &samlp.LogoutRequestType{}
	switch encoding {
	case "":
		reader := flate.NewReader(bytes.NewReader(reqBytes))
		decoder := xml.NewDecoder(reader)
		if err = decoder.Decode(req); err != nil {
			return nil, err
		}
	case EncodingDeflate:
		reader := flate.NewReader(bytes.NewReader(reqBytes))
		decoder := xml.NewDecoder(reader)
		if err = decoder.Decode(req); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown encoding")
	}

	return req, nil
}

func getLogoutRequestFromRequest(r *http.Request) (*LogoutRequestForm, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	request := &LogoutRequestForm{
		LogoutRequest: r.Form.Get("SAMLRequest"),
		Encoding:      r.Form.Get("SAMLEncoding"),
		RelayState:    r.Form.Get("RelayState"),
	}

	return request, nil
}
