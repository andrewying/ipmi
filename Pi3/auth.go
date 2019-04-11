/*
 * Copyright (c) Andrew Ying 2019.
 *
 * This file is part of the Intelligent Platform Management Interface (IPMI) software.
 * IPMI is licensed under the API Copyleft License. A copy of the license is available
 * at LICENSE.md.
 *
 * As far as the law allows, this software comes as is, without any warranty or
 * condition, and no contributor will be liable to anyone for any damages related
 * to this software or this license, under any kind of legal claim.
 */

package main

import (
	"errors"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type JWTMiddleware struct {
	// Signing algorithm
	SigningAlgorithm string
	PubKeyPath       string
	PrivKeyPath      string
	// Server public key
	PubKey interface{}
	// Server private key, should be of type *rsa.PrivKey or *ecdsa.PrivKey
	PrivKey interface{}
	// Map in the format of (identity) => (publicKey) of the list of authorised keys
	AuthorisedKeys map[string]interface{}
	AuthnTimeout   time.Duration
	SessionTimeout time.Duration
	Leeway         time.Duration
	Unauthorized   func(int, *gin.Context)
}

var (
	ErrInvalidAlg         = errors.New("signing algorithm is invalid")
	ErrHMACAlg            = errors.New("HMAC algorithms are not accepted")
	ErrMissingPubKey      = errors.New("public key is required")
	ErrMissingPrivKey     = errors.New("private key is required")
	ErrInvalidExpDuration = errors.New("expiration is longer than the permitted duration")
)

func (m *JWTMiddleware) MiddlewareInit() error {
	switch strings.ToUpper(m.SigningAlgorithm) {
	case "RS256":
	case "RS384":
	case "RS512":
	case "ES256":
	case "ES384":
	case "ES512":
		break
	case "HS256":
	case "HS384":
	case "HS512":
		return ErrHMACAlg
	default:
		return ErrInvalidAlg
	}

	if m.PubKeyPath != "" && m.PubKey == nil {
		keyData, err := ioutil.ReadFile(m.PubKeyPath)
		if err != nil {
			return ErrMissingPubKey
		}

		key, err := m.parsePublicKey(keyData)
		if err != nil {
			return ErrMissingPubKey
		}

		m.PubKey = key
	}

	if m.PubKey == nil {
		return ErrMissingPubKey
	}

	if m.PrivKeyPath != "" && m.PrivKey == nil {
		keyData, err := ioutil.ReadFile(m.PrivKeyPath)
		if err != nil {
			return ErrMissingPrivKey
		}

		key, err := m.parsePrivateKey(keyData)
		if err != nil {
			return ErrMissingPrivKey
		}

		m.PrivKey = key
	}

	if m.PrivKey == nil {
		return ErrMissingPrivKey
	}

	return nil
}

// Validate authentication request for a validly signed JWT
func (m *JWTMiddleware) ValidateAuthnRequest(t string) (bool, error) {
	token, err := jws.ParseJWT([]byte(t))
	if err != nil {
		return false, err
	}

	claims := token.Claims()
	issuer, _ := claims.Issuer()

	validator := jws.NewValidator(
		jws.Claims{},
		m.Leeway,
		m.Leeway,
		func(claims jws.Claims) error {
			exp, _ := claims.Expiration()
			iss, _ := claims.IssuedAt()

			expectedExp := iss.Add(m.AuthnTimeout)
			if expectedExp.Before(exp) {
				return ErrInvalidExpDuration
			}

			return nil
		},
	)

	err = token.Validate(
		m.AuthorisedKeys[issuer],
		jws.GetSigningMethod(m.SigningAlgorithm),
		validator,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Generate session token, in the form of a valid JWT to be signed by the user.
func (m *JWTMiddleware) GetSessionToken() string {
	now := time.Now()

	claim := jws.Claims{}
	claim.SetIssuedAt(now)
	claim.SetNotBefore(now)
	claim.SetExpiration(now.Add(m.SessionTimeout))

	token := jws.NewJWT(claim, jws.GetSigningMethod(m.SigningAlgorithm))
	bytes, err := token.Serialize(m.PrivKey)
	if err != nil {
		panic(err)
	}

	return string(bytes[:])
}

func (m *JWTMiddleware) ValidateSessionToken(t jwt.JWT) (bool, error) {
	validator := jws.NewValidator(
		jws.Claims{},
		m.Leeway,
		m.Leeway,
		func(claims jws.Claims) error {
			return nil
		},
	)

	err := t.Validate(
		m.PubKey,
		jws.GetSigningMethod(m.SigningAlgorithm),
		validator,
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *JWTMiddleware) LoginHandler(c *gin.Context) {
	token, present := c.GetPostForm("token")
	if !present {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "token missing from login request",
		})
		return
	}

	auth, err := m.ValidateAuthnRequest(token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid JWT",
		})
		return
	}

	if !auth {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthenticated",
		})
		return
	}

	session := m.GetSessionToken()
	c.JSON(http.StatusOK, gin.H{
		"token": session,
	})
}

func (m *JWTMiddleware) Authenticated(c *gin.Context) {
	headers := c.Request.Header
	t := headers.Get("Authorization")

	if t == "" {
		c.Redirect(http.StatusTemporaryRedirect, "auth/login")
		return
	}

	token, err := jws.ParseJWT([]byte(t))
	if err != nil {
		m.Unauthorized(http.StatusBadRequest, c)
		return
	}

	claims := token.Claims()
	if claims.Get("iat") == nil || claims.Get("exp") == nil ||
		claims.Get("sub") == nil {
		m.Unauthorized(http.StatusForbidden, c)
		return
	}

	status, err := m.ValidateSessionToken(token)
	if err != nil || !status {
		m.Unauthorized(http.StatusForbidden, c)
		return
	}
}

func (m *JWTMiddleware) AuthenticatedFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.Authenticated(c)
	}
}

func (m *JWTMiddleware) parsePublicKey(k []byte) (interface{}, error) {
	switch strings.ToUpper(m.SigningAlgorithm) {
	case "RS256":
	case "RS384":
	case "RS512":
		return crypto.ParseRSAPublicKeyFromPEM(k)
	case "ES256":
	case "ES384":
	case "ES512":
		return crypto.ParseECPublicKeyFromPEM(k)
	}

	return nil, ErrInvalidAlg
}

func (m *JWTMiddleware) parsePrivateKey(k []byte) (interface{}, error) {
	switch strings.ToUpper(m.SigningAlgorithm) {
	case "RS256":
	case "RS384":
	case "RS512":
		return crypto.ParseRSAPrivateKeyFromPEM(k)
	case "ES256":
	case "ES384":
	case "ES512":
		return crypto.ParseECPrivateKeyFromPEM(k)
	}

	return nil, ErrInvalidAlg
}
