package saml

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/saml/xml/protocol/samlp"
	"net/http"
	"time"
)

func (p *IdentityProvider) logoutHandleFunc(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		logging.Log("SAML-91pokk").Error(err)
		http.Error(w, fmt.Errorf("failed to parse form: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	request := r.Form.Get("SAMLRequest")
	encoding := r.Form.Get("SAMLEncoding")
	relayState := r.Form.Get("RelayState")

	logoutRequest, err := decodeLogoutRequest(encoding, request)
	if err != nil {
		logging.Log("SAML-892umn").Error(err)
		if err := sendBackLogoutResponse(
			p.logoutTemplate,
			w,
			relayState,
			"",
			makeUnsupportedlLogoutResponse(
				logoutRequest,
				"",
				p.EntityID,
				fmt.Errorf("failed to decode request: %w", err).Error(),
			),
		); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
		return
	}

	if err := validateLogoutRequest(logoutRequest); err != nil {
		logging.Log("SAML-892u3n").Error(err)
		if err := sendBackLogoutResponse(
			p.logoutTemplate,
			w,
			relayState,
			"",
			makeDeniedLogoutResponse(
				logoutRequest,
				"",
				p.EntityID,
				fmt.Errorf("failed to validate request: %w", err).Error(),
			),
		); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
		return
	}

	sp, err := p.storage.GetEntityByID(r.Context(), logoutRequest.Issuer.Text)
	if err != nil {
		logging.Log("SAML-292u3n").Error(err)
		if err := sendBackLogoutResponse(
			p.logoutTemplate,
			w,
			relayState,
			"",
			makeDeniedLogoutResponse(
				logoutRequest,
				"",
				p.EntityID,
				fmt.Errorf("unknown service provider: %w", err).Error(),
			),
		); err != nil {
			http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		}
		return
	}

	logoutURL := ""
	if sp.metadata.SPSSODescriptor.SingleLogoutService != nil {
		for _, url := range sp.metadata.SPSSODescriptor.SingleLogoutService {
			logoutURL = url.Location
			break
		}
	}

	if err := sendBackLogoutResponse(
		p.logoutTemplate,
		w,
		relayState,
		logoutURL,
		makeSuccessfulLogoutResponse(
			logoutRequest,
			logoutURL,
			p.EntityID,
		),
	); err != nil {
		logging.Log("SAML-846u3n").Error(err)
		http.Error(w, fmt.Errorf("failed to send response: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	logging.Log("SAML-892u3n").Info(fmt.Sprintf("logout request for user %s", logoutRequest.NameID.Text))
}

func decodeLogoutRequest(encoding string, message string) (*samlp.LogoutRequest, error) {
	reqBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil, err
	}

	req := &samlp.LogoutRequest{}
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

func validateLogoutRequest(request *samlp.LogoutRequest) error {
	now := time.Now().UTC()
	if request.NotOnOrAfter != "" {
		//TODO
		t, err := time.Parse("", request.NotOnOrAfter)
		if err != nil {
			return fmt.Errorf("failed to parse NotOnOrAfter: %w", err)
		}
		if t.After(now) {
			return fmt.Errorf("on or after time given by NotOnOrAfter")
		}
	}
	if request.NameID == nil || request.NameID.Text == "" {
		return fmt.Errorf("no nameID provided")
	}
	return nil
}
