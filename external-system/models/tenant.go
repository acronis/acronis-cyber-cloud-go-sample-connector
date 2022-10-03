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

// The `models` package provides data and handler objects
package models

import (
	"errors"
	"time"

	"gorm.io/datatypes"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"
)

// Used to hold API response of Tenant data object
type Tenant struct {
	IDNo            string         `json:"id" description:"Unique identifier" gorm:"primaryKey"`
	ParentID        string         `json:"parent_id" description:"ID of parent tenant"`
	VersionNumber   int64          `json:"version" description:"Auto-incremented entity version"`
	CreatedAt       time.Time      `json:"created_at" description:"Creation timestamp of tenant"`
	UpdatedAt       time.Time      `json:"updated_at" description:"Update timestamp of tenant"`
	TenantName      string         `json:"name" description:"Human-readable name that will be displayed to the users"`
	Kind            string         `json:"kind" description:"Kind (type) of tenant"`
	Enabled         bool           `json:"enabled"`
	CustomerType    string         `json:"customer_type,omitempty" description:"enterprise, consumer, small_office"`
	CustomerID      string         `json:"customer_id" description:"ID from external system used for reporting purposes"`
	BrandID         *int64         `json:"brand_id" description:"Brand id for old API"`
	BrandUUID       string         `json:"brand_uuid" description:"Brand ID for new API"`
	InternalFlag    *string        `json:"internal_flag" description:"Internal flag"`
	TenantLanguage  string         `json:"language" description:"Tenant preferred language"`
	OwnerID         string         `json:"owner_id" description:"Identifier of personal tenant owner user."`
	HasChildren     bool           `json:"has_children" description:"Flag to indicate if tenant has children"`
	DefaultIdpID    *string        `json:"default_idp_id" description:"Default identity provider ID"`
	UpdateLock      datatypes.JSON `json:"update_lock" description:"Tenant that owns the lock"`
	AncestralAccess bool           `json:"ancestral_access" description:"Indicates whether tenant's indirect ancestors have access to it"`
	MfaStatus       string         `json:"mfa_status" description:"Multi-factor authentication status"`
	PricingMode     string         `json:"pricing_mode" description:"Mode of tenant's pricing"`
	Contact         datatypes.JSON `json:"contact" description:"Contact details of a tenant"`
	Contacts        datatypes.JSON `json:"contacts" description:"List of contact details of a tenant"`
}

// Retrieve tenant objects from database
func GetTenants(offset, limit int) ([]Tenant, error) {
	tenants := []Tenant{}

	if err := config.DBConn.Limit(limit).Offset(offset).Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

// Get tenant object for given tenantId
func GetTenant(tenantID string) (*Tenant, error) {
	tenant := Tenant{}

	if err := config.DBConn.Where("id_no=?", tenantID).Find(&tenant).Error; err != nil {
		return &tenant, err
	}

	if tenant.IDNo == "" {
		return &tenant, errors.New(ErrItemNotFound)
	}

	return &tenant, nil
}

// Create or Update tenant based on tenant existence
func CreateOrUpdate(tenant *Tenant) (bool, error) {
	isCreated := false
	_, err := GetTenant(tenant.IDNo)
	if err != nil && err.Error() != ErrItemNotFound {
		return isCreated, err
	} else if err != nil && err.Error() == ErrItemNotFound {
		isCreated = true
	}

	return isCreated, config.DBConn.Save(&tenant).Error
}

// Delete tenant object based on tenant Id
func DeleteTenant(id string) error {
	tenant := &Tenant{IDNo: id}
	return config.DBConn.Delete(&tenant).Error
}
