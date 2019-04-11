/*
 * Copyright (c) Andrew Ying 2019.
 *
 * This file is part of the Intelligent Platform Management Interface (IPMI) software.
 * IPMI is free software. You can use, share, and build it under the terms of the
 * API Copyleft License.
 *
 * As far as the law allows, this software comes as is, without any warranty or
 * condition, and no contributor will be liable to anyone for any damages related
 * to this software or this license, under any kind of legal claim.
 *
 * A copy of the API Copyleft License is available at <LICENSE.md>.
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
	identity string
	token    string
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

	auth, err := m.ValidateAuthnRequest(r.token)
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
		c.Redirect(http.StatusTemporaryRedirect, "auth/login")
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
