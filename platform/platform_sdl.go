/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package platform

import (
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

var keyMapping = map[sdl.Keycode]int{
	sdl.K_UP:    KEY_UP,
	sdl.K_DOWN:  KEY_DOWN,
	sdl.K_LEFT:  KEY_LEFT,
	sdl.K_RIGHT: KEY_RIGHT,
}

func Init() error {
	runtime.LockOSThread()
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
		if key, ok := keyMapping[t.Keysym.Sym]; ok {
			ev.Key = key
			return ev
		}
	case *sdl.KeyDownEvent:
		ev := &KeyDownEvent{}
		if key, ok := keyMapping[t.Keysym.Sym]; ok {
			ev.Key = key
			return ev
		}
	}

	return nil
}
