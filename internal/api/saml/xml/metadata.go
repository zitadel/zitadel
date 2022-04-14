package xml

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"github.com/caos/zitadel/internal/api/saml/xml/md"
	"io"
	"net/http"
)

func ReadMetadataFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error while reading metadata with statusCode: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	bodyEncoded := base64.StdEncoding.EncodeToString(body)
	return []byte(bodyEncoded), nil
}

func ParseMetadataXmlIntoStruct(xmlData []byte) (*md.EntityDescriptorType, error) {
	xmlDataDecoded := make([]byte, 0)
	xmlDataDecoded, err := base64.StdEncoding.DecodeString(string(xmlData))
	if err != nil {
		return nil, err
	}

	metadata := &md.EntityDescriptorType{}
	if err := xml.Unmarshal(xmlDataDecoded, metadata); err != nil {
		return nil, err
	}
	return metadata, nil
}

func GetCertsFromKeyDescriptors(keyDescs []md.KeyDescriptorType) []string {
	certStrs := []string{}
	if keyDescs == nil {
		return certStrs
	}
	for _, keyDescriptor := range keyDescs {
		for _, x509Data := range keyDescriptor.KeyInfo.X509Data {
			if len(x509Data.X509Certificate) != 0 {
				switch keyDescriptor.Use {
				case "", "signing":
					certStrs = append(certStrs, x509Data.X509Certificate)
				}
			}
		}
	}
	return certStrs
}
