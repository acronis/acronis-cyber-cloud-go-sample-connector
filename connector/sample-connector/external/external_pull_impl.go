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

package external

import (
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/accclient"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/core"
)

// GetActiveTenantIDs returns tenantIDs from external-system
// This function is required for reconciliation purpose whereby connector needs to compare
// list of tenants between Acronis cloud and external-system and perform synchronizations.
// offset and limit are intended to provide simple mechanism for connector to pull tenants in pages.
// offset indicates starting index of tenants list
// limit indicates how many tenants to be requested starting from "offset"
// The function only returns tenantIDs as a slice of strings.
// when number of tenantIDs returned is less than "limit", it indicates no more pages to be requested.
func (external *SampleExternalSystem) GetActiveTenantIDs(offset, limit int) ([]string, error) {
	extTenants, err := external.client.GetTenants(offset, limit)
	if err != nil {
		return nil, err
	}

	accTenants := make([]string, len(extTenants))
	for i := range extTenants {
		accTenants[i] = extTenants[i].IDNo
	}

	return accTenants, nil
}

// CheckTenantExist checks whether the tenant with the specified tenantID exists in the external system
func (external *SampleExternalSystem) CheckTenantExist(tenantID string) (bool, error) {
	return external.client.CheckTenantExist(tenantID)
}

// GetActiveOfferingItemIDs returns list of offering items from external system
// This function is required for reconciliation purpose whereby connector needs to compare
// list of offering items between Acronis cloud and external-system and perform synchronizations.
// offset and limit are intended to provide simple mechanism for connector to pull offering items in pages.
// offset indicates starting index of offering items list
// limit indicates how many offering items to be requested starting from "offset"
// The function only returns pairs of <offeringItemName, tenantID> as a slice of core.OfferingItemID objects
// when number of objects returned is less than "limit", it indicates no more pages to be requested.
func (external *SampleExternalSystem) GetActiveOfferingItemIDs(offset, limit int) ([]core.OfferingItemID, error) {
	extOIs, err := external.client.GetOfferingItems(offset, limit)
	if err != nil {
		return nil, err
	}

	OIs := make([]core.OfferingItemID, len(extOIs))
	for i := range extOIs {
		if extOIs[i].Status > 0 {
			OIs = append(OIs, core.OfferingItemID{
				OfferingItemName: extOIs[i].Name,
				TenantID:         extOIs[i].TenantID,
			})
		}
	}
	return OIs, nil
}

// GetActiveUserIDs returns users information from external-system
// This function is required for reconciliation purpose whereby connector needs to compare
// list of users between Acronis cloud and external-system and perform synchronizations.
// offset and limit are intended to provide simple mechanism for connector to pull users in pages.
// offset indicates starting index of users list
// limit indicates how many users to be requested starting from "offset"
// The function only returns userIDs as a slice of strings
// when number of user IDs returned is less than "limit", it indicates no more pages to be requested.
func (external *SampleExternalSystem) GetActiveUserIDs(offset, limit int) ([]string, error) {
	extUsers, err := external.client.GetUsers(offset, limit)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, len(extUsers))
	for i := range extUsers {
		userIDs[i] = extUsers[i].ID
	}

	return userIDs, nil
}

// GetActiveAccessPolicyIDs returns access policy IDs from external system which are still active
// This function is required for reconciliation purpose whereby connector needs to compare
// list of access policies between Acronis cloud and external-system and perform synchronizations.
// offset and limit are intended to provide simple mechanism for connector to pull access policies in pages.
// offset indicates starting index of access policies list
// limit indicates how many access policies to be requested starting from "offset"
// The function only returns access policy IDs as a slice of strings
// when number of access policy IDs returned is less than "limit", it indicates no more pages to be requested.
func (external *SampleExternalSystem) GetActiveAccessPolicyIDs(offset, limit int) ([]string, error) {
	extAPs, err := external.client.GetAccessPolicies(offset, limit)
	if err != nil {
		return nil, err
	}

	apIDs := make([]string, len(extAPs))
	for i := range extAPs {
		apIDs[i] = extAPs[i].ID
	}

	return apIDs, nil
}

// GetUsages returns usages from external-system
// These usages will be sync-ed into Acronis cloud, and the returned objects should follow the structure of accclient.Usage
// offset and limit are intended to provide simple mechanism for connector to pull usages in pages.
// offset indicates starting index of usages list
// limit indicates how many usages to be requested starting from "offset"
func (external *SampleExternalSystem) GetUsages(offset, limit int) ([]accclient.Usage, error) {
	extUsage, err := external.client.GetUsages(offset, limit)
	if err != nil {
		return nil, err
	}

	accUsage := make([]accclient.Usage, len(extUsage))
	for i := range extUsage {
		accUsage[i] = accclient.Usage{
			ResourceID:   extUsage[i].ResourceID,
			UsageType:    extUsage[i].UsageType,
			TenantID:     extUsage[i].TenantID,
			OfferingItem: extUsage[i].OfferingItem,
			InfraID:      extUsage[i].InfraID,
			UsageValue:   extUsage[i].UsageValue,
		}
	}

	return accUsage, nil
}
