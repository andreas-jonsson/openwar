/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

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

	BlitPaletted(dp image.Point, src *image.Paletted)
	BlitImage(dp image.Point, src *image.Paletted, pal color.Palette)
	Blit(dp image.Point, src *image.Paletted, sr image.Rectangle, pal color.Palette)
}
