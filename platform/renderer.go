/*
Copyright (C) 2016-2017 Andreas T Jonsson

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
)

type Renderer interface {
	Clear()
	Present()
	Shutdown()
	ToggleFullscreen()
	SetWindowTitle(title string)
	BackBuffer() draw.Image

	DrawRect(dest image.Rectangle, c color.Color, fill bool)

	BlitPaletted(dp image.Point, src *image.Paletted)
	BlitImage(dp image.Point, src *image.Paletted, pal color.Palette)
	Blit(dp image.Point, src *image.Paletted, sr image.Rectangle, pal color.Palette)
}
