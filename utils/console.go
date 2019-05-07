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
	"errors"
	"fmt"
	"github.com/tarm/serial"
	"regexp"
	"time"
)

type SerialConsole struct {
	Device        string
	Authenticated bool
	config        *serial.Config
	connection    *serial.Port
}

var (
	PromptPattern             = regexp.MustCompile(`\S+ login:\S*`)
	AuthenticatedPattern      = regexp.MustCompile(`\S+@\S+:\S+\$`)
	ErrConsoleNotResponding   = errors.New("console not responding")
	ErrConsoleUnauthenticated = errors.New("console not authenticated")
	ErrLoginInvalid           = errors.New("username and/or password is incorrect")
)

// NewConsole creates a new serial console session
func NewConsole(device string) (*SerialConsole, error) {
	config := &serial.Config{
		Name:        device,
		Baud:        115200,
		ReadTimeout: time.Second * 5,
	}
	session, err := serial.OpenPort(config)
	if err != nil {
		return nil, err
	}

	return &SerialConsole{
		Device:        device,
		Authenticated: false,
		config:        config,
		connection:    session,
	}, nil
}

// Authenticate checks that the console session is authenticated
func (c *SerialConsole) Authenticate(username string, password string) error {
	var err error

	buf := make([]byte, 128)
	n, _ := c.connection.Read(buf)

	// If no prompts are available, write EOL to get a prompt
	for i := 0; n == 0; i++ {
		if i == 5 {
			return ErrConsoleNotResponding
		}

		_ = c.connection.Flush()
		_, err = c.connection.Write([]byte("\\r\\n"))

		if err != nil {
			return err
		}

		n, _ = c.connection.Read(buf)
	}

	if !PromptPattern.Match(buf) && !AuthenticatedPattern.Match(buf) {
		_ = c.connection.Flush()
		_, err = c.connection.Write([]byte("\\r\\n"))

		if err != nil {
			return err
		}

		n, _ = c.connection.Read(buf)
	}

	// Login is requested by console
	for j := 0; PromptPattern.Match(buf); j++ {
		if j == 5 {
			return ErrLoginInvalid
		}

		c.Authenticated = false
		time.Sleep(time.Second * 10)

		_ = c.connection.Flush()
		_, err = c.connection.Write(
			[]byte(fmt.Sprintf("%s\\r\\n", username)),
		)

		if err != nil {
			return err
		}

		_ = c.connection.Flush()
		_, err = c.connection.Write(
			[]byte(fmt.Sprintf("%s\\r\\n", password)),
		)

		if err != nil {
			return err
		}

		n, _ = c.connection.Read(buf)
	}

	// Session is already authenticated
	if AuthenticatedPattern.Match(buf) {
		c.Authenticated = true
	}

	return nil
}

// Write writes to serial console
func (c *SerialConsole) Write(s string) error {
	if !c.Authenticated {
		return ErrConsoleUnauthenticated
	}

	_, err := c.connection.Write([]byte(s))
	if err != nil {
		return err
	}

	return nil
}

// Close closes the serial console session
func (c *SerialConsole) Close() error {
	return c.connection.Close()
}
