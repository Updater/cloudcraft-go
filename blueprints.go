package cloudcraft

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

const blueprintBasePath = "blueprint"

// BlueprintsService is an interface for interfacing with the Blueprints
// endpoints of the Cloudcraft API
// See: https://developers.cloudcraft.co/#dbc3d135-6447-47f2-b043-bae65b722246
type BlueprintsService interface {
	List(context.Context) ([]Blueprint, *Response, error)
	Get(context.Context, string) (*Blueprint, *Response, error)
	Create(context.Context, *BlueprintCreateRequest) (*Blueprint, *Response, error)
	Update(context.Context, string, *BlueprintUpdateRequest) (*Blueprint, *Response, error)
	Delete(context.Context, string) (*Response, error)
	Export(context.Context, string, *BlueprintExportRequest) (*BlueprintImage, *Response, error)
}

// BlueprintsServiceOp handles communication with the Blueprint related methods of the
// Cloudcraft API.
type BlueprintsServiceOp struct {
	client *Client
}

var _ BlueprintsService = &BlueprintsServiceOp{}

type BlueprintData struct {
	Grid           string                   `json:"grid,omitempty"`
	LinkKey        string                   `json:"linkKey,omitempty"`
	Name           string                   `json:"name,omitempty"`
	Text           []map[string]interface{} `json:"text,omitempty"`
	Edges          []map[string]interface{} `json:"edges,omitempty"`
	Icons          []map[string]interface{} `json:"icons,omitempty"`
	Nodes          []map[string]interface{} `json:"nodes,omitempty"`
	Groups         []map[string]interface{} `json:"groups,omitempty"`
	Images         []map[string]interface{} `json:"images,omitempty"`
	Surfaces       []map[string]interface{} `json:"surfaces,omitempty"`
	Connectors     []map[string]interface{} `json:"connectors,omitempty"`
	DisabledLayers []map[string]interface{} `json:"disabledLayers,omitempty"`
}

// Blueprint represents a Cloudcraft Blueprint
type Blueprint struct {
	Id         string         `json:"id,omitempty"`
	Name       string         `json:"name,omitempty"`
	CreatedAt  time.Time      `json:"createdAt,omitempty"`
	UpdatedAt  time.Time      `json:"createdAt,omitempty"`
	CreatorId  string         `json:"CreatorId,omitempty"`
	LastUserId string         `json:"LastUserId,omitempty"`
	Data       *BlueprintData `json:data,omitempty`
}

type BlueprintExportParameters struct {
	Grid        bool    `url:"grid"`
	Height      int     `url:"height"`
	Landscape   bool    `url:"landscape"`
	PaperSize   string  `url:"paperSize"`
	Scale       float32 `url:"scale"`
	Transparent bool    `url:"transparent"`
	Width       int     `url:"width"`
}

type BlueprintImage struct {
	ContentType      string
	Content          *bytes.Buffer
	ExportParameters *BlueprintExportParameters
}

// Convert Blueprint to a string
func (d Blueprint) String() string {
	return Stringify(d)
}

type BlueprintsRoot struct {
	Blueprints []Blueprint `json:"blueprints"`
}

// BlueprintCreateRequest represents a request to create a Blueprint.
type BlueprintCreateRequest struct {
	Data *BlueprintData `json:"data"`
}

func (d BlueprintCreateRequest) String() string {
	return Stringify(d)
}

// BlueprintCreateRequest represents a request to update a Blueprint.
type BlueprintUpdateRequest struct {
	Data *BlueprintData `json:"data"`
}

func (d BlueprintUpdateRequest) String() string {
	return Stringify(d)
}

type BlueprintExportRequest struct {
	Format           string
	ExportParameters *BlueprintExportParameters
}

func (d BlueprintExportRequest) String() string {
	return Stringify(d)
}

// List all Blueprints.
func (s *BlueprintsServiceOp) List(ctx context.Context) ([]Blueprint, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, blueprintBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(BlueprintsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.Blueprints, resp, err
}

// Get individual Blueprint.
func (s *BlueprintsServiceOp) Get(ctx context.Context, blueprintId string) (*Blueprint, *Response, error) {
	if blueprintId == "" {
		return nil, nil, NewArgError("blueprintId", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", blueprintBasePath, blueprintId)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	blueprint := new(Blueprint)
	resp, err := s.client.Do(ctx, req, blueprint)
	if err != nil {
		return nil, resp, err
	}

	return blueprint, resp, err
}

// Create Blueprint
func (s *BlueprintsServiceOp) Create(ctx context.Context, createRequest *BlueprintCreateRequest) (*Blueprint, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := blueprintBasePath

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	blueprint := new(Blueprint)
	resp, err := s.client.Do(ctx, req, blueprint)
	if err != nil {
		return nil, resp, err
	}

	return blueprint, resp, err
}

// Update Blueprint
func (s *BlueprintsServiceOp) Update(ctx context.Context, blueprintId string, updateRequest *BlueprintUpdateRequest) (*Blueprint, *Response, error) {
	if blueprintId == "" {
		return nil, nil, NewArgError("blueprintId", "cannot be empty")
	}

	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	path := blueprintBasePath

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	blueprint := new(Blueprint)
	resp, err := s.client.Do(ctx, req, blueprint)
	if err != nil {
		return nil, resp, err
	}

	return blueprint, resp, err
}

// Delete Blueprint.
func (s *BlueprintsServiceOp) Delete(ctx context.Context, blueprintId string) (*Response, error) {
	if blueprintId == "" {
		return nil, NewArgError("blueprintId", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", blueprintBasePath, blueprintId)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

// var validPaperSizes = []string{"Letter", "Legal", "Tabloid", "Ledger", "A0", "A1", "A2", "A3", "A4", "A5"}

// imageMediaType = "image/svg+xml, image/png, application/pdf, application/xml, application/json"
// Export Blueprint.
func (s *BlueprintsServiceOp) Export(ctx context.Context, blueprintId string, exportRequest *BlueprintExportRequest) (*BlueprintImage, *Response, error) {
	if blueprintId == "" {
		return nil, nil, NewArgError("blueprintId", "cannot be empty")
	}

	path, err := addOptions(fmt.Sprintf("%s/%s/%s", blueprintBasePath, blueprintId, exportRequest.Format), exportRequest.ExportParameters)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	blueprintImage := new(BlueprintImage)
	resp, err := s.client.Do(ctx, req, blueprintImage)
	if err != nil {
		return nil, resp, err
	}

	return blueprintImage, resp, err
}
