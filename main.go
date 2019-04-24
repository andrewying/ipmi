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
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/adsisto/adsisto/auth"
	"github.com/adsisto/adsisto/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-webpack/webpack"
	"github.com/hashicorp/logutils"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unicode"
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

	file, err := os.OpenFile(
		config.GetString("app.log.file"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666,
	)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	filter := &logutils.LevelFilter{
		Levels: []logutils.LogLevel{"DEBUG", "INFO", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(
			strings.ToUpper(config.GetString("app.log.level")),
		),
		Writer: file,
	}
	log.SetOutput(filter)

	r := gin.Default()
	if !isDev {
		gin.SetMode(gin.ReleaseMode)
	}

	loadAssets(r, isDev)

	r.GET("/", HomeRenderer)
	authRoutes(r)

	s := &utils.HidStream{
		Device: config.GetString("usb.hid_device"),
	}
	r.GET("api/keystrokes", s.WebsocketHandler)
	images := &ImagesUploader{}
	r.POST("api/images", images.UploadHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] Error when starting server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("[INFO] Gracefully shutdown server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERROR] Error during server shutdown: %s\n", err)
	}

	<-ctx.Done()
	log.Println("[INFO] Server exiting")
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
	storeConfig := config.GetStringMapString("keys.store_config")
	parsedStoreConfig := make(map[string]string)

	for index, value := range storeConfig {
		index = strings.TrimSpace(index)
		chars := []rune(index)
		buffer := make([]rune, 0, len(index))

		var prev rune
		for i, curr := range chars {
			if i == 0 {
				buffer = append(buffer, unicode.ToUpper(curr))
			} else if curr != '_' {
				if prev == '_' {
					buffer = append(buffer, unicode.ToUpper(curr))
				} else {
					buffer = append(buffer, unicode.ToLower(curr))
				}
			}
			prev = curr
		}

		parsedStoreConfig[string(buffer)] = value
	}

	m := &auth.JWTMiddleware{
		PubKeyPath:       config.GetString("keys.server.public"),
		PrivKeyPath:      config.GetString("keys.server.private"),
		SigningAlgorithm: config.GetString("jwt.algorithm"),
		Interface:        config.GetString("keys.store"),
		InterfaceConfig:  parsedStoreConfig,
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
