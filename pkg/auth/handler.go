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
	"fmt"
	"github.com/SermoDigital/jose/jws"
	"github.com/adsisto/adsisto/pkg/response"
	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"regexp"
)

type authRequest struct {
	Token string `json:"token" validate:"required"`
}

var (
	validate      = validator.New()
	HeaderPattern = `Bearer ([A-Za-z0-9\-\._~\+\/]+=*)$`
)

func (m *JWTMiddleware) AuthHandler(c iris.Context) {
	r := &authRequest{}
	if err := c.ReadJSON(r); err != nil {
		c.StatusCode(http.StatusBadRequest)
		response.JSON(c, iris.Map{
			"code":    http.StatusBadRequest,
			"message": "invalid authentication request",
		})
		return
	}
	if err := validate.Struct(r); err != nil {
		c.StatusCode(http.StatusBadRequest)
		response.JSON(c, iris.Map{
			"code":    http.StatusBadRequest,
			"message": "invalid authentication request",
		})
		return
	}

	if err := c.ReadJSON(r); err != nil {
		c.StatusCode(http.StatusBadRequest)
		response.JSON(c, iris.Map{
			"code":    http.StatusBadRequest,
			"message": "invalid authentication request",
		})
		return
	}

	auth, err := m.ValidateAuthnRequest(r.Token)
	if err != nil {
		switch err {
		case ErrInvalidToken:
			c.StatusCode(http.StatusBadRequest)
			response.JSON(c, iris.Map{
				"code":    http.StatusBadRequest,
				"message": "invalid JWT",
			})
		default:
			c.StatusCode(http.StatusInternalServerError)
			response.JSON(c, iris.Map{
				"code":    http.StatusInternalServerError,
				"message": fmt.Sprint(err),
			})
		}

		return
	}

	if !auth {
		c.StatusCode(http.StatusUnauthorized)
		response.JSON(c, iris.Map{
			"code":    http.StatusUnauthorized,
			"message": "unauthenticated",
		})
		return
	}

	session, err := m.GetSessionToken()
	if err != nil {
		c.StatusCode(http.StatusInternalServerError)
		response.JSON(c, iris.Map{
			"code":    http.StatusInternalServerError,
			"message": fmt.Sprint(err),
		})
		return
	}

	response.JSON(c, iris.Map{
		"token": session,
	})
}

func (m *JWTMiddleware) Authenticated(c iris.Context) {
	t := c.GetHeader("Authorization")

	cookie := c.GetCookie(m.CookieName)

	if t == "" && cookie == "" {
		c.Redirect("/auth/login", http.StatusTemporaryRedirect)
		return
	}

	var jwt string

	if t != "" {
		pattern := regexp.MustCompile(HeaderPattern)
		match := pattern.FindStringSubmatch(t)
		if len(match) == 0 {
			m.Unauthorised(http.StatusBadRequest, c)
			return
		}

		jwt = match[0]
	} else {
		jwt = cookie
	}

	token, err := jws.ParseJWT([]byte(jwt))
	if err != nil {
		m.Unauthorised(http.StatusBadRequest, c)
		return
	}

	claims := token.Claims()
	if claims.Get("iat") == nil || claims.Get("exp") == nil ||
		claims.Get("sub") == nil {
		m.Unauthorised(http.StatusForbidden, c)
		return
	}

	status, err := m.ValidateSessionToken(token)
	if err != nil || !status {
		m.Unauthorised(http.StatusForbidden, c)
		return
	}

	c.Values().SetImmutable("authKey", claims)
	c.Next()
}
