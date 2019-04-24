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
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

type authRequest struct {
	Token string `json:"token" binding:"required"`
}

var HeaderPattern = `Bearer ([A-Za-z0-9\-\._~\+\/]+=*)$`

func (m *JWTMiddleware) AuthHandler(c *gin.Context) {
	r := &authRequest{}
	err := c.ShouldBindJSON(r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "invalid authentication request",
		})
	}

	auth, err := m.ValidateAuthnRequest(r.Token)
	if err != nil {
		switch err {
		case ErrInvalidToken:
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "invalid JWT",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": fmt.Sprint(err),
			})
		}

		return
	}

	if !auth {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "unauthenticated",
		})
		return
	}

	session, err := m.GetSessionToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": fmt.Sprint(err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": session,
	})
}

func (m *JWTMiddleware) Authenticated(c *gin.Context) {
	headers := c.Request.Header
	t := headers.Get("Authorization")

	cookie, _ := c.Request.Cookie(m.CookieName)

	if t == "" && cookie.String() == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/auth/login")
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
		jwt = cookie.String()
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
}
