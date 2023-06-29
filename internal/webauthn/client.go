package webauthn

import (
	"fmt"

	"github.com/descope/virtualwebauthn"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type Client struct {
	rp         virtualwebauthn.RelyingParty
	auth       virtualwebauthn.Authenticator
	credential virtualwebauthn.Credential
}

func NewClient(name, domain, origin string) *Client {
	rp := virtualwebauthn.RelyingParty{
		Name:   name,
		ID:     domain,
		Origin: origin,
	}
	return &Client{
		rp:         rp,
		auth:       virtualwebauthn.NewAuthenticator(),
		credential: virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2),
	}
}

func (c *Client) CreateAttestationResponse(optionsPb *structpb.Struct) (*structpb.Struct, error) {
	options, err := protojson.Marshal(optionsPb)
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAttestationResponse: %w", err)
	}
	parsedAttestationOptions, err := virtualwebauthn.ParseAttestationOptions(string(options))
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAttestationResponse: %w", err)
	}
	resp := new(structpb.Struct)
	err = protojson.Unmarshal([]byte(virtualwebauthn.CreateAttestationResponse(
		c.rp, c.auth, c.credential, *parsedAttestationOptions,
	)), resp)
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAttestationResponse: %w", err)
	}
	return resp, nil
}

func (c *Client) CreateAssertionResponse(optionsPb *structpb.Struct) (*structpb.Struct, error) {
	options, err := protojson.Marshal(optionsPb)
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAssertionResponse: %w", err)
	}
	parsedAssertionOptions, err := virtualwebauthn.ParseAssertionOptions(string(options))
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAssertionResponse: %w", err)
	}
	resp := new(structpb.Struct)
	err = protojson.Unmarshal([]byte(virtualwebauthn.CreateAssertionResponse(
		c.rp, c.auth, c.credential, *parsedAssertionOptions,
	)), resp)
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAssertionResponse: %w", err)
	}
	return resp, nil
}
