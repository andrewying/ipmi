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
	"crypto/x509"
	"encoding/base64"
	"gopkg.in/go-playground/validator.v9"
	"log"
)

func (m *JWTMiddleware) uniqueIdentityValidator(field validator.FieldLevel) bool {
	value, err := m.AuthorisedKeys.Get(field.Field().String())
	if value != "" {
		return false
	}

	if err != ErrKeyNotFound {
		log.Printf("[ERROR] Validation with authorised key store failed: %s\n", err)
		return false
	}

	return true
}

func (m *JWTMiddleware) existsIdentityValidator(field validator.FieldLevel) bool {
	value, err := m.AuthorisedKeys.Get(field.Field().String())
	if value == "" {
		if err != ErrKeyNotFound {
			log.Printf("[ERROR] Validation with authorised key store failed: %s\n", err)
			return false
		}

		return false
	}

	return true
}

func PublicKeyValidator(field validator.FieldLevel) bool {
	decoded, err := base64.StdEncoding.DecodeString(field.Field().String())
	if err != nil {
		return false
	}

	if _, err := x509.ParsePKCS1PublicKey(decoded); err == nil {
		return true
	}
	if _, err := x509.ParsePKIXPublicKey(decoded); err == nil {
		return true
	}

	return false
}
