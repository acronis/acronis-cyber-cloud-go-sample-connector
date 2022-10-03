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
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/storage"
)

var sampleTenant = "t1"

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

func CreateSampleOfferingItem(tenantID, name string) error {
	oi := GetSampleOfferingItem(tenantID, name)
	isCreated := false
	var err error
	if isCreated, err = CreateOrUpdateOfferingItem(&oi); err != nil {
		return err
	}
	if !isCreated {
		return errors.New("Failed to create item")
	}

	return nil
}

func CreateSampleTenant(tenantID string) error {
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
	if isCreated, err = CreateOrUpdateTenant(&tenant); err != nil {
		return err
	}

	if !isCreated {
		return errors.New("Failed to create tenant")
	}

	return nil
}

func DeleteSampleOfferingItem(tenantID, name string) error {
	if err := DeleteOfferingItemByTenantIDAndName(tenantID, name); err != nil {
		return err
	}

	return nil
}

func checkContinue(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestDeleteOfferingItemByTenantIDAndName(t *testing.T) {
	// Setup environment
	checkContinue(CreateSampleOfferingItem(sampleTenant, "i1"), t)
	checkContinue(CreateSampleOfferingItem(sampleTenant, "i2"), t)
	checkContinue(CreateSampleOfferingItem(sampleTenant, "i3"), t)

	type args struct {
		name     string
		tenantID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success case",
			args: args{
				name:     "i1",
				tenantID: sampleTenant,
			},
			wantErr: false,
		},
		{
			name: "failure case wrong name",
			args: args{
				name:     "i5",
				tenantID: sampleTenant,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteOfferingItemByTenantIDAndName(tt.args.tenantID, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DeleteOfferingItemByNameAndTenantID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	// Cleanup
	checkContinue(DeleteSampleOfferingItem(sampleTenant, "i2"), t)
	checkContinue(DeleteSampleOfferingItem(sampleTenant, "i3"), t)
}

func TestGetOfferingItemByNameAndTenantID(t *testing.T) {
	// environment setup
	checkContinue(CreateSampleOfferingItem(sampleTenant, "n1"), t)
	checkContinue(CreateSampleOfferingItem(sampleTenant, "n2"), t)
	checkContinue(CreateSampleOfferingItem(sampleTenant, "n3"), t)

	type args struct {
		name     string
		tenantID string
	}

	tests := []struct {
		name        string
		args        args
		wantName    string
		wanTenantID string
		wantErr     bool
	}{
		{
			name: "Success case",
			args: args{
				name:     "n1",
				tenantID: sampleTenant,
			},
			wantName:    "n1",
			wanTenantID: sampleTenant,
			wantErr:     false,
		},
		{
			name: "Failure case",
			args: args{
				name:     "n5",
				tenantID: sampleTenant,
			},
			wantName:    "",
			wanTenantID: "",
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOfferingItemByTenantIDAndName(tt.args.tenantID, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOfferingItemByTenantIDAndName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantName != got.Name {
				t.Errorf("GetOfferingItemByTenantIDAndName() = %v, want %v", got.Name, tt.wantName)
			}

			if tt.wanTenantID != got.TenantID {
				t.Errorf("GetOfferingItemByTenantIDAndName() = %v, want %v", got.TenantID, tt.wanTenantID)
			}
		})
	}

	// Cleanup
	checkContinue(DeleteSampleOfferingItem(sampleTenant, "n1"), t)
	checkContinue(DeleteSampleOfferingItem(sampleTenant, "n2"), t)
	checkContinue(DeleteSampleOfferingItem(sampleTenant, "n3"), t)
}

func TestCreateOrUpdateOfferingItem(t *testing.T) {
	oi1 := GetSampleOfferingItem(sampleTenant, "oi1")

	type args struct {
		oi *models.OfferingItem
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "success created",
			args: args{
				oi: &oi1,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "success updated",
			args: args{
				oi: &oi1,
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateOrUpdateOfferingItem(tt.args.oi)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOrUpdateOfferingItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateOrUpdateOfferingItem() = %v, want %v", got, tt.want)
			}
		})
	}

	// Cleanup environment
	checkContinue(DeleteSampleOfferingItem(sampleTenant, "oi1"), t)
}

func TestGetOfferingItems(t *testing.T) {
	cleanupOIs()

	// Setup environment
	checkContinue(CreateSampleOfferingItem(sampleTenant, "item1"), t)
	checkContinue(CreateSampleOfferingItem(sampleTenant, "item2"), t)
	checkContinue(CreateSampleOfferingItem(sampleTenant, "item3"), t)

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
			name: "Success case offset 0 limit 1",
			args: args{
				offset: 0,
				limit:  1,
			},
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name: "Success case offset 0 limit 3",
			args: args{
				offset: 0,
				limit:  3,
			},
			expectedCount: 3,
			wantErr:       false,
		},
		{
			name: "Success case offset 1 limit 1",
			args: args{
				offset: 1,
				limit:  1,
			},
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name: "Success case offset 1 limit 4",
			args: args{
				offset: 1,
				limit:  4,
			},
			expectedCount: 2,
			wantErr:       false,
		},
		{
			name: "Success case offset 2 limit 4",
			args: args{
				offset: 2,
				limit:  4,
			},
			expectedCount: 1,
			wantErr:       false,
		},
		{
			name: "Success case offset 3 limit 4",
			args: args{
				offset: 3,
				limit:  4,
			},
			expectedCount: 0,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetOfferingItems(tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOfferingItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.expectedCount != len(got) {
				t.Errorf("GetOfferingItems() = %v, want %v", len(got), tt.expectedCount)
			}
		})
	}

	// Cleanup environment
	checkContinue(DeleteSampleOfferingItem(sampleTenant, "item1"), t)
	checkContinue(DeleteSampleOfferingItem(sampleTenant, "item2"), t)
	checkContinue(DeleteSampleOfferingItem(sampleTenant, "item3"), t)
}

func cleanupOIs() {
	config.DBConn.Exec("DELETE FROM offering_items")
}

func TestMain(m *testing.M) {
	// Tests setup
	dbConfig := config.GetConfig()
	dbConfig.LoadEnvVar()
	storage.SetupDB(dbConfig)

	// Create sample tenant
	if err := CreateSampleTenant(sampleTenant); err != nil {
		fmt.Printf("Failed to create tenant %v\n", err)
		os.Exit(1)
	}

	// Tests run
	exitVal := m.Run()

	// Tests cleanup
	if err := config.DBConn.Migrator().DropTable(&models.OfferingItem{}, &models.Tenant{}, &models.User{}); err != nil {
		fmt.Printf("Cleanup failed %v\n", err)
		os.Exit(1)
	}
	os.Exit(exitVal)
}
