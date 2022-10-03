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
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/controllers"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/models"
)

type WebHandler struct {
	htmlTemplates map[string]*template.Template
	authHandler   *AuthHandler
}

func NewWebHandler(authHandler *AuthHandler) *WebHandler {
	listOfHTMLToParse := []string{
		"index.html",
		"user.html",
	}

	htmlTemplates := make(map[string]*template.Template)

	for _, htmlFile := range listOfHTMLToParse {
		tmpl, err := template.ParseFiles(filepath.Join(authHandler.config.WebUIDirectory, htmlFile))
		if err != nil {
			// Ignore error and continue parsing
			log.Printf("Failed to load html file %s: %v", htmlFile, err)
			continue
		}
		htmlTemplates[htmlFile] = tmpl
	}

	return &WebHandler{
		htmlTemplates: htmlTemplates,
		authHandler:   authHandler,
	}
}

// IndexPage handles GET /ui/index route
func (wh *WebHandler) IndexPage(w http.ResponseWriter, r *http.Request) {
	if tmpl, ok := wh.htmlTemplates["index.html"]; ok {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// UserPage handles GET /ui/user route. Display current user information
func (wh *WebHandler) UserPage(w http.ResponseWriter, r *http.Request) {
	if tmpl, ok := wh.htmlTemplates["user.html"]; ok {
		session, _ := wh.authHandler.sessionStore.Get(r, "auth")

		var (
			user models.User
			err  error
		)

		authenticated, ok := session.Values["authenticated"].(bool)
		if !ok {
			authenticated = false
		}

		if authenticated {
			// Get user
			user, err = controllers.GetUser(session.Values["userid"].(string))
			if err != nil {
				respondWithError(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		data := struct {
			models.User
			Authenticated bool
		}{
			User:          user,
			Authenticated: authenticated,
		}

		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}
