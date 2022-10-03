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

// The `handlers` package provides handler functions for routes
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/controllers"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
)

// CreateOrUpdateOfferingItem function to create new offering item
func CreateOrUpdateOfferingItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tenantID := vars["tenant_id"]
	newOI := models.OfferingItem{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&newOI); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	newOI.TenantID = tenantID
	isCreated, err := controllers.CreateOrUpdateOfferingItem(&newOI)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if isCreated {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

// GetOfferingItems function handler for get offering items
func GetOfferingItems(w http.ResponseWriter, r *http.Request) {
	limit, err := getOptionalIntQueryParam(r, "limit")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	if limit == 0 {
		limit = DefaultLimitValue
	}

	offset, err := getOptionalIntQueryParam(r, "offset")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	offeringItems, err := controllers.GetOfferingItems(offset, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, offeringItems)
}

// GetOfferingItem function get offering item
func GetOfferingItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["offering_item_name"]
	tenantID := vars["tenant_id"]
	offeringItem, err := controllers.GetOfferingItemByTenantIDAndName(tenantID, name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, offeringItem)
}

// DeleteOfferingItem handler for delete offering item
func DeleteOfferingItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["offering_item_name"]
	tenantID := vars["tenant_id"]

	err := controllers.DeleteOfferingItemByTenantIDAndName(tenantID, name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
