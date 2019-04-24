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
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/gpio/gpioreg"
	"periph.io/x/periph/host"
	"time"
)

type Config struct {
	InputPins  []string
	OutputPins []string
}

var (
	ErrPinUndefined = errors.New("GPIO pin not defined in configuration")
)

func (c *Config) SetupPins() []error {
	var ers []error

	if _, err := host.Init(); err != nil {
		ers = append(ers, err)
		return ers
	}

	for _, pin := range c.OutputPins {
		p := gpioreg.ByName(pin)
		err := p.Out(gpio.Low)
		if err != nil {
			ers = append(ers, err)
		}
	}
	return ers
}

func (c *Config) TogglePin(pin string) error {
	found := false

	for _, p := range c.OutputPins {
		if pin == p {
			found = true
			break
		}
	}

	if !found {
		return ErrPinUndefined
	}

	if _, err := host.Init(); err != nil {
		return err
	}

	p := gpioreg.ByName(pin)
	err := p.Out(gpio.High)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 2)
	err = p.Out(gpio.Low)
	if err != nil {
		return err
	}

	return nil
}
