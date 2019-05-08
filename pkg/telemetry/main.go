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

package telemetry

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"os"
)

var (
	ClientInstance    = &Client{}
	ErrUuidGeneration = errors.New("unable to generate new UUID")
)

type Client struct {
	Uuid    string
	OptedIn bool
}

// New starts a telemetry client instance.
func New(dataDir string) {
	file, err := os.OpenFile(
		dataDir+"telemetry.json",
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666,
	)
	if err != nil {
		log.Println("[ERROR] Unable to open telemetry preferences file")
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Printf("[ERROR] Unable to process telemetry preferences file: %s\n", err)
		return
	}

	if info.Size() == 0 {
		if err = GenerateUuid(); err != nil {
			log.Printf("[ERROR] %s\n", err)
			return
		}

		ClientInstance.OptedIn = false
		log.Printf("[INFO] Generated new telemetry UUID %s\n", ClientInstance.Uuid)

		encoded, err := json.Marshal(ClientInstance)
		if err != nil {
			log.Printf("[ERROR] Unable to generate new telemetry preferences file: %s\n", err)
			return
		}

		if _, err = file.Write(encoded); err != nil {
			log.Printf("[ERROR] %s\n", err)
			return
		}
	} else {
		buf := make([]byte, info.Size())
		if _, err = file.Read(buf); err != nil {
			log.Printf("[ERROR] %s\n", err)
			return
		}

		if err = json.Unmarshal(buf, ClientInstance); err != nil {
			log.Printf("[ERROR] Unable to parse telemetry preferences file: %s\n", err)
			ClientInstance = &Client{}
			return
		}

		if ClientInstance.Uuid == "" {
			if err = GenerateUuid(); err != nil {
				log.Printf("[ERROR] %s\n", err)
				return
			}
		}
	}

	if err = SetupApiClient(); err != nil {
		log.Printf("[ERROR] Unable to set up telemetry client: %s\n", err)
		return
	}
}

func GenerateUuid() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return ErrUuidGeneration
	}

	ClientInstance.Uuid = id.String()
	return nil
}
