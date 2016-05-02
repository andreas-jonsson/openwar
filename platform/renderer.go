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
	BackBuffer() draw.Image
	Blit(src image.Image, sp image.Point)
	BlitPal(src *image.Paletted, pal color.Palette, sp image.Point)
}
