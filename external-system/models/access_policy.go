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

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"
)

/*
 * Structure for the AccessPolicy type
 */
type AccessPolicy struct {
	ID          string    `json:"id" column:"id"  gorm:"primaryKey"`  // Access policy unique identifier
	Version     int64     `json:"version" column:"version"`           // Auto-incremented entity version
	CreatedAt   time.Time `json:"created_at" column:"created_at"`     // RFC3339 Formatted date
	UpdatedAt   time.Time `json:"updated_at" column:"updated_at"`     // RFC3339 Formatted date
	TrusteeID   string    `json:"trustee_id" column:"trustee_id"`     // Unique identifier of the Subject for whom access policy is granted.
	TrusteeType string    `json:"trustee_type" column:"trustee_type"` // Type of the Subject for whom access policy is granted.
	IssuerID    string    `json:"issuer_id" column:"issuer_id"`       // Unique identifier of access policy Issuer.
	TenantID    string    `json:"tenant_id" column:"tenant_id"`       // Unique identifier of the Subject for whom access policy is granted.
	RoleID      string    `json:"role_id" column:"role_id"`           // Name of user role in current implementation.
}

// Retrieve access policy objects from database
func GetAccessPolicies(offset, limit int) ([]AccessPolicy, error) {
	accessPolicies := []AccessPolicy{}

	if err := config.DBConn.Limit(limit).Offset(offset).Find(&accessPolicies).Error; err != nil {
		return nil, err
	}
	return accessPolicies, nil
}

// Get access policy object for given id
func GetAccessPolicy(id string) (*AccessPolicy, error) {
	accessPolicy := AccessPolicy{}

	if err := config.DBConn.Where("id=?", id).Find(&accessPolicy).Error; err != nil {
		return &accessPolicy, err
	}

	if accessPolicy.ID == "" {
		return &accessPolicy, errors.New(ErrItemNotFound)
	}

	return &accessPolicy, nil
}

// Create or Update access policy based on existence
func CreateOrUpdateAccessPolicy(accessPolicy *AccessPolicy) (bool, error) {
	isCreated := false
	_, err := GetAccessPolicy(accessPolicy.ID)
	if err != nil && err.Error() != ErrItemNotFound {
		return isCreated, err
	} else if err != nil && err.Error() == ErrItemNotFound {
		isCreated = true
	}

	return isCreated, config.DBConn.Save(&accessPolicy).Error
}

// Delete access policy object based on Id
func DeleteAccessPolicy(id string) error {
	accessPolicy := &AccessPolicy{ID: id}
	return config.DBConn.Delete(&accessPolicy).Error
}
