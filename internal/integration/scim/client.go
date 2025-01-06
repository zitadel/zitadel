package scim

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path"

	"github.com/zitadel/logging"
	"google.golang.org/grpc/metadata"

	zhttp "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/scim/middleware"
	"github.com/zitadel/zitadel/internal/api/scim/resources"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
)

type Client struct {
	Users *ResourceClient
}

type ResourceClient struct {
	client       *http.Client
	baseUrl      string
	resourceName string
}

type ScimError struct {
	Schemas       []string            `json:"schemas"`
	ScimType      string              `json:"scimType"`
	Detail        string              `json:"detail"`
	Status        string              `json:"status"`
	ZitadelDetail *ZitadelErrorDetail `json:"urn:ietf:params:scim:api:zitadel:messages:2.0:ErrorDetail,omitempty"`
}

type ZitadelErrorDetail struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func NewScimClient(target string) *Client {
	target = "http://" + target + schemas.HandlerPrefix
	client := &http.Client{}
	return &Client{
		Users: &ResourceClient{
			client:       client,
			baseUrl:      target,
			resourceName: "Users",
		},
	}
}

func (c *ResourceClient) Create(ctx context.Context, orgID string, body []byte) (*resources.ScimUser, error) {
	user := new(resources.ScimUser)
	err := c.doWithBody(ctx, http.MethodPost, orgID, "", bytes.NewReader(body), user)
	return user, err
}

func (c *ResourceClient) doWithBody(ctx context.Context, method, orgID, url string, body io.Reader, responseEntity interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, c.buildURL(orgID, url), body)
	if err != nil {
		return err
	}

	req.Header.Set(zhttp.ContentType, middleware.ContentTypeScim)
	return c.doReq(req, responseEntity)
}

func (c *ResourceClient) doReq(req *http.Request, responseEntity interface{}) error {
	addTokenAsHeader(req)

	resp, err := c.client.Do(req)
	defer func() {
		err := resp.Body.Close()
		logging.OnError(err).Error("Failed to close response body")
	}()

	if err != nil {
		return err
	}

	if (resp.StatusCode / 100) != 2 {
		return readScimError(resp)
	}

	if responseEntity == nil {
		return nil
	}

	err = readJson(responseEntity, resp)
	return err
}

func addTokenAsHeader(req *http.Request) {
	md, ok := metadata.FromOutgoingContext(req.Context())
	if !ok {
		return
	}

	req.Header.Set("Authorization", md.Get("Authorization")[0])
}

func readJson(entity interface{}, resp *http.Response) error {
	defer func(body io.ReadCloser) {
		err := body.Close()
		logging.OnError(err).Panic("Failed to close response body")
	}(resp.Body)

	err := json.NewDecoder(resp.Body).Decode(entity)
	logging.OnError(err).Panic("Failed decoding entity")
	return err
}

func readScimError(resp *http.Response) error {
	scimErr := new(ScimError)
	readErr := readJson(scimErr, resp)
	logging.OnError(readErr).Panic("Failed reading scim error")
	return scimErr
}

func (c *ResourceClient) buildURL(orgID, segment string) string {
	if segment == "" {
		return c.baseUrl + "/" + path.Join(orgID, c.resourceName)
	}

	return c.baseUrl + "/" + path.Join(orgID, c.resourceName, segment)
}

func (err *ScimError) Error() string {
	return "scim error: " + err.Detail
}
