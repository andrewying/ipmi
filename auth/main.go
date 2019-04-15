/*
 * Adsisto
 * Copyright (c) 2019 Andrew Ying
 *
 * This program is free software: you can redistribute it and/or modify it under
 * the terms of version 3 of the GNU General Public License as published by the
 * Free Software Foundation.
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
	"errors"
	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type JWTMiddleware struct {
	// Signing algorithm
	SigningAlgorithm string
	// Path to server public key
	PubKeyPath string
	// Path to server private key
	PrivKeyPath string
	// Server public key
	PubKey interface{}
	// Server private key, should be of type *rsa.PrivKey or *ecdsa.PrivKey
	PrivKey interface{}
	// Name of the authorised key interface
	Interface string
	// Interface configuration
	InterfaceConfig map[string]string
	// An AuthorisedKeyInterface instance
	AuthorisedKeys AuthorisedKeysInterface
	CookieName     string
	AuthnTimeout   time.Duration
	SessionTimeout time.Duration
	Leeway         time.Duration
	Unauthorised   func(int, *gin.Context)
}

type AuthorisedKeysInterface interface {
	New(map[string]string)
	Get(...interface{}) (string, error)
}

var (
	// Mapping between name and instances of AuthorisedKeysInterface
	keysInterfaces = map[string]AuthorisedKeysInterface{
		"mysql": &MysqlKeyStore{},
	}

	ErrInvalidAlg         = errors.New("signing algorithm is invalid")
	ErrHMACAlg            = errors.New("HMAC algorithms are not accepted")
	ErrMissingPubKey      = errors.New("public key is required")
	ErrMissingPrivKey     = errors.New("private key is required")
	ErrInvalidExpDuration = errors.New("expiration is longer than the permitted duration")
	ErrInvalidToken       = errors.New("invalid JWT")
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

	if m.AuthorisedKeys == nil {
		m.AuthorisedKeys = keysInterfaces[m.Interface]
		m.AuthorisedKeys.New(m.InterfaceConfig)
	}

	return nil
}

// Validate authentication request for a validly signed JWT
func (m *JWTMiddleware) ValidateAuthnRequest(t string) (bool, error) {
	token, err := jws.ParseJWT([]byte(t))
	if err != nil {
		return false, ErrInvalidToken
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

	key, err := m.AuthorisedKeys.Get(issuer)
	if key == "" {
		return false, nil
	}
	if err != nil {
		log.Print(err)
		return false, err
	}

	err = token.Validate(
		key,
		jws.GetSigningMethod(m.SigningAlgorithm),
		validator,
	)
	if err != nil {
		log.Print(err)
		return false, nil
	}

	return true, nil
}

// Generate session token, in the form of a valid JWT to be signed by the user.
func (m *JWTMiddleware) GetSessionToken() (string, error) {
	now := time.Now()

	claim := jws.Claims{}
	claim.SetIssuedAt(now)
	claim.SetNotBefore(now)
	claim.SetExpiration(now.Add(m.SessionTimeout))

	token := jws.NewJWT(claim, jws.GetSigningMethod(m.SigningAlgorithm))
	bytes, err := token.Serialize(m.PrivKey)
	if err != nil {
		return "", err
	}

	return string(bytes[:]), nil
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
