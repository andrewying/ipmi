/*
 * Adsisto
 * Copyright (c) 2019 Andrew Ying
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of version 3 of the GNU General Public License as published by the
 * Free Software Foundation. In addition, this program is also subject to certain
 * additional terms available at <SUPPLEMENT.md>.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT ANY
 * WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR
 * A PARTICULAR PURPOSE.  See the GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/SermoDigital/jose/jws"
	"github.com/adsisto/adsisto/pkg/response"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"regexp"
)

type authRequest struct {
	Token string `json:"token" validate:"required"`
}

const (
	claimsKey     int = 0
	HeaderPattern     = `Bearer ([A-Za-z0-9\-\._~\+\/]+=*)$`
)

var (
	validate = validator.New()
)

func (m *JWTMiddleware) AuthHandler(w http.ResponseWriter, r *http.Request) {
	auth := &authRequest{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(auth); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "invalid authentication request",
		})
		return
	}
	if err := validate.Struct(auth); err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"code":    http.StatusBadRequest,
			"message": "invalid authentication request",
		})
		return
	}

	key, err := m.ValidateAuthnRequest(auth.Token)
	if err != nil {
		switch err {
		case ErrInvalidToken:
			response.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "invalid JWT",
			})
		default:
			response.JSON(w, http.StatusInternalServerError, map[string]interface{}{
				"code":    http.StatusInternalServerError,
				"message": fmt.Sprint(err),
			})
		}

		return
	}

	if key == nil {
		response.JSON(w, http.StatusUnauthorized, map[string]interface{}{
			"code":    http.StatusUnauthorized,
			"message": "unauthenticated",
		})
		return
	}

	session, err := m.GetSessionToken(key)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": fmt.Sprint(err),
		})
		return
	}

	cookie := &http.Cookie{
		Name:  m.CookieName,
		Value: session,
	}
	http.SetCookie(w, cookie)

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"code":  http.StatusOK,
		"token": session,
	})
}

func (m *JWTMiddleware) Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := r.Header.Get("Authorization")

		cookie, _ := r.Cookie(m.CookieName)

		if t == "" && cookie.String() == "" {
			http.Redirect(w, r, "/auth/login", http.StatusTemporaryRedirect)
			return
		}

		var jwt string

		if t != "" {
			pattern := regexp.MustCompile(HeaderPattern)
			match := pattern.FindStringSubmatch(t)
			if len(match) == 0 {
				m.Unauthorised(http.StatusBadRequest, w)
				return
			}

			jwt = match[0]
		} else {
			jwt = cookie.String()
		}

		token, err := jws.ParseJWT([]byte(jwt))
		if err != nil {
			m.Unauthorised(http.StatusBadRequest, w)
			return
		}

		claims := token.Claims()
		if claims.Get("iat") == nil || claims.Get("exp") == nil ||
			claims.Get("sub") == nil {
			m.Unauthorised(http.StatusForbidden, w)
			return
		}

		status, err := m.ValidateSessionToken(token)
		if err != nil || !status {
			m.Unauthorised(http.StatusForbidden, w)
			return
		}

		ctx := r.Context()
		r = r.WithContext(context.WithValue(ctx, claimsKey, claims))
		next.ServeHTTP(w, r)
	})
}
