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
	"strings"
)

// Maps the physical keys with their corresponding byte values
var KeyMap = map[string]string{
	"A":          "0x04",
	"B":          "0x05",
	"C":          "0x06",
	"D":          "0x07",
	"E":          "0x08",
	"F":          "0x09",
	"G":          "0x0a",
	"H":          "0x0b",
	"I":          "0x0c",
	"J":          "0x0d",
	"K":          "0x0e",
	"L":          "0x0f",
	"M":          "0x10",
	"N":          "0x11",
	"O":          "0x12",
	"P":          "0x13",
	"Q":          "0x14",
	"R":          "0x15",
	"S":          "0x16",
	"T":          "0x17",
	"U":          "0x18",
	"V":          "0x19",
	"W":          "0x1a",
	"X":          "0x1b",
	"Y":          "0x1c",
	"Z":          "0x1d",
	"1":          "0x1e",
	"2":          "0x1f",
	"3":          "0x20",
	"4":          "0x21",
	"5":          "0x22",
	"6":          "0x23",
	"7":          "0x24",
	"8":          "0x25",
	"9":          "0x26",
	"0":          "0x27",
	"ENTER":      "0x28",
	"ESC":        "0x29",
	"BACKSPACE":  "0x2a",
	"TAB":        "0x2b",
	"SPACE":      "0x2c",
	"MINUS":      "0x2d",
	"EQUAL":      "0x2e",
	"LEFTBRACE":  "0x2f",
	"RIGHTBRACE": "0x30",
	"BACKSLASH":  "0x31",
	"HASHTILDE":  "0x32",
	"SEMICOLON":  "0x33",
	"APOSTROPHE": "0x34",
	"GRAVE":      "0x35",
	"COMMA":      "0x36",
	"DOT":        "0x37",
	"SLASH":      "0x38",
	"CAPSLOCK":   "0x39",
	"F1":         "0x3a",
	"F2":         "0x3b",
	"F3":         "0x3c",
	"F4":         "0x3d",
	"F5":         "0x3e",
	"F6":         "0x3f",
	"F7":         "0x40",
	"F8":         "0x41",
	"F9":         "0x42",
	"F10":        "0x43",
	"F11":        "0x44",
	"F12":        "0x45",
	"SYSRQ":      "0x46",
	"SCROLLLOCK": "0x47",
	"PAUSE":      "0x48",
	"INSERT":     "0x49",
	"HOME":       "0x4a",
	"PAGEUP":     "0x4b",
	"DELETE":     "0x4c",
	"END":        "0x4d",
	"PAGEDOWN":   "0x4e",
	"RIGHT":      "0x4f",
	"LEFT":       "0x50",
	"DOWN":       "0x51",
	"UP":         "0x52",
	"NUMLOCK":    "0x53",
	"KPSLASH":    "0x54",
	"KPASTERISK": "0x55",
	"KPMINUS":    "0x56",
	"KPPLUS":     "0x57",
	"KPENTER":    "0x58",
	"KPDOT":      "0x63",
	"102ND":      "0x64",
	"COMPOSE":    "0x65",
	"POWER":      "0x66",
	"KPEQUAL":    "0x67",
	"F13":        "0x68",
	"F14":        "0x69",
	"F15":        "0x6a",
	"F16":        "0x6b",
	"F17":        "0x6c",
	"F18":        "0x6d",
	"F19":        "0x6e",
	"F20":        "0x6f",
	"F21":        "0x70",
	"F22":        "0x71",
	"F23":        "0x72",
	"F24":        "0x73",
	"OPEN":       "0x74",
	"HELP":       "0x75",
	"PROPS":      "0x76",
	"FRONT":      "0x77",
	"STOP":       "0x78",
	"AGAIN":      "0x79",
	"UNDO":       "0x7a",
	"CUT":        "0x7b",
	"COPY":       "0x7c",
	"PASTE":      "0x7d",
	"FIND":       "0x7e",
	"MUTE":       "0x7f",
	"VOLUMEUP":   "0x80",
	"VOLUMEDOWN": "0x81",
	"KPCOMMA":    "0x85",
	"CTRL":       "0xe0",
	"SHIFT":      "0xe1",
	"ALT":        "0xe2",
}

var AliasMap = map[string]string{
	// Obsolete mapping
	"DECIMAL":   "KPDOT",
	"SUBTRACT":  "KPMINUS",
	"MULTIPLY":  "KPASTERISK",
	"ADD":       "KPPLUS",
	"DIVIDE":    "KPSLASH",
	"SEPARATOR": "KPCOMMA",
	// Obsolete mapping ends
	"+":          "PLUS",
	".":          "DOT",
	"-":          "MINUS",
	"=":          "EQUAL",
	" ":          "SPACE",
	"[":          "LEFTBRACE",
	"]":          "RIGHTBRACE",
	"\\":         "BACKSLASH",
	"#":          "HASHTILDE",
	";":          "SEMICOLON",
	"'":          "APOSTROPHE",
	"`":          "GRAVE",
	",":          "COMMA",
	"/":          "SLASH",
	"ARROWDOWN":  "DOWN",
	"ARROWLEFT":  "LEFT",
	"ARROWRIGHT": "RIGHT",
	"ARROWUP":    "UP",
}

func (m *StreamMessage) ParseMessage() {
	m.Key = strings.ToUpper(m.Key)

	if value, found := KeyMap[m.Key]; found {
		m.Key = value
		return
	}

	if alias, found := AliasMap[m.Key]; found {
		m.Key = KeyMap[alias]
		return
	}

	m.Key = ""
}

func (m *StreamMessage) GenerateHID() [8]byte {
	var array [8]byte

	if m.Key == "" {
		return array
	}

	bytes, err := hex.DecodeString(m.Key)
	if err != nil {
		return array
	}

	array[2] = bytes[0]

	switch {
	case m.Ctrl:
		array[0] = 0x01
	case m.Shift:
		array[0] = 0x02
	case m.Alt:
		array[0] = 0x04
	case m.Meta:
		array[0] = 0x08
	default:
		array[0] = 0x00
	}

	return array
}
