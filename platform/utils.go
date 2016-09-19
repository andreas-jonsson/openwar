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
)

func BlitPaletted(backBuffer *image.RGBA, dp image.Point, src *image.Paletted) {
	Blit(backBuffer, dp, src, src.Bounds(), src.Palette)
}

func BlitImage(backBuffer *image.RGBA, dp image.Point, src *image.Paletted, pal color.Palette) {
	Blit(backBuffer, dp, src, src.Bounds(), pal)
}

func Blit(backBuffer *image.RGBA, dp image.Point, src *image.Paletted, sr image.Rectangle, pal color.Palette) {
	bbMaxBounds := backBuffer.Bounds().Max

	min := sr.Min
	max := sr.Max

	if max.X > bbMaxBounds.X {
		max.X = bbMaxBounds.X
	}

	if max.Y > bbMaxBounds.Y {
		max.Y = bbMaxBounds.Y
	}

	srWidth := src.Bounds().Max.X
	sPix := src.Pix
	dPix := backBuffer.Pix

	for y, dy := min.Y, 0; y < max.Y; y++ {
		for x, dx := min.X, 0; x < max.X; x++ {
			i := sPix[y*srWidth+x]
			c := pal[i]

			if r, g, b, a := c.RGBA(); a > 0 {
				offset := (dy+dp.Y)*bbMaxBounds.X*4 + (dx+dp.X)*4

				dPix[offset] = byte(r)
				dPix[offset+1] = byte(g)
				dPix[offset+2] = byte(b)
				dPix[offset+3] = 0xFF
			}
			dx++
		}
		dy++
	}
}
