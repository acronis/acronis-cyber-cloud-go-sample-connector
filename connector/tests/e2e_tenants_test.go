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

// +build e2e

package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/accclient"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/core"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/connector/sample-connector/external"
	extclient "github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/client"
)

func TestTenantAndOfferingItems(t *testing.T) {
	if clientID == "" || clientSecret == "" || dcURL == "" {
		t.Fatal("Please provide clientID, clientSecret and dcURL in tools.go")
	}

	if applicationType == "" || roleName == "" || len(offeringItemNames) == 0 {
		t.Fatal("Please provide information about ISV application in tools.go")
	}

	// Initialize the client, which manages the authenticated session
	ctx := context.Background()
	httpDefaultClient := &http.Client{Timeout: 20 * time.Second}
	httpClient := getTestHTTPClient(
		ctx,
		clientID,
		clientSecret,
		dcURL,
		httpDefaultClient,
	)

	client := accclient.NewClient(httpClient, dcURL)

	// Get tenant ID
	tenantID, err := client.GetRegistrationTenantID(ctx, dcURL, clientID)
	if err != nil {
		t.Fatalf("Failed to get tenant ID: %v", err)
	}

	// Get application details from server
	appResp, err := client.GetApplications(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var app *accclient.Application = nil
	for i := range appResp.Items {
		if appResp.Items[i].Type == applicationType {
			app = appResp.Items[i]
			break
		}
	}
	if app == nil {
		t.Fatalf("%v is not available as application: %v", applicationType, err)
	}

	// Initialize request to create a new tenant
	lang := "en"
	tReq := &accclient.TenantPostRequest{
		ParentID: tenantID,
		Kind:     "customer",
		Name:     getTimestampedName("test"),
		Language: &lang,
	}

	tenantObj, err := client.CreateTenant(ctx, tReq)
	if err != nil {
		t.Fatalf("Failed to create tenant: %v", err)
	}
	defer cleanUpTenant(ctx, t, tenantObj.ID, client)

	// Craft PUT request for offering items
	var offPutItems []*accclient.OfferingItemTenantPut
	for i := range offeringItemNames {
		enabled := int8(1)
		item := &accclient.OfferingItemTenantPut{
			ApplicationID: app.ID,
			Name:          offeringItemNames[i],
			Status:        &enabled,
			InfraID:       nil,
			Quota:         nil,
		}
		offPutItems = append(offPutItems, item)
	}
	offeringReq := &accclient.OfferingItemsTenantPutRequest{OfferingItems: offPutItems}
	err = client.UpdateTenantOfferingItems(ctx, tenantObj.ID, offeringReq)
	if err != nil {
		t.Fatalf("Failed to PUT offering items: %v", err)
	}

	// Sleep to ensure tenant creation is synced
	initialBackOff := updateInterval + 5
	time.Sleep(time.Second * time.Duration(initialBackOff))

	// Setup client for external system
	extClient := extclient.NewClient(http.DefaultClient, externalSystemURL)
	externalClient := external.NewExternalSystem(extClient)

	var isTenantInExt bool

	const defaultRetries = 4
	for i := 0; i < defaultRetries; i++ {

		// Check that tenant exists in external-system database
		isTenantInExt, err = externalClient.CheckTenantExist(tenantObj.ID)
		if err != nil {
			t.Fatalf("Failed to query external system for tenant: %v", err)
		}
		if !isTenantInExt {
			// Attempt to retry because connector might be busy syncing
			if i+1 < defaultRetries {
				initialBackOff *= 2
				time.Sleep(time.Second * time.Duration(initialBackOff))

				continue
			}

			t.Fatalf("Tenant %v is not created in external system.", tenantObj.ID)
		}
	}

	// Check that offering item exists in external-system database
	if err = checkSyncedOfferingItems(externalClient, tenantObj.ID, len(offeringItemNames)); err != nil {
		t.Fatal(err)
	}

	// delete tenant
	cleanUpTenant(ctx, t, tenantObj.ID, client)

	time.Sleep(time.Second * (updateInterval + 5))

	// check that tenant is deleted on external-system
	isTenantInExt, err = externalClient.CheckTenantExist(tenantObj.ID)
	if err != nil {
		t.Fatalf("Failed to query external system for tenant: %v", err)
	}
	if isTenantInExt {
		t.Fatalf("Tenant %v is not deleted from external system.", tenantObj.ID)
	}

	// Check that offering item no longer exists in external-system database
	if err = checkSyncedOfferingItems(externalClient, tenantObj.ID, 0); err != nil {
		t.Fatal(err)
	}
}

func cleanUpTenant(ctx context.Context, t *testing.T, tenantID string, client *accclient.Client) {
	tenantObj, err := client.GetTenant(ctx, tenantID)
	if err == nil {
		// Tenant must be disabled before they can be deleted
		tenantActiveStatus := false
		editTenantReq := &accclient.TenantPutRequest{Version: tenantObj.Version, Enabled: &tenantActiveStatus}
		err = client.UpdateTenant(ctx, tenantObj.ID, editTenantReq)
		if err != nil {
			t.Fatalf("Failed to disable tenant: %v", err)
		}
		err = client.DeleteTenant(ctx, tenantObj.ID)
		if err != nil {
			t.Fatalf("Failed to delete tenant: %v", err)
		}
	}
}

func checkSyncedOfferingItems(externalClient core.ExternalSystemClient, newTenantID string, expectedOfferingItemCount int) error {
	offset := 0
	pageSize := 50
	offeringItemCount := 0
	for ; ; offset += pageSize {
		offeringItems, err := externalClient.GetActiveOfferingItemIDs(offset, pageSize)
		if err != nil {
			return fmt.Errorf("Failed to get offering items from external system: %v", err)
		}
		for _, item := range offeringItems {
			if item.TenantID == newTenantID {
				offeringItemCount++
			}
		}
		if len(offeringItems) < pageSize {
			break
		}
	}
	if offeringItemCount != expectedOfferingItemCount {
		return fmt.Errorf("Expected %v offering items in external system, got %v", expectedOfferingItemCount, offeringItemCount)
	}
	return nil
}
