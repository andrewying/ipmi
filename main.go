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

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/andrewying/ipmi/auth"
	"github.com/andrewying/ipmi/hid"
	"github.com/gin-gonic/gin"
	"github.com/go-webpack/webpack"
	"github.com/spf13/viper"
	"html/template"
	"time"
)

var (
	assetDirs  = []string{"css", "images", "js"}
	isDev      bool
	appName    string
	domain     string
	cookieName string
	config     *viper.Viper
)

var (
	ErrMissingConfig = errors.New("config file is required")
)

func main() {
	isDev = *flag.Bool("dev", false, "development mode")
	// configPath := flag.String("config", "", "path to config file")

	r := gin.Default()
	if !isDev {
		gin.SetMode(gin.ReleaseMode)
	}

	// loadConfig(*configPath)
	loadAssets(r, isDev)

	r.GET("/", HomeRenderer)
	authRoutes(r)

	s := &hid.Stream{}
	r.GET("api/keystrokes", s.WebsocketHandler)

	r.Run()
}

func loadConfig(path string) {
	viper.SetConfigName("config")

	if path != "" {
		viper.AddConfigPath(path)
	}
	viper.AddConfigPath("/etc/ipmi")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(ErrMissingConfig)
	}

	config = viper.GetViper()
	appName = config.GetString("app.name")
	domain = config.GetString("app.domain")
	cookieName = config.GetString("app.cookie_name")
}

func authRoutes(r *gin.Engine) {
	m := &auth.JWTMiddleware{
		PubKeyPath:       config.GetString("keys.public"),
		PrivKeyPath:      config.GetString("keys.private"),
		SigningAlgorithm: config.GetString("jwt.algorithm"),
		CookieName:       cookieName,
		AuthnTimeout:     time.Minute * time.Duration(config.GetFloat64("jwt.authn_timeout")),
		SessionTimeout:   time.Minute * time.Duration(config.GetFloat64("jwt.session_timeout")),
		Leeway:           time.Second * time.Duration(config.GetFloat64("jwt.leeway")),
	}

	err := m.MiddlewareInit()
	if err != nil {
		panic(err)
	}

	r.GET("auth/login", LoginRenderer)
	r.POST("auth/login", m.AuthHandler)
}

func loadAssets(r *gin.Engine, dev bool) *gin.Engine {
	webpack.FsPath = "./public"
	webpack.WebPath = "/"
	webpack.Plugin = "manifest"
	webpack.Init(dev)

	r.SetFuncMap(template.FuncMap{
		"asset": webpack.AssetHelper,
	})

	for i := 0; i < len(assetDirs); i++ {
		r.Static(
			fmt.Sprintf("/%s", assetDirs[i]),
			fmt.Sprintf("./public/%s", assetDirs[i]),
		)
	}
	r.LoadHTMLGlob("./src/*.tmpl")

	return r
}
