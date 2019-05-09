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
	"github.com/adsisto/adsisto/pkg/response"
	"github.com/kataras/iris"
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

func (m *JWTMiddleware) IndexHandler(c iris.Context) {
	keys, err := m.AuthorisedKeys.GetAll()
	if err != nil {
		c.StatusCode(http.StatusInternalServerError)
		response.JSON(c, iris.Map{
			"code":    http.StatusInternalServerError,
			"message": "unable to retrieve authorised keys",
		})
		return
	}

	response.JSON(c, iris.Map{
		"code": http.StatusOK,
		"keys": keys,
	})
}

func (m *JWTMiddleware) InsertHandler(c iris.Context) {
	instance := &newKeyInstance{}
	err := c.ReadForm(instance)
	if err != nil {
		invalidUserInput(c)
		return
	}

	if err := m.Validator.Struct(instance); err != nil {
		invalidUserInput(c)
		return
	}

	if instance.AccessLevel == "" {
		instance.AccessLevel = "0"
	}

	err = m.AuthorisedKeys.Insert(
		instance.Identity,
		instance.PublicKey,
		instance.AccessLevel,
	)
	if err != nil {
		if err == ErrMethodNotImplemented {
			c.StatusCode(http.StatusBadRequest)
			response.JSON(c, iris.Map{
				"code":    http.StatusBadRequest,
				"message": "method not implemented",
			})
			return
		}

		c.StatusCode(http.StatusInternalServerError)
		response.JSON(c, iris.Map{
			"code":    http.StatusInternalServerError,
			"message": "unable to insert new public key",
		})
		return
	}

	response.JSON(c, iris.Map{
		"code": http.StatusNoContent,
	})
}

func (m *JWTMiddleware) UpdateHandler(c iris.Context) {
	instance := &existingKeyInstance{}
	err := c.ReadForm(instance)
	if err != nil {
		invalidUserInput(c)
		return
	}

	if err := m.Validator.Struct(instance); err != nil {
		invalidUserInput(c)
		return
	}

	err = m.AuthorisedKeys.Update(
		instance.Identity,
		instance.PublicKey,
		instance.AccessLevel,
	)
	if err != nil {
		if err == ErrMethodNotImplemented {
			c.StatusCode(http.StatusBadRequest)
			response.JSON(c, iris.Map{
				"code":    http.StatusBadRequest,
				"message": "method not implemented",
			})
			return
		}

		c.StatusCode(http.StatusInternalServerError)
		response.JSON(c, iris.Map{
			"code":    http.StatusInternalServerError,
			"message": "unable to update new public key",
		})
		return
	}

	response.JSON(c, iris.Map{
		"code": http.StatusNoContent,
	})
}

func invalidUserInput(c iris.Context) {
	c.StatusCode(http.StatusBadRequest)
	response.JSON(c, iris.Map{
		"code":    http.StatusBadRequest,
		"message": "invalid user inputs",
	})
}

func (m *JWTMiddleware) DeleteHandler(c iris.Context) {
	instance := &deleteKeyInstance{}
	err := c.ReadForm(instance)
	if err != nil {
		invalidUserInput(c)
		return
	}

	if err := m.Validator.Struct(instance); err != nil {
		invalidUserInput(c)
		return
	}

	if err := m.AuthorisedKeys.Delete(instance.Identity); err != nil {
		if err == ErrMethodNotImplemented {
			c.StatusCode(http.StatusBadRequest)
			response.JSON(c, iris.Map{
				"code":    http.StatusBadRequest,
				"message": "method not implemented",
			})
			return
		}

		c.StatusCode(http.StatusInternalServerError)
		response.JSON(c, iris.Map{
			"code":    http.StatusInternalServerError,
			"message": "unable to update new public key",
		})
		return
	}

	response.JSON(c, iris.Map{
		"code": http.StatusNoContent,
	})
}
