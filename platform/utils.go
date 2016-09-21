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
	"path"
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

func drawRect(backBuffer *image.RGBA, dest image.Rectangle, c color.Color) {
	//TODO Optimize this!
	draw.DrawMask(backBuffer, backBuffer.Bounds(), &image.Uniform{c}, image.ZP, (*noneSolidRect)(&dest), image.ZP, draw.Over)
}

func blitPaletted(backBuffer *image.RGBA, dp image.Point, src *image.Paletted) {
	blit(backBuffer, dp, src, src.Bounds(), src.Palette)
}

func blitImage(backBuffer *image.RGBA, dp image.Point, src *image.Paletted, pal color.Palette) {
	blit(backBuffer, dp, src, src.Bounds(), pal)
}

func blit(backBuffer *image.RGBA, dp image.Point, src *image.Paletted, sr image.Rectangle, pal color.Palette) {
	bbMaxBounds := backBuffer.Bounds().Max

	min := sr.Min
	max := sr.Max

	srWidth := src.Bounds().Max.X
	sPix := src.Pix
	dPix := backBuffer.Pix

	//TODO Optimize this code!

	for y, dy := min.Y, 0; y < max.Y; y++ {
		if dy+dp.Y < 0 || dy+dp.Y >= bbMaxBounds.Y {
			continue
		}

		for x, dx := min.X, 0; x < max.X; x++ {
			if dx+dp.X < 0 || dx+dp.X >= bbMaxBounds.X {
				continue
			}

			i := sPix[y*srWidth+x]
			c := pal[i]

			if r, g, b, a := c.RGBA(); a > 0 {
				offset := (dy+dp.Y)*bbMaxBounds.X*4 + (dx+dp.X)*4

				dPix[offset] = byte(r)
				dPix[offset+1] = byte(g)
				dPix[offset+2] = byte(b)
				//dPix[offset+3] = 0xFF
			}
			dx++
		}
		dy++
	}
}

var DataPath string

func RootJoin(p ...string) string {
	return path.Join(DataPath, path.Join(p...))
}
