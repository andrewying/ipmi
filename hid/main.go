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

package hid

import (
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"strings"
)

type Stream struct {
	Device string
}

type StreamMessage struct {
	Key   string
	Ctrl  bool
	Shift bool
	Alt   bool
	Meta  bool
}

func (s *Stream) WebsocketHandler(c *gin.Context) {
	upgrader := websocket.Upgrader{}
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(s.Device)
	if err != nil {
		panic(err)
	}

	defer ws.Close()
	defer file.Close()

	for {
		message := StreamMessage{}
		err := ws.ReadJSON(message)
		if err != nil {
			log.Print(err)
		}

		message.ParseMessage()
		if message.Key != "" {
			bytes := message.GenerateHID()
			bytesEncoded := hex.EncodeToString(bytes[:])
			bytesEncoded = strings.Replace(bytesEncoded, "0x", "\\x", -1)

			command := fmt.Sprintf("printf \"%%b\" '%v' | hid-ops keyboard", bytesEncoded)
			_, err = file.Write([]byte(command))
			if err != nil {
				log.Print(err)
			}
		}
	}
}
