// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
)

// Error defines model for error.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Healthcheck defines model for healthcheck.
type Healthcheck = string

// Item defines model for item.
type Item struct {
	CreatedAt     *time.Time `json:"created_at,omitempty"`
	Id            *int64     `json:"id,omitempty"`
	Name          string     `json:"name"`
	Quantity      int32      `json:"quantity"`
	ReservationId *int64     `json:"reservation_id,omitempty"`
	ReservedAt    *time.Time `json:"reserved_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty"`
}

// CreateItemJSONBody defines parameters for CreateItem.
type CreateItemJSONBody = Item

// CreateItemJSONRequestBody defines body for CreateItem for application/json ContentType.
type CreateItemJSONRequestBody = CreateItemJSONBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// IsAlive request
	IsAlive(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// ListItems request
	ListItems(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// CreateItem request with any body
	CreateItemWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	CreateItem(ctx context.Context, body CreateItemJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// IsReady request
	IsReady(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) IsAlive(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewIsAliveRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) ListItems(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewListItemsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateItemWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateItemRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) CreateItem(ctx context.Context, body CreateItemJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewCreateItemRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) IsReady(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewIsReadyRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewIsAliveRequest generates requests for IsAlive
func NewIsAliveRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/alive")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewListItemsRequest generates requests for ListItems
func NewListItemsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/items")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewCreateItemRequest calls the generic CreateItem builder with application/json body
func NewCreateItemRequest(server string, body CreateItemJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewCreateItemRequestWithBody(server, "application/json", bodyReader)
}

// NewCreateItemRequestWithBody generates requests for CreateItem with any type of body
func NewCreateItemRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/items")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewIsReadyRequest generates requests for IsReady
func NewIsReadyRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/ready")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// IsAlive request
	IsAliveWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*IsAliveResponse, error)

	// ListItems request
	ListItemsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListItemsResponse, error)

	// CreateItem request with any body
	CreateItemWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateItemResponse, error)

	CreateItemWithResponse(ctx context.Context, body CreateItemJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateItemResponse, error)

	// IsReady request
	IsReadyWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*IsReadyResponse, error)
}

type IsAliveResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Healthcheck
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r IsAliveResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r IsAliveResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type ListItemsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *[]Item
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r ListItemsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r ListItemsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type CreateItemResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON201      *Item
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r CreateItemResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r CreateItemResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type IsReadyResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *Healthcheck
	JSONDefault  *Error
}

// Status returns HTTPResponse.Status
func (r IsReadyResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r IsReadyResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// IsAliveWithResponse request returning *IsAliveResponse
func (c *ClientWithResponses) IsAliveWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*IsAliveResponse, error) {
	rsp, err := c.IsAlive(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseIsAliveResponse(rsp)
}

// ListItemsWithResponse request returning *ListItemsResponse
func (c *ClientWithResponses) ListItemsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*ListItemsResponse, error) {
	rsp, err := c.ListItems(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseListItemsResponse(rsp)
}

// CreateItemWithBodyWithResponse request with arbitrary body returning *CreateItemResponse
func (c *ClientWithResponses) CreateItemWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*CreateItemResponse, error) {
	rsp, err := c.CreateItemWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateItemResponse(rsp)
}

func (c *ClientWithResponses) CreateItemWithResponse(ctx context.Context, body CreateItemJSONRequestBody, reqEditors ...RequestEditorFn) (*CreateItemResponse, error) {
	rsp, err := c.CreateItem(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseCreateItemResponse(rsp)
}

// IsReadyWithResponse request returning *IsReadyResponse
func (c *ClientWithResponses) IsReadyWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*IsReadyResponse, error) {
	rsp, err := c.IsReady(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseIsReadyResponse(rsp)
}

// ParseIsAliveResponse parses an HTTP response from a IsAliveWithResponse call
func ParseIsAliveResponse(rsp *http.Response) (*IsAliveResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &IsAliveResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Healthcheck
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseListItemsResponse parses an HTTP response from a ListItemsWithResponse call
func ParseListItemsResponse(rsp *http.Response) (*ListItemsResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &ListItemsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest []Item
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseCreateItemResponse parses an HTTP response from a CreateItemWithResponse call
func ParseCreateItemResponse(rsp *http.Response) (*CreateItemResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &CreateItemResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 201:
		var dest Item
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON201 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ParseIsReadyResponse parses an HTTP response from a IsReadyWithResponse call
func ParseIsReadyResponse(rsp *http.Response) (*IsReadyResponse, error) {
	bodyBytes, err := ioutil.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &IsReadyResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest Healthcheck
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && true:
		var dest Error
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSONDefault = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /alive)
	IsAlive(w http.ResponseWriter, r *http.Request)

	// (GET /items)
	ListItems(w http.ResponseWriter, r *http.Request)

	// (POST /items)
	CreateItem(w http.ResponseWriter, r *http.Request)

	// (GET /ready)
	IsReady(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

// IsAlive operation middleware
func (siw *ServerInterfaceWrapper) IsAlive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.IsAlive(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// ListItems operation middleware
func (siw *ServerInterfaceWrapper) ListItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListItems(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// CreateItem operation middleware
func (siw *ServerInterfaceWrapper) CreateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateItem(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

// IsReady operation middleware
func (siw *ServerInterfaceWrapper) IsReady(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.IsReady(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/alive", wrapper.IsAlive)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/items", wrapper.ListItems)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/items", wrapper.CreateItem)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/ready", wrapper.IsReady)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xVQW/bPAz9Kwa/7+jFbjvs4Fu3U4ACK3YdgkKzaVudLakSHcwo/N8HUknmZG4PS7Du",
	"FCkUycfHR/oZSts7a9BQgOIZQtlir+SI3lvPB+etQ08a5e/SVsi/tfW9IihAG4IUaHQYL9ighymFHkNQ",
	"jbzdGQN5bRqYphQ8Pg3aYwXF18PDNMbeHILZb49YEsdqUXXUli2W3xfipaAJ+wWoHhVh9cAwZ4ArRfiO",
	"dM8ZParqs+lGKMgPmC5Erk6L/fD+Zb9Z+Ub1UnuvftyhaaiF4irPU+i1OdwX8j0NypCm8TTrzTWIr+6H",
	"fu45y+gxoN8q0tY8/DHsGORM0gZXnUn8iUSEzBk5v2uEPbSpLSesMJReOyYCCrjvFDGEpLY+EUlo0yTK",
	"VEmng5xZPYE1rKnjmOvdfYs+xBhXq3yVc2XWoVFOQwE3q3x1Ayk4Ra2oLVOd3krLG5S6WYrSjXXFQcOt",
	"2IViZ02IGr3O8zhVhtCIm3Ku06U4Zo+B0+/Hkk//e6yhgP+yX3Ob7YY2m0+JEHJMxO39OtEhiTDFWquh",
	"o4uljwtjIfHekAKpJnA/50g3bMhiD14i704H2nflLPoOaV4rRLbJdJCY8l6NS3WxfhJbJ7PB22npb7Lb",
	"oEGvy0QeJHt2juiOoDZTCs6GBXo/yaJkgiHOHQb6aKvxYuAjo1Mc66P2XV08xyk/p91Jdp+FNx6BfU9Y",
	"/LwMx1c2xxexv/nmiDD/qc0xTT8DAAD//1IjtmTCCAAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
