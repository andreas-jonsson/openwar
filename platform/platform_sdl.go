/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package platform

import (
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

const screenScale = 200.0 / 240.0

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
	sdl.MOUSEMOTION:     MouseMotion,
	sdl.MOUSEWHEEL:      MouseWheel,
}

func init() {
	runtime.LockOSThread()
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
		ev.Y = int(float32(t.Y) * screenScale)
		ev.Type = int(t.Type)
		return ev
	case *sdl.MouseMotionEvent:
		ev := &MouseMotionEvent{}
		ev.X = int(t.X)
		ev.Y = int(float32(t.Y) * screenScale)
		ev.XRel = int(t.XRel)
		ev.YRel = int(float32(t.YRel) * screenScale)
		return ev
	}

	return nil
}
