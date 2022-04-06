package xml

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
	"github.com/caos/zitadel/internal/api/saml/xml/soap"
	"net/http"
	"strings"
)

const (
	EncodingDeflate = "urn:oasis:names:tc:SAML:2.0:bindings:URL-Encoding:DEFLATE"
)

func WriteXML(w http.ResponseWriter, body interface{}) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	encoder := xml.NewEncoder(w)

	err = encoder.Encode(body)
	if err != nil {
		return err
	}
	err = encoder.Flush()
	return err
}

func DecodeAuthNRequest(encoding string, message string) (*samlp.AuthnRequestType, error) {
	reqBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode: %w", err)
	}

	req := &samlp.AuthnRequestType{}
	switch encoding {
	case EncodingDeflate:
		reader := flate.NewReader(bytes.NewReader(reqBytes))
		decoder := xml.NewDecoder(reader)
		if err = decoder.Decode(req); err != nil {
			return nil, fmt.Errorf("failed to defalte decode: %w", err)
		}
	default:
		reader := flate.NewReader(bytes.NewReader(reqBytes))
		decoder := xml.NewDecoder(reader)
		if err = decoder.Decode(req); err != nil {
			if err := xml.Unmarshal(reqBytes, req); err != nil {
				return nil, fmt.Errorf("failed to unmarshal: %w", err)
			}
		}
	}

	return req, nil
}

func DecodeAttributeQuery(request string) (*samlp.AttributeQueryType, error) {
	decoder := xml.NewDecoder(strings.NewReader(request))
	var attrEnv soap.AttributeQueryEnvelope
	err := decoder.Decode(&attrEnv)
	if err != nil {
		return nil, err
	}

	return attrEnv.Body.AttributeQuery, nil
}

func DecodeLogoutRequest(encoding string, message string) (*samlp.LogoutRequestType, error) {
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
