package cloudcraft

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const userBasePath = "user"

// UsersService is an interface for interfacing with the Users
// endpoints of the Cloudcraft API
// See: https://developers.cloudcraft.co/#398fa0e6-3139-41e6-a5c2-3b9a31e15d6d
type UsersService interface {
	Me(context.Context) (*User, *Response, error)
}

// UsersServiceOp handles communication with the User related methods of the
// Cloudcraft API.
type UsersServiceOp struct {
	client *Client
}

var _ UsersService = &UsersServiceOp{}

// User represents a Cloudcraft User
type User struct {
	ID         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"createdAt,omitempty"`
	CreatorId  string    `json:"CreatorId,omitempty"`
	LastUserId string    `json:"LastUserId,omitempty"`
}

// Convert User to a string
func (d User) String() string {
	return Stringify(d)
}

// Get an individual user. Currently only "me" supported.
func (s *UsersServiceOp) Get(ctx context.Context, userID string) (*User, *Response, error) {
	if userID == "" {
		return nil, nil, NewArgError("userID", "cannot be empty")
	}

	path := fmt.Sprintf("%s/%d", userBasePath, userID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	user := new(User)
	resp, err := s.client.Do(ctx, req, user)
	if err != nil {
		return nil, resp, err
	}

	return user, resp, err
}

func (s *UsersServiceOp) Me(ctx context.Context) (*User, *Response, error) {
	return s.Get(ctx, "me")
}
