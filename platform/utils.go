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
	"math"
	"path"
	"sync/atomic"
)

type noneSolidRect image.Rectangle

func (r *noneSolidRect) ColorModel() color.Model {
	return color.AlphaModel
}

func (r *noneSolidRect) Bounds() image.Rectangle {
	return *(*image.Rectangle)(r)
}

func (r *noneSolidRect) At(x, y int) color.Color {
	if x == r.Max.X-1 || x == r.Min.X || y == r.Max.Y-1 || y == r.Min.Y {
		return color.Alpha{0xFF}
	}
	return color.Alpha{0}
}

func drawRect(backBuffer *image.RGBA, dest image.Rectangle, c color.Color, fill bool) {
	if fill {
		draw.Draw(backBuffer, dest, &image.Uniform{c}, image.ZP, draw.Over)
	} else {
		//TODO Optimize this!
		draw.DrawMask(backBuffer, backBuffer.Bounds(), &image.Uniform{c}, image.ZP, (*noneSolidRect)(&dest), image.ZP, draw.Over)
	}
}

func blitPaletted(backBuffer *image.RGBA, dp image.Point, src *image.Paletted) {
	blit(backBuffer, dp, src, src.Bounds(), src.Palette)
}

func blitImage(backBuffer *image.RGBA, dp image.Point, src *image.Paletted, pal color.Palette) {
	blit(backBuffer, dp, src, src.Bounds(), pal)
}

func blit(backBuffer *image.RGBA, dp image.Point, src *image.Paletted, sr image.Rectangle, pal color.Palette) {
	srcImage := *src
	srcImage.Palette = pal

	r := image.Rectangle{dp, dp.Add(sr.Size())}
	draw.Draw(backBuffer, r, &srcImage, sr.Min, draw.Over)
}

var ConfigPath string

func CfgRootJoin(p ...string) string {
	return path.Join(ConfigPath, path.Join(p...))
}

var idCounter uint64

func NewId64() uint64 {
	if idCounter == math.MaxUint64 {
		log.Panicln("id space exhausted")
	}
	return atomic.AddUint64(&idCounter, 1) - 1
}
