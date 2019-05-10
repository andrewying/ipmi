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
	"fmt"
	"github.com/adsisto/adsisto/pkg/response"
	"github.com/go-webpack/webpack/reader/manifest"
	"log"
	"net/http"
)

func HomeRenderer(w http.ResponseWriter, r *http.Request) {
	if pusher, ok := w.(http.Pusher); !*isDev && ok {
		assets, err := manifest.Read("./public")
		if err != nil {
			log.Print("Failed to push assets.")
		}

		pushAssets(pusher, assets)
	}

	response.View(w, templates.Lookup("index.tmpl"), map[string]string{
		"name":       appName,
		"domain":     domain,
		"cookieName": cookieName,
	})
}

func LoginRenderer(w http.ResponseWriter, r *http.Request) {
	response.View(w, templates.Lookup("login.tmpl"), map[string]string{
		"name":       fmt.Sprintf("Login - %s", appName),
		"domain":     domain,
		"cookieName": cookieName,
	})
}

func pushAssets(pusher http.Pusher, assets map[string][]string) {
	for _, files := range assets {
		for _, file := range files {
			if err := pusher.Push(file, nil); err != nil {
				log.Printf("Failed to push: %v", err)
			}
		}
	}
}
