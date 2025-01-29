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
	Users *ResourceClient[resources.ScimUser]
}

type ResourceClient[T any] struct {
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

type ListRequest struct {
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
		Users: &ResourceClient[resources.ScimUser]{
			client:       client,
			baseUrl:      target,
			resourceName: "Users",
		},
	}
}

func (c *ResourceClient[T]) Create(ctx context.Context, orgID string, body []byte) (*T, error) {
	return c.doWithBody(ctx, http.MethodPost, orgID, "", bytes.NewReader(body))
}

func (c *ResourceClient[T]) Replace(ctx context.Context, orgID, id string, body []byte) (*T, error) {
	return c.doWithBody(ctx, http.MethodPut, orgID, id, bytes.NewReader(body))
}

func (c *ResourceClient[T]) Update(ctx context.Context, orgID, id string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, c.buildURL(orgID, id), bytes.NewReader(body))
	if err != nil {
		return err
	}

	return c.doReq(req, nil)
}

func (c *ResourceClient[T]) List(ctx context.Context, orgID string, req *ListRequest) (*ListResponse[*T], error) {
	if req.SendAsPost {
		listReq, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		return c.doWithListResponse(ctx, http.MethodPost, orgID, ".search", bytes.NewReader(listReq))
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

	return c.doWithListResponse(ctx, http.MethodGet, orgID, "?"+query.Encode(), nil)
}

func (c *ResourceClient[T]) Get(ctx context.Context, orgID, resourceID string) (*T, error) {
	return c.doWithBody(ctx, http.MethodGet, orgID, resourceID, nil)
}

func (c *ResourceClient[T]) Delete(ctx context.Context, orgID, id string) error {
	return c.do(ctx, http.MethodDelete, orgID, id)
}

func (c *ResourceClient[T]) do(ctx context.Context, method, orgID, url string) error {
	req, err := http.NewRequestWithContext(ctx, method, c.buildURL(orgID, url), nil)
	if err != nil {
		return err
	}

	return c.doReq(req, nil)
}

func (c *ResourceClient[T]) doWithListResponse(ctx context.Context, method, orgID, url string, body io.Reader) (*ListResponse[*T], error) {
	req, err := http.NewRequestWithContext(ctx, method, c.buildURL(orgID, url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set(zhttp.ContentType, middleware.ContentTypeScim)
	response := new(ListResponse[*T])
	return response, c.doReq(req, response)
}

func (c *ResourceClient[T]) doWithBody(ctx context.Context, method, orgID, url string, body io.Reader) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.buildURL(orgID, url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set(zhttp.ContentType, middleware.ContentTypeScim)
	responseEntity := new(T)
	return responseEntity, c.doReq(req, responseEntity)
}

func (c *ResourceClient[T]) doReq(req *http.Request, responseEntity interface{}) error {
	addTokenAsHeader(req)

	resp, err := c.client.Do(req)

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

func (c *ResourceClient[T]) buildURL(orgID, segment string) string {
	if segment == "" || strings.HasPrefix(segment, "?") {
		return c.baseUrl + "/" + path.Join(orgID, c.resourceName) + segment
	}

	return c.baseUrl + "/" + path.Join(orgID, c.resourceName, segment)
}

func (err *ScimError) Error() string {
	return "scim error: " + err.Detail
}
