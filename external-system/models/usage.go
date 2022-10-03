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
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"
)

// Usage orm model
type Usage struct {
	ID           uint    `gorm:"primaryKey"`
	ResourceID   *string `json:"resourceId,omitempty"`
	UsageType    *string `json:"usageType,omitempty"`
	TenantID     *string `json:"tenantId,omitempty"`
	OfferingItem *string `json:"offeringItem,omitempty"`
	InfraID      *string `json:"infraId,omitempty"`
	UsageValue   int64   `json:"usageValue"`
}

// GetUsages gets list of usages from db
func GetUsages(offset, limit int) ([]Usage, error) {
	usages := make([]Usage, 0, limit)
	if err := config.DBConn.Limit(limit).Offset(offset).Find(&usages).Error; err != nil {
		return nil, err
	}
	return usages, nil
}

// CreateOrUpdateUsage creates usage in db if it doesn't exist, updates if exist
func CreateOrUpdateUsage(usage *Usage) (bool, error) {
	var isCreated bool

	var retrievedUsage Usage
	if err := config.DBConn.Find(&retrievedUsage, "id = ?", usage.ID).Error; err != nil {
		return false, err
	}

	if retrievedUsage.ID != usage.ID {
		isCreated = true
	}

	if err := config.DBConn.Save(&usage).Error; err != nil {
		return false, err
	}

	return isCreated, nil
}
