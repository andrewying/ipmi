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
	"crypto/tls"
	"github.com/certifi/gocertifi"
	"net/http"
)

var (
	// ServerURL is the URL to the remote telemetry server
	ServerURL string
	// ApiKey is the public API key for accessing the telemetry server
	ApiKey string
	// ApiClient is the HTTP client instance
	ApiClient *http.Client
)

// SetupApiClient sets up the HTTP client instance
func SetupApiClient() error {
	certPool, err := gocertifi.CACerts()
	if err != nil {
		return err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: certPool,
		},
	}
	ApiClient = &http.Client{
		Transport: transport,
	}

	return nil
}
