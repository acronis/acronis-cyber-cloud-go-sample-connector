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

package client

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/server/handlers"
)

func TestClient_CreateOrUpdateTenant(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)

	type args struct {
		tenant models.Tenant
	}
	tests := []struct {
		name     string
		args     args
		isCreate bool
		wantErr  bool
	}{
		{
			name: "successful creation of tenant",
			args: args{
				tenant: GetMockTenant("t1"),
			},
			isCreate: true,
			wantErr:  false,
		},
		{
			name: "successful update of tenant",
			args: args{
				tenant: GetMockTenant("t1"),
			},
			isCreate: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateOrUpdateTenant(&tt.args.tenant)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateOrUpdateTenant() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.isCreate {
				t.Errorf("Client.CreateOrUpdateTenant() = %v, want %v", got, tt.isCreate)
			}
		})
	}

	// test cleanup
	DeleteMockTenant(client, "t1", t)
}
func TestClient_DeleteTenant(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)
	CreateMockTenant(client, "t1", t)
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful deletion",
			args: args{
				id: "t1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := client.DeleteTenant(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteTenant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetTenants(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)
	for i := 1; i <= handlers.DefaultLimitValue; i++ {
		tenantID := "t" + strconv.Itoa(i)
		CreateMockTenant(client, tenantID, t)
	}

	type args struct {
		offset int
		limit  int
	}
	tests := []struct {
		name          string
		args          args
		expectedCount int
		wantErr       bool
	}{
		{
			name: "Number of tenants should be equal to default limit value when there are enough records",
			args: args{
				offset: 0,
			},
			expectedCount: handlers.DefaultLimitValue,
			wantErr:       false,
		},
		{
			name: "success offset 0 limit 5",
			args: args{
				offset: 0,
				limit:  5,
			},
			expectedCount: 5,
			wantErr:       false,
		},
		{
			name: "success offset 5 limit 10",
			args: args{
				offset: 5,
				limit:  10,
			},
			expectedCount: 5,
			wantErr:       false,
		},
		{
			name: "success offset 11 limit 5",
			args: args{
				offset: 11,
				limit:  5,
			},
			expectedCount: 0,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetTenants(tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetTenants() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.expectedCount != len(got) {
				t.Errorf("Client.GetTenants() = %v, want %v", len(got), tt.expectedCount)
			}
		})
	}

	for i := 1; i <= handlers.DefaultLimitValue; i++ {
		tenantID := "t" + strconv.Itoa(i)
		DeleteMockTenant(client, tenantID, t)
	}
}

func TestClient_GetTenant(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)
	CreateMockTenant(client, "t1", t)
	CreateMockTenant(client, "t2", t)

	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		wantID  string
		wantErr bool
	}{
		{
			name: "success get item",
			args: args{
				id: "t1",
			},
			wantID:  "t1",
			wantErr: false,
		},
		{
			name: "Error failed to get item",
			args: args{
				id: "t3",
			},
			wantID:  "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetTenant(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetTenant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.IDNo != tt.wantID {
				t.Errorf("Client.GetTenant() = %v, want %v", got.IDNo, tt.wantID)
			}
		})
	}
	DeleteMockTenant(client, "t1", t)
	DeleteMockTenant(client, "t2", t)
}

func GetMockTenant(id string) models.Tenant {
	tenant := models.Tenant{
		IDNo:            id,
		ParentID:        "ParentID",
		VersionNumber:   0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		TenantName:      "TenantName",
		Kind:            "Kind",
		Enabled:         false,
		CustomerType:    "CustomerType",
		CustomerID:      "",
		BrandID:         nil,
		BrandUUID:       "BrandUUID",
		InternalFlag:    nil,
		TenantLanguage:  "TenantLanguage",
		OwnerID:         "",
		HasChildren:     false,
		DefaultIdpID:    nil,
		UpdateLock:      nil,
		AncestralAccess: true,
		MfaStatus:       "MfaStatus",
		PricingMode:     "PricingMode",
		Contact:         nil,
	}

	return tenant
}

func CreateMockTenant(client *Client, id string, t *testing.T) {
	tenant := GetMockTenant(id)
	isCreated := false
	var err error
	isCreated, err = client.CreateOrUpdateTenant(&tenant)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !isCreated {
		t.Error("Failed to create tenant")
	}
}

func DeleteMockTenant(client *Client, id string, t *testing.T) {
	if err := client.DeleteTenant(id); err != nil {
		t.Errorf(err.Error())
	}
}
