package cloudcraft

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

const awsAccountBasePath = "aws/account"

// AwsAccountsService is an interface for interfacing with the AwsAccounts
// endpoints of the Cloudcraft API
// See: https://developers.cloudcraft.co/#dbc3d135-6447-47f2-b043-bae65b722246
type AwsAccountsService interface {
	List(context.Context) ([]AwsAccount, *Response, error)
	Get(context.Context, string) (*AwsAccount, *Response, error)
	Create(context.Context, *AwsAccountCreateOrUpdateRequest) (*AwsAccount, *Response, error)
	Update(context.Context, string, *AwsAccountCreateOrUpdateRequest) (*AwsAccount, *Response, error)
	Delete(context.Context, string) (*Response, error)
	Snapshot(context.Context, string, *AwsAccountSnapshotRequest) (*AwsAccountSnapshot, *Response, error)
	IamParameters(context.Context) (*AwsAccountIamParameters, *Response, error)
}

// AwsAccountsServiceOp handles communication with the AwsAccount related methods of the
// Cloudcraft API.
type AwsAccountsServiceOp struct {
	client *Client
}

var _ AwsAccountsService = &AwsAccountsServiceOp{}

type AwsAccountDataTextMapPos struct {
	RelTo  string `json:"relTo,omitempty"`
	offset []int  `json:"offset,omitempty"`
}

type AwsAccountDataText struct {
	Id        string                   `json:"id,omitempty"`
	Text      string                   `json:"text,omitempty"`
	Type      string                   `json:"type,omitempty"`
	Color     string                   `json:"color,omitempty"`
	TextSize  int                      `json:"textSize,omitempty"`
	Direction string                   `json:"direction,omitempty"`
	Isometric bool                     `json:"isometric,omitempty"`
	MapPos    AwsAccountDataTextMapPos `json:"mapPos,omitempty"`
}

type AwsAccountDataEdge struct {
	Id     string `json:"id,omitempty"`
	To     string `json:"to,omitempty"`
	From   string `json:"from,omitempty"`
	Type   string `json:"type,omitempty"`
	Color  string `json:"color,omitempty"`
	Width  int    `json:"width,omitempty"`
	Dashed bool   `json:"dashed,omitempty"`
}

type AwsAccountData struct {
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

// AwsAccount represents a Cloudcraft AwsAccount
type AwsAccount struct {
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	CreatorId  string    `json:"CreatorId,omitempty"`
	ExternalId string    `json:"externalId"`
	Id         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	RoleArn    string    `json:"roleArn,omitempty"`
	UpdatedAt  time.Time `json:"createdAt,omitempty"`
}

type AwsAccountSnapshotParameters struct {
	Autoconnect bool     `url:"autoconnect,omitempty"`
	Exclude     []string `url:"exclude,omitempty,comma"`
	Filter      string   `url:"filter,omitempty"`
	Grid        bool     `url:"grid,omitempty"`
	Height      int      `url:"height,omitempty"`
	Label       bool     `url:"label,omitempty"`
	Landscape   bool     `url:"landscape,omitempty"`
	PaperSize   string   `url:"paperSize,omitempty"`
	Projection  string   `url:"projection,omitempty"`
	Scale       float32  `url:"scale,omitempty"`
	Transparent bool     `url:"transparent,omitempty"`
	Width       int      `url:"width,omitempty"`
}

type AwsAccountSnapshot struct {
	ContentType        string
	Content            *bytes.Buffer
	SnapshotParameters *AwsAccountSnapshotParameters
}

// Convert AwsAccount to a string
func (d AwsAccount) String() string {
	return Stringify(d)
}

type AwsAccountsRoot struct {
	AwsAccounts []AwsAccount `json:"accounts"`
}

type AwsAccountCreateOrUpdateRequest struct {
	Name    string `json:"name"`
	RoleArn string `json:"roleArn"`
}

func (d AwsAccountCreateOrUpdateRequest) String() string {
	return Stringify(d)
}

type AwsAccountSnapshotRequest struct {
	Format             string
	Region             string
	SnapshotParameters *AwsAccountSnapshotParameters
}

func (d AwsAccountSnapshotRequest) String() string {
	return Stringify(d)
}

type AwsAccountIamParameters struct {
	AccountId     string `json:"accountId"`
	ExternalId    string `json:"externalId"`
	AwsConsoleUrl string `json:"awsConsoleUrl"`
}

func (d AwsAccountIamParameters) String() string {
	return Stringify(d)
}

// List all AwsAccounts.
func (s *AwsAccountsServiceOp) List(ctx context.Context) ([]AwsAccount, *Response, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, awsAccountBasePath, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(AwsAccountsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}
	return root.AwsAccounts, resp, err
}

// Get individual AwsAccount.
func (s *AwsAccountsServiceOp) Get(ctx context.Context, awsAccountID string) (*AwsAccount, *Response, error) {
	if awsAccountID == "" {
		return nil, nil, NewArgError("awsAccountID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", awsAccountBasePath, awsAccountID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	awsAccount := new(AwsAccount)
	resp, err := s.client.Do(ctx, req, awsAccount)
	if err != nil {
		return nil, resp, err
	}

	return awsAccount, resp, err
}

// Create AwsAccount
func (s *AwsAccountsServiceOp) Create(ctx context.Context, createRequest *AwsAccountCreateOrUpdateRequest) (*AwsAccount, *Response, error) {
	if createRequest == nil {
		return nil, nil, NewArgError("createRequest", "cannot be nil")
	}

	path := awsAccountBasePath

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	awsAccount := new(AwsAccount)
	resp, err := s.client.Do(ctx, req, awsAccount)
	if err != nil {
		return nil, resp, err
	}

	return awsAccount, resp, err
}

// Update AwsAccount
func (s *AwsAccountsServiceOp) Update(ctx context.Context, awsAccountID string, updateRequest *AwsAccountCreateOrUpdateRequest) (*AwsAccount, *Response, error) {
	if awsAccountID == "" {
		return nil, nil, NewArgError("awsAccountID", "cannot be empty")
	}

	if updateRequest == nil {
		return nil, nil, NewArgError("updateRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%s", awsAccountBasePath, awsAccountID)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, updateRequest)
	if err != nil {
		return nil, nil, err
	}

	awsAccount := new(AwsAccount)
	resp, err := s.client.Do(ctx, req, awsAccount)
	if err != nil {
		return nil, resp, err
	}

	return awsAccount, resp, err
}

// Delete AwsAccount.
func (s *AwsAccountsServiceOp) Delete(ctx context.Context, awsAccountID string) (*Response, error) {
	if awsAccountID == "" {
		return nil, NewArgError("awsAccountID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%s", awsAccountBasePath, awsAccountID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)

	return resp, err
}

// valid PaperSizes: "Letter", "Legal", "Tabloid", "Ledger", "A0", "A1", "A2", "A3", "A4", "A5"
// Format: One of "json", "svg", "png", "pdf", "mxGraph"

// Snapshot AwsAccount.
func (s *AwsAccountsServiceOp) Snapshot(ctx context.Context, awsAccountID string, snapshotRequest *AwsAccountSnapshotRequest) (*AwsAccountSnapshot, *Response, error) {
	if awsAccountID == "" {
		return nil, nil, NewArgError("awsAccountID", "cannot be empty")
	}

	path, err := addOptions(fmt.Sprintf("%s/%s/%s/%s", awsAccountBasePath, awsAccountID, snapshotRequest.Region, snapshotRequest.Format), snapshotRequest.SnapshotParameters)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	awsAccountSnapshot := new(AwsAccountSnapshot)
	resp, err := s.client.Do(ctx, req, awsAccountSnapshot)
	if err != nil {
		return nil, resp, err
	}

	return awsAccountSnapshot, resp, err
}

// Get AwsAccount IAM Parameters.
func (s *AwsAccountsServiceOp) IamParameters(ctx context.Context) (*AwsAccountIamParameters, *Response, error) {
	path := fmt.Sprintf("%s/iamParameters", awsAccountBasePath)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	awsAccountIamParameters := new(AwsAccountIamParameters)
	resp, err := s.client.Do(ctx, req, awsAccountIamParameters)
	if err != nil {
		return nil, resp, err
	}

	return awsAccountIamParameters, resp, err
}
