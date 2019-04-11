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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-webpack/webpack/reader/manifest"
	"log"
	"net/http"
)

func HomeRenderer(c *gin.Context) {
	if pusher := c.Writer.Pusher(); !isDev && pusher != nil {
		assets, err := manifest.Read("./public")
		if err != nil {
			log.Print("Failed to push assets.")
		}

		pushAssets(pusher, assets)
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"name": "Cam FM",
	})
}

func LoginRenderer(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"name": fmt.Sprintf("Login - %s", " Cam FM"),
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
