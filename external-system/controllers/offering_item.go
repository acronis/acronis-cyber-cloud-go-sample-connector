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

// "Package controllers provides actions to access data layer"
package controllers

import (
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
)

// CreateOrUpdateOfferingItem function creates new offering item if it doesn't exist
// else update the existing offering item
func CreateOrUpdateOfferingItem(oi *models.OfferingItem) (bool, error) {
	return oi.CreateOrUpdate()
}

// GetOfferingItems function fetches offeringItems
func GetOfferingItems(offset, limit int) ([]models.OfferingItem, error) {
	return models.GetOfferingItems(offset, limit)
}

// DeleteOfferingItemByTenantIDAndName function deletes offeringItem for a tenantID with the given item name
func DeleteOfferingItemByTenantIDAndName(tenantID, name string) error {
	return models.DeleteOfferingItemByTenantIDAndName(tenantID, name)
}

// GetOfferingItemByTenantIDAndName function returns offering item for a tenantID with the given item name
func GetOfferingItemByTenantIDAndName(tenantID, name string) (models.OfferingItem, error) {
	return models.GetOfferingItemByTenantIDAndName(tenantID, name)
}
