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
	"net/http"
	"time"
)

// AccessPolicyList represents a list of access policy to be assigned to a particular user/trutee
type AccessPolicyList struct {
	Items []*AccessPolicy `json:"items"`
}

// AccessPolicy represents an access policy assigned to a particular user/trutee
type AccessPolicy struct {
	ID          string          `json:"id"`                   // Access policy unique identifier
	Version     int64           `json:"version"`              // Auto-incremented entity version
	CreatedAt   time.Time       `json:"created_at"`           // RFC3339 Formatted date
	UpdatedAt   time.Time       `json:"updated_at"`           // RFC3339 Formatted date
	DeletedAt   *time.Time      `json:"deleted_at,omitempty"` // RFC3339 Formatted date
	TrusteeID   string          `json:"trustee_id"`           // Unique identifier of the Subject for whom access policy is granted.
	TrusteeType TrusteeTypeEnum `json:"trustee_type"`         // Type of the Subject for whom access policy is granted.
	IssuerID    string          `json:"issuer_id"`
	TenantID    string          `json:"tenant_id"`
	RoleID      RoleIDEnum      `json:"role_id"` // Name of user role in current implementation. Will be changed to UUID in next version
	Resource    *Resource       `json:"resource,omitempty"`
}

// TrusteeTypeEnum represents the type of trustee to be assigned an access policy
type TrusteeTypeEnum string

const (
	TrusteeTypeUser      TrusteeTypeEnum = "user"
	TrusteeTypeUserGroup TrusteeTypeEnum = "user_group"
	TrusteeTypeClient    TrusteeTypeEnum = "client"
)

// RoleIDEnum represents the type of role to be assigned to a trustee for a particular access policy
type RoleIDEnum string

const (
	RoleIDRootAdmin       RoleIDEnum = "root_admin"
	RoleIDPartnerAdmin    RoleIDEnum = "partner_admin"
	RoleIDCompanyAdmin    RoleIDEnum = "company_admin"
	RoleIDUnitAdmin       RoleIDEnum = "unit_admin"
	RoleIDReadOnlyAdmin   RoleIDEnum = "readonly_admin"
	RoleIDBackupAdmin     RoleIDEnum = "backup_admin"
	RoleIDBackupUser      RoleIDEnum = "backup_user"
	RoleIDSyncShareAdmin  RoleIDEnum = "sync_share_admin"
	RoleIDSyncShareUser   RoleIDEnum = "sync_share_user"
	RoleIDSyncShareGuest  RoleIDEnum = "sync_share_guest"
	RoleIDMonitoringAdmin RoleIDEnum = "monitoring_admin"
)

// Resource represents a structure for the custom type Resource
type Resource struct {
	ResourceID       string `json:"resource_id"`        // Unique identifier of resource.
	ResourceServerID string `json:"resource_server_id"` // Unique identifier of resource server.
	ScopeType        string `json:"scope_type"`         // Type of scope
}

// UpdateAccessPolicyResponse represents the response from UpdateAccessPolicy request
type UpdateAccessPolicyResponse struct {
	Response
	AccessPolicyList
}

func parseUpdateAccessPolicyResponse(r *http.Response) (*UpdateAccessPolicyResponse, error) {
	var apl UpdateAccessPolicyResponse
	if err := json.NewDecoder(r.Body).Decode(&apl); err != nil {
		return nil, err
	}

	apl.StatusCode = r.StatusCode
	apl.HTTPHeader = r.Header

	return &apl, nil
}
