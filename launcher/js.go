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

package launcher

import (
	"bytes"

	"github.com/andreas-jonsson/openwar/game"
	"github.com/andreas-jonsson/openwar/resource"
	"github.com/gopherjs/gopherjs/js"
)

func Start() {
	js.Global.Call("addEventListener", "load", func() { go load() })
}

func load() {
	cfg := &game.Config{
		Fullscreen: false,
		Widescreen: false,
		WC2Input:   true,
	}

	_, data := openFile()
	readSeeker := bytes.NewReader(data)

	if war, err := resource.OpenArchiveFrom(readSeeker, int64(len(data))); err == nil {
		game.Start(cfg, war)
	} else {
		js.Global.Call("alert", err.Error())
	}
}

func openFile() (string, []byte) {
	document := js.Global.Get("document")

	inputElem := document.Call("createElement", "input")
	inputElem.Call("setAttribute", "type", "file")
	inputElem.Call("setAttribute", "accept", ".war")

	document.Get("body").Call("appendChild", inputElem)

	filec := make(chan *js.Object, 1)
	inputElem.Set("onchange", func(event *js.Object) {
		filec <- inputElem.Get("files").Index(0)
	})

	file := <-filec
	name := file.Get("name").String()
	reader := js.Global.Get("FileReader").New()

	bufc := make(chan []byte, 1)
	reader.Set("onloadend", func(event *js.Object) {
		bufc <- js.Global.Get("Uint8Array").New(reader.Get("result")).Interface().([]byte)
	})
	reader.Call("readAsArrayBuffer", file)
	data := <-bufc

	div := document.Call("getElementById", "upload_text")
	document.Get("body").Call("removeChild", inputElem)
	document.Get("body").Call("removeChild", div)
	return name, data
}
