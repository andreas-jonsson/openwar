// +build mobile

/*
Copyright (C) 2016-2018 Andreas T Jonsson

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
)

var (
	ExternalBackBuffer *image.RGBA
	PaintEventChan     = make(chan chan struct{})
)

type mobileRenderer struct {
	backBuffer *image.RGBA
	width, height,
	logicalHeight int
}

func NewRenderer(w, h int, data ...interface{}) (Renderer, error) {
	r := mobileRenderer{
		width:         320,
		height:        200,
		logicalHeight: 240,
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
	return &r, nil
}

func (p *mobileRenderer) ToggleFullscreen() {
}

func (p *mobileRenderer) Clear() {
	draw.Draw(p.backBuffer, p.backBuffer.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 255}}, image.ZP, draw.Src)
}

func (p *mobileRenderer) Present() {
	if cb, ok := <-PaintEventChan; ok {
		draw.Draw(ExternalBackBuffer, p.backBuffer.Bounds(), p.backBuffer, image.ZP, draw.Src)
		cb <- struct{}{}
	}
}

func (p *mobileRenderer) Shutdown() {
}

func (p *mobileRenderer) SetWindowTitle(title string) {
}

func (p *mobileRenderer) BackBuffer() draw.Image {
	return p.backBuffer
}

func (p *mobileRenderer) DrawRect(dest image.Rectangle, c color.Color, fill bool) {
	drawRect(p.backBuffer, dest, c, fill)
}

func (p *mobileRenderer) BlitPaletted(dp image.Point, src *image.Paletted) {
	blit(p.backBuffer, dp, src, src.Bounds(), src.Palette)
}

func (p *mobileRenderer) BlitImage(dp image.Point, src *image.Paletted, pal color.Palette) {
	blit(p.backBuffer, dp, src, src.Bounds(), pal)
}

func (p *mobileRenderer) Blit(dp image.Point, src *image.Paletted, sr image.Rectangle, pal color.Palette) {
	blit(p.backBuffer, dp, src, sr, pal)
}
