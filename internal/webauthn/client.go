package webauthn

import (
	"fmt"

	"github.com/descope/virtualwebauthn"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

type Client struct {
	rp             virtualwebauthn.RelyingParty
	auth           virtualwebauthn.Authenticator
	authVerifyUser virtualwebauthn.Authenticator
	credential     virtualwebauthn.Credential
}

func NewClient(name, domain, origin string) *Client {
	rp := virtualwebauthn.RelyingParty{
		Name:   name,
		ID:     domain,
		Origin: origin,
	}
	return &Client{
		rp: rp,
		auth: virtualwebauthn.NewAuthenticatorWithOptions(virtualwebauthn.AuthenticatorOptions{
			UserNotVerified: true,
		}),
		authVerifyUser: virtualwebauthn.NewAuthenticator(),
		credential:     virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2),
	}
}

func (c *Client) CreateAttestationResponse(optionsPb *structpb.Struct) (*structpb.Struct, error) {
	options, err := protojson.Marshal(optionsPb)
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAttestationResponse: %w", err)
	}
	attestationResponse, err := c.CreateAttestationResponseData(options)
	if err != nil {
		return nil, err
	}
	resp := new(structpb.Struct)
	err = protojson.Unmarshal(attestationResponse, resp)
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAttestationResponse: %w", err)
	}
	return resp, nil
}

func (c *Client) CreateAttestationResponseData(options []byte) ([]byte, error) {
	parsedAttestationOptions, err := virtualwebauthn.ParseAttestationOptions(string(options))
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAttestationResponse: %w", err)
	}
	return []byte(virtualwebauthn.CreateAttestationResponse(
		c.rp, c.auth, c.credential, *parsedAttestationOptions,
	)), nil
}

func (c *Client) CreateAssertionResponse(optionsPb *structpb.Struct, verifyUser bool) (*structpb.Struct, error) {
	options, err := protojson.Marshal(optionsPb)
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAssertionResponse: %w", err)
	}
	parsedAssertionOptions, err := virtualwebauthn.ParseAssertionOptions(string(options))
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAssertionResponse: %w", err)
	}
	authenticator := c.auth
	if verifyUser {
		authenticator = c.authVerifyUser
	}
	resp := new(structpb.Struct)
	err = protojson.Unmarshal([]byte(virtualwebauthn.CreateAssertionResponse(
		c.rp, authenticator, c.credential, *parsedAssertionOptions,
	)), resp)
	if err != nil {
		return nil, fmt.Errorf("webauthn.Client.CreateAssertionResponse: %w", err)
	}
	return resp, nil
}
