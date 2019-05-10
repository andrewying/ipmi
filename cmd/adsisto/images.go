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

package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/adsisto/adsisto/pkg/response"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

type ImagesUploader struct {
	UploadDir string
}

func (h *ImagesUploader) UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Unable to process file: %s", err),
		})
		return
	}

	log.Printf("[INFO] Processing uploaded file %s\n", header.Filename)

	if filepath.Ext(header.Filename) != "iso" {
		response.JSON(w, http.StatusNotAcceptable, map[string]string{
			"error": "File uploaded must be a .iso file",
		})
		return
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("Unable to process file: %s", err),
		})
		return
	}

	hash := sha256.New()
	_, err = hash.Write(content)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Unable to save file to server",
		})
		return
	}

	err = ioutil.WriteFile(
		fmt.Sprintf("%s/%s.iso", h.UploadDir, hash.Sum(nil)),
		content,
		0644,
	)
	if err != nil {
		log.Printf("[ERROR] Unable to save uploaded image file: %s\n", err)
		response.JSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Unable to save file to server",
		})
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{
		"file": fmt.Sprintf("%s.iso", hash.Sum(nil)),
	})
}
