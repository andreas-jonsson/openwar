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

import (
	"image"
	"image/color"
	"image/draw"
	"log"
	"strconv"

	"github.com/gopherjs/gopherjs/js"
)

type jsRenderer struct {
	backBuffer       *image.RGBA
	canvas, document *js.Object

	width, height,
	logicalHeight int
}

func NewRenderer(w, h int, data ...interface{}) (Renderer, error) {
	r := jsRenderer{
		width:         320,
		height:        200,
		logicalHeight: 240,
		document:      js.Global.Get("document"),
	}

	for i := 0; i < len(data); i++ {
		handled := true
		p := data[i]

		ps, ok := p.(string)
		if ok {
			switch ps {
			case "fullscreen":
			case "widescreen":
				r.width = 640
				r.height = 300
				r.logicalHeight = 360
			case "title":
				i++
				r.SetWindowTitle(data[i].(string))
			default:
				handled = false
			}
		}

		if !handled {
			log.Println("invalid parameter passed to renderer:", p)
		}
	}

	r.backBuffer = image.NewRGBA(image.Rect(0, 0, r.width, r.height))

	canvas := r.document.Call("createElement", "canvas")
	canvas.Call("setAttribute", "width", strconv.Itoa(r.width))
	canvas.Call("setAttribute", "height", strconv.Itoa(r.height))
	canvas.Get("style").Set("width", strconv.Itoa(w)+"px")
	canvas.Get("style").Set("height", strconv.Itoa(h)+"px")
	r.document.Get("body").Call("appendChild", canvas)
	r.canvas = canvas

	return &r, nil
}

func (p *jsRenderer) ToggleFullscreen() {
}

func (p *jsRenderer) Clear() {
	draw.Draw(p.backBuffer, p.backBuffer.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)
}

func (p *jsRenderer) Present() {
	ctx := p.canvas.Call("getContext", "2d")
	img := ctx.Call("getImageData", 0, 0, p.width, p.height)
	data := img.Get("data")

	arrBuf := js.Global.Get("ArrayBuffer").New(data.Length())
	buf8 := js.Global.Get("Uint8ClampedArray").New(arrBuf)
	buf32 := js.Global.Get("Uint32Array").New(arrBuf)

	buf := buf32.Interface().([]uint)

	for y := 0; y < p.height; y++ {
		for x := 0; x < p.width; x++ {
			r, g, b, _ := p.backBuffer.At(x, y).RGBA()
			buf[y*p.width+x] = 0x00000000 | (uint(b) << 16) | (uint(g) << 8) | uint(r)
			//buf[y*p.width+x] = 0xFF0000FF
		}
	}

	data.Call("set", buf8)
	ctx.Call("putImageData", img, 0, 0)
}

func (p *jsRenderer) Shutdown() {
}

func (p *jsRenderer) SetWindowTitle(title string) {
	js.Global.Get("document").Set("title", title)
}

func (p *jsRenderer) BackBuffer() draw.Image {
	return p.backBuffer
}

func (p *jsRenderer) BlitPaletted(dp image.Point, src *image.Paletted) {
	p.BlitImage(dp, src, src.Palette)
}

func (p *jsRenderer) BlitImage(dp image.Point, src *image.Paletted, pal color.Palette) {
	p.Blit(dp, src, src.Bounds(), pal)
}

func (p *jsRenderer) Blit(dp image.Point, src *image.Paletted, sr image.Rectangle, pal color.Palette) {
	min := sr.Min
	max := sr.Max

	for y, dy := min.Y, 0; y < max.Y; y++ {
		for x, dx := min.X, 0; x < max.X; x++ {
			c := pal[src.ColorIndexAt(x, y)]
			if _, _, _, a := c.RGBA(); a > 0 {
				p.backBuffer.Set(dx+dp.X, dy+dp.Y, c)
			}
			dx++
		}
		dy++
	}
}
