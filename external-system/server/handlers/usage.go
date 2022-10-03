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

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/controllers"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
)

// FetchUsages handles GET /usages route
func FetchUsages(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))

	if limit == 0 {
		limit = DefaultLimitValue
	}

	usages, err := controllers.GetUsages(offset, limit)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, usages)
}

// CreateOrUpdateUsage handles POST /usages route
func CreateOrUpdateUsage(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var usage models.Usage

	if err := decoder.Decode(&usage); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	isCreated, err := controllers.CreateOrUpdateUsage(&usage)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	if isCreated {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
