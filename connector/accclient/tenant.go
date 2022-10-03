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

const customLayout = "2006-01-02T15:04:05"
const nullString = "null"

// Tenant represents the tenant information return by the API
type Tenant struct {
	// Unique identifier
	ID string `json:"id"`

	// Auto-incremented entity version
	Version int64 `json:"version"`

	// RFC3339 Formatted date or custom layout YYYY-MM-DDTHH:MM:SS
	// Creation timestamp
	CreatedAt CustomTime `json:"created_at"`

	// RFC3339 Formatted date or custom layout YYYY-MM-DDTHH:MM:SS
	// Last update timestamp, if tenant has just been created - then is equal to created_at,
	// if tenant has just been deleted - then is equal to deleted_at
	UpdatedAt CustomTime `json:"updated_at"`

	// RFC3339 Formatted date or custom layout YYYY-MM-DDTHH:MM:SS
	// Soft deletion timestamp
	DeletedAt CustomTime `json:"deleted_at"`

	// Human-readable name that will be displayed to the users
	Name string `json:"name"`

	// Business type of the tenant. Current implementations supports following values: enterprise, consumer, small_office
	CustomerType string `json:"customer_type"`

	ParentID string `json:"parent_id"` // ID of parent tenant

	// List of Tenant ID's that has following relationship - next tenant in the list is parent for previous one.
	// First tenant in the list is parent for current tenant.
	Path []string `json:"path"`

	// Kind (type) of the tenant in hierarchy. Current implementations supports following values: root, partner, folder, customer, unit
	Kind string `json:"kind"`

	Contact Contact `json:"contact"`

	// Empty array by default. Will be populated with all referencing contacts if query param 'with_contacts' is provided.
	Contacts []Contact `json:"contacts"`

	// Offering items, empty by default
	OfferingItems []OfferingItem `json:"offering_items"`

	// Flag, indicates whether the tenant is enabled or disabled
	Enabled bool `json:"enabled"`

	// ID from external system; for reporting purposes. This field can have a null value.
	CustomerID *string `json:"customer_id,omitempty"`

	// Brand id for API v1
	BrandID *int64 `json:"brand_id,omitempty"`

	// Brand ID for API v2
	BrandUUID   string  `json:"brand_uuid"`
	InternalTag *string `json:"internal_tag,omitempty"`

	// Tenant`s preferred language
	Language string `json:"language"`

	// Identifier of personal tenant owner user.
	OwnerID *string `json:"owner_id,omitempty"`

	// Flag, indicates if tenant has children
	HasChildren bool `json:"has_children"`

	// Default identity provider ID
	DefaultIDPID *string    `json:"default_idp_id"`
	UpdateLock   UpdateLock `json:"update_lock"`

	// Indicates whether tenant's indirect ancestors have access to it
	AncestralAccess bool `json:"ancestral_access"`

	// Multi-factor authentication status for tenant
	MFAStatus string `json:"mfa_status"`

	// Mode of tenant's pricing
	PricingMode PricingModeType `json:"pricing_mode"`
}

// Contact represents the contact info of the tenant
type Contact struct {
	// Unique identifier
	ID string `json:"id"`
	// RFC3339 Formatted date. Date and time when contact was created
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// RFC3339 Formatted date. Last update timestamp, if contact has just been created - then is equal to created_at
	UpdatedAt        *string       `json:"updated_at,omitempty"`
	Types            []ContactType `json:"types"`
	Email            *string       `json:"email,omitempty"`
	Address1         *string       `json:"address1,omitempty"`
	Address2         *string       `json:"address2,omitempty"`
	Country          *string       `json:"country,omitempty"`
	State            *string       `json:"state,omitempty"`
	City             *string       `json:"city,omitempty"`
	Zipcode          *string       `json:"zipcode,omitempty"`
	Phone            *string       `json:"phone,omitempty"`
	Firstname        *string       `json:"firstname,omitempty"`
	Lastname         *string       `json:"lastname,omitempty"`
	Title            *string       `json:"title,omitempty"`
	Website          *string       `json:"website,omitempty"`
	Industry         *string       `json:"industry,omitempty"`
	OrganizationSize *string       `json:"organization_size,omitempty"`
	EmailConfirmed   *bool         `json:"email_confirmed,omitempty"`
	ExternalID       *string       `json:"external_id,omitempty"`
}

// UpdateLock provides the details of the lock
type UpdateLock struct {
	// If true, updating the tenant via API is not allowed to anyone except for users of the tenant that owns the lock
	Enabled bool `json:"enabled"`

	// ID of tenant that owns the lock (only users of this tenant can update the tenant in question).
	OwnerID *string `json:"owner_id,omitempty"`
}

// ContactType represents the different types of contact
type ContactType string

const (
	ContactTypeLegal     ContactType = "legal"
	ContactTypePrimary   ContactType = "primary"
	ContactTypeBilling   ContactType = "billing"
	ContactTypeTechnical ContactType = "technical"
)

// TenantGetRequest represents the input params for the Get Tenant request
type TenantGetRequest struct {
	// UUIDs is mutually exclusive with 'SubTreeRootID' and 'ParentID'.
	UUIDs []string

	// ParentID is UUID of a tenant which child tenants will be fetched.
	// Is mutually exclusive with 'SubTreeRootID' and 'UUIDs'.
	ParentID string

	// Filter to fetch tenants hierarchy starting from (inclusive) the
	// specified one. Sorting by tenant level is always assumed. Mutually exclusive
	// with 'uuids' and 'parent_id'
	SubTreeRootID string

	// Filter to fetch tenants which were updated later than the specified timestamp
	UpdatedSince *time.Time

	// Limit sets the number of elements in current users page of the response.
	Limit *uint

	// After is a cursor to fetch the next users page. The cursor encodes all the filtering and sorting arguments,
	// thus client does not need to provide all them for the next page, only cursor should be provided.
	After string

	// LevelOfDetail is a predefined level of details for the user object to return.
	LevelOfDetail TenantLOD

	// WithContacts is a flag instructing to display all referencing contacts in response.
	// Is mutually exclusive with 'LevelOfDetail'.
	WithContacts *bool

	// WithOfferingItems can be set to true to embed offering items changes for each tenant
	WithOfferingItems *bool

	// AllowDeleted can be set to true to include tenants and offering items which are disabled
	AllowDeleted bool
}

func (t *TenantGetRequest) getQueryParam() url.Values {
	params := url.Values{}

	if len(t.UUIDs) > 0 {
		params.Set("uuids", strings.Join(t.UUIDs, ","))
	}
	if t.ParentID != "" {
		params.Set("parent_id", t.ParentID)
	}
	if t.SubTreeRootID != "" {
		params.Set("subtree_root_id", t.SubTreeRootID)
	}
	if t.UpdatedSince != nil {
		params.Set("updated_since", t.UpdatedSince.Format(time.RFC3339))
	}
	if t.Limit != nil {
		params.Set("limit", fmt.Sprintf("%d", *t.Limit))
	}
	if t.After != "" {
		params.Set("after", t.After)
	}
	if t.LevelOfDetail != "" {
		params.Set("lod", string(t.LevelOfDetail))
	}
	if t.WithContacts != nil {
		params.Set("with_contacts", fmt.Sprintf("%v", *t.WithContacts))
	}
	if t.WithOfferingItems != nil {
		params.Set("with_offering_items", fmt.Sprintf("%v", *t.WithOfferingItems))
	}

	params.Set("allow_deleted", fmt.Sprintf("%v", t.AllowDeleted))

	return params
}

// TenantLOD represents the available level of details
type TenantLOD string

const (
	TenantLODStamps         TenantLOD = "stamps"
	TenantLODBasic          TenantLOD = "basic"
	TenantLODFullWOContacts TenantLOD = "full_without_contacts"
	TenantLODFull           TenantLOD = "full"
)

// PricingModeTyperepresents the different pricing mode available
type PricingModeType string

const (
	PricingModeTrial      PricingModeType = "trial"
	PricingModeProduction PricingModeType = "production"
	PricingModeSuspended  PricingModeType = "suspended"
)

// CustomTime format that supports both RFC3339 and
// custom layout YYYY-MM-DDTHH:MMSS
type CustomTime struct {
	time.Time
}

// nolint:staticcheck // need to support custom layout parsing here, disable https://staticcheck.io/docs/checks#SA1002
func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == nullString {
		ct.Time = time.Time{}
		return
	}

	ct.Time, err = time.Parse(time.RFC3339, s)
	if err != nil {
		ct.Time, err = time.Parse(customLayout, s)
	}

	return
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	return ct.Time.MarshalJSON()
}

// TenantGetResponse represents the response from the Get Tenant API
type TenantGetResponse struct {
	Response
	Pagination
	Timestamp CustomTime `json:"timestamp"`
	Items     []Tenant   `json:"items"`
}

func parseTenantGetResponse(r *http.Response) (*TenantGetResponse, error) {
	var t TenantGetResponse
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}

	t.StatusCode = r.StatusCode
	t.HTTPHeader = r.Header

	return &t, nil
}

// SingleTenantResponse represents the returned tenant object for GET and POST requests of single tenant
// Currently it's only used for end-to-end test purpose, which only require ID and Version fields.
type SingleTenantResponse struct {
	Response
	ID      string `json:"id"`
	Version int64  `json:"version"`
}

func parseTenantResponse(r *http.Response) (*SingleTenantResponse, error) {
	var t SingleTenantResponse
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return nil, err
	}

	t.StatusCode = r.StatusCode
	t.HTTPHeader = r.Header

	return &t, nil
}

// TenantPostRequest represents the input params for the POST Tenant request
type TenantPostRequest struct {
	// Human-readable name that will be displayed to the users
	Name string `json:"name"`

	// Parent tenant ID in UUID format
	ParentID string `json:"parent_id"`

	// Current implementations supports following values: root, partner, folder, customer, unit
	Kind string `json:"kind"`

	Contact *Contact `json:"contact,omitempty"`

	// Flag, indicates whether the tenant is enabled or disabled
	Enabled *bool `json:"enabled,omitempty"`

	// ID from external system; for reporting purposes
	CustomerID *string `json:"customer_id,omitempty"`

	// Internal tag. This field can have a null value.
	InternalTag *string `json:"internal_tag,omitempty"`

	// Preferred locale, represented with 2 characters, e.g. "en"
	Language *string `json:"language,omitempty"`

	// Default identity provider ID
	DefaultIdpID *string `json:"default_idp_id,omitempty"`

	UpdateLock *UpdateLock `json:"update_lock,omitempty"`

	// Access to newly created tenant from ancestors
	AncestralAccess *bool `json:"ancestral_access,omitempty"`
}

// TenantPutRequest represents the input params for the PUT Tenant request
type TenantPutRequest struct {
	// Human-readable name that will be displayed to the users
	Name *string `json:"name,omitempty"`

	// Current implementations supports following values: enterprise, consumer, small_office
	CustomerType *string `json:"customer_type,omitempty"`

	// Parent tenant ID in UUID format
	ParentID *string `json:"parent_id,omitempty"`

	// Current implementations supports following values: root, partner, folder, customer, unit
	Kind *string `json:"kind,omitempty"`

	Contact *Contact `json:"contact,omitempty"`

	// Flag, indicates whether the tenant is enabled or disabled
	Enabled *bool `json:"enabled,omitempty"`

	// ID from external system; for reporting purposes. This field can have a null value.
	CustomerID *string `json:"customer_id,omitempty"`

	// Tenant`s version
	Version int64 `json:"version"`

	// Deprecated field. Brand cannot be changed using this
	BrandID *int64 `json:"brand_id,omitempty"`

	// Internal tag. This field can have a null value.
	InternalTag *string `json:"internal_tag,omitempty"`

	// Preferred locale, represented with 2 characters, e.g. "en"
	Language *string `json:"language,omitempty"`

	// Default identity provider ID
	DefaultIdpID *string `json:"default_idp_id,omitempty"`

	UpdateLock *UpdateLock `json:"update_lock,omitempty"`

	// Indicates whether tenant's indirect ancestors have access to it
	AncestralAccess *bool `json:"ancestral_access,omitempty"`
}
