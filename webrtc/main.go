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

package webrtc

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/adsisto/adsisto/webrtc/gst"
	"github.com/pion/webrtc"
	"math/rand"
)

type Config struct {
	StunServer string
	Source     string
	connection *webrtc.PeerConnection
	track      *webrtc.Track
}

func (c *Config) StartConnection() error {
	stunUrl := fmt.Sprintf("stun:%s", c.StunServer)

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{stunUrl},
			},
		},
	}
	connection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return err
	}

	c.connection = connection
	return nil
}

func (c *Config) StreamStart() (string, error) {
	track, err := c.connection.NewTrack(
		webrtc.DefaultPayloadTypeVP8, rand.Uint32(), "video", "pion1",
	)
	if err != nil {
		return "", err
	}
	c.track = track

	_, err = c.connection.AddTrack(track)
	if err != nil {
		return "", err
	}

	offer, err := c.connection.CreateOffer(nil)
	if err != nil {
		return "", err
	}

	err = c.connection.SetLocalDescription(offer)
	if err != nil {
		return "", err
	}

	encoded, err := json.Marshal(offer)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encoded), nil
}

func (c *Config) RemoteCallback(answer webrtc.SessionDescription) error {
	err := c.connection.SetRemoteDescription(answer)
	if err != nil {
		return err
	}

	gst.CreatePipeline(webrtc.VP8, []*webrtc.Track{c.track}, c.Source).Start()

	return nil
}
