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

package console

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
	PromptPattern           = regexp.MustCompile(`\S+ login:\S*`)
	AuthenticatedPattern    = regexp.MustCompile(`\S+@\S+:\S+\$`)
	ErrConsoleNotResponding = errors.New("console not responding")
	ErrLoginInvalid         = errors.New("username and/or password is incorrect")
)

// Create a new serial console session
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

// Checks that the console session is authenticated
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

// Close serial console session
func (c *SerialConsole) Close() error {
	return c.connection.Close()
}
