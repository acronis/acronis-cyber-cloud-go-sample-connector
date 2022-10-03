// Copyright (c) 2021 Acronis International GmbH
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package accclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type User struct {
	// Unique identifier
	ID string `json:"id"`

	// Auto-incremented entity version
	Version int `json:"version"`

	// ID of tenant this user belongs to
	TenantID string `json:"tenant_id"`

	// User's login
	Login string `json:"login"`

	Contact Contact `json:"contact"`

	// Flag, indicates whether the user has been activated or not
	Activated bool `json:"activated"`

	// Flag, indicates whether the user is enabled or disabled
	Enabled bool `json:"enabled"`

	// Date and time when user was created
	CreatedAt time.Time `json:"created_at"`

	// Last update timestamp, if user has just been created - then is equal to created_at,
	// if user has just been deleted - then is equal to deleted_at
	UpdatedAt time.Time `json:"updated_at"`

	// Soft deletion timestamp
	DeletedAt time.Time `json:"deleted_at"`

	// Preferred locale
	Language string `json:"language"`

	// Identity provider UUID
	IdpID string `json:"idp_id"`

	// User's ID in external identity provider (e.g. SID in AD)
	ExternalID string `json:"external_id"`

	// UUID of user's personal tenant.
	PersonalTenantID *string `json:"personal_tenant_id,omitempty"`

	BusinessType []BusinessType `json:"buyer"`

	Notifications []UserNotification `json:"notifications"`

	// Multi-factor authentication status for user
	MFAStatus string `json:"mfa_status"`

	// Access policies related to the user
	AccessPolicies []AccessPolicy `json:"access_policies"`
}

// UserGetRequest represents the input params for the Get Users request
type UserGetRequest struct {
	// UUIDs is mutually exclusive with 'SubTreeRootTenantID', 'ExternalIDs' and 'TenantID'. Maximum of 100 uuids.
	UUIDs []string

	// ExternalIDs is mutually exclusive with 'SubTreeRootTenantID' and 'UUIDs'. Maximum of 100 ids
	ExternalIDs []string

	// TenantID is a filter to fetch users for the specified tenant. Required when searching
	// by 'external_id'. Is mutually exclusive with 'SubTreeRootTenantID' and 'UUIDs'
	TenantID string

	// SubTreeRootTenantID is a filter to fetch users for tenants hierarchy starting from (inclusive) the specified one.
	// Is mutually exclusive with 'UUIDs' and 'TenantID'
	SubTreeRootTenantID string

	// UpdatedSince is a filter to fetch users which were updated later than the specified timestamp
	UpdatedSince *time.Time

	// Limit sets the number of elements in current users page of the response.
	Limit *uint

	// After is a cursor to fetch the next users page. The cursor encodes all the filtering and sorting arguments,
	// thus client does not need to provide all them for the next page, only cursor should be provided.
	After string

	// LevelOfDetail is a predefined level of details for the user object to return.
	LevelOfDetail UserLOD

	// WithAccessPolicies can be set to true to embed access policies changes for users
	WithAccessPolicies *bool

	// AllowDeleted can be set to true to include users and access policies which are disabled
	AllowDeleted bool
}

func (u *UserGetRequest) getQueryParam() url.Values {
	params := url.Values{}

	if len(u.UUIDs) > 0 {
		params.Set("uuids", strings.Join(u.UUIDs, ","))
	}
	if len(u.ExternalIDs) > 0 {
		params.Set("external_ids", strings.Join(u.ExternalIDs, ","))
	}
	if u.TenantID != "" {
		params.Set("tenant_id", u.TenantID)
	}
	if u.SubTreeRootTenantID != "" {
		params.Set("subtree_root_tenant_id", u.SubTreeRootTenantID)
	}
	if u.UpdatedSince != nil {
		params.Set("updated_since", u.UpdatedSince.Format(time.RFC3339))
	}
	if u.Limit != nil {
		params.Set("limit", fmt.Sprintf("%d", *u.Limit))
	}
	if u.After != "" {
		params.Set("after", u.After)
	}
	if u.LevelOfDetail != "" {
		params.Set("lod", string(u.LevelOfDetail))
	}
	if u.WithAccessPolicies != nil {
		params.Set("with_access_policies", fmt.Sprintf("%v", *u.WithAccessPolicies))
	}

	params.Set("allow_deleted", fmt.Sprintf("%v", u.AllowDeleted))

	return params
}

// UserLOD represents the available level of details
type UserLOD string

const (
	UserLODStamps UserLOD = "stamps"
	UserLODBasic  UserLOD = "basic"
	UserLODFull   UserLOD = "full"
)

// BusinessType represents the type of business
type BusinessType string

const (
	BusinessTypeBuyer BusinessType = "buyer"
)

// UserNotification represents the different user notifications
type UserNotification string

const (
	UserNotificationMaintenance       UserNotification = "maintenance"
	UserNotificationQuota             UserNotification = "quota"
	UserNotificationReports           UserNotification = "reports"
	UserNotificationBackupError       UserNotification = "backup_error"
	UserNotificationBackupWarning     UserNotification = "backup_warning"
	UserNotificationBackupInfo        UserNotification = "backup_info"
	UserNotificationBackupDailyReport UserNotification = "backup_daily_report"
	UserNotificationBackupCritical    UserNotification = "backup_critical"
)

// UserGetResponse represents the response from the Get Users API
type UserGetResponse struct {
	Response
	Pagination
	Timestamp time.Time `json:"timestamp"`
	Items     []User    `json:"items"`
}

func parseUserGetResponse(r *http.Response) (*UserGetResponse, error) {
	var u UserGetResponse
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return nil, err
	}

	u.StatusCode = r.StatusCode
	u.HTTPHeader = r.Header

	return &u, nil
}

// =================================
// Post User
// =================================

// UserPost represents user object that can be created on Acronis cloud
type UserPost struct {
	TenantID      string              `json:"tenant_id"`
	Login         *string             `json:"login,omitempty"`
	ExternalID    *string             `json:"external_id,omitempty"`
	IdpID         *string             `json:"idp_id,omitempty"`
	Contact       *Contact            `json:"contact,omitempty"`
	Enabled       *bool               `json:"enabled,omitempty"`
	Language      *string             `json:"language,omitempty"`
	BusinessTypes *[]BusinessType     `json:"business_types,omitempty"`
	Notifications *[]UserNotification `json:"notifications,omitempty"`
}

// UserPostResponse represents minimal response format returned by Post /users API
// Currently it's only used for autotest purpose
type UserPostResponse struct {
	Response
	ID string `json:"id"`
}

func parseUserPostRespose(r *http.Response) (*UserPostResponse, error) {
	var u UserPostResponse
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		return nil, err
	}

	u.StatusCode = r.StatusCode
	u.HTTPHeader = r.Header

	return &u, nil
}
