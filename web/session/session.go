//-----------------------------------------------------------------------------
// Copyright (c) 2020 Detlef Stern
//
// This file is part of zettelstore.
//
// Zettelstore is free software: you can redistribute it and/or modify it under
// the terms of the GNU Affero General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// Zettelstore is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with Zettelstore. If not, see <http://www.gnu.org/licenses/>.
//-----------------------------------------------------------------------------

// Package session provides utilities for using sessions.
package session

import (
	"context"
	"net/http"
	"strings"
	"time"

	"zettelstore.de/z/auth/token"
	"zettelstore.de/z/config"
	"zettelstore.de/z/domain"
	"zettelstore.de/z/usecase"
)

const sessionName = "zsession"

// SetToken sets the session cookie for later user identification.
func SetToken(w http.ResponseWriter, token []byte, d time.Duration) {
	cookie := http.Cookie{
		Name:     sessionName,
		Value:    string(token),
		Path:     config.URLPrefix(),
		Secure:   config.SecureCookie(),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	if config.PersistentCookie() && d > 0 {
		cookie.Expires = time.Now().Add(d).Add(30 * time.Second).UTC()
	}
	http.SetCookie(w, &cookie)
}

// ClearToken invalidates the session cookie by sending an empty one.
func ClearToken(ctx context.Context, w http.ResponseWriter) context.Context {
	if w != nil {
		SetToken(w, nil, 0)
	}
	return updateContext(ctx, nil, nil)
}

// Handler enriches the request context with optional user information.
type Handler struct {
	next         http.Handler
	getUserByZid usecase.GetUserByZid
}

// NewHandler creates a new handler.
func NewHandler(next http.Handler, getUserByZid usecase.GetUserByZid) *Handler {
	return &Handler{
		next:         next,
		getUserByZid: getUserByZid,
	}
}

type ctxKeyType struct{}

var ctxKey ctxKeyType

// AuthData stores all relevant authentication data for a context.
type AuthData struct {
	User    *domain.Meta
	Token   []byte
	Now     time.Time
	Issued  time.Time
	Expires time.Time
}

// GetAuthData returns the full authentication data from the context.
func GetAuthData(ctx context.Context) *AuthData {
	data, ok := ctx.Value(ctxKey).(*AuthData)
	if ok {
		return data
	}
	return nil

}

// GetUser returns the user meta data from the context, if there is one. Else return nil.
func GetUser(ctx context.Context) *domain.Meta {
	if data := GetAuthData(ctx); data != nil {
		return data.User
	}
	return nil
}

func updateContext(ctx context.Context, user *domain.Meta, data *token.Data) context.Context {
	if data == nil {
		return context.WithValue(ctx, ctxKey, &AuthData{User: user})
	}
	return context.WithValue(ctx, ctxKey, &AuthData{User: user, Token: data.Token, Now: data.Now, Issued: data.Issued, Expires: data.Expires})
}

// ServeHTTP processes one HTTP request.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	k := token.KindJSON
	t := getHeaderToken(r)
	if t == nil {
		k = token.KindHTML
		t = getSessionToken(r)
	}
	if t == nil {
		h.next.ServeHTTP(w, r)
		return
	}
	tokenData, err := token.CheckToken(t, k)
	if err != nil {
		h.next.ServeHTTP(w, r)
		return
	}
	ctx := r.Context()
	user, err := h.getUserByZid.Run(ctx, tokenData.Zid, tokenData.Ident)
	if err != nil {
		h.next.ServeHTTP(w, r)
		return
	}
	h.next.ServeHTTP(w, r.WithContext(updateContext(ctx, user, &tokenData)))
}

func getSessionToken(r *http.Request) []byte {
	cookie, err := r.Cookie(sessionName)
	if err != nil {
		return nil
	}
	return []byte(cookie.Value)
}

func getHeaderToken(r *http.Request) []byte {
	h := r.Header["Authorization"]
	if h == nil {
		return nil
	}

	// “Multiple message-header fields with the same field-name MAY be
	// present in a message if and only if the entire field-value for that
	// header field is defined as a comma-separated list.”
	// — “Hypertext Transfer Protocol” RFC 2616, subsection 4.2
	auth := strings.Join(h, ", ")

	const prefix = "Bearer "
	// RFC 2617, subsection 1.2 defines the scheme token as case-insensitive.
	if len(auth) < len(prefix) || !strings.EqualFold(auth[:len(prefix)], prefix) {
		return nil
	}
	return []byte(auth[len(prefix):])
}
