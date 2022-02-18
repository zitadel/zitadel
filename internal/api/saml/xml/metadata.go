package xml

import (
	"encoding/base64"
	"encoding/xml"
	"github.com/caos/zitadel/internal/api/saml/xml/metadata/md"
	"io"
	"net/http"
)

func ReadMetadataFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyEncoded := base64.StdEncoding.EncodeToString(body)
	return []byte(bodyEncoded), nil
}

func ParseMetadataXmlIntoStruct(xmlData []byte) (*md.EntityDescriptor, error) {
	xmlDataDecoded := make([]byte, 0)
	xmlDataDecoded, err := base64.StdEncoding.DecodeString(string(xmlData))
	if err != nil {
		return nil, err
	}

	metadata := &md.EntityDescriptor{}
	if err := xml.Unmarshal(xmlDataDecoded, metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}
