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

package models

import (
	"errors"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"
)

var (
	// ErrItemNotFound returned when item is not found
	ErrItemNotFound = "ErrorItemNotFound"
)

// OfferingItem struct
type OfferingItem struct {
	Quota
	ID              uint64  `json:"id" column:"id"  gorm:"primaryKey"`
	ApplicationID   string  `json:"application_id" column:"application_id"`
	Name            string  `json:"name" column:"name" gorm:"primaryKey"`
	Edition         *string `json:"edition,omitempty" column:"edition"`
	UsageName       string  `json:"usage_name" column:"usage_name"`
	TenantID        string  `json:"tenant_id" column:"tenant_id" gorm:"primaryKey"`
	UpdatedAt       string  `json:"updated_at" column:"updated_at"`
	Status          int64   `json:"status" column:"status"`
	Locked          *bool   `json:"locked,omitempty" column:"locked"`
	Type            string  `json:"type" column:"type"`
	InfraID         *string `json:"infra_id,omitempty" column:"infra_id"`
	MeasurementUnit string  `json:"measurement_unit" column:"measurement_unit"`
}

// Quota struct
type Quota struct {
	Value   *float64 `json:"value,omitempty" column:"value"`
	Overage *float64 `json:"overage,omitempty" column:"overage"`
	Version float64  `json:"version" column:"version"`
}

// Create function creates new offeringItem
func (oi *OfferingItem) Create() error {
	tx := config.DBConn.Create(&oi)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected != 1 {
		return errors.New(ErrItemNotFound)
	}

	return nil
}

// Update function updates offeringItem
func (oi *OfferingItem) Update() error {
	tx := config.DBConn.Model(&OfferingItem{}).Where("name = ? AND tenant_id = ?", oi.Name, oi.TenantID).Updates(&oi)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected != 1 {
		return errors.New("failed to update offering Item")
	}

	return nil
}

// CreateOrUpdate function creates offeringitem if it doesn't exist
// Updates if it exists
// Name and tenantID fields are used to identify if the offering item exists
func (oi *OfferingItem) CreateOrUpdate() (bool, error) {
	isCreated := false
	_, err := GetOfferingItemByTenantIDAndName(oi.TenantID, oi.Name)
	if err != nil && err.Error() != ErrItemNotFound {
		return isCreated, err
	}

	// Create new item if not found
	if err != nil && err.Error() == ErrItemNotFound {
		// Create new item
		if err := oi.Create(); err != nil {
			return isCreated, err
		}

		isCreated = true
		return isCreated, nil
	}

	if err := oi.Update(); err != nil {
		return isCreated, err
	}

	return isCreated, nil
}

// DeleteOfferingItemByTenantIDAndName function deletes an offering item for a tenant with the given item name
func DeleteOfferingItemByTenantIDAndName(tenantID, name string) error {
	var oi OfferingItem
	tx := config.DBConn.Where("tenant_id = ? AND name = ?", tenantID, name).Delete(oi)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected != 1 {
		return errors.New("failed to delete offering item")
	}

	return nil
}

// GetOfferingItems function fetches the offering items with limit
func GetOfferingItems(offset, limit int) ([]OfferingItem, error) {
	var OfferingItems []OfferingItem
	if err := config.DBConn.Offset(offset).Limit(limit).Find(&OfferingItems).Error; err != nil {
		return nil, err
	}

	return OfferingItems, nil
}

// GetOfferingItemByTenantIDAndName function gets offering item by name and its tenantID
func GetOfferingItemByTenantIDAndName(tenantID, name string) (OfferingItem, error) {
	var oi OfferingItem
	if err := config.DBConn.Where("tenant_id = ? AND name = ?", tenantID, name).Find(&oi).Error; err != nil {
		return OfferingItem{}, err
	}

	if oi.Name == "" {
		return OfferingItem{}, errors.New(ErrItemNotFound)
	}

	return oi, nil
}
