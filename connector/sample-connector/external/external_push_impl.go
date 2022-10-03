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
	"encoding/json"
	"fmt"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/accclient"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/core"
	extclient "github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/client"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
)

// SampleExternalSystem is sample implementation of how we interact with external system
// it implements core.ExternalSystemClient interface.
// This implementation will sync changes from Acronis cloud to a sample external-system server which will
// eventually store the data inside a postgres database.
type SampleExternalSystem struct {
	client *extclient.Client
}

// NewExternalSystem returns SampleExternalSystem as the implementation of core.ExternalSystemClient
func NewExternalSystem(client *extclient.Client) core.ExternalSystemClient {
	return &SampleExternalSystem{client}
}

// CreateOrUpdateTenant handles tenant changes from connector.
// The input parameter is tenant object from Acronis cloud.
// it returns a boolean value indicating whether a new tenant object is created on external-system, and error value if any.
// Failed operation or successful "update" operation on this particular tenant will return false
func (external *SampleExternalSystem) CreateOrUpdateTenant(tenant *accclient.Tenant) (bool, error) {
	contact, contactErr := json.Marshal(tenant.Contact)
	if contactErr != nil {
		return false, fmt.Errorf("failed to process tenant's contact: %v", contactErr)
	}
	contacts, contactsErr := json.Marshal(tenant.Contacts)
	if contactsErr != nil {
		return false, fmt.Errorf("failed to process tenant contacts: %v", contactsErr)
	}

	updateLock, err := json.Marshal(tenant.UpdateLock)
	if err != nil {
		return false, fmt.Errorf("failed to process tenant's UpdateLock: %v", err)
	}

	externalTenant := models.Tenant{
		IDNo:            tenant.ID,
		ParentID:        tenant.ParentID,
		VersionNumber:   tenant.Version,
		CreatedAt:       tenant.CreatedAt.Time,
		UpdatedAt:       tenant.UpdatedAt.Time,
		TenantName:      tenant.Name,
		Kind:            tenant.Kind,
		Enabled:         tenant.Enabled,
		CustomerType:    tenant.CustomerType,
		CustomerID:      getStringSafe(tenant.CustomerID),
		BrandID:         tenant.BrandID,
		InternalFlag:    tenant.InternalTag,
		BrandUUID:       tenant.BrandUUID,
		TenantLanguage:  tenant.Language,
		OwnerID:         getStringSafe(tenant.OwnerID),
		HasChildren:     tenant.HasChildren,
		DefaultIdpID:    tenant.DefaultIDPID,
		UpdateLock:      updateLock,
		AncestralAccess: tenant.AncestralAccess,
		MfaStatus:       tenant.MFAStatus,
		PricingMode:     string(tenant.PricingMode),
		Contact:         contact,
		Contacts:        contacts,
	}

	return external.client.CreateOrUpdateTenant(&externalTenant)
}

// DeleteTenant handles tenant deletion from connector.
// The input parameter accepts Acronis tenantID which should be deleted from external-system.
func (external *SampleExternalSystem) DeleteTenant(tenantID string) error {
	return external.client.DeleteTenant(tenantID)
}

// CreateOrUpdateOfferingItem handles offering item changes from connector
// The input parameter is offering item object from Acronis cloud.
// it returns a boolean value indicating whether a new offering item object is created on external-system, and error value if any.
// Failed operation or successful "update" operation on this particular offering item will return false
func (external *SampleExternalSystem) CreateOrUpdateOfferingItem(item *accclient.OfferingItem) (bool, error) {
	externalItem := models.OfferingItem{
		Quota: models.Quota{
			Value:   item.Quota.Value,
			Overage: item.Quota.Overage,
			Version: item.Quota.Version,
		},
		ApplicationID:   item.ApplicationID,
		Name:            item.Name,
		Edition:         item.Edition,
		UsageName:       item.UsageName,
		TenantID:        item.TenantID,
		UpdatedAt:       item.UpdatedAt.String(),
		Status:          int64(item.Status),
		Locked:          &item.Locked,
		Type:            item.Type,
		InfraID:         &item.InfraID,
		MeasurementUnit: item.MeasurementUnit,
	}

	return external.client.CreateOrUpdateOfferingItem(&externalItem)
}

// DeleteOfferingItem handles offering item deletion from connector
// The input parameter accepts Acronis offering item ID which should be deleted from external-system
func (external *SampleExternalSystem) DeleteOfferingItem(itemID core.OfferingItemID) error {
	return external.client.DeleteOfferingItem(itemID.TenantID, itemID.OfferingItemName)
}

// CreateOrUpdateUser handles user changes from connector
// The input parameter is user object from Acronis cloud.
// it returns a boolean value indicating whether a new user object is created on external-system, and error value if any.
// Failed operation or successful "update" operation on this particular user will return false
func (external *SampleExternalSystem) CreateOrUpdateUser(user *accclient.User) (bool, error) {
	contact, err := json.Marshal(user.Contact)
	if err != nil {
		return false, fmt.Errorf("failed to process user's contact: %v", err)
	}

	externalUser := models.User{
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		ID:        user.ID,
		TenantID:  user.TenantID,
		Login:     user.Login,
		Contact:   string(contact),
		Activated: user.Activated,
		Enabled:   user.Enabled,
	}

	return external.client.CreateOrUpdateUser(&externalUser)
}

// DeleteUser handles user deletion from connector
// The input parameter accepts Acronis user ID which should be deleted from extern-system
func (external *SampleExternalSystem) DeleteUser(userID string) error {
	return external.client.DeleteUser(userID)
}

// CreateOrUpdateAccessPolicy handles access policy changes from connector
// The input parameter is user policy object from Acronis cloud.
// it returns a boolean value indicating whether a new user policy object is created on external-system, and error value if any.
// Failed operation or successful "update" operation on this particular user policy will return false
func (external *SampleExternalSystem) CreateOrUpdateAccessPolicy(accessPolicy *accclient.AccessPolicy) (bool, error) {
	externalAP := models.AccessPolicy{
		ID:          accessPolicy.ID,
		Version:     accessPolicy.Version,
		CreatedAt:   accessPolicy.CreatedAt,
		UpdatedAt:   accessPolicy.UpdatedAt,
		TrusteeID:   accessPolicy.TrusteeID,
		TrusteeType: string(accessPolicy.TrusteeType),
		IssuerID:    accessPolicy.IssuerID,
		TenantID:    accessPolicy.TenantID,
		RoleID:      string(accessPolicy.RoleID),
	}

	return external.client.CreateOrUpdateAccessPolicy(&externalAP)
}

// DeleteAccessPolicy handles access policy deletion from connector
// The input parameter accepts Acronis access policy ID which should be deleted from external-system
func (external *SampleExternalSystem) DeleteAccessPolicy(accessPolicyID string) error {
	return external.client.DeleteAccessPolicy(accessPolicyID)
}

// ===================
// helper functions
// ===================

func getStringSafe(val *string) string {
	if val != nil {
		return *val
	}
	return ""
}
