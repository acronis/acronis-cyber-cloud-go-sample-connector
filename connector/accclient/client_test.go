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

package accclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func getTestServer(statusCode int, body interface{}) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respBody, err := json.Marshal(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(statusCode)
		if _, err := w.Write(respBody); err != nil {
			return
		}
	}))

	return srv
}

func getTestEchoServer(statusCode int) *httptest.Server {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(statusCode)
		if _, err := w.Write(body); err != nil {
			return
		}
	}))

	return srv
}

var (
	testTenant1 = Tenant{
		ID:           "dbcacdd4-17f5-4678-8175-f6c35c23fb2d",
		Version:      2,
		Name:         "test_customer_1",
		CustomerType: "consumer",
		ParentID:     "bebc44f6-bb27-4be5-9a80-76998a8ddbc2",
		Kind:         "customer",
		Contact: Contact{
			ID: "d1dd354e-5c46-4ee0-84eb-401c59f45daa",
		},
		Enabled:      true,
		HasChildren:  false,
		DefaultIDPID: nil,
		UpdateLock: UpdateLock{
			Enabled: false,
		},
		AncestralAccess: true,
		MFAStatus:       "enabled",
		PricingMode:     PricingModeTrial,
	}
)

func TestClient_GetTenants(t *testing.T) {
	type args struct {
		ctx    context.Context
		getReq *TenantGetRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *TenantGetResponse
		wantErr bool
	}{
		{
			name: "Successful",
			args: args{
				ctx: context.Background(),
				getReq: &TenantGetRequest{
					SubTreeRootID: "bebc44f6-bb27-4be5-9a80-76998a8ddbc2",
					LevelOfDetail: TenantLODFull,
				},
			},
			want: &TenantGetResponse{
				Response: Response{
					StatusCode: http.StatusOK,
				},
				Items: []Tenant{
					testTenant1,
				},
			},
			wantErr: false,
		},
		{
			name: "Internal Server Error",
			args: args{
				ctx: context.Background(),
				getReq: &TenantGetRequest{
					SubTreeRootID: "bebc44f6-bb27-4be5-9a80-76998a8ddbc2",
					LevelOfDetail: TenantLODFull,
				},
			},
			want: &TenantGetResponse{
				Response: Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: true,
		},
		{
			name: "Unauthorized Error",
			args: args{
				ctx: context.Background(),
				getReq: &TenantGetRequest{
					SubTreeRootID: "bebc44f6-bb27-4be5-9a80-76998a8ddbc2",
					LevelOfDetail: TenantLODFull,
				},
			},
			want: &TenantGetResponse{
				Response: Response{
					StatusCode: http.StatusUnauthorized,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := getTestServer(tt.want.StatusCode, tt.want)
			defer server.Close()

			client := NewClient(server.Client(), server.URL)

			got, err := client.GetTenants(tt.args.ctx, tt.args.getReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetTenants() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				var clientErr Error
				if errors.As(err, &clientErr) {
					if clientErr.Code != fmt.Sprintf("%d", tt.want.StatusCode) {
						t.Errorf("Client.GetTenants() error statuscode = %v, wantErr %d", clientErr.Code, tt.want.StatusCode)
					}
					return
				}
				t.Errorf("Client.GetTenants() not client error: %v", err)
				return
			}

			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("Client.GetTenants() statuscode mismatched = %v, want %v", got.StatusCode, tt.want.StatusCode)
				return
			}

			if !reflect.DeepEqual(got.Items, tt.want.Items) {
				t.Errorf("Client.GetTenants() items mismatched = %v, want %v", got.Items, tt.want.Items)
			}
		})
	}
}

var (
	testOfferingItem1 = OfferingItem{
		ApplicationID:   "6e6d758d-8e74-3ae3-ac84-50eb0dff12eb",
		Name:            "pw_storage",
		TenantID:        "c9c46ef9-1fc1-4002-8862-5db4c08b8b3d",
		UsageName:       "storage",
		Status:          1,
		Locked:          false,
		Type:            "infra",
		InfraID:         "1cd43915-6fa9-4485-9365-fb99a90ce453",
		MeasurementUnit: "bytes",
	}
)

func TestClient_GetOfferingItems(t *testing.T) {
	type args struct {
		ctx    context.Context
		getReq *OfferingItemsGetRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *OfferingItemsGetResponse
		wantErr bool
	}{
		{
			name: "Successful",
			args: args{
				ctx: context.Background(),
				getReq: &OfferingItemsGetRequest{
					SubTreeRootTenantID: "dbcacdd4-17f5-4678-8175-f6c35c23fb2d",
				},
			},
			want: &OfferingItemsGetResponse{
				Response: Response{
					StatusCode: http.StatusOK,
				},
				Items: []OfferingItem{
					testOfferingItem1,
				},
			},
			wantErr: false,
		},
		{
			name: "Internal Server Error",
			args: args{
				ctx: context.Background(),
				getReq: &OfferingItemsGetRequest{
					SubTreeRootTenantID: "dbcacdd4-17f5-4678-8175-f6c35c23fb2d",
				},
			},
			want: &OfferingItemsGetResponse{
				Response: Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: true,
		},
		{
			name: "Unauthorized Error",
			args: args{
				ctx: context.Background(),
				getReq: &OfferingItemsGetRequest{
					SubTreeRootTenantID: "dbcacdd4-17f5-4678-8175-f6c35c23fb2d",
				},
			},
			want: &OfferingItemsGetResponse{
				Response: Response{
					StatusCode: http.StatusUnauthorized,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := getTestServer(tt.want.StatusCode, tt.want)
			defer server.Close()

			client := NewClient(server.Client(), server.URL)

			got, err := client.GetOfferingItems(tt.args.ctx, tt.args.getReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetOfferingItems() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				var clientErr Error
				if errors.As(err, &clientErr) {
					if clientErr.Code != fmt.Sprintf("%d", tt.want.StatusCode) {
						t.Errorf("Client.GetOfferingItems() error statuscode = %v, wantErr %d", clientErr.Code, tt.want.StatusCode)
					}
					return
				}
				t.Errorf("Client.GetOfferingItems() not client error: %v", err)
				return
			}

			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("Client.GetOfferingItems() statuscode mismatched = %v, want %v", got.StatusCode, tt.want.StatusCode)
				return
			}

			if !reflect.DeepEqual(got.Items, tt.want.Items) {
				t.Errorf("Client.GetOfferingItems() items mismatched = %v, want %v", got.Items, tt.want.Items)
			}
		})
	}
}

var (
	testUser1 = User{
		ID:       "6e6d758d-8e74-3ae3-ac84-50eb0dff12eb",
		Version:  2,
		TenantID: "c9c46ef9-1fc1-4002-8862-5db4c08b8b3d",
		Login:    "login",
		Contact: Contact{
			ID: "894fad52-0cfb-11eb-adc1-0242ac120002",
		},
		Activated:     true,
		Enabled:       true,
		Language:      "lang",
		IdpID:         "855d26d0-21e0-4ba4-bb84-b396ce4db9a3",
		ExternalID:    "extid",
		BusinessType:  []BusinessType{},
		Notifications: []UserNotification{},
		MFAStatus:     "disabled",
	}
)

func TestClient_GetUsers(t *testing.T) {
	type args struct {
		ctx    context.Context
		getReq *UserGetRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *UserGetResponse
		wantErr bool
	}{
		{
			name: "Successful",
			args: args{
				ctx: context.Background(),
				getReq: &UserGetRequest{
					SubTreeRootTenantID: "dbcacdd4-17f5-4678-8175-f6c35c23fb2d",
				},
			},
			want: &UserGetResponse{
				Response: Response{
					StatusCode: http.StatusOK,
				},
				Items: []User{
					testUser1,
				},
			},
			wantErr: false,
		},
		{
			name: "Internal Server Error",
			args: args{
				ctx: context.Background(),
				getReq: &UserGetRequest{
					SubTreeRootTenantID: "dbcacdd4-17f5-4678-8175-f6c35c23fb2d",
				},
			},
			want: &UserGetResponse{
				Response: Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: true,
		},
		{
			name: "Unauthorized Error",
			args: args{
				ctx: context.Background(),
				getReq: &UserGetRequest{
					SubTreeRootTenantID: "dbcacdd4-17f5-4678-8175-f6c35c23fb2d",
				},
			},
			want: &UserGetResponse{
				Response: Response{
					StatusCode: http.StatusUnauthorized,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := getTestServer(tt.want.StatusCode, tt.want)
			defer server.Close()

			client := NewClient(server.Client(), server.URL)

			got, err := client.GetUsers(tt.args.ctx, tt.args.getReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				var clientErr Error
				if errors.As(err, &clientErr) {
					if clientErr.Code != fmt.Sprintf("%d", tt.want.StatusCode) {
						t.Errorf("Client.GetUsers() error statuscode = %v, wantErr %d", clientErr.Code, tt.want.StatusCode)
					}
					return
				}
				t.Errorf("Client.GetUsers() not client error: %v", err)
				return
			}

			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("Client.GetUsers() statuscode mismatched = %v, want %v", got.StatusCode, tt.want.StatusCode)
				return
			}

			if !reflect.DeepEqual(got.Items, tt.want.Items) {
				t.Errorf("Client.GetUsers() items mismatched = \ngot  %v, \nwant %v", got.Items, tt.want.Items)
			}
		})
	}
}

var (
	testTenantID1              = "dbcacdd4-17f5-4678-8175-f6c35c23fb2d"
	testOfferingItemStr1       = "test_offering_item"
	testUsageType1             = "count"
	testUsageValue1      int64 = 21

	testUsageReq1 = Usage{
		TenantID:     &testTenantID1,
		OfferingItem: &testOfferingItemStr1,
		UsageType:    &testUsageType1,
		UsageValue:   testUsageValue1,
	}
	testUsageResp1 = UsagesResponse{
		TenantID:     &testTenantID1,
		OfferingItem: &testOfferingItemStr1,
		UsageType:    &testUsageType1,
	}
)

func TestClient_UpdateUsages(t *testing.T) {
	type args struct {
		ctx    context.Context
		getReq *UsagesPutRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *UsagesPutResponse
		wantErr bool
	}{
		{
			name: "Successful",
			args: args{
				ctx: context.Background(),
				getReq: &UsagesPutRequest{
					Items: []Usage{
						testUsageReq1,
					},
				},
			},
			want: &UsagesPutResponse{
				Response: Response{
					StatusCode: http.StatusOK,
				},
				Items: []UsagesResponse{
					testUsageResp1,
				},
			},
			wantErr: false,
		},
		{
			name: "Internal Server Error",
			args: args{
				ctx: context.Background(),
				getReq: &UsagesPutRequest{
					Items: []Usage{
						testUsageReq1,
					},
				},
			},
			want: &UsagesPutResponse{
				Response: Response{
					StatusCode: http.StatusInternalServerError,
				},
			},
			wantErr: true,
		},
		{
			name: "Unauthorized Error",
			args: args{
				ctx: context.Background(),
				getReq: &UsagesPutRequest{
					Items: []Usage{
						testUsageReq1,
					},
				},
			},
			want: &UsagesPutResponse{
				Response: Response{
					StatusCode: http.StatusUnauthorized,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := getTestEchoServer(tt.want.StatusCode)
			defer server.Close()

			client := NewClient(server.Client(), server.URL)

			got, err := client.UpdateUsages(tt.args.ctx, tt.args.getReq)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.UpdateUsages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				var clientErr Error
				if errors.As(err, &clientErr) {
					if clientErr.Code != fmt.Sprintf("%d", tt.want.StatusCode) {
						t.Errorf("Client.UpdateUsages() error statuscode = %v, wantErr %d", clientErr.Code, tt.want.StatusCode)
					}
					return
				}
				t.Errorf("Client.UpdateUsages() not client error: %v", err)
				return
			}

			if got.StatusCode != tt.want.StatusCode {
				t.Errorf("Client.UpdateUsages() statuscode mismatched = %v, want %v", got.StatusCode, tt.want.StatusCode)
				return
			}

			if !reflect.DeepEqual(got.Items, tt.want.Items) {
				t.Errorf("Client.UpdateUsages() items mismatched = \ngot  %v, \nwant %v", got.Items, tt.want.Items)
			}
		})
	}
}
