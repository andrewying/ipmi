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
	"encoding/json"
	"github.com/adsisto/adsisto/pkg/response"
	"net/http"
)

type newKeyInstance struct {
	Identity    string `json:"identity" validate:"required,email,uniqueIdentity"`
	PublicKey   string `json:"publicKey" validate:"required"`
	AccessLevel string `json:"accessLevel" validate:"omitempty,gte=0"`
}

type existingKeyInstance struct {
	Identity    string `json:"identity" validate:"required,email,existsIdentity"`
	PublicKey   string `json:"publicKey" validate:"required"`
	AccessLevel string `json:"accessLevel" validate:"required,gte=0"`
}

type deleteKeyInstance struct {
	Identity string `json:"identity" validate:"required,email,existsIdentity"`
}

func (m *JWTMiddleware) IndexHandler(w http.ResponseWriter, r *http.Request) {
	keys, err := m.AuthorisedKeys.GetAll()
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "unable to retrieve authorised keys",
		})
		return
	}

	response.JSON(w, http.StatusOK, map[string]interface{}{
		"code": http.StatusOK,
		"keys": keys,
	})
}

func (m *JWTMiddleware) InsertHandler(w http.ResponseWriter, r *http.Request) {
	instance := &newKeyInstance{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(instance); err != nil {
		invalidUserInput(w)
		return
	}

	if err := m.Validator.Struct(instance); err != nil {
		invalidUserInput(w)
		return
	}

	if instance.AccessLevel == "" {
		instance.AccessLevel = "0"
	}

	err := m.AuthorisedKeys.Insert(
		instance.Identity,
		instance.PublicKey,
		instance.AccessLevel,
	)
	if err != nil {
		if err == ErrMethodNotImplemented {
			response.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "method not implemented",
			})
			return
		}

		response.JSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "unable to insert new public key",
		})
		return
	}

	response.JSON(w, http.StatusNoContent, map[string]interface{}{
		"code": http.StatusNoContent,
	})
}

func (m *JWTMiddleware) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	instance := &existingKeyInstance{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(instance); err != nil {
		invalidUserInput(w)
		return
	}

	if err := m.Validator.Struct(instance); err != nil {
		invalidUserInput(w)
		return
	}

	err := m.AuthorisedKeys.Update(
		instance.Identity,
		instance.PublicKey,
		instance.AccessLevel,
	)
	if err != nil {
		if err == ErrMethodNotImplemented {
			response.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "method not implemented",
			})
			return
		}

		response.JSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "unable to update public key",
		})
		return
	}

	response.JSON(w, http.StatusNoContent, map[string]interface{}{
		"code": http.StatusNoContent,
	})
}

func invalidUserInput(w http.ResponseWriter) {
	response.JSON(w, http.StatusBadRequest, map[string]interface{}{
		"code":    http.StatusBadRequest,
		"message": "invalid user inputs",
	})
}

func (m *JWTMiddleware) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	instance := &deleteKeyInstance{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(instance); err != nil {
		invalidUserInput(w)
		return
	}

	if err := m.Validator.Struct(instance); err != nil {
		invalidUserInput(w)
		return
	}

	if err := m.AuthorisedKeys.Delete(instance.Identity); err != nil {
		if err == ErrMethodNotImplemented {
			response.JSON(w, http.StatusBadRequest, map[string]interface{}{
				"code":    http.StatusBadRequest,
				"message": "method not implemented",
			})
			return
		}

		response.JSON(w, http.StatusInternalServerError, map[string]interface{}{
			"code":    http.StatusInternalServerError,
			"message": "unable to delete public key",
		})
		return
	}

	response.JSON(w, http.StatusNoContent, map[string]interface{}{
		"code": http.StatusNoContent,
	})
}
