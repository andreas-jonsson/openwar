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
	"log"
	"os"
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
	isInvalid := func() bool {
		_, err := os.Stat(path.Join(DataPath, "icon.png"))
		return os.IsNotExist(err)
	}

	wd, _ := os.Getwd()
	DataPath = path.Join(wd, "data")

	if isInvalid() {
		DataPath = "/usr/local/share/openwar"
		if runtime.GOOS == "darwin" {
			if isInvalid() {
				DataPath = path.Join(sdl.GetBasePath(), "data")
			}
		}
	}

	if isInvalid() {
		log.Panicln("Could not locate data folder!")
	}

	log.Println("Data path:", DataPath)
}

func Init() error {
	return sdl.Init(sdl.INIT_EVERYTHING)
}

func Shutdown() {
	sdl.Quit()
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
		ev.Type = int(t.Type)
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
