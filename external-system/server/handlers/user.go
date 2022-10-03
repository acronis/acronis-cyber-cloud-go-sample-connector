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

	"github.com/gorilla/mux"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/controllers"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
)

// FetchUsers handles GET /users route
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.FormValue("offset"))
	limit, _ := strconv.Atoi(r.FormValue("limit"))

	if limit == 0 {
		limit = DefaultLimitValue
	}

	users, err := controllers.GetUsers(offset, limit)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

// FetchUser handles GET /users/id route
func FetchUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user, err := controllers.GetUser(id)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, user)
}

// DeleteUser handles DELETE /users/id route
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := controllers.DeleteUser(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	w.WriteHeader(http.StatusNoContent)
}

// CreateOrUpdateUser handles POST /user route
func CreateOrUpdateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var user models.User

	if err := decoder.Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	isCreated, err := controllers.CreateOrUpdateUser(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	if isCreated {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}
