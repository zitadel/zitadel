package scim

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/zitadel/logging"
	"google.golang.org/grpc/metadata"

	zhttp "github.com/zitadel/zitadel/internal/api/http"
	"github.com/zitadel/zitadel/internal/api/scim/middleware"
	"github.com/zitadel/zitadel/internal/api/scim/resources"
	"github.com/zitadel/zitadel/internal/api/scim/schemas"
)

type Client struct {
	client  *http.Client
	baseURL string
	Users   *ResourceClient[resources.ScimUser]
}

type ResourceClient[T any] struct {
	client       *http.Client
	baseURL      string
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

type ListRequest struct {
	Schemas []schemas.ScimSchemaType `json:"schemas"`

	Count *int `json:"count,omitempty"`

	// StartIndex An integer indicating the 1-based index of the first query result.
	StartIndex *int `json:"startIndex,omitempty"`

	// Filter a scim filter expression to filter the query result.
	Filter *string `json:"filter,omitempty"`

	SortBy    *string               `json:"sortBy,omitempty"`
	SortOrder *ListRequestSortOrder `json:"sortOrder,omitempty"`

	SendAsPost bool
}

type ListRequestSortOrder string

const (
	ListRequestSortOrderAsc ListRequestSortOrder = "ascending"
	ListRequestSortOrderDsc ListRequestSortOrder = "descending"
)

type ListResponse[T any] struct {
	Schemas      []schemas.ScimSchemaType `json:"schemas"`
	ItemsPerPage int                      `json:"itemsPerPage"`
	TotalResults int                      `json:"totalResults"`
	StartIndex   int                      `json:"startIndex"`
	Resources    []T                      `json:"Resources"`
}

type BulkRequest struct {
	Schemas      []schemas.ScimSchemaType `json:"schemas"`
	FailOnErrors *int                     `json:"failOnErrors"`
	Operations   []*BulkRequestOperation  `json:"Operations"`
}

type BulkRequestOperation struct {
	Method string          `json:"method"`
	BulkID string          `json:"bulkId"`
	Path   string          `json:"path"`
	Data   json.RawMessage `json:"data"`
}

type BulkResponse struct {
	Schemas    []schemas.ScimSchemaType `json:"schemas"`
	Operations []*BulkResponseOperation `json:"Operations"`
}

type BulkResponseOperation struct {
	Method   string     `json:"method"`
	BulkID   string     `json:"bulkId,omitempty"`
	Location string     `json:"location,omitempty"`
	Response *ScimError `json:"response,omitempty"`
	Status   string     `json:"status"`
}

const (
	listQueryParamSortBy     = "sortBy"
	listQueryParamSortOrder  = "sortOrder"
	listQueryParamCount      = "count"
	listQueryParamStartIndex = "startIndex"
	listQueryParamFilter     = "filter"
)

func NewScimClient(target string) *Client {
	target = "http://" + target + schemas.HandlerPrefix
	client := &http.Client{}
	return &Client{
		client:  client,
		baseURL: target,
		Users: &ResourceClient[resources.ScimUser]{
			client:       client,
			baseURL:      target,
			resourceName: "Users",
		},
	}
}

func (c *Client) GetServiceProviderConfig(ctx context.Context, orgID string) ([]byte, error) {
	return c.getWithRawResponse(ctx, orgID, "/ServiceProviderConfig")
}

func (c *Client) GetSchemas(ctx context.Context, orgID string) ([]byte, error) {
	return c.getWithRawResponse(ctx, orgID, "/Schemas")
}

func (c *Client) GetSchema(ctx context.Context, orgID, schemaID string) ([]byte, error) {
	return c.getWithRawResponse(ctx, orgID, "/Schemas/"+schemaID)
}

func (c *Client) GetResourceTypes(ctx context.Context, orgID string) ([]byte, error) {
	return c.getWithRawResponse(ctx, orgID, "/ResourceTypes")
}

func (c *Client) GetResourceType(ctx context.Context, orgID, name string) ([]byte, error) {
	return c.getWithRawResponse(ctx, orgID, "/ResourceTypes/"+name)
}

func (c *Client) Bulk(ctx context.Context, orgID string, body []byte) (*BulkResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+orgID+"/Bulk", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set(zhttp.ContentType, middleware.ContentTypeScim)
	resp := new(BulkResponse)
	return resp, doReq(c.client, req, resp)
}

func (c *ResourceClient[T]) Create(ctx context.Context, orgID string, body []byte) (*T, error) {
	return c.doWithBody(ctx, http.MethodPost, orgID, "", bytes.NewReader(body))
}

func (c *ResourceClient[T]) Replace(ctx context.Context, orgID, id string, body []byte) (*T, error) {
	return c.doWithBody(ctx, http.MethodPut, orgID, id, bytes.NewReader(body))
}

func (c *ResourceClient[T]) Update(ctx context.Context, orgID, id string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.buildResourceURL(orgID, id), bytes.NewReader(body))
	if err != nil {
		return err
	}

	return doReq(c.client, req, nil)
}

func (c *ResourceClient[T]) List(ctx context.Context, orgID string, req *ListRequest) (*ListResponse[*T], error) {
	listResponse := new(ListResponse[*T])

	if req.SendAsPost {
		listReq, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}

		err = c.doWithResponse(ctx, http.MethodPost, orgID, ".search", bytes.NewReader(listReq), listResponse)
		return listResponse, err
	}

	query, err := url.ParseQuery("")
	if err != nil {
		return nil, err
	}

	if req.SortBy != nil {
		query.Set(listQueryParamSortBy, *req.SortBy)
	}

	if req.SortOrder != nil {
		query.Set(listQueryParamSortOrder, string(*req.SortOrder))
	}

	if req.Count != nil {
		query.Set(listQueryParamCount, strconv.Itoa(*req.Count))
	}

	if req.StartIndex != nil {
		query.Set(listQueryParamStartIndex, strconv.Itoa(*req.StartIndex))
	}

	if req.Filter != nil {
		query.Set(listQueryParamFilter, *req.Filter)
	}

	err = c.doWithResponse(ctx, http.MethodGet, orgID, "?"+query.Encode(), nil, listResponse)
	return listResponse, err
}

func (c *ResourceClient[T]) Get(ctx context.Context, orgID, resourceID string) (*T, error) {
	return c.doWithBody(ctx, http.MethodGet, orgID, resourceID, nil)
}

func (c *ResourceClient[T]) Delete(ctx context.Context, orgID, id string) error {
	return c.do(ctx, http.MethodDelete, orgID, id)
}

func (c *ResourceClient[T]) do(ctx context.Context, method, orgID, url string) error {
	req, err := http.NewRequestWithContext(ctx, method, c.buildResourceURL(orgID, url), nil)
	if err != nil {
		return err
	}

	return doReq(c.client, req, nil)
}

func (c *ResourceClient[T]) doWithResponse(ctx context.Context, method, orgID, url string, body io.Reader, response interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, c.buildResourceURL(orgID, url), body)
	if err != nil {
		return err
	}

	req.Header.Set(zhttp.ContentType, middleware.ContentTypeScim)
	return doReq(c.client, req, response)
}

func (c *ResourceClient[T]) doWithBody(ctx context.Context, method, orgID, url string, body io.Reader) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.buildResourceURL(orgID, url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set(zhttp.ContentType, middleware.ContentTypeScim)
	responseEntity := new(T)
	return responseEntity, doReq(c.client, req, responseEntity)
}

func (c *Client) getWithRawResponse(ctx context.Context, orgID, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/"+orgID+url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := resp.Body.Close()
		logging.OnError(err).Error("Failed to close response body")
	}()

	if (resp.StatusCode / 100) != 2 {
		return nil, readScimError(resp)
	}

	return io.ReadAll(resp.Body)
}

func doReq(client *http.Client, req *http.Request, responseEntity interface{}) error {
	addTokenAsHeader(req)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	defer func() {
		err := resp.Body.Close()
		logging.OnError(err).Error("Failed to close response body")
	}()

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

func (c *ResourceClient[T]) buildResourceURL(orgID, segment string) string {
	if segment == "" || strings.HasPrefix(segment, "?") {
		return c.baseURL + "/" + path.Join(orgID, c.resourceName) + segment
	}

	return c.baseURL + "/" + path.Join(orgID, c.resourceName, segment)
}

func (err *ScimError) Error() string {
	return "scim error: " + err.Detail
}
