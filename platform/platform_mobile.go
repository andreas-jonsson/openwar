// +build mobile

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
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/touch"
)

const maxEvents = 128

var (
	InputEventChan = make(chan interface{}, maxEvents)
	sizeEvent      size.Event
)

func Init() error {
	idCounter = 0
	return nil
}

func Shutdown() {
}

func Mouse() MouseState {
	return MouseState{}
}

func PollEvent() Event {
	select {
	case ev, ok := <-InputEventChan:
		if ok {
			switch e := ev.(type) {
			case size.Event:
				sizeEvent = e
			case touch.Event:
				ws := 320 / float32(sizeEvent.WidthPx)
				hs := 200 / float32(sizeEvent.HeightPx)

				if e.Type == touch.TypeBegin {
					return &MouseButtonEvent{X: int(e.X * ws), Y: int(e.Y * hs), Button: 0, Type: MouseButtonDown}
				} else if e.Type == touch.TypeEnd {
					return &MouseButtonEvent{X: int(e.X * ws), Y: int(e.Y * hs), Button: 0, Type: MouseButtonUp}
				} else {
					return &MouseMotionEvent{X: int(e.X * ws), Y: int(e.Y * hs)}
				}
			}
		} else {
			return QuitEvent{}
		}
	default:
	}
	return nil
}
