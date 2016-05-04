/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

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
