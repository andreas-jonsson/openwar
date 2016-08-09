/*
Copyright (C) 2016 Andreas T Jonsson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package platform

const (
	KeyUnknown = iota
	KeyUp
	KeyDown
	KeyLeft
	KeyRight
	KeyEsc
	KeySpace
)

const (
	KeyModNone   = 0x0000
	KeyModLshift = 0x0001
	KeyModRshift = 0x0002
	KeyModLctrl  = 0x0040
	KeyModRctrl  = 0x0080
	KeyModLalt   = 0x0100
	KeyModRalt   = 0x0200
	KeyModLgui   = 0x0400
	KeyModRgui   = 0x0400
	KeyModCaps   = 0x2000
)

const (
	KeyModAlt   = KeyModLalt | KeyModRalt
	KeyModGui   = KeyModLgui | KeyModRgui
	KeyModCtrl  = KeyModLctrl | KeyModRctrl
	KeyModShift = KeyModLshift | KeyModRshift
)

const (
	MouseButtonDown = iota
	MouseButtonUp
	MouseMotion
	MouseWheel
)

type (
	Event     interface{}
	QuitEvent struct{}

	KeyUpEvent struct {
		Rune     rune
		Key, Mod int
	}

	KeyDownEvent struct {
		Rune     rune
		Key, Mod int
	}

	MouseMotionEvent struct {
		X, Y, XRel, YRel int
	}

	MouseButtonEvent struct {
		X, Y, Button, Type int
	}
)
