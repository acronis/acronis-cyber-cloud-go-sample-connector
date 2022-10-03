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
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/server"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/storage"
)

var oiTenant = "oi_t1"
var testServerAddr string

func CreateSampleTenant(client *Client, tenantID string) error {
	tenant := models.Tenant{
		IDNo:           tenantID,
		ParentID:       "parentID",
		VersionNumber:  1,
		TenantName:     "name",
		Kind:           "kind",
		Enabled:        true,
		CustomerType:   "type",
		CustomerID:     "customer",
		BrandUUID:      "branchUUID",
		TenantLanguage: "english",
		OwnerID:        "owner",
		HasChildren:    false,
	}

	isCreated := false
	var err error
	if isCreated, err = client.CreateOrUpdateTenant(&tenant); err != nil {
		return err
	}

	if !isCreated {
		return errors.New("Failed to create tenant")
	}

	return nil
}

func GetSampleOfferingItem(tenantID, name string) models.OfferingItem {
	oi := models.OfferingItem{
		ApplicationID:   "applicationID",
		Name:            name,
		Edition:         nil,
		UsageName:       "usagename",
		TenantID:        tenantID,
		UpdatedAt:       "updated",
		Status:          1,
		Locked:          nil,
		Type:            "type",
		InfraID:         nil,
		MeasurementUnit: "measurement",
		Quota: models.Quota{
			Value:   nil,
			Overage: nil,
			Version: 0.0,
		},
	}

	return oi
}

func CreateSampleOfferingItem(client *Client, tenantID, name string) error {
	oi := GetSampleOfferingItem(tenantID, name)
	isCreated := false
	var err error
	isCreated, err = client.CreateOrUpdateOfferingItem(&oi)
	if err != nil {
		return err
	}

	if !isCreated {
		return errors.New("Failed to create item")
	}

	return nil
}

func DeleteSampleOfferingItem(client *Client, tenantID, name string) error {
	if err := client.DeleteOfferingItem(tenantID, name); err != nil {
		return err
	}

	return nil
}

func checkContinue(err error, t *testing.T) {
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestClient_CreateOrUpdateOfferingItem(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)
	checkContinue(CreateSampleTenant(client, oiTenant), t)

	type args struct {
		oi models.OfferingItem
	}
	tests := []struct {
		name     string
		args     args
		isCreate bool
		wantErr  bool
	}{
		{
			name: "successful creation of offering item",
			args: args{
				oi: GetSampleOfferingItem(oiTenant, "o1"),
			},
			isCreate: true,
			wantErr:  false,
		},
		{
			name: "successful update of offering item",
			args: args{
				oi: GetSampleOfferingItem(oiTenant, "o1"),
			},
			isCreate: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.CreateOrUpdateOfferingItem(&tt.args.oi)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.CreateOrUpdateOfferingItem() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.isCreate {
				t.Errorf("Client.CreateOrUpdateOfferingItem() = %v, want %v", got, tt.isCreate)
			}
		})
	}

	// test cleanup
	checkContinue(DeleteSampleOfferingItem(client, oiTenant, "o1"), t)
	DeleteMockTenant(client, oiTenant, t)
}

func TestClient_DeleteOfferingItem(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)
	checkContinue(CreateSampleTenant(client, oiTenant), t)
	checkContinue(CreateSampleOfferingItem(client, oiTenant, "d1"), t)

	type args struct {
		tenantID string
		name     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "successful deletion",
			args: args{
				tenantID: oiTenant,
				name:     "d1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := client.DeleteOfferingItem(tt.args.tenantID, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Client.DeleteOfferingItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	DeleteMockTenant(client, oiTenant, t)
}

func TestClient_GetOfferingItems(t *testing.T) {
	cleanupOIs()
	client := NewClient(http.DefaultClient, testServerAddr)
	checkContinue(CreateSampleTenant(client, oiTenant), t)
	checkContinue(CreateSampleOfferingItem(client, oiTenant, "g1"), t)
	checkContinue(CreateSampleOfferingItem(client, oiTenant, "g2"), t)
	checkContinue(CreateSampleOfferingItem(client, oiTenant, "g3"), t)

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
			name: "success offset 0 limit 1",
			args: args{
				offset: 0,
				limit:  1,
			},
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name: "success offset 0 limit 3",
			args: args{
				offset: 0,
				limit:  3,
			},
			expectedCount: 3,
			wantErr:       false,
		},
		{
			name: "success offset 2 limit 1",
			args: args{
				offset: 2,
				limit:  1,
			},
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name: "success offset 2 limit 3",
			args: args{
				offset: 2,
				limit:  3,
			},
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name: "success offset 3 limit 1",
			args: args{
				offset: 3,
				limit:  1,
			},
			expectedCount: 0,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetOfferingItems(tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetOfferingItems() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.expectedCount != len(got) {
				t.Errorf("Client.GetOfferingItems() = %v, want %v", len(got), tt.expectedCount)
			}
		})
	}

	checkContinue(DeleteSampleOfferingItem(client, oiTenant, "g1"), t)
	checkContinue(DeleteSampleOfferingItem(client, oiTenant, "g2"), t)
	checkContinue(DeleteSampleOfferingItem(client, oiTenant, "g3"), t)
	DeleteMockTenant(client, oiTenant, t)
}

func TestClient_GetOfferingItem(t *testing.T) {
	client := NewClient(http.DefaultClient, testServerAddr)
	checkContinue(CreateSampleTenant(client, oiTenant), t)
	checkContinue(CreateSampleOfferingItem(client, oiTenant, "p1"), t)
	checkContinue(CreateSampleOfferingItem(client, oiTenant, "p2"), t)

	type args struct {
		tenantID string
		name     string
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantErr  bool
	}{
		{
			name: "success get item",
			args: args{
				tenantID: oiTenant,
				name:     "p1",
			},
			wantName: "p1",
			wantErr:  false,
		},
		{
			name: "Error failed to get item",
			args: args{
				tenantID: oiTenant,
				name:     "p3",
			},
			wantName: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.GetOfferingItem(tt.args.tenantID, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetOfferingItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil && got.Name != tt.wantName {
				t.Errorf("Client.GetOfferingItem() = %v, want %v", got.Name, tt.wantName)
			}
		})
	}
	checkContinue(DeleteSampleOfferingItem(client, oiTenant, "p1"), t)
	checkContinue(DeleteSampleOfferingItem(client, oiTenant, "p2"), t)
	DeleteMockTenant(client, oiTenant, t)
}

func cleanupOIs() {
	config.DBConn.Exec("DELETE FROM offering_items")
}

func TestMain(m *testing.M) {
	// Setup the database
	appConfig := config.GetConfig()
	appConfig.LoadEnvVar()
	storage.SetupDB(appConfig)

	exitVal := func() int {
		// Server setup
		app := &server.App{}
		if err := app.InitializeRoutes(appConfig); err != nil {
			fmt.Printf("Error initializing routes: %v\n", err)
			os.Exit(1)
		}

		testServer := httptest.NewServer(app.Router)
		defer testServer.Close()

		testServerAddr = testServer.URL

		// Run the tests
		return m.Run()
	}()

	// Cleanup
	if err := config.DBConn.Migrator().DropTable(&models.OfferingItem{}, &models.Tenant{}, &models.AccessPolicy{}); err != nil {
		fmt.Println("Cleanup failed")
		os.Exit(1)
	}
	os.Exit(exitVal)
}
