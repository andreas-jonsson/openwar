// +build js

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

import "github.com/gopherjs/gopherjs/js"

var (
	eventQueue = make(chan Event)

	mouseXPrev, mouseYPrev,
	mouseXRel, mouseYRel,
	mouseX, mouseY int
)

func setupCanvasInput(canvas *js.Object, w, h, resX, resY int) {
	rect := canvas.Call("getBoundingClientRect")
	left := rect.Get("left").Float()
	top := rect.Get("top").Float()

	canvas.Call("addEventListener", "mousemove", func(event *js.Object) {
		mouseX = int(((event.Get("clientX").Float() - left) / float64(w)) * float64(resX))
		mouseY = int(((event.Get("clientY").Float() - top) / float64(h)) * float64(resY))

		mouseXRel += mouseX - mouseXPrev
		mouseYRel += mouseY - mouseYPrev

		mouseXPrev = mouseX
		mouseYPrev = mouseY
	})

	canvas.Call("addEventListener", "mouseup", func(event *js.Object) {
		go func() {
			eventQueue <- &MouseButtonEvent{
				X:      mouseX,
				Y:      mouseY,
				Button: event.Get("button").Int(),
				Type:   MouseButtonUp,
			}
		}()
	})

	canvas.Call("addEventListener", "mousedown", func(event *js.Object) {
		go func() {
			eventQueue <- &MouseButtonEvent{
				X:      mouseX,
				Y:      mouseY,
				Button: event.Get("button").Int(),
				Type:   MouseButtonDown,
			}
		}()
	})
}

func Init() error {
	return nil
}

func Shutdown() {
}

func PollEvent() Event {
	if mouseXRel != 0 || mouseYRel != 0 {
		ev := &MouseMotionEvent{
			X:    mouseX,
			Y:    mouseY,
			XRel: mouseXRel,
			YRel: mouseYRel,
		}

		mouseXRel = 0
		mouseYRel = 0
		return ev
	}

	select {
	case ev := <-eventQueue:
		return ev
	default:
		return nil
	}
}
