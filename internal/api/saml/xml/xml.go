package xml

import (
	"bufio"
	"bytes"
	"compress/flate"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/caos/zitadel/internal/api/saml/xml/samlp"
	"github.com/caos/zitadel/internal/api/saml/xml/soap"
	"github.com/caos/zitadel/internal/api/saml/xml/xml_dsig"
	"net/http"
	"strings"
)

const (
	EncodingDeflate = "urn:oasis:names:tc:SAML:2.0:bindings:URL-Encoding:DEFLATE"
)

func Marshal(data interface{}) (string, error) {
	var xmlbuff bytes.Buffer

	memWriter := bufio.NewWriter(&xmlbuff)
	_, err := memWriter.Write([]byte(xml.Header))
	if err != nil {
		return "", err
	}

	encoder := xml.NewEncoder(memWriter)
	err = encoder.Encode(data)
	if err != nil {
		return "", err
	}

	err = memWriter.Flush()
	if err != nil {
		return "", err
	}

	return xmlbuff.String(), nil
}

func DeflateAndBase64(data []byte) ([]byte, error) {
	b := &bytes.Buffer{}
	w1 := base64.NewEncoder(base64.StdEncoding, b)
	defer w1.Close()

	w2, _ := flate.NewWriter(w1, 1)
	defer w2.Close()

	bw := bufio.NewWriter(w1)
	if _, err := bw.Write(data); err != nil {
		return nil, err
	}
	if err := bw.Flush(); err != nil {
		return nil, err
	}

	if err := w2.Flush(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func WriteXMLMarshalled(w http.ResponseWriter, body interface{}) error {
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

func Write(w http.ResponseWriter, body []byte) error {
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}

	_, err = w.Write(body)
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

func DecodeSignature(encoding string, message string) (*xml_dsig.SignatureType, error) {
	retBytes := []byte(message)
	/*retBytes, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode: %w", err)
	}*/

	ret := &xml_dsig.SignatureType{}
	switch encoding {
	case EncodingDeflate:
		reader := flate.NewReader(bytes.NewReader(retBytes))
		decoder := xml.NewDecoder(reader)
		if err := decoder.Decode(ret); err != nil {
			return nil, fmt.Errorf("failed to defalte decode: %w", err)
		}
	default:
		reader := flate.NewReader(bytes.NewReader(retBytes))
		decoder := xml.NewDecoder(reader)
		if err := decoder.Decode(ret); err != nil {
			if err := xml.Unmarshal(retBytes, ret); err != nil {
				return nil, fmt.Errorf("failed to unmarshal: %w", err)
			}
		}
	}

	return ret, nil
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
