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

func TestUserAndAccessPolicy(t *testing.T) {
	if clientID == "" || clientSecret == "" || dcURL == "" {
		t.Fatal("Please provide clientID, clientSecret and dcURL in tools.go")
	}

	ctx := context.Background()
	httpDefaultClient := &http.Client{Timeout: 120 * time.Second}
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

	// create user on Acronis
	userName := fmt.Sprintf("autotest_%v", time.Now().UnixNano())
	language := "en"
	email := "test_gmail@gmail.com"
	firstName := "test"
	lastName := "auto"

	user, err := client.CreateUser(ctx, &accclient.UserPost{
		TenantID: tenantID,
		Login:    &userName,
		Language: &language,
		Contact: &accclient.Contact{
			Email:     &email,
			Firstname: &firstName,
			Lastname:  &lastName,
		},
	})
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// assign access policy (specified in tools.go)
	accessPolicyList := accclient.AccessPolicyList{[]*accclient.AccessPolicy{
		{
			RoleID:      roleName,
			ID:          "00000000-0000-0000-0000-000000000000",
			IssuerID:    tenantID,
			TenantID:    tenantID,
			TrusteeID:   user.ID,
			TrusteeType: accclient.TrusteeTypeUser,
			Version:     0,
		},
	}}
	updateResp, err := client.UpdateAccessPolicy(ctx, user.ID, &accessPolicyList)
	if err != nil || updateResp.StatusCode != http.StatusOK {
		t.Fatalf("Access policy update failed with error %v and status code ", err)
	}

	// Setup client for external system
	extClient := extclient.NewClient(http.DefaultClient, externalSystemURL)
	externalClient := external.NewExternalSystem(extClient)

	defer func() {
		// Delete user on Acronis
		if err := client.DeleteUser(ctx, user.ID); err != nil {
			t.Fatalf("Failed to delete user %v", err)
		}

		time.Sleep(time.Second * (updateInterval + 5))
		// check that user is deleted on external system
		userFound, err := checkUserExistInExtSystem(externalClient, user.ID)
		if err != nil {
			t.Fatalf("Failed to check user in external system: %v", err)
		}
		if userFound {
			t.Fatal("user is not deleted on external-system")
		}

		// Check that access policies are also deleted.
		err, syncedPoliciesCount := getNumberOfSyncedAccessPolicies(externalClient, &updateResp.Items)
		if err != nil {
			t.Fatal(err)
		}
		if syncedPoliciesCount != 0 {
			t.Fatalf("Expected 0 access policies in external system, got %v", syncedPoliciesCount)
		}
	}()

	// wait for connector to sync data to external-system
	time.Sleep(time.Second * (updateInterval + 5))

	// check that new user ID exists in external-system database
	userFound, err := checkUserExistInExtSystem(externalClient, user.ID)
	if err != nil {
		t.Fatalf("Failed to check user in external system: %v", err)
	}
	if !userFound {
		t.Fatal("user is not synced correctly to external-system")
	}

	// Check that access policies are also synced.
	err, syncedPoliciesCount := getNumberOfSyncedAccessPolicies(externalClient, &updateResp.Items)
	if err != nil {
		t.Fatal(err)
	}
	if syncedPoliciesCount != len(updateResp.Items) {
		t.Fatalf("Expected %v access policies in external system, got %v", len(updateResp.Items), syncedPoliciesCount)
	}
}

func checkUserExistInExtSystem(externalClient core.ExternalSystemClient, userID string) (bool, error) {
	offset := 0
	pageSize := 50
	userFound := false
	for ; ; offset += pageSize {
		ids, err := externalClient.GetActiveUserIDs(offset, pageSize)
		if err != nil {
			return false, err
		}
		for _, id := range ids {
			if id == userID {
				userFound = true
				break
			}
		}
		if len(ids) < pageSize {
			break
		}
	}

	return userFound, nil
}

func getNumberOfSyncedAccessPolicies(externalClient core.ExternalSystemClient, policies *[]*accclient.AccessPolicy) (error, int) {
	offset := 0
	pageSize := 50
	syncedPoliciesCount := 0
	for ; ; offset += pageSize {
		syncedPolicies, err := externalClient.GetActiveAccessPolicyIDs(offset, pageSize)
		if err != nil {
			return fmt.Errorf("Failed to get access policies from external system: %v", err), 0
		}
		// For each synced policy, does the id match with any existing policies?
		for _, item := range syncedPolicies {
			for _, policy := range *policies {
				if policy.ID == item {
					syncedPoliciesCount++
				}
			}
		}
		if len(syncedPolicies) < pageSize {
			break
		}
	}
	return nil, syncedPoliciesCount
}
