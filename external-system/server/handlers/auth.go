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
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	oidc "github.com/coreos/go-oidc"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/config"
	"github.com/acronis/acronis-cyber-cloud-go-sample-connector/external-system/controllers"
)

// AuthHandler implements handlers related to authorization.
type AuthHandler struct {
	config       *config.Config
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
	sessionStore *sessions.CookieStore
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(appConfig *config.Config) (*AuthHandler, error) {
	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, appConfig.AuthSettings.IDPAddress)
	if err != nil {
		return nil, fmt.Errorf("error creating provider. %w", err)
	}
	oidcConfig := &oidc.Config{
		ClientID: appConfig.AuthSettings.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)

	oauth2Config := &oauth2.Config{
		ClientID:     appConfig.AuthSettings.ClientID,
		ClientSecret: appConfig.AuthSettings.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  appConfig.AuthSettings.RedirectURL,
		Scopes:       []string{oidc.ScopeOpenID},
	}

	return &AuthHandler{
		config:       appConfig,
		oauth2Config: oauth2Config,
		verifier:     verifier,
		sessionStore: sessions.NewCookieStore([]byte(appConfig.AuthSettings.SessionSecret)),
	}, nil
}

// AuthCallback handles GET /auth/acronis/callback route. Handles oauth2 callback and validates the response.
// Redirects to user page after successful validation.
func (a *AuthHandler) AuthCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	session, _ := a.sessionStore.Get(r, "auth")

	if state, ok := session.Values["state"]; ok {
		if r.URL.Query().Get("state") != state {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "state not set", http.StatusBadRequest)
		return
	}

	oauth2Token, err := a.oauth2Config.Exchange(ctx, r.URL.Query().Get("code"))
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
		return
	}

	idToken, err := a.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get user
	sub := idToken.Subject

	log.Printf("Log in user: %v", sub)
	user, err := controllers.GetUser(sub)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Save username in session
	session.Values["authenticated"] = true
	session.Values["username"] = user.Login
	session.Values["userid"] = user.ID
	session.Values["state"] = "" // remove state
	// Save the session before we write to the response/return from the handler.
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/ui/user", http.StatusTemporaryRedirect)
}

// AuthInitiateRedirect handles GET /auth/acronis/redirect route. Redirects the user to authorization endpoint
func (a *AuthHandler) AuthInitiateRedirect(w http.ResponseWriter, r *http.Request) {
	state := generateState()

	// Save the state in user session
	session, _ := a.sessionStore.Get(r, "auth")
	session.Values["state"] = state
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	url := a.oauth2Config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// AuthInitiate handles GET /auth/acronis/initiate route. Verify params and redirects to logon flow
func (a *AuthHandler) AuthInitiate(w http.ResponseWriter, r *http.Request) {
	iss := r.URL.Query().Get("iss")

	// Verify that issuer is valid
	if iss != a.config.AuthSettings.IDPAddress {
		http.Error(w, "invalid issuer", http.StatusBadRequest)
	}

	// Ensure that tenant exist
	tenantID := r.URL.Query().Get("tenant_id")
	if _, err := controllers.GetTenant(tenantID); err != nil {
		http.Error(w, "invalid tenant id", http.StatusBadRequest)
	}

	http.Redirect(w, r, "/auth/acronis/redirect", http.StatusTemporaryRedirect)
}

// Logout handles POST /auth/logout route. Logs the user out and empty user auth session.
// Redirects back to user page.
func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := a.sessionStore.Get(r, "auth")
	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Values["userid"] = ""
	// Save the session before we write to the response/return from the handler.
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/ui/user", http.StatusSeeOther)
}

func generateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
