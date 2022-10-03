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

func TestClient_CreateOrUpdateAccessPolicy(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)

	type args struct {
		accessPolicy models.AccessPolicy
	}
	tests := []struct {
		name     string
		args     args
		isCreate bool
		wantErr  bool
	}{
		{
			name: "successful creation of access policy",
			args: args{
				accessPolicy: GetMockAccessPolicy("a1"),
			},
			isCreate: true,
			wantErr:  false,
		},
		{
			name: "successful update of access policy",
			args: args{
				accessPolicy: GetMockAccessPolicy("a1"),
			},
			isCreate: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateOrUpdateAccessPolicy(&tt.args.accessPolicy)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateOrUpdateAccessPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.isCreate {
				t.Errorf("Client.CreateOrUpdateAccessPolicy() = %v, want %v", got, tt.isCreate)
			}
		})
	}

	// test cleanup
	DeleteMockAccessPolicy(client, "a1", t)
}

func TestClient_DeleteAccessPolicy(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)
	CreateMockAccessPolicy(client, "a1", t)
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
				id: "a1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := client.DeleteAccessPolicy(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteAccessPolicy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetAccessPolicies(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)
	for i := 1; i <= handlers.DefaultLimitValue; i++ {
		id := "a" + strconv.Itoa(i)
		CreateMockAccessPolicy(client, id, t)
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
			name: "Number of access policies should be equal to default limit value when there are enough records",
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
			got, err := client.GetAccessPolicies(tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAccessPolicies() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.expectedCount != len(got) {
				t.Errorf("Client.GetAccessPolicies() = %v, want %v", len(got), tt.expectedCount)
			}
		})
	}

	for i := 1; i <= handlers.DefaultLimitValue; i++ {
		id := "a" + strconv.Itoa(i)
		DeleteMockAccessPolicy(client, id, t)
	}
}

func TestClient_GetAccessPolicy(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)
	CreateMockAccessPolicy(client, "a1", t)
	CreateMockAccessPolicy(client, "a2", t)

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
				id: "a1",
			},
			wantID:  "a1",
			wantErr: false,
		},
		{
			name: "Error failed to get item",
			args: args{
				id: "a3",
			},
			wantID:  "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetAccessPolicy(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetAccessPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.ID != tt.wantID {
				t.Errorf("Client.GetAccessPolicy() = %v, want %v", got.ID, tt.wantID)
			}
		})
	}
	DeleteMockAccessPolicy(client, "a1", t)
	DeleteMockAccessPolicy(client, "a2", t)
}

func GetMockAccessPolicy(id string) models.AccessPolicy {
	accessPolicy := models.AccessPolicy{
		ID:          id,
		Version:     0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		TrusteeID:   "trusteeID",
		TrusteeType: "trusteeType",
		IssuerID:    "issuerID",
		TenantID:    "tenantID",
		RoleID:      "roleID",
	}

	return accessPolicy
}

func CreateMockAccessPolicy(client *Client, id string, t *testing.T) {
	accessPolicy := GetMockAccessPolicy(id)
	isCreated := false
	var err error
	isCreated, err = client.CreateOrUpdateAccessPolicy(&accessPolicy)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !isCreated {
		t.Error("Failed to create access policy")
	}
}

func DeleteMockAccessPolicy(client *Client, id string, t *testing.T) {
	if err := client.DeleteAccessPolicy(id); err != nil {
		t.Errorf(err.Error())
	}
}
