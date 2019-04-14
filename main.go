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

package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/adsisto/adsisto/auth"
	"github.com/adsisto/adsisto/hid"
	"github.com/gin-gonic/gin"
	"github.com/go-webpack/webpack"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"os"
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
	configPath := flag.String("config", "", "path to config file")
	loadConfig(*configPath)

	file, err := os.OpenFile(config.GetString("app.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	log.SetOutput(file)

	r := gin.Default()
	if !isDev {
		gin.SetMode(gin.ReleaseMode)
	}

	loadAssets(r, isDev)

	r.GET("/", HomeRenderer)
	authRoutes(r)

	s := &hid.Stream{
		Device: config.GetString("usb.hid_device"),
	}
	r.GET("api/keystrokes", s.WebsocketHandler)

	r.Run()
}

func loadConfig(path string) {
	viper.SetConfigName("config")

	if path != "" {
		viper.AddConfigPath(path)
	}
	viper.AddConfigPath("/etc/adsisto")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(ErrMissingConfig)
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
		AuthnTimeout:     time.Minute * time.Duration(config.GetInt("jwt.authn_timeout")),
		SessionTimeout:   time.Minute * time.Duration(config.GetInt("jwt.session_timeout")),
		Leeway:           time.Second * time.Duration(config.GetInt("jwt.leeway")),
	}

	err := m.MiddlewareInit()
	if err != nil {
		log.Panic(err)
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
