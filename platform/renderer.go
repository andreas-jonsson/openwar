/* Any copyright is dedicated to the Public Domain.
 * http://creativecommons.org/publicdomain/zero/1.0/ */

package platform

import (
	"image"
	"image/draw"
)

type Renderer interface {
	Clear()
	Present()
	Shutdown()
	BackBuffer() draw.Image
	Blit(src image.Image, sp image.Point)
}
