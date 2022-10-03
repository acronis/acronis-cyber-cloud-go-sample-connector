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

/* Fetch list of access policies with limit and offset as a query params.
   Accessible for http request routes */
func FetchAccessPolicies(w http.ResponseWriter, r *http.Request) {
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

	accessPolicies, err := controllers.GetAccessPolicies(offset, limit)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, accessPolicies)
}

/* Fetch Access policy based on id */
func FetchAccessPolicy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	accessPolicy, err := controllers.GetAccessPolicy(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, accessPolicy)
}

// Create or update access policy in external system database
func CreateOrUpdateAccessPolicy(w http.ResponseWriter, r *http.Request) {
	accessPolicy := models.AccessPolicy{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&accessPolicy); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	isCreated, createErr := controllers.CreateOrUpdateAccessPolicy(&accessPolicy)
	if createErr != nil {
		respondWithError(w, http.StatusInternalServerError, createErr.Error())
		return
	}

	if isCreated {
		respondWithJSON(w, http.StatusCreated, nil)
	} else {
		respondWithJSON(w, http.StatusNoContent, nil)
	}
}

// Delete access policy in external system database
func DeleteAccessPolicy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	deleteErr := controllers.DeleteAccessPolicy(id)
	if deleteErr != nil {
		respondWithError(w, http.StatusInternalServerError, deleteErr.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
