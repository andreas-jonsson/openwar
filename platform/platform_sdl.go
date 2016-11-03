// +build !js,!mobile

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

import (
	"os"
	"os/user"
	"path"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	ScreenScale4x3  = 200.0 / 240.0
	ScreenScale16x9 = 200.0 / 234.0
)

var ScreenScale = ScreenScale4x3

var keyMapping = map[sdl.Keycode]int{
	sdl.K_UP:     KeyUp,
	sdl.K_DOWN:   KeyDown,
	sdl.K_LEFT:   KeyLeft,
	sdl.K_RIGHT:  KeyRight,
	sdl.K_ESCAPE: KeyEsc,
	sdl.K_SPACE:  KeySpace,
}

var keyModMapping = map[uint16]int{
	sdl.KMOD_NONE: KeyModNone,

	sdl.KMOD_CAPS:  KeyModCaps,
	sdl.KMOD_ALT:   KeyModAlt,
	sdl.KMOD_CTRL:  KeyModCtrl,
	sdl.KMOD_SHIFT: KeyModShift,
	sdl.KMOD_GUI:   KeyModGui,

	sdl.KMOD_RALT:   KeyModRalt,
	sdl.KMOD_RCTRL:  KeyModRctrl,
	sdl.KMOD_RSHIFT: KeyModRshift,
	sdl.KMOD_RGUI:   KeyModRgui,

	sdl.KMOD_LALT:   KeyModLalt,
	sdl.KMOD_LCTRL:  KeyModLctrl,
	sdl.KMOD_LSHIFT: KeyModLshift,
	sdl.KMOD_LGUI:   KeyModLgui,
}

var mouseMapping = map[int]int{
	sdl.MOUSEBUTTONDOWN: MouseButtonDown,
	sdl.MOUSEBUTTONUP:   MouseButtonUp,
	sdl.MOUSEWHEEL:      MouseWheel,
}

func init() {
	if runtime.GOOS == "windows" {
		ConfigPath = path.Join(os.Getenv("LOCALAPPDATA"), "OpenWar")
	} else {
		if usr, err := user.Current(); err == nil {
			ConfigPath = path.Join(usr.HomeDir, ".openwar")
		}
	}
	os.MkdirAll(ConfigPath, 0755)
}

func Init() error {
	idCounter = 0
	// Looks like we can't reinitialize SDL after all. :(
	if sdl.WasInit(0) == 0 {
		return sdl.Init(sdl.INIT_EVERYTHING)
	}
	return nil
}

func Shutdown() {
	//sdl.Quit()
}

func Mouse() MouseState {
	x, y, buttons := sdl.GetMouseState()

	window := sdl.GetMouseFocus()
	if window == nil {
		return MouseState{}
	}

	w, h := window.GetSize()
	y = int(200.0 * (float64(y) / float64(h)))
	xw := (float64(x) / float64(w))

	if ScreenScale == ScreenScale4x3 {
		x = int(320.0 * xw)
	} else {
		x = int(416.0 * xw)
	}

	left := (buttons & sdl.ButtonLMask()) != 0
	middle := (buttons & sdl.ButtonMMask()) != 0
	right := (buttons & sdl.ButtonRMask()) != 0

	return MouseState{X: x, Y: y, Buttons: [3]bool{left, middle, right}}
}

func PollEvent() Event {
	event := sdl.PollEvent()
	if event == nil {
		return nil
	}

	switch t := event.(type) {
	case *sdl.QuitEvent:
		return &QuitEvent{}
	case *sdl.KeyUpEvent:
		ev := &KeyUpEvent{}
		if mod, ok := keyModMapping[t.Keysym.Mod]; ok {
			ev.Mod = mod
		} else {
			ev.Mod = KeyModNone
		}

		if key, ok := keyMapping[t.Keysym.Sym]; ok {
			ev.Key = key
			ev.Rune = rune(t.Keysym.Unicode)
		} else {
			ev.Key = KeyUnknown
		}
		return ev
	case *sdl.KeyDownEvent:
		ev := &KeyDownEvent{}
		if mod, ok := keyModMapping[t.Keysym.Mod]; ok {
			ev.Mod = mod
		} else {
			ev.Mod = KeyModNone
		}

		if key, ok := keyMapping[t.Keysym.Sym]; ok {
			ev.Key = key
			ev.Rune = rune(t.Keysym.Unicode)
		} else {
			ev.Key = KeyUnknown
		}
		return ev
	case *sdl.MouseButtonEvent:
		ev := &MouseButtonEvent{}
		ev.Button = int(t.Button)
		ev.X = int(t.X)
		ev.Y = int(float64(t.Y) * ScreenScale)

		switch t.Type {
		case sdl.MOUSEBUTTONDOWN:
			ev.Type = MouseButtonDown
		case sdl.MOUSEBUTTONUP:
			ev.Type = MouseButtonUp
		case sdl.MOUSEWHEEL:
			ev.Type = MouseWheel
		}
		return ev
	case *sdl.MouseMotionEvent:
		ev := &MouseMotionEvent{}
		ev.X = int(t.X)
		ev.Y = int(float64(t.Y) * ScreenScale)
		ev.XRel = int(t.XRel)
		ev.YRel = int(float64(t.YRel) * ScreenScale)
		return ev
	}

	return nil
}
