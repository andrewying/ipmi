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

package utils

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/websocket"
	"log"
	"os"
	"strings"
)

type HidStream struct {
	Device string
}

type HidStreamMessage struct {
	Key   string
	Ctrl  bool
	Shift bool
	Alt   bool
	Meta  bool
}

func (s *HidStream) WebsocketHandler() context.Handler {
	ws := websocket.New(websocket.Config{})
	ws.OnConnection(s.receiveInput)

	return ws.Handler()
}

func (s *HidStream) receiveInput(c websocket.Connection) {
	file, err := os.Create(s.Device)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	message := HidStreamMessage{}

	c.OnMessage(func(data []byte) {
		err := json.Unmarshal(data, message)
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
	})
}
